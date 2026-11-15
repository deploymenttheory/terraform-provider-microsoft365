# Returns the notification, enablement and authentication context rules to their defaults.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "11111111-1111-1111-1111-111111111111"

  enablement_rules = {
    Enablement_EndUser_Assignment = {
      enabled_rules = ["Justification"]
    }
  }

  notification_rules = {
    Notification_Approver_EndUser_Assignment = {
      notification_level            = "All"
      is_default_recipients_enabled = true
    }
  }

  authentication_context_rules = {
    AuthenticationContext_EndUser_Assignment = {
      is_enabled = false
    }
  }
}
