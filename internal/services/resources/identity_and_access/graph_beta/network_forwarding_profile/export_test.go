package graphBetaNetworkForwardingProfile

import (
	"sync"
	"time"
)

// ShortenCreatedLookup makes the post-504 lookup fast in unit tests.
func ShortenCreatedLookup() func() {
	attempts, duplicates, interval := createdLookupAttempts, duplicateLookupAttempts, createdLookupInterval
	createdLookupAttempts, duplicateLookupAttempts, createdLookupInterval = 2, 2, 10*time.Millisecond
	return func() {
		createdLookupAttempts, duplicateLookupAttempts, createdLookupInterval = attempts, duplicates, interval
	}
}

// ResetClaimedProfiles forgets claimed IDs, since the mock reuses one ID.
func ResetClaimedProfiles() { claimedProfiles = sync.Map{} }

// ClaimProfile marks an ID as created by another resource in this apply.
func ClaimProfile(id string) { claimProfile(id) }
