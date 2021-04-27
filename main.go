package main

import (
	"crypto/rand"
	"flag"
	"math/big"
	"time"

	"github.com/andlabs/ui"
)

var (
	max = flag.Int64("max", 10, "")
)

func init() {
	flag.Parse()
}

func main() {
	ui.Main(setupUI)
}

func setupUI() {
	win := ui.NewWindow("trainer", 0, 0, true)
	win.OnClosing(func(w *ui.Window) bool {
		ui.Quit()
		return true
	})
	win.SetMargined(true)

	ui.OnShouldQuit(func() bool {
		win.Destroy()
		return true
	})

	vb := ui.NewVerticalBox()

	params := ui.NewGroup("Parameters")
	form := ui.NewForm()

	dur := formSlider(form, "One show time(ms)", 5000)
	num := formSlider(form, "Shows number", 50)

	params.SetChild(form)

	vb.Append(params, true)

	hb := ui.NewHorizontalBox()
	number := ui.NewGroup("Number")
	ovb := ui.NewVerticalBox()

	outHand := &areaHandler{text: ""}
	outArea := ui.NewArea(outHand)
	outArea.Disable()

	ovb.Append(outArea, true)

	buttonGo := ui.NewButton("GO!")
	sum := big.NewInt(0)
	buttonGo.OnClicked(func(b *ui.Button) {
		go func() {
			ticker := time.NewTicker(time.Duration(dur.Value()) * time.Millisecond)

			sum.SetInt64(0)
			for i := 0; i < num.Value(); i++ {
				n, _ := rand.Int(rand.Reader, big.NewInt(*max))
				outHand.text = n.String()
				outArea.QueueRedrawAll()
				sum.Add(sum, n)

				<-ticker.C
			}

			ticker.Stop()
		}()
	})
	ovb.Append(buttonGo, true)
	number.SetChild(ovb)
	hb.Append(number, true)

	answer := ui.NewGroup("Answer")
	ivb := ui.NewVerticalBox()

	ansbox := ui.NewSpinbox(0, 1<<31-1)
	ivb.Append(ansbox, true)

	buttonTest := ui.NewButton("Test")
	buttonTest.OnClicked(func(b *ui.Button) {
		if ansbox.Value() != int(sum.Int64()) {
			ui.MsgBoxError(win, "Wrong", "Correct answer: "+sum.String())
		} else {
			ui.MsgBox(win, "Right", "Correct answer: "+sum.String())
		}
	})
	ivb.Append(buttonTest, true)

	answer.SetChild(ivb)

	hb.Append(answer, true)
	vb.Append(hb, true)

	win.SetChild(vb)
	win.Show()
}

func formSlider(form *ui.Form, label string, max int) (s *ui.Slider) {
	s = ui.NewSlider(1, max)
	form.Append(label, s, true)
	return
}

type areaHandler struct {
	text string
}

func (h *areaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	layout := ui.DrawNewTextLayout(&ui.DrawTextLayoutParams{
		String:      ui.NewAttributedString(h.text),
		DefaultFont: &ui.FontDescriptor{Family: "arial", Size: ui.TextSize(p.AreaHeight / 2), Weight: 20, Italic: 5, Stretch: 1},
		Width:       p.AreaWidth,
	})

	p.Context.Text(layout, 0, 0)
	layout.Free()
}

func (areaHandler) MouseEvent(*ui.Area, *ui.AreaMouseEvent) {}
func (areaHandler) MouseCrossed(*ui.Area, bool)             {}
func (areaHandler) DragBroken(*ui.Area)                     {}
func (areaHandler) KeyEvent(*ui.Area, *ui.AreaKeyEvent) bool {
	return false
}
