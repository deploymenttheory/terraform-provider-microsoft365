terraform {
  required_providers {
    microsoft365 = {
      source = "registry.terraform.io/terracurl/microsoft365"
    }
  }
}

provider "microsoft365" {
  use_cli = true
}

# Example: Look up iOS mobile app configuration by display name
data "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "by_name" {
  display_name = "Sales App Configuration"
}

# Example: Look up iOS mobile app configuration by ID
data "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "by_id" {
  id = "12345678-1234-1234-1234-123456789012"
}

# Example: Use data source to reference configuration details
output "config_details" {
  value = {
    id                   = data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.by_name.id
    display_name         = data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.by_name.display_name
    description          = data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.by_name.description
    targeted_mobile_apps = data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.by_name.targeted_mobile_apps
    version              = data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.by_name.version
    created_date_time    = data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.by_name.created_date_time
    last_modified        = data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.by_name.last_modified_date_time
  }
}

# Example: Show configuration settings
output "config_settings" {
  value = [
    for setting in data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.by_name.settings : {
      key       = setting.app_config_key
      type      = setting.app_config_key_type
      value     = setting.app_config_key_value
    }
  ]
}

# Example: Show configuration assignments
output "config_assignments" {
  value = [
    for assignment in data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.by_name.assignments : {
      id          = assignment.id
      target_type = assignment.target.odata_type
      group_id    = try(assignment.target.group_id, null)
    }
  ]
}

# Example: Reference XML settings (note: this is sensitive data)
output "has_xml_settings" {
  value = data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.by_name.encoded_setting_xml != null
  description = "Whether the configuration has XML settings defined"
}

# Example: Use data source with custom timeouts
data "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "with_timeouts" {
  display_name = "iOS App Config with Timeouts"
  
  timeouts {
    read = "5m"
  }
}

# Example: Use data source to copy configuration to another resource
data "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "source_config" {
  display_name = "Source iOS Config"
}

resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "copied_config" {
  display_name = "${data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.source_config.display_name} - Copy"
  description  = "Copied from: ${data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.source_config.description}"
  
  targeted_mobile_apps = data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.source_config.targeted_mobile_apps
  
  # Copy settings if they exist
  dynamic "settings" {
    for_each = data.microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.source_config.settings
    content {
      app_config_key       = settings.value.app_config_key
      app_config_key_type  = settings.value.app_config_key_type
      app_config_key_value = settings.value.app_config_key_value
    }
  }
}