resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "basic" {
  display_name = "Windows 10/11 - Basic Compliance Policy"
  description  = "Basic Windows device compliance policy requiring BitLocker, Secure Boot and a password"

  # Password requirements
  password_required                          = true
  password_block_simple                      = true
  password_required_to_unlock_from_idle      = true
  password_minimum_length                    = 8
  password_minimum_character_set_count       = 3
  password_required_type                     = "alphanumeric"
  password_minutes_of_inactivity_before_lock = 15

  # Security requirements
  bit_locker_enabled         = true
  secure_boot_enabled        = true
  code_integrity_enabled     = true
  storage_require_encryption = true

  # Defender requirements
  defender_enabled      = true
  rtP_enabled           = true
  antivirus_required    = true
  anti_spyware_required = true

  # OS version requirements
  os_minimum_version = "10.0.19041.0"

  # Role scope tags
  role_scope_tag_ids = ["0"]

  # Non-compliance actions
  scheduled_actions_for_rule {
    rule_name = "PasswordRequired"
    scheduled_action_configurations {
      action_type        = "block"
      grace_period_hours = 24
    }
  }
}

resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "custom_compliance" {
  display_name = "Windows 10/11 - Custom Compliance Policy"
  description  = "Windows device compliance policy with custom compliance script"

  # Password requirements
  password_required                     = true
  password_block_simple                 = true
  password_required_to_unlock_from_idle = true
  password_minimum_length               = 8
  password_required_type                = "alphanumeric"

  # Security requirements
  storage_require_encryption = true
  active_firewall_required   = true
  tpm_required               = true
  antivirus_required         = true
  anti_spyware_required      = true

  # Custom compliance script
  custom_compliance_required = true
  device_compliance_policy_script {
    device_compliance_script_id = "8c3d2ec3-3e63-4df3-8265-69bbba1e53e5"
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

  # Role scope tags
  role_scope_tag_ids = ["0"]

  # Non-compliance actions
  scheduled_actions_for_rule {
    rule_name = "PasswordRequired"
    scheduled_action_configurations {
      action_type              = "block"
      grace_period_hours       = 6
      notification_template_id = ""
    }
  }
}


resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "wsl" {
  display_name = "Windows 10/11 - WSL Compliance Policy"
  description  = "Windows device compliance policy with WSL distribution requirements"

  # Password requirements
  password_required                     = true
  password_block_simple                 = true
  password_required_to_unlock_from_idle = true
  password_minimum_length               = 8
  password_minimum_character_set_count  = 3
  password_required_type                = "alphanumeric"

  # Security requirements
  bit_locker_enabled  = true
  secure_boot_enabled = true
  tpm_required        = true

  # WSL distributions
  wsl_distributions {
    distribution       = "Ubuntu"
    minimum_os_version = "20.04"
    maximum_os_version = "22.04"
  }

  wsl_distributions {
    distribution       = "Debian"
    minimum_os_version = "11.0"
    maximum_os_version = "12.0"
  }

  # Valid OS build ranges
  valid_operating_system_build_ranges {
    low_os_version  = "10.0.19041.0"
    high_os_version = "10.0.22631.3155"
  }

  # Role scope tags
  role_scope_tag_ids = ["0"]

  # Non-compliance actions
  scheduled_actions_for_rule {
    rule_name = "PasswordRequired"
    scheduled_action_configurations {
      action_type              = "block"
      grace_period_hours       = 24
      notification_template_id = ""
    }

    scheduled_action_configurations {
      action_type                  = "notification"
      grace_period_hours           = 0
      notification_template_id     = "00000000-0000-0000-0000-000000000000"
      notification_message_cc_list = []
    }
  }
} 