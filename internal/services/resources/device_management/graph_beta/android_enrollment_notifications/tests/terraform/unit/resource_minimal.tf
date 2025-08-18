resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "unit" {
  display_name  = "Unit Test - Android Enrollment Notifications"
  description   = "Minimal configuration for unit testing"
  platform_type = "androidForWork"
  default_locale = "en-US"

  notification_templates = ["email", "push"]

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}