# Get the WSL networking configuration setting
data "microsoft365_graph_beta_device_management_group_policy_category" "wsl_networking" {
  setting_name = "Configure default networking mode" // Define the group policy item you wish to return

  timeouts ={
    read = "5m"
  }
}

# Output the complete data structure
output "wsl_networking_setting" {
  description = "Complete group policy setting information from all three API calls"
  value = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking
}

# Access group policy top tier list 
output "category_info" {
  description = "Category information from the first API call"
  value = {
    id               = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.category.id
    display_name     = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.category.display_name
    is_root          = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.category.is_root
    ingestion_source = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.category.ingestion_source
    parent_category  = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.category.parent
  }
}

# Access group policy definition by id (from 2nd API call)
output "definition_info" {
  description = "Detailed policy definition from the second API call"
  value = {
    id                       = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.id
    display_name             = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.display_name
    explain_text             = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.explain_text
    category_path            = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.category_path
    class_type               = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.class_type
    policy_type              = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.policy_type
    version                  = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.version
    has_related_definitions  = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.has_related_definitions
    group_policy_category_id = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.group_policy_category_id
    min_device_csp_version   = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.min_device_csp_version
    min_user_csp_version     = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.min_user_csp_version
    supported_on             = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.supported_on
    last_modified_date_time  = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.definition.last_modified_date_time
  }
}

# Access presentation information (from 3rd API call) - now properly populated!
output "presentation_info" {
  description = "Presentation configuration from the third API call - dropdown with options"
  value = {
    presentation_id          = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].id
    odata_type              = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].odata_type
    label                   = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].label
    required                = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].required
    last_modified_date_time = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].last_modified_date_time
    
    # Dropdown-specific properties.
    default_item            = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].default_item
    available_options       = data.microsoft365_graph_beta_device_management_group_policy_category.wsl_networking.presentations[0].items
  }
}