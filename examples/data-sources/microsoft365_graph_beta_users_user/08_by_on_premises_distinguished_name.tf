# Example 8: Look up a user by on_premises_distinguished_name (DN)

data "microsoft365_graph_beta_users_user" "by_dn" {
  on_premises_distinguished_name = "CN=John Doe,OU=Users,DC=contoso,DC=com"
}

output "user_id" {
  value = data.microsoft365_graph_beta_users_user.by_dn.items[0].id
}

output "display_name" {
  value = data.microsoft365_graph_beta_users_user.by_dn.items[0].display_name
}

output "on_premises_domain_name" {
  value = data.microsoft365_graph_beta_users_user.by_dn.items[0].on_premises_domain_name
}
