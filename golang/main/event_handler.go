//
// event_handler.go
// 게임 내 모든 사용자 입력(마우스, 키보드)을 처리하는 핸들러입니다.
//
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//
// (g *Game) HandleGameInput
// 게임 플레이 중 발생하는 모든 입력을 관리합니다.
// ModeView일 때는 카메라 조작을, ModeNormal일 때는 캐릭터 조작을 담당합니다.
// g Game 구조체의 포인터
//
func (g *Game) HandleGameInput() {
	
	m := g.worldMap
    
	centerX := float32(ScreenWidth) / 2
	
	//
	// VIEW_MODE 관찰 모드 로직
	// 드래그 활성화 및 모드 해제 조건 체크
	//
    if g.currentMode == ModeView {
        m.UpdateCamera() // 드래그 활성화

        // 우클릭 또는 ESC 누르면 모드 해제 및 카메라 복귀
        if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) || 
           inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
            g.currentMode = ModeNormal
            m.isDragging = false 
			
            m.CenterCameraOnPlayer()
        }
        return 
    }

	//
	// NORMAL_MODE 일반 모드 로직
	// 마우스 클릭을 통한 버튼 상호작용 및 캐릭터 이동
	//
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
	
	if m.MovePlayer(targetQ, targetR) {
		// 1. 이동에 성공했다면, 새로운 타일의 실제 '세계 좌표'를 가져옵니다.
		// map_system.go에 정의된 getTileScreenPos를 활용하거나 새로 계산합니다.
		spacingX := float32(40) * 1.5
		spacingY := float32(40) * 1.73205
		newX := float32(targetQ) * spacingX
		newY := float32(targetR) * spacingY
		if targetQ%2 != 0 {
			newY += spacingY / 2
		}

		// 2. 트윈의 목표값을 갱신합니다. 이제 플레이어는 이곳을 향해 미끄러집니다.
		p.TweenX.SetTarget(newX)
		p.TweenY.SetTarget(newY)

		m.UpdateVisibility()
		m.TurnCount++
	}

	// 스페이스바 입력을 통한 턴 경과 처리
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		m.TurnCount++
		m.UpdateVision(3)
		m.CalculateReachable(2)
	}
}