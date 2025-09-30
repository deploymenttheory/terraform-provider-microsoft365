data "microsoft365_utility_itunes_app_metadata" "edge" {
  search_term  = "microsoft edge"
  country_code = "us"
}

resource "microsoft365_graph_beta_device_and_app_management_ios_store_app" "example" {
  display_name  = data.microsoft365_utility_itunes_app_metadata.edge.results[0].track_name
  description   = data.microsoft365_utility_itunes_app_metadata.edge.results[0].description
  publisher     = data.microsoft365_utility_itunes_app_metadata.edge.results[0].seller_name
  app_store_url = data.microsoft365_utility_itunes_app_metadata.edge.results[0].artist_view_url

  applicable_device_type = {
    ipad            = true
    iphone_and_ipod = true
  }

  minimum_supported_operating_system = {
    v14_0 = true
  }

  # Optional fields
  information_url         = "https://example.com/app-info"
  privacy_information_url = "https://example.com/privacy"
  owner                   = "IT Department"
  developer               = data.microsoft365_utility_itunes_app_metadata.edge.results[0].artist_name
  notes                   = "Managed by Terraform - Version: ${data.microsoft365_utility_itunes_app_metadata.edge.results[0].version}"
  is_featured             = false
  role_scope_tag_ids      = ["0"]

  # App icon (optional)
  app_icon = {
    icon_url_source = data.microsoft365_utility_itunes_app_metadata.edge.results[0].artwork_url_512
  }

  # Categories (optional)
  categories = ["Productivity", "Business"]

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
} 