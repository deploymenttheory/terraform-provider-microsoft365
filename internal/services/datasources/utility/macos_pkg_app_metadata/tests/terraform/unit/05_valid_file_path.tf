# Unit Test: Valid file path format (validates regex, will fail when trying to read non-existent file)
data "microsoft365_utility_macos_pkg_app_metadata" "test" {
  installer_file_path_source = "/path/to/nonexistent/app.pkg"
}

