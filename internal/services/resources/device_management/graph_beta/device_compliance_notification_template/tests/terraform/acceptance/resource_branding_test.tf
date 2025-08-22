resource "random_integer" "acc_test_suffix" {
  min = 1000
  max = 9999
}

resource "microsoft365_graph_beta_device_management_device_compliance_notification_template" "branding_test" {
  display_name     = "Acc Test Branding - ${random_integer.acc_test_suffix.result}"
  branding_options = ["includeCompanyLogo", "includeCompanyName", "includeContactInformation"]

  role_scope_tag_ids = ["0"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Compliance Notification"
      message_template = "Your device requires attention to maintain compliance with company policies."
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