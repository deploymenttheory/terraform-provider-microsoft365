# Basic device enrollment notification configuration with email only
resource "microsoft365_graph_beta_device_management_device_enrollment_notification_configuration" "email_only" {
  display_name           = "Email Enrollment Notification"
  description            = "Email notification for device enrollment"
  template_types         = ["email"]
  priority               = 1
  default_locale         = "en-US"
  
  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Push notification only configuration
resource "microsoft365_graph_beta_device_management_device_enrollment_notification_configuration" "push_only" {
  display_name           = "Push Enrollment Notification"
  description            = "Push notification for device enrollment"
  template_types         = ["push"]
  priority               = 2
  default_locale         = "en-US"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Combined email and push notification configuration
resource "microsoft365_graph_beta_device_management_device_enrollment_notification_configuration" "combined" {
  display_name           = "Combined Enrollment Notifications"
  description            = "Both email and push notifications for device enrollment"
  template_types         = ["email", "push"]
  priority               = 3
  default_locale         = "en-US"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
} 