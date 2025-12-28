package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SchedulerHandler struct {
	schedulerJobs []SchedulerJob
	jobAdded      chan SchedulerJob
	ctx           context.Context
	cancel        context.CancelFunc
	stop          chan struct{}
	wg            sync.WaitGroup
}

func NewSchedulerHandler() *SchedulerHandler {
	ctx, cancel := context.WithCancel(context.Background())
	return &SchedulerHandler{
		schedulerJobs: make([]SchedulerJob, 0),
		jobAdded:      make(chan SchedulerJob),
		ctx:           ctx,
		cancel:        cancel,
		stop:          make(chan struct{}),
		wg:            sync.WaitGroup{},
	}
}

func (s *SchedulerHandler) AddJob(job SchedulerJob) {
	s.schedulerJobs = append(s.schedulerJobs, job)
}

func (s *SchedulerHandler) Run() {
	s.wg.Add(len(s.schedulerJobs))

	fmt.Println("Run job with", len(s.schedulerJobs))
	for _, job := range s.schedulerJobs {
		fmt.Println("Init job which has wait time", job.WaitTime().String())
		go s.executeJob(job)
	}

	for {
		select {
		case addedJob := <-s.jobAdded:
			s.wg.Add(1)
			fmt.Println("Add job name", addedJob.GetName())
			go s.executeJob(addedJob)
		case <-s.stop:
			return
		}
	}
}

func (s *SchedulerHandler) executeJob(job SchedulerJob) {
	defer s.wg.Done()

	for {
		toExecuteTime := job.WaitTime()
		timer := time.NewTimer(toExecuteTime)
		select {
		case <-timer.C:
			job.Execute()
			timer.Stop()
		case <-s.ctx.Done():
			timer.Stop()
			return
		}

		if job.IsOneTime() {
			fmt.Println("Job done", job.GetName())
			return
		}
	}
}

func (s *SchedulerHandler) waitUntilJobDone(timeout time.Duration) {

	waitDone := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(waitDone)
	}()

	select {
	case <-time.After(timeout):
		fmt.Println("Cancel all running")
		s.cancel()
		return
	case <-waitDone:
		fmt.Println("All job done")
		return
	}
}

func (s *SchedulerHandler) Stop(timeout time.Duration) {
	close(s.stop)
	s.waitUntilJobDone(timeout)
}
