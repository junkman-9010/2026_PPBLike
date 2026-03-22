package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// MonsterState: 몬스터의 현재 행동 상태를 정의합니다.
type MonsterState int

const (
	Idle    MonsterState = iota // 대기 및 랜덤 이동
	Chasing                     // 플레이어 추격
)

// MonsterStats: 몬스터의 능력치를 관리합니다.
type MonsterStats struct {
	MaxVision int
	MaxMove   int
	HP        int
	MaxHP     int
}

// Monster: 게임 내 적 캐릭터 구조체입니다.
type Monster struct {
	Q, R   int
	Stats  MonsterStats
	State  MonsterState
	Alive  bool
}

// NewMonster: 새로운 몬스터 인스턴스를 생성합니다.
func NewMonster(q, r int) *Monster {
	return &Monster{
		Q:     q,
		R:     r,
		Alive: true,
		State: Idle,
		Stats: MonsterStats{
			MaxVision: 3, 
			MaxMove:   2,
			HP:        20,
			MaxHP:     20,
		},
	}
}

// UpdateAI: 몬스터의 AI 로직을 업데이트합니다. (매 턴 호출)
func (m *Monster) UpdateAI(mMap *HexMap) {
	if !m.Alive {
		return
	}

	// 1. 시야 체크: 플레이어가 몬스터의 시야 범위 내에 있는지 확인
	playerKey := fmt.Sprintf("%d,%d", mMap.player.Q, mMap.player.R)
	visibleTiles := mMap.CalculateMonsterVision(m) // map.go 또는 map_system.go에 정의된 함수 사용

	if _, canSee := visibleTiles[playerKey]; canSee {
		m.State = Chasing
	} else {
		m.State = Idle
	}

	// 2. 이동 경로 결정
	var targetQ, targetR int
	if m.State == Chasing {
		// 추격 모드: 플레이어와 가장 가까워지는 인접 타일 탐색
		targetQ, targetR = m.getNextStepTowards(mMap, mMap.player.Q, mMap.player.R)
	} else {
		// 대기 모드: 인접한 타일 중 하나로 랜덤 이동
		neighbors := mMap.GetNeighbors(m.Q, m.R)
		if len(neighbors) > 0 {
			idx := rand.Intn(len(neighbors))
			targetQ, targetR = neighbors[idx][0], neighbors[idx][1]
		}
	}

	// 3. 지형 비용 및 이동 가능 여부 확인 후 좌표 업데이트
	if mMap.CanMoveMonster(m, targetQ, targetR) {
		m.Q, m.R = targetQ, targetR
	}
}

// getNextStepTowards: 목표 좌표(플레이어)로 향하는 최적의 인접 타일을 반환합니다.
func (m *Monster) getNextStepTowards(mMap *HexMap, pq, pr int) (int, int) {
	bestQ, bestR := m.Q, m.R
	// 단순 맨해튼 거리를 이용한 거리 계산
	minDist := math.Abs(float64(m.Q-pq)) + math.Abs(float64(m.R-pr))

	for _, n := range mMap.GetNeighbors(m.Q, m.R) {
		// 이동 불가능한 지형(바다, 산)은 AI 경로에서 제외
		tile := mMap.tiles[n[1]][n[0]]
		if tile.Terrain == Ocean || tile.Terrain == Mountain {
			continue
		}
		
		dist := math.Abs(float64(n[0]-pq)) + math.Abs(float64(n[1]-pr))
		if dist < minDist {
			minDist = dist
			bestQ, bestR = n[0], n[1]
		}
	}
	return bestQ, bestR
}

// Draw: 몬스터를 화면에 렌더링합니다.
func (m *Monster) Draw(screen *ebiten.Image, offsetX, offsetY float32) {
	if !m.Alive {
		return
	}

	// 헥사곤 좌표를 화면 좌표로 변환
	spacingX := float32(HexRadius) * 1.5
	spacingY := float32(HexRadius) * 1.73205
	posX := float32(m.Q)*spacingX + offsetX
	posY := float32(m.R)*spacingY + offsetY
	if m.Q%2 != 0 {
		posY += spacingY / 2
	}

	// 상태에 따른 색상 변화 (일반: 어두운 빨강, 추격: 밝은 빨강)
	clr := color.RGBA{180, 40, 40, 255}
	if m.State == Chasing {
		clr = color.RGBA{255, 0, 0, 255}
	}

	// 몬스터 본체 그리기 (사각형 형태)
	size := float32(16)
	vector.DrawFilledRect(screen, posX-(size/2), posY-(size/2), size, size, clr, true)
	
	// 체력바 표시 (옵션: 머리 위에 작은 바 추가)
	/*
	if m.Stats.HP < m.Stats.MaxHP {
		u.drawMonsterHPBar(screen, posX, posY, m)
	}
	*/
}

// drawMonsterHPBar: 몬스터 머리 위에 간이 체력바를 그립니다.
/*
func (u *UI) drawMonsterHPBar(screen *ebiten.Image, x, y float32, m *Monster) {
	barW := float32(20)
	barH := float32(3)
	bx := x - (barW / 2)
	by := y - 15

	// 배경 (검정)
	vector.DrawFilledRect(screen, bx, by, barW, barH, color.RGBA{0, 0, 0, 255}, true)
	// 잔여 체력 (초록)
	ratio := float32(m.Stats.HP) / float32(m.Stats.MaxHP)
	vector.DrawFilledRect(screen, bx, by, barW*ratio, barH, color.RGBA{50, 255, 50, 255}, true)
}
*/