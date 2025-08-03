resource "microsoft365_graph_beta_device_management_device_category" "maximal" {
  display_name        = "Test Maximal Device Category - Unique"
  description         = "Maximal device category for testing with all features"
  role_scope_tag_ids  = ["0", "1"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}