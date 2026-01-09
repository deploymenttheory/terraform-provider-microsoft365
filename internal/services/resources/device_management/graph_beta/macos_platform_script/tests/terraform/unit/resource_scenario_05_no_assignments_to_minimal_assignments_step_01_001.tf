# Scenario 5: Step 1 - Resource with no assignments
resource "microsoft365_graph_beta_device_management_macos_platform_script" "add_minimal_assignment" {
  display_name   = "unit-test-add-minimal-assignment"
  file_name      = "test_minimal_assignment.sh"
  script_content = "#!/bin/bash\necho 'Test'\nexit 0"
  run_as_account = "system"

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}
