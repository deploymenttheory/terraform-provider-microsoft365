# Example: Basic Notification Message Template
resource "microsoft365_graph_beta_device_management_notification_message_template" "basic" {
  display_name    = "Basic Compliance Notification"
  description     = "Basic notification template for device compliance violations"
  default_locale  = "en-US"
  branding_options = "includeCompanyLogo"
  
  role_scope_tag_ids = ["0"]

  localized_notification_messages = [
    {
      locale           = "en-US"
      subject          = "Device Compliance Issue Detected"
      message_template = "Hello {UserName},\n\nYour device '{DeviceName}' has been found to be non-compliant with company policies. Please take action to resolve the following issues:\n\n{ComplianceReasons}\n\nFor assistance, please contact IT support.\n\nThank you,\nIT Security Team"
      is_default       = true
    }
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example: Multi-language Notification Message Template
resource "microsoft365_graph_beta_device_management_notification_message_template" "multilingual" {
  display_name    = "Multi-language Compliance Notification"
  description     = "Notification template with multiple language support"
  default_locale  = "en-US"
  branding_options = "includeCompanyLogo"
  
  role_scope_tag_ids = ["0"]

  localized_notification_messages = [
    {
      locale           = "en-US"
      subject          = "Device Compliance Issue"
      message_template = "Hello {UserName},\n\nYour device '{DeviceName}' is not compliant. Please resolve: {ComplianceReasons}\n\nContact IT for help.\n\nIT Security Team"
      is_default       = true
    },
    {
      locale           = "es-ES"
      subject          = "Problema de Cumplimiento del Dispositivo"
      message_template = "Hola {UserName},\n\nTu dispositivo '{DeviceName}' no cumple las normas. Por favor resuelve: {ComplianceReasons}\n\nContacta con IT para ayuda.\n\nEquipo de Seguridad IT"
      is_default       = false
    },
    {
      locale           = "fr-FR"
      subject          = "Problème de Conformité de l'Appareil"
      message_template = "Bonjour {UserName},\n\nVotre appareil '{DeviceName}' n'est pas conforme. Veuillez résoudre: {ComplianceReasons}\n\nContactez l'IT pour aide.\n\nÉquipe de Sécurité IT"
      is_default       = false
    }
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example: Advanced Notification Template with Full Branding
resource "microsoft365_graph_beta_device_management_notification_message_template" "advanced" {
  display_name    = "Advanced Compliance Notification"
  description     = "Advanced notification template with comprehensive branding and device details"
  default_locale  = "en-US"
  branding_options = "includeCompanyLogo"
  
  role_scope_tag_ids = ["0", "1"]

  localized_notification_messages = [
    {
      locale           = "en-US"
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

# Output examples
output "basic_template_id" {
  description = "ID of the basic notification message template"
  value       = microsoft365_graph_beta_device_management_notification_message_template.basic.id
}

output "multilingual_template_id" {
  description = "ID of the multi-language notification message template"
  value       = microsoft365_graph_beta_device_management_notification_message_template.multilingual.id
}

output "advanced_template_id" {
  description = "ID of the advanced notification message template"
  value       = microsoft365_graph_beta_device_management_notification_message_template.advanced.id
}