package main

import "encoding/json"

type BloodPressurePoint struct {
	Id          int      `json:"id"`
	Systolic    float32  `json:"systolic"`
	Diastolic   float32  `json:"diastolic"`
	Timestamp   string   `json:"ts"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

func (b *BloodPressurePoint) Bytes() ([]byte, error) {
	j, err := json.Marshal(b)
	return j, err
}

type BloodPressureRange struct {
	Samples int                  `json:"samples"`
	Points  []BloodPressurePoint `json:"points"`
}

func (p *BloodPressureRange) Bytes() ([]byte, error) {
	j, e := json.Marshal(p)
	return j, e
}
