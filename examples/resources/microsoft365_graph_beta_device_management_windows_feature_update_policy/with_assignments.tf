# ==============================================================================
# Windows Feature Update Policy with Assignments
# ==============================================================================
# This example demonstrates how to deploy a Windows feature update policy with
# group-based assignments. It includes both inclusion and exclusion groups.

# Windows Feature Update Policy with multiple assignments
resource "microsoft365_graph_beta_device_management_windows_feature_update_policy" "with_assignments" {
  display_name                                            = "Windows 11 25H2 - With Assignments"
  description                                             = "Feature update deployment with targeted assignments"
  feature_update_version                                  = "Windows 11, version 25H2"
  install_feature_updates_optional                        = true
  install_latest_windows10_on_windows11_ineligible_device = true

  role_scope_tag_ids = ["0", "1"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "44444444-4444-4444-4444-444444444444"
    }
  ]
}
