resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Dependencies - the app the assignment hangs off, and its target group
#
# The app is looked up rather than created. Intune leaves a newly created app in
# publishingState "processing", publishingState is read-only on both POST and
# PATCH, and the service refuses to assign an app that is not published -
#
#   400 "Invalid operation: app's PublishingState is not 'Published'."
#
# so a scenario that creates its own app cannot reliably assign to it. Assignment
# behaviour is what these scenarios cover; app creation is not.
# ==============================================================================

data "microsoft365_graph_beta_device_and_app_management_mobile_app" "app_011" {
  display_name = "Microsoft Desktop App Installer Preview"
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_011_1" {
  display_name     = "acc-test-group-011-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-011-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for windows_universal_app_x mobile app assignment"
  hard_delete      = true
}

# ==============================================================================
# Pause - the group must have propagated before Intune will accept an assignment
# that references it, otherwise the create fails with a 400 naming an unknown
# group.
# ==============================================================================

resource "time_sleep" "wait_for_dependencies_011" {
  depends_on = [microsoft365_graph_beta_groups_group.acc_test_group_011_1]

  create_duration = "25s"
}

# ==============================================================================
# Scenario 11: windows_universal_app_x settings against a Universal AppX app
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_011" {
  mobile_app_id = data.microsoft365_graph_beta_device_and_app_management_mobile_app.app_011.items[0].id
  intent        = "required"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = microsoft365_graph_beta_groups_group.acc_test_group_011_1.id
  }

  settings = {
    windows_universal_app_x = {
      use_device_context = true
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }

  depends_on = [time_sleep.wait_for_dependencies_011]
}
