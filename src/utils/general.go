package utils

import (
	"time"
)

func ToRef[T any](val T) *T {
	return &val
}

type Interval struct {
	Stop func()
}

func NewInterval(f func(), intervalMs int) *Interval {
	ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)
	go func() {
		for range ticker.C {
			f()
		}
	}()
	return &Interval{
		Stop: func() {
			ticker.Stop()
		},
	}
}

func NewTimeout(f func(), timeoutMs int) *time.Timer {
	return time.AfterFunc(time.Duration(timeoutMs)*time.Millisecond, f)
}
