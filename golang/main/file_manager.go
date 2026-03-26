// file_manager.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ConfigData는 파일에 저장될 설정 항목들을 정의합니다.
type ConfigData struct {
	ResolutionIndex int `json:"resolution_index"`
}

const configFileName = "config.json"

// SaveConfig: 현재 설정을 파일로 저장합니다.
func SaveConfig() error {
	config := ConfigData{
		ResolutionIndex: CurrentResIndex,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFileName, data, 0644)
}

// LoadConfig: 파일에서 설정을 불러와 전역 변수에 적용합니다.
func LoadConfig() {
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		// 파일이 없으면 기본값 유지
		return
	}

	data, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return
	}

	var config ConfigData
	if err := json.Unmarshal(data, &config); err != nil {
		return
	}

	// 불러온 값을 전역 변수에 적용
	if config.ResolutionIndex >= 0 && config.ResolutionIndex < len(SupportedResolutions) {
		CurrentResIndex = config.ResolutionIndex
		ScreenWidth = SupportedResolutions[CurrentResIndex].Width
		ScreenHeight = SupportedResolutions[CurrentResIndex].Height
	}
}