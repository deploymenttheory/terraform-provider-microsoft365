# ==============================================================================
# User-Driven Mode - With Maximal Assignments
# ==============================================================================
# This example demonstrates how to configure policy assignments using both
# group-based targeting and all licensed users. Assignments determine which
# users will have this policy applied during device enrollment.

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "ud_max_assign" {
  name                  = "Autopilot DPP - User-Driven with Assignments"
  description           = "User-driven mode with maximal assignments demonstrating group and all licensed users targeting"
  role_scope_tag_ids    = ["0"]
  device_security_group = "12345678-1234-1234-1234-123456789012" # Replace with your device security group ID

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_1" # Enhanced
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0" # User-driven
    account_type    = "enrollment_autopilot_dpp_accountype_1"     # Standard user
  }

  # Configure the out-of-box experience
  oobe_settings = {
    timeout_in_minutes   = 90
    custom_error_message = "Please contact your IT administrator for assistance."
    allow_skip           = true
    allow_diagnostics    = true
  }

  # Allow specific apps during device preparation
  allowed_apps = [
    {
      app_id   = "12345678-1234-1234-1234-123456789012" # Replace with your WinGet app ID
      app_type = "winGetApp"
    }
  ]

  # Allow specific scripts during device preparation
  allowed_scripts = [
    "87654321-4321-4321-4321-210987654321" # Replace with your Windows Platform Script ID
  ]

  # Configure policy assignments
  # This demonstrates multiple assignment types:
  # - All licensed users (applies to all users with Intune licenses)
  # - Specific security groups (targeted deployment)
  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "11111111-1111-1111-1111-111111111111" # Replace with your group ID
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "22222222-2222-2222-2222-222222222222" # Replace with your group ID
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "33333333-3333-3333-3333-333333333333" # Replace with your group ID
    },
  ]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
