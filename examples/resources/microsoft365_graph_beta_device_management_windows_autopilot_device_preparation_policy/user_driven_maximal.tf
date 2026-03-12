# ==============================================================================
# User-Driven Mode - Maximal Configuration
# ==============================================================================
# This example demonstrates a user-driven deployment with enhanced mode features,
# OOBE settings, apps, and scripts. User-driven mode provides more control and
# customization options for the end-user experience.

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "ud_maximal" {
  name                  = "Autopilot DPP - User-Driven Maximal"
  description           = "User-driven mode maximal configuration with enhanced mode features"
  role_scope_tag_ids    = ["0"]
  device_security_group = "12345678-1234-1234-1234-123456789012" # Replace with your device security group ID

  deployment_settings = {
    deployment_mode = "enrollment_autopilot_dpp_deploymentmode_1" # Enhanced
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_0" # User-driven
    account_type    = "enrollment_autopilot_dpp_accountype_1"     # Standard user
  }

  # Configure the out-of-box experience
  oobe_settings = {
    timeout_in_minutes   = 120
    custom_error_message = "Please contact your IT administrator for assistance with device setup."
    allow_skip           = true
    allow_diagnostics    = true
  }

  allowed_apps = [
    {
      app_id   = "12345678-1234-1234-1234-123456789012" # Replace with your WinGet app ID
      app_type = "winGetApp"
    },
    {
      app_id   = "12345678-1234-1234-1234-234567890123" # Replace with your WinGet app ID
      app_type = "win32LobApp"
    },
    {
      app_id   = "12345678-1234-1234-1234-345678901234" # Replace with your Win32CatalogApp app ID
      app_type = "win32CatalogApp"
    },
    {
      app_id   = "12345678-1234-1234-1234-456789012345" # Replace with your OfficeSuiteApp app ID
      app_type = "officeSuiteApp"
    },
  ]

  allowed_scripts = [
    "87654321-4321-4321-4321-210987654321",
    "87654321-4321-4321-4321-210987654322",
  ]

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
