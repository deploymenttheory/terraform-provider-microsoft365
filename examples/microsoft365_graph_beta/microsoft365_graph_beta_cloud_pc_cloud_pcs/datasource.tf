# Example: Retrieve all Cloud PCs

data "microsoft365_graph_beta_cloud_pc_cloud_pcs" "all" {
  filter_type = "all"
}

# Output: List all Cloud PC IDs
output "all_cloud_pc_ids" {
  value = [for pc in data.microsoft365_graph_beta_cloud_pc_cloud_pcs.all.items : pc.id]
}

# Output: Show all details for the first Cloud PC (if present)
output "first_cloud_pc_details" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pcs.all.items[0]
}

# Example: Retrieve a specific Cloud PC by ID
data "microsoft365_graph_beta_cloud_pc_cloud_pcs" "by_id" {
  filter_type  = "id"
  filter_value = "662009bc-7732-4f6f-8726-25883518ffff" # Replace with an actual ID
}

output "cloud_pc_by_id" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pcs.by_id.items[0]
}

# Example: Retrieve Cloud PCs by display name substring
data "microsoft365_graph_beta_cloud_pc_cloud_pcs" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Demo" # This will match Cloud PCs with "Demo" in their name
}

output "cloud_pcs_by_display_name" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pcs.by_display_name.items
}

# Example: Retrieve Cloud PCs by user principal name
data "microsoft365_graph_beta_cloud_pc_cloud_pcs" "by_upn" {
  filter_type  = "user_principal_name"
  filter_value = "user@contoso.com" # Replace with an actual UPN
}

output "cloud_pcs_by_upn" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pcs.by_upn.items
}

# Example: Retrieve Cloud PCs by status
data "microsoft365_graph_beta_cloud_pc_cloud_pcs" "by_status" {
  filter_type  = "status"
  filter_value = "provisioned" # Valid values include provisioned, provisioning, failed, etc.
}

output "provisioned_cloud_pcs" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pcs.by_status.items
}

# Example: Retrieve Cloud PCs by product type (using server-side filtering)
data "microsoft365_graph_beta_cloud_pc_cloud_pcs" "by_product_type" {
  filter_type  = "product_type"
  filter_value = "enterprise"
}

output "enterprise_cloud_pcs" {
  value = data.microsoft365_graph_beta_cloud_pc_cloud_pcs.by_product_type.items
}

# Example: Show Cloud PC status summary
output "cloud_pc_status_summary" {
  value = {
    total_count        = length(data.microsoft365_graph_beta_cloud_pc_cloud_pcs.all.items)
    provisioned_count  = length([for pc in data.microsoft365_graph_beta_cloud_pc_cloud_pcs.all.items : pc if pc.status == "provisioned"])
    provisioning_count = length([for pc in data.microsoft365_graph_beta_cloud_pc_cloud_pcs.all.items : pc if pc.status == "provisioning"])
    failed_count       = length([for pc in data.microsoft365_graph_beta_cloud_pc_cloud_pcs.all.items : pc if pc.status == "failed"])
  }
} 