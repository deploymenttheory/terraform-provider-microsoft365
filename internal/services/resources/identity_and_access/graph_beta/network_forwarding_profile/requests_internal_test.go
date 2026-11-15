package graphBetaNetworkForwardingProfile

import (
	"testing"

	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

func TestCreatedLookupAttemptsFor(t *testing.T) {
	cases := map[string]struct {
		info kiotaerrors.GraphErrorInfo
		want int
	}{
		"gateway timeout":       {kiotaerrors.GraphErrorInfo{StatusCode: 504, IsODataError: true}, createdLookupAttempts},
		"connection reset":      {kiotaerrors.GraphErrorInfo{StatusCode: 400, ErrorCode: "URLError", Category: kiotaerrors.CategoryNetwork}, createdLookupAttempts},
		"non-graph 500":         {kiotaerrors.GraphErrorInfo{StatusCode: 500}, createdLookupAttempts},
		"duplicate after retry": {kiotaerrors.GraphErrorInfo{StatusCode: 400, IsODataError: true}, duplicateLookupAttempts},
		"conflict":              {kiotaerrors.GraphErrorInfo{StatusCode: 409, IsODataError: true}, duplicateLookupAttempts},
		"graph 500":             {kiotaerrors.GraphErrorInfo{StatusCode: 500, IsODataError: true}, 1},
		"forbidden":             {kiotaerrors.GraphErrorInfo{StatusCode: 403, IsODataError: true}, 1},
	}
	for name, tc := range cases {
		if got := createdLookupAttemptsFor(tc.info); got != tc.want {
			t.Errorf("%s: attempts = %d, expected %d", name, got, tc.want)
		}
	}
}
