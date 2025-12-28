package scheduler

import (
	"fmt"
	"time"
)

const (
	REMINDER_JOB_TYPE = 0
)

type SchedulerJob interface {
	GetName() string
	IsOneTime() bool
	Execute()
	WaitTime() time.Duration
}

type job struct {
	name        string
	hour        int
	minute      int
	executeFunc func()
	repeatly    bool
}

type OneTimeJob struct {
	*job
}

type DailyJob struct {
	*job
}

type MonthlyJob struct {
	*job
	date int
}

func (j *job) WaitTime() time.Duration {
	refTime := time.Now()
	if refTime.Hour() > j.hour || (refTime.Hour() <= j.hour && refTime.Minute() > j.minute) {
		refTime = refTime.Add(time.Second * OneDay_sec)
	}
	targetTime := getNextTargetTime(refTime, j.hour, j.minute)
	return time.Until(targetTime)
}

func (j *job) GetName() string {
	toExecuteTime := j.WaitTime()
	return fmt.Sprintf("Job: %s will run in %v", j.name, toExecuteTime.String())
}

func (j *job) IsOneTime() bool {
	return j.repeatly
}

func (j *job) Execute() {
	j.executeFunc()
}

func NewOneTimeJob(name string, hour int, minute int, callback func()) *OneTimeJob {
	return &OneTimeJob{
		job: &job{
			name:        name,
			hour:        hour,
			minute:      minute,
			repeatly:    false,
			executeFunc: callback,
		},
	}
}

func NewDailyJob(name string, hour int, minute int, callback func()) *DailyJob {
	return &DailyJob{
		job: &job{
			name:        name,
			hour:        hour,
			minute:      minute,
			repeatly:    true,
			executeFunc: callback,
		},
	}
}

func NewMonthlyJob(name string, date int, hour int, minute int, callback func()) *MonthlyJob {
	return &MonthlyJob{
		job: &job{
			name:        name,
			hour:        hour,
			minute:      minute,
			repeatly:    true,
			executeFunc: callback,
		},
		date: date,
	}
}

func (d *DailyJob) WaitTime() time.Duration {
	refTime := time.Now()
	if refTime.Hour() > d.hour || (refTime.Hour() <= d.hour && refTime.Minute() > d.minute) {
		refTime = refTime.Add(time.Second * OneDay_sec)
	}
	targetTime := getNextTargetTime(refTime, d.hour, d.minute)
	return time.Until(targetTime)
}

func (d *MonthlyJob) WaitTime() time.Duration {
	refTime := time.Now()
	if refTime.Hour() > d.hour || (refTime.Hour() <= d.hour && refTime.Minute() > d.minute) {
		refTime = refTime.Add(time.Second * OneDay_sec)
	}
	targetTime := getNextTargetTime(refTime, d.hour, d.minute)
	return time.Until(targetTime)
}

func (d *OneTimeJob) WaitTime() time.Duration {
	refTime := time.Now()
	if refTime.Hour() > d.hour || (refTime.Hour() <= d.hour && refTime.Minute() > d.minute) {
		refTime = refTime.Add(time.Second * OneDay_sec)
	}
	targetTime := getNextTargetTime(refTime, d.hour, d.minute)
	return time.Until(targetTime)
}
