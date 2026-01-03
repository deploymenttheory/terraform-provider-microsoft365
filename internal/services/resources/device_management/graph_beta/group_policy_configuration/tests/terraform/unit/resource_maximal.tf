resource "microsoft365_graph_beta_device_management_group_policy_configuration" "maximal" {
  display_name        = "Maximal Group Policy Configuration"
  description         = "This is a comprehensive test configuration"
  role_scope_tag_ids  = ["0", "1", "2"]
}

