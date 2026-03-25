// animation.go
package main

import (
	"math"
)

// Tween은 시작값에서 목표값까지 부드럽게 수치를 변화시키는 구조체입니다.
type Tween struct {
	Current float32 // 현재 화면에 표시될 값
	Target  float32 // 최종 도달해야 할 목표 값
	Speed   float32 // 변화 속도 (0.05 ~ 0.2 권장)
}

// NewTween: 새로운 트윈 객체를 생성합니다.
func NewTween(initialValue float32, speed float32) *Tween {
	return &Tween{
		Current: initialValue,
		Target:  initialValue,
		Speed:   speed,
	}
}

// Update: 매 프레임마다 Current를 Target에 가깝게 이동시킵니다.
func (t *Tween) Update() {
	diff := t.Target - t.Current
	if math.Abs(float64(diff)) < 0.01 {
		t.Current = t.Target
		return
	}
	t.Current += diff * t.Speed
}

// SetTarget: 새로운 목표값을 설정합니다.
func (t *Tween) SetTarget(target float32) {
	t.Target = target
}

// ---------------------------------------------------------
// UI 애니메이션용: 시간에 따라 진동하는 값 (호버 효과 등)
// ---------------------------------------------------------

type Pulse struct {
	Value  float32
	timer  float64
	Speed  float64
	Range  float32
}

func NewPulse(speed float64, pulseRange float32) *Pulse {
	return &Pulse{Speed: speed, Range: pulseRange}
}

func (p *Pulse) Update() {
	p.timer += p.Speed
	// Sin 함수를 이용해 -Range ~ +Range 사이를 반복함
	p.Value = float32(math.Sin(p.timer)) * p.Range
}