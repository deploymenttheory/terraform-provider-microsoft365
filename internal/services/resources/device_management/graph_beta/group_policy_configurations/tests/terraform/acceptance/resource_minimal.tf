resource "random_uuid" "test" {}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test" {
  display_name = "Test Acceptance Group Policy Configuration - ${random_uuid.test.result}"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
