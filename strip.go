package main

import (
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type Effect func(pixels []color.Color)

type Strip struct {
	m      sync.Mutex
	pixels []color.Color
	effect Effect
	done   chan (struct{})
}

func NewStrip(size uint) *Strip {
	return &Strip{
		pixels: make([]color.Color, size),
		effect: EffectRandom,
		done:   make(chan (struct{})),
	}
}

func (s *Strip) Size() int {
	return len(s.pixels)
}

func (s *Strip) SetEffect(effect Effect) {
	s.m.Lock()
	defer s.m.Unlock()
	s.effect = effect
}

func (s *Strip) Pixels() []color.Color {
	s.m.Lock()
	defer s.m.Unlock()
	return s.pixels
}

func (s *Strip) Run(callback func()) {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.m.Lock()
				s.effect(s.pixels)
				s.m.Unlock()
				callback()
			case <-s.done:
				return
			}
		}
	}()
}

func (s *Strip) Stop() {
	close(s.done)
}

func EffectRandom(pixels []color.Color) {
	for i := 0; i < len(pixels); i++ {
		pixels[i] = color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		}
	}
}

func EffectRainbow(pixels []color.Color) {
	t := int(time.Now().UnixNano() / 10000000)
	for i := 0; i < len(pixels); i++ {
		hue := ((i * 360 / len(pixels)) + t*2) % 360
		pixels[i] = colorful.Hsv(float64(hue), 1, 1)
	}
}

func EffectWave(pixels []color.Color) {
	t := float64(time.Now().UnixNano() / 100000000)
	for i := 0; i < len(pixels); i++ {
		l := math.Sin(t*0.5+float64(i)/10.0)/2.0 + 0.5
		pixels[i] = colorful.Hsv(180, 1, l)
	}
}

func EffectColoredWave(pixels []color.Color) {
	t := float64(time.Now().UnixNano() / 100000000)
	for i := 0; i < len(pixels); i++ {
		h := ((i * 360 / len(pixels)) + int(t)*2) % 360
		l := math.Sin(t*0.5+float64(i)/10.0)/2.0 + 0.5
		pixels[i] = colorful.Hsv(float64(h), 1, l)
	}
}

func EffectCombined(effects ...Effect) Effect {
	return func(pixels []color.Color) {
		t := int(time.Now().Unix() / 5)
		effects[t%len(effects)](pixels)
	}
}

func EffectBlend(pixels []color.Color) {
	t := float64(time.Now().UnixNano() / 100000000)
	alpha := math.Sin(t*0.1)/2.0 + 0.5

	p1 := make([]color.Color, len(pixels))
	EffectRandom(p1)
	p2 := make([]color.Color, len(pixels))
	EffectWave(p2)

	for i := 0; i < len(pixels); i++ {
		pixels[i] = colorBlend(p1[i], p2[i], alpha)
	}
}

func colorBlend(c1, c2 color.Color, alpha float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return color.RGBA64{
		uint16(float64(r1)*alpha + float64(r2)*(1-alpha)),
		uint16(float64(g1)*alpha + float64(g2)*(1-alpha)),
		uint16(float64(b1)*alpha + float64(b2)*(1-alpha)),
		uint16(float64(a1)*alpha + float64(a2)*(1-alpha)),
	}
}

func GetEffectByName(name string) Effect {
	effects := make(map[string]Effect)
	effects["Random"] = EffectRandom
	effects["Rainbow"] = EffectRainbow
	effects["Wave"] = EffectWave
	effects["Blend"] = EffectBlend
	effects["Colored Wave"] = EffectColoredWave
	effects["Combined"] = EffectCombined(
		EffectRandom,
		EffectRainbow,
		EffectWave,
		EffectColoredWave,
	)
	return effects[name]
}
