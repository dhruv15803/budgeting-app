package worker

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
)

type RecurringScheduler struct {
	cron   *cron.Cron
	run    func(ctx context.Context) (int, error)
	logger func(format string, args ...interface{})
}

// NewRecurringScheduler wires up a cron job that invokes `run` on the given schedule.
// `schedule` uses 5-field cron syntax (minute hour dom month dow). Runs in UTC.
func NewRecurringScheduler(schedule string, run func(ctx context.Context) (int, error), logger func(format string, args ...interface{})) (*RecurringScheduler, error) {
	c := cron.New(cron.WithLocation(time.UTC))
	s := &RecurringScheduler{cron: c, run: run, logger: logger}

	if _, err := c.AddFunc(schedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		n, err := run(ctx)
		if err != nil {
			if logger != nil {
				logger("recurring generator failed: %v", err)
			}
			return
		}
		if logger != nil {
			logger("recurring generator created %d expenses", n)
		}
	}); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *RecurringScheduler) Start() { s.cron.Start() }

func (s *RecurringScheduler) Stop() {
	<-s.cron.Stop().Done()
}

// RunNow executes the generator immediately, bypassing the cron schedule.
// Useful for startup catch-up and tests.
func (s *RecurringScheduler) RunNow(ctx context.Context) (int, error) {
	return s.run(ctx)
}
