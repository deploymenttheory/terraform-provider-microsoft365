resource "microsoft365_graph_beta_device_and_app_management_browser_site_list" "example" {
  display_name = "Example Browser Site List"
  description  = "This is an example browser site list for Internet Explorer Mode"

  # Optional: Define custom timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}