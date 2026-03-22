package main

import (
	"image/color"
	//"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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
	}
}

// Draw: 맵 오프셋을 반영하여 플레이어를 화면에 그립니다.
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

// ConsumeFood: 턴 경과에 따른 식량 소모 및 페널티 처리
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