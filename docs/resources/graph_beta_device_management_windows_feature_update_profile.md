---
page_title: "microsoft365_graph_beta_device_management_windows_feature_update_profile Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages Windows feature update profiles using the /deviceManagement/windowsFeatureUpdateProfiles endpoint. Feature update profiles control major Windows version deployments (like Windows 11 24H2) with rollout scheduling, device eligibility rules, and deployment timing to ensure controlled OS upgrades across managed devices.
---

# microsoft365_graph_beta_device_management_windows_feature_update_profile (Resource)

Manages Windows feature update profiles using the `/deviceManagement/windowsFeatureUpdateProfiles` endpoint. Feature update profiles control major Windows version deployments (like Windows 11 24H2) with rollout scheduling, device eligibility rules, and deployment timing to ensure controlled OS upgrades across managed devices.

## Microsoft Documentation

- [windowsFeatureUpdateProfile resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsfeatureupdateprofile?view=graph-rest-beta)
- [Create windowsFeatureUpdateProfile](https://learn.microsoft.com/en-us/graph/api/intune-softwareupdate-windowsfeatureupdateprofile-create?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.14.1-alpha | Experimental | Initial release |

## Example Usage

```terraform
resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "example" {
  display_name                                            = "Windows 11 22H2 Deployment x"
  description                                             = "Feature update profile for Windows 11 22H2"
  feature_update_version                                  = "Windows 11, version 22H2"
  install_latest_windows10_on_windows11_ineligible_device = true
  install_feature_updates_optional                        = true
  role_scope_tag_ids                                      = ["8", "9"]

  // rollout_settings = Make update available gradually
  rollout_settings = {
    offer_start_date_time_in_utc = "2025-05-01T00:00:00Z"
    offer_end_date_time_in_utc   = "2025-06-30T23:59:59Z"
    offer_interval_in_days       = 7
  }

  // Optional assignment blocks
  assignment {
    target = "include"
    group_ids = [
      "11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555"
    ]
  }

  assignment {
    target = "exclude"
    group_ids = [
      "11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555"
    ]
  }

  # Optional timeout block
  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }

}

resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "example_2" {
  display_name                                            = "Windows 11 22H2 Deployment y"
  description                                             = "Feature update profile for Windows 11 22H2"
  feature_update_version                                  = "Windows 11, version 22H2"
  install_latest_windows10_on_windows11_ineligible_device = true
  install_feature_updates_optional                        = true
  role_scope_tag_ids                                      = ["8", "9"]

  // rollout_settings = Make update available on a specific date
  rollout_settings = {
    offer_start_date_time_in_utc = "2025-05-01T00:00:00Z"
  }

  // Optional assignment blocks
  assignment {
    target = "include"
    group_ids = [
      //"11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555"
    ]
  }

  assignment {
    target = "exclude"
    group_ids = [
      "11111111-2222-3333-4444-555555555555",
      //"11111111-2222-3333-4444-555555555555"
    ]
  }

  # Optional timeout block
  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

resource "microsoft365_graph_beta_device_management_windows_feature_update_profile" "example_3" {
  display_name                                            = "Windows 11 22H2 Deployment z"
  description                                             = "Feature update profile for Windows 11 22H2"
  feature_update_version                                  = "Windows 11, version 22H2"
  install_latest_windows10_on_windows11_ineligible_device = false
  install_feature_updates_optional                        = true
  role_scope_tag_ids                                      = ["8", "9"]

  // include no rollout_settings block to make Make update available as soon as possible

  // Optional assignment blocks
  assignment {
    target = "include"
    group_ids = [
      //"11111111-2222-3333-4444-555555555555",
      "11111111-2222-3333-4444-555555555555"
    ]
  }

  assignment {
    target = "exclude"
    group_ids = [
      "11111111-2222-3333-4444-555555555555",
      //"11111111-2222-3333-4444-555555555555"
    ]
  }

  # Optional timeout block
  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name of the profile.
- `feature_update_version` (String) The feature update version that will be deployed to the devices targeted by this profile. Valid values are: "Windows 11, version 24H2", "Windows 11, version 23H2", "Windows 11, version 22H2", "Windows 10, version 22H2". By selecting this Feature update to deploy you are agreeing that when applying this operating system to a device either (1) the applicable Windows license was purchased though volume licensing, or (2) that you are authorized to bind your organization and are accepting on its behalf the relevant Microsoft Software License Terms to be found here https://go.microsoft.com/fwlink/?linkid=2171206.

### Optional

- `assignment` (Block List) Assignments for Windows Quality Update policies, specifying groups to include or exclude. (see [below for nested schema](#nestedblock--assignment))
- `description` (String) The description of the profile which is specified by the user.
- `install_feature_updates_optional` (Boolean) If true, the Windows 11 update will become optional
- `install_latest_windows10_on_windows11_ineligible_device` (Boolean) Specifies whether Windows 10 devices that are not eligible for Windows 11 are offered the latest Windows 10 feature updates. Changes to this field require the resource to be replaced.
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `rollout_settings` (Attributes) The windows update rollout settings, including offer start date time, offer end date time, and days between each set of offers. (see [below for nested schema](#nestedatt--rollout_settings))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date time that the profile was created.
- `deployable_content_display_name` (String) Friendly display name of the quality update profile deployable content
- `end_of_support_date` (String) The last supported date for a feature update
- `id` (String) The Identifier of the entity.
- `last_modified_date_time` (String) The date time that the profile was last modified.

<a id="nestedblock--assignment"></a>
### Nested Schema for `assignment`

Required:

- `group_ids` (Set of String) Set of Microsoft Entra ID group IDs to apply for this assignment.
- `target` (String) Specifies whether the assignment is 'include' or 'exclude'.


<a id="nestedatt--rollout_settings"></a>
### Nested Schema for `rollout_settings`

Optional:

- `offer_end_date_time_in_utc` (String) The UTC offer end date time of the rollout.
- `offer_interval_in_days` (Number) The number of days between each set of offers.
- `offer_start_date_time_in_utc` (String) The UTC offer start date time of the rollout.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
# {resource_id}
terraform import microsoft365_graph_beta_device_and_app_management_windows_feature_update_profile.example windows-feature-update-profile-id
```