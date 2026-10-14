# Every omitted/null top-level setting must remain inherited from the tenant
# defaults. Terraform removes partner overrides introduced outside Terraform.
resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings" "inherited" {
  tenant_id                   = "12345678-1234-1234-1234-123456789012"
  enforce_default_inheritance = true
  tenant_restrictions         = null
}
