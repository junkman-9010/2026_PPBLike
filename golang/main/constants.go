package main

import "image/color"

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	HexRadius    = 80
)

type Scene int
const (
	SceneMenu Scene = iota
	SceneGame
	SceneOption
)

const (
	ModeNormal = iota
	ModeView
)

// map.go에서 이쪽으로 이동하여 통합 관리
type TerrainType int
const (
	Ocean TerrainType = iota
	Plains
	DeepForest
	Desert
	Mountain
)

var (
	TerrainColors = map[TerrainType]color.RGBA{
		Ocean:      {10, 50, 100, 255},
		Plains:     {100, 150, 50, 255},
		DeepForest: {30, 80, 30, 255},
		Desert:     {220, 190, 100, 255},
		Mountain:   {100, 90, 80, 255},
	}
	VillageColor = color.RGBA{180, 100, 50, 255}
)

var terrainStats = map[TerrainType]TerrainStats{
	Ocean:      {99, 1},
	Plains:     {1, 1},
	DeepForest: {2, 2},
	Desert:     {1, 1},
	Mountain:   {99, 2},
}