
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

var (
	whiteImage = ebiten.NewImage(3, 3)
	subImage   *ebiten.Image
)

func init() {
	whiteImage.Fill(color.White)
	subImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
}

// 지형별 이동/시야 비용 정의
type TerrainStats struct {
	MoveCost   int
	VisionCost int
}

type Tile struct {
	Q, R       int
	Terrain    TerrainType
	HasVillage bool
}

type HexMap struct {
	width, height int
	tiles         [][]Tile
	noise         opensimplex.Noise // [필수] map_system.go에서 사용

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

// Draw 로직은 기존과 동일하되 TerrainColors(대문자) 확인
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

// drawTile: 개별 육각형 타일을 그립니다.
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

func (m *HexMap) drawHexOutline(screen *ebiten.Image, x, y, radius float32, clr color.Color, width float32) {
	for i := 0; i < 6; i++ {
		a1 := float64(i) * 3.14159 / 3
		a2 := float64(i+1) * 3.14159 / 3
		vector.StrokeLine(screen, x+radius*float32(math.Cos(a1)), y+radius*float32(math.Sin(a1)), 
			x+radius*float32(math.Cos(a2)), y+radius*float32(math.Sin(a2)), width, clr, false)
	}
}

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