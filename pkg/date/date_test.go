package date

import "testing"

func TestGetTodayStartEnd(t *testing.T) {
	start, end := GetTodayStartEnd()
	t.Logf("start: %s", start)
	t.Logf("end: %s", end)
}
