package main

import "image/color"

// Resolution: 해상도 정보를 담는 구조체
type Resolution struct {
	Width  int
	Height int
	Name   string
}

// 지원하는 해상도 리스트
var SupportedResolutions = []Resolution{
	{1280, 720,  "1280x720 (16:9)"},
	{1600, 900,  "1600x900 (16:9)"},
	{1920, 1080, "1920x1080 (16:9)"},
}

// 전역 상태로 관리할 현재 해상도 (기본값: 1280x720)
var (
	CurrentResIndex = 0
	// 초기화를 위해 변수로 선언 (기존 상수를 대체)
	ScreenWidth  = SupportedResolutions[0].Width
	ScreenHeight = SupportedResolutions[0].Height
)

const (
	HexRadius = 60
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

// TerrainType 및 색상/스태츠 로직은 기존과 동일하게 유지
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

type TerrainStats struct {
	MoveCost   int
	VisionCost int
}

var terrainStats = map[TerrainType]TerrainStats{
	Ocean:      {99, 1},
	Plains:     {1, 1},
	DeepForest: {2, 2},
	Desert:     {1, 1},
	Mountain:   {99, 3},
}