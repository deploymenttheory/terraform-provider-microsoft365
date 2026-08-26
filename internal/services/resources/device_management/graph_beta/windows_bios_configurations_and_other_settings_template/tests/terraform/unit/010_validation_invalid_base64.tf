resource "microsoft365_graph_beta_device_management_windows_bios_configurations_and_other_settings_template" "test_010" {
  display_name               = "unit-test-windows-bios-configuration-010-invalid-base64"
  file_name                  = "test-bios-010.cctk"
  configuration_file_content = "[cctk] this is raw text, not base64"
}
