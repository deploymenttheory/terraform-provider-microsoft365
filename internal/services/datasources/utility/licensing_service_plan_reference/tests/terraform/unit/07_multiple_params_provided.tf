# Unit Test: Multiple search parameters provided (should fail)
data "microsoft365_utility_licensing_service_plan_reference" "test" {
  product_name = "Microsoft 365 E3"
  string_id    = "ENTERPRISEPACK"
}

