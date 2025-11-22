# Example 3: Security Group with Dynamic Device Membership
# Creates a security group that automatically includes devices based on a membership rule.
# Ideal for device management scenarios like Conditional Access or Intune policies.
resource "microsoft365_graph_beta_groups_group" "security_dynamic_devices" {
  display_name                     = "Corporate Managed Devices"
  mail_nickname                    = "corporate-devices"
  mail_enabled                     = false
  security_enabled                 = true
  description                      = "Security group containing all corporate managed devices"
  group_types                      = ["DynamicMembership"]
  membership_rule                  = "(device.accountEnabled -eq true)"
  membership_rule_processing_state = "On"
}

