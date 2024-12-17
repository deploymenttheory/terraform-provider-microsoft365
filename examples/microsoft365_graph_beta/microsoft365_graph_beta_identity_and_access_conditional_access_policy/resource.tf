resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "example_policy" {
  display_name = "test"
  state        = "disabled"

  conditions = {
    applications = {
      include_applications = ["All"]
      exclude_applications = []
      include_user_actions = []
      application_filter   = null
    }

    users = {
      include_users  = ["All"]
      exclude_users  = ["11111111-1111-1111-1111-111111111111"]
      include_groups = []
      exclude_groups = ["11111111-1111-1111-1111-111111111111"]
      exclude_roles = [
        "11111111-1111-1111-1111-111111111111",
        "11111111-1111-1111-1111-111111111111"
      ]
      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["b2bCollaborationGuest", "b2bCollaborationMember"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    platforms = {
      include_platforms = ["iOS", "windows", "windowsPhone"]
      exclude_platforms = []
    }

    locations = {
      include_locations = [
        "11111111-1111-1111-1111-111111111111",
        "11111111-1111-1111-1111-111111111111"
      ]
      exclude_locations = []
    }

    client_app_types = ["browser", "mobileAppsAndDesktopClients", "exchangeActiveSync", "other"]

    devices = {
      device_filter = {
        mode = "include"
        rule = "device.deviceId -eq \"thing\""
      }
      include_devices = []
      exclude_devices = []
    }

    user_risk_levels    = ["high"]
    sign_in_risk_levels = ["none"]

    authentication_flows = {
      transfer_methods = ["deviceCodeFlow", "authenticationTransfer"]
    }
  }

  grant_controls = {
    operator          = "AND"
    built_in_controls = ["mfa", "approvedApplication"]
  }

  session_controls = {
    cloud_app_security = {
      is_enabled              = true
      cloud_app_security_type = "monitorOnly"
    }

    sign_in_frequency = {
      is_enabled          = true
      type                = "hours"
      value               = 5
      frequency_interval  = "timeBased"
      authentication_type = "primaryAndSecondaryAuthentication"
    }

    persistent_browser = {
      is_enabled = true
      mode       = "always"
    }

    continuous_access_evaluation = {
      mode = "strictLocation"
    }

    disable_resilience_defaults = true
  }

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}