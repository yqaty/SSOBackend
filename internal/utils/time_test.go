package utils

import (
	"testing"
	"time"
)

func TestJudgeJoinTime(t *testing.T) {
	timeStr1 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	timeStr2 := time.Date(2023, 9, 10, 0, 0, 0, 0, time.UTC)
	timeStr3 := time.Date(2021, 8, 31, 0, 0, 0, 0, time.UTC)

	t.Log("time =", JudgeJoinTime(timeStr1))
	t.Log("time =", JudgeJoinTime(timeStr2))
	t.Log("time =", JudgeJoinTime(timeStr3))
}
