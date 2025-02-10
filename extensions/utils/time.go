package utils

import "time"

func TimeSince(t time.Time) float64 {
	return float64(time.Since(t)) / float64(time.Second)
}
