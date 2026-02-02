# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "test_id_multi" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Dependencies - Application with wait
# ==============================================================================


resource "microsoft365_graph_beta_applications_application" "test_multi" {
  display_name = "acc-test-app-uri-multi-${random_string.test_id_multi.result}"
  description  = "Application for multiple identifier URI acceptance test"
  hard_delete  = true
}

resource "time_sleep" "wait_for_app_multi" {
  depends_on      = [microsoft365_graph_beta_applications_application.test_multi]
  create_duration = "15s"
}

# ==============================================================================
# Test Case: Multiple Application Identifier URIs
# ==============================================================================

resource "microsoft365_graph_beta_applications_application_identifier_uri" "test_uri1" {
  application_id = microsoft365_graph_beta_applications_application.test_multi.id
  identifier_uri = "api://${microsoft365_graph_beta_applications_application.test_multi.app_id}"

  depends_on = [time_sleep.wait_for_app_multi]
}

resource "microsoft365_graph_beta_applications_application_identifier_uri" "test_uri2" {
  application_id = microsoft365_graph_beta_applications_application.test_multi.id
  identifier_uri = "https://deploymenttheory.com/${random_string.test_id_multi.result}" # must be a valid domain on the test tenant

  depends_on = [microsoft365_graph_beta_applications_application_identifier_uri.test_uri1]
}
