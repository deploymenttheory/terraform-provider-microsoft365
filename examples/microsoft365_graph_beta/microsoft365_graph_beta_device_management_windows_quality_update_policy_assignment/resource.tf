# Basic group assignment with rollout settings
resource "microsoft365_graph_beta_device_management_windows_quality_update_policy_assignment" "group_example" {
  windows_quality_update_policy_id = "00000000-0000-0000-0000-000000000001"

  target = {
    target_type = "groupAssignment"
    group_id    = "00000000-0000-0000-0000-000000000002"
  }

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# SCCM collection assignment with comprehensive settings
resource "microsoft365_graph_beta_device_management_windows_quality_update_policy_assignment" "sccm_example" {
  windows_quality_update_policy_id = "00000000-0000-0000-0000-000000000004"

  target = {
    target_type   = "configurationManagerCollection"
    collection_id = "MEMABCDEF01"
  }

}

# Exclusion group assignment (minimal configuration)
resource "microsoft365_graph_beta_device_management_windows_quality_update_policy_assignment" "exclusion_example" {
  windows_quality_update_policy_id = "00000000-0000-0000-0000-000000000007"

  target = {
    target_type = "exclusionGroupAssignment"
    group_id    = "00000000-0000-0000-0000-000000000008"
  }

}