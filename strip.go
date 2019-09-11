package main

import (
	"sync"
	"math"
	"math/rand"
	"image/color"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type Effect func(pixels []color.Color)

type Strip struct {
	m sync.Mutex
	pixels []color.Color
	effect Effect
	done chan(struct{})
}

func NewStrip(size uint) *Strip {
	return &Strip{
		pixels: make([]color.Color, size),
		effect: EffectRandom,
		done: make(chan(struct{})),
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
	go func () {
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
		hue := ((i * 360 / len(pixels)) + t * 2) % 360
		pixels[i] = colorful.Hsv(float64(hue), 1, 1)
	}
}

func EffectWave(pixels []color.Color) {
	t := float64(time.Now().UnixNano() / 100000000)
	for i := 0; i < len(pixels); i++ {
		l := math.Sin(t * 0.5 + float64(i) / 10.0) / 2.0 + 0.5
		pixels[i] = colorful.Hsv(180, 1, l)
	}
}

func EffectCombined(effects ...Effect) Effect {
	return func(pixels []color.Color) {
		t := int(time.Now().Unix() / 5)
		effects[t%len(effects)](pixels)
	}
}

func GetEffectByName(name string) Effect {
	effects := make(map[string]Effect)
	effects["Random"] = EffectRandom
	effects["Rainbow"] = EffectRainbow
	effects["Wave"] = EffectWave
	effects["Combined"] = EffectCombined(
		EffectRandom,
		EffectRainbow,
		EffectWave,
	)
	return effects[name]
}
