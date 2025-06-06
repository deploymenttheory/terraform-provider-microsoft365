resource "microsoft365_graph_beta_device_management_device_enrollment_notification_configuration" "example" {
  display_name   = "Test Notification Configuration"
  description    = "Test notification configuration for device enrollment"
  template_types = ["push", "email"]

  branding_options = [
    "includeCompanyLogo",
    "includeCompanyName",
    "includeContactInformation",
    "includeCompanyPortalLink",
    "includeDeviceDetails"
  ]

  push_localized_message = {
    locale           = "en-us"
    subject          = "Device Enrolled Successfully"
    message_template = "Your device has been successfully enrolled in the organization!"
    is_default       = true
  }

  email_localized_message = {
    locale           = "en-us"
    subject          = "Device Enrollment Complete"
    message_template = "Your device enrollment process is now complete. You can start using your device."
    is_default       = true
  }

  assignments = [
    {
      target = {
        target_type = "group"
        group_id    = "b8c661c2-fa9a-4351-af86-adc1729c343f"
      }
    },
    {
      target = {
        target_type = "group"
        group_id    = "51a96cdd-4b9b-4849-b416-8c94a6d88797"
      }
    },

    # Or assign to all users
    # {
    #   target = {
    #     target_type = "allLicensedUsers"
    #   }
    # }
  ]
}