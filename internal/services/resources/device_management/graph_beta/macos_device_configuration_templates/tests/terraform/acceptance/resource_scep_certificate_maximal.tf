resource "random_string" "scep_cert_suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

resource "microsoft365_graph_beta_groups_group" "scep_cert_group_1" {
  display_name     = "acc-test-macos-scep-cert-group-1-${random_string.scep_cert_suffix.result}"
  mail_nickname    = "acc-test-scep-cert-1-${random_string.scep_cert_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "scep_cert_group_2" {
  display_name     = "acc-test-macos-scep-cert-group-2-${random_string.scep_cert_suffix.result}"
  mail_nickname    = "acc-test-scep-cert-2-${random_string.scep_cert_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "scep_cert_group_3" {
  display_name     = "acc-test-macos-scep-cert-group-3-${random_string.scep_cert_suffix.result}"
  mail_nickname    = "acc-test-scep-cert-3-${random_string.scep_cert_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "scep_cert_group_4" {
  display_name     = "acc-test-macos-scep-cert-group-4-${random_string.scep_cert_suffix.result}"
  mail_nickname    = "acc-test-scep-cert-4-${random_string.scep_cert_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "scep_cert_example" {
  display_name = "acc-test-macOS-scep-cert-${random_string.scep_cert_suffix.result}"
  description  = "SCEP certificate profile for device authentication"

  scep_certificate = {
    deployment_channel                = "deviceChannel"
    renewal_threshold_percentage      = 20
    certificate_store                 = "machine"
    certificate_validity_period_scale = "years"
    certificate_validity_period_value = 1
    subject_name_format               = "custom"
    subject_name_format_string        = "CN={{AAD_Device_ID}}"
    root_certificate_odata_bind       = microsoft365_graph_beta_device_management_macos_device_configuration_templates.trusted_cert_example.id
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

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.scep_cert_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.scep_cert_group_2.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.scep_cert_group_3.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.scep_cert_group_4.id
    }
  ]

  depends_on = [
    microsoft365_graph_beta_groups_group.scep_cert_group_1,
    microsoft365_graph_beta_groups_group.scep_cert_group_2,
    microsoft365_graph_beta_groups_group.scep_cert_group_3,
    microsoft365_graph_beta_groups_group.scep_cert_group_4
  ]

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}
