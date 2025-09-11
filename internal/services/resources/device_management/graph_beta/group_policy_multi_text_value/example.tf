# Example usage of the group_policy_multi_text_value resource with simplified auto-discovery

# Minimal configuration - auto-discovers everything
resource "microsoft365_graph_beta_device_management_group_policy_multi_text_value" "trusted_sites" {
  group_policy_configuration_id = "12345678-1234-1234-1234-123456789012"
  
  # Just specify the policy name and class type - IDs are auto-discovered
  policy_name = "Allow automatic full screen on specified sites"
  class_type  = "machine"
  
  values = [
    "https://intranet.company.com",
    "https://portal.company.com", 
    "https://sharepoint.company.com",
    "https://wiki.company.com"
  ]
  
  timeouts {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Advanced configuration with presentation selection
resource "microsoft365_graph_beta_device_management_group_policy_multi_text_value" "asr_rules" {
  group_policy_configuration_id = "12345678-1234-1234-1234-123456789012"
  
  policy_name        = "Configure Attack Surface Reduction rules"
  class_type         = "machine"
  presentation_index = 0  # Use first suitable presentation (default)
  
  values = [
    "7674ba52-37eb-4a4f-a9a1-f0f9a1619a2c",  # Block Adobe Reader from creating child processes
    "d4f940ab-401b-4efc-aadc-ad5f3c50688a",  # Block all Office applications from creating child processes
    "9e6c4e1f-7d60-472f-ba1a-a39ef669e4b2"   # Block credential stealing from the Windows local security authority subsystem
  ]
}

# Backward compatibility - you can still use the old explicit ID approach
resource "microsoft365_graph_beta_device_management_group_policy_multi_text_value" "explicit_ids" {
  group_policy_configuration_id    = "12345678-1234-1234-1234-123456789012"
  group_policy_definition_value_id = "87654321-4321-4321-4321-210987654321"  # Optional if using policy_name
  presentation_id                  = "abcdef12-3456-7890-abcd-ef1234567890"  # Optional if using policy_name
  
  values = ["value1", "value2", "value3"]
}
