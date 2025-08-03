resource "microsoft365_graph_beta_device_management_device_category" "test" {
  display_name       = "Test Acceptance Device Category - Updated"
  description        = "Updated description for acceptance testing"
  role_scope_tag_ids = ["0", "1", "2"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}