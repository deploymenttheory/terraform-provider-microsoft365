---
page_title: "microsoft365_graph_beta_identity_and_access_conditional_access_policy Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages Microsoft 365 Conditional Access Policies using the /identity/conditionalAccess/policies endpoint. Conditional Access policies define the conditions under which users can access cloud apps.
---

# microsoft365_graph_beta_identity_and_access_conditional_access_policy (Resource)

Manages Microsoft 365 Conditional Access Policies using the `/identity/conditionalAccess/policies` endpoint. Conditional Access policies define the conditions under which users can access cloud apps.

## Microsoft Documentation

- [conditionalAccessPolicy resource type](https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesspolicy?view=graph-rest-beta)
- [Create conditionalAccessPolicy](https://learn.microsoft.com/en-us/graph/api/conditionalaccessroot-post-policies?view=graph-rest-beta)
- [Update conditionalAccessPolicy](https://learn.microsoft.com/en-us/graph/api/conditionalaccesspolicy-update?view=graph-rest-beta)
- [Delete conditionalAccessPolicy](https://learn.microsoft.com/en-us/graph/api/conditionalaccesspolicy-delete?view=graph-rest-beta)
- [Conditional Access documentation](https://learn.microsoft.com/en-us/azure/active-directory/conditional-access/)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `Policy.ReadWrite.ConditionalAccess`, `Policy.Read.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.19.0-alpha | Experimental | Initial release |
| v0.34.0-alpha | Experimental | Numerous bug fixes and added graceful 429 error handling |

## Example Usage

```terraform
# ==============================================================================
# Conditional Access Policies - Baseline (48 Policies)
# ==============================================================================
# This file defines 48 conditional access policies for Entra ID. Policies are
# organized by category:
# - CAD: Device-based policies (19 policies)
# - CAL: Location-based policies (6 policies)
# - CAP: Platform-based policies (4 policies)
# - CAU: User-based policies (19 policies)
#
# IMPORTANT: All policies are set to "enabledForReportingButNotEnforced" by
# default (except CAU011 which is "disabled"). Review and test thoroughly
# before changing state to "enabled".
#
# Break Glass Account Exclusion: All policies exclude the break glass accounts
# group to ensure emergency access is maintained.
# ==============================================================================

# ==============================================================================
# Device-Based Policies (CAD) - 19 Policies
# ==============================================================================
# CAD001: macOS Device Compliance
# Grants macOS access to Office 365 for all users when using modern auth clients
# and device is compliant.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad001_macos_compliant" {
  display_name = "acc-CAD001-O365: Grant macOS access for All users when Modern Auth Clients and Compliant-v1.1"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad001_exclude.id
      ]
      include_roles = []
      exclude_roles = []
      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["macOS"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice"]
    custom_authentication_factors = []
  }


}

# CAD002: Windows Device Compliance
# Grants Windows access to Office 365 for all users when using modern auth clients
# and device is compliant.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad002_windows_compliant" {
  display_name = "CAD002-O365: Grant Windows access for All users when Modern Auth Clients and Compliant-v1.1"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad002_exclude.id
      ]
      include_roles = []
      exclude_roles = []
      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["windows"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "domainJoinedDevice"]
    custom_authentication_factors = []
  }


}

# CAD003: iOS and Android Device Compliance or App Protection
# Grants iOS and Android access to Office 365 for all users when using modern auth
# clients and device has app protection policy or is compliant.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad003_mobile_compliant_or_app_protection" {
  display_name = "CAD003-O365: Grant iOS and Android access for All users when Modern Auth Clients and AppProPol or Compliant-v1.3"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad003_exclude.id
      ]
      include_roles = []
      exclude_roles = []
      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["android", "iOS"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "compliantApplication"]
    custom_authentication_factors = []
  }


}

# CAD004: Require MFA on Non-Compliant Devices via Browser
# Requires MFA for all users accessing Office 365 via browser when device is non-compliant.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad004_browser_noncompliant_mfa" {
  display_name = "CAD004-O365: Grant Require MFA for All users when Browser and Non-Compliant-v1.3"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad004_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    devices = {
      device_filter = {
        mode = "exclude"
        rule = "device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""
      }
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
    # Note: Source uses custom authentication strength ID "eaedd457-3e01-413b-a02e-417489193d1d" (Custom MFA)
    # Replace with your tenant-specific custom authentication strength ID or use built-in:
    authentication_strength = {
      id = "00000000-0000-0000-0000-000000000002" # Multifactor authentication (Changed from custom MFA ID)
    }
  }


}

# CAD005: Block Unsupported Device Platforms
# Blocks access to Office 365 for unsupported device platforms.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad005_block_unsupported_platforms" {
  display_name = "CAD005-O365: Block access for unsupported device platforms for All users when Modern Auth Clients-v1.1"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad005_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["all"]
      exclude_platforms = ["android", "iOS", "windows", "macOS", "linux"]
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }


}

# CAD006: Block Downloads on Unmanaged Devices
# Session control to block downloads on unmanaged devices for Office 365.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad006_session_block_download_unmanaged" {
  display_name = "CAD006-O365: Session block download on unmanaged device for All users when Browser and Modern App Clients and Non-Compliant-v1.5"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad006_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    devices = {
      device_filter = {
        mode = "exclude"
        rule = "device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""
      }
    }

    sign_in_risk_levels = []
  }

  session_controls = {
    application_enforced_restrictions = {
      is_enabled = true
    }
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
  }


}

# CAD007: Sign-in Frequency for Mobile Apps on Non-Compliant Devices
# Sets sign-in frequency to 7 days for Office 365 mobile apps on non-compliant iOS/Android devices.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad007_mobile_signin_frequency" {
  display_name = "CAD007-O365: Session set Sign-in Frequency for Apps for All users when Modern Auth Clients and Non-Compliant-v1.2"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad007_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["android", "iOS"]
      exclude_platforms = []
    }

    devices = {
      device_filter = {
        mode = "exclude"
        rule = "device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""
      }
    }

    sign_in_risk_levels = []
  }

  session_controls = {
    sign_in_frequency = {
      value               = 7
      type                = "days"
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "timeBased"
      is_enabled          = true
    }
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
  }


}

# CAD008: Sign-in Frequency for Browser on Non-Compliant Devices
# Sets sign-in frequency to 1 hour for all apps accessed via browser on non-compliant devices.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad008_browser_signin_frequency" {
  display_name = "CAD008-All: Session set Sign-in Frequency for All users when Browser and Non-Compliant-v1.1"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad008_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    devices = {
      device_filter = {
        mode = "exclude"
        rule = "device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""
      }
    }

    sign_in_risk_levels = []
  }

  session_controls = {
    sign_in_frequency = {
      value               = 1
      type                = "days"
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "timeBased"
      is_enabled          = true
    }
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
  }


}

# CAD009: Disable Browser Persistence on Non-Compliant Devices
# Disables persistent browser session for all apps on non-compliant devices.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad009_disable_browser_persistence" {
  display_name = "CAD009-All: Session disable browser persistence for All users when Browser and Non-Compliant-v1.2"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad009_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    devices = {
      device_filter = {
        mode = "exclude"
        rule = "device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""
      }
    }

    sign_in_risk_levels = []
  }

  session_controls = {
    persistent_browser = {
      mode       = "never"
      is_enabled = true
    }
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
  }


}

# CAD010: Require MFA for Device Registration/Join
# Requires MFA when users register or join devices to Entra ID.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad010_device_registration_mfa" {
  display_name = "CAD010-RJD: Require MFA for device join or registration when Browser and Modern Auth Clients-v1.1"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["all"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad010_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = []
      exclude_applications                            = []
      include_user_actions                            = ["urn:user:registerdevice"]
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }


}

# CAD011: Linux Device Compliance
# Grants Linux access to Office 365 for all users (excluding guests) when using
# modern auth clients and device is compliant.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad011_linux_compliant" {
  display_name = "CAD011-O365: Grant Linux access for All users when Modern Auth Clients and Compliant-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = ["GuestsOrExternalUsers"]
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad001_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["linux"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice"]
    custom_authentication_factors = []
  }


}

# CAD012: Admin Access on Compliant Devices
# Grants access for admin users to all apps when using compliant devices.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad012_admin_compliant_access" {
  display_name = "CAD012-All: Grant access for Admin users when Browser and Modern Auth Clients and Compliant-v1.1"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad012_exclude.id
      ]
      include_roles = [
        data.microsoft365_graph_beta_identity_and_access_role_definitions.application_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.application_developer.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.authentication_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.authentication_extensibility_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.b2c_ief_keyset_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.billing_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.cloud_application_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.cloud_device_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.user_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.compliance_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.directory_writers.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.security_reader.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.global_reader.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.exchange_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.helpdesk_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.hybrid_identity_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.insights_business_leader.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.intune_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.knowledge_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.privileged_authentication_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.privileged_role_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.reports_reader.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.search_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.sharepoint_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.teams_administrator.items[0].id,
        data.microsoft365_graph_beta_identity_and_access_role_definitions.security_administrator.items[0].id,
      ]
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "domainJoinedDevice"]
    custom_authentication_factors = []
  }


}

# CAD013: Selected Apps - Compliant Device Requirement
# Requires compliant device for access to selected applications.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad013_selected_apps_compliant" {
  display_name = "CAD013-Selected: Grant access for All users when Browser and Modern Auth Clients and Compliant-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad013_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications = [
        "a4f2693f-129c-4b96-982b-2c364b8314d7", # Specific application 1
        "499b84ac-1321-427f-aa17-267ca6975798", # Specific application 2
        "996def3d-b36c-4153-8607-a6fd3c01b89f", # Specific application 3
        "797f4846-ba00-4fd7-ba43-dac1f8f63013"  # Specific application 4
      ]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["all"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "domainJoinedDevice"]
    custom_authentication_factors = []
  }

  timeouts = {
    create = "150s"
    read   = "150s"
    update = "150s"
    delete = "150s"
  }
}

# CAD014: Edge App Protection on Windows
# Requires app protection policy for Edge browser on Windows for Office 365 access.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad014_edge_app_protection_windows" {
  display_name = "CAD014-O365: Require App Protection Policy for Edge on Windows for All users when Browser and Non-Compliant-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cad014_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad014_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["windows"]
      exclude_platforms = []
    }

    devices = {
      device_filter = {
        mode = "exclude"
        rule = "device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""
      }
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantApplication"]
    custom_authentication_factors = []
  }


}

# CAD015: Compliant Device for Windows and macOS Browser Access
# Requires compliant device for all users accessing all apps via browser on Windows/macOS.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad015_windows_macos_browser_compliant" {
  display_name = "CAD015-All: Grant access for All users when Browser and Modern Auth Clients and Compliant on Windows and macOS-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cad015_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad015_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["windows", "macOS"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "domainJoinedDevice"]
    custom_authentication_factors = []
  }


}

# CAD016: Token Protection for EXO/SPO/CloudPC on Windows
# Requires token protection for Exchange Online, SharePoint Online, and Cloud PC on Windows.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad016_token_protection_windows" {
  display_name = "CAD016-EXO_SPO_CloudPC: Require token protection when Modern Auth Clients on Windows-v1.2"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cad016_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad016_exclude.id
      ]
      include_roles = []
      exclude_roles = []

      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications = [
        data.microsoft365_graph_beta_applications_service_principal.azure_virtual_desktop.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.windows_365.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.windows_cloud_login.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.office_365_exchange_online.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.office_365_sharepoint_online.items[0].app_id,
      ]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["windows"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }


}

# CAD017: Selected Apps - Mobile App Protection or Compliance
# Requires app protection policy or device compliance for selected apps on iOS/Android.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad017_selected_mobile_app_protection" {
  display_name = "CAD017-Selected: Grant iOS and Android access for All users when Modern Auth Clients and AppProPol or Compliant-v1.1"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = [microsoft365_graph_beta_groups_group.cad017_include.id]
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad017_exclude.id
      ]
      include_roles = []
      exclude_roles = []

      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications                            = ["None"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["android", "iOS"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "compliantApplication"]
    custom_authentication_factors = []
  }


}

# CAD018: Cloud PC - Mobile App Protection or Compliance
# Requires app protection policy or device compliance for Cloud PC access on iOS/Android.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad018_cloudpc_mobile_app_protection" {
  display_name = "CAD018-CloudPC: Grant iOS and Android access for All users when Modern Auth Clients and AppProPol or Compliant-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad018_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications = [
        data.microsoft365_graph_beta_applications_service_principal.azure_virtual_desktop.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.microsoft_remote_desktop.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.windows_365.items[0].app_id,
        data.microsoft365_graph_beta_applications_service_principal.windows_cloud_login.items[0].app_id
      ]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["android", "iOS"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "compliantApplication"]
    custom_authentication_factors = []
  }


}

# CAD019: Intune Enrollment - MFA and Sign-in Frequency
# Requires MFA and sets sign-in frequency to every time for Intune enrollment.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad019_intune_enrollment_mfa" {
  display_name = "CAD019-Intune: Require MFA and set sign-in frequency to every time-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad019_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications = [
        data.microsoft365_graph_beta_applications_service_principal.microsoft_intune_enrollment.items[0].app_id
      ]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
    authentication_strength = {
      id = "00000000-0000-0000-0000-000000000002" # multifactor_authentication
    }
  }

  session_controls = {
    sign_in_frequency = {
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "everyTime"
      is_enabled          = true
      # Note: type and value are not set when frequency_interval is "everyTime"
    }
  }


}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `conditions` (Attributes) Conditions that must be met for the policy to apply. (see [below for nested schema](#nestedatt--conditions))
- `display_name` (String) The display name for the Conditional Access policy.
- `grant_controls` (Attributes) Controls for granting access. (see [below for nested schema](#nestedatt--grant_controls))
- `state` (String) Specifies the state of the policy. Possible values are: enabled, disabled, enabledForReportingButNotEnforced.

### Optional

- `partial_enablement_strategy` (String) Strategy for partial enablement of the policy.
- `session_controls` (Attributes) Controls for managing user sessions. (see [below for nested schema](#nestedatt--session_controls))
- `template_id` (String) ID of the template this policy is derived from.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `created_date_time` (String) The creation date and time of the policy.
- `deleted_date_time` (String) The deletion date and time of the policy, if applicable.
- `id` (String) String (identifier)
- `modified_date_time` (String) The last modified date and time of the policy.

<a id="nestedatt--conditions"></a>
### Nested Schema for `conditions`

Required:

- `applications` (Attributes) Applications and user actions included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--applications))
- `client_app_types` (Set of String) Client application types included in the policy. Possible values are: all, browser, mobileAppsAndDesktopClients, exchangeActiveSync, other.
- `sign_in_risk_levels` (Set of String) Sign-in risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.
- `users` (Attributes) Users, groups, and roles included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--users))

Optional:

- `client_applications` (Attributes) Client applications configuration for the conditional access policy. (see [below for nested schema](#nestedatt--conditions--client_applications))
- `device_states` (Attributes) Device states included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--device_states))
- `devices` (Attributes) Devices included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--devices))
- `locations` (Attributes) Locations included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--locations))
- `platforms` (Attributes) Platforms included in and excluded from the policy. (see [below for nested schema](#nestedatt--conditions--platforms))
- `service_principal_risk_levels` (Set of String) Service principal risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.
- `times` (Attributes) Times and days when the policy applies. (see [below for nested schema](#nestedatt--conditions--times))
- `user_risk_levels` (Set of String) User risk levels included in the policy. Possible values are: low, medium, high, hidden, none, unknownFutureValue.

<a id="nestedatt--conditions--applications"></a>
### Nested Schema for `conditions.applications`

Required:

- `exclude_applications` (Set of String) Applications to exclude from the policy. For empty requests, use []
- `include_applications` (Set of String) Applications to include in the policy. Can use the special value 'All' to include all applications.
- `include_authentication_context_class_references` (Set of String) Authentication context secures data and actions in applications, including custom applications, line-of-business (LOB) applications, SharePoint, and applications protected by Microsoft Defender for Cloud Apps. Can be predefined builtin contexts: `require_trusted_device` (or c1), `require_terms_of_use` (or c2), `require_trusted_location` (or c3), `require_strong_authentication` (or c4), `required_trust_type:azure_ad_joined` (or c5), `require_access_from_an_approved_app` (or c6), `required_trust_type:hybrid_azure_ad_joined` (or c7) or custom authentication context class references in the format 'c' followed by a number from 8 through to 99 (e.g., c1, c8, c10, c25, c99). Learn more here 'https://learn.microsoft.com/en-us/entra/identity/conditional-access/concept-conditional-access-cloud-apps#authentication-context'.
- `include_user_actions` (Set of String) User actions to include in the policy.

Optional:

- `application_filter` (Attributes) Configure app filters you want to policy to apply to. Using custom security attributes you can use the rule builder or rule syntax text box to create or edit the filter rules. this feature is currently in preview, only attributes of type String are supported. Attributes of type Integer or Boolean are not currently supported. Learn more here 'https://learn.microsoft.com/en-us/entra/identity/conditional-access/concept-filter-for-applications'. (see [below for nested schema](#nestedatt--conditions--applications--application_filter))
- `global_secure_access` (Attributes) Global Secure Access settings for the conditional access policy. (see [below for nested schema](#nestedatt--conditions--applications--global_secure_access))

<a id="nestedatt--conditions--applications--application_filter"></a>
### Nested Schema for `conditions.applications.application_filter`

Required:

- `mode` (String) Mode of the filter. Possible values are: include, exclude.
- `rule` (String) Rule syntax for the filter.


<a id="nestedatt--conditions--applications--global_secure_access"></a>
### Nested Schema for `conditions.applications.global_secure_access`



<a id="nestedatt--conditions--users"></a>
### Nested Schema for `conditions.users`

Required:

- `exclude_groups` (Set of String) Groups to exclude from the policy.
- `exclude_roles` (Set of String) Microsoft Entra tenant roles to exclude from the policy.
- `exclude_users` (Set of String) Users to exclude from the policy. Can use special values like 'GuestsOrExternalUsers'.
- `include_groups` (Set of String) Groups to include in the policy.
- `include_roles` (Set of String) Roles to include in the policy.
- `include_users` (Set of String) Users to include in the policy. Can use special values like 'All', 'None', or 'GuestsOrExternalUsers'.

Optional:

- `exclude_guests_or_external_users` (Attributes) Configuration for excluding guests or external users. (see [below for nested schema](#nestedatt--conditions--users--exclude_guests_or_external_users))
- `include_guests_or_external_users` (Attributes) Configuration for including guests or external users. (see [below for nested schema](#nestedatt--conditions--users--include_guests_or_external_users))

<a id="nestedatt--conditions--users--exclude_guests_or_external_users"></a>
### Nested Schema for `conditions.users.exclude_guests_or_external_users`

Required:

- `external_tenants` (Attributes) Configuration for external tenants. (see [below for nested schema](#nestedatt--conditions--users--exclude_guests_or_external_users--external_tenants))

Optional:

- `guest_or_external_user_types` (Set of String) Types of guests or external users to exclude. Possible values are: InternalGuest, B2bCollaborationGuest, B2bCollaborationMember, B2bDirectConnectUser, OtherExternalUser, ServiceProvider.

<a id="nestedatt--conditions--users--exclude_guests_or_external_users--external_tenants"></a>
### Nested Schema for `conditions.users.exclude_guests_or_external_users.external_tenants`

Required:

- `membership_kind` (String) Kind of membership. Possible values are: all, enumerated, unknownFutureValue.

Optional:

- `members` (Set of String) The list of tenant IDs for external tenants.



<a id="nestedatt--conditions--users--include_guests_or_external_users"></a>
### Nested Schema for `conditions.users.include_guests_or_external_users`

Required:

- `external_tenants` (Attributes) Configuration for external tenants. (see [below for nested schema](#nestedatt--conditions--users--include_guests_or_external_users--external_tenants))

Optional:

- `guest_or_external_user_types` (Set of String) Types of guests or external users to include. Possible values are: InternalGuest, B2bCollaborationGuest, B2bCollaborationMember, B2bDirectConnectUser, OtherExternalUser, ServiceProvider.

<a id="nestedatt--conditions--users--include_guests_or_external_users--external_tenants"></a>
### Nested Schema for `conditions.users.include_guests_or_external_users.external_tenants`

Required:

- `membership_kind` (String) Kind of membership. Possible values are: all, enumerated, unknownFutureValue.

Optional:

- `members` (Set of String) The list of Microsoft Entra organization tenant IDs for external tenants to exclude from the CA policy.




<a id="nestedatt--conditions--client_applications"></a>
### Nested Schema for `conditions.client_applications`

Required:

- `include_service_principals` (Set of String) Service principals to include in the policy. Can use the special value 'ServicePrincipalsInMyTenant' to include all service principals.

Optional:

- `exclude_service_principals` (Set of String) Service principals to exclude from the policy.


<a id="nestedatt--conditions--device_states"></a>
### Nested Schema for `conditions.device_states`

Optional:

- `exclude_states` (Set of String) Device states to exclude from the policy.
- `include_states` (Set of String) Device states to include in the policy.


<a id="nestedatt--conditions--devices"></a>
### Nested Schema for `conditions.devices`

Optional:

- `device_filter` (Attributes) Filter that defines the devices the policy applies to. (see [below for nested schema](#nestedatt--conditions--devices--device_filter))
- `exclude_device_states` (Set of String) Device states to exclude from the policy.
- `exclude_devices` (Set of String) Devices to exclude from the policy.
- `include_device_states` (Set of String) Device states to include in the policy.
- `include_devices` (Set of String) Devices to include in the policy.

<a id="nestedatt--conditions--devices--device_filter"></a>
### Nested Schema for `conditions.devices.device_filter`

Required:

- `mode` (String) Mode of the filter. Possible values are: include, exclude.
- `rule` (String) Rule syntax for the filter.



<a id="nestedatt--conditions--locations"></a>
### Nested Schema for `conditions.locations`

Required:

- `exclude_locations` (Set of String) Named locations to exclude from the policy. Can use special values like 'AllTrusted' or provide guid'sof named locations.
- `include_locations` (Set of String) Named locations to include in the policy. Can use special values like 'All' or 'AllTrusted' 'or provide guid'sof named locations.


<a id="nestedatt--conditions--platforms"></a>
### Nested Schema for `conditions.platforms`

Required:

- `include_platforms` (Set of String) Platforms to include in the policy.

Optional:

- `exclude_platforms` (Set of String) Platforms to exclude from the policy.


<a id="nestedatt--conditions--times"></a>
### Nested Schema for `conditions.times`

Optional:

- `all_day` (Boolean) Whether the policy applies all day.
- `end_time` (String) End time for the policy.
- `excluded_ranges` (Set of String) Time ranges when the policy does not apply.
- `included_ranges` (Set of String) Time ranges when the policy applies.
- `start_time` (String) Start time for the policy.
- `time_zone` (String) Time zone for the policy times.



<a id="nestedatt--grant_controls"></a>
### Nested Schema for `grant_controls`

Required:

- `built_in_controls` (Set of String) List of built-in controls required by the policy. Possible values are: block, mfa, compliantDevice, domainJoinedDevice, approvedApplication, compliantApplication, passwordChange, unknownFutureValue.
- `custom_authentication_factors` (Set of String) Custom authentication factors for granting access.
- `operator` (String) Operator to apply to the controls. Possible values are: AND, OR. When setting a singular operator, use 'OR'.

Optional:

- `authentication_strength` (Attributes) Authentication strength is a Conditional Access control that specifies which combinations of authentication methods can be used to access a resource. Users can satisfy the strength requirements by authenticating with any of the allowed combinations. read more here 'https://learn.microsoft.com/en-us/entra/identity/authentication/concept-authentication-strengths'. (see [below for nested schema](#nestedatt--grant_controls--authentication_strength))
- `terms_of_use` (Set of String) Terms of use required for granting access.

<a id="nestedatt--grant_controls--authentication_strength"></a>
### Nested Schema for `grant_controls.authentication_strength`

Required:

- `id` (String) ID of the authentication strength policy. Can be a GUID or predefined built-in values: 'multifactor_authentication' (maps to '00000000-0000-0000-0000-000000000002'), 'passwordless_mfa' (maps to '00000000-0000-0000-0000-000000000003'), or 'phishing_resistant_mfa' (maps to '00000000-0000-0000-0000-000000000004').

Optional:

- `allowed_combinations` (Set of String) The allowed authentication method combinations that satisfy the authentication strength policy.
- `description` (String) Description of the authentication strength policy.
- `display_name` (String) Display name of the authentication strength policy.
- `policy_type` (String) Type of the policy. Possible values are: builtIn, custom.
- `requirements_satisfied` (String) Requirements satisfied by the policy.

Read-Only:

- `created_date_time` (String) Creation date and time of the authentication strength policy.
- `modified_date_time` (String) Last modified date and time of the authentication strength policy.



<a id="nestedatt--session_controls"></a>
### Nested Schema for `session_controls`

Optional:

- `application_enforced_restrictions` (Attributes) Application enforced restrictions for the session. (see [below for nested schema](#nestedatt--session_controls--application_enforced_restrictions))
- `cloud_app_security` (Attributes) Cloud app security controls for the session. (see [below for nested schema](#nestedatt--session_controls--cloud_app_security))
- `continuous_access_evaluation` (Attributes) Continuous access evaluation controls for the session. (see [below for nested schema](#nestedatt--session_controls--continuous_access_evaluation))
- `disable_resilience_defaults` (Boolean) Whether to disable resilience defaults.
- `global_secure_access_filtering_profile` (Attributes) Global Secure Access filtering profile for the session. (see [below for nested schema](#nestedatt--session_controls--global_secure_access_filtering_profile))
- `persistent_browser` (Attributes) Persistent browser controls for the session. (see [below for nested schema](#nestedatt--session_controls--persistent_browser))
- `secure_sign_in_session` (Attributes) Secure sign-in session controls. (see [below for nested schema](#nestedatt--session_controls--secure_sign_in_session))
- `sign_in_frequency` (Attributes) Sign-in frequency controls for the session. (see [below for nested schema](#nestedatt--session_controls--sign_in_frequency))

<a id="nestedatt--session_controls--application_enforced_restrictions"></a>
### Nested Schema for `session_controls.application_enforced_restrictions`

Required:

- `is_enabled` (Boolean) Whether application enforced restrictions are enabled.


<a id="nestedatt--session_controls--cloud_app_security"></a>
### Nested Schema for `session_controls.cloud_app_security`

Required:

- `cloud_app_security_type` (String) Type of cloud app security control. Possible values are: blockDownloads, mcasConfigured, monitorOnly, unknownFutureValue.
- `is_enabled` (Boolean) Whether cloud app security controls are enabled.


<a id="nestedatt--session_controls--continuous_access_evaluation"></a>
### Nested Schema for `session_controls.continuous_access_evaluation`

Required:

- `mode` (String) Mode for continuous access evaluation. Possible values are: disabled, basic, strict.


<a id="nestedatt--session_controls--global_secure_access_filtering_profile"></a>
### Nested Schema for `session_controls.global_secure_access_filtering_profile`

Required:

- `is_enabled` (Boolean) Whether global secure access filtering controls are enabled.
- `profile_id` (String) ID of the global secure access filtering profile.


<a id="nestedatt--session_controls--persistent_browser"></a>
### Nested Schema for `session_controls.persistent_browser`

Required:

- `is_enabled` (Boolean) Whether persistent browser controls are enabled.
- `mode` (String) Mode for persistent browser. Possible values are: always, never.


<a id="nestedatt--session_controls--secure_sign_in_session"></a>
### Nested Schema for `session_controls.secure_sign_in_session`

Required:

- `is_enabled` (Boolean) Whether secure sign-in session controls are enabled.


<a id="nestedatt--session_controls--sign_in_frequency"></a>
### Nested Schema for `session_controls.sign_in_frequency`

Required:

- `is_enabled` (Boolean) Whether sign-in frequency controls are enabled.

Optional:

- `authentication_type` (String) Authentication type for sign-in frequency. Possible values are: primaryAndSecondaryAuthentication, secondaryAuthentication.
- `frequency_interval` (String) Frequency interval for sign-in frequency. Possible values are: timeBased, everyTime.
- `type` (String) Type of sign-in frequency control. Possible values are: days, hours. Not used when frequency_interval is everyTime.
- `value` (Number) Value for the sign-in frequency. Not used when frequency_interval is everyTime.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Important Notes

### Policy States
- **enabled**: The policy is active and will be enforced
- **disabled**: The policy exists but is not enforced
- **enabledForReportingButNotEnforced**: The policy will be evaluated and logged but not enforced (report-only mode)

### Applications
- Use `"All"` to target all cloud applications
- Use `"Office365"` to target all Office 365 applications
- Use specific application IDs for targeted policies
- Application filters support complex OData expressions for fine-grained control

### Users and Groups
- Use `"All"` to target all users
- Use `"GuestsOrExternalUsers"` to target external users
- Specify user, group, or role object IDs for targeted policies
- Emergency access accounts should always be excluded from blocking policies

### Locations
- Named locations must be created in Azure AD before referencing
- Use `"All"` for all locations or `"AllTrusted"` for all trusted locations
- IP-based and country-based locations are supported

### Client App Types
- `browser`: Web browsers
- `mobileAppsAndDesktopClients`: Mobile apps and desktop clients
- `exchangeActiveSync`: Exchange ActiveSync clients
- `other`: Other clients including legacy authentication

### Grant Controls
- **Operator**: `AND` requires all controls, `OR` requires any control
- **Built-in Controls**: `block`, `mfa`, `compliantDevice`, `domainJoinedDevice`, `approvedApplication`, `compliantApplication`, `passwordChange`
- **Authentication Strength**: Reference to custom authentication strength policies

### Session Controls
- **Application Restrictions**: Control access to specific applications
- **Cloud App Security**: Integration with Microsoft Defender for Cloud Apps
- **Sign-in Frequency**: Control how often users must re-authenticate
- **Persistent Browser**: Control browser session persistence
- **Continuous Access Evaluation**: Real-time policy evaluation

### Device Filters
- Support complex OData expressions for device-based conditions
- Common filters include device compliance, trust type, and device attributes
- Use `include` mode to target devices matching the filter
- Use `exclude` mode to exclude devices matching the filter

### Risk-based Policies
- **User Risk Levels**: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`
- **Sign-in Risk Levels**: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`
- Requires Azure AD Identity Protection licenses

### Best Practices
- Always exclude emergency access accounts from blocking policies
- Test policies in report-only mode before enabling enforcement
- Use specific targeting rather than broad "All" assignments when possible
- Monitor policy impact through Azure AD sign-in logs
- Implement a phased rollout for new policies
- Document policy purpose and expected behavior

### Common Policy Scenarios
- **Block Legacy Authentication**: Target legacy client app types with block control
- **Require MFA for Admins**: Target administrative roles with MFA requirement
- **Device Compliance**: Require compliant or domain-joined devices for access
- **Location-based Access**: Block or require additional controls based on location
- **Risk-based Access**: Respond to user or sign-in risk with appropriate controls

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import using composite ID format: {policy_id}/{condition_id}
terraform import microsoft365_graph_beta_identity_and_access_conditional_access_policy.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
``` 