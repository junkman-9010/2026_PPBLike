// Package main은 게임 'Pixel Princess Blitz'의 진입점(Entry Point)이며
// 게임 인스턴스 초기화 및 메인 루프 실행을 담당합니다.
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// main 함수는 윈도우 설정(해상도, 제목)을 수행하고, 
// Game 구조체를 초기화한 뒤 Ebiten 엔진의 메인 루프를 시작합니다.
func main() {
	// ScreenWidth와 ScreenHeight 상수를 이용하여 윈도우 크기를 설정합니다.
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

// Layout은 Ebiten 엔진에서 화면의 논리적 해상도를 정의합니다. 
// 현재 설정된 ScreenWidth와 ScreenHeight를 반환하여 고정 해상도를 유지합니다.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
