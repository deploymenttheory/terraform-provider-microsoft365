resource "microsoft365_graph_beta_applications_application" "managed_sp_app" {
  display_name = "my-managed-service-principal"
  description  = "Service principal managed by another service principal"
}

resource "microsoft365_graph_beta_applications_service_principal" "managed_sp" {
  app_id = microsoft365_graph_beta_applications_application.managed_sp_app.app_id
}

resource "microsoft365_graph_beta_applications_application" "manager_app" {
  display_name = "sp-manager"
  description  = "Service principal that manages other service principals"
}

resource "microsoft365_graph_beta_applications_service_principal" "manager_sp" {
  app_id = microsoft365_graph_beta_applications_application.manager_app.app_id
}

# Assign service principal as another service principal's owner
resource "microsoft365_graph_beta_applications_service_principal_owner" "sp_owner" {
  service_principal_id = microsoft365_graph_beta_applications_service_principal.managed_sp.id
  owner_id             = microsoft365_graph_beta_applications_service_principal.manager_sp.id
  owner_object_type    = "ServicePrincipal"
}
