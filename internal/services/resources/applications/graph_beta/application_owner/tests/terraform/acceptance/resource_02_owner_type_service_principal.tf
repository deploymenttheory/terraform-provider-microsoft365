# Acceptance test: Application Owner with Service Principal owner type
# Full dependency chain: random_string -> applications (app + owner app) -> service_principal -> application_owner

resource "random_string" "test_id_sp" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test_app_sp" {
  display_name = "acc-test-app-owner-sp-${random_string.test_id_sp.result}"
  description  = "Application for service principal owner assignment acceptance test"
  hard_delete  = true
}

# Create a separate application that will become a service principal (owner)
resource "microsoft365_graph_beta_applications_application" "test_owner_app" {
  display_name = "acc-test-owner-app-${random_string.test_id_sp.result}"
  description  = "Application for creating service principal owner"
  hard_delete  = true
}

# Create service principal from the owner application
resource "microsoft365_graph_beta_applications_service_principal" "test_owner_sp" {
  app_id      = microsoft365_graph_beta_applications_application.test_owner_app.app_id
  hard_delete = true
}

resource "time_sleep" "wait_for_resources" {
  depends_on = [
    microsoft365_graph_beta_applications_application.test_app_sp,
    microsoft365_graph_beta_applications_service_principal.test_owner_sp
  ]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_applications_application_owner" "test_service_principal" {
  application_id    = microsoft365_graph_beta_applications_application.test_app_sp.id
  owner_id          = microsoft365_graph_beta_applications_service_principal.test_owner_sp.id
  owner_object_type = "ServicePrincipal"

  depends_on = [time_sleep.wait_for_resources]
}
