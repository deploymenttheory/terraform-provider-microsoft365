resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "maximal" {
  display_name       = "AccTest-Maximal-GPC-${random_string.suffix.result}"
  description        = "Acceptance test for maximal group policy configuration"
  role_scope_tag_ids = ["0"]
}

