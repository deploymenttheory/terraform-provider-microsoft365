resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "Test Acceptance Role Scope Tag"
  description  = ""

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}