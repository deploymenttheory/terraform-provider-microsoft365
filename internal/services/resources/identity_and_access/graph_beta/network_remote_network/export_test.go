package graphBetaIdentityAndAccessNetworkRemoteNetwork

import "time"

// SetDeletionPollInterval shortens the deletion poll for unit tests.
func SetDeletionPollInterval(d time.Duration) (restore func()) {
	previous := deletionPollInterval
	deletionPollInterval = d
	return func() { deletionPollInterval = previous }
}
