package main

import "github.com/hajimehoshi/ebiten/v2"

// UI 구조체 정의
type UI struct {
	// 향후 UI 상태(애니메이션 타이머, 선택된 슬롯 등)를 관리할 필드를 추가할 수 있습니다.
}

// NewUI: UI 인스턴스를 생성하여 반환합니다. (오류 해결 지점)
func NewUI() *UI {
	return &UI{}
}

// Draw: 상단과 하단 UI를 순서대로 호출하는 메인 그리기 함수
func (u *UI) Draw(screen *ebiten.Image, g *Game) {
	// ui_header.go에 정의된 함수 호출
	u.drawTopUI(screen, g)

	// ui_footer.go에 정의된 함수 호출
	u.drawBottomUI(screen, g)
}