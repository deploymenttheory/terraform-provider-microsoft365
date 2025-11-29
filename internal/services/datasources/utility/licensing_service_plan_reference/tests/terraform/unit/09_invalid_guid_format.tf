# Unit Test: Invalid GUID format (should fail)
data "microsoft365_utility_licensing_service_plan_reference" "test" {
  guid = "not-a-valid-guid"
}

