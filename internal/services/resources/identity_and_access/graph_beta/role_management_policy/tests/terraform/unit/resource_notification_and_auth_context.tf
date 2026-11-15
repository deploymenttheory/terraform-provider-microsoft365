# Manages a notification rule, an enablement rule and the authentication context rule.
resource "microsoft365_graph_beta_identity_and_access_role_management_policy" "test" {
  role_definition_id = "11111111-1111-1111-1111-111111111111"

  enablement_rules = {
    Enablement_EndUser_Assignment = {
      enabled_rules = ["Justification", "Ticketing"]
    }
  }

  notification_rules = {
    Notification_Approver_EndUser_Assignment = {
      notification_level            = "Critical"
      is_default_recipients_enabled = false
      notification_recipients       = ["pim-approvers@example.com"]
    }
  }

  authentication_context_rules = {
    AuthenticationContext_EndUser_Assignment = {
      is_enabled  = true
      claim_value = "c4"
    }
  }
}
