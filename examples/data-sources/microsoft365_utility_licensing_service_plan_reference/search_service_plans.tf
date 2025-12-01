# Example: Find which products include Exchange Online service plans
data "microsoft365_utility_licensing_service_plan_reference" "exchange_plans" {
  service_plan_name = "Exchange Online"
}

# Output the service plans and which products include them
output "exchange_service_plans" {
  value = [
    for plan in data.microsoft365_utility_licensing_service_plan_reference.exchange_plans.matching_service_plans : {
      service_plan_id   = plan.id
      service_plan_name = plan.name
      included_in       = [for sku in plan.included_in_skus : sku.product_name]
    }
  ]
}

