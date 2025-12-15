# Example: Get tenant information by tenant ID
# This example demonstrates how to retrieve tenant information using a tenant ID (GUID)
# Useful when you know the tenant ID and need to validate or retrieve tenant details

data "microsoft365_graph_beta_identity_and_access_tenant_information" "by_tenant_id" {
  filter_type  = "tenant_id"
  filter_value = "6babcaad-604b-40ac-a9d7-9fd97c0b779f" # Replace with your target tenant ID

  timeouts = {
    read = "1m"
  }
}

# Output tenant details
output "tenant_info_by_id" {
  value = {
    tenant_id             = data.microsoft365_graph_beta_identity_and_access_tenant_information.by_tenant_id.tenant_id
    display_name          = data.microsoft365_graph_beta_identity_and_access_tenant_information.by_tenant_id.display_name
    default_domain_name   = data.microsoft365_graph_beta_identity_and_access_tenant_information.by_tenant_id.default_domain_name
    federation_brand_name = data.microsoft365_graph_beta_identity_and_access_tenant_information.by_tenant_id.federation_brand_name
  }
  description = "Tenant information retrieved by tenant ID"
}

