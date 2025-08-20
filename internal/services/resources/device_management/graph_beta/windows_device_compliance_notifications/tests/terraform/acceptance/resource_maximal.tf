resource "random_integer" "acc_test_suffix" {
  min = 1000
  max = 9999
}

resource "microsoft365_graph_beta_device_management_windows_device_compliance_notifications" "maximal" {
  display_name     = "Acc Test Maximal - ${random_integer.acc_test_suffix.result}"
  branding_options = ["includeCompanyLogo", "includeCompanyName", "includeContactInformation", "includeCompanyPortalLink", "includeDeviceDetails"]

  role_scope_tag_ids = [microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Compliance Issue Detected"
      message_template = <<-EOT
        Dear {UserName},

        Your device '{DeviceName}' has been found to be non-compliant with company policies. 
        Please take action to resolve the following issues:

        {ComplianceReasons}

        For assistance, please contact IT support.

        Thank you,
        IT Security Team
      EOT
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
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "5m"
  }
}