// file_manager.go
package main

import (
	"fmt"

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
		fmt.Println("[Config] 설정 파일이 없습니다. 기본값을 사용합니다.")
		return
	}

	data, err := ioutil.ReadFile(configFileName)
	if err != nil {
		fmt.Printf("[Config] 파일을 읽는 중 오류 발생: %v\n", err)
		return
	}

	var config ConfigData
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Printf("[Config] JSON 파싱 오류: %v\n", err)
		return
	}

	// 불러온 값을 전역 변수에 적용
	if config.ResolutionIndex >= 0 && config.ResolutionIndex < len(SupportedResolutions) {
		CurrentResIndex = config.ResolutionIndex
		ScreenWidth = SupportedResolutions[CurrentResIndex].Width
		ScreenHeight = SupportedResolutions[CurrentResIndex].Height
		
		// 성공 로그 출력
		fmt.Printf("[Config] 설정 로드 완료: %s (%dx%d)\n", 
			SupportedResolutions[CurrentResIndex].Name, ScreenWidth, ScreenHeight)
	} else {
		fmt.Println("[Config] 잘못된 해상도 인덱스입니다. 기본값으로 복구합니다.")
	}
}