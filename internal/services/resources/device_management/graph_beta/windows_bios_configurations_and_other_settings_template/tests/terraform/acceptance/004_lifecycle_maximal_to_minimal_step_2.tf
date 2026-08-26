
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_bios_configurations_and_other_settings_template" "test_004" {
  display_name                  = "acc-test-windows-bios-configuration-004-lifecycle-${random_string.test_suffix.result}"
  description                   = "Reduced to a minimal BIOS configuration template"
  file_name                     = "test-bios-004.cctk"
  configuration_file_content    = "W2NjdGtdCjtEZWxsIENvbW1hbmQgfCBDb25maWd1cmUgLSBtaW5pbWFsIHVuaXQgdGVzdCBwYWNrYWdlClNlY3VyZUJvb3Q9RW5hYmxlZApUcG09T24K"
  hardware_configuration_format = "dell"
  per_device_password_disabled  = false
  role_scope_tag_ids            = ["0"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
