# Role Scope Tags used for acceptance testing
# These role scope tags serve as dependencies.

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "acc_test_role_scope_tag_1" {
  display_name = "acc-test-role-scope-tag-1-${random_string.suffix.result}"
  description  = "Test role scope tag for acceptance testing"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "acc_test_role_scope_tag_2" {
  # This resource depends on the first one to avoid state locking issues.
  depends_on = [
    microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1
  ]
  display_name = "acc-test-role-scope-tag-2-${random_string.suffix.result}"
  description  = "Test role scope tag for acceptance testing"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}