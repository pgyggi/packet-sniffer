package common

import (
	"time"
)

func CalDates(indexprefix string, fromSec int64, toSec int64) []string {

	var days []string
	from := time.Unix(fromSec, 0)
	to := time.Unix(toSec, 0)

	fromStr := from.Format("2006-01-02")
	toStr := to.Format("2006-01-02")
	days = append(days, indexprefix+fromStr)

	if fromStr == toStr {
		return days
	} else {

		if from.After(to) {
			return days
		}

		i := 1
		for {
			t := 24 * i
			tmp := from.Add(time.Duration(t) * time.Hour)
			prefix := tmp.Format("2006-01-02")
			days = append(days, indexprefix+prefix)

			if prefix == toStr {
				break
			}
			i++
		}
	}

	return days
}
