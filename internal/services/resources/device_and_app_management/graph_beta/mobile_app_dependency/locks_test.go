package graphBetaMobileAppDependency

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestUnitAppRelationshipLockSerializesWriters(t *testing.T) {
	key := "aabbccdd-1111-1111-1111-111111111111"
	var writers sync.WaitGroup
	var count int
	for index := range 50 {
		writers.Add(1)
		go func() {
			defer writers.Done()
			appID := key
			if index%2 == 0 {
				appID = strings.ToUpper(key)
			}
			unlock, err := lockAppRelationships(context.Background(), appID)
			if err != nil {
				t.Error(err)
				return
			}
			count++
			unlock()
		}()
	}
	writers.Wait()
	if count != 50 {
		t.Fatalf("completed writes = %d, want 50", count)
	}
}

func TestUnitAppRelationshipLocksAllowIndependentAppsAndCancellation(t *testing.T) {
	unlock, err := lockAppRelationships(context.Background(), "app-one")
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	otherUnlock, err := lockAppRelationships(ctx, "app-two")
	if err != nil {
		t.Fatalf("an unrelated app was blocked: %v", err)
	}
	otherUnlock()
	canceled, stop := context.WithCancel(context.Background())
	stop()
	if _, err := lockAppRelationships(canceled, "app-one"); err == nil {
		t.Fatal("waiting for an app lock ignored cancellation")
	}
}
