package main

import (
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Menu struct {
	selection int
	options   []string
}

func NewMenu() *Menu {
	return &Menu{
		options: []string{"START GAME", "OPTIONS", "EXIT"},
	}
}

func (m *Menu) Update() Scene {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.selection = (m.selection - 1 + len(m.options)) % len(m.options)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.selection = (m.selection + 1) % len(m.options)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch m.selection {
		case 0: return SceneGame
		case 1: return SceneOption
		case 2: return -1
		}
	}
	return SceneMenu
}

func (m *Menu) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 15, 255})
	
	// Menu Box
	bx, by, bw, bh := float32(500), float32(280), float32(280), float32(180)
	vector.StrokeRect(screen, bx, by, bw, bh, 2, color.RGBA{80, 80, 120, 255}, false)
	
	ebitenutil.DebugPrintAt(screen, "PIXEL PRINCESS BLITZ", 560, 220)

	for i, option := range m.options {
		str := option
		if i == m.selection {
			str = "> " + option
			vector.DrawFilledRect(screen, bx+5, by+float32(30+(i*45)), bw-10, 25, color.RGBA{255, 255, 255, 20}, false)
		}
		ebitenutil.DebugPrintAt(screen, str, 580, 320+(i*45))
	}
}