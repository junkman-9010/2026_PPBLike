// ui.go
//
package main

import "github.com/hajimehoshi/ebiten/v2"

// UI는 게임 화면에 표시되는 상단과 하단 UI를 관리하는 구조체입니다.
type UI struct {
	// 향후 UI 상태(애니메이션 타이머, 선택된 슬롯 등)를 관리할 필드를 추가할 수 있습니다.
}

// NewUI는 UI 인스턴스를 생성하여 반환합니다.
// UI가 초기화되는 과정에서 특별한 처리가 필요하면 여기에 구현하면 됩니다.
// 반환값은 생성된 UI 포인터이며, nil 값이 반환되는 경우는 없습니다.
func NewUI() *UI {
	return &UI{}
}

// Draw는 화면에 UI를 그리는 메인 함수입니다.
// screen : ebiten의 화면 버퍼입니다.
// g      : 현재 게임 상태를 담고 있는 구조체입니다.
//
// Draw 함수는 내부적으로 상단 UI와 하단 UI를 순차적으로 호출합니다.
// 각 하위 함수는 ui_header.go와 ui_footer.go에 구현되어 있습니다.
func (u *UI) Draw(screen *ebiten.Image, g *Game) {
	// ui_header.go에 정의된 함수 호출
	u.drawTopUI(screen, g)

	// ui_footer.go에 정의된 함수 호출
	u.drawBottomUI(screen, g)
}