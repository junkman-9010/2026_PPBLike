// Package main은 Pixel Princess Blitz 게임의 진입점입니다.
// Ebiten v2를 사용해 게임을 초기화하고 실행합니다.
// Go 1.20 문법을 사용합니다.
package main

import (
    "log"

    "github.com/hajimehoshi/ebiten/v2"
)

// main 함수는 화면 크기와 타이틀을 설정하고,
// Game 구조체를 초기화한 뒤 ebiten.RunGame을 호출해
// 게임 루프를 시작합니다.
// 오류가 발생하면 로그를 기록하고 프로그램을 종료합니다.
func main() {
	
	// 1. 설정 파일 먼저 로드
    LoadConfig() 

    // 2. 로드된 ScreenWidth/Height로 윈도우 설정
    ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
    ebiten.SetWindowTitle("Pixel Princess Blitz")
	
    // ScreenWidth와 ScreenHeight를 이용해 윈도우 크기를 설정합니다.
    ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
    ebiten.SetWindowTitle("Pixel Princess Blitz")

    // 게임에서 사용하는 씬, 메뉴, 월드맵, UI 등의 구성 요소를 생성하여
    // Game 인스턴스를 초기화합니다.
    game := &Game{
        currentScene: SceneMenu,      // 시작 화면을 메뉴로 설정
        menu:         NewMenu(),       // 메뉴 시스템 생성
        worldMap:     NewHexMap(40, 25), // 40x25 크기의 헥사곤 맵 생성
        ui:           NewUI(),         // UI 레이어 생성
    }

    // RunGame을 통해 게임의 Update와 Draw 루프를 실행하며,
    // 오류 발생 시 로그를 남기고 종료합니다.
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}

// Layout은 Ebiten 엔진이 사용할 논리적 화면 크기를 반환합니다.
// 고정 해상도(ScreenWidth, ScreenHeight)를 사용합니다.
// 이 함수는 Ebiten에 의해 화면 레이아웃을 결정할 때 호출됩니다.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return ScreenWidth, ScreenHeight
}