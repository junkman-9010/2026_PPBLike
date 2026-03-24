// ui_footer.go
//
// 이 파일은 게임 화면 하단 UI(바닥 UI)를 그리는 함수들을 모아 두었습니다.
// 주로 인벤토리, 스탯바, HP/MP 바 등을 화면에 렌더링합니다.
// 각 함수는 `UI` 구조체의 메서드로 구현되어 있으며, 화면 버퍼와
// 현재 게임 상태를 담은 `Game` 포인터를 인자로 받습니다.
//
// godoc 주석은 모든 함수 앞에 `//` 로 시작하며, 한글로 작성되었습니다.
//

package main

import (
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawInventoryBar
//
// 인벤토리 바(장비, 스킬 키 등)를 그립니다.
// 인벤토리 슬롯은 간단히 텍스트로 표시되며, 나중에 실제 아이템 이미지를 넣을 수 있습니다.
//
//   screen : 화면 버퍼
//   x, y   : 인벤토리 바 시작 좌표
//   totalW : 인벤토리 바 전체 너비
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

// drawBottomUI
//
// 화면 하단 바를 그립니다. 배경판, 스탯바, 인벤토리 바, HP/MP 바를
// 한 화면에 묶어 그리도록 설계되었습니다.
//
//   screen : ebiten.Image 타입의 화면 버퍼
//   g      : 현재 게임 전역 상태를 담은 Game 객체
func (u *UI) drawBottomUI(screen *ebiten.Image, g *Game) {
	p := g.worldMap.player
	totalWidth := float32(430)
	padding := float32(15)
	
	bgW := totalWidth + (padding * 2)
	bgH := float32(140)
	bx := float32(ScreenWidth)/2 - (bgW / 2)
	by := float32(ScreenHeight) - bgH - 10

	// 전체 배경 판
	vector.DrawFilledRect(screen, bx, by, bgW, bgH, color.RGBA{20, 20, 20, 220}, true)
	vector.StrokeRect(screen, bx, by, bgW, bgH, 2, color.RGBA{80, 80, 80, 255}, true)

	contentX := bx + padding
	contentY := by + 15

	// 1. 상단: 스테미너(ST)와 경험치(EXP) 바 (얇게 표시)
	u.drawStatBar(screen, contentX, contentY, totalWidth, 12, 
		p.Stamina, p.MaxStamina, color.RGBA{255, 215, 0, 255}, 
		p.Exp, p.MaxExp, color.RGBA{50, 205, 50, 255}, "ST", "EXP")

	// 2. 중단: 인벤토리 및 스킬 슬롯
	u.drawInventoryBar(screen, contentX, contentY + 25, totalWidth)

	// 3. 하단: HP와 MP 바 (가독성을 위해 두껍게 표시)
	u.drawStatBar(screen, contentX, contentY + 85, totalWidth, 22, 
		p.HP, p.MaxHP, color.RGBA{200, 50, 50, 255}, 
		p.MP, p.MaxMP, color.RGBA{65, 105, 225, 255}, "HP", "MP")
}