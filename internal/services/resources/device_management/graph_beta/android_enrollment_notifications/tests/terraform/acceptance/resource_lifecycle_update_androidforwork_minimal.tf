resource "random_uuid" "lifecycle" {}

resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "lifecycle" {
  display_name   = "Acceptance - Android Enrollment Notifications (AndroidForWork Minimal Updated) - ${random_uuid.lifecycle.result}"
  description    = "Updated minimal configuration for AndroidForWork platform type acceptance testing"
  platform_type  = "androidForWork"
  default_locale = "en-US"

  notification_templates = ["email", "push"]

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }

  lifecycle {
    ignore_changes = [
      role_scope_tag_ids,
      created_date_time,
      last_modified_date_time,
      version
    ]
  }
}