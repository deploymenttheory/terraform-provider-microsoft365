resource "random_uuid" "lifecycle" {}

resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "lifecycle" {
  display_name   = "Acceptance - Android Enrollment Notifications (Android Minimal) - ${random_uuid.lifecycle.result}"
  description    = "Minimal configuration for Android platform type acceptance testing"
  platform_type  = "android"
  default_locale = "en-US"

  notification_templates = ["email"]

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }

}