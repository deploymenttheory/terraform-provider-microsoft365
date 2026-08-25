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

data "microsoft365_graph_beta_device_and_app_management_mobile_app" "company_portal_005" {
  display_name = "Microsoft Intune Company Portal"
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
# Pause - the group must have propagated before Intune will accept an assignment
# that references it, otherwise the create fails with a 400 naming an unknown
# group.
# ==============================================================================

resource "time_sleep" "wait_for_dependencies_005" {
  depends_on = [microsoft365_graph_beta_groups_group.acc_test_group_005_1]

  create_duration = "25s"
}

# ==============================================================================
# Scenario 5 Step 1: required intent with all settings
#
# intent carries RequiresReplace, so step 2 destroys and recreates the
# assignment. Terraform destroys before creating, so the replacement does not
# collide with the inclusion intent the original holds on the same group.
# ==============================================================================

resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "test_005" {
  mobile_app_id = data.microsoft365_graph_beta_device_and_app_management_mobile_app.company_portal_005.items[0].id
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
