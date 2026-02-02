resource "microsoft365_graph_beta_applications_application" "example" {
  display_name = "my-application"
  description  = "Application with service principal"
}

resource "microsoft365_graph_beta_applications_service_principal" "example" {
  app_id = microsoft365_graph_beta_applications_application.example.app_id
}

resource "microsoft365_graph_beta_users_user" "sp_owner" {
  display_name        = "Service Principal Owner"
  user_principal_name = "sp.owner@mycompany.com"
  mail_nickname       = "sp.owner"
  account_enabled     = true
  password_profile = {
    password                           = "TempP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Assign user as service principal owner
resource "microsoft365_graph_beta_applications_service_principal_owner" "user_owner" {
  service_principal_id = microsoft365_graph_beta_applications_service_principal.example.id
  owner_id             = microsoft365_graph_beta_users_user.sp_owner.id
  owner_object_type    = "User"
}
