package main

import (
	"fmt"

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

// 별도의 OptionUpdate 함수 (Game.Update에서 SceneOption일 때 호출)
func (g *Game) HandleOptionInput() {
    // 좌우 키로 해상도 변경
    if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
        idx := (CurrentResIndex - 1 + len(SupportedResolutions)) % len(SupportedResolutions)
        UpdateConfigValue(idx)
        ebiten.SetWindowSize(SupportedResolutions[idx].Width, SupportedResolutions[idx].Height)
    }
    if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
        idx := (CurrentResIndex + 1) % len(SupportedResolutions)
        UpdateConfigValue(idx)
        ebiten.SetWindowSize(SupportedResolutions[idx].Width, SupportedResolutions[idx].Height)
    }

    // ESC로 메인 메뉴 복귀
    if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
        g.currentScene = SceneMenu
    }
}

func (m *Menu) DrawOptions(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 30, 255}) // 옵션 배경색 (약간 진한 남색)

	// 안내 문구 출력
	ebitenutil.DebugPrintAt(screen, "--- OPTIONS ---", ScreenWidth/2-50, ScreenHeight/2-60)
	
	// 현재 해상도 표시
	resText := fmt.Sprintf("Resolution: < %s >", SupportedResolutions[CurrentResIndex].Name)
	ebitenutil.DebugPrintAt(screen, resText, ScreenWidth/2-100, ScreenHeight/2)

	ebitenutil.DebugPrintAt(screen, "Press LEFT/RIGHT to Change", ScreenWidth/2-100, ScreenHeight/2+40)
	ebitenutil.DebugPrintAt(screen, "Press ENTER or ESC to Save & Back", ScreenWidth/2-110, ScreenHeight/2+80)
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