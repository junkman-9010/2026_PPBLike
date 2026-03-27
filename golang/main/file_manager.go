// file_manager.go
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// SaveJSON: 데이터를 JSON 파일로 저장합니다.
func SaveJSON(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, jsonData, 0644)
}

// LoadJSON: JSON 파일을 읽어 구조체에 채웁니다.
func LoadJSON(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// FileExists: 파일 존재 여부를 확인합니다.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
