resource "random_string" "application_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name = "acc-test-app-empty-${random_string.application_suffix.result}"
  description  = "Empty collection acceptance test"
  hard_delete  = true

  api = {
    known_client_applications   = []
    oauth2_permission_scopes    = []
    pre_authorized_applications = []
  }
  app_roles                = []
  required_resource_access = []
}
