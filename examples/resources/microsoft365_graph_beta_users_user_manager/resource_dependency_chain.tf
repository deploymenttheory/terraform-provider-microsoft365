# Dependency chain example - create users and assign manager relationship

# Create the manager user
resource "microsoft365_graph_beta_users_user" "manager" {
  display_name        = "Jane Smith"
  account_enabled     = true
  user_principal_name = "jane.smith@contoso.com"
  mail_nickname       = "janesmith"
  hard_delete         = true
  password_profile = {
    password                           = "SecurePassword123!"
    force_change_password_next_sign_in = true
  }
}

# Create the employee user
resource "microsoft365_graph_beta_users_user" "employee" {
  display_name        = "John Doe"
  account_enabled     = true
  user_principal_name = "john.doe@contoso.com"
  mail_nickname       = "johndoe"
  hard_delete         = true
  password_profile = {
    password                           = "SecurePassword123!"
    force_change_password_next_sign_in = true
  }
}

# Assign the manager to the employee
resource "microsoft365_graph_beta_users_user_manager" "employee_manager" {
  user_id    = microsoft365_graph_beta_users_user.employee.id
  manager_id = microsoft365_graph_beta_users_user.manager.id
}

