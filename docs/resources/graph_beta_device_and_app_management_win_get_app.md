---
page_title: "microsoft365_graph_beta_device_and_app_management_win_get_app Resource - terraform-provider-microsoft365"
subcategory: "Intune"
description: |-
  Manages an Intune Microsoft Store app (new) resource aka winget, using the mobileapps graph beta API.
---

# microsoft365_graph_beta_device_and_app_management_win_get_app (Resource)

Manages an Intune Microsoft Store app (new) resource aka winget, using the mobileapps graph beta API.

## Example Usage

```terraform
resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "whatsapp" {
  package_identifier              = "9NKSQGP7F2NH" # The unique identifier for the app obtained from msft app store
  automatically_generate_metadata = true

  # Install experience settings
  install_experience = {
    run_as_account = "user" # Can be 'system' or 'user'
  }

  role_scope_tag_ids = ["0"]

  # Optional fields
  is_featured             = true
  privacy_information_url = "https://privacy.example.com"
  information_url         = "https://info.example.com"
  owner                   = "example-owner"
  developer               = "example-developer"
  notes                   = "Some relevant notes for this app."

  # Optional: Define custom timeouts
  timeouts = {
    create = "10m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "Adobe_Acrobat_Reader_DC" {
  package_identifier              = "xpdp273c0xhqh2" # The unique identifier for the app obtained from msft app store
  automatically_generate_metadata = false
  display_name                    = "Adobe Acrobat Reader DC"
  description                     = "Adobe Acrobat Reader DC is the free, trusted standard for viewing, printing, signing, and annotating PDFs. It's the only PDF viewer that can open and interact with all types of PDF content – including forms and multimedia."
  publisher                       = "Adobe Inc."
  large_icon = {
    type  = "image/png"
    value = filebase64("${path.module}/Adobe_Reader_XI_icon.png")
  }
  # Install experience settings
  install_experience = {
    run_as_account = "user" # Can be 'system' or 'user'
  }

  role_scope_tag_ids = ["0"]

  # Optional fields
  is_featured             = true
  privacy_information_url = "https://privacy.example.com"
  information_url         = "https://info.example.com"
  owner                   = "example-owner"
  developer               = "example-developer"
  notes                   = "Some relevant notes for this app."

  # Optional: Define custom timeouts
  timeouts = {
    create = "10m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_and_app_management_win_get_app" "visual_studio_code" {
  package_identifier = "XP9KHM4BK9FZ7Q" # The unique identifier for the app obtained from msft app store

  # Install experience settings
  install_experience = {
    run_as_account = "user" # Can be 'system' or 'user'
  }

  role_scope_tag_ids = ["0"]

  # Optional fields
  is_featured             = true
  privacy_information_url = "https://privacy.example.com"
  information_url         = "https://info.example.com"
  owner                   = "example-owner"
  developer               = "example-developer"
  notes                   = "Some relevant notes for this app."

  # Optional: Define custom timeouts
  timeouts = {
    create = "10m"
    update = "10m"
    delete = "5m"
  }

  # App assignments configuration
  assignments = {
    id            = "assignment_id" # Read-only, typically auto-generated
    mobile_app_id = "app_id_value"  # The ID of the app being assigned

    mobile_app_assignments = [
      # 2 Assignments for "available" intent
      {
        intent = "required" # Possible values: available, required, uninstall, availableWithoutEnrollment
        source = "direct"   # Possible values: direct, policySets

        target = {
          target_type                                      = "groupAssignmentTarget" # Possible values: groupAssignmentTarget, allLicensedUsersAssignmentTarget, etc.
          group_id                                         = "group_id_value"
          device_and_app_management_assignment_filter_id   = "filter_id_value"
          device_and_app_management_assignment_filter_type = "include" # Possible values: include, exclude, none
          is_exclusion_group                               = false
        }

        settings = {
          notifications = "showAll" # Possible values: showAll, showReboot, hideAll
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
          restart_settings = {
            grace_period_in_minutes                         = 15
            countdown_display_before_restart_in_minutes     = 5
            restart_notification_snooze_duration_in_minutes = 10
          }
        }
      },
      {
        id        = "assignment_2"
        intent    = "available"
        source    = "policySets"
        source_id = "source_id_2"

        target = {
          target_type                                      = "microsoft.graph.groupAssignmentTarget"
          group_id                                         = "group_id_available_2"
          device_and_app_management_assignment_filter_id   = "filter_id_2"
          device_and_app_management_assignment_filter_type = "include"
          is_exclusion_group                               = false
        }

        settings = {
          notifications = "showReboot"
        }
      },

      # 2 Assignments for "required" intent
      {
        id        = "assignment_3"
        intent    = "required"
        source    = "direct"
        source_id = "source_id_3"

        target = {
          target_type                                      = "microsoft.graph.groupAssignmentTarget"
          group_id                                         = "group_id_required_1"
          device_and_app_management_assignment_filter_id   = "filter_id_3"
          device_and_app_management_assignment_filter_type = "exclude"
          is_exclusion_group                               = false
        }

        settings = {
          notifications = "hideAll"
          install_time_settings = {
            use_local_time     = true
            deadline_date_time = "2024-12-31T23:59:59Z"
          }
        }
      },
      {
        id        = "assignment_4"
        intent    = "required"
        source    = "policySets"
        source_id = "source_id_4"

        target = {
          target_type                                      = "microsoft.graph.groupAssignmentTarget"
          group_id                                         = "group_id_required_2"
          device_and_app_management_assignment_filter_id   = "filter_id_4"
          device_and_app_management_assignment_filter_type = "include"
          is_exclusion_group                               = true
        }

        settings = {
          notifications = "showAll"
        }
      },

      # 2 Assignments for "uninstall" intent
      {
        id        = "assignment_5"
        intent    = "uninstall"
        source    = "direct"
        source_id = "source_id_5"

        target = {
          target_type                                      = "microsoft.graph.groupAssignmentTarget"
          group_id                                         = "group_id_uninstall_1"
          device_and_app_management_assignment_filter_id   = "filter_id_5"
          device_and_app_management_assignment_filter_type = "none"
          is_exclusion_group                               = false
        }

        settings = {
          notifications = "showReboot"
          restart_settings = {
            grace_period_in_minutes                         = 15
            countdown_display_before_restart_in_minutes     = 5
            restart_notification_snooze_duration_in_minutes = 10
          }
        }
      },
      {
        id        = "assignment_6"
        intent    = "uninstall"
        source    = "policySets"
        source_id = "source_id_6"

        target = {
          target_type                                      = "microsoft.graph.groupAssignmentTarget"
          group_id                                         = "group_id_uninstall_2"
          device_and_app_management_assignment_filter_id   = "filter_id_6"
          device_and_app_management_assignment_filter_type = "exclude"
          is_exclusion_group                               = true
        }

        settings = {
          notifications = "hideAll"
        }
      }
    ]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `automatically_generate_metadata` (Boolean) When set to `true`, the provider will automatically fetch metadata from the Microsoft Store for Business using the package identifier. This will populate the `display_name`, `description`, `publisher`, and 'icon' fields.
- `install_experience` (Attributes) The install experience settings associated with this application. (see [below for nested schema](#nestedatt--install_experience))
- `package_identifier` (String) The **unique package identifier** for the WinGet/Microsoft Store app from the storefront.

For example, for the app Microsoft Edge Browser URL [https://apps.microsoft.com/detail/xpfftq037jwmhs?hl=en-us&gl=US](https://apps.microsoft.com/detail/xpfftq037jwmhs?hl=en-us&gl=US), the package identifier is `xpfftq037jwmhs`.

**Important notes:**
- This identifier is **required** at creation time.
- It **cannot be modified** for existing Terraform-deployed WinGet/Microsoft Store apps.

attempting to modify this value will result in a failed request.

### Optional

- `assignments` (Attributes) Configuration for Mobile App Assignment, including settings and targets for Microsoft Intune. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) A detailed description of the WinGet/ Microsoft Store for Business app.This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.
- `developer` (String) The developer of the app.
- `display_name` (String) The title of the WinGet app imported from the Microsoft Store for Business.This field value must match the expected title of the app in the Microsoft Store for Business associated with the `package_identifier`.This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.
- `information_url` (String) The more information Url.
- `is_featured` (Boolean) The value indicating whether the app is marked as featured by the admin.
- `large_icon` (Attributes) The large icon for the WinGet app, automatically downloaded and set from the Microsoft Store for Business. (see [below for nested schema](#nestedatt--large_icon))
- `notes` (String) Notes for the app.
- `owner` (String) The owner of the app.
- `privacy_information_url` (String) The privacy statement Url.
- `publisher` (String) The publisher of the WinGet/ Microsoft Store for Business app.This field is automatically populated based on the package identifier when `automatically_generate_metadata` is set to true.
- `role_scope_tag_ids` (List of String) List of scope tag ids for this mobile app.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The date and time the app was created. This property is read-only.
- `dependent_app_count` (Number) The total number of dependencies the child app has. This property is read-only.
- `id` (String) The unique graph guid that identifies this resource.Assigned at time of resource creation. This property is read-only.
- `is_assigned` (Boolean) The value indicating whether the app is assigned to at least one group. This property is read-only.
- `last_modified_date_time` (String) The date and time the app was last modified. This property is read-only.
- `manifest_hash` (String) Hash of package metadata properties used to validate that the application matches the metadata in the source repository.
- `publishing_state` (String) The publishing state for the app. The app cannot be assigned unless the app is published. Possible values are: notPublished, processing, published.
- `superseded_app_count` (Number) The total number of apps this app is directly or indirectly superseded by. This property is read-only.
- `superseding_app_count` (Number) The total number of apps this app directly or indirectly supersedes. This property is read-only.
- `upload_state` (Number) The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.

<a id="nestedatt--install_experience"></a>
### Nested Schema for `install_experience`

Required:

- `run_as_account` (String) The account type (System or User) that actions should be run as on target devices. Required at creation time.


<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `mobile_app_id` (String) The ID of the mobile app associated with this assignment.

Optional:

- `mobile_app_assignments` (Attributes List) List of assignments for the mobile app. (see [below for nested schema](#nestedatt--assignments--mobile_app_assignments))

Read-Only:

- `id` (String) Key of the entity. This property is read-only.

<a id="nestedatt--assignments--mobile_app_assignments"></a>
### Nested Schema for `assignments.mobile_app_assignments`

Required:

- `intent` (String) The install intent defined by the admin. Possible values are: available, required, uninstall, availableWithoutEnrollment.
- `target` (Attributes) The target group assignment defined by the admin. (see [below for nested schema](#nestedatt--assignments--mobile_app_assignments--target))

Optional:

- `settings` (Attributes) The settings for target assignment defined by the admin. (see [below for nested schema](#nestedatt--assignments--mobile_app_assignments--settings))

Read-Only:

- `id` (String) Key of the assignment entity. This property is read-only.
- `source` (String) The resource type which is the source for the assignment. Possible values are: direct, policySets. This property is read-only.
- `source_id` (String) The identifier of the source of the assignment. This property is read-only.

<a id="nestedatt--assignments--mobile_app_assignments--target"></a>
### Nested Schema for `assignments.mobile_app_assignments.target`

Required:

- `target_type` (String) The type of target assignment. Possible values include groupAssignmentTarget, allLicensedUsersAssignmentTarget, etc.

Optional:

- `device_and_app_management_assignment_filter_id` (String) The ID of the filter for the target assignment.
- `device_and_app_management_assignment_filter_type` (String) The type of filter for the target assignment. Possible values are: none, include, exclude.
- `group_id` (String) The ID of the target group.
- `is_exclusion_group` (Boolean) Indicates whether this is an exclusion group.


<a id="nestedatt--assignments--mobile_app_assignments--settings"></a>
### Nested Schema for `assignments.mobile_app_assignments.settings`

Optional:

- `install_time_settings` (Attributes) Settings related to install time. (see [below for nested schema](#nestedatt--assignments--mobile_app_assignments--settings--install_time_settings))
- `notifications` (String) The notification settings for the assignment.
- `restart_settings` (Attributes) Settings related to restarts after installation. (see [below for nested schema](#nestedatt--assignments--mobile_app_assignments--settings--restart_settings))

<a id="nestedatt--assignments--mobile_app_assignments--settings--install_time_settings"></a>
### Nested Schema for `assignments.mobile_app_assignments.settings.install_time_settings`

Optional:

- `deadline_date_time` (String) The time at which the app should be installed.
- `use_local_time` (Boolean) Whether the local device time or UTC time should be used when determining the deadline times.


<a id="nestedatt--assignments--mobile_app_assignments--settings--restart_settings"></a>
### Nested Schema for `assignments.mobile_app_assignments.settings.restart_settings`

Optional:

- `countdown_display_before_restart_in_minutes` (Number) The number of minutes before the restart time to display the countdown dialog for pending restarts.
- `grace_period_in_minutes` (Number) The number of minutes to wait before restarting the device after an app installation.
- `restart_notification_snooze_duration_in_minutes` (Number) The number of minutes to snooze the restart notification dialog when the snooze button is selected.





<a id="nestedatt--large_icon"></a>
### Nested Schema for `large_icon`

Optional:

- `type` (String) The MIME type of the app's large icon, automatically populated based on the `package_identifier` when `automatically_generate_metadata` is true. Example: `image/png`
- `value` (String) The icon image to use for the winget app. This field is automatically populated based on the `package_identifier` when `automatically_generate_metadata` is set to true.


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
# Using the provider-default project ID, the import ID is:
# {resource_id}
terraform import microsoft365_graph_beta_device_and_app_win_get_app.example win-get-app-id
```

