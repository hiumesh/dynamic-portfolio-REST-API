package observability

import (
	"context"
	"sync"

	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
)

var (
	cleanupWaitGroup sync.WaitGroup
)

// WaitForCleanup waits until all observability long-running goroutines shut
// down cleanly or until the provided context signals done.
func WaitForCleanup(ctx context.Context) {
	utilities.WaitForCleanup(ctx, &cleanupWaitGroup)
}
