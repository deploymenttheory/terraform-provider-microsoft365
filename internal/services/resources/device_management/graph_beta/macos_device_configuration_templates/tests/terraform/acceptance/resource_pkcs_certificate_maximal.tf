resource "random_string" "pkcs_cert_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "pkcs_cert_group_1" {
  display_name     = "acc-test-macos-pkcs-cert-group-1-${random_string.pkcs_cert_suffix.result}"
  mail_nickname    = "acc-test-pkcs-cert-1-${random_string.pkcs_cert_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "pkcs_cert_group_2" {
  display_name     = "acc-test-macos-pkcs-cert-group-2-${random_string.pkcs_cert_suffix.result}"
  mail_nickname    = "acc-test-pkcs-cert-2-${random_string.pkcs_cert_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "pkcs_cert_group_3" {
  display_name     = "acc-test-macos-pkcs-cert-group-3-${random_string.pkcs_cert_suffix.result}"
  mail_nickname    = "acc-test-pkcs-cert-3-${random_string.pkcs_cert_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "pkcs_cert_group_4" {
  display_name     = "acc-test-macos-pkcs-cert-group-4-${random_string.pkcs_cert_suffix.result}"
  mail_nickname    = "acc-test-pkcs-cert-4-${random_string.pkcs_cert_suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "microsoft365_graph_beta_device_management_macos_device_configuration_templates" "pkcs_cert_example" {
  display_name = "acc-test-macOS-pkcs-cert-${random_string.pkcs_cert_suffix.result}"
  description  = "PKCS certificate profile for user authentication"

  pkcs_certificate = {
    deployment_channel                = "deviceChannel"
    renewal_threshold_percentage      = 20
    certificate_store                 = "machine"
    certificate_validity_period_scale = "months"
    certificate_validity_period_value = 12
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

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.pkcs_cert_group_1.id
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.pkcs_cert_group_2.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.pkcs_cert_group_3.id
    },
    {
      type     = "exclusionGroupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.pkcs_cert_group_4.id
    }
  ]

  depends_on = [
    microsoft365_graph_beta_groups_group.pkcs_cert_group_1,
    microsoft365_graph_beta_groups_group.pkcs_cert_group_2,
    microsoft365_graph_beta_groups_group.pkcs_cert_group_3,
    microsoft365_graph_beta_groups_group.pkcs_cert_group_4
  ]

  timeouts = {
    create = "50s"
    read   = "5m"
    update = "30m"
    delete = "30m"
  }
}
