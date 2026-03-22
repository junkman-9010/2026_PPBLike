package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawTopUI: 상단 바 전체와 내부 요소들을 그립니다.
func (u *UI) drawTopUI(screen *ebiten.Image, g *Game) {
	// 1. 상단 바 배경 (살짝 투명한 검은색)
	vector.DrawFilledRect(screen, 0, 0, float32(ScreenWidth), 60, color.RGBA{0, 0, 0, 100}, true)

	centerX := float32(ScreenWidth) / 2

	// 2. 좌측 상태 (중앙 구체 왼쪽 배치)
	// 골드 아이콘 및 수치
	u.drawIconButton(screen, centerX-220, 10, "GOLD", fmt.Sprintf("%d", g.worldMap.player.Gold), color.RGBA{218, 165, 32, 255})
	
	// 식량(Food) 아이콘 및 수치 (10/10 형식)
	u.drawIconButton(screen, centerX-150, 10, "FOOD", fmt.Sprintf("%d/%d", g.worldMap.player.Food, g.worldMap.player.MaxFood), color.RGBA{139, 69, 19, 255})
	
	// 캐릭터 초상화 (HP가 낮으면 붉은 테두리 연출 가능)
	u.drawPortrait(screen, centerX-85, 5, g.worldMap.player.HP < 20)

	// 3. 중앙 시간 장치 (DAY 및 시간대 표시)
	u.drawCenterTime(screen, g.worldMap)

	// 4. 우측 메뉴 버튼 (중앙 구체 오른쪽 배치)
	u.drawMenuButton(screen, centerX+55, 10, "VIEW", color.RGBA{150, 50, 50, 255})
	u.drawMenuButton(screen, centerX+120, 10, "STATS", color.RGBA{50, 100, 100, 255})
	u.drawMenuButton(screen, centerX+185, 10, "INVT", color.RGBA{80, 80, 80, 255})
	u.drawMenuButton(screen, centerX+250, 10, "DECK", color.RGBA{100, 80, 40, 255})
}

// drawCenterTime: 중앙의 황금색 구체와 시간 정보를 그립니다.
func (u *UI) drawCenterTime(screen *ebiten.Image, m *HexMap) {
	
	timeStepRaw, hour := m.GetTimeContext() // 1. 변수 이름을 잠시 바꿉니다.
	timeStep := int(timeStepRaw)           // 2. 확실하게 숫자로 변환합니다.
	day := (m.TurnCount / 24) + 1
	cx := float32(ScreenWidth) / 2

	var bubbleColor color.RGBA
	var timeLabel string

	// timeStep을 int로 확실히 받아서 처리합니다.
	switch int(timeStep) {
	case 0:
		bubbleColor = color.RGBA{255, 230, 150, 200}
		timeLabel = "MORNING"
	case 1:
		bubbleColor = color.RGBA{255, 180, 50, 220}
		timeLabel = "AFTERNOON"
	case 2:
		bubbleColor = color.RGBA{200, 100, 50, 200}
		timeLabel = "EVENING"
	case 3:
		bubbleColor = color.RGBA{40, 40, 100, 220}
		timeLabel = "NIGHT"
	default:
		bubbleColor = color.RGBA{255, 180, 50, 220}
		timeLabel = "UNKNOWN"
	}

	vector.DrawFilledCircle(screen, cx, 30, 35, bubbleColor, true)
	vector.StrokeCircle(screen, cx, 30, 35, 2, color.RGBA{218, 165, 32, 255}, true)
	
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("DAY %d", day), int(cx)-18, 15)
	ebitenutil.DebugPrintAt(screen, timeLabel, int(cx)-28, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%02d:00", hour), int(cx)-15, 45)
}

// drawIconButton: 아이콘과 텍스트가 결합된 형태를 그립니다.
func (u *UI) drawIconButton(screen *ebiten.Image, x, y float32, label, value string, clr color.RGBA) {
	vector.DrawFilledRect(screen, x, y, 60, 40, color.RGBA{40, 40, 40, 200}, true)
	vector.DrawFilledRect(screen, x+2, y+2, 15, 15, clr, true) // 아이콘 대용 사각형
	ebitenutil.DebugPrintAt(screen, label, int(x)+20, int(y)+2)
	ebitenutil.DebugPrintAt(screen, value, int(x)+5, int(y)+22)
}

// drawPortrait: 캐릭터 초상화 영역을 그립니다.
func (u *UI) drawPortrait(screen *ebiten.Image, x, y float32, isLowHP bool) {
	borderColor := color.RGBA{150, 150, 150, 255}
	if isLowHP {
		borderColor = color.RGBA{255, 0, 0, 255} // 위급 상황 시 붉은색
	}
	vector.DrawFilledRect(screen, x, y, 50, 50, color.RGBA{60, 60, 60, 255}, true)
	vector.StrokeRect(screen, x, y, 50, 50, 2, borderColor, true)
}

// drawMenuButton: 상단 우측의 메뉴 버튼들을 그립니다.
func (u *UI) drawMenuButton(screen *ebiten.Image, x, y float32, text string, clr color.RGBA) {
	vector.DrawFilledRect(screen, x, y, 60, 40, color.RGBA{40, 40, 40, 200}, true)
	vector.DrawFilledRect(screen, x, y+35, 60, 5, clr, true) // 하단 포인트 색상
	ebitenutil.DebugPrintAt(screen, text, int(x)+10, int(y)+12)
}