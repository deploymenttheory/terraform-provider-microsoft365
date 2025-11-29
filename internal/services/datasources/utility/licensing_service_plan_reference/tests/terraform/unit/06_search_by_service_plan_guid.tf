# Unit Test: Search for service plan by GUID
data "microsoft365_utility_licensing_service_plan_reference" "test" {
  service_plan_guid = "113feb6c-3fe4-4440-bddc-54d774bf0318"
}

