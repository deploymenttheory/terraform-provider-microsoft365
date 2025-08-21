# Example: English Notification Template with Full Branding
resource "microsoft365_graph_beta_device_management_device_compliance_notification_templates" "english" {
  display_name     = "English Compliance Notification"
  branding_options = ["includeCompanyLogo", "includeCompanyName", "includeContactInformation"]

  role_scope_tag_ids = ["0", "1"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Immediate Action Required: Device Compliance"
      message_template = <<-EOT
        Dear {UserName},

        Your device '{DeviceName}' (Serial: {DeviceSerialNumber}) has been identified as non-compliant with our security policies.

        Compliance Issues:
        {ComplianceReasons}

        Required Actions:
        1. Update your device to the latest security patches
        2. Enable BitLocker encryption if not already enabled
        3. Ensure Windows Defender is active and up-to-date
        4. Contact IT support if you need assistance

        Device Details:
        - Device Name: {DeviceName}
        - Operating System: {DeviceOSVersion}
        - Last Check-in: {LastCheckInTime}

        Failure to address these issues within 24 hours may result in restricted access to company resources.

        For immediate assistance, contact:
        - IT Helpdesk: +1-555-0123
        - Email: it-support@company.com
        - Portal: https://company.com/it-support

        Best regards,
        IT Security Team
        Company Name
      EOT
      is_default       = true
    }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "5m"
  }
}