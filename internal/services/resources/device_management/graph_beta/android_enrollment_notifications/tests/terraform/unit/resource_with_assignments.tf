resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "with_assignments" {
  display_name  = "Unit Test - Android Enrollment Notifications with Assignments"
  description   = "Configuration for unit testing with assignments"
  platform_type = "androidForWork"
  default_locale = "en-US"

  notification_templates = ["email", "push"]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "12345678-1234-1234-1234-123456789abc"
    },
    {
      type     = "allLicensedUsersAssignmentTarget"
    }
  ]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
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