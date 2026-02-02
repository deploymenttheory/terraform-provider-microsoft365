resource "microsoft365_graph_beta_applications_application" "managed_app" {
  display_name = "my-managed-application"
  description  = "Application managed by a service principal"
}

resource "microsoft365_graph_beta_applications_application" "manager_app" {
  display_name = "app-manager"
  description  = "Service principal that manages other applications"
}

resource "microsoft365_graph_beta_applications_service_principal" "manager_sp" {
  app_id = microsoft365_graph_beta_applications_application.manager_app.app_id
}

# Assign service principal as application owner
resource "microsoft365_graph_beta_applications_application_owner" "sp_owner" {
  application_id    = microsoft365_graph_beta_applications_application.managed_app.id
  owner_id          = microsoft365_graph_beta_applications_service_principal.manager_sp.id
  owner_object_type = "ServicePrincipal"
}
