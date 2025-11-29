# Acceptance Test: Search for product by string ID
data "microsoft365_utility_licensing_service_plan_reference" "test" {
  string_id = "SPE_E3_RPA1" # Microsoft 365 E3 (no Teams)
}

