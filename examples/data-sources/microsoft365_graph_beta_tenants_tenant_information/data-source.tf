# Application permission: CrossTenantInformation.ReadBasic.All.
data "microsoft365_graph_beta_tenants_tenant_information" "by_id" {
  tenant_id = "6babcaad-604b-40ac-a9d7-9fd97c0b779f"
}
data "microsoft365_graph_beta_tenants_tenant_information" "by_domain" {
  domain_name = "contoso.com"
}
output "tenant_information" {
  value = data.microsoft365_graph_beta_tenants_tenant_information.by_domain.items[0]
}
