resource "microsoft365_graph_beta_device_management_device_category" "test" {
  display_name = "Test Acceptance Device Category"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}