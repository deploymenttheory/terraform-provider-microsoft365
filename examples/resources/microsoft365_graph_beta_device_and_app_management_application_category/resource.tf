resource "microsoft365_graph_beta_device_and_app_management_application_category" "example" {
  display_name = "Business Apps"

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}