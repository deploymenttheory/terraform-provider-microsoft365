---
page_title: "microsoft365_graph_beta_device_management_app_control_for_business_policy Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages App Control for Business configuration policies with custom XML content using the /deviceManagement/configurationPolicies endpoint. This resource allows you to deploy custom App Control for Business policies by providing XML policy content directly.
---

# microsoft365_graph_beta_device_management_app_control_for_business_policy (Resource)

Manages App Control for Business configuration policies with custom XML content using the `/deviceManagement/configurationPolicies` endpoint. This resource allows you to deploy custom App Control for Business policies by providing XML policy content directly.

## Microsoft Documentation

- [App Control for Business](https://learn.microsoft.com/en-us/windows/security/application-security/application-control/app-control-for-business/appcontrol-and-applocker-overview)
- [deviceManagementConfigurationPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta)
- [Create deviceManagementConfigurationPolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfigv2-devicemanagementconfigurationpolicy-create?view=graph-rest-beta)
- [Update deviceManagementConfigurationPolicy](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfigv2-devicemanagementconfigurationpolicy-update?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.28.0-alpha | Testing | Initial release |

## Example Usage

```terraform
# example of inline XML policy
resource "microsoft365_graph_beta_device_management_app_control_for_business_policy" "inline" {
  name        = "unit-test-app-control-policy-inline"
  description = "unit-test-app-control-policy-inline"

  policy_xml = <<-EOT
<?xml version="1.0" encoding="utf-8"?>
<SiPolicy xmlns="urn:schemas-microsoft-com:sipolicy" PolicyType="Base Policy">
  <VersionEx>1.0.3.0</VersionEx>
  <PolicyID>{264C0644-19BE-418F-BAED-29E5E36250AD}</PolicyID>
  <BasePolicyID>{264C0644-19BE-418F-BAED-29E5E36250AD}</BasePolicyID>
  <PlatformID>{2E07F7E4-194C-4D20-B7C9-6F44A6C5A234}</PlatformID>
  <Rules>
    <Rule>
      <Option>Enabled:Unsigned System Integrity Policy</Option>
    </Rule>
    <Rule>
      <Option>Enabled:Advanced Boot Options Menu</Option>
    </Rule>
	<Rule>
      <Option>Enabled:UMCI</Option>
    </Rule>
	<Rule>
      <Option>Enabled:Update Policy No Reboot</Option>
    </Rule>
    <Rule>
      <Option>Enabled:Inherit Default Policy</Option>
    </Rule>
    <Rule>
      <Option>Enabled:Revoked Expired As Unsigned</Option>
    </Rule>
    <Rule>
      <Option>Disabled:Script Enforcement</Option>
    </Rule>
  </Rules>
  <!--EKUS-->
  <EKUs />
  <!--File Rules-->
  <FileRules>
    <Allow ID="ID_ALLOW_A_1" FileName="*" />
    <Allow ID="ID_ALLOW_A_2" FileName="*" />
  </FileRules>
  <!--Signers-->
  <Signers />
  <!--Driver Signing Scenarios-->
  <SigningScenarios>
    <SigningScenario Value="131" ID="ID_SIGNINGSCENARIO_DRIVERS_1" FriendlyName="Auto generated policy on 08-17-2015">
      <ProductSigners>
        <FileRulesRef>
          <FileRuleRef RuleID="ID_ALLOW_A_1" />
        </FileRulesRef>
      </ProductSigners>
    </SigningScenario>
    <SigningScenario Value="12" ID="ID_SIGNINGSCENARIO_WINDOWS" FriendlyName="Auto generated policy on 08-17-2015">
      <ProductSigners>
        <FileRulesRef>
          <FileRuleRef RuleID="ID_ALLOW_A_2" />
        </FileRulesRef>
      </ProductSigners>
    </SigningScenario>
  </SigningScenarios>
  <UpdatePolicySigners />
  <CiSigners />
  <HvciOptions>1</HvciOptions>
  <Settings>
	<Setting Provider="AllHostIds" Key="AllKeys" ValueName="EnterpriseDefinedClsId">
	  <Value>
        <Boolean>true</Boolean>
      </Value>
	</Setting>
    <Setting Provider="PolicyInfo" Key="Information" ValueName="Name">
      <Value>
        <String>AllowAllEnableHVCI</String>
      </Value>
    </Setting>
    <Setting Provider="PolicyInfo" Key="Information" ValueName="Id">
      <Value>
        <String>022422</String>
      </Value>
    </Setting>
  </Settings>
</SiPolicy>
  EOT

  # Role scope tags for testing
  role_scope_tag_ids = ["0", "1", "2"]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "33333333-3333-3333-3333-333333333333"
      filter_id   = "44444444-4444-4444-4444-444444444444"
      filter_type = "include"
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "55555555-5555-5555-5555-555555555555"
      filter_type = "exclude"
    }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}

# App Control for Business XML Policy configuration referencing a file
resource "microsoft365_graph_beta_device_management_app_control_for_business_policy" "file" {
  name        = "app-control-policy-file"
  description = "app-control-policy-file"

  # XML policy referencing a file
  policy_xml = file("${path.module}/AllowAll_EnableHVCI.xml")

  role_scope_tag_ids = ["0", "1", "2"]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.group_1.id
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.assignment_filter_1.id
      filter_type = "include"
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.assignment_filter_2.id
      filter_type = "exclude"
    }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the App Control for Business policy.
- `policy_xml` (String) The XML content of the App Control for Business policy. When you create policies for use with App Control for Business, start from an existing base policy and then add or remove rules to build your own custom policy. Windows includes several example policies that you can use. These example policies are provided by microsoft 'as-is'. You should thoroughly test the policies you deploy using safe deployment methods. These base policies can be found on Windows 11 22H2 and later devices. The locations of these policies can be found [here](https://learn.microsoft.com/en-us/windows/security/application-security/application-control/app-control-for-business/design/example-appcontrol-base-policies). For more information on policy rules and file rules, please see the [Understand App Control for Business policy rules and file rules](https://learn.microsoft.com/en-us/windows/security/application-security/application-control/app-control-for-business/design/select-types-of-rules-to-create).However, if you prefer an easier method, you should try the community-based tool [AppControl Manager](https://github.com/HotCakeX/Harden-Windows-Security/wiki/AppControl-Manager)

### Optional

- `assignments` (Attributes Set) Assignments for the device configuration. Each assignment specifies the target group and schedule for script execution. Supports group filters. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Optional description of the resource. Maximum length is 1500 characters.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this App Control for Business policy.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier of the app control for business policy.

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.

Optional:

- `filter_id` (String) ID of the filter to apply to the assignment.
- `filter_type` (String) Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.
- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

- **App Control for Business Design Guide**: [here](https://learn.microsoft.com/en-us/windows/security/application-security/application-control/app-control-for-business/design/appcontrol-design-guide)

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# {resource_id}
terraform import microsoft365_graph_beta_device_management_app_control_for_business_policy.example 00000000-0000-0000-0000-000000000000
```
