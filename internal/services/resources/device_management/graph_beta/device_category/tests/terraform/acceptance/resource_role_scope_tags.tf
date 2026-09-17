resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_device_category" "role_tags" {
  display_name       = "Test Role Scope Tags Device Category - ${random_string.test_suffix.result}"
  role_scope_tag_ids = ["0", "1", "2"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}
