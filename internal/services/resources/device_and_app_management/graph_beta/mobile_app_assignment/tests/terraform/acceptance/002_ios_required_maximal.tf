resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Dependencies - the iOS store app the assignment hangs off, and its target group
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_ios_store_app" "acc_test_app_002" {
  display_name  = "acc-test-ios-store-app-002-${random_string.test_suffix.result}"
  description   = "Test iOS store app for mobile app assignment required intent"
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

resource "microsoft365_graph_beta_groups_group" "acc_test_group_002_1" {
  display_name     = "acc-test-group-002-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-002-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for mobile app assignment required intent"
  hard_delete      = true
}

# ==============================================================================
# Pause - the app and group must have propagated before Intune will accept an
# assignment that references them, otherwise the create fails with a 400 naming
# an unknown group or app.
# ==============================================================================

resource "time_sleep" "wait_for_dependencies_002" {
  depends_on = [
    microsoft365_graph_beta_device_and_app_management_ios_store_app.acc_test_app_002,
    microsoft365_graph_beta_groups_group.acc_test_group_002_1,
  ]

  create_duration = "25s"
}

# ==============================================================================
# Mobile App Assignment - Scenario 2: required intent, maximal settings
#
# required is the one intent that accepts all three settings the service
# constrains, so each is set explicitly here and must reach the API unchanged.
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_002" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_ios_store_app.acc_test_app_002.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = microsoft365_graph_beta_groups_group.acc_test_group_002_1.id
    device_and_app_management_assignment_filter_type = "none"
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

  depends_on = [time_sleep.wait_for_dependencies_002]
}
