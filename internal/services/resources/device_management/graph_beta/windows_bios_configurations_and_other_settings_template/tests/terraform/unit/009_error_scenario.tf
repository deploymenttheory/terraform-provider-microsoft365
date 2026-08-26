resource "microsoft365_graph_beta_device_management_windows_bios_configurations_and_other_settings_template" "test_009" {
  display_name               = "unit-test-windows-bios-configuration-009-error"
  file_name                  = "test-bios-009.cctk"
  configuration_file_content = "W2NjdGtdCjtEZWxsIENvbW1hbmQgfCBDb25maWd1cmUgLSBtaW5pbWFsIHVuaXQgdGVzdCBwYWNrYWdlClNlY3VyZUJvb3Q9RW5hYmxlZApUcG09T24K"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
