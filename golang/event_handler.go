package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// HandleGameInput: 게임 플레이 중 발생하는 모든 입력을 관리합니다.
func (g *Game) HandleGameInput() {
	m := g.worldMap

	// 1. 마우스 왼쪽 버튼 클릭 (이동 및 선택)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// 드래그 중이 아니었을 때만 클릭으로 인정 (중요!)
		if !m.isDragging {
			cx, cy := ebiten.CursorPosition()
			q, r := m.ScreenToTile(float32(cx), float32(cy))

			// 타일 선택 로직
			if q >= 0 && q < m.width && r >= 0 && r < m.height {
				// 이미 선택된 타일을 다시 클릭하면 이동
				if m.selectedQ == q && m.selectedR == r {
					m.MovePlayerToSelected()
				} else {
					// 새로운 타일 선택
					m.selectedQ, m.selectedR = q, r
				}
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