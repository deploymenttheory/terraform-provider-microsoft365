resource "microsoft365_graph_beta_device_management_device_category" "minimal" {
  display_name = "Test Minimal Device Category - Unique"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}