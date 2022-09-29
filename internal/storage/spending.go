package storage

import (
	"time"
)

type Category int
type ReportType int

const (
	Food Category = iota
	Other
)

func (c Category) String() string {
	r := ""
	switch c {
	case Food:
		r = "food"
	case Other:
		r = "other"
	}
	return r
}

const (
	Week ReportType = iota
	Month
	Year
)

var allSpendigs = make(map[Category][]*Model)

type Model struct {
	Value    float32
	Category Category
	Date     time.Time
}

func New(val float32, cat Category) *Model {
	return &Model{Value: val, Category: cat, Date: time.Now()}
}

func (m *Model) Save() {
	allSpendigs[m.Category] = append(allSpendigs[m.Category], m)
}

func groupBy(startAt time.Time, endAt time.Time) map[Category]float32 {
	result := make(map[Category]float32)
	for cat, ms := range allSpendigs {
		for i := 0; i < len(ms); i++ {
			if (ms[i].Date.After(startAt) && ms[i].Date.Before(endAt)) ||
				(ms[i].Date.Equal(startAt) || ms[i].Date.Equal(endAt)) {
				result[cat] += ms[i].Value
			}
		}
	}
	return result
}

func GetStatsBy(rt ReportType) (time.Time, time.Time, map[Category]float32) {
	endAt := time.Now()
	var startAt time.Time
	switch rt {
	case Week:
		startAt = endAt.AddDate(0, 0, -7)
	case Month:
		startAt = endAt.AddDate(0, -1, 0)
	case Year:
		startAt = endAt.AddDate(-1, 0, 0)
	}
	return startAt, endAt, groupBy(startAt, endAt)
}
