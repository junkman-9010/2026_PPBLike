package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) HandleGameInput() {
	m := g.worldMap
	p := m.player
	centerX := float32(ScreenWidth) / 2

	// 1. VIEW 모드 처리
	if g.currentMode == ModeView {
		m.UpdateCamera()
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) ||
			inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.currentMode = ModeNormal
			m.isDragging = false
			m.CenterCameraOnPlayer()
		}
		return
	}

	// 2. NORMAL 모드 처리 (마우스 클릭)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		fcx, fcy := float32(cx), float32(cy)

		// 상단 VIEW 버튼 클릭 판정
		if fcx >= centerX+55 && fcx <= centerX+115 && fcy >= 10 && fcy <= 50 {
			g.currentMode = ModeView
			return
		}

		// 타일 선택 및 이동 로직
		q, r := m.ScreenToTile(fcx, fcy)
		if q >= 0 && q < m.width && r >= 0 && r < m.height {
			// 이미 선택된 타일을 다시 클릭하면 이동
				if m.selectedQ == q && m.selectedR == r {
				// 1. 함수를 먼저 실행합니다.
				m.MovePlayerToSelected() 
				
				// 2. 이동 후의 새로운 좌표를 기반으로 트윈 목적지를 설정합니다.
				spacingX := float32(HexRadius) * 1.5
				spacingY := float32(HexRadius) * 1.73205
				targetX := float32(m.player.Q) * spacingX
				targetY := float32(m.player.R) * spacingY
				if m.player.Q%2 != 0 {
					targetY += spacingY / 2
				}
				
				p.TweenX.SetTarget(targetX)
				p.TweenY.SetTarget(targetY)
			} else {
				// 처음 클릭하면 선택만 함
				m.selectedQ, m.selectedR = q, r
			}
		}
	}
}