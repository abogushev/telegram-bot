package model

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type Report struct {
	Start time.Time                  `json:"start"`
	End   time.Time                  `json:"end"`
	Data  map[string]decimal.Decimal `json:"data,omitempty"`
}

func NewReport(start time.Time, end time.Time, data map[string]decimal.Decimal) *Report {
	return &Report{Start: start, End: end, Data: data}
}
func FromJSON(data string) ([]Report, error) {
	result := make([]Report, 0)
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ToJSON(arr []Report) (string, error) {
	marshal, err := json.Marshal(&arr)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
