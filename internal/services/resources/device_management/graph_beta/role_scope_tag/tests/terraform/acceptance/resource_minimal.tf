resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
  numeric = true
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "acc-test-role-scope-tag-minimal-${random_string.suffix.result}"
  description  = "acc-test-role-scope-tag-minimal-${random_string.suffix.result}"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}