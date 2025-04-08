resource "microsoft365_graph_beta_device_and_app_management_device_category" "example" {
  display_name  = "Corporate Tablets"
  description   = "This category represents company-owned tablets"
  
  role_scope_tag_ids = ["8", "9"]

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}