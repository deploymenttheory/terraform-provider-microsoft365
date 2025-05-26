# Example 1: Platform Restrictions Configuration
resource "microsoft365_graph_beta_device_management_device_enrollment_configuration" "platform_restrictions" {
  display_name                           = "Corporate Platform Restrictions"
  description                           = "Platform restrictions for corporate device enrollment"
  device_enrollment_configuration_type  = "platformRestrictions"
  priority                             = 10

  platform_restriction {
    platform_type = "ios"
    
    restriction {
      platform_blocked                     = false
      personal_device_enrollment_blocked   = true
      os_minimum_version                   = "15.0"
      os_maximum_version                   = "17.0"
      blocked_manufacturers                = []
      blocked_skus                        = []
    }
  }

  assignment {
    target_group_id = "12345678-1234-1234-1234-123456789012"
  }

  timeouts {
    create = "10m"
    read   = "5m"
    update = "10m"
    delete = "5m"
  }
}

# Example 2: Windows Hello for Business Configuration
resource "microsoft365_graph_beta_device_management_device_enrollment_configuration" "windows_hello" {
  display_name                           = "Corporate Windows Hello for Business"
  description                           = "Windows Hello for Business settings for corporate devices"
  device_enrollment_configuration_type  = "windowsHelloForBusiness"
  priority                             = 5

  windows_hello_for_business {
    state                               = "enabled"
    pin_minimum_length                  = 6
    pin_maximum_length                  = 12
    pin_uppercase_characters_usage      = "allowed"
    pin_lowercase_characters_usage      = "required"
    pin_special_characters_usage        = "allowed"
    security_device_required            = true
    unlock_with_biometrics_enabled      = true
    remote_passport_enabled             = false
    pin_previous_block_count            = 5
    pin_expiration_in_days              = 365
    enhanced_biometrics_state           = "enabled"
    security_key_for_sign_in           = "enabled"
    enhanced_sign_in_security          = 1
  }

  assignment {
    target_group_id = "87654321-4321-4321-4321-210987654321"
  }
}

# Example 3: Device Enrollment Limit Configuration
resource "microsoft365_graph_beta_device_management_device_enrollment_configuration" "enrollment_limit" {
  display_name                           = "User Device Limit"
  description                           = "Limit the number of devices users can enroll"
  device_enrollment_configuration_type  = "limit"
  priority                             = 15

  device_enrollment_limit {
    limit = 5
  }

  assignment {
    target_group_id = "11111111-2222-3333-4444-555555555555"
  }
}

# Example 4: Windows 10 Enrollment Completion Page Configuration
resource "microsoft365_graph_beta_device_management_device_enrollment_configuration" "enrollment_status_page" {
  display_name                           = "Corporate Enrollment Status Page"
  description                           = "Custom enrollment status page for Windows 10 devices"
  device_enrollment_configuration_type  = "windows10EnrollmentCompletionPageConfiguration"
  priority                             = 20

  windows10_enrollment_completion_page {
    show_installation_progress              = true
    block_device_setup_retry_by_user       = false
    allow_device_reset_on_install_failure  = true
    allow_log_collection_on_install_failure = true
    custom_error_message                   = "If you encounter issues during setup, please contact IT support at it-help@company.com"
    install_progress_timeout_in_minutes    = 90
    allow_device_use_on_install_failure    = false
    selected_mobile_app_ids                = [
      "app-id-1",
      "app-id-2"
    ]
    allow_non_blocking_app_installation           = true
    install_quality_updates                       = true
    track_install_progress_for_autopilot_only    = false
    disable_user_status_tracking_after_first_user = false
  }

  assignment {
    target_group_id = "99999999-8888-7777-6666-555555555555"
  }
}

# Example 5: Enrollment Notifications Configuration
resource "microsoft365_graph_beta_device_management_device_enrollment_configuration" "enrollment_notifications" {
  display_name                           = "Enrollment Email Notifications"
  description                           = "Email notifications for device enrollment"
  device_enrollment_configuration_type  = "enrollmentNotificationsConfiguration"
  priority                             = 25

  enrollment_notifications {
    platform_type                      = "allPlatforms"
    template_type                      = "email"
    notification_message_template_id   = "template-id-12345"
    notification_templates             = ["template1", "template2"]
    branding_options                   = [
      "includeCompanyLogo",
      "includeCompanyName",
      "includeContactInformation"
    ]
    default_locale                     = "en-US"
    include_company_portal_link        = true
    send_push_notification             = false
    notification_title                 = "Welcome to Company Devices"
    notification_body                  = "Your device has been successfully enrolled in our mobile device management system."
    notification_sender                = "IT Support"
  }

  assignment {
    target_group_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
  }
}

# Example 6: Device Co-management Authority Configuration
resource "microsoft365_graph_beta_device_management_device_enrollment_configuration" "comanagement" {
  display_name                           = "ConfigMgr Co-management"
  description                           = "Co-management settings for Configuration Manager integration"
  device_enrollment_configuration_type  = "deviceComanagementAuthorityConfiguration"
  priority                             = 30

  device_comanagement_authority {
    managed_device_authority                           = 1
    install_configuration_manager_agent               = true
    configuration_manager_agent_command_line_argument = "/mp:sccm.company.com /sitecode:ABC"
  }

  assignment {
    target_group_id = "cccccccc-dddd-eeee-ffff-000000000000"
  }
}

# Example 7: Multiple Platform Restrictions
resource "microsoft365_graph_beta_device_management_device_enrollment_configuration" "android_restrictions" {
  display_name                           = "Android Device Restrictions"
  description                           = "Restrictions for Android device enrollment"
  device_enrollment_configuration_type  = "singlePlatformRestriction"
  priority                             = 12

  platform_restriction {
    platform_type = "android"
    
    restriction {
      platform_blocked                     = false
      personal_device_enrollment_blocked   = true
      os_minimum_version                   = "10.0"
      os_maximum_version                   = ""
      blocked_manufacturers                = ["manufacturer1", "manufacturer2"]
      blocked_skus                        = ["sku1", "sku2", "sku3"]
    }
  }

  assignment {
    target_group_id = "dddddddd-eeee-ffff-1111-222222222222"
  }
}

# Example 8: Default Windows 10 Enrollment Completion Page
resource "microsoft365_graph_beta_device_management_device_enrollment_configuration" "default_esp" {
  display_name                           = "Default ESP Settings"
  description                           = "Default enrollment status page configuration"
  device_enrollment_configuration_type  = "defaultWindows10EnrollmentCompletionPageConfiguration"
  priority                             = 1

  default_windows10_enrollment_completion_page {
    allow_devices_for_users                        = true
    show_installation_progress                     = true
    allow_device_reset_on_install_failure         = false
    allow_log_collection_on_install_failure       = true
    custom_error_message                          = "Please contact your administrator for assistance."
    install_progress_timeout_in_minutes           = 60
    selected_mobile_app_ids                       = []
    track_install_progress_for_autopilot_only     = true
    disable_user_status_tracking_after_first_user = true
  }

  assignment {
    target_group_id = "ffffffff-1111-2222-3333-444444444444"
  }
}

# Example 9: Comprehensive Configuration with Multiple Assignments
resource "microsoft365_graph_beta_device_management_device_enrollment_configuration" "comprehensive_config" {
  display_name                           = "Executive Device Configuration"
  description                           = "Comprehensive device enrollment configuration for executives"
  device_enrollment_configuration_type  = "platformRestrictions"
  priority                             = 1
  role_scope_tag_ids                    = ["0", "executive-scope-tag"]

  platform_restriction {
    platform_type = "ios"
    
    restriction {
      platform_blocked                     = false
      personal_device_enrollment_blocked   = false
      os_minimum_version                   = "16.0"
      os_maximum_version                   = ""
      blocked_manufacturers                = []
      blocked_skus                        = []
    }
  }

  # Multiple assignments to different groups
  assignment {
    target_group_id = "executives-group-id"
  }

  assignment {
    target_group_id = "senior-management-group-id"
  }

  timeouts {
    create = "15m"
    read   = "10m"
    update = "15m"
    delete = "10m"
  }
}