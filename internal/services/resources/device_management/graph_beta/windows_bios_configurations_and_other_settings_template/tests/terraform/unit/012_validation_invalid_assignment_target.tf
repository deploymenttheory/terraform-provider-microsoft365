resource "microsoft365_graph_beta_device_management_windows_bios_configurations_and_other_settings_template" "test_012" {
  display_name               = "unit-test-windows-bios-configuration-012-invalid-assignment"
  file_name                  = "test-bios-012.cctk"
  configuration_file_content = "W2NjdGtdCjtEZWxsIENvbW1hbmQgfCBDb25maWd1cmUgLSBtaW5pbWFsIHVuaXQgdGVzdCBwYWNrYWdlClNlY3VyZUJvb3Q9RW5hYmxlZApUcG09T24K"

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]
}
