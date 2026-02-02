resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-application"
  description  = "Application with user owner"
}

resource "microsoft365_graph_beta_users_user" "app_owner" {
  display_name        = "Application Owner"
  user_principal_name = "app.owner@mycompany.com"
  mail_nickname       = "app.owner"
  account_enabled     = true
  password_profile = {
    password                           = "TempP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Assign user as application owner
resource "microsoft365_graph_beta_applications_application_owner" "user_owner" {
  application_id    = microsoft365_graph_beta_applications_application.example.id
  owner_id          = microsoft365_graph_beta_users_user.app_owner.id
  owner_object_type = "User"
}
