resource "microsoft365_graph_beta_device_management_windows_bios_configurations_and_other_settings_template" "test_011" {
  display_name                  = "unit-test-windows-bios-configuration-011-invalid-format"
  file_name                     = "test-bios-011.cctk"
  configuration_file_content    = "W2NjdGtdCjtEZWxsIENvbW1hbmQgfCBDb25maWd1cmUgLSBtaW5pbWFsIHVuaXQgdGVzdCBwYWNrYWdlClNlY3VyZUJvb3Q9RW5hYmxlZApUcG09T24K"
  hardware_configuration_format = "lenovo"
}
