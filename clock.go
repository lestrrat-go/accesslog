package accesslog

import "time"

type Clock interface {
	Now() time.Time
}

// StaticClock is a Clock that always returns the same time. It's only used for
// testing.
type StaticClock time.Time

func (c StaticClock) Now() time.Time {
	return time.Time(c)
}

// SystemClock is a wrapper around time.Now().
type SystemClock struct{}

func (SystemClock) Now() time.Time {
	return time.Now()
}
