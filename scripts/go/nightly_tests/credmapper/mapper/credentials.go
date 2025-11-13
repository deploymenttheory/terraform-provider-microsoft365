package mapper

// ServiceCredential represents the environment variable names for a service's credentials.
type ServiceCredential struct {
	ClientIDVar     string
	ClientSecretVar string
}

// Service name constants.
const (
	ServiceApplications           = "applications"
	ServiceBackupStorage          = "backup_storage"
	ServiceDeviceAndAppMgmt       = "device_and_app_management"
	ServiceDeviceMgmt             = "device_management"
	ServiceGroups                 = "groups"
	ServiceIdentityAccess         = "identity_and_access"
	ServiceM365Admin              = "m365_admin"
	ServiceUsers                  = "users"
	ServiceWindows365             = "windows_365"
	ServiceMultitenantMgmt        = "multitenant_management"
	ServiceUtility                = "utility"
)

// Environment variable names for service credentials.
const (
	EnvM365ClientID               = "M365_CLIENT_ID"
	EnvM365ClientSecret           = "M365_CLIENT_SECRET"

	EnvApplicationsClientID       = "M365_CLIENT_ID_APPLICATIONS"
	EnvApplicationsClientSecret   = "M365_CLIENT_SECRET_APPLICATIONS"

	EnvBackupStorageClientID      = "M365_CLIENT_ID_BACKUP_STORAGE"
	EnvBackupStorageClientSecret  = "M365_CLIENT_SECRET_BACKUP_STORAGE"

	EnvDeviceAppMgmtClientID      = "M365_CLIENT_ID_DEVICE_AND_APP_MGMT"
	EnvDeviceAppMgmtClientSecret  = "M365_CLIENT_SECRET_DEVICE_AND_APP_MGMT"

	EnvDeviceMgmtClientID         = "M365_CLIENT_ID_DEVICE_MGMT"
	EnvDeviceMgmtClientSecret     = "M365_CLIENT_SECRET_DEVICE_MGMT"

	EnvGroupsClientID             = "M365_CLIENT_ID_GROUPS"
	EnvGroupsClientSecret         = "M365_CLIENT_SECRET_GROUPS"

	EnvIdentityAccessClientID     = "M365_CLIENT_ID_IDENTITY_ACCESS"
	EnvIdentityAccessClientSecret = "M365_CLIENT_SECRET_IDENTITY_ACCESS"

	EnvM365AdminClientID          = "M365_CLIENT_ID_M365_ADMIN"
	EnvM365AdminClientSecret      = "M365_CLIENT_SECRET_M365_ADMIN"

	EnvUsersClientID              = "M365_CLIENT_ID_USERS"
	EnvUsersClientSecret          = "M365_CLIENT_SECRET_USERS"

	EnvWindows365ClientID         = "M365_CLIENT_ID_WINDOWS_365"
	EnvWindows365ClientSecret     = "M365_CLIENT_SECRET_WINDOWS_365"

	EnvMultitenantMgmtClientID    = "M365_CLIENT_ID_MULTITENANT_MGMT"
	EnvMultitenantMgmtClientSecret= "M365_CLIENT_SECRET_MULTITENANT_MGMT"

	EnvUtilityClientID            = "M365_CLIENT_ID_UTILITY"
	EnvUtilityClientSecret        = "M365_CLIENT_SECRET_UTILITY"
)

// serviceCredentials maps service names to their credential environment variable names.
var serviceCredentials = map[string]ServiceCredential{
	ServiceApplications: {
		ClientIDVar:     EnvApplicationsClientID,
		ClientSecretVar: EnvApplicationsClientSecret,
	},
	ServiceBackupStorage: {
		ClientIDVar:     EnvBackupStorageClientID,
		ClientSecretVar: EnvBackupStorageClientSecret,
	},
	ServiceDeviceAndAppMgmt: {
		ClientIDVar:     EnvDeviceAppMgmtClientID,
		ClientSecretVar: EnvDeviceAppMgmtClientSecret,
	},
	ServiceDeviceMgmt: {
		ClientIDVar:     EnvDeviceMgmtClientID,
		ClientSecretVar: EnvDeviceMgmtClientSecret,
	},
	ServiceGroups: {
		ClientIDVar:     EnvGroupsClientID,
		ClientSecretVar: EnvGroupsClientSecret,
	},
	ServiceIdentityAccess: {
		ClientIDVar:     EnvIdentityAccessClientID,
		ClientSecretVar: EnvIdentityAccessClientSecret,
	},
	ServiceM365Admin: {
		ClientIDVar:     EnvM365AdminClientID,
		ClientSecretVar: EnvM365AdminClientSecret,
	},
	ServiceUsers: {
		ClientIDVar:     EnvUsersClientID,
		ClientSecretVar: EnvUsersClientSecret,
	},
	ServiceWindows365: {
		ClientIDVar:     EnvWindows365ClientID,
		ClientSecretVar: EnvWindows365ClientSecret,
	},
	ServiceMultitenantMgmt: {
		ClientIDVar:     EnvMultitenantMgmtClientID,
		ClientSecretVar: EnvMultitenantMgmtClientSecret,
	},
	ServiceUtility: {
		ClientIDVar:     EnvUtilityClientID,
		ClientSecretVar: EnvUtilityClientSecret,
	},
}

// GetServiceCredential returns the credential configuration for a service.
func GetServiceCredential(service string) (ServiceCredential, bool) {
	cred, ok := serviceCredentials[service]
	return cred, ok
}

// SupportedServices returns a list of all supported service names.
func SupportedServices() []string {
	services := make([]string, 0, len(serviceCredentials))
	for service := range serviceCredentials {
		services = append(services, service)
	}
	return services
}
