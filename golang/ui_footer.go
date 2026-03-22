package main

import (
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (u *UI) drawBottomUI(screen *ebiten.Image, g *Game) {
	p := g.worldMap.player
	totalWidth := float32(430)
	padding := float32(15)
	
	bgW := totalWidth + (padding * 2)
	bgH := float32(140)
	bx := float32(ScreenWidth)/2 - (bgW / 2)
	by := float32(ScreenHeight) - bgH - 10

	// 배경 판
	vector.DrawFilledRect(screen, bx, by, bgW, bgH, color.RGBA{20, 20, 20, 200}, true)
	vector.StrokeRect(screen, bx, by, bgW, bgH, 2, color.RGBA{60, 60, 60, 255}, true)

	contentX := bx + padding
	contentY := by + 15

	// 1. 스테미너 / 경험치
	u.drawStatBar(screen, contentX, contentY, totalWidth, 8, 
		p.Stamina, p.MaxStamina, color.RGBA{255, 215, 0, 255}, 
		p.Exp, p.MaxExp, color.RGBA{50, 205, 50, 255}, "ST", "EXP")

	// 2. 인벤토리 바
	u.drawInventoryBar(screen, contentX, contentY + 25, totalWidth)

	// 3. HP / MP
	u.drawStatBar(screen, contentX, contentY + 85, totalWidth, 22, 
		p.HP, p.MaxHP, color.RGBA{200, 50, 50, 255}, 
		p.MP, p.MaxMP, color.RGBA{65, 105, 225, 255}, "HP", "MP")
}

func (u *UI) drawInventoryBar(screen *ebiten.Image, x, y, totalW float32) {
	slotSize, spacing := float32(40), float32(8)
	u.drawSlot(screen, x, y, "X", "")
	u.drawSlot(screen, x+slotSize+spacing, y, "Z", "")
	
	skillStartX := x + (slotSize+spacing)*2 + 10
	skillKeys := []string{"A", "S", "D", "F", "G"}
	for i, key := range skillKeys {
		u.drawSlot(screen, skillStartX + float32(i)*(slotSize+spacing), y, key, "")
	}
	u.drawSlot(screen, x + totalW - slotSize, y, "V", "")
}