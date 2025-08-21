resource "random_integer" "acc_test_suffix" {
  min = 1000
  max = 9999
}

resource "microsoft365_graph_beta_device_management_device_compliance_notification_templates" "minimal" {
  display_name     = "Acc Test Minimal - ${random_integer.acc_test_suffix.result}"
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