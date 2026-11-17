# Pre-logon Private Access profile: connects devices to Private Access before a
# user signs in to Windows, for example to reach a domain controller.
resource "microsoft365_graph_beta_identity_and_access_network_device_forwarding_profile" "pre_logon" {
  name        = "Pre-logon - domain controllers"
  description = "Private Access before Windows sign-in"
  priority    = 180
  state       = "enabled"

  assignment_rules = {
    membership = "enumerated"
    group_ids  = ["00000000-0000-0000-0000-000000000001"]
    platforms  = ["windows"]
  }
}
