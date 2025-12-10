# Minimal example - assign a manager to a user using existing user IDs
resource "microsoft365_graph_beta_users_user_manager" "example" {
  user_id    = "00000000-0000-0000-0000-000000000001" # The employee
  manager_id = "00000000-0000-0000-0000-000000000002" # The manager
}

