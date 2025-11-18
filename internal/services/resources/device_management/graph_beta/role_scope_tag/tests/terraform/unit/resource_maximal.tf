resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = true
  lower   = true
  numeric = true
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "maximal" {
  display_name = "unit-test-role-scope-tag-maximal-${random_string.suffix.result}"
  description  = "unit-test-role-scope-tag-maximal-${random_string.suffix.result}"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    }
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}