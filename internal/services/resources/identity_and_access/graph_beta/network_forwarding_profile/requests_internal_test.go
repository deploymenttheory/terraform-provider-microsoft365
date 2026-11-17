package graphBetaNetworkForwardingProfile

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

func TestCreatedLookupAttemptsFor(t *testing.T) {
	ctx := context.Background()
	transport := func(msg string) kiotaerrors.GraphErrorInfo {
		err := &url.Error{Op: "Post", URL: "https://graph.microsoft.com/beta", Err: errors.New(msg)}
		return kiotaerrors.GraphError(ctx, fmt.Errorf("create forwarding profile: %w", err))
	}
	cases := map[string]struct {
		info kiotaerrors.GraphErrorInfo
		want int
	}{
		"connection reset":      {transport("read: connection reset by peer"), createdLookupAttempts},
		"unexpected EOF":        {transport("unexpected EOF"), createdLookupAttempts},
		"client timeout":        {transport("context deadline exceeded"), createdLookupAttempts},
		"gateway timeout":       {kiotaerrors.GraphErrorInfo{StatusCode: 504, IsODataError: true}, createdLookupAttempts},
		"non-graph 500":         {kiotaerrors.GraphErrorInfo{StatusCode: 500}, createdLookupAttempts},
		"duplicate after retry": {kiotaerrors.GraphErrorInfo{StatusCode: 400, IsODataError: true}, duplicateLookupAttempts},
		"conflict":              {kiotaerrors.GraphErrorInfo{StatusCode: 409, IsODataError: true}, duplicateLookupAttempts},
		"graph 500":             {kiotaerrors.GraphErrorInfo{StatusCode: 500, IsODataError: true}, 1},
		"forbidden":             {kiotaerrors.GraphErrorInfo{StatusCode: 403, IsODataError: true}, 1},
	}
	for name, tc := range cases {
		if got := createdLookupAttemptsFor(tc.info); got != tc.want {
			t.Errorf("%s: attempts = %d (status %d, code %q), expected %d",
				name, got, tc.info.StatusCode, tc.info.ErrorCode, tc.want)
		}
	}
}

// Only one of several concurrent callers can claim a profile.
func TestTryClaimIsExclusive(t *testing.T) {
	id := "ABCDEF00-0000-0000-0000-000000000001"
	t.Cleanup(func() { claimedProfiles.Delete(strings.ToLower(id)) })
	var wins atomic.Int32
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if tryClaim(id) {
				wins.Add(1)
			}
		}()
	}
	wg.Wait()
	if wins.Load() != 1 {
		t.Fatalf("expected exactly one claim, got %d", wins.Load())
	}
	if tryClaim(strings.ToLower(id)) {
		t.Fatal("claim must ignore ID case")
	}
}
