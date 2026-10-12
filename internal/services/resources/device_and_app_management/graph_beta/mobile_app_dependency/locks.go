package graphBetaMobileAppDependency

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

var dependencyAppLocks sync.Map

// lockAppRelationships holds the app lock through read-back so sibling resources see the completed
// write before editing the same collection. Waiting respects operation timeouts.
func lockAppRelationships(ctx context.Context, sourceID string) (func(), error) {
	value, _ := dependencyAppLocks.LoadOrStore(strings.ToLower(sourceID), make(chan struct{}, 1))
	lock := value.(chan struct{})
	select {
	case lock <- struct{}{}:
		return func() { <-lock }, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("wait for app relationship lock: %w", ctx.Err())
	}
}
