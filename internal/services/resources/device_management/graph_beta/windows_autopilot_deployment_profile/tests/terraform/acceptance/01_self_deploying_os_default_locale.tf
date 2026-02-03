# Acceptance test: Windows Autopilot - User Driven with OS Default Locale
# Full dependency chain: random_string -> groups -> time_sleep -> autopilot_profile

resource "random_string" "test_id_01" {
  length  = 8
  special = false
  upper   = false
}

# Test Group 1 - Assignment Target
resource "microsoft365_graph_beta_groups_group" "test_group_1" {
  display_name     = "acc-test-autopilot-01-group1-${random_string.test_id_01.result}"
  description      = "Test group for autopilot deployment profile acceptance test"
  mail_nickname    = "acc-autopilot-01-g1-${random_string.test_id_01.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

# Test Group 2 - Assignment Target
resource "microsoft365_graph_beta_groups_group" "test_group_2" {
  display_name     = "acc-test-autopilot-01-group2-${random_string.test_id_01.result}"
  description      = "Test group for autopilot deployment profile acceptance test"
  mail_nickname    = "acc-autopilot-01-g2-${random_string.test_id_01.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

# Test Group 3 - Exclusion Target
resource "microsoft365_graph_beta_groups_group" "test_group_3" {
  display_name     = "acc-test-autopilot-01-group3-${random_string.test_id_01.result}"
  description      = "Test group for autopilot deployment profile acceptance test"
  mail_nickname    = "acc-autopilot-01-g3-${random_string.test_id_01.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

# Wait for eventual consistency after group creation
resource "time_sleep" "wait_for_groups" {
  depends_on = [
    microsoft365_graph_beta_groups_group.test_group_1,
    microsoft365_graph_beta_groups_group.test_group_2,
    microsoft365_graph_beta_groups_group.test_group_3
  ]

  create_duration = "30s"
}

# User-Driven Deployment Profile Example
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "user_driven" {
  display_name                                 = "acc test user driven autopilot profile with os default locale"
  description                                  = "user driven autopilot profile with os default locale"
  device_name_template                         = "thing-%RAND:5%"
  locale                                       = "os-default"
  preprovisioning_allowed                      = true
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0", "1"]
  device_join_type                             = "microsoft_entra_joined"
  hybrid_azure_ad_join_skip_connectivity_check = false // always false when using microsoft_entra_joined

  out_of_box_experience_setting = {
    device_usage_type               = "singleUser"
    privacy_settings_hidden         = true       // Privacy settings
    eula_hidden                     = true       // Microsoft Software License Terms
    user_type                       = "standard" // standard or administrator
    keyboard_selection_page_skipped = true       // Automatically configure keyboard
  }

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.test_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.test_group_2.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.test_group_3.id
    }
  ]

  depends_on = [time_sleep.wait_for_groups]
}