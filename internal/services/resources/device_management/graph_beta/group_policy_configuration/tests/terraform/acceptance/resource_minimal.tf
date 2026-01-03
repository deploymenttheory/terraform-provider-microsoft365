resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "minimal" {
  display_name = "AccTest-Minimal-GPC-${random_string.suffix.result}"
}

