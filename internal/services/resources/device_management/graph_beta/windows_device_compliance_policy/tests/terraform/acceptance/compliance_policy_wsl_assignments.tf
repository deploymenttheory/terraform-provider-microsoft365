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

# Test Group 1
resource "microsoft365_graph_beta_groups_group" "acc_test_group_1" {
  display_name     = "acc-test-group-1-${random_string.suffix.result}"
  description      = "Test group for m365 tf provider acceptance tests"
  mail_nickname    = "acc-test-1-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Test Group 2
resource "microsoft365_graph_beta_groups_group" "acc_test_group_2" {
  display_name     = "acc-test-group-2-${random_string.suffix.result}"
  description      = "Test group for m365 tf provider acceptance tests"
  mail_nickname    = "acc-test-2-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Test Group 3
resource "microsoft365_graph_beta_groups_group" "acc_test_group_3" {
  display_name     = "acc-test-group-3-${random_string.suffix.result}"
  description      = "Test group for m365 tf provider acceptance tests"
  mail_nickname    = "acc-test-3-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Test Group 4
resource "microsoft365_graph_beta_groups_group" "acc_test_group_4" {
  display_name     = "acc-test-group-4-${random_string.suffix.result}"
  description      = "Test group for m365 tf provider acceptance tests"
  mail_nickname    = "acc-test-4-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  visibility       = "Private"
  hard_delete      = true

  timeouts = {
    create = "60s"
    read   = "30s"
    update = "30s"
    delete = "60s"
  }
}

# Test Group 5 - Microsoft 365 Group - mail-enabled (for notifications)
resource "microsoft365_graph_beta_groups_group" "acc_test_group_5" {
  display_name     = "acc-test-group-5-mail-enabled-${random_string.suffix.result}"
  description      = "Test group for m365 tf provider acceptance tests"
  mail_nickname    = "acc-test-5-${random_string.suffix.result}"
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
  display_name     = "acc-test-dcnt-wsl-assignments-${random_string.suffix.result}"
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
    microsoft365_graph_beta_groups_group.acc_test_group_1,
    microsoft365_graph_beta_groups_group.acc_test_group_2,
    microsoft365_graph_beta_groups_group.acc_test_group_3,
    microsoft365_graph_beta_groups_group.acc_test_group_4,
    microsoft365_graph_beta_groups_group.acc_test_group_5,
    microsoft365_graph_beta_device_management_device_compliance_notification_template.acc_test_device_compliance_notification_template
  ]
}

# ==============================================================================
# Windows Device Compliance Policy
# ==============================================================================

resource "microsoft365_graph_beta_device_management_windows_device_compliance_policy" "wsl_assignments" {
  display_name       = "acc-test-wdcp-wsl-assignments-${random_string.suffix.result}"
  description        = "acc-test-wdcp-wsl-assignments-${random_string.suffix.result}"
  
  depends_on = [time_sleep.wait_for_dependencies]
  role_scope_tag_ids = ["0"]

  wsl_distributions = [
    {
      distribution       = "Ubuntu"
      minimum_os_version = "1.0"
      maximum_os_version = "1.0"
    },
    {
      distribution       = "redhat"
      minimum_os_version = "1.0"
      maximum_os_version = "1.0"
    }
  ]

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
          notification_message_cc_list = [microsoft365_graph_beta_groups_group.acc_test_group_5.id]
        },
        {
          action_type        = "retire"
          grace_period_hours = 48
        },
      ]
    }
  ]

  # Assignments
  assignments = [
    # Optional: Assignment targeting all devices without filter
    {
      type        = "allDevicesAssignmentTarget"
      filter_type = "none"
    },
    # Optional: Assignment targeting all licensed users without filter
    {
      type        = "allLicensedUsersAssignmentTarget"
      filter_type = "none"
    },
    # Optional: Assignment targeting a specific group without filter
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_1.id
      filter_type = "none"
    },
    # Optional: Assignment targeting a specific group without filter
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_2.id
      filter_type = "none"
    },
    # Optional: Exclusion group assignments
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_4.id
    },
  ]

}
