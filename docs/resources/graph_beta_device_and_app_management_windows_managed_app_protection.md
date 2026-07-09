---
page_title: "microsoft365_graph_beta_device_and_app_management_windows_managed_app_protection Resource - terraform-provider-microsoft365"
subcategory: "Device and App Management"
description: |-
  Manages Windows Mobile Application Management (MAM) app protection policies in Microsoft Intune.
---

# microsoft365_graph_beta_device_and_app_management_windows_managed_app_protection

Manages Windows Mobile Application Management (MAM) app protection policies in Microsoft Intune. These policies control how managed apps handle corporate data on Windows devices, including data transfer restrictions, clipboard behaviour, and threat response actions.

## API Permissions

The following API permissions are required:

| Permission | Type |
|---|---|
| `DeviceManagementApps.Read.All` | Read |
| `DeviceManagementApps.ReadWrite.All` | Read/Write |

## Example Usage

```hcl
resource "microsoft365_graph_beta_device_and_app_management_windows_managed_app_protection" "example" {
  display_name = "Example Windows MAM Policy"
  description  = "Windows MAM policy managed by Terraform."

  allowed_inbound_data_transfer_sources        = "none"
  allowed_outbound_data_transfer_destinations  = "none"
  allowed_outbound_clipboard_sharing_level     = "none"

  print_blocked = true

  maximum_allowed_device_threat_level      = "notConfigured"
  mobile_threat_defense_remediation_action = "block"

  period_offline_before_wipe_is_enforced = "P90D"
  period_offline_before_access_check     = "P30D"
}
```

## Argument Reference

### Required

- `display_name` - (Required) Policy display name. Must be unique within your tenant.

### Optional

- `description` - (Optional) The policy's description. Defaults to `""`.
- `role_scope_tag_ids` - (Optional) List of scope tag IDs for this entity instance.
- `print_blocked` - (Optional) When `true`, printing is blocked from managed apps. Defaults to `false`.
- `allowed_inbound_data_transfer_sources` - (Optional) Sources from which data is allowed to be transferred. Possible values: `allApps`, `none`. Defaults to `allApps`.
- `allowed_outbound_clipboard_sharing_level` - (Optional) Clipboard sharing level across org and non-org resources. Possible values: `anyDestinationAnySource`, `none`, `orgDestinationAnySource`, `orgDestinationOrgSource`. Defaults to `anyDestinationAnySource`.
- `allowed_outbound_data_transfer_destinations` - (Optional) Destinations to which data is allowed to be transferred. Possible values: `allApps`, `none`. Defaults to `allApps`.
- `app_action_if_unable_to_authenticate_user` - (Optional) Action when user cannot authenticate. Possible values: `block`, `wipe`, `warn`, `blockWhenSettingIsSupported`.
- `maximum_allowed_device_threat_level` - (Optional) Maximum allowed device threat level. Possible values: `notConfigured`, `secured`, `low`, `medium`, `high`. Defaults to `notConfigured`.
- `mobile_threat_defense_remediation_action` - (Optional) Action if threat threshold is not met. Possible values: `block`, `wipe`. Defaults to `block`.
- `minimum_required_sdk_version` - (Optional) Minimum SDK version required to access company data.
- `minimum_wipe_sdk_version` - (Optional) Minimum SDK version before app is wiped.
- `minimum_required_os_version` - (Optional) Minimum OS version required to access company data.
- `minimum_warning_os_version` - (Optional) Minimum OS version before warning is shown.
- `minimum_wipe_os_version` - (Optional) Minimum OS version before app is wiped.
- `minimum_required_app_version` - (Optional) Minimum app version required to access company data.
- `minimum_warning_app_version` - (Optional) Minimum app version before warning is shown.
- `minimum_wipe_app_version` - (Optional) Minimum app version before app is wiped.
- `maximum_required_os_version` - (Optional) Maximum OS version allowed to access company data.
- `maximum_warning_os_version` - (Optional) Maximum OS version before warning is shown.
- `maximum_wipe_os_version` - (Optional) Maximum OS version before app is wiped.
- `period_offline_before_wipe_is_enforced` - (Optional) Time offline before managed data is wiped. ISO 8601 duration. Defaults to `P90D`.
- `period_offline_before_access_check` - (Optional) Time offline before access is checked. ISO 8601 duration. Defaults to `P30D`.

## Attribute Reference

- `id` - The unique identifier of the policy.
- `created_date_time` - The date and time the policy was created.
- `last_modified_date_time` - The date and time the policy was last modified.
- `version` - Version of the entity.
- `is_assigned` - Whether the policy is deployed to any groups.
- `deployed_app_count` - Number of apps the policy is deployed to.

## Import

```shell
terraform import microsoft365_graph_beta_device_and_app_management_windows_managed_app_protection.example <policy-id>
```