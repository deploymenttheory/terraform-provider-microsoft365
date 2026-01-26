# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

# Test Group 5 - Microsoft 365 Group - mail-enabled (for notifications)
resource "microsoft365_graph_beta_groups_group" "acc_test_group" {
  display_name     = "acc-test-group-mail-enabled-${random_string.suffix.result}"
  description      = "Test group for m365 tf provider acceptance tests"
  mail_nickname    = "acc-test-${random_string.suffix.result}"
  mail_enabled     = true
  security_enabled = false
  group_types      = ["Unified"]
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# ==============================================================================
# Windows Device Compliance Script Dependency
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_device_compliance_script" "acc_test_windows_device_compliance_script" {
  display_name             = "acc-test-windows-device-compliance-script-${random_string.suffix.result}"
  description              = "acc-test-windows-device-compliance-script"
  publisher                = "Acceptance Test Publisher"
  detection_script_content = "Get-Process | Select-Object -First 10"
  run_as_account           = "system"
  enforce_signature_check  = false
  run_as_32_bit            = false

  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}

# ==============================================================================
# Device Compliance Notification Template Dependency
# ==============================================================================

resource "microsoft365_graph_beta_device_management_device_compliance_notification_template" "acc_test_device_compliance_notification_template" {
  display_name     = "acc-test-dcnt-custom-compliance-${random_string.suffix.result}"
  branding_options = ["includeCompanyLogo", "includeCompanyName", "includeContactInformation", "includeCompanyPortalLink", "includeDeviceDetails"]

  role_scope_tag_ids = [0]

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

# ==============================================================================
# Time Sleep for Eventual Consistency
# ==============================================================================

resource "time_sleep" "wait_for_dependencies" {
  create_duration = "30s"

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group,
    microsoft365_graph_beta_device_management_windows_device_compliance_script.acc_test_windows_device_compliance_script,
    microsoft365_graph_beta_device_management_device_compliance_notification_template.acc_test_device_compliance_notification_template
  ]
}

# ==============================================================================
# Windows Device Compliance Policy
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "custom_compliance" {

  depends_on = [time_sleep.wait_for_dependencies]

  display_name       = "acc-test-wdcp-custom-compliance-${random_string.suffix.result}"
  description        = "acc-test-wdcp-custom-compliance-${random_string.suffix.result}"
  role_scope_tag_ids = ["0"]

  # Custom compliance script
  custom_compliance_required = true
  device_compliance_policy_script = {
    device_compliance_script_id = microsoft365_graph_beta_device_management_windows_device_compliance_script.acc_test_windows_device_compliance_script.id
    rules_content = jsonencode({
      "Rules" : [
        {
          "SettingName" : "BiosVersion",
          "Operator" : "GreaterEquals",
          "DataType" : "Version",
          "Operand" : "2.3",
          "MoreInfoUrl" : "https://bing.com",
          "RemediationStrings" : [
            {
              "Language" : "en_US",
              "Title" : "BIOS Version needs to be upgraded to at least 2.3. Value discovered was {ActualValue}.",
              "Description" : "BIOS must be updated. Please refer to the link above"
            },
            {
              "Language" : "de_DE",
              "Title" : "BIOS-Version muss auf mindestens 2.3 aktualisiert werden. Der erkannte Wert lautet {ActualValue}.",
              "Description" : "BIOS muss aktualisiert werden. Bitte beziehen Sie sich auf den obigen Link"
            }
          ]
        },
        {
          "SettingName" : "TPMChipPresent",
          "Operator" : "IsEquals",
          "DataType" : "Boolean",
          "Operand" : true,
          "MoreInfoUrl" : "https://bing.com",
          "RemediationStrings" : [
            {
              "Language" : "en_US",
              "Title" : "TPM chip must be enabled.",
              "Description" : "TPM chip must be enabled. Please refer to the link above"
            },
            {
              "Language" : "de_DE",
              "Title" : "TPM-Chip muss aktiviert sein.",
              "Description" : "TPM-Chip muss aktiviert sein. Bitte beziehen Sie sich auf den obigen Link"
            }
          ]
        },
        {
          "SettingName" : "Manufacturer",
          "Operator" : "IsEquals",
          "DataType" : "String",
          "Operand" : "Microsoft Corporation",
          "MoreInfoUrl" : "https://bing.com",
          "RemediationStrings" : [
            {
              "Language" : "en_US",
              "Title" : "Only Microsoft devices are supported.",
              "Description" : "You are not currently using a Microsoft device."
            },
            {
              "Language" : "de_DE",
              "Title" : "Nur Microsoft-Geräte werden unterstützt.",
              "Description" : "Sie verwenden derzeit kein Microsoft-Gerät."
            }
          ]
        }
      ]
    })
  }


  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 12
        },
        {
          action_type                  = "notification"
          grace_period_hours           = 24
          notification_template_id     = microsoft365_graph_beta_device_management_device_compliance_notification_template.acc_test_device_compliance_notification_template.id
          notification_message_cc_list = [microsoft365_graph_beta_groups_group.acc_test_group.id]
        },
        {
          action_type        = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

}