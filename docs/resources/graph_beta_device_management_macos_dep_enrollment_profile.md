---
page_title: "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages a macOS Automated Device Enrollment (DEP/ADE) enrollment profile using the /deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles endpoint with the #microsoft.graph.depMacOSEnrollmentProfile OData type. This profile drives automated (low-touch) macOS enrollment: skipping Setup Assistant panes, auto-creating the local admin account, and gating the desktop until MDM configuration finishes (await_device_configured). Note: fully hands-off ("zero-touch") provisioning is only approached when enrolling without user affinity (requires_user_authentication = false) over a wired network; Apple keeps some early Setup Assistant panes (such as network and region/language) non-skippable, so at least minimal physical interaction may still be required on a freshly wiped device.
---

# microsoft365_graph_beta_device_management_macos_dep_enrollment_profile (Resource)

Manages a macOS Automated Device Enrollment (DEP/ADE) enrollment profile using the `/deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles` endpoint with the `#microsoft.graph.depMacOSEnrollmentProfile` OData type. This profile drives automated (low-touch) macOS enrollment: skipping Setup Assistant panes, auto-creating the local admin account, and gating the desktop until MDM configuration finishes (`await_device_configured`). Note: fully hands-off ("zero-touch") provisioning is only approached when enrolling without user affinity (`requires_user_authentication = false`) over a wired network; Apple keeps some early Setup Assistant panes (such as network and region/language) non-skippable, so at least minimal physical interaction may still be required on a freshly wiped device.

## Microsoft Documentation

- [depMacOSEnrollmentProfile resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-depmacosenrollmentprofile?view=graph-rest-beta)
- [depEnrollmentBaseProfile resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-depenrollmentbaseprofile?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `DeviceManagementServiceConfig.Read.All`
- `DeviceManagementServiceConfig.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.56.0-alpha | Experimental | Initial release |

## Known Issues

### Retrieval-triggered password rotation

Live `curl.exe` checks on 2026-11-06 found that the enrollment-profile endpoint returns
`depProfileDelayAutoRotationSetting: null` after POST, PATCH, individual GET and filtered LIST,
even when retrieval-triggered rotation values are submitted. The rotation period in days persists.
The live API metadata still lists the original retrieval fields; no equivalent replacement
was found in the ADE settings catalog.

The examples and acceptance lifecycle therefore configure only `auto_rotation_period_in_days`.
Configuring `on_retrieval_auto_rotate_password_enabled` or
`on_retrieval_delay_auto_rotate_password_in_hours` on an affected tenant can produce an
inconsistent-result error after apply. Both attributes are deprecated in the Terraform schema. Their state names and types remain
unchanged, so existing state requires no field migration. Scheduled rotation has different
semantics and is not a replacement for retrieval-triggered rotation.

## Example Usage

The four examples correspond to configurations captured from the Intune enrollment-profile
wizard. Modern and legacy authentication both use user affinity. The acceptance lifecycle
updates one profile through minimum, modern authentication, legacy authentication, without
user affinity, and back to minimum. All 27 current Setup Assistant skip keys are supported.

Intune also offers newer enrollment policies backed by `configurationPolicies`. Those are
managed by `microsoft365_graph_beta_device_management_macos_device_enrollment_policy`;
these examples configure the existing DEP enrollment-profile endpoint.

### Minimum

```terraform
# Configuration verified against the Intune enrollment-profile wizard.
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "minimal" {
  display_name                                               = "macOS DEP - minimal"
  description                                                = "macOS DEP enrollment with minimal"
  requires_user_authentication                               = false
  configuration_web_url                                      = false
  supervised_mode_enabled                                    = true
  support_department                                         = "IT Support"
  is_mandatory                                               = true
  profile_removal_disabled                                   = false
  location_disabled                                          = false
  support_phone_number                                       = "+1-555-0100"
  restore_blocked                                            = false
  apple_id_disabled                                          = false
  terms_and_conditions_disabled                              = false
  touch_id_disabled                                          = false
  apple_pay_disabled                                         = false
  siri_disabled                                              = false
  diagnostics_disabled                                       = false
  privacy_pane_disabled                                      = false
  screen_time_screen_disabled                                = false
  file_vault_disabled                                        = false
  icloud_diagnostics_disabled                                = false
  icloud_storage_disabled                                    = false
  choose_your_lock_screen_disabled                           = false
  display_tone_setup_disabled                                = false
  registration_disabled                                      = false
  accessibility_screen_disabled                              = false
  auto_unlock_with_watch_disabled                            = false
  await_device_configured                                    = false
  enable_authentication_via_company_portal                   = false
  require_company_portal_on_setup_assistant_enrolled_devices = false
  skip_primary_setup_account_creation                        = true
  set_primary_setup_account_as_regular_user                  = true
  dont_auto_populate_primary_account_info                    = true
  primary_account_full_name                                  = ""
  primary_account_user_name                                  = ""
  enable_restrict_editing                                    = false
  admin_account_user_name                                    = ""
  admin_account_full_name                                    = ""
  hide_admin_account                                         = false
  appearance_disabled                                        = false
  enable_lockdown_mode_disabled                              = false
  software_update_disabled                                   = false
  additional_privacy_settings_disabled                       = false
  update_completed_disabled                                  = false
  terms_of_address_disabled                                  = false
  intelligence_disabled                                      = false
  os_showcase_disabled                                       = false
  app_store_disabled                                         = false
  welcome_screen_disabled                                    = false
  pass_code_disabled                                         = false
  zoom_disabled                                              = false

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

### Modern authentication with user affinity

```terraform
# Configuration verified against the Intune enrollment-profile wizard.
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "modern_auth" {
  display_name                                               = "macOS DEP - modern auth"
  description                                                = "macOS DEP enrollment with modern auth"
  requires_user_authentication                               = true
  configuration_web_url                                      = true
  supervised_mode_enabled                                    = true
  support_department                                         = "IT Support"
  is_mandatory                                               = true
  profile_removal_disabled                                   = true
  location_disabled                                          = true
  support_phone_number                                       = "+1-555-0100"
  restore_blocked                                            = true
  apple_id_disabled                                          = true
  terms_and_conditions_disabled                              = true
  touch_id_disabled                                          = true
  apple_pay_disabled                                         = true
  siri_disabled                                              = true
  diagnostics_disabled                                       = true
  privacy_pane_disabled                                      = true
  screen_time_screen_disabled                                = true
  file_vault_disabled                                        = true
  icloud_diagnostics_disabled                                = true
  icloud_storage_disabled                                    = true
  choose_your_lock_screen_disabled                           = true
  display_tone_setup_disabled                                = true
  registration_disabled                                      = true
  accessibility_screen_disabled                              = true
  auto_unlock_with_watch_disabled                            = true
  await_device_configured                                    = true
  enable_authentication_via_company_portal                   = false
  require_company_portal_on_setup_assistant_enrolled_devices = false
  skip_primary_setup_account_creation                        = false
  set_primary_setup_account_as_regular_user                  = true
  dont_auto_populate_primary_account_info                    = false
  primary_account_full_name                                  = "Acceptance User"
  primary_account_user_name                                  = "accuser"
  enable_restrict_editing                                    = true
  admin_account_user_name                                    = "accadmin"
  admin_account_full_name                                    = "Acceptance Administrator"
  hide_admin_account                                         = true
  appearance_disabled                                        = true
  enable_lockdown_mode_disabled                              = true
  software_update_disabled                                   = true
  additional_privacy_settings_disabled                       = true
  update_completed_disabled                                  = true
  terms_of_address_disabled                                  = true
  intelligence_disabled                                      = true
  os_showcase_disabled                                       = true
  app_store_disabled                                         = true
  welcome_screen_disabled                                    = false
  pass_code_disabled                                         = false
  zoom_disabled                                              = false

  admin_account_password_rotation = {
    auto_rotation_period_in_days = 90
  }

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

### Legacy authentication with user affinity

```terraform
# Configuration verified against the Intune enrollment-profile wizard.
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "legacy_auth" {
  display_name                                               = "macOS DEP - legacy auth"
  description                                                = "macOS DEP enrollment with legacy auth"
  requires_user_authentication                               = true
  configuration_web_url                                      = false
  supervised_mode_enabled                                    = true
  support_department                                         = "IT Support"
  is_mandatory                                               = true
  profile_removal_disabled                                   = true
  location_disabled                                          = true
  support_phone_number                                       = "+1-555-0100"
  restore_blocked                                            = true
  apple_id_disabled                                          = true
  terms_and_conditions_disabled                              = true
  touch_id_disabled                                          = true
  apple_pay_disabled                                         = true
  siri_disabled                                              = true
  diagnostics_disabled                                       = true
  privacy_pane_disabled                                      = true
  screen_time_screen_disabled                                = true
  file_vault_disabled                                        = true
  icloud_diagnostics_disabled                                = true
  icloud_storage_disabled                                    = true
  choose_your_lock_screen_disabled                           = true
  display_tone_setup_disabled                                = true
  registration_disabled                                      = true
  accessibility_screen_disabled                              = true
  auto_unlock_with_watch_disabled                            = true
  await_device_configured                                    = true
  enable_authentication_via_company_portal                   = false
  require_company_portal_on_setup_assistant_enrolled_devices = false
  skip_primary_setup_account_creation                        = false
  set_primary_setup_account_as_regular_user                  = false
  dont_auto_populate_primary_account_info                    = true
  primary_account_full_name                                  = ""
  primary_account_user_name                                  = ""
  enable_restrict_editing                                    = false
  admin_account_user_name                                    = "accadmin"
  admin_account_full_name                                    = "Acceptance Administrator"
  hide_admin_account                                         = false
  appearance_disabled                                        = true
  enable_lockdown_mode_disabled                              = true
  software_update_disabled                                   = true
  additional_privacy_settings_disabled                       = true
  update_completed_disabled                                  = true
  terms_of_address_disabled                                  = true
  intelligence_disabled                                      = true
  os_showcase_disabled                                       = true
  app_store_disabled                                         = true
  welcome_screen_disabled                                    = false
  pass_code_disabled                                         = false
  zoom_disabled                                              = false

  admin_account_password_rotation = {
    auto_rotation_period_in_days = 180
  }

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

### Without user affinity

```terraform
# Configuration verified against the Intune enrollment-profile wizard.
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "without_user_affinity" {
  display_name                                               = "macOS DEP - without user affinity"
  description                                                = "macOS DEP enrollment with without user affinity"
  requires_user_authentication                               = false
  configuration_web_url                                      = false
  supervised_mode_enabled                                    = true
  support_department                                         = "IT Support"
  is_mandatory                                               = true
  profile_removal_disabled                                   = true
  location_disabled                                          = true
  support_phone_number                                       = "+1-555-0100"
  restore_blocked                                            = true
  apple_id_disabled                                          = true
  terms_and_conditions_disabled                              = true
  touch_id_disabled                                          = true
  apple_pay_disabled                                         = true
  siri_disabled                                              = true
  diagnostics_disabled                                       = true
  privacy_pane_disabled                                      = true
  screen_time_screen_disabled                                = true
  file_vault_disabled                                        = true
  icloud_diagnostics_disabled                                = true
  icloud_storage_disabled                                    = true
  choose_your_lock_screen_disabled                           = true
  display_tone_setup_disabled                                = true
  registration_disabled                                      = true
  accessibility_screen_disabled                              = true
  auto_unlock_with_watch_disabled                            = true
  await_device_configured                                    = true
  enable_authentication_via_company_portal                   = false
  require_company_portal_on_setup_assistant_enrolled_devices = false
  skip_primary_setup_account_creation                        = false
  set_primary_setup_account_as_regular_user                  = false
  dont_auto_populate_primary_account_info                    = true
  primary_account_full_name                                  = ""
  primary_account_user_name                                  = ""
  enable_restrict_editing                                    = false
  admin_account_user_name                                    = ""
  admin_account_full_name                                    = ""
  hide_admin_account                                         = false
  appearance_disabled                                        = true
  enable_lockdown_mode_disabled                              = true
  software_update_disabled                                   = true
  additional_privacy_settings_disabled                       = true
  update_completed_disabled                                  = true
  terms_of_address_disabled                                  = true
  intelligence_disabled                                      = true
  os_showcase_disabled                                       = true
  app_store_disabled                                         = true
  welcome_screen_disabled                                    = false
  pass_code_disabled                                         = false
  zoom_disabled                                              = false

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Name of the profile displayed in Intune.

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `accessibility_screen_disabled` (Boolean) Indicates if the Accessibility screen is disabled.
- `additional_privacy_settings_disabled` (Boolean) Hide the AdditionalPrivacySettings Setup Assistant pane. Maps to the `AdditionalPrivacySettings` enabled skip key.
- `admin_account_full_name` (String) The full name for the auto-created local admin account.
- `admin_account_password` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) The password for the auto-created local admin account. This is a write-only, sensitive value and is never returned by the Graph API.
- `admin_account_password_rotation` (Attributes) Settings for local admin account password automatic rotation (depProfileAdminAccountPasswordRotationSetting). (see [below for nested schema](#nestedatt--admin_account_password_rotation))
- `admin_account_user_name` (String) The user name (short name) for the auto-created local admin account.
- `app_store_disabled` (Boolean) Hide the AppStore Setup Assistant pane. Maps to the `AppStore` enabled skip key.
- `appearance_disabled` (Boolean) Hide the Appearance Setup Assistant pane. Maps to the `Appearance` enabled skip key.
- `apple_id_disabled` (Boolean) Indicates if the Apple ID setup pane is disabled.
- `apple_pay_disabled` (Boolean) Indicates if the Apple Pay setup pane is disabled.
- `auto_advance_setup_enabled` (Boolean) Indicates if Setup Assistant will automatically advance through its screens.
- `auto_unlock_with_watch_disabled` (Boolean) Indicates if the Unlock With Watch screen is disabled.
- `await_device_configured` (Boolean) Indicates if the device will need to wait for configured confirmation (the desktop is gated until MDM configuration finishes). Maps to `waitForDeviceConfiguredConfirmation`.
- `choose_your_lock_screen_disabled` (Boolean) Indicates if the Choose Your Lock Screen screen is disabled.
- `configuration_web_url` (Boolean) Indicates if the admin-assisted setup assistant login (web-based authentication) URL is used. Cannot be true when `use_platform_sso_during_setup_assistant` is true.
- `dep_onboarding_settings_id` (String) Identifier of the parent depOnboardingSetting (Apple ABM/ASM ADE token) that contains this macOS DEP enrollment profile. If omitted, the provider auto-resolves it by listing `/deviceManagement/depOnboardingSettings` and selecting the Apple ADE/ABM token (`tokenType` = `dep`) or ASM token (`tokenType` = `appleSchoolManager`), ignoring Apple Configurator tokens. Set this explicitly only if your tenant has more than one Apple DEP token (auto-resolution is ambiguous in that case).
- `description` (String) Description of the profile. Maximum length is 1500 characters.
- `device_name_template` (String) Sets a literal or name pattern for the device name.
- `diagnostics_disabled` (Boolean) Indicates if the Diagnostics setup pane is disabled.
- `display_tone_setup_disabled` (Boolean) Indicates if the DisplayTone setup screen is disabled.
- `dont_auto_populate_primary_account_info` (Boolean) Indicates whether Setup Assistant will auto-populate the primary account information.
- `enable_authentication_via_company_portal` (Boolean) Indicates to authenticate with the Company Portal instead of Apple Setup Assistant.
- `enable_lockdown_mode_disabled` (Boolean) Hide the EnableLockdownMode Setup Assistant pane. Maps to the `EnableLockdownMode` enabled skip key.
- `enable_restrict_editing` (Boolean) Indicates whether the user will be blocked from editing the account.
- `enrollment_time_azure_ad_group_ids` (Set of String) List of enrollment-time Microsoft Entra (Azure AD) group GUIDs to be associated with the profile.
- `file_vault_disabled` (Boolean) Indicates if FileVault is disabled.
- `hide_admin_account` (Boolean) Indicates whether the local admin account should be hidden.
- `icloud_diagnostics_disabled` (Boolean) Indicates if the iCloud Analytics screen is disabled.
- `icloud_storage_disabled` (Boolean) Indicates if the iCloud Documents and Desktop screen is disabled.
- `intelligence_disabled` (Boolean) Hide the Intelligence Setup Assistant pane. Maps to the `Intelligence` enabled skip key.
- `is_mandatory` (Boolean) Indicates if the profile is mandatory.
- `location_disabled` (Boolean) Indicates if the Location service setup pane is disabled.
- `os_showcase_disabled` (Boolean) Hide the OSShowcase Setup Assistant pane. Maps to the `OSShowcase` enabled skip key.
- `pass_code_disabled` (Boolean) Indicates if the Passcode setup pane is disabled.
- `primary_account_full_name` (String) The full name for the primary account.
- `primary_account_user_name` (String) The account name (short name) for the primary account.
- `privacy_pane_disabled` (Boolean) Indicates if the Privacy screen is disabled.
- `profile_removal_disabled` (Boolean) Indicates if the profile removal option is disabled.
- `registration_disabled` (Boolean) Indicates if registration is disabled.
- `request_requires_network_tether` (Boolean) Indicates if the device is network-tethered to run the command.
- `require_company_portal_on_setup_assistant_enrolled_devices` (Boolean) Indicates that the Company Portal is required on setup assistant enrolled devices.
- `requires_user_authentication` (Boolean) Indicates if the profile requires user authentication.
- `restore_blocked` (Boolean) Indicates if the Restore setup pane is blocked.
- `screen_time_screen_disabled` (Boolean) Indicates if the Screen Time setup screen is disabled.
- `set_primary_setup_account_as_regular_user` (Boolean) Indicates whether Setup Assistant will set the primary account as a regular (non-admin) user.
- `siri_disabled` (Boolean) Indicates if the Siri setup pane is disabled.
- `skip_primary_setup_account_creation` (Boolean) Indicates whether Setup Assistant will skip the user interface for primary account setup.
- `software_update_disabled` (Boolean) Hide the SoftwareUpdate Setup Assistant pane. Maps to the `SoftwareUpdate` enabled skip key.
- `supervised_mode_enabled` (Boolean) Supervised mode. True to enable, false otherwise.
- `support_department` (String) Support department information.
- `support_phone_number` (String) Support phone number.
- `terms_and_conditions_disabled` (Boolean) Indicates if the 'Terms and Conditions' setup pane is disabled.
- `terms_of_address_disabled` (Boolean) Hide the TermsOfAddress Setup Assistant pane. Maps to the `TermsOfAddress` enabled skip key.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `touch_id_disabled` (Boolean) Indicates if the Touch ID setup pane is disabled.
- `update_completed_disabled` (Boolean) Hide the UpdateCompleted Setup Assistant pane. Maps to the `UpdateCompleted` enabled skip key.
- `use_platform_sso_during_setup_assistant` (Boolean) Indicates whether Platform SSO is used as part of device enrollment during Setup Assistant. Cannot be true when `configuration_web_url` is true.
- `welcome_screen_disabled` (Boolean) Indicates if the Get Started (Welcome) setup pane is disabled. macOS 15 and later.
- `zoom_disabled` (Boolean) Indicates if the Zoom setup pane is disabled.

### Read-Only

- `configuration_endpoint_url` (String) Configuration endpoint url to use for enrollment. Generated by Intune.
- `enabled_skip_keys` (Set of String) Computed, read-only set of Setup Assistant skip keys (Apple `SkipKeys`) that the provider sends to Graph. This is derived from the individual `*_disabled` boolean attributes; do not set it directly. Includes the Privacy and Registration panes supported by Intune.
- `id` (String) The unique identifier of the enrollment profile. Format is `{depOnboardingSettingsId}_{profileId}`.
- `is_default` (Boolean) Indicates if this is the default profile.

<a id="nestedatt--admin_account_password_rotation"></a>
### Nested Schema for `admin_account_password_rotation`

Optional:

- `auto_rotation_period_in_days` (Number) The number of days between automatic admin account password rotations.
- `on_retrieval_auto_rotate_password_enabled` (Boolean, Deprecated) Deprecated: Microsoft Graph does not persist this setting. Remove it from configuration. Use `auto_rotation_period_in_days` for scheduled rotation; there is no verified retrieval-triggered replacement.
- `on_retrieval_delay_auto_rotate_password_in_hours` (Number, Deprecated) Deprecated: Microsoft Graph does not persist this setting. Remove it from configuration. Use `auto_rotation_period_in_days` for scheduled rotation; there is no verified retrieval-triggered replacement.


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
#!/bin/bash

# {depOnboardingSettingsId}_{enrollmentProfileId}
terraform import microsoft365_graph_beta_device_management_macos_dep_enrollment_profile.example 00000000-0000-0000-0000-000000000000_11111111-1111-1111-1111-111111111111
```
