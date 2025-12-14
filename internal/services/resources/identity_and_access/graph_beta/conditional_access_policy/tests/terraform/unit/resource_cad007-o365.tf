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
        "22222222-2222-2222-2222-222222222222",
        "33333333-3333-3333-3333-333333333333"
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

    sign_in_risk_levels           = []
    user_risk_levels              = []
    service_principal_risk_levels = []
    agent_id_risk_levels          = []
    insider_risk_levels           = []
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

