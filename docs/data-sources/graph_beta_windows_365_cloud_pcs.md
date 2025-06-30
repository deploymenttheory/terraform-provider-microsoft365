---
page_title: "microsoft365_graph_beta_windows_365_cloud_pcs Data Source - terraform-provider-microsoft365"
subcategory: "Windows 365"
description: |-
  Retrieves Cloud PC devices from Microsoft Intune using the /deviceManagement/virtualEndpoint/cloudPCs endpoint. Supports filtering by all, id, display_name, user_principal_name, status, or product_type for comprehensive Cloud PC management.
---

# microsoft365_graph_beta_windows_365_cloud_pcs (Data Source)

Retrieves Cloud PC devices from Microsoft Intune using the `/deviceManagement/virtualEndpoint/cloudPCs` endpoint. Supports filtering by all, id, display_name, user_principal_name, status, or product_type for comprehensive Cloud PC management.

This data source allows you to list and filter Cloud PCs in your tenant, providing details such as status, assigned user, provisioning policy, service plan, and more.

## Microsoft Documentation

- [List cloudPCs](https://learn.microsoft.com/en-us/graph/api/virtualendpoint-list-cloudpcs?view=graph-rest-beta)

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `CloudPC.Read.All`, `CloudPC.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.18.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
# Example: Retrieve all Cloud PCs
data "microsoft365_graph_beta_windows_365_cloud_pcs" "all" {
  filter_type = "all"
}

output "all_cloud_pcs_full" {
  value = [
    for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : {
      id                       = pc.id
      display_name             = pc.display_name
      aad_device_id            = pc.aad_device_id
      image_display_name       = pc.image_display_name
      managed_device_id        = pc.managed_device_id
      managed_device_name      = pc.managed_device_name
      provisioning_policy_id   = pc.provisioning_policy_id
      provisioning_policy_name = pc.provisioning_policy_name
      on_premises_connection_name = pc.on_premises_connection_name
      service_plan_id          = pc.service_plan_id
      service_plan_name        = pc.service_plan_name
      service_plan_type        = pc.service_plan_type
      status                   = pc.status
      user_principal_name      = pc.user_principal_name
      last_modified_date_time  = pc.last_modified_date_time
      status_detail_code       = pc.status_detail_code
      status_detail_message    = pc.status_detail_message
      grace_period_end_date_time = pc.grace_period_end_date_time
      provisioning_type        = pc.provisioning_type
      device_region_name       = pc.device_region_name
      disk_encryption_state    = pc.disk_encryption_state
      product_type             = pc.product_type
      user_account_type        = pc.user_account_type
      enable_single_sign_on    = pc.enable_single_sign_on
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
  filter_type   = "odata"
  odata_filter  = "status eq 'provisioned'"
  odata_select  = "id,displayName,status,userPrincipalName"
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
  filter_type   = "odata"
  odata_filter  = "status eq 'provisioned' and productType eq 'enterprise'"
  odata_select  = "id,displayName,status,userPrincipalName,serviceplanName"
  odata_top     = 10
}

# Output: Basic information about all Cloud PCs
output "all_cloud_pcs_basic_info" {
  value = [
    for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : {
      id               = pc.id
      display_name     = pc.display_name
      status           = pc.status
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
      id                        = pc.id
      display_name              = pc.display_name
      provisioning_policy_name  = pc.provisioning_policy_name
      service_plan_name         = pc.service_plan_name
      service_plan_type         = pc.service_plan_type
      image_display_name        = pc.image_display_name
      provisioning_type         = pc.provisioning_type
      device_region_name        = pc.device_region_name
      disk_encryption_state     = pc.disk_encryption_state
      product_type              = pc.product_type
      user_account_type         = pc.user_account_type
      enable_single_sign_on     = pc.enable_single_sign_on
    }
  ]
}

# Output: Cloud PC device IDs for integration with other systems
output "cloud_pc_device_ids" {
  value = [
    for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : {
      cloud_pc_id        = pc.id
      aad_device_id      = pc.aad_device_id
      managed_device_id  = pc.managed_device_id
      managed_device_name = pc.managed_device_name
    }
  ]
}

# Output: Cloud PC status information for monitoring
output "cloud_pc_status_info" {
  value = [
    for pc in data.microsoft365_graph_beta_windows_365_cloud_pcs.all.items : {
      id                      = pc.id
      display_name            = pc.display_name
      status                  = pc.status
      status_detail_code      = pc.status_detail_code
      status_detail_message   = pc.status_detail_message
      last_modified_date_time = pc.last_modified_date_time
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
      id                       = pc.id
      display_name             = pc.display_name
      user_principal_name      = pc.user_principal_name
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filter_type` (String) Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `odata`. Use 'all' to retrieve all Cloud PCs, 'id' to retrieve a specific Cloud PC by its unique identifier, 'display_name' to filter by name, or 'odata' to use advanced OData query parameters.

### Optional

- `filter_value` (String) Value to filter by. Not required when filter_type is 'all'. For 'id', provide the Cloud PC ID. For other filters, provide the appropriate value to match.
- `odata_count` (Boolean) OData $count query parameter to include a count of items. Only applicable when filter_type is 'odata'.
- `odata_filter` (String) OData $filter query parameter. Only applicable when filter_type is 'odata'. Example: "status eq 'provisioned'"
- `odata_select` (String) OData $select query parameter to specify which fields to return. Only applicable when filter_type is 'odata'. Example: "id,displayName,status"
- `odata_top` (Number) OData $top query parameter to limit the number of items returned. Only applicable when filter_type is 'odata'.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `items` (Attributes List) The list of Cloud PCs that match the filter criteria. (see [below for nested schema](#nestedatt--items))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- `aad_device_id` (String) The Azure AD device ID associated with the Cloud PC.
- `device_region_name` (String) The Azure region where the Cloud PC is deployed.
- `disk_encryption_state` (String) The disk encryption state of the Cloud PC (e.g., encryptedUsingPlatformManagedKey).
- `display_name` (String) The display name of the Cloud PC.
- `enable_single_sign_on` (Boolean) Indicates whether single sign-on is enabled for the Cloud PC.
- `grace_period_end_date_time` (String) The date and time when the grace period for the Cloud PC ends.
- `id` (String) The unique identifier for the Cloud PC.
- `image_display_name` (String) The display name of the image used for the Cloud PC.
- `last_modified_date_time` (String) The date and time when the Cloud PC was last modified.
- `managed_device_id` (String) The managed device ID associated with the Cloud PC.
- `managed_device_name` (String) The name of the managed device associated with the Cloud PC.
- `on_premises_connection_name` (String) The name of the on-premises connection used for the Cloud PC.
- `product_type` (String) The product type of the Cloud PC (e.g., enterprise).
- `provisioning_policy_id` (String) The ID of the provisioning policy used for the Cloud PC.
- `provisioning_policy_name` (String) The name of the provisioning policy used for the Cloud PC.
- `provisioning_type` (String) The type of provisioning used for the Cloud PC (e.g., dedicated).
- `service_plan_id` (String) The ID of the service plan associated with the Cloud PC.
- `service_plan_name` (String) The name of the service plan associated with the Cloud PC.
- `service_plan_type` (String) The type of service plan associated with the Cloud PC.
- `status` (String) The current status of the Cloud PC (e.g., provisioned, provisioning, failed).
- `status_detail_code` (String) The error/warning code associated with the Cloud PC status.
- `status_detail_message` (String) The status message associated with the error code.
- `user_account_type` (String) The account type of the user on provisioned Cloud PCs (e.g., standardUser, administrator).
- `user_principal_name` (String) The user principal name (UPN) of the user assigned to the Cloud PC. 