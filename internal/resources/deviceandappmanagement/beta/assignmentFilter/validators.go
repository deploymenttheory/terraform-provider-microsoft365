package graphBetaAssignmentFilter

/* platform type validator */

var validPlatformTypes = []string{
	"android",
	"androidForWork",
	"iOS",
	"macOS",
	"windowsPhone81",
	"windows81AndLater",
	"windows10AndLater",
	"androidWorkProfile",
	"unknown",
	"androidAOSP",
	"androidMobileApplicationManagement",
	"iOSMobileApplicationManagement",
	"windowsMobileApplicationManagement",
}

var validAssignmentFilterManagementTypes = []string{
	"devices",
	"apps",
	"unknownFutureValue",
}
