resource "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "mozilla_firefox" {
  display_name            = "Firefox 136.0.pkg"
  description             = "test"
  publisher               = "Example Publisher"
  is_featured             = true
  privacy_information_url = "https://example.com/privacy"
  information_url         = "https://example.com/info"
  owner                   = "Example Owner"
  developer               = "Example Developer"
  notes                   = "This is a macOS PKG application managed through Terraform."
  role_scope_tag_ids      = [8, 9]

  categories = [
    microsoft365_graph_beta_device_and_app_management_application_category.web_browser.id, // custom category
    "Business",                                                                            // built-in example
    "Productivity",
  ]

  app_icon = {
    icon_file_path_source = "/local/path/Firefox_logo.png"
  }

  app_installer = {
    installer_file_path_source = "/local/path/Firefox_136.0.pkg"
  }


  macos_pkg_app = {
    ignore_version_detection = true

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
        group_id    = "11111111-2222-3333-4444-555555555555"
      }
    },

    # Assignment 2: Another exclusion group with available intent
    {
      intent = "available"
      source = "direct"
      target = {
        target_type = "exclusionGroupAssignment"
        group_id    = "11111111-2222-3333-4444-555555555555"
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
        group_id    = "11111111-2222-3333-4444-555555555555"
      }
    },

    # Assignment 6: Another group assignment with required intent
    {
      intent = "required"
      source = "direct"
      target = {
        target_type = "groupAssignment"
        group_id    = "11111111-2222-3333-4444-555555555555"
      }
    }
  ]

  # Optional: Add timeouts
  timeouts = {
    create = "3m"
    read   = "20s"
    update = "3m"
    delete = "20s"
  }
}