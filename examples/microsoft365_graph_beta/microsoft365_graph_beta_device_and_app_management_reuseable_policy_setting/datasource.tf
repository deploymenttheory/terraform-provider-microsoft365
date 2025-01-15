// Data Source: Reusable Policy Settings
// Basic usage: lookup by display name
data "microsoft365_graph_beta_device_and_app_management_reuseable_policy_setting" "example" {
  display_name = "epm certificate"
}

// Output to verify data source
output "reuseable_policy_settings_details" {
  value = {
    id                                     = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.id
    display_name                           = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.display_name
    description                            = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.description
    settings                               = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.settings
    created_date_time                      = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.created_date_time
    last_modified_date_time                = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.last_modified_date_time
    version                                = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.version
    referencing_configuration_policies     = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.referencing_configuration_policies
    referencing_configuration_policy_count = data.microsoft365_graph_beta_device_and_app_management_reuseable_policy_settings.example.referencing_configuration_policy_count
  }
}
