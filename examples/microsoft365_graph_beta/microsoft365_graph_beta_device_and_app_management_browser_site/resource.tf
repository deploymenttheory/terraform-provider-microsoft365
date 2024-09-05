resource "microsoft365_graph_beta_device_and_app_management_browser_site" "example" {
  allow_redirect     = true
  comment            = "Example browser site for Contoso"
  compatibility_mode = "default"
  merge_type         = "default"
  target_environment = "internetExplorerMode"
  web_url            = "https://www.contoso.com"

  timeouts {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}