resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Dependencies - the app the assignment hangs off, and its target group
#
# The app is the tenant's Microsoft Intune Company Portal, looked up by display
# name through the mobile_app data source. It is a permanent fixture of the tenant and
# is already published, which matters: Intune leaves a newly created iOS store app
# in publishingState "processing", publishingState is read-only on both POST and
# PATCH, and the service refuses to assign an app that is not published -
#
#   400 "Invalid operation: app's PublishingState is not 'Published'."
#
# so a scenario that creates its own app cannot reliably assign to it. Assignment
# behaviour is what these scenarios cover; app creation is not.
# ==============================================================================

data "microsoft365_graph_beta_device_and_app_management_mobile_app" "company_portal_001" {
  display_name = "Microsoft Intune Company Portal"
}

resource "microsoft365_graph_beta_groups_group" "acc_test_group_001_1" {
  display_name     = "acc-test-group-001-1-${random_string.test_suffix.result}"
  mail_nickname    = "acc-test-group-001-1-${random_string.test_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Test group 1 for mobile app assignment available intent"
  hard_delete      = true
}

# ==============================================================================
# Pause - the group must have propagated before Intune will accept an assignment
# that references it, otherwise the create fails with a 400 naming an unknown
# group.
# ==============================================================================

resource "time_sleep" "wait_for_dependencies_001" {
  depends_on = [microsoft365_graph_beta_groups_group.acc_test_group_001_1]

  create_duration = "25s"
}

# ==============================================================================
# Scenario 1: available intent, minimal settings
#
# is_removable is omitted. The Intune service rejects isRemovable for any intent
# other than required, so the provider must not send the field at all. This is
# the configuration reported in issue #3692 and it fails on main.
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_001" {
  mobile_app_id = data.microsoft365_graph_beta_device_and_app_management_mobile_app.company_portal_001.items[0].id
  intent        = "available"
  source        = "direct"

  target = {
    target_type = "groupAssignment"
    group_id    = microsoft365_graph_beta_groups_group.acc_test_group_001_1.id
  }

  settings = {
    ios_store = {}
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }

  depends_on = [time_sleep.wait_for_dependencies_001]
}
