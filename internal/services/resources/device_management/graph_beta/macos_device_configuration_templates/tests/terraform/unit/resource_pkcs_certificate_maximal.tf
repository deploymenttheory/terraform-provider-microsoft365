
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "pkcs_cert_example" {
  display_name = "unit-test-macOS-pkcs-certificate-example"
  description  = "PKCS certificate profile for user authentication"

  pkcs_certificate = {
    deployment_channel                = "userChannel"
    renewal_threshold_percentage      = 30
    certificate_store                 = "user"
    certificate_validity_period_scale = "months"
    certificate_validity_period_value = 12
    subject_name_format               = "commonNameIncludingEmail"
    subject_name_format_string        = "CN={{UserName}},E={{EmailAddress}},O=Example Corp"
    certification_authority           = "ExampleCA\\ExampleCA-CA"
    certification_authority_name      = "ExampleCA-CA"
    certificate_template_name         = "UserAuthentication"

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
      filter_id   = "00000000-0000-0000-0000-000000000003"
      filter_type = "include"
    },
    {
      type        = "exclusionGroupAssignmentTarget"
      group_id    = "00000000-0000-0000-0000-000000000002"
      filter_id   = "00000000-0000-0000-0000-000000000003"
      filter_type = "exclude"
    }
  ]

  timeouts = {
    create = "30m"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}