resource "microsoft365_graph_beta_device_management_windows_quality_update_expedite_policy" "exclusion_assignment" {
  display_name = "Test Exclusion Assignment Windows Quality Update Expedite Policy - Unique"

  assignments = [
    { type = "exclusionGroupAssignmentTarget", group_id = "33333333-3333-3333-3333-333333333333" }
  ]
}


