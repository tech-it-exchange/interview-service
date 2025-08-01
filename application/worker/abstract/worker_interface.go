package abstract

import "context"

type WorkerInterface interface {
	Start(ctx context.Context)
}
