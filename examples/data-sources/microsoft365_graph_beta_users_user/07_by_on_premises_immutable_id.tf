# Example 7: Look up a user by on_premises_immutable_id (sourceAnchor)
# Useful for correlating cloud users with on-premises Active Directory accounts.

data "microsoft365_graph_beta_users_user" "by_immutable_id" {
  on_premises_immutable_id = "T0AbQ29udG9zb1VzZXI="
}

output "user_id" {
  value = data.microsoft365_graph_beta_users_user.by_immutable_id.items[0].id
}

output "display_name" {
  value = data.microsoft365_graph_beta_users_user.by_immutable_id.items[0].display_name
}

output "on_premises_sam_account_name" {
  value = data.microsoft365_graph_beta_users_user.by_immutable_id.items[0].on_premises_sam_account_name
}
