########################################################################################
# macOS PKG Assignment Examples
########################################################################################

# Resource for assigning a macos_pkg_app (company_portal) to all licensed users
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "company_portal_all_users" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_macos_pkg_app.company_portal.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "allLicensedUsers"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning a macos_pkg_app (company_portal) to all devices
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "company_portal_all_devices" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_macos_pkg_app.company_portal.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "allDevices"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning company_portal to a specific group with available install intent
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "company_portal_group_assignment_available" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_macos_pkg_app.company_portal.id
  intent        = "available"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "2c39cf3d-78ef-4227-acb1-3a14fc7fbb99"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning company_portal to a specific group with required install intent
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "company_portal_group_assignment_required" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_macos_pkg_app.company_portal.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "7e30b7f0-b2f1-4220-883f-f1d8066eef2d"
    device_and_app_management_assignment_filter_type = "none"
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

########################################################################################
# Win Get Assignment Examples
########################################################################################

# Resource for assigning a WinGet app (Firefox) to all licensed users
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "firefox_all_users" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win_get_app.example_firefox.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "allLicensedUsers"
    device_and_app_management_assignment_filter_type = "none"
  }

  settings = {
    win_get = {
      notifications = "showAll"
      install_time_settings = {
        use_local_time     = true
        deadline_date_time = "2025-06-01T18:00:00Z"
      }
      restart_settings = {
        grace_period_in_minutes                         = 240
        countdown_display_before_restart_in_minutes     = 30
        restart_notification_snooze_duration_in_minutes = 60
      }
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning a WinGet app (Firefox) to all devices
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "firefox_all_devices" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win_get_app.example_firefox.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "allDevices"
    device_and_app_management_assignment_filter_type = "none"
  }

  settings = {
    win_get = {
      notifications = "showAll"
      install_time_settings = {
        use_local_time     = true
        deadline_date_time = "2025-06-01T18:00:00Z"
      }
      restart_settings = {
        grace_period_in_minutes                         = 240
        countdown_display_before_restart_in_minutes     = 30
        restart_notification_snooze_duration_in_minutes = 60
      }
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning Firefox to a specific group with available install intent
# and assignment filters
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "firefox_group_assignment_available" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win_get_app.example_firefox.id
  intent        = "available"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "2c39cf3d-78ef-4227-acb1-3a14fc7fbb99"
    device_and_app_management_assignment_filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
    device_and_app_management_assignment_filter_type = "include"
  }

  settings = {
    win_get = {
      notifications = "hideAll"
      restart_settings = {
        grace_period_in_minutes                         = 240
        countdown_display_before_restart_in_minutes     = 30
        restart_notification_snooze_duration_in_minutes = 60
      }
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning Firefox to a specific group with uninstall install intent
# and assignment filters
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "firefox_group_assignment_uninstall" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win_get_app.example_firefox.id
  intent        = "uninstall"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "eadb85bd-6567-4db9-b65c-3f5070d83487"
    device_and_app_management_assignment_filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
    device_and_app_management_assignment_filter_type = "include"
  }

  settings = {
    win_get = {
      notifications = "hideAll"
      restart_settings = {
        grace_period_in_minutes                         = 240
        countdown_display_before_restart_in_minutes     = 30
        restart_notification_snooze_duration_in_minutes = 60
      }
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

# Resource for assigning Firefox to a specific group with required install intent
# and assignment filters
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "firefox_group_assignment_required" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_win_get_app.example_firefox.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "7e30b7f0-b2f1-4220-883f-f1d8066eef2d"
    device_and_app_management_assignment_filter_id   = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
    device_and_app_management_assignment_filter_type = "exclude"
  }

  settings = {
    win_get = {
      notifications = "hideAll"
      restart_settings = {
        grace_period_in_minutes                         = 240
        countdown_display_before_restart_in_minutes     = 30
        restart_notification_snooze_duration_in_minutes = 60
      }
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}

########################################################################################
# iOS Store App Assignment Examples
########################################################################################

# Resource for assigning a iOS Store app (Microsoft Edge) to a specific group with required install intent
# and assignment filters
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment" "ios_store_app_assignment" {
  mobile_app_id = microsoft365_graph_beta_device_and_app_management_ios_store_app.example.id
  intent        = "required"
  source        = "direct"

  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
    device_and_app_management_assignment_filter_id   = "471b28c1-8d90-49a2-b639-a47b5f84986d"
    device_and_app_management_assignment_filter_type = "include"
  }

  settings = {
    ios_store = {
      is_removable                = true
      prevent_managed_app_backup  = false
      uninstall_on_device_removal = true
      vpn_configuration_id        = "80f8c0a5-f3ec-4936-bcbc-420dc0ca3665"
    }
  }

  timeouts = {
    create = "30s"
    read   = "30s"
    update = "30s"
    delete = "30s"
  }
}