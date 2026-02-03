# Acceptance test: Windows Autopilot - HoloLens with Group Assignment
# Full dependency chain: random_string -> group -> time_sleep -> autopilot_profile

resource "random_string" "test_id_04" {
  length  = 8
  special = false
  upper   = false
}

# Test Group 1 - Assignment Target
resource "microsoft365_graph_beta_groups_group" "test_group_1" {
  display_name     = "acc-test-autopilot-04-group1-${random_string.test_id_04.result}"
  description      = "Test group for autopilot deployment profile acceptance test"
  mail_nickname    = "acc-autopilot-04-g1-${random_string.test_id_04.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

# Wait for eventual consistency after group creation
resource "time_sleep" "wait_for_groups" {
  depends_on = [microsoft365_graph_beta_groups_group.test_group_1]

  create_duration = "30s"
}

# HoloLens Deployment Profile Example
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "hololens_with_all_device_assignment" {
  display_name                                 = "acc_test_hololens_with_all_device_assignment"
  description                                  = "hololens autopilot profile with hk locale and group assignment"
  device_name_template                         = "thing-%RAND:2%"
  locale                                       = "zh-HK"
  preprovisioning_allowed                      = false
  device_type                                  = "holoLens"
  hardware_hash_extraction_enabled             = false
  role_scope_tag_ids                           = ["0"]
  device_join_type                             = "microsoft_entra_joined"
  hybrid_azure_ad_join_skip_connectivity_check = false

  out_of_box_experience_setting = {
    device_usage_type               = "shared"
    privacy_settings_hidden         = true
    eula_hidden                     = true
    user_type                       = "standard"
    keyboard_selection_page_skipped = true
  }

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.test_group_1.id
    }
  ]

  depends_on = [time_sleep.wait_for_groups]
}