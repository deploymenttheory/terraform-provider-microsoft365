resource "random_id" "test_008" {
  byte_length = 4
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test_008" {
  display_name = "acc-test-gpd-errors-${random_id.test_008.hex}"
  description  = "Acceptance test for group policy definition validation errors"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

