package graphBetaIdentityAndAccessB2bManagementPolicy

import "time"

// SetNotFoundConfirmationForTesting shortens the refresh not-found confirmation window so unit
// tests can exercise a genuine deletion without waiting a minute. It returns a restore function.
func SetNotFoundConfirmationForTesting(window, interval time.Duration) (restore func()) {
	previousWindow, previousInterval := notFoundConfirmationWindow, notFoundConfirmationInterval
	notFoundConfirmationWindow, notFoundConfirmationInterval = window, interval
	return func() {
		notFoundConfirmationWindow, notFoundConfirmationInterval = previousWindow, previousInterval
	}
}
