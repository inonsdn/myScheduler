package scheduler

import "time"

const (
	OneDay_sec = 24 * 60 * 60
)

func getTargetTime(hour int, minute int) time.Time {
	now := time.Now()
	targetTime := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		hour,
		minute,
		0,              // seconds
		0,              // nanoseconds
		now.Location(), // use the current local timezone
	)
	return targetTime
}

func getNextTargetTime(refTime time.Time, hour int, minute int) time.Time {
	targetTime := time.Date(
		refTime.Year(),
		refTime.Month(),
		refTime.Day(),
		hour,
		minute,
		0,                  // seconds
		0,                  // nanoseconds
		refTime.Location(), // use the current local timezone
	)
	return targetTime
}
