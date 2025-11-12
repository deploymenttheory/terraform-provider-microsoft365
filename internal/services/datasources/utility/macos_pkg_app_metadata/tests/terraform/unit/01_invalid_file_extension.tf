# Unit Test: Invalid file extension (should fail validation)
data "microsoft365_utility_macos_pkg_app_metadata" "test" {
  installer_file_path_source = "/path/to/app.dmg"
}

