resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "transition" {
  display_name = "AccTest-Transition-GPC-${random_string.suffix.result}"
}

