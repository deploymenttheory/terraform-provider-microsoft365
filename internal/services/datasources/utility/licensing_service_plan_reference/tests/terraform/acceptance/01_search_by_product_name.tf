# Acceptance Test: Search for products by name
data "microsoft365_utility_licensing_service_plan_reference" "test" {
  product_name = "Microsoft 365 E3"
}

