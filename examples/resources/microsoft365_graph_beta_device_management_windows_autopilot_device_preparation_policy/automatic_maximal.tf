# ==============================================================================
# Automatic Mode - Maximal Configuration
# ==============================================================================
# This example demonstrates an automatic deployment with apps and scripts.
# Automatic mode requires minimal configuration and is ideal for shared devices.

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "auto_maximal" {
  name               = "Autopilot DPP - Automatic Maximal"
  description        = "Automatic mode maximal configuration with apps and scripts"
  role_scope_tag_ids = ["0"]

  deployment_settings = {
    deployment_type = "enrollment_autopilot_dpp_deploymenttype_1" # Automatic
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

  timeouts = {
    create = "180s"
    read   = "30s"
    update = "180s"
    delete = "60s"
  }
}
