resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "block_legacy_authentication" {
  display_name = "Block Legacy Authentication"
  state        = "enabled"

  conditions = {
    applications = {
      include_applications = ["All"]
      exclude_applications = []
      include_user_actions = []
      application_filter   = null
    }

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = ["11111111-1111-1111-1111-111111111111"] # Emergency access group
      exclude_roles = [
        "62e90394-69f5-4237-9190-012177145e10" # Global Administrator
      ]
      exclude_guests_or_external_users = null
    }

    platforms = {
      include_platforms = []
      exclude_platforms = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = [
        "11111111-1111-1111-1111-111111111111" # Trusted office locations
      ]
    }

    client_app_types = ["exchangeActiveSync", "other"]

    devices = {
      device_filter   = null
      include_devices = []
      exclude_devices = []
    }

    user_risk_levels    = []
    sign_in_risk_levels = []

    authentication_flows = null
  }

  grant_controls = {
    operator          = "OR"
    built_in_controls = ["block"]
  }

  session_controls = null

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "require_mfa_for_admins" {
  display_name = "Require MFA for Admin Roles"
  state        = "enabled"
  conditions = {
    applications = {
      include_applications = ["All"]
      exclude_applications = []
      include_user_actions = []
      application_filter   = null
    }

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = ["11111111-1111-1111-1111-111111111111"] # Emergency access group
      include_roles = [
        "62e90394-69f5-4237-9190-012177145e10", # Global Administrator
        "194ae4cb-b126-40b2-bd5b-6091b380977d", # Security Administrator
        "729827e3-9c14-49f7-bb1b-9608f156bbb8"  # Helpdesk Administrator
      ]
      exclude_roles                    = []
      exclude_guests_or_external_users = null
    }

    platforms = {
      include_platforms = []
      exclude_platforms = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = []
    }

    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    devices = {
      device_filter   = null
      include_devices = []
      exclude_devices = []
    }

    user_risk_levels    = []
    sign_in_risk_levels = []

    authentication_flows = null
  }

  grant_controls = {
    operator          = "AND"
    built_in_controls = ["mfa"]
  }

  session_controls = {
    sign_in_frequency = {
      is_enabled          = true
      type                = "hours"
      value               = 4
      frequency_interval  = "timeBased"
      authentication_type = "primaryAndSecondaryAuthentication"
    }
  }

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "compliant_device_policy" {
  display_name = "Require Compliant or Hybrid Joined Device"
  state        = "enabled"
  conditions = {
    applications = {
      include_applications = ["Office365"]
      exclude_applications = []
      include_user_actions = []
      application_filter   = null
    }

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = ["11111111-1111-1111-1111-111111111111"] # Emergency access group
      exclude_roles = [
        "62e90394-69f5-4237-9190-012177145e10" # Global Administrator
      ]
      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["b2bCollaborationGuest", "b2bCollaborationMember"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    platforms = {
      include_platforms = ["windows", "macOS", "iOS", "android"]
      exclude_platforms = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = [
        "11111111-1111-1111-1111-111111111111" # Trusted office locations
      ]
    }

    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    devices = {
      device_filter = {
        mode = "exclude"
        rule = "device.isCompliant -eq True or device.trustType -eq \"Hybrid Azure AD joined\""
      }
      include_devices = []
      exclude_devices = []
    }

    user_risk_levels    = []
    sign_in_risk_levels = []

    authentication_flows = null
  }

  grant_controls = {
    operator          = "OR"
    built_in_controls = ["compliantDevice", "domainJoinedDevice"]
  }

  session_controls = {
    cloud_app_security = {
      is_enabled              = true
      cloud_app_security_type = "monitorOnly"
    }

    continuous_access_evaluation = {
      mode = "strictLocation"
    }
  }

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "high_risk_sign_in_policy" {
  display_name = "High Risk Sign-in Policy"
  state        = "enabled"
  conditions = {
    applications = {
      include_applications = ["All"]
      exclude_applications = []
      include_user_actions = []
      application_filter   = null
    }

    users = {
      include_users                    = ["All"]
      exclude_users                    = []
      include_groups                   = []
      exclude_groups                   = ["11111111-1111-1111-1111-111111111111"] # Emergency access group
      exclude_roles                    = []
      exclude_guests_or_external_users = null
    }

    platforms = {
      include_platforms = []
      exclude_platforms = []
    }

    locations = {
      include_locations = ["All"]
      exclude_locations = []
    }

    client_app_types = ["browser", "mobileAppsAndDesktopClients", "exchangeActiveSync", "other"]

    devices = {
      device_filter   = null
      include_devices = []
      exclude_devices = []
    }

    user_risk_levels    = []
    sign_in_risk_levels = ["high"]

    authentication_flows = null
  }

  grant_controls = {
    operator          = "AND"
    built_in_controls = ["mfa", "passwordChange"]
  }

  session_controls = {
    sign_in_frequency = {
      is_enabled          = true
      type                = "hours"
      value               = 1
      frequency_interval  = "timeBased"
      authentication_type = "primaryAndSecondaryAuthentication"
    }

    persistent_browser = {
      is_enabled = true
      mode       = "never"
    }
  }

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}