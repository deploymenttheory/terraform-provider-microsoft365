resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_groups_group" "join" {
  display_name     = "acc-test-device-registration-join-${random_string.suffix.result}"
  mail_nickname    = "acc-test-drp-join-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Acceptance test group allowed to join devices to Microsoft Entra ID"
  hard_delete      = true
}

resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "test" {
  azure_ad_join = {
    allowed_to_join = {
      type   = "enumerated"
      groups = [microsoft365_graph_beta_groups_group.join.id]
    }
  }
}
