resource "random_string" "user_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "random_string" "manager_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create the user who will have a manager assigned
resource "microsoft365_graph_beta_users_user" "employee" {
  display_name        = "acc-test-employee-${random_string.user_suffix.result}"
  user_principal_name = "acc-test-employee-${random_string.user_suffix.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-employee-${random_string.user_suffix.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

# Create the user who will be the manager
resource "microsoft365_graph_beta_users_user" "manager" {
  display_name        = "acc-test-manager-${random_string.manager_suffix.result}"
  user_principal_name = "acc-test-manager-${random_string.manager_suffix.result}@deploymenttheory.com"
  mail_nickname       = "acc-test-manager-${random_string.manager_suffix.result}"
  account_enabled     = true
  hard_delete         = true
  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = false
  }
}

# Assign the manager relationship
resource "microsoft365_graph_beta_users_user_manager" "test" {
  user_id    = microsoft365_graph_beta_users_user.employee.id
  manager_id = microsoft365_graph_beta_users_user.manager.id
}

