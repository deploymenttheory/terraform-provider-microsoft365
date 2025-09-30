# Example usage of the iTunes app metadata data source
# This data source allows you to query the iTunes App Store API for app metadata

# Search for Firefox in the UK App Store
data "microsoft365_utility_itunes_app_metadata" "firefox_uk" {
  search_term  = "firefox"
  country_code = "gb"
}

# Search for Microsoft Office in the US App Store
data "microsoft365_utility_itunes_app_metadata" "office_us" {
  search_term  = "microsoft office"
  country_code = "us"
}

# Output the bundle ID of the first Firefox app result
output "firefox_bundle_id" {
  value = data.microsoft365_utility_itunes_app_metadata.firefox_uk.results[0].bundle_id
}

# Output the app name and icon URL of the first Microsoft Office app result
output "office_app_name" {
  value = data.microsoft365_utility_itunes_app_metadata.office_us.results[0].track_name
}

output "office_app_icon" {
  value = data.microsoft365_utility_itunes_app_metadata.office_us.results[0].artwork_url_512
}

# Output additional fields for Microsoft Office app
output "office_app_version" {
  value = data.microsoft365_utility_itunes_app_metadata.office_us.results[0].version
}

output "office_app_minimum_os_version" {
  value = data.microsoft365_utility_itunes_app_metadata.office_us.results[0].minimum_os_version
}

output "office_app_vpp_device_based_licensed" {
  value = data.microsoft365_utility_itunes_app_metadata.office_us.results[0].is_vpp_device_based_licensed
}

output "office_app_release_notes" {
  value = data.microsoft365_utility_itunes_app_metadata.office_us.results[0].release_notes
} 