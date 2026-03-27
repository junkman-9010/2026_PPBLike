// menu.go
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Menu struct {
	options   []string
	selection int
}

func NewMenu() *Menu {
	return &Menu{
		options:   []string{"START GAME", "OPTIONS", "EXIT"},
		selection: 0,
	}
}

func (m *Menu) Update() Scene {
	// 위/아래 키로 선택 이동
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		m.selection = (m.selection + 1) % len(m.options)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		m.selection = (m.selection - 1 + len(m.options)) % len(m.options)
	}

	// 엔터 키로 선택 확정
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch m.selection {
		case 0: return SceneGame
		case 1: return SceneOption
		case 2: return -1 // 종료 시그널
		}
	}
	return SceneMenu
}

func (m *Menu) Draw(screen *ebiten.Image) {
    screen.Fill(color.RGBA{10, 10, 15, 255})
    
    sw := float32(ScreenWidth)
    sh := float32(ScreenHeight)

    // 1. 타이틀 위치 (중앙 상단 30% 지점)
    titleText := "PIXEL PRINCESS BLITZ"
    titleX := sw * 0.5 - 70 // 텍스트 폭 절반 보정
    titleY := sh * 0.3
    ebitenutil.DebugPrintAt(screen, titleText, int(titleX), int(titleY))

    // 2. 메뉴 박스 설정 (너비 30%, 높이 40%)
    bw, bh := sw * 0.3, sh * 0.4
    bx := (sw - bw) / 2
    by := (sh - bh) / 2 + (sh * 0.05) // 타이틀 아래로 약간 띄움

    vector.DrawFilledRect(screen, bx, by, bw, bh, color.RGBA{25, 25, 30, 255}, false)
    vector.StrokeRect(screen, bx, by, bw, bh, 1, color.RGBA{60, 60, 70, 255}, false)

    // 3. 메뉴 아이템 (박스 높이의 일정 비율 사용)
    itemH := bh / float32(len(m.options)+1) 
    for i, opt := range m.options {
        itemY := by + (itemH * float32(i)) + (itemH * 0.5)
        
        // 선택 하이라이트
        if m.selection == i {
            vector.StrokeRect(screen, bx+5, itemY-5, bw-10, itemH, 2, color.RGBA{200, 200, 255, 255}, false)
        }
        ebitenutil.DebugPrintAt(screen, opt, int(bx+(bw*0.2)), int(itemY+10))
    }
}