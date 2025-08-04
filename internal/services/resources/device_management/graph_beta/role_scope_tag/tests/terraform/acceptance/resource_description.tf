resource "microsoft365_graph_beta_device_management_role_scope_tag" "description" {
  display_name = "Test Description Role Scope Tag"
  description  = "This is a test role scope tag with description"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}