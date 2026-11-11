
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_quality_update_policy" "test_003" {
  display_name     = "acc-test-windows-quality-update-policy-003-lifecycle-${random_string.test_suffix.result}"
  description      = "Lifecycle Step 2: Updated to maximal configuration"
  hotpatch_enabled = true

  role_scope_tag_ids = ["0", microsoft365_graph_beta_device_management_role_scope_tag.test.id]

}


resource "microsoft365_graph_beta_device_management_role_scope_tag" "test" {
  display_name = "acc-test-quality-scope-${random_string.test_suffix.result}"
  description  = "Scope tag owned by the Windows quality update acceptance test"
}
