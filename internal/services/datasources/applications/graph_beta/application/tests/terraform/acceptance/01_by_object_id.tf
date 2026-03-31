resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name     = "acc-test-app-${random_string.suffix.result}"
  sign_in_audience = "AzureADMyOrg"

  tags = [
    "terraform-test",
    "datasource-test"
  ]
}

resource "time_sleep" "wait_for_app_propagation" {
  depends_on      = [microsoft365_graph_beta_applications_application.test]
  create_duration = "30s"
}

data "microsoft365_graph_beta_applications_application" "by_object_id" {
  object_id = microsoft365_graph_beta_applications_application.test.id

  depends_on = [time_sleep.wait_for_app_propagation]
}
