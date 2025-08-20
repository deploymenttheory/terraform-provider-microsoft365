
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "scep_cert_example" {
  display_name = "unit-test-macOS-scep-certificate-example"
  description  = "SCEP certificate profile for device authentication"

  scep_certificate = {
    deployment_channel                = "deviceChannel"
    renewal_threshold_percentage      = 20
    certificate_store                 = "machine"
    certificate_validity_period_scale = "years"
    certificate_validity_period_value = 2
    subject_name_format               = "custom"
    subject_name_format_string        = "CN={{DeviceName}},O=Example Corp,C=US"
    root_certificate_odata_bind       = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('87654321-4321-4321-4321-210987654321')"
    key_size                          = "size2048"
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

    allow_all_apps_access = false
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
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000004"
    }
  ]


  timeouts = {
    create = "30m"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}
