resource "microsoft365_graph_beta_device_and_app_management_browser_site" "example_site" {
  browser_site_list_assignment_id = microsoft365_graph_beta_device_and_app_management_browser_site_list.example.id
  web_url                         = "https://example.com"
  allow_redirect                  = true
  compatibility_mode              = "internetExplorer11"
  comment                         = "Example site for IE mode"
  target_environment              = "internetExplorerMode"
  merge_type                      = "noMerge"

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}