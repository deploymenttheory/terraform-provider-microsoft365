package graphBetaNetworkForwardingProfile

import "time"

// ShortenCreatedLookup makes the post-504 lookup fast in unit tests.
func ShortenCreatedLookup() func() {
	attempts, interval := createdLookupAttempts, createdLookupInterval
	createdLookupAttempts, createdLookupInterval = 2, 10*time.Millisecond
	return func() { createdLookupAttempts, createdLookupInterval = attempts, interval }
}
