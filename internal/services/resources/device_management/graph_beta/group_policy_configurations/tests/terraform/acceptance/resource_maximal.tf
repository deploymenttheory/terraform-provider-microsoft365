resource "random_uuid" "test" {}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "test" {
  display_name       = "Test Acceptance Group Policy Configuration - Updated - ${random_uuid.test.result}"
  description        = "Updated description for acceptance testing"
  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
