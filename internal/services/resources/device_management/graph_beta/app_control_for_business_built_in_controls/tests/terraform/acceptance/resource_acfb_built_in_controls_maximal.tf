resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "acc_test_group_004" {
  display_name     = "acc-test-group-004-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-004-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group for app control policy maximal configuration"
  hard_delete      = true
}

# ==============================================================================
# App Control Policy Resource - Maximal Configuration
# ==============================================================================

# Advanced App Control for Business configuration with multiple assignments
resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "advanced" {
  name        = "acc-test-app-control-for-business-built-in-controls-maximal-${random_string.test_suffix.result}"
  description = "acc-test-app-control-for-business-built-in-controls-maximal"

  # App Control settings
  enable_app_control                 = "audit"                                                                   # audit = logs but allows, enforce = blocks untrusted apps
  additional_rules_for_trusting_apps = ["trust_apps_with_good_reputation", "trust_apps_from_managed_installers"] # Both Microsoft Store and managed installer apps

  # Role scope tags
  role_scope_tag_ids = ["0"]

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group_004
  ]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_004.id
    },
    {
      type = "allDevicesAssignmentTarget"
    }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}