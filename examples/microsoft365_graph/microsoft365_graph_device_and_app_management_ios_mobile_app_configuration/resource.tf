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

# Example: iOS Mobile App Configuration with XML settings
resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "example_xml" {
  display_name = "Example iOS App Config - XML"
  description  = "Example iOS mobile app configuration using XML settings"
  
  # List of app bundle IDs this configuration applies to
  targeted_mobile_apps = [
    "com.example.myapp",
    "com.example.anotherapp"
  ]
  
  # Base64 encoded XML configuration
  encoded_setting_xml = base64encode(<<XML
<dict>
    <key>ServerURL</key>
    <string>https://api.example.com</string>
    <key>EnableDebugMode</key>
    <false/>
    <key>MaxRetries</key>
    <integer>3</integer>
</dict>
XML
  )
}

# Example: iOS Mobile App Configuration with structured settings
resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "example_settings" {
  display_name = "Example iOS App Config - Structured"
  description  = "Example iOS mobile app configuration using structured settings"
  
  targeted_mobile_apps = [
    "com.contoso.salesapp"
  ]
  
  # Structured app configuration settings
  settings {
    app_config_key       = "serverUrl"
    app_config_key_type  = "stringType"
    app_config_key_value = "https://api.contoso.com"
  }
  
  settings {
    app_config_key       = "enableSync"
    app_config_key_type  = "booleanType"
    app_config_key_value = "true"
  }
  
  settings {
    app_config_key       = "syncInterval"
    app_config_key_type  = "integerType"
    app_config_key_value = "300"
  }
  
  settings {
    app_config_key       = "apiVersion"
    app_config_key_type  = "realType"
    app_config_key_value = "2.5"
  }
}

# Example: iOS Mobile App Configuration with assignments
data "microsoft365_graph_beta_groups_group" "sales_team" {
  display_name = "Sales Team"
}

data "microsoft365_graph_beta_groups_group" "test_devices" {
  display_name = "Test Devices"
}

resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "example_assigned" {
  display_name = "Sales App Configuration"
  description  = "Configuration for the sales team iOS app"
  
  targeted_mobile_apps = [
    "com.contoso.salesapp"
  ]
  
  settings {
    app_config_key       = "environment"
    app_config_key_type  = "stringType"
    app_config_key_value = "production"
  }
  
  settings {
    app_config_key       = "cacheEnabled"
    app_config_key_type  = "booleanType"
    app_config_key_value = "true"
  }
  
  # Assign to all licensed users
  assignments {
    target {
      odata_type = "#microsoft.graph.allLicensedUsersAssignmentTarget"
    }
  }
  
  # Also assign to a specific group
  assignments {
    target {
      odata_type = "#microsoft.graph.groupAssignmentTarget"
      group_id   = data.microsoft365_graph_beta_groups_group.sales_team.id
    }
  }
  
  # Exclude test devices group
  assignments {
    target {
      odata_type = "#microsoft.graph.exclusionGroupAssignmentTarget"
      group_id   = data.microsoft365_graph_beta_groups_group.test_devices.id
    }
  }
}

# Example: Minimal iOS Mobile App Configuration
resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "example_minimal" {
  display_name = "Minimal iOS Config"
}

# Example: iOS Mobile App Configuration with custom timeouts
resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "example_timeouts" {
  display_name = "iOS Config with Custom Timeouts"
  
  targeted_mobile_apps = [
    "com.example.app"
  ]
  
  timeouts {
    create = "5m"
    read   = "2m"
    update = "5m"
    delete = "3m"
  }
}