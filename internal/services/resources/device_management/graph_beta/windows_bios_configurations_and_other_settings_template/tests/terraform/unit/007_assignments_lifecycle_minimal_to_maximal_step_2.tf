resource "microsoft365_graph_beta_device_management_windows_bios_configurations_and_other_settings_template" "test_007" {
  display_name               = "unit-test-windows-bios-configuration-007-assignments-lifecycle"
  file_name                  = "test-bios-007.cctk"
  configuration_file_content = "W2NjdGtdCjtEZWxsIENvbW1hbmQgfCBDb25maWd1cmUgLSBtaW5pbWFsIHVuaXQgdGVzdCBwYWNrYWdlClNlY3VyZUJvb3Q9RW5hYmxlZApUcG09T24K"

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "11111111-1111-1111-1111-111111111111"
      filter_id   = "44444444-4444-4444-4444-444444444444"
      filter_type = "include"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "66666666-6666-6666-6666-666666666666"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "77777777-7777-7777-7777-777777777777"
    }
  ]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
