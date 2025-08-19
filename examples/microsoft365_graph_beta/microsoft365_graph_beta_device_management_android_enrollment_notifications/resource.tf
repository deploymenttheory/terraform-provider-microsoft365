resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "email_minimal" {
  display_name     = "email minimal"
  description      = "minimal configuration for email"
  platform_type    = "androidForWork" // "androidForWork" , "android"
  default_locale   = "en-US"
  branding_options = ["none"]

  notification_templates = ["email"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "email_maximal" {
  display_name   = "email maximal"
  description    = "Complete configuration withall features"
  platform_type  = "androidForWork" // "androidForWork" , "android"
  default_locale = "en-US"
  branding_options = ["includeCompanyLogo",
    "includeCompanyName",
    "includeCompanyPortalLink",
    "includeContactInformation",
    "includeDeviceDetails"
  ]

  notification_templates = ["email"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "push_maximal" {
  display_name           = "push maximal"
  description            = "Complete push configuration"
  platform_type          = "androidForWork" // "androidForWork" , "android"
  default_locale         = "en-US"
  branding_options       = ["none"] // no branding options for push
  notification_templates = ["push"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "push"
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

resource "microsoft365_graph_beta_device_management_android_enrollment_notifications" "all" {
  display_name   = "configuration with all enrollment notification features"
  description    = "Complete configuration with all features"
  platform_type  = "androidForWork" // "androidForWork" , "android"
  default_locale = "en-US"
  branding_options = ["includeCompanyLogo",
    "includeCompanyName",
    "includeCompanyPortalLink",
    "includeContactInformation",
    "includeDeviceDetails"
  ]

  notification_templates = ["email", "push"]

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "email"
    },
    {
      locale           = "en-us"
      subject          = "Device Enrollment Required"
      message_template = "Please enroll your device into Intune using the Company Portal to access corporate resources."
      is_default       = true
      template_type    = "push"
    }
  ]

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_1.id
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]

  role_scope_tag_ids = ["0", "1"]

  timeouts = {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}