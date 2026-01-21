---
page_title: "microsoft365_graph_beta_windows_365_cloud_pc_role_definition Resource - terraform-provider-microsoft365"
subcategory: "Windows 365"

description: |-
  Manages custom role definitions for Windows 365 using the /roleManagement/cloudPC/roleDefinitions endpoint. This resource is used to define sets of permissions that can be assigned to administrators for Cloud PC management, policy configuration, and administrative functions.
---

# microsoft365_graph_beta_windows_365_cloud_pc_role_definition (Resource)

Manages custom role definitions for Windows 365 using the `/roleManagement/cloudPC/roleDefinitions` endpoint. This resource is used to define sets of permissions that can be assigned to administrators for Cloud PC management, policy configuration, and administrative functions.

## Microsoft Documentation

- [unifiedRoleDefinition resource type (beta)](https://learn.microsoft.com/en-us/graph/api/resources/unifiedroledefinition?view=graph-rest-beta)
- [Create roleDefinition (beta)](https://learn.microsoft.com/en-us/graph/api/rbacapplication-post-roledefinitions?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `RoleManagement.ReadWrite.CloudPC`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.25.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
resource "microsoft365_graph_beta_windows_365_cloud_pc_role_definition" "example" {
  display_name = "Windows 365 Cloud PC Role Definition"
  description  = "Custom role for Windows 365 Cloud PC administration with scoped permissions"

  role_permissions = [
    {
      allowed_resource_actions = [
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
        "Microsoft.CloudPC/AdminHighlights/Operate"
      ]
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `display_name` (String) Display Name of the Role definition.
- `role_permissions` (Attributes List) List of Role Permissions this role is allowed to perform. Not used for in-built Cloud PC role definitions. (see [below for nested schema](#nestedatt--role_permissions))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) Key of the entity. This is read-only and automatically generated.
- `is_built_in` (Boolean) Type of Role. Set to True if it is built-in, or set to False if it is a custom role definition.
- `is_built_in_role_definition` (Boolean) Type of Role. Set to True if it is built-in, or set to False if it is a custom role definition.

<a id="nestedatt--role_permissions"></a>
### Nested Schema for `role_permissions`

Optional:

- `allowed_resource_actions` (Set of String) Allowed actions for this role permission. This field is equivalent to 'actions' and can be used interchangeably. The API will consolidate values from both fields. Each action must start with 'Microsoft.CloudPC/'.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **Cloud PC Role Management**: This resource manages custom role definitions for Windows 365 Cloud PC administration.
- **Permission Control**: Defines granular permissions for Cloud PC operations, device management, and administrative functions.
- **Built-in Roles**: Cannot modify built-in roles such as "Cloud PC Administrator" and "Cloud PC Reader".
- **Role Permissions**: Supports Cloud PC-specific permissions starting with "Microsoft.CloudPC/".
- **Administrative Scope**: Role definitions apply to Windows 365 Cloud PC management within the tenant.
- **Security Model**: Integrates with Entra ID (Azure AD) role-based access control for secure permission management.
- **Permission Validation**: All specified permissions are validated against available Cloud PC operations.
- **Custom Roles**: Enables creation of custom roles tailored to specific organizational needs and responsibilities.

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_windows_365_cloud_pc_role_definition.example 00000000-0000-0000-0000-000000000000
```