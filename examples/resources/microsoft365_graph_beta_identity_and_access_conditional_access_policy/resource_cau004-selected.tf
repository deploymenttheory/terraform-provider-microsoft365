# CAU004: Route Through Microsoft Defender for Cloud Apps
# Routes browser traffic through Microsoft Defender for Cloud Apps (MDCA) for
# monitoring and control on non-compliant devices.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau004_mdca_route" {
  display_name = "CAU004-Selected: Session route through MDCA for All users when Browser on Non-Compliant-v1.2"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau004_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      # Note: Add specific application IDs, typically includes Office365
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
    cloud_app_security = {
      cloud_app_security_type = "mcasConfigured"
      is_enabled              = true
    }
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
  }
}

