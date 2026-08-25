resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Dependencies - the iOS store app the assignment hangs off, and its target group
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_ios_store_app" "acc_test_app_005" {
  display_name  = "acc-test-ios-store-app-005-${random_string.test_suffix.result}"
  description   = "Test iOS store app for mobile app assignment intent lifecycle"
  publisher     = "Terraform Provider Test"
  app_store_url = "https://apps.apple.com/us/app/microsoft-edge/id1288723196"

  applicable_device_type = {
    ipad            = true
    iphone_and_ipod = true
  }

  minimum_supported_operating_system = {
    v14_0 = true
  }

  timeouts = {
    create = "3m"
    read   = "1m"
    update = "3m"
    delete = "1m"
  }
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_005_1" {
  display_name     = "acc-test-group-005-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-005-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for mobile app assignment intent lifecycle"
  hard_delete      = true
}

# ==============================================================================
# Pause - the app and group must have propagated before Intune will accept an
# assignment that references them, otherwise the create fails with a 400 naming
# an unknown group or app.
# ==============================================================================

resource "time_sleep" "wait_for_dependencies_005" {
  depends_on = [
    microsoft365_graph_beta_device_and_app_management_ios_store_app.acc_test_app_005,
    microsoft365_graph_beta_groups_group.acc_test_group_005_1,
  ]

  create_duration = "30s"
}

# ==============================================================================
# Mobile App Assignment - Scenario 5 Step 1: required intent with all settings
#
# Both intent and settings carry RequiresReplace, so step 2 destroys and
# recreates the assignment rather than updating it in place. Terraform destroys
# before creating, so the replacement does not collide with the inclusion intent
# the original assignment holds on the same group.
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_005" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_ios_store_app.acc_test_app_005.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = microsoft365_graph_beta_groups_group.acc_test_group_005_1.id
  }

  settings = {
    ios_store = {
      is_removable                = true
      prevent_managed_app_backup  = true
      uninstall_on_device_removal = true
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }

  depends_on = [time_sleep.wait_for_dependencies_005]
}
