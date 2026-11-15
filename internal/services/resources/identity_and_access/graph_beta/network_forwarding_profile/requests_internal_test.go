package graphBetaNetworkForwardingProfile

import (
	"context"
	"errors"
	"fmt"
	"net/url"
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
