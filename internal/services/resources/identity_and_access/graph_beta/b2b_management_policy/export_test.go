package graphBetaIdentityAndAccessB2bManagementPolicy

import "time"

// SetNotFoundConfirmationForTesting overrides the not-found confirmation timing and returns a restore func.
func SetNotFoundConfirmationForTesting(window, interval time.Duration) (restore func()) {
	previousWindow, previousInterval := notFoundConfirmationWindow, notFoundConfirmationInterval
	notFoundConfirmationWindow, notFoundConfirmationInterval = window, interval
	return func() {
		notFoundConfirmationWindow, notFoundConfirmationInterval = previousWindow, previousInterval
	}
}
