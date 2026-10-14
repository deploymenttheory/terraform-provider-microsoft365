# Cross-tenant access policy can be imported using any ID (it will be normalized to "crossTenantAccessPolicy")
terraform import microsoft365_graph_beta_identity_and_access_cross_tenant_access_policy.example crossTenantAccessPolicy

# Set the Terraform-only destroy flag explicitly to avoid a subsequent state update.
terraform import microsoft365_graph_beta_identity_and_access_cross_tenant_access_policy.example 'crossTenantAccessPolicy:restore_defaults_on_destroy=true'
