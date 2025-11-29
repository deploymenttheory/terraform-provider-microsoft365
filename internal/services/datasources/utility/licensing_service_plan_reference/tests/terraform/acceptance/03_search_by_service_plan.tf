# Acceptance Test: Search for service plans and find which products include them
data "microsoft365_utility_licensing_service_plan_reference" "test" {
  service_plan_name = "Exchange Online"
}

