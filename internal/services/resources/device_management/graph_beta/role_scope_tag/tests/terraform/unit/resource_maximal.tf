resource "microsoft365_graph_beta_device_management_role_scope_tag" "maximal" {
  display_name = "Test Maximal Role Scope Tag - Unique"
  description  = "Maximal role scope tag for testing with all features"

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