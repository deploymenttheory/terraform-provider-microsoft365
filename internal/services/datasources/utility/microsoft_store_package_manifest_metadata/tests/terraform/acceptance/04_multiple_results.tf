# Acceptance Test: Search that returns multiple results
data "microsoft365_utility_microsoft_store_package_manifest_metadata" "test" {
  search_term = "Microsoft"
}

