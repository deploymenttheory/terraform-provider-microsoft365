# Example: Look up a license by its String ID
# String IDs are used in PowerShell v1.0 and the skuPartNumber property in Microsoft Graph
data "microsoft365_utility_licensing_service_plan_reference" "enterprise_pack" {
  string_id = "ENTERPRISEPACK" # Microsoft 365 E3
}

# Use the GUID in a license assignment
resource "microsoft365_graph_beta_users_user_license_assignment" "example" {
  user_id = "user-id-here"
  sku_id  = data.microsoft365_utility_licensing_service_plan_reference.enterprise_pack.matching_products[0].guid
}

