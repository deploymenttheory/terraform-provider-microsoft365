resource "microsoft365_graph_beta_device_management_windows_autopilot_device_preparation_policy" "example" {
  name                = "Windows Autopilot Device Preparation Policy"
  description         = "Example Windows Autopilot Device Preparation Policy"
  role_scope_tag_ids  = ["0"]
  device_security_group = "12345678-1234-1234-1234-123456789012" # ID of the assigned security device group that must have the 'Intune Provisioning Client' service principal (AppId: f1346770-5b25-470b-88bd-d5744ab7952c) as its owner

  deployment_settings = {
    deployment_mode  = "enrollment_autopilot_dpp_deploymentmode_0" # Standard mode
    deployment_type  = "enrollment_autopilot_dpp_deploymenttype_0" # User-driven
    join_type        = "enrollment_autopilot_dpp_jointype_0"       # Azure AD joined
    account_type     = "enrollment_autopilot_dpp_accountype_0"     # Standard user
  }

  oobe_settings = {
    timeout_in_minutes   = 60
    custom_error_message = "Contact your IT department for assistance."
    allow_skip           = false
    allow_diagnostics    = true
  }

  allowed_apps = [
    {
      app_id   = "12345678-1234-1234-1234-123456789012" # Example Windows Store app
      app_type = "winGetApp"
    },
    {
      app_id   = "23456789-2345-2345-2345-234567890123" # Example Win32 app
      app_type = "win32LobApp"
    },
    {
      app_id   = "34567890-3456-3456-3456-345678901234" # Example Office app
      app_type = "officeSuiteApp"
    },
    {
      app_id   = "45678901-4567-4567-4567-456789012345" # Example Universal app
      app_type = "windowsUniversalAppX"
    }
  ]

  allowed_scripts = [
    "12345678-1234-1234-1234-123456789012", # Example script ID
    "12345678-1234-1234-1234-123456789012",
  ]

  assignments = {
    include_group_ids = [
      "12345678-1234-1234-1234-123456789012", # Example group ID
      "12345678-1234-1234-1234-123456789012",
    ]
  }
} 