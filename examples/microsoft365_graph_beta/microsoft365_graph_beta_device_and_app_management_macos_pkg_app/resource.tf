# Example: macOS PKG App Resource
resource "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "google_chrome" {
  display_name            = "GoogleChrome.pkg"
  description             = "thing"
  publisher               = "Example Publisher"
  is_featured             = false
  privacy_information_url = "https://example.com/privacy"
  information_url         = "https://example.com/info"
  owner                   = "Example Owner"
  developer               = "Example Developer"
  notes                   = "This is a macOS PKG application managed through Terraform."
  role_scope_tag_ids      = [8, 9]

  app_icon = {
    icon_file_path = "C:\\your\\localpath\\chrome_logo.png"
  }

  categories = [
    {
      display_name = "Productivity"
    },
    {
      display_name = "Business"
    }
  ]

  macos_pkg_app = {
    package_installer_file_source = "C:\\your\\localpath\\GoogleChrome.pkg"
    ignore_version_detection      = true

    minimum_supported_operating_system = {
      v14_0 = true
    }

    pre_install_script = {
      script_content = base64encode("#!/bin/bash\necho macOS PKG Pre-install script example")
    }

    post_install_script = {
      script_content = base64encode("#!/bin/bash\necho macOS PKG Post-install script example")
    }
  }

  # App assignments configuration
  assignments = [
    
    # Assignment 1: Exclusion group with available intent
    {
      intent = "available"
      source = "direct"
      target = {
        target_type = "exclusionGroupAssignment"
        group_id    = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
      }
    },
    
    # Assignment 2: Another exclusion group with available intent
    {
      intent = "available"
      source = "direct"
      target = {
        target_type = "exclusionGroupAssignment"
        group_id    = "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
      }
    },

    # Assignment 3: All devices with required intent
    {
      intent = "required"
      source = "direct"
      target = {
        target_type = "allDevices"
      }
    },

    # Assignment 4: All licensed users with required intent
    {
      intent = "required"
      source = "direct"
      target = {
        target_type = "allLicensedUsers"
      }
    },

    # Assignment 5: Group assignment with required intent
    {
      intent = "required"
      source = "direct"
      target = {
        target_type = "groupAssignment"
        group_id    = "b15228f4-9d49-41ed-9b4f-0e7c721fd9c2"
      }
    },
    
    # Assignment 6: Another group assignment with required intent
    {
      intent = "required"
      source = "direct"
      target = {
        target_type = "groupAssignment"
        group_id    = "ea8e2fb8-e909-44e6-bae7-56757cf6f347"
      }
    }
  ]

  # Optional: Add timeouts
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "30m"
  }
}