// REF: https://learn.microsoft.com/en-us/graph/throttling-limits#assignment-service-limits

package retry

import (
	"context"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/rand"
)

// RetryableAssignmentOperation executes an assignment operation with specific rate limiting
func RetryableAssignmentOperation(ctx context.Context, operation string, fn func() error) error {
	var attempt int
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	const (
		tenSecondLimit = 500   // requests per 10 seconds per app per tenant
		hourlyLimit    = 15000 // requests per hour per app per tenant
		maxBackoff     = 10 * time.Second
		baseDelay      = 3 * time.Second
	)

	for {
		err := fn()
		if err == nil {
			return nil
		}

		graphError := errors.GraphError(ctx, err)
		if graphError.StatusCode != 429 {
			return err
		}

		// Parse throttle scope if available
		var throttleScope ThrottleScope
		if scope := graphError.Headers.Get("x-ms-throttle-scope"); len(scope) > 0 {
			throttleScope = parseThrottleScope(scope[0])
		}

		// Get throttle information
		var throttleInfo string
		if info := graphError.Headers.Get("x-ms-throttle-information"); len(info) > 0 {
			throttleInfo = info[0]
		}

		// Use Retry-After if provided, otherwise use exponential backoff
		var backoffDelay time.Duration
		if graphError.RetryAfter != "" {
			if seconds, err := time.ParseDuration(graphError.RetryAfter + "s"); err == nil {
				backoffDelay = seconds
			}
		}

		if backoffDelay == 0 {
			backoffDelay = baseDelay * time.Duration(1<<attempt)
			if backoffDelay > maxBackoff {
				backoffDelay = maxBackoff
			}
		}

		// Add jitter: randomly between 50-100% of calculated delay
		jitterDelay := backoffDelay/2 + time.Duration(r.Int63n(int64(backoffDelay/2)))
		attempt++

		logDetails := map[string]interface{}{
			"operation":      operation,
			"attempt":        attempt,
			"delay_seconds":  jitterDelay.Seconds(),
			"status_code":    graphError.StatusCode,
			"rate_limit_10s": tenSecondLimit,
			"rate_limit_1h":  hourlyLimit,
		}

		if throttleInfo != "" {
			logDetails["throttle_reason"] = throttleInfo
		}
		if throttleScope != (ThrottleScope{}) {
			logDetails["throttle_scope"] = throttleScope.Scope
			logDetails["throttle_limit"] = throttleScope.Limit
		}

		tflog.Info(ctx, "Microsoft Graph assignment rate limit encountered", logDetails)

		timer := time.NewTimer(jitterDelay)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			continue
		}
	}
}
