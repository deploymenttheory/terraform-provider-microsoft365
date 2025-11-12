# Unit Test: Valid URL format (validates regex, will fail when trying to download from fake URL)
data "microsoft365_utility_macos_pkg_app_metadata" "test" {
  installer_url_source = "https://nonexistent.example.com/app.pkg"
}

