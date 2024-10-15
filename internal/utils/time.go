package utils

import (
	"strconv"
	"time"
)

func JudgeJoinTime(time time.Time) string {
	// 2022-01-01 00:00:00 - 2022-07-01 23:59:59 2022S
	// 2022-07-01 00:00:00 - 2022-08-31 23:59:59 2022C
	// 2022-09-01 00:00:00 - 2022-12-31 23:59:59 2022A
	year, month, _ := time.Date()
	var season string
	if month >= 1 && month <= 6 {
		season = "S"
	} else if month == 7 || month == 8 {
		season = "C"
	} else {
		season = "A"
	}
	return strconv.Itoa(year) + season
}
