resource "microsoft365_graph_beta_device_management_windows_bios_configurations_and_other_settings_template" "test_006" {
  display_name                  = "unit-test-windows-bios-configuration-006-assignments-maximal"
  description                   = "Maximal test with multiple assignments"
  file_name                     = "test-bios-006.cctk"
  configuration_file_content    = "W2NjdGtdCjtEZWxsIENvbW1hbmQgfCBDb25maWd1cmUgLSBtYXhpbWFsIHVuaXQgdGVzdCBwYWNrYWdlCk51bUxvY2s9RW5hYmxlZApTZWN1cmVCb290PUVuYWJsZWQKVHBtPU9uClRwbUFjdGl2YXRpb249QWN0aXZhdGUKVXNiQm9vdD1EaXNhYmxlZApWaXJ0dWFsaXphdGlvbj1FbmFibGVkCldha2VPbkxhbj1MYW5Pbmx5Cg=="
  hardware_configuration_format = "dell"
  per_device_password_disabled  = true
  role_scope_tag_ids            = ["0", "1"]

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "11111111-1111-1111-1111-111111111111"
      filter_id   = "44444444-4444-4444-4444-444444444444"
      filter_type = "include"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "22222222-2222-2222-2222-222222222222"
      filter_id   = "55555555-5555-5555-5555-555555555555"
      filter_type = "exclude"
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
