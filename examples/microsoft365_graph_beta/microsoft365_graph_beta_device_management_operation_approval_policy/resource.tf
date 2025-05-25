# Basic Operation Approval Policy configuration
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "basic_approval" {
  display_name = "Basic App Approval Policy"
  description  = "Requires approval for application deployments"
  
  policy_set {
    policy_type = "app"
  }
  
  approver_group_ids = [
    "12345678-1234-1234-1234-123456789012"  # IT Administrators group
  ]
}

# Comprehensive Operation Approval Policy for scripts
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "script_approval" {
  display_name     = "PowerShell Script Approval Policy"
  description      = "Multi-level approval required for PowerShell script execution on managed devices"
  policy_type      = "script"
  policy_platform  = "windows10AndLater"
  
  policy_set {
    policy_type      = "script"
    policy_platform  = "windows10AndLater"
  }
  
  approver_group_ids = [
    "11111111-1111-1111-1111-111111111111",  # Security Team
    "22222222-2222-2222-2222-222222222222",  # IT Management
    "33333333-3333-3333-3333-333333333333"   # Compliance Team
  ]
  
  timeouts {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}