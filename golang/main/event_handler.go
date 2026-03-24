/**
 * @file event_handler.go
 * @brief 게임 내 모든 사용자 입력(마우스, 키보드)을 처리하는 핸들러입니다.
 */
 
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

/**
 * @fn (g *Game) HandleGameInput
 * @brief 게임 플레이 중 발생하는 모든 입력을 관리합니다.
 * @details ModeView일 때는 카메라 조작을, ModeNormal일 때는 캐릭터 조작을 담당합니다.
 * @param g Game 구조체의 포인터
 */
func (g *Game) HandleGameInput() {
	
	m := g.worldMap
    
	centerX := float32(ScreenWidth) / 2
	
	/**
	 * @section VIEW_MODE 관찰 모드 로직
     * 드래그 활성화 및 모드 해제 조건 체크
     */
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

	/**
	 * @section NORMAL_MODE 일반 모드 로직
     * 마우스 클릭을 통한 버튼 상호작용 및 캐릭터 이동
     */
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

	/**
	 * @brief 스페이스바 입력을 통한 턴 경과 처리
	 */
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		m.TurnCount++
		m.UpdateVision(3)
		m.CalculateReachable(2)
	}
}