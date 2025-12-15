# Example: Get tenant information by domain name
# This example demonstrates how to retrieve tenant information using a domain name
# Useful when you know the domain name and need to find the tenant ID or other tenant details

data "microsoft365_graph_beta_identity_and_access_tenant_information" "by_domain_name" {
  filter_type  = "domain_name"
  filter_value = "contoso.com" # Replace with your target domain name

  timeouts = {
    read = "1m"
  }
}

# Output tenant details
output "tenant_info_by_domain" {
  value = {
    tenant_id             = data.microsoft365_graph_beta_identity_and_access_tenant_information.by_domain_name.tenant_id
    display_name          = data.microsoft365_graph_beta_identity_and_access_tenant_information.by_domain_name.display_name
    default_domain_name   = data.microsoft365_graph_beta_identity_and_access_tenant_information.by_domain_name.default_domain_name
    federation_brand_name = data.microsoft365_graph_beta_identity_and_access_tenant_information.by_domain_name.federation_brand_name
  }
  description = "Tenant information retrieved by domain name"
}

