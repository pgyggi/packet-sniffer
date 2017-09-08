package common

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	end := time.Now()
	start := end.Add(-5 * time.Hour)
	CalDates(start.Unix(), end.Unix())
}
