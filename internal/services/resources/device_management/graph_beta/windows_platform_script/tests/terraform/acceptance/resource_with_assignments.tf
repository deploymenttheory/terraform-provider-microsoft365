resource "random_uuid" "with_assignments" {
}

resource "microsoft365_graph_beta_device_management_windows_platform_script" "with_assignments" {
  display_name           = "Acceptance - Windows Platform Script with Assignments"
  description            = "Test description for script with assignments"
  file_name              = "acceptance-assignments-script.ps1"
  script_content         = "Write-Host 'Script with assignments for acceptance testing'"
  run_as_account         = "system"
  enforce_signature_check = false
  run_as_32_bit          = false
  
  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    }
  ]

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }

  lifecycle {
    ignore_changes = [
      role_scope_tag_ids
    ]
  }
}