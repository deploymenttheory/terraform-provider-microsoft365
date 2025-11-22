# Example 5: Microsoft 365 Group with Dynamic User Membership
# Creates a Microsoft 365 group (formerly Office 365 group) with automatic membership.
# Includes Teams, SharePoint, Outlook, and other Microsoft 365 services.
resource "microsoft365_graph_beta_groups_group" "m365_dynamic_users" {
  display_name                     = "Marketing Department"
  mail_nickname                    = "marketing-dept"
  mail_enabled                     = true
  security_enabled                 = true
  group_types                      = ["Unified", "DynamicMembership"]
  membership_rule                  = "(user.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
  visibility                       = "Private"
}

