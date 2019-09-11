package main

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/canvas"
)

type LightsWidget struct {
	pos fyne.Position
	strip *Strip
}

func NewLightsWidget(strip *Strip) *LightsWidget {
	w := &LightsWidget{strip: strip}
	widget.Renderer(w).Layout(w.MinSize())
	return w
}

func (w *LightsWidget) SetStrip(strip *Strip) {
	w.strip = strip
}

func (w *LightsWidget) Size() fyne.Size {
	return fyne.NewSize((w.strip.Size() + 1) * 20 + 10, 50)
	
}

func (w *LightsWidget) Resize(size fyne.Size) {

}

func (w *LightsWidget) Position() fyne.Position {
	return w.pos
}

func (w *LightsWidget) Move(pos fyne.Position) {
	w.pos = pos
}

func (w *LightsWidget) MinSize() fyne.Size {
	return w.Size()
}

func (w *LightsWidget) Visible() bool {
	return true
}
func (w *LightsWidget) Show() {}
func (w *LightsWidget) Hide() {}


func (w *LightsWidget) CreateRenderer() fyne.WidgetRenderer {
	return &LightsWidgetRenderer{w: w}
}

type LightsWidgetRenderer struct {
	w *LightsWidget
	objects []fyne.CanvasObject
}

func (r *LightsWidgetRenderer) Layout(size fyne.Size) {
	r.Refresh()
}

func (r *LightsWidgetRenderer) MinSize() fyne.Size {
	return r.w.MinSize()
}

func (r *LightsWidgetRenderer) Refresh() {
	if r.objects == nil {
		var objs []fyne.CanvasObject
		for i, color := range r.w.strip.Pixels() {
			o := canvas.NewCircle(color)
			o.Resize(fyne.NewSize(10, 10))
			o.Move(fyne.NewPos(
				20+(i%splitLength)*20,
				20+(i/splitLength)*20,
			))
			objs = append(objs, o)
		}
		r.objects = objs
	} else {
		for i, color := range r.w.strip.Pixels() {
			circle := r.objects[i].(*canvas.Circle)
			circle.FillColor = color
		}
	}
	canvas.Refresh(r.w)
}

func (r *LightsWidgetRenderer) ApplyTheme() {}

func (r *LightsWidgetRenderer) BackgroundColor() color.Color {
	return color.Black
}

func (r *LightsWidgetRenderer) Objects() []fyne.CanvasObject{
	return r.objects
}
func (r *LightsWidgetRenderer) Destroy() {}