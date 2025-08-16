resource "microsoft365_graph_beta_device_and_app_management_ios_mobile_app_configuration" "role_scope_tags" {
  display_name       = "Test iOS Mobile App Configuration with Role Scope Tags"
  role_scope_tag_ids = ["0", "1", "2"]
}