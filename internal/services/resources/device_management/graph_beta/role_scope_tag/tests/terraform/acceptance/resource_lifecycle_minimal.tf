resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "lifecycle_first" {
  display_name     = "acc-test-rst-cycle-first-${random_string.suffix.result}"
  mail_nickname    = "rst-cycle-first-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "lifecycle_second" {
  display_name     = "acc-test-rst-cycle-second-${random_string.suffix.result}"
  mail_nickname    = "rst-cycle-second-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "lifecycle" {
  display_name = "acc-test-role-scope-tag-cycle-${random_string.suffix.result}"
  description  = "Role scope tag assignment lifecycle"
  assignments = [
    { type = "groupAssignmentTarget", group_id = microsoft365_graph_beta_groups_group.lifecycle_first.id }
  ]
}
