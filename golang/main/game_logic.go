package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	currentScene Scene
	currentMode  int // 추가: ModeNormal 또는 ModeView
	menu         *Menu
	worldMap     *HexMap
	ui           *UI
	
	optionUI *OptionUI // 추가
}

func (g *Game) Update() error {
	switch g.currentScene {
	case SceneMenu:
		next := g.menu.Update()
		if next == -1 {
			os.Exit(0)
		}
		g.currentScene = next

	case SceneGame:
		// 1. 카메라 업데이트 (드래그 로직 포함)
		//g.worldMap.UpdateCamera()

		// 2. [수정] 분리한 이벤트 핸들러 호출
		g.HandleGameInput()
		
		// 플레이어의 현재 위치(Current)를 목표 위치(Target)로 부드럽게 이동시킴
        g.worldMap.player.TweenX.Update()
        g.worldMap.player.TweenY.Update()

		// 3. ESC 메뉴
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.currentScene = SceneMenu
		}
	
	case SceneOption:
		if g.optionUI == nil { g.optionUI = NewOptionUI() } // 초기화
		g.optionUI.Update(g)
	}
	return nil
}

// handleGameInput: 인게임 내에서의 클릭 및 상호작용 로직 분리
func (g *Game) handleGameInput() {
	// 방향키를 이용한 선택 좌표 이동
	if inpututil.IsKeyJustPressed(ebiten.KeyUp)    { g.worldMap.selectedR-- }
	if inpututil.IsKeyJustPressed(ebiten.KeyDown)  { g.worldMap.selectedR++ }
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft)   { g.worldMap.selectedQ-- }
	if inpututil.IsKeyJustPressed(ebiten.KeyRight)  { g.worldMap.selectedQ++ }

	// 스페이스바: 선택한 좌표로 플레이어 이동 및 턴 소모
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.worldMap.MovePlayerToSelected()
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.currentScene {
	case SceneMenu:
		g.menu.Draw(screen)
	case SceneGame:
		// 맵과 엔티티 그리기
		//g.worldMap.Draw(screen)
		g.worldMap.Draw(screen, g.currentMode)
		// UI 레이어 그리기 (항상 최상단)
		g.ui.Draw(screen, g)
	case SceneOption:
        if g.optionUI != nil {
            g.optionUI.Draw(screen)
        }
	}
}