# Example 4: Role-Assignable Security Group
# Creates a security group that can be assigned to Entra ID roles.
# Note: Requires elevated permissions and visibility must be "Private".
# Once created, is_assignable_to_role cannot be changed.
resource "microsoft365_graph_beta_groups_group" "security_role_assignable" {
  display_name          = "Privileged Access Administrators"
  mail_nickname         = "privileged-admins"
  mail_enabled          = false
  security_enabled      = true
  description           = "Security group for privileged access administration"
  is_assignable_to_role = true
  visibility            = "Private"
}

