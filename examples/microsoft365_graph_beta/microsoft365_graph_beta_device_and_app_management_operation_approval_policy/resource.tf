# Basic example for device wipe approval
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "device_wipe_approval" {
  display_name = "Device Wipe Approval Policy"
  description  = "Requires approval before wiping corporate devices"
  policy_type  = "deviceWipe"
  
  policy_set = {
    policy_type = "deviceWipe"
  }
  
  approver_group_ids = [
    "12345678-1234-1234-1234-123456789012", # IT Security Team
    "87654321-4321-4321-4321-210987654321"  # Device Management Team
  ]
}

# Windows-specific compliance policy approval
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "windows_compliance_approval" {
  display_name    = "Windows Compliance Policy Approval"
  description     = "Approval required for Windows compliance policy changes"
  policy_type     = "compliancePolicy"
  policy_platform = "windows10AndLater"
  
  policy_set = {
    policy_type     = "compliancePolicy"
    policy_platform = "windows10AndLater"
  }
  
  approver_group_ids = [
    "11111111-1111-1111-1111-111111111111", # Compliance Team
    "22222222-2222-2222-2222-222222222222"  # Security Managers
  ]
}

# iOS/iPadOS app approval policy
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "ios_app_approval" {
  display_name    = "iOS App Deployment Approval"
  description     = "Requires approval before deploying apps to iOS/iPadOS devices"
  policy_type     = "app"
  policy_platform = "iOSiPadOS"
  
  policy_set = {
    policy_type     = "app"
    policy_platform = "iOSiPadOS"
  }
  
  approver_group_ids = [
    "33333333-3333-3333-3333-333333333333"  # Mobile App Team
  ]
}

# Script deployment approval for macOS
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "macos_script_approval" {
  display_name    = "macOS Script Approval Policy"
  description     = "All scripts deployed to macOS devices require prior approval"
  policy_type     = "script"
  policy_platform = "macOS"
  
  policy_set = {
    policy_type     = "script"
    policy_platform = "macOS"
  }
  
  approver_group_ids = [
    "44444444-4444-4444-4444-444444444444", # macOS Administrators
    "55555555-5555-5555-5555-555555555555"  # Security Team Lead
  ]
}

# Device retire approval for Android Enterprise
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "android_retire_approval" {
  display_name    = "Android Device Retirement Approval"
  policy_type     = "deviceRetire"
  policy_platform = "androidEnterprise"
  
  policy_set = {
    policy_type     = "deviceRetire"
    policy_platform = "androidEnterprise"
  }
  
  approver_group_ids = [
    "66666666-6666-6666-6666-666666666666"  # Android Device Team
  ]
}

# Endpoint security policy approval (platform agnostic)
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "endpoint_security_approval" {
  display_name    = "Endpoint Security Policy Approval"
  description     = "Multi-layered approval for critical endpoint security changes"
  policy_type     = "endpointSecurityPolicy"
  policy_platform = "notApplicable"  # Can also be omitted as this is the default
  
  policy_set = {
    policy_type     = "endpointSecurityPolicy"
    policy_platform = "notApplicable"
  }
  
  approver_group_ids = [
    "77777777-7777-7777-7777-777777777777", # Security Operations Center
    "88888888-8888-8888-8888-888888888888", # CISO Office
    "99999999-9999-9999-9999-999999999999"  # IT Leadership
  ]
}

# Device action approval with custom timeouts
resource "microsoft365_graph_beta_device_management_operation_approval_policy" "device_action_approval" {
  display_name = "Critical Device Actions Approval"
  description  = "Approval required for critical device management actions"
  policy_type  = "deviceAction"
  
  policy_set = {
    policy_type = "deviceAction"
  }
  
  approver_group_ids = [
    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
  ]
  
  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example using variables for reusability
variable "security_team_group_id" {
  description = "The Azure AD group ID for the security team"
  type        = string
  default     = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
}

variable "it_admin_group_id" {
  description = "The Azure AD group ID for IT administrators"
  type        = string
  default     = "cccccccc-cccc-cccc-cccc-cccccccccccc"
}

resource "microsoft365_graph_beta_device_management_operation_approval_policy" "configuration_policy_approval" {
  display_name = "Configuration Policy Approval"
  description  = "Requires approval for device configuration policy changes"
  policy_type  = "configurationPolicy"
  
  policy_set = {
    policy_type = "configurationPolicy"
  }
  
  approver_group_ids = [
    var.security_team_group_id,
    var.it_admin_group_id
  ]
}