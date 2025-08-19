
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "pkcs_cert_example" {
  display_name = "pkcs-cert"
  description  = "pkcs-cert"

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
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_3.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_4.id
    }
  ]

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}