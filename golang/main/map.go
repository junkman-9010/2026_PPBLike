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
	
	highlightPulse *Pulse
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

// isOffScreen은 주어진 좌표가 화면 가시 영역 밖에 있는지 확인합니다.
func (m *HexMap) isOffScreen(x, y, padding float32) bool {
	return x < -padding || x > float32(ScreenWidth)+padding ||
		y < -padding || y > float32(ScreenHeight)+padding
}

// Draw 메서드는 현재 화면에 보이는 모든 타일과 객체를
// 렌더링합니다. 화면 밖 타일은 그리지 않아 성능을 최적화합니다.
// 타일이 밝혀졌다면 기본 색상과 테두리를, 시야에 없으면 그림자를,
// reachableTiles 에 있는 타일은 노란색 테두리로, 선택된 타일은 빨간색 테두리로 표시합니다.
// 마지막으로 몬스터와 플레이어를 현재 시야에 따라 그립니다.
func (m *HexMap) Draw(screen *ebiten.Image, currentMode int) {
	spacingX := float32(HexRadius) * 1.5
	spacingY := float32(HexRadius) * 1.73205
	padding := float32(HexRadius) * 2

	mx, my := ebiten.CursorPosition()
    hoverQ, hoverR := m.ScreenToTile(float32(mx), float32(my))

    // 1단계: 모든 지형과 노란색(이동 가능) 테두리를 먼저 그립니다.
    for r := 0; r < m.height; r++ {
        for q := 0; q < m.width; q++ {
            posX, posY := m.getTileScreenPos(q, r, spacingX, spacingY)
            if m.isOffScreen(posX, posY, padding) { continue }

            key := fmt.Sprintf("%d,%d", q, r)
            if m.revealedTiles[key] {
                m.drawTile(screen, posX, posY, HexRadius, m.tiles[r][q])
                m.drawHexOutline(screen, posX, posY, HexRadius, color.RGBA{0, 0, 0, 50}, 1)

                // 노란색 테두리는 여기서 미리 다 그려둡니다.
                if currentMode == ModeNormal {
                    if _, reachable := m.reachableTiles[key]; reachable {
                        m.drawHexOutline(screen, posX, posY, HexRadius, color.RGBA{255, 255, 0, 150}, 3)
                    }
                }

                if !m.visibleTiles[key] {
                    m.drawHexShadow(screen, posX, posY, HexRadius, color.RGBA{0, 0, 0, 150})
                }
            }
        }
    }

    // 2단계: 모든 타일이 그려진 "위에" 빨간색 하이라이트만 단독으로 덧씌웁니다.
    // 마우스가 가리키는 타일(hoverQ, hoverR)이 유효한지 확인 후 그립니다.
    if hoverQ >= 0 && hoverQ < m.width && hoverR >= 0 && hoverR < m.height {
        key := fmt.Sprintf("%d,%d", hoverQ, hoverR)
        
        // 렌더링 조건 체크
        shouldDrawRed := false
        if currentMode == ModeView {
            shouldDrawRed = true // View 모드에선 무조건
        } else if _, reachable := m.reachableTiles[key]; reachable {
            shouldDrawRed = true // Normal 모드에선 이동 가능 지역일 때만
        }

        // 2. Draw 함수 내부
		if shouldDrawRed {
			m.highlightPulse.Update() // 애니메이션 값 업데이트
			posX, posY := m.getTileScreenPos(hoverQ, hoverR, spacingX, spacingY)
			
			// 기본 두께 5에 Pulse 값(±2)을 더해 3~7px 사이로 계속 변하게 함
			dynamicWidth := 5 + m.highlightPulse.Value
			m.drawHexOutline(screen, posX, posY, HexRadius, color.RGBA{255, 0, 0, 255}, dynamicWidth)
		}
    }
	
	// 4. 오브젝트 레이어
	for _, mon := range m.monsters {
		if m.visibleTiles[fmt.Sprintf("%d,%d", mon.Q, mon.R)] {
			mon.Draw(screen, m.offsetX, m.offsetY)
		}
	}
	
	m.player.Draw(screen, m.offsetX, m.offsetY)
}