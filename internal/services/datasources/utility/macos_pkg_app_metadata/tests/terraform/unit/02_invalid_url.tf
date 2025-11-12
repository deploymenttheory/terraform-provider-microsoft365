# Unit Test: Invalid URL format (should fail validation)
data "microsoft365_utility_macos_pkg_app_metadata" "test" {
  installer_url_source = "not-a-valid-url"
}

