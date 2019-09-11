package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	app := app.New()

	strip := NewStrip(150)
	// strip.SetEffect(EffectWave)
	strip.SetEffect(EffectCombined(
		EffectRandom,
		EffectRainbow,
		EffectWave,
	))

	lightsWidget := NewLightsWidget(strip)

	w := app.NewWindow("Golights")

	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Golights!"),
		lightsWidget,
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	ch, cleanup := dmx()
	defer cleanup()

	strip.Run(func() {
		widget.Refresh(lightsWidget)
		ch <- toDMX(strip.Pixels())
	})

	w.ShowAndRun()
}
