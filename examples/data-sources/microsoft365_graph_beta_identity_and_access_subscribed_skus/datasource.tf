# Example: Get all subscribed SKUs
data "microsoft365_graph_beta_identity_and_access_subscribed_skus" "all" {
  timeouts = {
    read = "30s"
  }
}

# Example: Filter by SKU part number
data "microsoft365_graph_beta_identity_and_access_subscribed_skus" "enterprise_premium" {
  sku_part_number = "ENTERPRISEPREMIUM"

  timeouts = {
    read = "30s"
  }
}

# Example: Filter by applies_to User
data "microsoft365_graph_beta_identity_and_access_subscribed_skus" "user_skus" {
  applies_to = "User"

  timeouts = {
    read = "30s"
  }
}

# Example: Get specific SKU by ID
data "microsoft365_graph_beta_identity_and_access_subscribed_skus" "specific_sku" {
  sku_id = "c7df2760-2c81-4ef7-b578-5b5392b571df"

  timeouts = {
    read = "30s"
  }
}

# Example: Filter by partial SKU part number match
data "microsoft365_graph_beta_identity_and_access_subscribed_skus" "premium_skus" {
  sku_part_number = "PREMIUM"

  timeouts = {
    read = "30s"
  }
}

# Output examples
output "all_skus_count" {
  value       = length(data.microsoft365_graph_beta_identity_and_access_subscribed_skus.all.subscribed_skus)
  description = "Total number of subscribed SKUs"
}

output "enterprise_premium_sku" {
  value       = data.microsoft365_graph_beta_identity_and_access_subscribed_skus.enterprise_premium.subscribed_skus
  description = "Enterprise Premium SKU details"
}

output "user_assignable_skus" {
  value = [
    for sku in data.microsoft365_graph_beta_identity_and_access_subscribed_skus.user_skus.subscribed_skus : {
      sku_part_number   = sku.sku_part_number
      consumed_units    = sku.consumed_units
      enabled_units     = sku.prepaid_units.enabled
      capability_status = sku.capability_status
    }
  ]
  description = "Summary of SKUs that can be assigned to users"
}

output "specific_sku_service_plans" {
  value       = length(data.microsoft365_graph_beta_identity_and_access_subscribed_skus.specific_sku.subscribed_skus) > 0 ? data.microsoft365_graph_beta_identity_and_access_subscribed_skus.specific_sku.subscribed_skus[0].service_plans : []
  description = "Service plans for the specific SKU"
}

output "premium_skus_summary" {
  value = [
    for sku in data.microsoft365_graph_beta_identity_and_access_subscribed_skus.premium_skus.subscribed_skus : {
      name               = sku.sku_part_number
      total_licenses     = sku.prepaid_units.enabled
      used_licenses      = sku.consumed_units
      available_licenses = sku.prepaid_units.enabled - sku.consumed_units
    }
  ]
  description = "License usage summary for Premium SKUs"
} 