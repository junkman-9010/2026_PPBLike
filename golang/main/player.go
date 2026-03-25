//
// player.go
// 플레이어 캐릭터의 데이터 구조와 메서드를 정의합니다.
// 
package main

import (
	"image/color"
	//"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//
// Player struct
// brief 플레이어의 능력치 및 위치 정보를 담는 구조체입니다.
//
type Player struct {
	Q, R int
	
	// 기본 능력치
	HP, MaxHP     int
	MP, MaxMP     int
	Gold          int
	Food, MaxFood int
	
	// 행동 관련 데이터
	Stamina, MaxStamina int
	Exp, MaxExp         int
	Level               int
	
	// 애니메이션을 위한 트윈 객체
	TweenX *Tween
	TweenY *Tween

	Color color.RGBA
}

func NewPlayer(q, r int) *Player {
	return &Player{
		Q: q, R: r,
		HP: 100, MaxHP: 100,
		MP: 50, MaxMP: 50,
		Gold: 100,
		Food: 10, MaxFood: 10,
		Stamina: 100, MaxStamina: 100,
		Exp: 0, MaxExp: 100,
		Level: 1,
		Color: color.RGBA{255, 200, 0, 255},
		TweenX: NewTween(startX, 0.1), // 속도 0.1 (숫자가 작을수록 묵직하게 이동)
		TweenY: NewTween(startY, 0.1),
	}
}

//
// (p *Player) ConsumeFood 턴 경과에 따라 식량을 소모하며, 식량이 부족할 경우 체력을 감소시킵니다.
// @param amount 소모할 식량 양
//
func (p *Player) ConsumeFood(amount int) {
	if p.Food > 0 {
		p.Food -= amount
	} else {
		// 식량이 없으면 체력 감소
		p.HP -= 5
		if p.HP < 0 { p.HP = 0 }
	}
}

// AddExp: 경험치 획득 및 레벨업 체크
func (p *Player) AddExp(amount int) {
	p.Exp += amount
	if p.Exp >= p.MaxExp {
		p.LevelUp()
	}
}

func (p *Player) LevelUp() {
	p.Level++
	p.Exp -= p.MaxExp
	p.MaxExp = int(float64(p.MaxExp) * 1.2) // 다음 레벨 경험치 증가
	p.HP = p.MaxHP                          // 레벨업 시 체력 회복
	p.MP = p.MaxMP
}

// [추가] 매 프레임 위치를 업데이트하는 함수
func (p *Player) UpdateAnimation() {
	p.TweenX.Update()
	p.TweenY.Update()
}

//
// (p *Player) Draw 플레이어를 화면에 렌더링합니다.
// screen 렌더링 타겟 이미지
// offsetX 카메라 X 오프셋
// offsetY 카메라 Y 오프셋
//
func (p *Player) Draw(screen *ebiten.Image, offsetX, offsetY float32) {
	// 헥사곤 좌표 -> 화면 좌표 변환 (map_system.go의 로직과 일치해야 함)
	spacingX := float32(HexRadius) * 1.5
	spacingY := float32(HexRadius) * 1.73205
	
	posX := float32(p.Q)*spacingX + offsetX
	posY := float32(p.R)*spacingY + offsetY
	if p.Q%2 != 0 {
		posY += spacingY / 2
	}

	// 플레이어 캐릭터 본체 (원형)
	vector.DrawFilledCircle(screen, posX, posY, float32(HexRadius)*0.6, p.Color, true)
	// 외곽선 추가
	vector.StrokeCircle(screen, posX, posY, float32(HexRadius)*0.6, 2, color.White, true)
}