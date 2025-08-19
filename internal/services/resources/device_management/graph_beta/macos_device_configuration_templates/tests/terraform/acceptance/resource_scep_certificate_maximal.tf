
resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "scep_cert_example" {
  display_name = "scep-example"
  description  = "scep-example"

  scep_certificate = {
    deployment_channel                = "deviceChannel"
    renewal_threshold_percentage      = 20
    certificate_store                 = "machine"
    certificate_validity_period_scale = "years"
    certificate_validity_period_value = 1
    subject_name_format               = "custom"
    subject_name_format_string        = "CN={{AAD_Device_ID}}"
    root_certificate_odata_bind       = "https://graph.microsoft.com/beta/deviceManagement/deviceConfigurations('afeb20c3-48a0-4bbd-8191-c9c9f2fd62d2')"
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
