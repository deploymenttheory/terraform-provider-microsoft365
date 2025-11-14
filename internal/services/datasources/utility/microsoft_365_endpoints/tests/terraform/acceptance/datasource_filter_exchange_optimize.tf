data "microsoft365_utility_microsoft_365_endpoints" "test" {
  instance      = "worldwide"
  service_areas = ["Exchange"]
  categories    = ["Optimize"]
}

