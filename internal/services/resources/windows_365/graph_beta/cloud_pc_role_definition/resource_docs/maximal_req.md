w365

POST https://graph.microsoft.com/beta/roleManagement/cloudPC/roleDefinitions

// Code snippets are only available for the latest major version. Current major version is $v0.*

// Dependencies
import (
	  "context"
	  msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	  graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	  //other-imports
)

requestBody := graphmodels.NewUnifiedRoleDefinition()
id := ""
requestBody.SetId(&id) 
description := "test"
requestBody.SetDescription(&description) 
displayName := "test"
requestBody.SetDisplayName(&displayName) 


unifiedRolePermission := graphmodels.NewUnifiedRolePermission()
condition := null
unifiedRolePermission.SetCondition(&condition) 
allowedResourceActions := []string {
	"Microsoft.CloudPC/CloudPCs/Read",
	"Microsoft.CloudPC/CloudPCs/Reprovision",
	"Microsoft.CloudPC/CloudPCs/Resize",
	"Microsoft.CloudPC/CloudPCs/EndGracePeriod",
	"Microsoft.CloudPC/CloudPCs/Restore",
	"Microsoft.CloudPC/CloudPCs/Reboot",
	"Microsoft.CloudPC/CloudPCs/Rename",
	"Microsoft.CloudPC/CloudPCs/Troubleshoot",
	"Microsoft.CloudPC/CloudPCs/ModifyDiskEncryptionType",
	"Microsoft.CloudPC/CloudPCs/ChangeUserAccountType",
	"Microsoft.CloudPC/CloudPCs/PlaceUnderReview",
	"Microsoft.CloudPC/CloudPCs/RetryPartnerAgentInstallation",
	"Microsoft.CloudPC/CloudPCs/ApplyCurrentProvisioningPolicy",
	"Microsoft.CloudPC/CloudPCs/CreateSnapshot",
	"Microsoft.CloudPC/CloudPCs/PowerOn",
	"Microsoft.CloudPC/CloudPCs/PowerOff",
	"Microsoft.CloudPC/CloudPCs/DisasterRecoveryFailover",
	"Microsoft.CloudPC/CloudPCs/DisasterRecoveryFailback",
	"Microsoft.CloudPC/CloudPCs/Start",
	"Microsoft.CloudPC/CloudPCs/Stop",
	"Microsoft.CloudPC/CloudPCs/GetCloudPcLaunchInfo",
	"Microsoft.CloudPC/CloudPCs/ReinstallAgent",
	"Microsoft.CloudPC/CloudPCs/CheckAgentStatus",
	"Microsoft.CloudPC/CloudPCs/RetrieveAgentStatus",
	"Microsoft.CloudPC/CloudPCs/Provision",
	"Microsoft.CloudPC/CloudPCs/Deprovision",
	"Microsoft.CloudPC/DeviceImages/Create",
	"Microsoft.CloudPC/DeviceImages/Delete",
	"Microsoft.CloudPC/DeviceImages/Read",
	"Microsoft.CloudPC/OnPremisesConnections/Create",
	"Microsoft.CloudPC/OnPremisesConnections/Delete",
	"Microsoft.CloudPC/OnPremisesConnections/Read",
	"Microsoft.CloudPC/OnPremisesConnections/Update",
	"Microsoft.CloudPC/OnPremisesConnections/RunHealthChecks",
	"Microsoft.CloudPC/OnPremisesConnections/UpdateAdDomainPassword",
	"Microsoft.CloudPC/ProvisioningPolicies/Assign",
	"Microsoft.CloudPC/ProvisioningPolicies/Apply",
	"Microsoft.CloudPC/ProvisioningPolicies/Create",
	"Microsoft.CloudPC/ProvisioningPolicies/Delete",
	"Microsoft.CloudPC/ProvisioningPolicies/Read",
	"Microsoft.CloudPC/ProvisioningPolicies/Update",
	"Microsoft.CloudPC/UserSettings/Assign",
	"Microsoft.CloudPC/UserSettings/Create",
	"Microsoft.CloudPC/UserSettings/Delete",
	"Microsoft.CloudPC/UserSettings/Read",
	"Microsoft.CloudPC/UserSettings/Update",
	"Microsoft.CloudPC/Roles/Read",
	"Microsoft.CloudPC/Roles/Create",
	"Microsoft.CloudPC/Roles/Update",
	"Microsoft.CloudPC/Roles/Delete",
	"Microsoft.CloudPC/RoleAssignments/Create",
	"Microsoft.CloudPC/RoleAssignments/Update",
	"Microsoft.CloudPC/RoleAssignments/Delete",
	"Microsoft.CloudPC/AuditData/Read",
	"Microsoft.CloudPC/SupportedRegion/Read",
	"Microsoft.CloudPC/ServicePlan/Read",
	"Microsoft.CloudPC/Snapshot/Read",
	"Microsoft.CloudPC/Snapshot/Share",
	"Microsoft.CloudPC/Snapshot/Import",
	"Microsoft.CloudPC/Snapshot/PurgeImportedSnapshot",
	"Microsoft.CloudPC/OrganizationSettings/Read",
	"Microsoft.CloudPC/OrganizationSettings/Update",
	"Microsoft.CloudPC/ExternalPartnerSettings/Read",
	"Microsoft.CloudPC/ExternalPartnerSettings/Create",
	"Microsoft.CloudPC/ExternalPartnerSettings/Update",
	"Microsoft.CloudPC/PerformanceReports/Read",
	"Microsoft.CloudPC/SharedUseServicePlans/Read",
	"Microsoft.CloudPC/FrontLineServicePlans/Read",
	"Microsoft.CloudPC/SharedUseLicenseUsageReports/Read",
	"Microsoft.CloudPC/FrontlineReports/Read",
	"Microsoft.CloudPC/CrossRegionDisasterRecovery/Read",
	"Microsoft.CloudPC/BulkActions/Read",
	"Microsoft.CloudPC/BulkActions/Write",
	"Microsoft.CloudPC/ActionStatus/Read",
	"Microsoft.CloudPC/InaccessibleReports/Read",
	"Microsoft.CloudPC/MaintenanceWindows/Assign",
	"Microsoft.CloudPC/MaintenanceWindows/Create",
	"Microsoft.CloudPC/MaintenanceWindows/Delete",
	"Microsoft.CloudPC/MaintenanceWindows/Read",
	"Microsoft.CloudPC/MaintenanceWindows/Update",
	"Microsoft.CloudPC/DeviceRecommendation/Read",
	"Microsoft.CloudPC/CloudApps/Read",
	"Microsoft.CloudPC/CloudApps/Publish",
	"Microsoft.CloudPC/CloudApps/Update",
	"Microsoft.CloudPC/CloudApps/Reset",
	"Microsoft.CloudPC/CloudApps/Unpublish",
	"Microsoft.CloudPC/Settings/Assign",
	"Microsoft.CloudPC/Settings/Create",
	"Microsoft.CloudPC/Settings/Read",
	"Microsoft.CloudPC/Settings/Update",
	"Microsoft.CloudPC/Settings/Delete",
	"Microsoft.CloudPC/AdminHighlights/Operate",
}
unifiedRolePermission.SetAllowedResourceActions(allowedResourceActions)

rolePermissions := []graphmodels.UnifiedRolePermissionable {
	unifiedRolePermission,
}
requestBody.SetRolePermissions(rolePermissions)

// To initialize your graphClient, see https://learn.microsoft.com/en-us/graph/sdks/create-client?from=snippets&tabs=go
roleDefinitions, err := graphClient.RoleManagement().CloudPC().RoleDefinitions().Post(context.Background(), requestBody, nil)


windows autopatch

{"Name":"test",
"Description":"test",
"Permissions":[
	"Microsoft.Autopatch_Roles_Read","Microsoft.Autopatch_Roles_Create","Microsoft.Autopatch_Roles_Edit","Microsoft.Autopatch_Roles_Delete","Microsoft.Autopatch_Roles_Assign","Microsoft.Autopatch_Reports_Read","Microsoft.Autopatch_Reports_ActionDiscoverDevices","Microsoft.Autopatch_Reports_ActionAssignRing","Microsoft.Autopatch_Reports_ActionExcludeDevices","Microsoft.Autopatch_Reports_ActionRestoreExcludedDevices","Microsoft.Autopatch_AutopatchGroups_Read","Microsoft.Autopatch_AutopatchGroups_Create","Microsoft.Autopatch_AutopatchGroups_Edit","Microsoft.Autopatch_AutopatchGroups_Delete","Microsoft.Autopatch_Messages_Read","Microsoft.Autopatch_SupportRequests_Read"
],
"ScopeTags":["0"]
}