// Copyright 2022 Mark Mandriota. All right reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"math"
	"math/rand"
	"time"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	run(app.New().NewWindow("Number Trainer"))
}

func run(w fyne.Window) {
	timing := binding.NewFloat()
	timing.Set(1)
	timingP := newParameterSliderWithData("FPS", 1, 24, timing)

	amount := binding.NewFloat()
	amount.Set(3)
	amountP := newParameterSliderWithData("AMOUNT", 1, 12, amount)

	digits := binding.NewFloat()
	digits.Set(1)
	digitsP := newParameterSliderWithData("DIGITS", 1, 12, digits)

	sign := binding.NewBool()
	signP := widget.NewCheckWithData("SIGN", sign)

	random := binding.NewInt()
	answer := binding.NewInt()

	viewer := widget.NewLabelWithData(binding.IntToString(random))
	tester := widget.NewButton("TRY", func() {
		maxN := math.Pow10(int(try(digits.Get())))
		sign := try(sign.Get())

		ticker := time.NewTicker(time.Duration(float64(time.Second) / try(timing.Get())))
		defer ticker.Stop()

		n, sum := 0, 0
		for i := int(try(amount.Get())); i > 0; i-- {
			if sign {
				n = rand.Intn(int(maxN)*2) - int(maxN)
			} else {
				n = rand.Intn(int(maxN))
			}

			random.Set(0)
			random.Set(n)
			sum += n

			<-ticker.C

			viewer.SetText("")

			<-time.After(time.Duration(float64(time.Second) / try(timing.Get()) / 2))
		}

		answerP := widget.NewEntryWithData(binding.IntToString(answer))

		dialog.NewForm("SOLUTION", "", "", []*widget.FormItem{
			widget.NewFormItem("ANSWER", answerP),
		}, func(b bool) {
			if b {
				if asum := try(answer.Get()); asum == sum {
					viewer.SetText("RIGHT!")
				} else if asum < sum {
					viewer.SetText("TOO SMALL..")
				} else {
					viewer.SetText("TOO LARGE..")
				}
			}
		}, w).Show()
	})

	w.SetContent(container.NewVBox(
		timingP,
		layout.NewSpacer(),
		amountP,
		layout.NewSpacer(),
		digitsP,
		layout.NewSpacer(),
		container.NewGridWithColumns(3, signP, viewer, tester),
	))

	w.ShowAndRun()
}

func newParameterSliderWithData(name string, min, max float64, data binding.Float) fyne.CanvasObject {
	return container.NewGridWithColumns(2,
		widget.NewSliderWithData(min, max, data),
		container.NewGridWithColumns(2,
			widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "%.0f")),
			widget.NewLabel(name),
		),
	)
}

func try[T any](v T, e error) T {
	if e != nil {
		panic(e)
	}

	return v
}
