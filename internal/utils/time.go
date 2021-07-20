package utils

import "time"

// UnixNow returns current unix seconds timestamp as time.Time
func UnixNow() time.Time {
	return time.Unix(time.Now().Unix(), 0)
}

// UnixUTCNow returns current unix seconds timestamp as time.Time in UTC timezone
func UnixUTCNow() time.Time {
	return UnixNow().UTC()
}
