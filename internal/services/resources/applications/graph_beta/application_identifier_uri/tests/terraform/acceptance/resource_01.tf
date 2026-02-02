# Acceptance test: Application Identifier URI

resource "random_string" "test_id" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name = "acc-test-app-uri-${random_string.test_id.result}"
  description  = "Application for identifier URI acceptance test"
  hard_delete  = true
}

resource "time_sleep" "wait_for_app" {
  depends_on      = [microsoft365_graph_beta_applications_application.test]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_applications_application_identifier_uri" "test" {
  application_id = microsoft365_graph_beta_applications_application.test.id
  identifier_uri = "api://${microsoft365_graph_beta_applications_application.test.app_id}"

  depends_on = [time_sleep.wait_for_app]
}
