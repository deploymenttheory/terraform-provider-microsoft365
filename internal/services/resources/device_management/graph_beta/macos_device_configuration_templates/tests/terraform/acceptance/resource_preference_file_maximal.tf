resource "random_string" "preference_file_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "preference_file_group_1" {
  display_name     = "acc-test-macos-preference-file-group-1-${random_string.preference_file_suffix.result}"
  mail_nickname    = "acc-test-preference-file-1-${random_string.preference_file_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "preference_file_group_2" {
  display_name     = "acc-test-macos-preference-file-group-2-${random_string.preference_file_suffix.result}"
  mail_nickname    = "acc-test-preference-file-2-${random_string.preference_file_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "preference_file_group_3" {
  display_name     = "acc-test-macos-preference-file-group-3-${random_string.preference_file_suffix.result}"
  mail_nickname    = "acc-test-preference-file-3-${random_string.preference_file_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "preference_file_group_4" {
  display_name     = "acc-test-macos-preference-file-group-4-${random_string.preference_file_suffix.result}"
  mail_nickname    = "acc-test-preference-file-4-${random_string.preference_file_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "preference_file_example" {
  display_name = "acc-test-macOS-preference-file-${random_string.preference_file_suffix.result}"
  description  = "Configure Safari browser settings via preference file"

  preference_file = {
    file_name         = "com.apple.Safari.plist"
    bundle_id         = "com.apple.Safari"
    configuration_xml = <<-EOT
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
      <plist version="1.0">
      <dict>
          <key>HomePage</key>
          <string>https://www.example.com</string>
          <key>AutoOpenSafeDownloads</key>
          <false/>
          <key>DefaultBrowserPromptingState</key>
          <integer>2</integer>
      </dict>
      </plist>
    EOT
  }

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.preference_file_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.preference_file_group_2.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.preference_file_group_3.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.preference_file_group_4.id
    }
  ]

  depends_on = [
    microsoft365_graph_beta_groups_group.preference_file_group_1,
    microsoft365_graph_beta_groups_group.preference_file_group_2,
    microsoft365_graph_beta_groups_group.preference_file_group_3,
    microsoft365_graph_beta_groups_group.preference_file_group_4
  ]

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}
