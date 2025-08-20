
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "preference_file_example" {
  display_name = "unit-test-macOS-preference-file-example"
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

  role_scope_tag_ids = ["00000000-0000-0000-0000-000000000001"]

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000002"
      filter_id   = "00000000-0000-0000-0000-000000000003"
      filter_type = "include"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000002"
      filter_id   = "00000000-0000-0000-0000-000000000003"
      filter_type = "exclude"
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000002"
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000004"
    }
  ]

  timeouts = {
    create = "30m"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}