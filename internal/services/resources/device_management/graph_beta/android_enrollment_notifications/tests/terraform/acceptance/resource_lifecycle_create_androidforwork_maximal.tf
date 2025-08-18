
resource "random_uuid" "lifecycle" {}

resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "lifecycle" {
  display_name     = "Acceptance - Android Enrollment Notifications (AndroidForWork Maximal) - ${random_uuid.lifecycle.result}"
  description      = "Maximal configuration for AndroidForWork platform type acceptance testing"
  platform_type    = "androidForWork"
  default_locale   = "en-US"
  branding_options = "includeCompanyLogo,includeCompanyName,includeContactInformation"

  notification_templates = ["email", "push"]

  localized_notification_messages = [
    {
      locale           = "en-US"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your AndroidForWork device to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
    {
      locale           = "en-US"
      subject          = "Device Enrollment"
      message_template = "Enroll your AndroidForWork device now"
      is_default       = true
      template_type    = "push"
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = azuread_group.acc_test_group_1.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = azuread_group.acc_test_group_2.id
    }
  ]

  role_scope_tag_ids = [
    microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id,
    microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_2.id
  ]

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }

  lifecycle {
    ignore_changes = [
      created_date_time,
      last_modified_date_time,
      version
    ]
  }
}