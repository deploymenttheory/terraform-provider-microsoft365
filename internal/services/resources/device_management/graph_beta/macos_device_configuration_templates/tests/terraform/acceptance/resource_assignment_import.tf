resource "random_string" "assignment_import" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "assignment_import" {
  display_name     = "acc-test-4873-macos-group-${random_string.assignment_import.result}"
  mail_nickname    = "acc-test-4873-macos-${random_string.assignment_import.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}
# Example 1: macOS Custom Configuration Template
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "test" {
  display_name = "acc-test-4873-macos-${random_string.assignment_import.result}"
  description  = "Example custom configuration template for macOS devices"

  custom_configuration = {
    deployment_channel = "deviceChannel"
    payload_file_name  = "com.example.custom.mobileconfig"
    payload_name       = "Custom Configuration Example"
    payload            = <<-EOT
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
      <plist version="1.0">
      <dict>
          <key>PayloadContent</key>
          <array>
              <dict>
                  <key>PayloadDisplayName</key>
                  <string>Custom Example Configuration</string>
                  <key>PayloadIdentifier</key>
                  <string>com.example.custom.settings</string>
                  <key>PayloadType</key>
                  <string>com.example.custom</string>
                  <key>PayloadUUID</key>
                  <string>12345678-1234-1234-1234-123456789012</string>
                  <key>PayloadVersion</key>
                  <integer>1</integer>
                  <key>ExampleSetting</key>
                  <true/>
              </dict>
          </array>
          <key>PayloadDisplayName</key>
          <string>Custom Configuration Example</string>
          <key>PayloadIdentifier</key>
          <string>com.example.custom</string>
          <key>PayloadType</key>
          <string>Configuration</string>
          <key>PayloadUUID</key>
          <string>87654321-4321-4321-4321-210987654321</string>
          <key>PayloadVersion</key>
          <integer>1</integer>
      </dict>
      </plist>
    EOT

  }

  role_scope_tag_ids = ["0"]
  assignments = [{
    type     = "groupAssignmentTarget"
    group_id = microsoft365_graph_beta_groups_group.assignment_import.id
  }]
}
