package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"errors"
	"testing"
)

func TestLockIpApplicationSegmentWritesReturnsWhenContextIsCancelled(t *testing.T) {
	unlock, err := lockIpApplicationSegmentWrites(context.Background(), "lock-test-application")
	if err != nil {
		t.Fatalf("first lock returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := lockIpApplicationSegmentWrites(ctx, "lock-test-application"); !errors.Is(err, context.Canceled) {
		t.Fatalf("lock while held = %v, expected context.Canceled", err)
	}

	unlock()
	unlock, err = lockIpApplicationSegmentWrites(context.Background(), "lock-test-application")
	if err != nil {
		t.Fatalf("lock after unlock returned error: %v", err)
	}
	unlock()
}
