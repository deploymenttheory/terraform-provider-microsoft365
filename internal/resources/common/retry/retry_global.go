// REF: https://learn.microsoft.com/en-us/graph/throttling

package retry

import (
	"context"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/exp/rand"
)

// ThrottleScope represents the scope of throttling from x-ms-throttle-scope header
type ThrottleScope struct {
	Scope         string
	Limit         string
	ApplicationID string
	ResourceID    string
}

// parseThrottleScope parses the x-ms-throttle-scope header
func parseThrottleScope(scope string) ThrottleScope {
	parts := strings.Split(scope, "/")
	if len(parts) != 4 {
		return ThrottleScope{}
	}
	return ThrottleScope{
		Scope:         parts[0],
		Limit:         parts[1],
		ApplicationID: parts[2],
		ResourceID:    parts[3],
	}
}

// RetryableOperation executes an operation with automatic retry on rate limit errors
func RetryableOperation(ctx context.Context, operation string, fn func() error) error {
	var attempt int
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

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

		const maxBackoff = 10 * time.Second
		baseDelay := 2 * time.Second

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
			"operation":     operation,
			"attempt":       attempt,
			"delay_seconds": jitterDelay.Seconds(),
			"status_code":   graphError.StatusCode,
		}

		if throttleInfo != "" {
			logDetails["throttle_reason"] = throttleInfo
		}
		if throttleScope != (ThrottleScope{}) {
			logDetails["throttle_scope"] = throttleScope.Scope
			logDetails["throttle_limit"] = throttleScope.Limit
		}

		tflog.Info(ctx, "Microsoft Graph rate limit encountered", logDetails)

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
