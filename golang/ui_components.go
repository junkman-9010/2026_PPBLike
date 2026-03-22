package main

import (
	"fmt"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawStatBar: 좌우로 분할된 게이지 바를 그립니다.
func (u *UI) drawStatBar(screen *ebiten.Image, x, y, w, h float32, val1, max1 int, col1 color.RGBA, val2, max2 int, col2 color.RGBA, label1, label2 string) {
	spacing := float32(10)
	halfW := (w - spacing) / 2

	// 좌측 바
	vector.DrawFilledRect(screen, x, y, halfW, h, color.RGBA{30, 30, 30, 255}, true)
	if max1 > 0 {
		ratio1 := float32(val1) / float32(max1)
		vector.DrawFilledRect(screen, x, y, halfW*ratio1, h, col1, true)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s %d/%d", label1, val1, max1), int(x)+5, int(y)+2)

	// 우측 바
	rightX := x + halfW + spacing
	vector.DrawFilledRect(screen, rightX, y, halfW, h, color.RGBA{30, 30, 30, 255}, true)
	if max2 > 0 {
		ratio2 := float32(val2) / float32(max2)
		vector.DrawFilledRect(screen, rightX, y, halfW*ratio2, h, col2, true)
	}
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s %d/%d", label2, val2, max2), int(rightX)+5, int(y)+2)
}

// drawSlot: 개별 인벤토리 슬롯을 그립니다.
func (u *UI) drawSlot(screen *ebiten.Image, x, y float32, key, label string) {
	vector.StrokeRect(screen, x, y, 40, 40, 2, color.RGBA{100, 100, 100, 255}, true)
	ebitenutil.DebugPrintAt(screen, key, int(x)+15, int(y)+45)
}