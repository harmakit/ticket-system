package scheduler

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"ticket-system/booking/internal/service"
)

type Scheduler struct {
	gcs gocron.Scheduler
	bls service.BusinessLogic
}

func Init(ctx context.Context, bls service.BusinessLogic) (Scheduler, error) {
	var s Scheduler
	var err error

	gcs, err := gocron.NewScheduler()
	if err != nil {
		return s, err
	}

	s.gcs = gcs
	s.bls = bls

	err = initJobs(ctx, s)
	if err != nil {
		return s, err
	}

	s.gcs.Start()

	return s, nil
}

func (s Scheduler) Start() {
	s.gcs.Start()
}

func (s Scheduler) Shutdown() error {
	return s.gcs.Shutdown()
}

func initJobs(ctx context.Context, s Scheduler) error {
	for _, job := range getJobs(ctx, s) {
		_, err := s.gcs.NewJob(job.d, job.t, job.o...)
		if err != nil {
			return err
		}
	}
	return nil
}
