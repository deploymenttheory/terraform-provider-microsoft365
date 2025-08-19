# Terraform resource configuration for Microsoft 365 Graph Beta Device Management macOS Configuration Templates

# Example 1: macOS Custom Configuration Template
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "custom_configuration_example" {
  display_name = "macos custom mobileconfig example"
  description  = "macos custom mobileconfig example"

  custom_configuration = {
    deployment_channel  = "deviceChannel"
    payload_file_name   = "com.example.custom.mobileconfig"
    payload_name        = "Custom Configuration Example"
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

  role_scope_tag_ids = ["00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"]

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000001"
      filter_id   = "00000000-0000-0000-0000-000000000002"
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
      group_id    = "00000000-0000-0000-0000-000000000002"
    }
  ]

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}

# Example 2: macOS Preference File Configuration
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "preference_file_example" {
  display_name = "macos preference file example"
  description  = "macos preference file example"

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

  role_scope_tag_ids = ["00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"]

 assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000001"
      filter_id   = "00000000-0000-0000-0000-000000000002"
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
      group_id    = "00000000-0000-0000-0000-000000000002"
    }
  ]
  
  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}

# Example 3: macOS Trusted Root Certificate
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "trusted_cert_example" {
  display_name = "macos trusted root certificate example"
  description  = "macos trusted root certificate example"

  trusted_certificate = {
    deployment_channel        = "deviceChannel"
    cert_file_name           = "MicrosoftRootCertificateAuthority2011.cer"
    trusted_root_certificate = filebase64("MicrosoftRootCertificateAuthority2011.cer")
  }

  role_scope_tag_ids = ["00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"]

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000001"
      filter_id   = "00000000-0000-0000-0000-000000000002"
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
      group_id    = "00000000-0000-0000-0000-000000000002"
    }
  ]

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}

# Example 4: macOS SCEP Certificate Profile
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "scep_cert_example" {
  display_name = "macos scep certificate example"
  description  = "macos scep certificate example"

  scep_certificate = {
    deployment_channel                  = "deviceChannel"
    renewal_threshold_percentage        = 20
    certificate_store                   = "machine"
    certificate_validity_period_scale   = "years"
    certificate_validity_period_value   = 1
    subject_name_format                 = "custom"
    subject_name_format_string          = "CN={{AAD_Device_ID}}"
    root_certificate_odata_bind         = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('00000000-0000-0000-0000-000000000001')"
    key_size                           = "size4096"
    key_usage                          = ["digitalSignature", "keyEncipherment"]
    
    custom_subject_alternative_names = [
      {
        san_type = "userPrincipalName"
        name     = "some-upn"
      },
      {
        san_type = "emailAddress"
        name     = "some-email"
      },
      {
        san_type = "domainNameService"
        name     = "some-dns"
      },
      {
        san_type = "universalResourceIdentifier"
        name     = "some-uri"
      }
    ]
    
    extended_key_usages = [
      {
        name              = "Any Purpose"
        object_identifier = "2.5.29.37.0"
      },
      {
        name              = "Client Authentication"
        object_identifier = "1.3.6.1.5.5.7.3.2"
      },
      {
        name              = "Secure Email"
        object_identifier = "1.3.6.1.5.5.7.3.4"
      },
      {
        name              = "custom"
        object_identifier = "7.01.4"
      }
    ]
    
    scep_server_urls = [
      "https://something.com",
      "https://something2.com"
    ]
    
    allow_all_apps_access = true
  }

  role_scope_tag_ids = ["00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"]

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000001"
      filter_id   = "00000000-0000-0000-0000-000000000002"
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
      group_id    = "00000000-0000-0000-0000-000000000002"
    }
  ]

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}



# Example 5: macOS PKCS Certificate Profile
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "pkcs_cert_example" {
  display_name = "macos pkcs certificate example"
  description  = "macos pkcs certificate example"

  pkcs_certificate = {
    deployment_channel                  = "deviceChannel"
    renewal_threshold_percentage        = 20
    certificate_store                   = "machine"
    certificate_validity_period_scale   = "years"
    certificate_validity_period_value   = 1
    subject_name_format                 = "custom"
    subject_name_format_string          = "CN={{AAD_Device_ID}}"
    certification_authority             = "some-auth"
    certification_authority_name        = "some-name"
    certificate_template_name           = "some-template-name"
    
    custom_subject_alternative_names = [
      {
        san_type = "emailAddress"
        name     = "some-email"
      },
      {
        san_type = "userPrincipalName"
        name     = "some-upn"
      },
      {
        san_type = "domainNameService"
        name     = "some-dns"
      },
      {
        san_type = "universalResourceIdentifier"
        name     = "some-uri"
      },
      {
        san_type = "customAzureADAttribute"
        name     = "some-custom-att"
      },
      {
        san_type = "emailAddress"
        name     = "some-other-email"
      }
    ]
    
    allow_all_apps_access = true
  }

  role_scope_tag_ids = ["00000000-0000-0000-0000-000000000001"]

  assignments = [
    {
      type        = "groupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000001"
      filter_id   = "00000000-0000-0000-0000-000000000002"
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
      group_id    = "00000000-0000-0000-0000-000000000002"
    }
  ]

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}