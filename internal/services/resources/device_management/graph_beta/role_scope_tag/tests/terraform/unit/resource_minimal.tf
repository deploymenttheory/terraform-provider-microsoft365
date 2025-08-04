resource "microsoft365_graph_beta_device_management_role_scope_tag" "minimal" {
  display_name = "Test Minimal Role Scope Tag - Unique"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}