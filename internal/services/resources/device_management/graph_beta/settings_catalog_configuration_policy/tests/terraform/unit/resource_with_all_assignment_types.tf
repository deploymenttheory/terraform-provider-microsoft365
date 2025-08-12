resource "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy" "all_assignment_types" {
  name         = "Test All Assignment Types Settings Catalog Policy - Unit"
  description  = "Settings catalog configuration policy with all assignment types for unit testing"
  platforms    = "macOS"
  technologies = ["mdm", "appleRemoteManagement"]

   template_reference = {
    template_id = ""
  }

  configuration_policy = {
    name     = "Test All Assignment Types Settings Catalog Policy Configuration - Unit"
    settings = []
  }

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "11111111-1111-1111-1111-111111111111"
      filter_type = "include"
      filter_id   = "2222222-2222-2222-2222-222222222222"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "33333333-3333-3333-3333-333333333333"
      filter_type = "include"
      filter_id   = "4444444-4444-4444-4444-444444444444"
    },
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "include"
      filter_id   = "5555555-5555-5555-5555-555555555555"
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "include"
      filter_id   = "6666666-6666-6666-6666-666666666666"
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = "7777777-7777-7777-7777-777777777777"
      filter_type = "include"
      filter_id   = "8888888-8888-8888-8888-888888888888"
    }
  ]

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}