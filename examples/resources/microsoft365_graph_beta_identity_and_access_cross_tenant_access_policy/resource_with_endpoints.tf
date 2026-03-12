resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_policy" "example" {
  allowed_cloud_endpoints = [
    "microsoftonline.us",
    "partner.microsoftonline.cn",
  ]
}
