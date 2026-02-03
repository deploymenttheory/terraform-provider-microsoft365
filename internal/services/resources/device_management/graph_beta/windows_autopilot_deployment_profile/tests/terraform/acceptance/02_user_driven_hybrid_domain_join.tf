# Acceptance test: Windows Autopilot - User Driven Hybrid Domain Join
# Full dependency chain: random_string -> group -> time_sleep -> autopilot_profile

resource "random_string" "test_id_02" {
  length  = 8
  special = false
  upper   = false
}

# Test Group 1 - Assignment Target
resource "microsoft365_graph_beta_groups_group" "test_group_1" {
  display_name     = "acc-test-autopilot-02-group1-${random_string.test_id_02.result}"
  description      = "Test group for autopilot deployment profile acceptance test"
  mail_nickname    = "acc-autopilot-02-g1-${random_string.test_id_02.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

# Wait for eventual consistency after group creation
resource "time_sleep" "wait_for_groups" {
  depends_on = [microsoft365_graph_beta_groups_group.test_group_1]

  create_duration = "30s"
}

# User-Driven with Japanese Language and Allow Pre-provisioned Deployment
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "user_driven_japanese_preprovisioned_with_assignments" {
  display_name                                 = "acc_test_user_driven_japanese_preprovisioned"
  description                                  = "user driven autopilot profile with japanese locale and allow pre provisioned deployment"
  locale                                       = "ja-JP"
  preprovisioning_allowed                      = true
  device_type                                  = "windowsPc"
  hardware_hash_extraction_enabled             = true
  role_scope_tag_ids                           = ["0"]
  device_join_type                             = "microsoft_entra_hybrid_joined"
  hybrid_azure_ad_join_skip_connectivity_check = true

  out_of_box_experience_setting = {
    device_usage_type               = "singleUser"
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