# Acceptance test: Service Principal - Minimal
# Full dependency chain: random_string -> application -> service_principal

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name = "acc-test-sp-minimal-${random_string.test_id.result}"
  description  = "Application for service principal minimal acceptance test"
  hard_delete  = true
}

resource "time_sleep" "wait_for_app" {
  depends_on      = [microsoft365_graph_beta_applications_application.test]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_applications_service_principal" "test_minimal" {
  app_id = microsoft365_graph_beta_applications_application.test.app_id

  depends_on = [time_sleep.wait_for_app]
}
