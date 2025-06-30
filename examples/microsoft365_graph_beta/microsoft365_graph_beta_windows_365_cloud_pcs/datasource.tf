# Example: Retrieve all Cloud PCs
data "microsoft365_graph_beta_windows_365_cloud_pcs" "all" {
  filter_type = "all"
}

output "all_cloud_pcs_full" {
  value = [
    for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : {
      id                          = pc.id
      display_name                = pc.display_name
      aad_device_id               = pc.aad_device_id
      image_display_name          = pc.image_display_name
      managed_device_id           = pc.managed_device_id
      managed_device_name         = pc.managed_device_name
      provisioning_policy_id      = pc.provisioning_policy_id
      provisioning_policy_name    = pc.provisioning_policy_name
      on_premises_connection_name = pc.on_premises_connection_name
      service_plan_id             = pc.service_plan_id
      service_plan_name           = pc.service_plan_name
      service_plan_type           = pc.service_plan_type
      status                      = pc.status
      user_principal_name         = pc.user_principal_name
      last_modified_date_time     = pc.last_modified_date_time
      status_detail_code          = pc.status_detail_code
      status_detail_message       = pc.status_detail_message
      grace_period_end_date_time  = pc.grace_period_end_date_time
      provisioning_type           = pc.provisioning_type
      device_region_name          = pc.device_region_name
      disk_encryption_state       = pc.disk_encryption_state
      product_type                = pc.product_type
      user_account_type           = pc.user_account_type
      enable_single_sign_on       = pc.enable_single_sign_on
    }
  ]
}

# Example: Retrieve a specific Cloud PC by ID
data "microsoft365_graph_beta_windows_365_cloud_pcs" "by_id" {
  filter_type  = "id"
  filter_value = "662009bc-7732-4f6f-8726-25883518ffff" # Replace with an actual Cloud PC ID
}

# Example: Retrieve Cloud PCs by display name substring
data "microsoft365_graph_beta_windows_365_cloud_pcs" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Demo" # This will match Cloud PCs with "Demo" in their name
}

# Example: Retrieve Cloud PCs by user principal name (using OData)
data "microsoft365_graph_beta_windows_365_cloud_pcs" "by_upn" {
  filter_type  = "odata"
  odata_filter = "contains(userPrincipalName, 'user@contoso.com')" # Replace with an actual UPN
}

# Example: Retrieve Cloud PCs by status (using OData)
data "microsoft365_graph_beta_windows_365_cloud_pcs" "by_status" {
  filter_type  = "odata"
  odata_filter = "status eq 'provisioned'" # Valid values include provisioned, provisioning, failed, etc.
}

# Example: Retrieve Cloud PCs by product type (using OData filtering)
data "microsoft365_graph_beta_windows_365_cloud_pcs" "by_product_type" {
  filter_type  = "odata"
  odata_filter = "productType eq 'enterprise'"
}

# Example: Using OData to filter Cloud PCs with advanced criteria
data "microsoft365_graph_beta_windows_365_cloud_pcs" "by_odata_filter" {
  filter_type  = "odata"
  odata_filter = "status eq 'provisioned' and contains(displayName, 'Demo')"
}

# Example: Using OData to select specific fields (reduces state file size)
data "microsoft365_graph_beta_windows_365_cloud_pcs" "by_odata_select" {
  filter_type  = "odata"
  odata_filter = "status eq 'provisioned'"
  odata_select = "id,displayName,status,userPrincipalName"
}

# Example: Using OData to limit the number of results
data "microsoft365_graph_beta_windows_365_cloud_pcs" "by_odata_top" {
  filter_type  = "odata"
  odata_filter = "status eq 'provisioned'"
  odata_top    = 5
}

# Example: Using OData to include count in the response
data "microsoft365_graph_beta_windows_365_cloud_pcs" "by_odata_count" {
  filter_type  = "odata"
  odata_filter = "status eq 'provisioned'"
  odata_count  = true
}

# Example: Using OData with multiple parameters for optimized queries
data "microsoft365_graph_beta_windows_365_cloud_pcs" "optimized_query" {
  filter_type  = "odata"
  odata_filter = "status eq 'provisioned' and productType eq 'enterprise'"
  odata_select = "id,displayName,status,userPrincipalName,serviceplanName"
  odata_top    = 10
}

# Output: Basic information about all Cloud PCs
output "all_cloud_pcs_basic_info" {
  value = [
    for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : {
      id                  = pc.id
      display_name        = pc.display_name
      status              = pc.status
      user_principal_name = pc.user_principal_name
    }
  ]
}

# Output: Detailed information about a specific Cloud PC (first one from the list)
output "first_cloud_pc_details" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items[0]
}

# Output: Cloud PC configuration details
output "cloud_pc_configuration_details" {
  value = [
    for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : {
      id                       = pc.id
      display_name             = pc.display_name
      provisioning_policy_name = pc.provisioning_policy_name
      service_plan_name        = pc.service_plan_name
      service_plan_type        = pc.service_plan_type
      image_display_name       = pc.image_display_name
      provisioning_type        = pc.provisioning_type
      device_region_name       = pc.device_region_name
      disk_encryption_state    = pc.disk_encryption_state
      product_type             = pc.product_type
      user_account_type        = pc.user_account_type
      enable_single_sign_on    = pc.enable_single_sign_on
    }
  ]
}

# Output: Cloud PC device IDs for integration with other systems
output "cloud_pc_device_ids" {
  value = [
    for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : {
      cloud_pc_id         = pc.id
      aad_device_id       = pc.aad_device_id
      managed_device_id   = pc.managed_device_id
      managed_device_name = pc.managed_device_name
    }
  ]
}

# Output: Cloud PC status information for monitoring
output "cloud_pc_status_info" {
  value = [
    for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : {
      id                         = pc.id
      display_name               = pc.display_name
      status                     = pc.status
      status_detail_code         = pc.status_detail_code
      status_detail_message      = pc.status_detail_message
      last_modified_date_time    = pc.last_modified_date_time
      grace_period_end_date_time = pc.grace_period_end_date_time
    }
  ]
}

# Output: Cloud PC status summary for reporting
output "cloud_pc_status_summary" {
  value = {
    total_count        = length(data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items)
    provisioned_count  = length([for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : pc if pc.status == "provisioned"])
    provisioning_count = length([for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : pc if pc.status == "provisioning"])
    failed_count       = length([for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : pc if pc.status == "failed"])
  }
}

# Example: Finding Cloud PCs that need attention (failed or have warnings)
output "cloud_pcs_needing_attention" {
  value = [
    for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : {
      id                    = pc.id
      display_name          = pc.display_name
      status                = pc.status
      status_detail_code    = pc.status_detail_code
      status_detail_message = pc.status_detail_message
      user_principal_name   = pc.user_principal_name
    } if pc.status == "failed" || pc.status_detail_code != ""
  ]
}

# Example: Cloud PCs approaching grace period end
output "cloud_pcs_approaching_grace_period_end" {
  value = [
    for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : {
      id                         = pc.id
      display_name               = pc.display_name
      user_principal_name        = pc.user_principal_name
      grace_period_end_date_time = pc.grace_period_end_date_time
    } if pc.grace_period_end_date_time != ""
  ]
}

# Example: Group Cloud PCs by service plan for capacity planning
output "cloud_pcs_by_service_plan" {
  value = {
    for plan in distinct([for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : pc.service_plan_name if pc.service_plan_name != ""]) :
    plan => [
      for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items :
      {
        id           = pc.id
        display_name = pc.display_name
        status       = pc.status
      } if pc.service_plan_name == plan
    ]
  }
}

# Example: Group Cloud PCs by region for geographic distribution analysis
output "cloud_pcs_by_region" {
  value = {
    for region in distinct([for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : pc.device_region_name if pc.device_region_name != ""]) :
    region => length([
      for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : pc
      if pc.device_region_name == region
    ])
  }
}

# Example: Output from optimized OData query with selected fields
output "optimized_cloud_pcs" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pcs.optimized_query.items
}

# Example: Using OData select to minimize state file size
output "minimal_state_cloud_pcs" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pcs.by_odata_select.items
} 