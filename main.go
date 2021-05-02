package main

//#include <windows.h>
import "C"
import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"
	"unsafe"

	. "github.com/andlabs/ui"
)

func main() {
	Main(setupUI)
}

func setupUI() {
	win := mainWindows()

	content := NewVerticalBox()
	params := NewGroup("Parameters")
	form := NewForm()

	dur := formSlider(form, "Showing time", "%dms", 100, 5000)
	num := formSlider(form, "Shows number", "%dx", 1, 50)
	dig := formSlider(form, "Digits number", "%dx", 1, 10)

	params.SetChild(form)
	content.Append(params, true)

	hb := NewHorizontalBox()
	number := NewGroup("Number")
	ovb := NewVerticalBox()

	outHand := &areaHandler{text: "0"}
	outArea := NewArea(outHand)
	outArea.Disable()

	ovb.Append(outArea, true)

	sum := big.NewInt(0)
	buttonGo := NewButton("Go!")
	buttonTest := NewButton("Test")
	buttonGo.OnClicked(func(b *Button) {
		buttonGo.Disable()
		buttonTest.Disable()

		go func() {
			ticker := time.NewTicker(time.Duration(dur.Value()) * time.Millisecond)

			sum.SetInt64(0)
			for i := 0; i < num.Value(); i++ {
				n, _ := rand.Int(rand.Reader, big.NewInt(int64(math.Pow10(dig.Value()))))
				outHand.text = ""
				outArea.QueueRedrawAll()
				<-time.After(50 * time.Millisecond)
				outHand.text = n.String()
				outArea.QueueRedrawAll()
				sum.Add(sum, n)

				<-ticker.C
			}

			ticker.Stop()

			buttonGo.Enable()
			buttonTest.Enable()
		}()
	})

	ovb.Append(buttonGo, true)
	number.SetChild(ovb)
	hb.Append(number, true)

	answer := NewGroup("Answer")
	ivb := NewVerticalBox()

	ansbox := NewSpinbox(0, 1<<31-1)
	ivb.Append(ansbox, true)

	buttonTest.OnClicked(func(b *Button) {
		switch ansbox.Value() {
		case int(sum.Int64()):
			MsgBox(win, "Right", "Correct answer: "+sum.String())
		default:
			MsgBoxError(win, "Wrong", "Correct answer: "+sum.String())
		}
	})

	ivb.Append(buttonTest, true)
	answer.SetChild(ivb)
	hb.Append(answer, true)
	content.Append(hb, true)
	win.SetChild(content)
	C.SetWindowPos(C.HWND(unsafe.Pointer(win.Handle())), C.HWND_TOPMOST, C.int(0), C.int(0), C.int(0), C.int(0), C.SWP_NOSIZE|C.SWP_NOMOVE|C.SWP_SHOWWINDOW|C.SWP_ASYNCWINDOWPOS)
}

func mainWindows() *Window {
	win := NewWindow("summation trainer", 0, 0, true)
	win.OnClosing(func(w *Window) bool {
		Quit()
		return true
	})
	win.SetMargined(true)

	OnShouldQuit(func() bool {
		win.Destroy()
		return true
	})

	return win
}

func formSlider(form *Form, label, format string, min, max int) (s *Slider) {
	vb := NewVerticalBox()
	s = NewSlider(min, max)
	l := NewLabel(fmt.Sprintf(format, s.Value()))
	s.OnChanged(func(s *Slider) {
		l.SetText(fmt.Sprintf(format, s.Value()))
	})
	vb.Append(l, false)
	vb.Append(s, false)
	form.Append(label, vb, true)
	return
}

type areaHandler struct {
	text string
}

func (h *areaHandler) Draw(a *Area, p *AreaDrawParams) {
	layout := DrawNewTextLayout(&DrawTextLayoutParams{
		String:      NewAttributedString(h.text),
		DefaultFont: &FontDescriptor{Family: "Arial", Size: TextSize(p.AreaHeight / 2), Weight: 20, Italic: 5, Stretch: 1},
		Width:       p.AreaWidth,
	})

	p.Context.Text(layout, p.AreaWidth/2-p.AreaHeight/3*(float64(len([]rune(h.text)))/2), 0)
	layout.Free()
}

func (areaHandler) MouseEvent(a *Area, e *AreaMouseEvent) {}

func (areaHandler) MouseCrossed(a *Area, b bool) {}

func (areaHandler) DragBroken(a *Area) {}

func (areaHandler) KeyEvent(a *Area, k *AreaKeyEvent) bool {
	return false
}
