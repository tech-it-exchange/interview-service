package worker

import "context"

type WorkersInterface interface {
	StartWorkers(ctx context.Context)
}
