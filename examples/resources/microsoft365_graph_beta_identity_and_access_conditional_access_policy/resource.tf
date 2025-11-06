# ==============================================================================
# Conditional Access Policies - Baseline Device-based policies
# ==============================================================================
# This file defines 19 conditional access policies for Entra ID. Policies are
# organized by category:
# - CAD: Device-based policies (19 policies)
#
# IMPORTANT: All policies are set to "enabledForReportingButNotEnforced" by
# default. Review and test thoroughly before changing state to "enabled".
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# CAD006: Block Downloads on Unmanaged Devices
# Session control to block downloads on unmanaged devices for Office 365.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad006_block_download_unmanaged" {
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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
        "9b895d92-2cd3-44c7-9d02-a6ac2d5ea5c3", # Global Administrator
        "cf1c38e5-3621-4004-a7cb-879624dced7c", # Application Administrator
        "c4e39bd9-1100-46d3-8c65-fb160da0071f", # Authentication Administrator
        "25a516ed-2fa0-40ea-a2d0-12923a21473a", # Billing Administrator
        "aaf43236-0c0d-4d5f-883a-6955382ac081", # Cloud Application Administrator
        "b0f54661-2d74-4c50-afa3-1ec803f12efe", # Security Operator
        "158c047a-c907-4556-b7ef-446551a6b5f7", # Conditional Access Administrator
        "7698a772-787b-4ac8-901f-60d6b08affd2", # Cloud Device Administrator
        "17315797-102d-40b4-93e0-432062caca18", # User Administrator
        "b1be1c3e-b65d-4f19-8427-f6fa0d97feb9", # Compliance Administrator
        "9360feb5-f418-4baa-8175-e2a00bac4301", # Directory Writers
        "29232cdf-9323-42fd-ade2-1d097af3e4de", # Security Reader
        "f2ef992c-3afb-46b9-b7cf-a126ee74c451", # Global Reader
        "62e90394-69f5-4237-9190-012177145e10", # Exchange Administrator
        "729827e3-9c14-49f7-bb1b-9608f156bbb8", # Helpdesk Administrator
        "8ac3fc64-6eca-42ea-9e69-59f4c7b60eb2", # Hybrid Identity Administrator
        "3a2c62db-5318-420d-8d74-23affee5d9d5", # Insights Business Leader
        "966707d0-3269-4727-9be2-8c3a10f19b9d", # Intune Administrator
        "7be44c8a-adaf-4e2a-84d6-ab2649e08a13", # Knowledge Administrator
        "e8611ab8-c189-46e8-94e1-60213ab1f814", # Privileged Authentication Administrator
        "194ae4cb-b126-40b2-bd5b-6091b380977d", # Privileged Role Administrator
        "5f2222b1-57c3-48ba-8ad5-d4759f1fde6f", # Reports Reader
        "5d6b6bb7-de71-4623-b4af-96380a352509", # Search Administrator
        "f28a1f50-f6e7-4571-818b-6a12f2af6b6c", # SharePoint Administrator
        "69091246-20e8-4a56-aa4d-066075b2a7a8", # Teams Administrator
        "fe930be7-5e62-47db-91af-98c3a49a38b1"  # Security Administrator
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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
        "11111111-1111-1111-1111-111111111111", # Specific application 1
        "22222222-2222-2222-2222-222222222222", # Specific application 2
        "33333333-3333-3333-3333-333333333333", # Specific application 3
        "44444444-4444-4444-4444-444444444444"  # Specific application 4
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
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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
        rule = "device.isCompliant -eq True"
      }
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["compliantApplication"]
    custom_authentication_factors = []
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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
    }

    applications = {
      # Exchange Online, SharePoint Online, Windows Cloud PC
      include_applications = [
        "00000002-0000-0ff1-ce00-000000000000", # Exchange Online
        "00000003-0000-0ff1-ce00-000000000000", # SharePoint Online  
        "0af06dc6-e4b5-4f28-818e-e78e62d137a5"  # Windows Cloud PC
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
    built_in_controls             = []
    custom_authentication_factors = []
    authentication_strength = {
      id = "00000000-0000-0000-0000-000000000004" # Token protection strength ID
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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
      include_applications                            = ["0af06dc6-e4b5-4f28-818e-e78e62d137a5"] # Windows Cloud PC
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
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
      include_applications                            = ["d4ebce55-015a-49b5-a083-c84d1797ae8c"] # Microsoft Intune Enrollment
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
      id = "00000000-0000-0000-0000-000000000002" # Built-in MFA strength
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

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
