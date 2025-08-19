
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "preference_file_example" {
  display_name = "acc-test-macOS-preference-file-example"
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

  role_scope_tag_ids = [microsoft365_graph_beta_device_management_role_scope_tag.acc_test_role_scope_tag_1.id]

 assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_1.id
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_1.id
      filter_type = "include"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_2.id
      filter_id   = microsoft365_graph_beta_device_management_assignment_filter.acc_test_assignment_filter_2.id
      filter_type = "exclude"
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = microsoft365_graph_beta_groups_group.acc_test_group_4.id
    }
  ]
  
  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}