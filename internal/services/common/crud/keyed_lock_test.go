package crud

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflogtest"
)

func TestAcquireKeyedLock_ContentionSerialisesWriters(t *testing.T) {
	const writers = 50
	key := "aabbccdd-1111-1111-1111-111111111111"
	var active, maxActive int32
	var count int
	var wg sync.WaitGroup
	for index := range writers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			writerKey := key
			if index%2 == 0 {
				writerKey = strings.ToUpper(key)
			}
			unlock, err := AcquireKeyedLock(context.Background(), "test_contention", writerKey)
			if err != nil {
				t.Error(err)
				return
			}
			defer unlock()
			now := atomic.AddInt32(&active, 1)
			for {
				seen := atomic.LoadInt32(&maxActive)
				if now <= seen || atomic.CompareAndSwapInt32(&maxActive, seen, now) {
					break
				}
			}
			count++ // unsynchronised on purpose: the race detector fails the test if writers overlap
			time.Sleep(time.Millisecond)
			atomic.AddInt32(&active, -1)
		}()
	}
	wg.Wait()
	if count != writers {
		t.Fatalf("completed writes = %d, want %d", count, writers)
	}
	if maxActive != 1 {
		t.Fatalf("max concurrent writers = %d, want 1", maxActive)
	}
}

func TestAcquireKeyedLock_IndependentScopesAndKeysDoNotBlock(t *testing.T) {
	unlock, err := AcquireKeyedLock(context.Background(), "test_independent", "key-one")
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	otherKey, err := AcquireKeyedLock(ctx, "test_independent", "key-two")
	if err != nil {
		t.Fatalf("a different key was blocked: %v", err)
	}
	otherKey()
	otherScope, err := AcquireKeyedLock(ctx, "test_independent_other", "key-one")
	if err != nil {
		t.Fatalf("the same key in a different scope was blocked: %v", err)
	}
	otherScope()
}

func TestAcquireKeyedLock_CancellationWhileWaiting(t *testing.T) {
	unlock, err := AcquireKeyedLock(context.Background(), "test_cancel", "key")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		waiterUnlock, waitErr := AcquireKeyedLock(ctx, "test_cancel", "KEY")
		if waiterUnlock != nil {
			waiterUnlock()
		}
		done <- waitErr
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("waiter did not return after cancellation")
	}

	unlock()
	assertKeyedLockFree(t, "test_cancel", "key")
}

func TestAcquireKeyedLock_DeadlineWhileWaiting(t *testing.T) {
	unlock, err := AcquireKeyedLock(context.Background(), "test_deadline", "key")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	waiterUnlock, err := AcquireKeyedLock(ctx, "test_deadline", "key")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context.DeadlineExceeded", err)
	}
	if waiterUnlock != nil {
		t.Fatal("expected a nil unlock func on failure")
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("waited %s, past the deadline", elapsed)
	}

	unlock()
	assertKeyedLockFree(t, "test_deadline", "key")
}

func TestAcquireKeyedLock_DoneContextDoesNotTakeFreeLock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for range 20 {
		unlock, err := AcquireKeyedLock(ctx, "test_done", "key")
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
		if unlock != nil {
			t.Fatal("expected a nil unlock func on failure")
		}
	}

	expired, stop := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer stop()
	if _, err := AcquireKeyedLock(expired, "test_done", "key"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context.DeadlineExceeded", err)
	}
	assertKeyedLockFree(t, "test_done", "key")
}

func TestAcquireKeyedLock_UnlockIsIdempotent(t *testing.T) {
	unlock, err := AcquireKeyedLock(context.Background(), "test_idempotent", "key")
	if err != nil {
		t.Fatal(err)
	}
	unlock()

	held, err := AcquireKeyedLock(context.Background(), "test_idempotent", "key")
	if err != nil {
		t.Fatal(err)
	}
	defer held()

	// A second call to the first unlock must not release the lock now held by another writer.
	unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if _, err := AcquireKeyedLock(ctx, "test_idempotent", "key"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("a repeated unlock released another writer's lock: %v", err)
	}
}

func TestAcquireKeyedLock_LogsLongWaits(t *testing.T) {
	unlock, err := AcquireKeyedLock(context.Background(), "test_log", "key")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		time.Sleep(keyedLockWaitLogThreshold + 100*time.Millisecond)
		unlock()
	}()

	var output bytes.Buffer
	ctx := tflogtest.RootLogger(context.Background(), &output)
	waiterUnlock, err := AcquireKeyedLock(ctx, "test_log", "key")
	if err != nil {
		t.Fatal(err)
	}
	waiterUnlock()
	if !strings.Contains(output.String(), "for the test_log lock on key") {
		t.Fatalf("expected a long wait to be logged, got %q", output.String())
	}

	output.Reset()
	quick, err := AcquireKeyedLock(ctx, "test_log", "key")
	if err != nil {
		t.Fatal(err)
	}
	quick()
	if output.Len() != 0 {
		t.Fatalf("expected no log for an uncontended lock, got %q", output.String())
	}
}

func assertKeyedLockFree(t *testing.T, scope, key string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	unlock, err := AcquireKeyedLock(ctx, scope, key)
	if err != nil {
		t.Fatalf("lock was left held: %v", err)
	}
	unlock()
}
