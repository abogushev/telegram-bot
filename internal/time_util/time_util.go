package time_util

import "time"

var DateTemplate = "02-01-2006"

func DatesEq(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func TimeToDate(t time.Time) string {
	return t.Format(DateTemplate)
}

func DateToTime(d string) (time.Time, error) {
	return time.Parse(DateTemplate, d)
}
