package common

import (
	"fmt"
)

func TransformToFloat64(number interface{}) (float64, error) {

	var v float64
	switch i := number.(type) {
	case float64:
		v = i
	case float32:
		v = float64(i)
	case int64:
		v = float64(i)
	case int32:
		v = float64(i)
	case int:
		v = float64(i)
	default:
		return -1, fmt.Errorf("Transform value error")
	}
	return v, nil
}

type Int64arr []int64

func (a Int64arr) Len() int           { return len(a) }
func (a Int64arr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Int64arr) Less(i, j int) bool { return a[i] < a[j] }
