---
page_title: "microsoft365_graph_beta_device_management_macos_device_configuration_templates Resource - terraform-provider-microsoft365"
subcategory: "Device Management"

description: |-
  Manages macOS configuration templates in Microsoft Intune. This resource creates device configurations for macOS devices including custom configuration profiles, preference files, trusted certificates, and certificate profiles (SCEP/PKCS).
---

# microsoft365_graph_beta_device_management_macos_device_configuration_templates (Resource)

Manages macOS configuration templates in Microsoft Intune. This resource creates device configurations for macOS devices including custom configuration profiles, preference files, trusted certificates, and certificate profiles (SCEP/PKCS).

## Microsoft Documentation

- [macOS Custom Configuration resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoscustomconfiguration?view=graph-rest-beta)
- [macOS PKCS Certificate resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macospkcscertificateprofile?view=graph-rest-beta)
- [macOS SCEP Certificate resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macosscepcertificateprofile?view=graph-rest-beta)
- [macOS Trusted Certificate resource type](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macostrustedrootcertificate?view=graph-rest-beta)
- [macOS Preference File resource type](https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-macoscustomappconfiguration-create?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this resource.

### Microsoft Graph

- **Application**: `DeviceManagementConfiguration.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.27.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Terraform resource configuration for Microsoft 365 Graph Beta Device Management macOS Configuration Templates

# Example 1: macOS Custom Configuration Template
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "custom_configuration_example" {
  display_name = "macos custom mobileconfig example"
  description  = "macos custom mobileconfig example"

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
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
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
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
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
    deployment_channel       = "deviceChannel"
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
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
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
    deployment_channel                = "deviceChannel"
    renewal_threshold_percentage      = 20
    certificate_store                 = "machine"
    certificate_validity_period_scale = "years"
    certificate_validity_period_value = 1
    subject_name_format               = "custom"
    subject_name_format_string        = "CN={{AAD_Device_ID}}"
    root_certificate_odata_bind       = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('00000000-0000-0000-0000-000000000001')"
    key_size                          = "size4096"
    key_usage                         = ["digitalSignature", "keyEncipherment"]

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
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
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
    deployment_channel                = "deviceChannel"
    renewal_threshold_percentage      = 20
    certificate_store                 = "machine"
    certificate_validity_period_scale = "years"
    certificate_validity_period_value = 1
    subject_name_format               = "custom"
    subject_name_format_string        = "CN={{AAD_Device_ID}}"
    certification_authority           = "some-auth"
    certification_authority_name      = "some-name"
    certificate_template_name         = "some-template-name"

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
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    }
  ]

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name for the macOS configuration template.

### Optional

- `assignments` (Attributes Set) Assignments for the device configuration. Each assignment specifies the target group and schedule for script execution. Supports group filters. (see [below for nested schema](#nestedatt--assignments))
- `custom_configuration` (Attributes) The custom configuration template allows IT admins to assign settings that aren't built into Intune yet. For macOS devices, you can import a .mobileconfig file that you created using Profile Manager or a different tool. (see [below for nested schema](#nestedatt--custom_configuration))
- `description` (String) The description for the macOS configuration template.
- `pkcs_certificate` (Attributes) PKCS certificate profile configuration for macOS devices. (see [below for nested schema](#nestedatt--pkcs_certificate))
- `preference_file` (Attributes) Configure a preference file that uses the standard property list (.plist) format to define preferences for apps and the device. (see [below for nested schema](#nestedatt--preference_file))
- `role_scope_tag_ids` (Set of String) Set of scope tag IDs for this Settings Catalog template profile.
- `scep_certificate` (Attributes) SCEP certificate profile configuration for macOS devices. (see [below for nested schema](#nestedatt--scep_certificate))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `trusted_certificate` (Attributes) Trusted root certificate configuration for macOS devices. (see [below for nested schema](#nestedatt--trusted_certificate))

### Read-Only

- `id` (String) The unique identifier for the macOS configuration template.

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `type` (String) Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.

Optional:

- `filter_id` (String) ID of the filter to apply to the assignment.
- `filter_type` (String) Type of filter to apply. Must be one of: 'include', 'exclude', or 'none'.
- `group_id` (String) The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.


<a id="nestedatt--custom_configuration"></a>
### Nested Schema for `custom_configuration`

Required:

- `deployment_channel` (String) Select the channel you want to use to deploy your configuration profile. If the channel doesn’t match what’s listed for the payload in Apple documentation, deployment could fail. The selected channel cannot be changed once the profile has been created. Possible values are: deviceChannel, userChannel.
- `payload` (String) The macOS configuration payload (.mobileconfig / .plist) file content.
- `payload_file_name` (String) The profile name displayed to users.
- `payload_name` (String) The name of the payload configuration.


<a id="nestedatt--pkcs_certificate"></a>
### Nested Schema for `pkcs_certificate`

Required:

- `subject_name_format` (String) Defaults to custom.

Optional:

- `allow_all_apps_access` (Boolean) Whether to allow all applications to access the certificate.
- `certificate_store` (String) The certificate store location. Possible values are: user, machine.
- `certificate_template_name` (String) The certificate template name for PKCS certificates.
- `certificate_validity_period_scale` (String) The certificate validity period scale. Possible values are: days, months, years.
- `certificate_validity_period_value` (Number) The certificate validity period value.
- `certification_authority` (String) The certification authority for PKCS certificates.
- `certification_authority_name` (String) The certification authority name for PKCS certificates.
- `custom_subject_alternative_names` (Attributes Set) Custom Subject Alternative Names for the certificate. (see [below for nested schema](#nestedatt--pkcs_certificate--custom_subject_alternative_names))
- `deployment_channel` (String) The deployment channel for the certificate. Possible values are: deviceChannel, userChannel.
- `renewal_threshold_percentage` (Number) The certificate renewal threshold percentage (1-99).
- `subject_name_format_string` (String) Custom format to use with SubjectNameFormat = Custom. Example: CN={{AAD_Device_ID}},O={{Organization}}

<a id="nestedatt--pkcs_certificate--custom_subject_alternative_names"></a>
### Nested Schema for `pkcs_certificate.custom_subject_alternative_names`

Required:

- `name` (String) The SAN value/name.
- `san_type` (String) The SAN type. Possible values are: emailAddress, userPrincipalName, customAzureADAttribute, domainNameService, universalResourceIdentifier.



<a id="nestedatt--preference_file"></a>
### Nested Schema for `preference_file`

Required:

- `bundle_id` (String) The bundle ID (Preference domain name) of the application this preference file applies to. Typically in the format com.company.appname.
- `configuration_xml` (String) The base64-encoded XML configuration content (.plist file content).
- `file_name` (String) The file name of the preference file (.plist file).


<a id="nestedatt--scep_certificate"></a>
### Nested Schema for `scep_certificate`

Required:

- `certificate_validity_period_scale` (String) The certificate validity period scale. Possible values are: days, months, years.
- `certificate_validity_period_value` (Number) The certificate validity period value.
- `deployment_channel` (String) The deployment channel for the certificate. Possible values are: deviceChannel, userChannel.
- `extended_key_usages` (Attributes Set) Extended key usage settings for the certificate. (see [below for nested schema](#nestedatt--scep_certificate--extended_key_usages))
- `key_size` (String) The key size in bits for the certificate.2048 is the recommended minimum key length.Possible values are: size1024, size2048, size4096.
- `key_usage` (Set of String) Key usage options for the certificate. Possible values are: keyEncipherment, digitalSignature.
- `renewal_threshold_percentage` (Number) The certificate renewal threshold percentage (1-99).
- `root_certificate_odata_bind` (String) Reference to the pre existing trusted root certificate configuration for the odata bind.Valid format is "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('00000000-0000-0000-0000-000000000000')".Or you can supply just the ID of the certificate configuration. e.g. '00000000-0000-0000-0000-000000000000'
- `subject_name_format` (String) Defaults to custom.
- `subject_name_format_string` (String) Select how Intune automatically creates the subject name in the certificate request. If the certificate is for a user, you can also include the user's email address in the subject name. Please review subject name documentation 'https://learn.microsoft.com/en-us/intune/intune-service/protect/certificates-profile-scep'on how to best use the Subject name format field.Custom. Example: CN={{AAD_Device_ID}},O={{Organization}}

Optional:

- `allow_all_apps_access` (Boolean) Whether to allow all applications to access the certificate.
- `certificate_store` (String) The certificate store location. Possible values are: user, machine.
- `custom_subject_alternative_names` (Attributes Set) Custom Subject Alternative Names for the certificate. (see [below for nested schema](#nestedatt--scep_certificate--custom_subject_alternative_names))
- `scep_server_urls` (Set of String) SCEP server URL(s) for certificate enrollment.

<a id="nestedatt--scep_certificate--extended_key_usages"></a>
### Nested Schema for `scep_certificate.extended_key_usages`

Required:

- `name` (String) The extended key usage name.
- `object_identifier` (String) The extended key usage object identifier (OID).


<a id="nestedatt--scep_certificate--custom_subject_alternative_names"></a>
### Nested Schema for `scep_certificate.custom_subject_alternative_names`

Required:

- `name` (String) The SAN value/name.
- `san_type` (String) The SAN type. Possible values are: emailAddress, userPrincipalName, customAzureADAttribute, domainNameService, universalResourceIdentifier.



<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--trusted_certificate"></a>
### Nested Schema for `trusted_certificate`

Required:

- `cert_file_name` (String) The file name of the certificate file (.cer file).
- `deployment_channel` (String) The deployment channel for the certificate. Possible values are: deviceChannel, userChannel.
- `trusted_root_certificate` (String) The base64-encoded trusted root certificate content. This should be a filebase64() encoded string. e.g filebase64("my-root-cert.cer")

## Important Notes

- **Assignment Required**: Policies must be assigned to device or user groups to be deployed.

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

# Example 1: Import Custom Configuration Template
print_status "Example 1: Importing a Custom Configuration Template"
echo "terraform import microsoft365_graph_beta_device_management_macos_device_configuration_templates.custom_config_example \"12345678-1234-1234-1234-123456789012\""
echo ""

# Example 2: Import Preference File Configuration
print_status "Example 2: Importing a Preference File Configuration"
echo "terraform import microsoft365_graph_beta_device_management_macos_device_configuration_templates.preference_file_example \"87654321-4321-4321-4321-210987654321\""
echo ""

# Example 3: Import Trusted Certificate Configuration
print_status "Example 3: Importing a Trusted Certificate Configuration"
echo "terraform import microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example \"11111111-2222-3333-4444-555555555555\""
echo ""

# Example 4: Import SCEP Certificate Profile
print_status "Example 4: Importing a SCEP Certificate Profile"
echo "terraform import microsoft365_graph_beta_device_management_macos_device_configuration_templates.scep_cert_example \"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee\""
echo ""

# Example 5: Import PKCS Certificate Profile
print_status "Example 5: Importing a PKCS Certificate Profile"
echo "terraform import microsoft365_graph_beta_device_management_macos_device_configuration_templates.pkcs_cert_example \"ffffffff-eeee-dddd-cccc-bbbbbbbbbbbb\""
echo ""
```

