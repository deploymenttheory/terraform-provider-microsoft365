resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "test" {
  display_name     = "acc-test-windows-certificate-${random_string.suffix.result}"
  mail_nickname    = "acc-test-windows-certificate-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  hard_delete      = true
}

resource "time_sleep" "wait_for_group" {
  depends_on      = [microsoft365_graph_beta_groups_group.test]
  create_duration = "15s"
}

resource "microsoft365_graph_beta_device_management_windows_trusted_root_certificate" "test" {
  depends_on               = [time_sleep.wait_for_group]
  display_name             = "acc-test-windows-certificate-${random_string.suffix.result}"
  description              = "Trusted certificate provider acceptance test"
  cert_file_name           = "root.cer"
  trusted_root_certificate = "CERTIFICATE_BASE64"
  destination_store        = "computerCertStoreRoot"
  assignments = [{
    type     = "groupAssignmentTarget"
    group_id = microsoft365_graph_beta_groups_group.test.id
  }]
}
