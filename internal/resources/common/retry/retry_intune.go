// REF: https://learn.microsoft.com/en-us/graph/throttling-limits#intune-service-limits

package retry

import (
	"context"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/rand"
)

// IntuneOperationType defines the type of Intune operation
type IntuneOperationType string

const (
	IntuneWrite IntuneOperationType = "Write" // POST, PUT, DELETE, PATCH
	IntuneRead  IntuneOperationType = "Read"  // GET and others
)

// RetryableIntuneOperation executes an Intune operation with specific rate limiting
func RetryableIntuneOperation(ctx context.Context, operation string, opType IntuneOperationType, fn func() error) error {
	var attempt int
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	const (
		// Write operations (POST, PUT, DELETE, PATCH)
		writePerAppLimit = 100 // requests per 20 seconds
		writeTenantLimit = 200 // requests per 20 seconds

		// General operations
		generalPerAppLimit = 1000 // requests per 20 seconds
		generalTenantLimit = 2000 // requests per 20 seconds

		maxBackoff = 10 * time.Second
		baseDelay  = 2 * time.Second
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

		// Enhanced logging with rate limit context
		logDetails := map[string]interface{}{
			"operation":      operation,
			"attempt":        attempt,
			"delay_seconds":  jitterDelay.Seconds(),
			"status_code":    graphError.StatusCode,
			"operation_type": string(opType),
		}

		if opType == IntuneWrite {
			logDetails["rate_limit_per_app"] = writePerAppLimit
			logDetails["rate_limit_tenant"] = writeTenantLimit
			logDetails["window_seconds"] = 20
		} else {
			logDetails["rate_limit_per_app"] = generalPerAppLimit
			logDetails["rate_limit_tenant"] = generalTenantLimit
			logDetails["window_seconds"] = 20
		}

		if throttleInfo != "" {
			logDetails["throttle_reason"] = throttleInfo
		}
		if throttleScope != (ThrottleScope{}) {
			logDetails["throttle_scope"] = throttleScope.Scope
			logDetails["throttle_limit"] = throttleScope.Limit
		}

		tflog.Info(ctx, "Intune service rate limit encountered", logDetails)

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
