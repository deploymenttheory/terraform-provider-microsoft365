# Basic Windows Autopilot Deployment Profile
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "entra_joined" {
  display_name = "Corporate Windows Autopilot Profile"
  description  = "Windows Autopilot deployment profile for corporate devices with OOBE customization"

  # Device join configuration
  device_join_type = "microsoft_entra_joined"

  # Device configuration
  device_type                      = "windowsPc"
  device_name_template             = "CORP-%SERIAL%"
  locale                           = "en-GB"
  preprovisioning_allowed          = true
  hardware_hash_extraction_enabled = true

  # Role scope tags
  role_scope_tag_ids = ["0", "9", "8"]

  # Azure AD configuration
  hybrid_azure_ad_join_skip_connectivity_check = false

  out_of_box_experience_setting = {
    privacy_settings_hidden         = false
    eula_hidden                     = false
    user_type                       = "administrator"
    device_usage_type               = "singleUser"
    keyboard_selection_page_skipped = true
    escape_link_hidden              = true
  }

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}

# Hybrid Domain Join Windows Autopilot Deployment Profile
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile" "hybrid_domain_join" {
  display_name = "hybrid domain join"
  description  = "test"

  # Device join configuration
  device_join_type = "microsoft_entra_hybrid_joined"

  # Device configuration
  device_type                      = "windowsPc"
  device_name_template             = ""
  locale                           = "os-default"
  preprovisioning_allowed          = true
  hardware_hash_extraction_enabled = true

  # Role scope tags
  role_scope_tag_ids = ["0"]

  # Azure AD configuration
  hybrid_azure_ad_join_skip_connectivity_check = true

  out_of_box_experience_setting = {
    privacy_settings_hidden         = false
    eula_hidden                     = false
    user_type                       = "administrator"
    device_usage_type               = "singleUser"
    keyboard_selection_page_skipped = true
    escape_link_hidden              = true
  }

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}
