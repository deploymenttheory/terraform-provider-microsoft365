data "microsoft365_utility_microsoft_365_endpoint_reference" "test" {
  instance      = "worldwide"
  service_areas = ["Exchange"]
  categories    = ["Optimize"]
}

