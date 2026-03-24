package main

import (
	"fmt"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawStatBar는 두 개의 수치 데이터(예: HP/MP 또는 ST/EXP)를 받아 하나의 통합된 바 형태로 렌더링합니다.
// 배경 바 위에 현재 값의 비율만큼 컬러 바를 채우고, 좌측에 라벨(HP, MP 등)을 표시합니다.
func (u *UI) drawStatBar(screen *ebiten.Image, x, y, width, height float32, 
	val1, max1 int, color1 color.RGBA, 
	val2, max2 int, color2 color.RGBA, 
	label1, label2 string) {
	
	// 각 바의 너비 (중앙 간격 제외하고 절반씩 배분)
	barW := (width / 2) - 10
	
	// 첫 번째 바 (왼쪽: 예 - HP, ST)
	u.renderSingleBar(screen, x, y, barW, height, val1, max1, color1, label1)
	
	// 두 번째 바 (오른쪽: 예 - MP, EXP)
	u.renderSingleBar(screen, x + barW + 20, y, barW, height, val2, max2, color2, label2)
}

// renderSingleBar는 단일 게이지 바를 화면에 그리는 내부 보조 함수입니다.
func (u *UI) renderSingleBar(screen *ebiten.Image, x, y, w, h float32, val, max int, clr color.RGBA, label string) {
	// 1. 배경 (어두운 회색)
	vector.DrawFilledRect(screen, x, y, w, h, color.RGBA{40, 40, 40, 255}, true)
	
	// 2. 현재 값 비율 계산
	ratio := float32(0)
	if max > 0 {
		ratio = float32(val) / float32(max)
	}
	if ratio > 1 { ratio = 1 }

	// 3. 게이지 채우기
	vector.DrawFilledRect(screen, x, y, w*ratio, h, clr, true)
	
	// 4. 테두리
	vector.StrokeRect(screen, x, y, w, h, 1, color.RGBA{100, 100, 100, 255}, true)

	// 5. 텍스트 표시 (라벨 및 수치)
	statText := fmt.Sprintf("%s: %d/%d", label, val, max)
	ebitenutil.DebugPrintAt(screen, statText, int(x)+5, int(y)+int(h/2)-7)
}

// drawSlot: 개별 인벤토리 슬롯을 그립니다.
func (u *UI) drawSlot(screen *ebiten.Image, x, y float32, key, label string) {
	vector.StrokeRect(screen, x, y, 40, 40, 2, color.RGBA{100, 100, 100, 255}, true)
	ebitenutil.DebugPrintAt(screen, key, int(x)+15, int(y)+45)
}