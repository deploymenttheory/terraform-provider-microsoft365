# Example 1: Basic Security Group with Assigned Membership
# Creates a standard security group where members are manually assigned.
# This is the most common type of security group used for access control.
resource "microsoft365_graph_beta_groups_group" "security_basic" {
  display_name     = "Engineering Team"
  mail_nickname    = "engineering-team"
  mail_enabled     = false
  security_enabled = true
  description      = "Security group for engineering team members"
  hard_delete      = true
}

