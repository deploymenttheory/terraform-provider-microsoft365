# Unit Test: Search for service plans by ID
data "microsoft365_utility_licensing_service_plan_reference" "test" {
  service_plan_id = "EXCHANGE"
}

