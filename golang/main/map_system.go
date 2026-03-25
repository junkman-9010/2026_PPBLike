// map_system.go
//
// 헥사곤 기반 게임 맵의 핵심 로직을 구현한 파일입니다.
// 전체 게임 흐름(맵 생성, 시야 계산, 이동, 카메라 조작 등)을 담당하며,
// 현재 플레이어와 몬스터의 행동을 제어합니다.
//
// godoc 형식의 한글 주석을 사용하여 각 함수·구조체·필드에 대한 설명을 추가했습니다.
// 이 파일을 별도 패키지로 분리하고 싶다면
// `package map` 으로 선언하고 필요한 부분만 export(대문자)하면 됩니다.
// 현재는 `package main` 으로 작성해 두었습니다.
//
package main

import (
	"fmt"
	
	"math"
	"math/rand"
	
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/ojrac/opensimplex-go"
)


//
// NewHexMap
// 새로운 헥사곤 맵 인스턴스를 생성하고 초기화합니다.
// w 맵의 가로 크기(열 수)
// h 맵의 세로 크기(행 수)
// 초기화된 HexMap 포인터
//
func NewHexMap(w, h int) *HexMap {
	m := &HexMap{
		width:          w,
		height:         h,
		noise:          opensimplex.NewNormalized(rand.Int63()),
		player:         NewPlayer(0, 0),
		visibleTiles:   make(map[string]bool),
		revealedTiles:  make(map[string]bool),
		reachableTiles: make(map[string]int),
		// [추가] Pulse 객체 생성 (속도 0.1, 범위 2.0px)
		highlightPulse: NewPulse(0.1, 2.0),
	}
	m.generateTerrain()
	
	// [해결 코드] 플레이어를 육지로 이동시킵니다.
	m.findSafeSpawnPoint()
	
	m.SpawnMonsters(5)
	m.UpdateVision(3)
	m.CalculateReachable(2)
	m.CenterCameraOnPlayer()
	
	return m
}

//
// MovePlayerToSelected
// 현재 선택된 타일로 플레이어를 이동시키고 턴을 소모합니다.
// 이동 후 시야 업데이트 및 몬스터 AI를 실행합니다.
//
func (m *HexMap) MovePlayerToSelected() bool{
	key := fmt.Sprintf("%d,%d", m.selectedQ, m.selectedR)
	if _, ok := m.reachableTiles[key]; ok {
		m.player.Q, m.player.R = m.selectedQ, m.selectedR
		m.TurnCount++
		m.UpdateVision(3)
		m.CalculateReachable(2)
		m.CenterCameraOnPlayer()
		for _, mon := range m.monsters {
			mon.UpdateAI(m)
		}
	}
	
	// 이동 확정
    m.player.Q = m.selectedQ
    m.player.R = m.selectedR
	
	// 2. [추가] 이동할 목적지의 실제 화면(픽셀) 좌표 계산
    // HexRadius와 간격 계산식을 constants.go와 동일하게 맞춥니다.
    spacingX := float32(HexRadius) * 1.5
    spacingY := float32(HexRadius) * 1.73205
    
    targetX := float32(m.player.Q) * spacingX
    targetY := float32(m.player.R) * spacingY
    if m.player.Q%2 != 0 {
        targetY += spacingY / 2
    }

    // 3. [핵심] Tween의 목표치(Target)를 설정합니다.
    // 이제 p.TweenX.Update()가 호출될 때마다 Current가 이 Target으로 서서히 변합니다.
    m.player.TweenX.SetTarget(targetX)
    m.player.TweenY.SetTarget(targetY)

    // 선택 해제 및 시야 업데이트
    m.selectedQ = -1
    m.selectedR = -1

    return true
}

//
// SpawnMonsters
// 맵 내에 바다가 아닌 안전한 스폰 지점을 찾아 플레이어를 배치합니다.
//
func (m *HexMap) SpawnMonsters(count int) {
	for i := 0; i < count; i++ {
		mq, mr := rand.Intn(m.width), rand.Intn(m.height)
		if m.tiles[mr][mq].Terrain != Ocean {
			m.monsters = append(m.monsters, NewMonster(mq, mr))
		}
	}
}

func (m *HexMap) CanMoveMonster(mon *Monster, q, r int) bool {
	if q < 0 || q >= m.width || r < 0 || r >= m.height { return false }
	return m.tiles[r][q].Terrain != Ocean && m.tiles[r][q].Terrain != Mountain
}

func (m *HexMap) CalculateMonsterVision(mon *Monster) map[string]bool {
	v := make(map[string]bool)
	v[fmt.Sprintf("%d,%d", mon.Q, mon.R)] = true
	for _, n := range m.GetNeighbors(mon.Q, mon.R) {
		v[fmt.Sprintf("%d,%d", n[0], n[1])] = true
	}
	return v
}

//
// GetTimeContext
// 현재 턴 수에 따른 시간 단계(Step)와 시간(Hour)을 반환합니다.
// timeStep 시간 단계 (0:아침, 1:오후, 2:저녁, 3:밤), hour 24시간제 시간
//
func (m *HexMap) GetTimeContext() (int, int) { // 반환 타입을 (int, int)로 설정
    hour := (m.TurnCount / 6) % 24
    
    // 시간대에 따른 숫자(0, 1, 2, 3) 반환
    timeStep := 0
    if hour >= 6 && hour < 12 {
        timeStep = 0 // Morning
    } else if hour >= 12 && hour < 18 {
        timeStep = 1 // Afternoon
    } else if hour >= 18 && hour < 24 {
        timeStep = 2 // Evening
    } else {
        timeStep = 3 // Night
    }
    
    return timeStep, hour
}

// generateTiles 은 opensimplex 노이즈를 이용해 맵 타일을 생성합니다.
func (m *HexMap) generateTerrain() {
	m.tiles = make([][]Tile, m.height)
	for r := 0; r < m.height; r++ {
		m.tiles[r] = make([]Tile, m.width)
		for q := 0; q < m.width; q++ {
			// 노이즈 값 가져오기 (0.0 ~ 1.0 사이)
			nx := float64(q)/float64(m.width) - 0.5
			ny := float64(r)/float64(m.height) - 0.5
			val := m.noise.Eval2(nx*3, ny*3) // 배율을 3 정도로 줘서 지형을 큼직하게 만듦

			var terrain TerrainType
			
			// [중요] 지형 결정 임계값 조정
			if val < 0.3 {
				terrain = Ocean      // 바다 (낮은 곳)
			} else if val < 0.5 {
				terrain = Plains     // 평원
			} else if val < 0.7 {
				terrain = DeepForest // 숲 (평원과 산 사이)
			} else if val < 0.85 {
				terrain = Desert     // 사막
			} else {
				terrain = Mountain   // 산 (가장 높은 곳)
			}

			m.tiles[r][q] = Tile{
				Q:       q,
				R:       r,
				Terrain: terrain,
			}
		}
	}
}

// UpdateVision: 시야 비용을 고려한 다익스트라 가시 범위 계산
func (m *HexMap) UpdateVision(maxVision int) {
	m.visibleTiles = make(map[string]bool)
	type Item struct{ q, r, cost int }
	queue := []Item{{m.player.Q, m.player.R, 0}}
	
	startKey := fmt.Sprintf("%d,%d", m.player.Q, m.player.R)
	m.visibleTiles[startKey], m.revealedTiles[startKey] = true, true

	for len(queue) > 0 {
		curr := queue[0]; queue = queue[1:]
		for _, n := range m.GetNeighbors(curr.q, curr.r) {
			newCost := curr.cost + terrainStats[m.tiles[n[1]][n[0]].Terrain].VisionCost
			if newCost > maxVision { continue }

			key := fmt.Sprintf("%d,%d", n[0], n[1])
			if !m.visibleTiles[key] {
				m.visibleTiles[key], m.revealedTiles[key] = true, true
				queue = append(queue, Item{n[0], n[1], newCost})
			}
		}
	}
}

// CalculateReachable: 이동 코스트를 고려한 다익스트라 이동 범위 계산
func (m *HexMap) CalculateReachable(maxMove int) {
	m.reachableTiles = make(map[string]int)
	type Item struct{ q, r, cost int }
	queue := []Item{{m.player.Q, m.player.R, 0}}
	m.reachableTiles[fmt.Sprintf("%d,%d", m.player.Q, m.player.R)] = 0

	for len(queue) > 0 {
		curr := queue[0]; queue = queue[1:]
		for _, n := range m.GetNeighbors(curr.q, curr.r) {
			newCost := curr.cost + terrainStats[m.tiles[n[1]][n[0]].Terrain].MoveCost
			if newCost > maxMove { continue }

			key := fmt.Sprintf("%d,%d", n[0], n[1])
			if oldCost, ok := m.reachableTiles[key]; !ok || newCost < oldCost {
				m.reachableTiles[key] = newCost
				queue = append(queue, Item{n[0], n[1], newCost})
			}
		}
	}
}

func (m *HexMap) GetNeighbors(q, r int) [][2]int {
	var neighbors [][2]int
	offsets := [][2]int{{1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {0, 1}}
	if q%2 != 0 { offsets = [][2]int{{1, 1}, {1, 0}, {0, -1}, {-1, 0}, {-1, 1}, {0, 1}} }
	for _, o := range offsets {
		nq, nr := q+o[0], r+o[1]
		if nq >= 0 && nq < m.width && nr >= 0 && nr < m.height {
			neighbors = append(neighbors, [2]int{nq, nr})
		}
	}
	return neighbors
}


/**
// CenterCameraOnPlayer
// 카메라의 초점을 플레이어 캐릭터의 중앙으로 이동시킵니다.
 */
func (m *HexMap) CenterCameraOnPlayer() {
	spacingX := float32(HexRadius) * 1.5
	spacingY := float32(HexRadius) * 1.73205
	px := float32(m.player.Q) * spacingX
	py := float32(m.player.R) * spacingY
	if m.player.Q%2 != 0 { py += spacingY / 2 }
	m.offsetX = float32(ScreenWidth)/2 - px
	m.offsetY = float32(ScreenHeight)/2 - py
}

// GetTileScreenPos 은 헥사곤 좌표(q, r)를 화면 픽셀 좌표로 변환합니다.
// sX, sY: 헥사곤 간격( HexRadius 기반)
// 반환값: pX, pY
func (m *HexMap) getTileScreenPos(q, r int, sX, sY float32) (float32, float32) {
	pX := float32(q)*sX + m.offsetX
	pY := float32(r)*sY + m.offsetY
	if q%2 != 0 { pY += sY / 2 }
	return pX, pY
}

// isOutsideScreen 은 주어진 좌표가 화면 밖에 있는지 여부를 판단합니다.
// x, y: 화면 좌표, padding: 허용 여백
func (m *HexMap) isOutsideScreen(x, y, padding float32) bool {
	return x < -padding || x > float32(ScreenWidth)+padding || y < -padding || y > float32(ScreenHeight)+padding
}

// UpdateCamera 은 마우스 입력에 따라 카메라를 이동시킵니다.
// 왼쪽 버튼이 눌리면 드래그를 시작하고, 놓으면 종료합니다.
func (m *HexMap) UpdateCamera() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if !m.isDragging {
			m.isDragging = true
			m.lastMouseX, m.lastMouseY = mx, my
		} else {
			m.offsetX += float32(mx - m.lastMouseX)
			m.offsetY += float32(my - m.lastMouseY)
			m.lastMouseX, m.lastMouseY = mx, my
		}
	} else { m.isDragging = false }
}

// ScreenToTile 은 화면 픽셀 좌표를 헥사곤 그리드 좌표(Q, R)로 변환합니다.
// x: 화면 X 좌표, y: 화면 Y 좌표
// 반환값: q(열), r(행)
// 범위를 벗어나면 (-1, -1) 반환
func (m *HexMap) ScreenToTile(x, y float32) (int, int) {
	// 카메라 오프셋 보정
	worldX := x - m.offsetX
	worldY := y - m.offsetY

	// 헥사곤 간격 계산 (constants.go의 HexRadius 기준)
	spacingX := float32(HexRadius) * 1.5
	spacingY := float32(HexRadius) * 1.73205

	// 대략적인 그리드 위치 계산
	q := int(math.Round(float64(worldX / spacingX)))
	
	adjY := worldY
	if q%2 != 0 {
		adjY -= spacingY / 2
	}
	r := int(math.Round(float64(adjY / spacingY)))

	// 맵 범위 체크
	if q < 0 || q >= m.width || r < 0 || r >= m.height {
		return -1, -1
	}

	return q, r
}

// findSafeSpawnPoint 은 맵의 중앙부를 기준으로 가장 먼저 발견되는
// 바다(Ocean)와 산(Mountain) 이외의 육지에 플레이어를 배치합니다.
// 플레이어가 이미 존재하는 경우 위치를 갱신합니다.
func (m *HexMap) findSafeSpawnPoint() {
	// 맵의 중앙부부터 탐색하여 가장 먼저 발견되는 육지(바다가 아닌 곳)에 배치
	for r := 0; r < m.height; r++ {
		for q := 0; q < m.width; q++ {
			t := m.tiles[r][q].Terrain
			if t != Ocean && t != Mountain { // 바다와 산이 아닌 곳(평원, 숲 등)을 찾음
				m.player.Q = q
				m.player.R = r
				return // 찾으면 즉시 종료
			}
		}
	}
}