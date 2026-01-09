# Scenario 6: Step 1 - Resource with no assignments
resource "microsoft365_graph_beta_device_management_macos_platform_script" "add_maximal_assignments" {
  display_name   = "unit-test-add-maximal-assignments"
  file_name      = "test_maximal_assignments.sh"
  script_content = "#!/bin/bash\necho 'Min Test'\nexit 0"
  run_as_account = "system"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
