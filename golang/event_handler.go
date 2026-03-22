package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// HandleGameInput: 게임 플레이 중 발생하는 모든 입력을 관리합니다.
func (g *Game) HandleGameInput() {
	
	m := g.worldMap
    
	centerX := float32(ScreenWidth) / 2

    // --- 1. 관찰 모드(ModeView)일 때 ---
    if g.currentMode == ModeView {
        m.UpdateCamera() // 드래그 활성화

        // 우클릭 또는 ESC 누르면 모드 해제
        if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) || 
           inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
            g.currentMode = ModeNormal
            m.isDragging = false 
			
			// 관찰 모드 종료 즉시 플레이어에게 카메라 복귀
            m.CenterCameraOnPlayer()
        }
        return 
    }

    // --- 2. 일반 모드(ModeNormal)일 때 ---
	// 일반 모드에서는 m.UpdateCamera()를 호출하지 않으므로 드래그가 작동하지 않습니다.
    if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
        cx, cy := ebiten.CursorPosition()
        fcx, fcy := float32(cx), float32(cy)

        // 상단 VIEW 버튼 클릭 판정 (centerX+55 위치)
        if fcx >= centerX+55 && fcx <= centerX+115 && fcy >= 10 && fcy <= 50 {
            g.currentMode = ModeView
            return
        }

        // 일반적인 타일 선택 및 이동
        q, r := m.ScreenToTile(fcx, fcy)
        if q >= 0 && q < m.width && r >= 0 && r < m.height {
            if m.selectedQ == q && m.selectedR == r {
                m.MovePlayerToSelected()
            } else {
                m.selectedQ, m.selectedR = q, r
            }
        }
    }

	// 2. 키보드 입력 (스페이스바: 턴 넘기기 등)
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		m.TurnCount++
		m.UpdateVision(3)
		m.CalculateReachable(2)
	}
}