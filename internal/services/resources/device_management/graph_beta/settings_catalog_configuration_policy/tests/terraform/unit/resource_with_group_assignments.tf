resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "group_assignments" {
  name         = "Test Group Assignments Settings Catalog Policy - Unit"
  description  = ""
  platforms    = "macOS"
  technologies = ["mdm", "appleRemoteManagement"]

  template_reference = {
    template_id = ""
  }

  configuration_policy = {
    name     = "Test Group Assignments Settings Catalog Policy Configuration - Unit"
    settings = []
  }

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "11111111-1111-1111-1111-111111111111"
      filter_type = "include"
      filter_id   = "22222222-2222-2222-2222-222222222222"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "33333333-3333-3333-3333-333333333333"
      filter_type = "include"
      filter_id   = "44444444-4444-4444-4444-444444444444"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}