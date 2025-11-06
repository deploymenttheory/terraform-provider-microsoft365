package generic_client

import (
	"context"
	"testing"
	"time"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/stretchr/testify/assert"
)

// TestGetRetryDelay tests the GetRetryDelay function
func TestGetRetryDelay(t *testing.T) {
	testCases := []struct {
		name        string
		errorInfo   errors.GraphErrorInfo
		attempt     int
		minExpected time.Duration
		maxExpected time.Duration
	}{
		{
			name: "With RetryAfter header",
			errorInfo: errors.GraphErrorInfo{
				RetryAfter: "30",
			},
			attempt:     1,
			minExpected: 30 * time.Second,         // Minimum is the API-specified value
			maxExpected: 37500 * time.Millisecond, // Maximum is base + 25% jitter (30s + 7.5s)
		},
		{
			name:        "First attempt without RetryAfter",
			errorInfo:   errors.GraphErrorInfo{},
			attempt:     1,
			minExpected: 750 * time.Millisecond, // 1s ± 25% jitter = 1s ± 0.25s
			maxExpected: 1250 * time.Millisecond,
		},
		{
			name:        "Second attempt without RetryAfter",
			errorInfo:   errors.GraphErrorInfo{},
			attempt:     2,
			minExpected: 3 * time.Second, // 4s ± 25% jitter = 4s ± 1s = 3s to 5s
			maxExpected: 5 * time.Second,
		},
		{
			name:        "Third attempt without RetryAfter",
			errorInfo:   errors.GraphErrorInfo{},
			attempt:     3,
			minExpected: 6750 * time.Millisecond, // 9s ± 25% jitter = 9s ± 2.25s = 6.75s to 11.25s
			maxExpected: 11250 * time.Millisecond,
		},
		{
			name:        "Max delay cap",
			errorInfo:   errors.GraphErrorInfo{},
			attempt:     100,             // Very large attempt number
			minExpected: 0,               // Could be capped and then have negative jitter
			maxExpected: 5 * time.Minute, // Should never exceed max cap
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			delay := GetRetryDelay(ctx, &tc.errorInfo, tc.attempt)

			// All cases should respect min and max bounds
			assert.GreaterOrEqual(t, delay, tc.minExpected, "Delay should be greater than or equal to min expected")
			assert.LessOrEqual(t, delay, tc.maxExpected, "Delay should be less than or equal to max expected")
		})
	}
}
