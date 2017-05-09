package model

import (
	"time"
)

// Now return the current time.Time with microsecond precision.
func Now() time.Time {
	now := time.Now()
	return time.Unix(now.Unix(), int64(now.Nanosecond())/1000*1000)
}
