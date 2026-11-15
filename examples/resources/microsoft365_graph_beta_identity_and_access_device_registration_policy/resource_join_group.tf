# Allow only members of a group to join devices to Microsoft Entra ID.
resource "microsoft365_graph_beta_groups_group" "device_joiners" {
  display_name     = "Device Joiners"
  mail_nickname    = "device-joiners"
  mail_enabled     = false
  security_enabled = true
  description      = "Users allowed to join devices to Microsoft Entra ID"
}

resource "microsoft365_graph_beta_identity_and_access_device_registration_policy" "example" {
  azure_ad_join = {
    allowed_to_join = {
      type   = "enumerated"
      groups = [microsoft365_graph_beta_groups_group.device_joiners.id]
    }
  }
}
