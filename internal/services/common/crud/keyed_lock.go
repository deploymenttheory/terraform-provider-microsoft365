package crud

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// keyedLockWaitLogThreshold is how long a writer waits before the wait is logged.
const keyedLockWaitLogThreshold = time.Second

// keyedLocks maps scope + lower-cased key to a one-slot channel that acts as the lock.
var keyedLocks sync.Map

// AcquireKeyedLock serialises writers that share a Graph collection, such as an app's
// relationships or a principal's license set. Writers with the same scope and key run
// one at a time; keys are compared case-insensitively, and different scopes or keys never
// block each other.
//
// It waits until the lock is free or ctx is done, so a writer never waits past its
// deadline or after cancellation. A context that is already done returns an error
// without taking the lock. Waits longer than one second are logged.
//
// On success it returns an unlock func that is safe to call more than once. On failure it
// returns a nil unlock func and an error that wraps ctx.Err().
func AcquireKeyedLock(ctx context.Context, scope, key string) (func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("wait for %s lock on %q: %w", scope, key, err)
	}

	value, _ := keyedLocks.LoadOrStore(scope+"\x00"+strings.ToLower(key), make(chan struct{}, 1))
	lock := value.(chan struct{})

	start := time.Now()
	select {
	case lock <- struct{}{}:
	case <-ctx.Done():
		tflog.Debug(ctx, fmt.Sprintf("Gave up waiting %s for the %s lock on %s", time.Since(start), scope, key))
		return nil, fmt.Errorf("wait for %s lock on %q: %w", scope, key, ctx.Err())
	}

	if waited := time.Since(start); waited > keyedLockWaitLogThreshold {
		tflog.Debug(ctx, fmt.Sprintf("Waited %s for the %s lock on %s", waited, scope, key))
	}

	var once sync.Once
	return func() { once.Do(func() { <-lock }) }, nil
}
