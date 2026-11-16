# Enable Microsoft Entra ID security defaults for the tenant.
#
# WARNING: this enforces security defaults, including MFA registration, for every user in
# the tenant. Graph rejects enabling security defaults while Conditional Access policies
# are enabled, so disable those policies first or keep is_enabled = false.
resource "microsoft365_graph_beta_identity_and_access_security_defaults_enforcement_policy" "example" {
  is_enabled = true

  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}
