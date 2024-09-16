package scheduler

import (
	"context"
	"github.com/go-co-op/gocron/v2"
	"ticket-system/booking/internal/service"
	"time"
)

type job struct {
	d gocron.JobDefinition
	t gocron.Task
	o []gocron.JobOption
}

func getJobs(ctx context.Context, s Scheduler) []job {
	var jobs []job
	jobs = append(jobs, removeExpiredBookingsJob(ctx, s))
	return jobs
}

func removeExpiredBookingsJob(ctx context.Context, scheduler Scheduler) job {
	f := func(bl service.BusinessLogic) error {
		return bl.RemoveExpiredBookings(ctx)
	}

	var j job
	j.d = gocron.DurationJob(time.Minute)
	j.t = gocron.NewTask(f, scheduler.bls)

	return j
}
