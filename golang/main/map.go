// Package main은 Pixel Princess Blitz 게임의 핵심 맵 시스템을 담당합니다.
// map.go에 정의된 구조체와 함수들은 헥사곤 지도를 그리며, 게임 내 시야·이동·타일 상태를 관리합니다.
package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ojrac/opensimplex-go"
)

// whiteImage는 3x3 픽셀 단위로 만들어진 흰색 이미지이며,
// 이후 SubImage 를 통해 1x1 픽셀 단위로 자른 이미지(subImage)와 함께
// 각 타일을 그리는 데 사용됩니다.
var (
	whiteImage = ebiten.NewImage(3, 3)
	subImage   *ebiten.Image
)

// init 함수는 전역 이미지 객체를 초기화합니다.
// whiteImage 를 흰색으로 채우고, 그 안에서 1x1 영역을 잘라 subImage 로 저장합니다.
func init() {
	whiteImage.Fill(color.White)
	subImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
}

// TerrainStats는 각 지형 타입이 갖는 이동 비용과 시야 비용을 정의합니다.
type TerrainStats struct {
	MoveCost   int
	VisionCost int
}

// Tile은 헥사곤 맵 상의 한 타일을 나타냅니다.
// Q, R 좌표는 엑스포넌트(좌표계) 기반이며, Terrain 은 지형 타입,
// HasVillage 은 해당 타일에 마을이 있는지를 표시합니다.
type Tile struct {
	Q, R       int
	Terrain    TerrainType
	HasVillage bool
}

// HexMap은 전체 맵을 나타내는 구조체이며, 타일 배열,
// 노이즈 생성기, 플레이어·몬스터 목록, 시야·이동 정보 등을 포함합니다.
type HexMap struct {
	width, height int
	tiles         [][]Tile
	noise         opensimplex.Noise // map_system.go에서 사용

	player   *Player
	monsters []*Monster

	reachableTiles map[string]int
	visibleTiles   map[string]bool
	revealedTiles  map[string]bool

	selectedQ, selectedR   int
	offsetX, offsetY       float32
	isDragging             bool
	lastMouseX, lastMouseY int
	TurnCount              int
}

// drawTile은 개별 육각형 타일을 화면에 그립니다.
// vector 패키지의 Path 를 이용해 6개의 꼭지점으로 육각형을 만들고,
// 해당 타일의 지형 색상을 채우며, 마을이 있으면 VillageColor 로 덮어씁니다.
func (m *HexMap) drawTile(screen *ebiten.Image, x, y, radius float32, tile Tile) {
	var path vector.Path
	for i := 0; i < 6; i++ {
		angle := float64(i) * 3.14159 / 3
		vx := x + radius*float32(math.Cos(angle))
		vy := y + radius*float32(math.Sin(angle))
		if i == 0 { path.MoveTo(vx, vy) } else { path.LineTo(vx, vy) }
	}
	path.Close()

	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	clr := TerrainColors[tile.Terrain]
	if tile.HasVillage { clr = VillageColor }

	rf, gf, bf, af := float32(clr.R)/255, float32(clr.G)/255, float32(clr.B)/255, float32(clr.A)/255
	for i := range vs {
		vs[i].SrcX, vs[i].SrcY = 1, 1
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = rf, gf, bf, af
	}
	screen.DrawTriangles(vs, is, subImage, &ebiten.DrawTrianglesOptions{})
}

// drawHexOutline는 육각형 타일 주변에 테두리를 그립니다.
// 각 변마다 vector.StrokeLine 를 사용해 색상·두께를 지정합니다.
func (m *HexMap) drawHexOutline(screen *ebiten.Image, x, y, radius float32, clr color.Color, width float32) {
	for i := 0; i < 6; i++ {
		a1 := float64(i) * 3.14159 / 3
		a2 := float64(i+1) * 3.14159 / 3
		vector.StrokeLine(screen, x+radius*float32(math.Cos(a1)), y+radius*float32(math.Sin(a1)), 
			x+radius*float32(math.Cos(a2)), y+radius*float32(math.Sin(a2)), width, clr, false)
	}
}

// drawHexShadow는 현재 시야에 없는 영역에 어두운 그림자를 적용합니다.
func (m *HexMap) drawHexShadow(screen *ebiten.Image, x, y, radius float32, clr color.Color) {
	var path vector.Path
	for i := 0; i < 6; i++ {
		angle := float64(i) * 3.14159 / 3
		if i == 0 { path.MoveTo(x+radius*float32(math.Cos(angle)), y+radius*float32(math.Sin(angle)))
		} else { path.LineTo(x+radius*float32(math.Cos(angle)), y+radius*float32(math.Sin(angle))) }
	}
	path.Close()
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	r, g, b, a := clr.RGBA()
	rf, gf, bf, af := float32(r>>8)/255, float32(g>>8)/255, float32(b>>8)/255, float32(a>>8)/255
	for i := range vs {
		vs[i].SrcX, vs[i].SrcY = 1, 1
		vs[i].ColorR, vs[i].ColorG, vs[i].ColorB, vs[i].ColorA = rf, gf, bf, af
	}
	screen.DrawTriangles(vs, is, subImage, &ebiten.DrawTrianglesOptions{})
}


// Draw 메서드는 현재 화면에 보이는 모든 타일과 객체를
// 렌더링합니다. 화면 밖 타일은 그리지 않아 성능을 최적화합니다.
// 타일이 밝혀졌다면 기본 색상과 테두리를, 시야에 없으면 그림자를,
// reachableTiles 에 있는 타일은 노란색 테두리로, 선택된 타일은 빨간색 테두리로 표시합니다.
// 마지막으로 몬스터와 플레이어를 현재 시야에 따라 그립니다.
func (m *HexMap) Draw(screen *ebiten.Image) {
	spacingX := float32(HexRadius) * 1.5
	spacingY := float32(HexRadius) * 1.73205
	padding := float32(HexRadius)

	// 1. 지형 및 테두리 그리기
	for r := 0; r < m.height; r++ {
		for q := 0; q < m.width; q++ {
			posX, posY := m.getTileScreenPos(q, r, spacingX, spacingY)
			
			// 화면 밖 타일은 그리지 않음 (최적화)
			if posX < -padding || posX > float32(ScreenWidth)+padding || 
			   posY < -padding || posY > float32(ScreenHeight)+padding {
				continue
			}

			key := fmt.Sprintf("%d,%d", q, r)
			// 한 번이라도 밝혀진 타일만 그림
			if m.revealedTiles[key] {
				m.drawTile(screen, posX, posY, HexRadius, m.tiles[r][q])
				
				// [핵심] 테두리가 사라졌다면 이 부분이 확실히 있어야 합니다.
				// 타일 기본 테두리 (검은색, 두께 1)
				m.drawHexOutline(screen, posX, posY, HexRadius, color.RGBA{0, 0, 0, 50}, 1)

				// 현재 시야에 없는 곳은 어둡게 처리
				if !m.visibleTiles[key] {
					m.drawHexShadow(screen, posX, posY, HexRadius, color.RGBA{0, 0, 0, 150})
				}
			}
		}
	}

	// 2. 이동 가능 범위 하이라이트 (노란색 테두리)
	for key := range m.reachableTiles {
		var q, r int
		fmt.Sscanf(key, "%d,%d", &q, &r)
		posX, posY := m.getTileScreenPos(q, r, spacingX, spacingY)
		m.drawHexOutline(screen, posX, posY, HexRadius, color.RGBA{255, 255, 0, 255}, 2)
	}

	// 3. 선택된 타일 표시 (빨간색 테두리)
	selX, selY := m.getTileScreenPos(m.selectedQ, m.selectedR, spacingX, spacingY)
	m.drawHexOutline(screen, selX, selY, HexRadius, color.RGBA{255, 50, 50, 255}, 3)
	
	// 4. 오브젝트 레이어
	for _, mon := range m.monsters {
		if m.visibleTiles[fmt.Sprintf("%d,%d", mon.Q, mon.R)] {
			mon.Draw(screen, m.offsetX, m.offsetY)
		}
	}
	m.player.Draw(screen, m.offsetX, m.offsetY)
}