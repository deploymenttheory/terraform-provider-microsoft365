# Example 1: Apply region configuration changes to existing Cloud PCs (default behavior)
action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "apply_region_changes" {
  provisioning_policy_id = "12345678-1234-1234-1234-123456789abc"
}

# Example 2: Apply region configuration explicitly
action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "apply_region_explicit" {
  provisioning_policy_id = "12345678-1234-1234-1234-123456789abc"
  policy_settings        = "region"
}

# Example 3: Apply single sign-on configuration changes to existing Cloud PCs
action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "apply_sso_changes" {
  provisioning_policy_id = "87654321-4321-4321-4321-cba987654321"
  policy_settings        = "singleSignOn"
}

# Example 4: Apply changes to Frontline shared Cloud PC with reserve percentage
action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "apply_frontline_with_reserve" {
  provisioning_policy_id = "abcdef12-3456-7890-abcd-ef1234567890"
  policy_settings        = "region"
  reserve_percentage     = 20 # Keep 20% of Cloud PCs available
}

# Example 5: Apply changes with custom timeouts
action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "apply_with_timeouts" {
  provisioning_policy_id = "fedcba98-7654-3210-fedc-ba9876543210"
  policy_settings        = "singleSignOn"

  timeouts = {
    run = "10m"
  }
}

# Example 6: Using with a provisioning policy resource reference
resource "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy" "example" {
  display_name        = "Example Provisioning Policy"
  description         = "Policy for IT department Cloud PCs"
  # ... other configuration
}

action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "apply_to_resource" {
  provisioning_policy_id = microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.example.id
  policy_settings        = "region"
}

# Example 7: Frontline shared Cloud PC - Maximum reserve percentage
action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "frontline_max_reserve" {
  provisioning_policy_id = "11223344-5566-7788-99aa-bbccddeeff00"
  policy_settings        = "region"
  reserve_percentage     = 99 # Keep 99% of Cloud PCs available (maximum value)
}

# Example 8: Frontline shared Cloud PC - Zero reserve percentage
action "microsoft365_graph_beta_device_management_windows_365_apply_cloud_pc_provisioning_policy" "frontline_no_reserve" {
  provisioning_policy_id = "00112233-4455-6677-8899-aabbccddeeff"
  policy_settings        = "region"
  reserve_percentage     = 0 # Apply to all Cloud PCs
}

