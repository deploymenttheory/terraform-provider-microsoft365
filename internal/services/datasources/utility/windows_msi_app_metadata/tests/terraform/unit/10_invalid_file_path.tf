data "microsoft365_utility_windows_msi_app_metadata" "test" {
  installer_file_path_source = "testdata/nonexistent.msi"
}
