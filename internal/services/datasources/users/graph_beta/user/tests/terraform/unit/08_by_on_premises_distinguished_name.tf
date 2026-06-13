data "microsoft365_graph_beta_users_user" "test" {
  on_premises_distinguished_name = "CN=Test UserOne,OU=Users,DC=contoso,DC=com"
}
