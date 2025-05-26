# Terms and Conditions configuration
resource "microsoft365_graph_beta_device_management_terms_and_conditions" "terms" {
  display_name = "Enterprise Security and Compliance Policy"
  description  = "Terms and conditions for enterprise device management and security compliance"
  title        = "Enterprise Device Management Agreement"

  body_text = <<-EOT
    ENTERPRISE DEVICE MANAGEMENT AGREEMENT
    
    This agreement governs the use of company-managed devices and access to corporate resources.
    
    ACCEPTABLE USE POLICY:
    • Devices must be used primarily for business purposes
    • Personal use is limited to incidental, reasonable use
    • All software installations require IT approval
    • Regular security updates must be maintained
    
    SECURITY REQUIREMENTS:
    • Strong passwords or biometric authentication required
    • Device encryption must remain enabled
    • VPN connection required for remote access
    • Immediate reporting of lost or stolen devices
    
    MONITORING AND COMPLIANCE:
    • Company reserves the right to monitor device usage
    • Regular compliance scans will be performed
    • Non-compliance may result in device restrictions
    
    DATA PROTECTION:
    • Company data remains property of the organization
    • Personal data will be separated where possible
    • Device wipe may be performed upon termination
    
    By accepting these terms, you acknowledge understanding and agree to comply with all stated policies.
  EOT

  acceptance_statement = <<-EOT
    I understand that by accepting these terms and conditions, I am agreeing to comply with all company policies regarding device usage, security requirements, and data protection. I acknowledge that violation of these terms may result in disciplinary action up to and including termination of employment and/or device access restrictions.
  EOT

  role_scope_tag_ids = ["0", "1", "2"] # Default scope + custom scopes

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}