package model

import (
	"encoding/json"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/time_util"
	"time"

	"github.com/shopspring/decimal"
)

type Report struct {
	UserId int64                      `json:"userId"`
	Start  time.Time                  `json:"start"`
	End    time.Time                  `json:"end"`
	Data   map[string]decimal.Decimal `json:"data,omitempty"`
}

func NewReport(userId int64, start time.Time, end time.Time, data map[string]decimal.Decimal) *Report {
	return &Report{UserId: userId, Start: start, End: end, Data: data}
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

type ReportRequest struct {
	UserId int64  `json:"userId"`
	Start  string `json:"start"`
	End    string `json:"end"`
}

func NewReportRequest(userId int64, start, end time.Time) *ReportRequest {
	return &ReportRequest{UserId: userId, Start: time_util.TimeToDate(start), End: time_util.TimeToDate(end)}
}
