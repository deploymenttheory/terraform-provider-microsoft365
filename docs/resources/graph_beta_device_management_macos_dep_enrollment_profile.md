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

## Example Usage

```terraform
# Example 1: True zero-touch macOS enrollment (enroll WITHOUT user affinity / "userless")
#
# This is the genuine zero-touch path: no user authentication, no Company Portal,
# the device provisions with no human interaction. Setup Assistant panes are skipped
# and the desktop is gated until MDM finishes (await_device_configured).
#
# NOTE: Local admin/user account auto-creation is NOT available in this flow. Per
# Microsoft, account creation requires user affinity + Setup Assistant auth +
# await_device_configured (see Example 2).
#
# Setup Assistant panes are skipped via the individual *_disabled booleans below; the
# provider derives the read-only enabled_skip_keys array from them.
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "zero_touch_userless" {
  display_name                 = "macOS DEP - Zero Touch (userless)"
  description                  = "Userless zero-touch macOS enrollment; skips Setup Assistant and gates the desktop until MDM finishes"
  requires_user_authentication = false
  supervised_mode_enabled      = true

  # Userless (and auto-advancing) profiles must be mandatory, or Graph rejects the request.
  is_mandatory = true

  # Gate the desktop until MDM configuration finishes (awaitDeviceConfigured)
  await_device_configured = true

  # Setup Assistant pane skip toggles (drive enabled_skip_keys)
  apple_id_disabled             = true
  apple_pay_disabled            = true
  terms_and_conditions_disabled = true
  diagnostics_disabled          = true
  display_tone_setup_disabled   = true
  siri_disabled                 = true
  file_vault_disabled           = true
  location_disabled             = true
  restore_blocked               = true
  screen_time_screen_disabled   = true
  icloud_storage_disabled       = true
  icloud_diagnostics_disabled   = true
  welcome_screen_disabled       = true

  # privacy_pane_disabled and registration_disabled work as boolean properties, but
  # their skip-key strings are rejected by Graph, so they are NOT added to
  # enabled_skip_keys (the provider handles this automatically).
  privacy_pane_disabled = true
  registration_disabled = true

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 2: User-affinity macOS enrollment with an auto-created local admin account (LAPS)
#
# Account auto-creation requires ALL of the following (enforced by Intune):
#   - user affinity (requires_user_authentication = true, Setup Assistant auth)
#   - await_device_configured = true
# The admin password is write-only/sensitive. When admin_account_password_rotation is
# set, Intune manages the password (LAPS-style) and rotates it automatically.
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "user_affinity_with_local_admin" {
  display_name                 = "macOS DEP - User Affinity + Local Admin"
  description                  = "Setup Assistant with modern authentication; auto-creates a hidden, LAPS-managed local admin account"
  requires_user_authentication = true
  supervised_mode_enabled      = true

  # Required for local account creation to take effect
  await_device_configured = true

  apple_id_disabled    = true
  diagnostics_disabled = true
  siri_disabled        = true
  restore_blocked      = true

  # Auto-create the local admin account (managed local user)
  admin_account_user_name = "ladmin"
  admin_account_full_name = "Local Administrator"
  admin_account_password  = var.local_admin_password # sensitive / write-only
  hide_admin_account      = true

  # Optional automatic admin password rotation (LAPS-style)
  admin_account_password_rotation = {
    auto_rotation_period_in_days                     = 30
    on_retrieval_auto_rotate_password_enabled        = true
    on_retrieval_delay_auto_rotate_password_in_hours = 24
  }

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example 3: Minimal profile
resource "microsoft365_graph_beta_device_management_macos_dep_enrollment_profile" "minimal" {
  display_name                 = "macOS DEP - Minimal"
  description                  = "Minimal macOS enrollment profile"
  requires_user_authentication = false
  is_mandatory                 = true # required because this profile is userless

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

variable "local_admin_password" {
  type      = string
  sensitive = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Name of the profile displayed in Intune.

### Optional

> **NOTE**: [Write-only arguments](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments) are supported in Terraform 1.11 and later.

- `accessibility_screen_disabled` (Boolean) Indicates if the Accessibility screen is disabled.
- `admin_account_full_name` (String) The full name for the auto-created local admin account.
- `admin_account_password` (String, Sensitive, [Write-only](https://developer.hashicorp.com/terraform/language/resources/ephemeral#write-only-arguments)) The password for the auto-created local admin account. This is a write-only, sensitive value and is never returned by the Graph API.
- `admin_account_password_rotation` (Attributes) Settings for local admin account password automatic rotation (depProfileAdminAccountPasswordRotationSetting). (see [below for nested schema](#nestedatt--admin_account_password_rotation))
- `admin_account_user_name` (String) The user name (short name) for the auto-created local admin account.
- `apple_id_disabled` (Boolean) Indicates if the Apple ID setup pane is disabled.
- `apple_pay_disabled` (Boolean) Indicates if the Apple Pay setup pane is disabled.
- `auto_advance_setup_enabled` (Boolean) Indicates if Setup Assistant will automatically advance through its screens.
- `auto_unlock_with_watch_disabled` (Boolean) Indicates if the Unlock With Watch screen is disabled.
- `await_device_configured` (Boolean) Indicates if the device will need to wait for configured confirmation (the desktop is gated until MDM configuration finishes). Maps to `waitForDeviceConfiguredConfirmation`.
- `choose_your_lock_screen_disabled` (Boolean) Indicates if the Choose Your Lock Screen screen is disabled.
- `configuration_web_url` (Boolean) Indicates if the admin-assisted setup assistant login (web-based authentication) URL is used. Cannot be true when `use_platform_sso_during_setup_assistant` is true.
- `dep_onboarding_settings_id` (String) Identifier of the parent depOnboardingSetting (Apple ABM/ASM ADE token) that contains this macOS DEP enrollment profile. If omitted, the provider resolves it from the `/deviceManagement` endpoint's `intuneAccountId`. On tenants with multiple DEP tokens (for example, a separate Apple Configurator token), that fallback may select the wrong token, so set this explicitly to the ABM/ADE token id. List your tokens with `GET /deviceManagement/depOnboardingSettings` and pick the one whose `tokenType` is the Apple ABM/ADE token (not `appleConfigurator`).
- `description` (String) Description of the profile. Maximum length is 1500 characters.
- `device_name_template` (String) Sets a literal or name pattern for the device name.
- `diagnostics_disabled` (Boolean) Indicates if the Diagnostics setup pane is disabled.
- `display_tone_setup_disabled` (Boolean) Indicates if the DisplayTone setup screen is disabled.
- `dont_auto_populate_primary_account_info` (Boolean) Indicates whether Setup Assistant will auto-populate the primary account information.
- `enable_authentication_via_company_portal` (Boolean) Indicates to authenticate with the Company Portal instead of Apple Setup Assistant.
- `enable_restrict_editing` (Boolean) Indicates whether the user will be blocked from editing the account.
- `enrollment_time_azure_ad_group_ids` (Set of String) List of enrollment-time Microsoft Entra (Azure AD) group GUIDs to be associated with the profile.
- `file_vault_disabled` (Boolean) Indicates if FileVault is disabled.
- `hide_admin_account` (Boolean) Indicates whether the local admin account should be hidden.
- `icloud_diagnostics_disabled` (Boolean) Indicates if the iCloud Analytics screen is disabled.
- `icloud_storage_disabled` (Boolean) Indicates if the iCloud Documents and Desktop screen is disabled.
- `is_mandatory` (Boolean) Indicates if the profile is mandatory.
- `location_disabled` (Boolean) Indicates if the Location service setup pane is disabled.
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
- `supervised_mode_enabled` (Boolean) Supervised mode. True to enable, false otherwise.
- `support_department` (String) Support department information.
- `support_phone_number` (String) Support phone number.
- `terms_and_conditions_disabled` (Boolean) Indicates if the 'Terms and Conditions' setup pane is disabled.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `touch_id_disabled` (Boolean) Indicates if the Touch ID setup pane is disabled.
- `use_platform_sso_during_setup_assistant` (Boolean) Indicates whether Platform SSO is used as part of device enrollment during Setup Assistant. Cannot be true when `configuration_web_url` is true.
- `welcome_screen_disabled` (Boolean) Indicates if the Get Started (Welcome) setup pane is disabled. macOS 15 and later.
- `zoom_disabled` (Boolean) Indicates if the Zoom setup pane is disabled.

### Read-Only

- `configuration_endpoint_url` (String) Configuration endpoint url to use for enrollment. Generated by Intune.
- `enabled_skip_keys` (Set of String) Computed, read-only set of Setup Assistant skip keys (Apple `SkipKeys`) that the provider sends to Graph. This is derived from the individual `*_disabled` boolean attributes; do not set it directly. Note: `Privacy` and `Registration` are intentionally omitted from this array because the Microsoft Graph API rejects those skip-key strings, even though the `privacy_pane_disabled` and `registration_disabled` boolean properties work correctly.
- `id` (String) The unique identifier of the enrollment profile. Format is `{depOnboardingSettingsId}_{profileId}`.
- `is_default` (Boolean) Indicates if this is the default profile.

<a id="nestedatt--admin_account_password_rotation"></a>
### Nested Schema for `admin_account_password_rotation`

Optional:

- `auto_rotation_period_in_days` (Number) The number of days between automatic admin account password rotations.
- `on_retrieval_auto_rotate_password_enabled` (Boolean) Indicates whether the password is automatically rotated after retrieval.
- `on_retrieval_delay_auto_rotate_password_in_hours` (Number) The delay in hours before automatically rotating the password after retrieval.


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
