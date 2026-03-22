package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Pixel Princess Blitz")

	// 게임 인스턴스 초기화
	game := &Game{
		currentScene: SceneMenu,
		menu:         NewMenu(),
		worldMap:     NewHexMap(40, 25),
		ui:           NewUI(),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Layout: 화면 크기 조정 (고정 해상도 사용)
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}