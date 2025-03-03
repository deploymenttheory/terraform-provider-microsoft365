resource "microsoft365_graph_beta_device_and_app_management_macos_pkg_app" "mozilla_firefox" {
  display_name            = "Firefox"
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
    //icon_file_path_source = "C:\\your\\filepath\\Firefox_logo_2019.png"
    // or
    icon_url_source = "https://upload.wikimedia.org/wikipedia/commons/1/16/Firefox_logo%2C_2017.png"
  }

  //categories = ["Productivity",  "Business"]

  macos_pkg_app = {
    //installer_file_path_source = "C:\\your\\filepath\\GoogleChrome.pkg"
    // or
    installer_url_source = "https://ftp.mozilla.org/pub/firefox/releases/136.0/mac/en-GB/Firefox%20136.0.pkg"
    ignore_version_detection      = true

    included_apps = [{
      bundle_id      = "org.mozilla.firefox"
      bundle_version = "136.0"
    }]

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