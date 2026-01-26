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
# Device Compliance Notification Template Dependency
# ==============================================================================

resource "microsoft365_graph_beta_device_management_device_compliance_notification_template" "acc_test_device_compliance_notification_template" {
  display_name     = "acc-test-dcnt-system-security-${random_string.suffix.result}"
  branding_options = ["includeCompanyLogo"]

  role_scope_tag_ids = ["0"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Compliance Required"
      message_template = "Please ensure your device meets the compliance requirements to access corporate resources."
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

# ==============================================================================
# Time Sleep for Eventual Consistency
# ==============================================================================

resource "time_sleep" "wait_for_dependencies" {
  create_duration = "30s"

  depends_on = [
    microsoft365_graph_beta_groups_group.acc_test_group,
    microsoft365_graph_beta_device_management_device_compliance_notification_template.acc_test_device_compliance_notification_template
  ]
}

# ==============================================================================
# Windows Device Compliance Policy
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "system_security" {
  display_name = "acc-test-wdcp-system-security-${random_string.suffix.result}"

  depends_on         = [time_sleep.wait_for_dependencies]
  description        = "acc-test-wdcp-system-security-${random_string.suffix.result}"
  role_scope_tag_ids = ["0"]

  # System Security Settings
  system_security = {
    active_firewall_required                         = true
    anti_spyware_required                            = true
    antivirus_required                               = true
    configuration_manager_compliance_required        = true
    defender_enabled                                 = true
    defender_version                                 = "1.0.0.0"
    device_threat_protection_enabled                 = true
    device_threat_protection_required_security_level = "medium"
    password_block_simple                            = true
    password_minimum_character_set_count             = 4
    password_required                                = true
    password_required_to_unlock_from_idle            = true
    password_required_type                           = "alphanumeric"
    rtp_enabled                                      = true
    signature_out_of_date                            = true
    storage_require_encryption                       = true
    tpm_required                                     = true
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
