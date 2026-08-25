resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Dependencies - an already published iOS store app, and the target group
#
# The app is looked up rather than created. Intune leaves a newly created iOS
# store app in publishingState "processing" - still processing more than an hour
# after creation when measured against a live tenant - publishingState is
# read-only on both POST and PATCH, and the service rejects an assignment for any
# app that is not published:
#
#   400 "Invalid operation: app's PublishingState is not 'Published'."
#
# so a scenario that creates its own app can never assign to it, whatever wait it
# uses. Assignment behaviour is what is under test here, not app creation.
# ==============================================================================

data "microsoft365_graph_beta_device_and_app_management_mobile_app" "ios_store_apps_005" {
  list_all        = true
  app_type_filter = "iosStoreApp"
}

locals {
  published_ios_store_app_id_005 = [
    for app in data.microsoft365_graph_beta_device_and_app_management_mobile_app.ios_store_apps_005.items :
    app.id if app.publishing_state == "published"
  ][0]
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
  mobile_app_id = local.published_ios_store_app_id_005
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
