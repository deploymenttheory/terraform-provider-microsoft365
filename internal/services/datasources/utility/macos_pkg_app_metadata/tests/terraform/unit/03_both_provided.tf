# Unit Test: Both installer_file_path_source and installer_url_source provided (should fail)
data "microsoft365_utility_macos_pkg_app_metadata" "test" {
  installer_file_path_source = "/path/to/app.pkg"
  installer_url_source       = "https://example.com/app.pkg"
}

