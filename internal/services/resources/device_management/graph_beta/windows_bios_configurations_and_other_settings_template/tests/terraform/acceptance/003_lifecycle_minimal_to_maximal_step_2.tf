
resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_windows_bios_configurations_and_other_settings_template" "test_003" {
  display_name                  = "acc-test-windows-bios-configuration-003-lifecycle-${random_string.test_suffix.result}"
  description                   = "Promoted to a maximal BIOS configuration template"
  file_name                     = "test-bios-003.cctk"
  configuration_file_content    = "W2NjdGtdCjtEZWxsIENvbW1hbmQgfCBDb25maWd1cmUgLSBtYXhpbWFsIHVuaXQgdGVzdCBwYWNrYWdlCk51bUxvY2s9RW5hYmxlZApTZWN1cmVCb290PUVuYWJsZWQKVHBtPU9uClRwbUFjdGl2YXRpb249QWN0aXZhdGUKVXNiQm9vdD1EaXNhYmxlZApWaXJ0dWFsaXphdGlvbj1FbmFibGVkCldha2VPbkxhbj1MYW5Pbmx5Cg=="
  hardware_configuration_format = "dell"
  per_device_password_disabled  = true
  role_scope_tag_ids            = ["0", "1"]

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
