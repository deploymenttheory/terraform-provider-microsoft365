resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "example" {
  display_name = "Maximal iOS App Config"
  description  = "A comprehensive iOS mobile app configuration with all possible settings"
  
  targeted_mobile_apps = [
    "com.example.app1",
    "com.example.app2",
    "com.example.app3"
  ]
  
  encoded_setting_xml = base64encode(<<XML
<?xml version="1.0" encoding="utf-8"?>
<configuration>
  <bundle_id>com.example.app</bundle_id>
  <server_url>https://api.example.com</server_url>
  <feature_flags>
    <flag name="new_ui" enabled="true"/>
    <flag name="beta_features" enabled="false"/>
  </feature_flags>
</configuration>
XML
  )
  
  settings {
    app_config_key       = "server_url"
    app_config_key_type  = "stringType"
    app_config_key_value = "https://api.example.com"
  }
  
  settings {
    app_config_key       = "timeout_seconds"
    app_config_key_type  = "integerType"
    app_config_key_value = "30"
  }
  
  settings {
    app_config_key       = "enable_analytics"
    app_config_key_type  = "booleanType"
    app_config_key_value = "true"
  }
  
  settings {
    app_config_key       = "api_version"
    app_config_key_type  = "realType"
    app_config_key_value = "2.5"
  }
  
  settings {
    app_config_key       = "access_token"
    app_config_key_type  = "tokenType"
    app_config_key_value = "{{token}}"
  }
  
  assignments {
    target {
      odata_type = "#microsoft.graph.allLicensedUsersAssignmentTarget"
    }
  }
  
  assignments {
    target {
      odata_type = "#microsoft.graph.groupAssignmentTarget"
      group_id   = "11111111-2222-3333-4444-555555555555"
    }
  }
  
  assignments {
    target {
      odata_type = "#microsoft.graph.exclusionGroupAssignmentTarget"
      group_id   = "66666666-7777-8888-9999-000000000000"
    }
  }
  
  role_scope_tag_ids = ["0", "1", "2"]
  
  timeouts {
    create = "30m"
    read   = "5m"
    update = "30m"
    delete = "10m"
  }
}