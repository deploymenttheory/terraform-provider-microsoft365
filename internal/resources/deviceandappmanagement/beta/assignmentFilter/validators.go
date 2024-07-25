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

func getValidAssignmentFilterTypes() []string {
	// This reflects the order in the SDK's String() method
	return []string{
		"unknown",
		"deviceConfigurationAndCompliance",
		"application",
		"androidEnterpriseApp",
		"enrollmentConfiguration",
		"groupPolicyConfiguration",
		"zeroTouchDeploymentDeviceConfigProfile",
		"androidEnterpriseConfiguration",
		"deviceFirmwareConfigurationInterfacePolicy",
		"resourceAccessPolicy",
		"win32app",
		"deviceManagmentConfigurationAndCompliancePolicy",
	}
}
