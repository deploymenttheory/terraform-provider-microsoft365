package graphBetaNetworkForwardingProfile

import "time"

// ShortenCreatedLookup makes the post-504 lookup fast in unit tests.
func ShortenCreatedLookup() func() {
	attempts, duplicates, interval := createdLookupAttempts, duplicateLookupAttempts, createdLookupInterval
	createdLookupAttempts, duplicateLookupAttempts, createdLookupInterval = 2, 2, 10*time.Millisecond
	return func() {
		createdLookupAttempts, duplicateLookupAttempts, createdLookupInterval = attempts, duplicates, interval
	}
}
