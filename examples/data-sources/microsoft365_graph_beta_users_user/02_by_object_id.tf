# Example 2: Look up a user by object_id
# This example shows a representative set of the available output attributes.

data "microsoft365_graph_beta_users_user" "by_object_id" {
  object_id = "12345678-1234-1234-1234-123456789012"
}

output "user_id" {
  description = "The unique identifier for the user object"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].id
}

output "display_name" {
  description = "The name displayed in the address book for the user"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].display_name
}

output "user_principal_name" {
  description = "The user principal name (UPN) of the user"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].user_principal_name
}

output "mail" {
  description = "The SMTP address for the user"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].mail
}

output "job_title" {
  description = "The user's job title"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].job_title
}

output "department" {
  description = "The name of the department in which the user works"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].department
}

output "account_enabled" {
  description = "Whether the account is enabled"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].account_enabled
}

output "on_premises_sync_enabled" {
  description = "Whether the user is synced from an on-premises Active Directory"
  value       = data.microsoft365_graph_beta_users_user.by_object_id.items[0].on_premises_sync_enabled
}
