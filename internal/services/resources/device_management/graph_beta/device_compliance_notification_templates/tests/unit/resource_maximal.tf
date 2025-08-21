# Example: Multilingual Notification Message Template
resource "microsoft365_graph_beta_device_management_device_compliance_notification_templates" "multilingual" {
  display_name     = "Multilingual Compliance Notification"
  branding_options = ["includeCompanyLogo", "includeCompanyName", "includeContactInformation"]

  role_scope_tag_ids = ["0"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Compliance Issue Detected"
      message_template = "Hello {UserName},\n\nYour device '{DeviceName}' has been found to be non-compliant with company policies. Please take action to resolve the following issues:\n\n{ComplianceReasons}\n\nFor assistance, please contact IT support.\n\nThank you,\nIT Security Team"
      is_default       = true
    },
    {
      locale           = "es-es"
      subject          = "Problema de Cumplimiento del Dispositivo"
      message_template = "Hola {UserName},\n\nTu dispositivo '{DeviceName}' no cumple las normas. Por favor resuelve: {ComplianceReasons}\n\nContacta con IT para ayuda.\n\nEquipo de Seguridad IT"
      is_default       = false
    },
    {
      locale           = "fr-fr"
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