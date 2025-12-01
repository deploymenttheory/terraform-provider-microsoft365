# Example: Search for Microsoft 365 E3 license by product name
data "microsoft365_utility_licensing_service_plan_reference" "m365_e3" {
  product_name = "Microsoft 365 E3"
}

# Output the matching product details
output "m365_e3_details" {
  value = {
    product_name = data.microsoft365_utility_licensing_service_plan_reference.m365_e3.matching_products[0].product_name
    string_id    = data.microsoft365_utility_licensing_service_plan_reference.m365_e3.matching_products[0].string_id
    guid         = data.microsoft365_utility_licensing_service_plan_reference.m365_e3.matching_products[0].guid
  }
}

# Output the service plans included in Microsoft 365 E3
output "m365_e3_service_plans" {
  value = data.microsoft365_utility_licensing_service_plan_reference.m365_e3.matching_products[0].service_plans_included
}

