# Example 6: Microsoft 365 Group with Role Assignment
# Creates a Microsoft 365 group that can be assigned to Entra ID roles.
# Combines collaboration features with privileged access management.
# Note: Requires elevated permissions and visibility must be "Private".
resource "microsoft365_graph_beta_groups_group" "m365_role_assignable" {
  display_name          = "Executive Leadership Team"
  mail_nickname         = "executive-team"
  mail_enabled          = true
  security_enabled      = true
  group_types           = ["Unified"]
  description           = "Microsoft 365 group for executive leadership"
  is_assignable_to_role = true
  visibility            = "Private"
  hard_delete           = true
}

