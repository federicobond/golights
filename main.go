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

	selectedEffect := "Combined"

	var effectsSelector *widget.Radio
	effectsSelector = widget.NewRadio([]string{
		"Random",
		"Rainbow",
		"Wave",
		"Colored Wave",
		"Combined",
	}, func(selected string) {
		if selected == "" {
			effectsSelector.SetSelected(selectedEffect)
			return
		}
		selectedEffect = selected
		strip.SetEffect(GetEffectByName(selected))
	})

	effectsSelector.SetSelected("Combined")

	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Golights!"),
		effectsSelector,
		lightsWidget,
		widget.NewButton("Quit", func() {
			strip.Stop()
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
