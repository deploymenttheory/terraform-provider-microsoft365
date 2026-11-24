package crud

import (
	"strings"
	"sync"
)

// claims maps scope + lower-cased key to an empty value once that key is claimed.
var claims sync.Map

// claimRegistryKey builds the map key for a scope and key; keys compare case-insensitively.
func claimRegistryKey(scope, key string) string {
	return scope + "\x00" + strings.ToLower(key)
}

// ClaimKey records key under scope and reports whether it was already recorded. The check
// and the record are one atomic step, so when several callers claim the same scope and key
// at once, exactly one of them sees alreadyClaimed == false.
//
// The registry is a process-wide record of claims, not a lock: entries are never released,
// so a claimed key stays claimed for the life of the provider process. Use it to stop two
// writers adopting the same remote object, or to stop a deleted object being adopted again
// from a stale list. Keys are compared case-insensitively, and scopes never see each
// other's claims.
func ClaimKey(scope, key string) (alreadyClaimed bool) {
	_, alreadyClaimed = claims.LoadOrStore(claimRegistryKey(scope, key), struct{}{})
	return alreadyClaimed
}

// RecordClaim records key under scope without checking it, such as after a create or delete.
// Like ClaimKey, the record is never released.
func RecordClaim(scope, key string) {
	claims.Store(claimRegistryKey(scope, key), struct{}{})
}

// ResetClaims forgets every claim recorded under scope. It is for tests only: provider code
// must never call it, since claims are meant to last for the life of the process. Unit tests
// that reuse one mock ID across test cases call it to start each case with no claims.
func ResetClaims(scope string) {
	prefix := scope + "\x00"
	claims.Range(func(key, _ any) bool {
		if strings.HasPrefix(key.(string), prefix) {
			claims.Delete(key)
		}
		return true
	})
}
