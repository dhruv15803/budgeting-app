package worker

import "context"

type VerificationJob struct {
	ToEmail         string
	VerificationURL string
}

type Queue struct {
	ch chan VerificationJob
}

func NewQueue(buffer int) *Queue {
	return &Queue{ch: make(chan VerificationJob, buffer)}
}

func (q *Queue) Submit(job VerificationJob) bool {
	select {
	case q.ch <- job:
		return true
	default:
		return false
	}
}

func RunVerificationMailWorker(ctx context.Context, q *Queue, send func(to string, verificationURL string) error, log func(format string, args ...interface{})) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-q.ch:
				if !ok {
					return
				}
				if err := send(job.ToEmail, job.VerificationURL); err != nil && log != nil {
					log("verification email send failed: %v", err)
				}
			}
		}
	}()
}
