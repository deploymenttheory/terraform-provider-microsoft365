---
page_title: "microsoft365_utility_itunes_app_metadata Data Source - terraform-provider-microsoft365"
subcategory: "Utilities"
description: |-
  Use this data source to query the iTunes App Store API for app metadata.
---

# microsoft365_utility_itunes_app_metadata

Use this data source to query the iTunes App Store API for app metadata. 
This data source allows you to search for apps by name and iTunes store country code, 
returning details like bundle ID and artwork URLs.

This data source can be used in conjunction with the the intunes vpp app resource
to populate metadata as part of the deployment of vpp apps from the iTunes App Store
with the intunes vpp app resource, 'microsoft365_graph_beta_device_and_app_management_macos_vpp_app'.

## Example Usage

```terraform
```terraform
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
```
```

## Argument Reference

* `search_term` - (Required) The search term to use when querying the iTunes App Store API.
* `country_code` - (Required) The two-letter country code for the App Store to search (e.g., 'us', 'gb', 'jp').

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The ID of this resource.
* `results` - A list of app results returned from the iTunes App Store API. Each result contains:
  * `track_id` - The unique identifier for the app.
  * `track_name` - The name of the app.
  * `bundle_id` - The bundle identifier of the app.
  * `artwork_url_60` - URL for the 60x60 app icon.
  * `artwork_url_100` - URL for the 100x100 app icon.
  * `artwork_url_512` - URL for the 512x512 app icon.
  * `seller_name` - The name of the app's seller/developer.
  * `primary_genre` - The primary genre of the app.
  * `description` - The description of the app.
  * `version` - The version of the app.
  * `price` - The price of the app in the local currency.
  * `formatted_price` - The formatted price of the app (e.g., 'Free', '$0.99').
  * `release_date` - The release date of the app.
  * `average_rating` - The average user rating of the app.
  * `artist_name` - The name of the artist/developer.
  * `minimum_os_version` - The minimum OS version required to run the app.
  * `content_advisory_rating` - The content advisory rating (e.g., '4+', '12+', '17+').
  * `is_vpp_device_based_licensed` - Whether the app supports VPP device-based licensing.
  * `release_notes` - Notes about the latest release of the app.
  * `currency` - The currency code for the price (e.g., 'USD', 'GBP', 'EUR').
  * `user_rating_count` - The number of user ratings for the app. 