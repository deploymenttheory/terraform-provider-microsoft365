data "microsoft365_utility_windows_msi_app_metadata" "firefox" {
  installer_url_source = "https://download.mozilla.org/?product=firefox-msi-latest-ssl&os=win64&lang=en-US"
} 