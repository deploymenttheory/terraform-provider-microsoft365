package generic_client

import (
	"context"
	"strings"
	"time"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// GetRetryDelay calculates the appropriate retry delay based on error information
func GetRetryDelay(ctx context.Context, errorInfo *errors.GraphErrorInfo, attempt int) time.Duration {
	var baseDelay time.Duration

	// Try to parse Retry-After header if available
	if errorInfo.RetryAfter != "" {
		retryAfter := strings.TrimSpace(errorInfo.RetryAfter)

		tflog.Info(ctx, "Retry-After header received from API", map[string]any{
			"raw_value":     errorInfo.RetryAfter,
			"trimmed_value": retryAfter,
			"attempt":       attempt,
		})

		// Try parsing as-is (e.g., "2s", "30s")
		if duration, err := time.ParseDuration(retryAfter); err == nil {
			baseDelay = duration
			tflog.Info(ctx, "Successfully parsed Retry-After as Go duration format", map[string]any{
				"parsed_value": baseDelay.String(),
				"method":       "direct_parse",
			})
		} else if duration, err := time.ParseDuration(retryAfter + "s"); err == nil {
			// Try parsing with "s" suffix for plain numbers (e.g., "2" -> "2s")
			baseDelay = duration
			tflog.Info(ctx, "Successfully parsed Retry-After as plain number", map[string]any{
				"parsed_value": baseDelay.String(),
				"method":       "add_s_suffix",
			})
		} else {
			// Handle "N seconds" format (e.g., "2 seconds" -> "2s")
			retryAfter = strings.TrimSuffix(retryAfter, " seconds")
			retryAfter = strings.TrimSuffix(retryAfter, " second")
			if duration, err := time.ParseDuration(retryAfter + "s"); err == nil {
				baseDelay = duration
				tflog.Info(ctx, "Successfully parsed Retry-After after stripping 'seconds' suffix", map[string]any{
					"parsed_value": baseDelay.String(),
					"method":       "strip_seconds_suffix",
				})
			} else {
				tflog.Warn(ctx, "Failed to parse Retry-After header, will use exponential backoff", map[string]any{
					"raw_value": errorInfo.RetryAfter,
					"attempted_values": []string{
						errorInfo.RetryAfter,
						errorInfo.RetryAfter + "s",
						strings.TrimSuffix(strings.TrimSuffix(errorInfo.RetryAfter, " seconds"), " second") + "s",
					},
				})
			}
		}
	} else {
		tflog.Debug(ctx, "No Retry-After header present in API response, will use exponential backoff", map[string]any{
			"attempt": attempt,
		})
	}

	// If we successfully parsed the Retry-After header, use it with positive jitter
	if baseDelay > 0 {
		// Add positive jitter (0% to +25%) for safety
		// The API's Retry-After is a MINIMUM - we should never wait less than specified
		maxJitter := time.Duration(float64(baseDelay) * 0.25)

		// Generate random jitter factor between 0.0 and 1.0 (positive only)
		randomFactor := float64(time.Now().UnixNano()%1000) / 1000.0
		jitterAdjustment := time.Duration(float64(maxJitter) * randomFactor)

		delay := baseDelay + jitterAdjustment

		tflog.Info(ctx, "Using API-provided Retry-After delay with positive jitter", map[string]any{
			"base_delay":        baseDelay.String(),
			"jitter_range":      "+0% to +25%",
			"jitter_adjustment": jitterAdjustment.String(),
			"final_delay":       delay.String(),
			"min_delay":         baseDelay.String(),
			"max_delay":         (baseDelay + maxJitter).String(),
			"attempt":           attempt,
		})

		return delay
	}

	// Fallback: Exponential backoff with jitter if no Retry-After header
	baseDelay = time.Second
	maxDelay := 5 * time.Minute

	delay := time.Duration(attempt*attempt) * baseDelay

	// Add jitter (±25%)
	jitter := time.Duration(float64(delay) * 0.25)

	// Generate random jitter factor between -1.0 and 1.0
	randomFactor := (float64(time.Now().UnixNano()%1000)/1000.0)*2.0 - 1.0
	jitterAdjustment := time.Duration(float64(jitter) * randomFactor)

	delay += jitterAdjustment

	// Ensure delay is not negative
	if delay < 0 {
		delay = baseDelay
	}

	// Apply maximum cap after jitter to ensure we never exceed maxDelay
	if delay > maxDelay {
		delay = maxDelay
	}

	tflog.Info(ctx, "Using exponential backoff with jitter (no Retry-After header)", map[string]any{
		"base_delay":                     baseDelay.String(),
		"calculated_delay_before_jitter": time.Duration(attempt*attempt) * baseDelay,
		"jitter_adjustment":              jitterAdjustment.String(),
		"final_delay":                    delay.String(),
		"max_delay":                      maxDelay.String(),
		"attempt":                        attempt,
	})

	return delay
}
