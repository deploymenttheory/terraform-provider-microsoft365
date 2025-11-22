# Example 2: Security Group with Dynamic User Membership
# Creates a security group that automatically adds/removes users based on a membership rule.
# Useful for automatically managing group membership based on user attributes.
resource "microsoft365_graph_beta_groups_group" "security_dynamic_users" {
  display_name                     = "Active Employees"
  mail_nickname                    = "active-employees"
  mail_enabled                     = false
  security_enabled                 = true
  description                      = "Security group containing all active employees"
  group_types                      = ["DynamicMembership"]
  membership_rule                  = "(user.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
}

