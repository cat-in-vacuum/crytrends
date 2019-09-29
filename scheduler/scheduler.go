package scheduler

import "github.com/jasonlvhit/gocron"

type Scheduler struct {
	Scheduler *gocron.Scheduler
}

type Job struct {
	Name string
	gocron.Job
}

func New() *Scheduler {
	return &Scheduler{
		Scheduler: gocron.NewScheduler(),
	}
}
