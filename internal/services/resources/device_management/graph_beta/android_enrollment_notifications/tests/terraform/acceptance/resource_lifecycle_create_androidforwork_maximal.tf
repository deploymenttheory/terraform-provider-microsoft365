
resource "random_uuid" "lifecycle" {}

# Create test groups for assignments
resource "random_uuid" "group_1" {}
resource "random_uuid" "group_2" {}

resource "azuread_group" "acc_test_group_1" {
  display_name     = "Test Group 1 - ${random_uuid.group_1.result}"
  security_enabled = true
}

resource "azuread_group" "acc_test_group_2" {
  display_name     = "Test Group 2 - ${random_uuid.group_2.result}" 
  security_enabled = true
}

resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "lifecycle" {
  display_name     = "Acceptance - Android Enrollment Notifications (AndroidForWork Maximal) - ${random_uuid.lifecycle.result}"
  description      = "Maximal configuration for AndroidForWork platform type acceptance testing"
  platform_type    = "androidForWork"
  default_locale   = "en-US"
  branding_options = ["includeCompanyLogo", "includeCompanyName", "includeContactInformation"]

  notification_templates = ["email", "push"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your AndroidForWork device to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
    {
      locale           = "en-us"
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

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }

}