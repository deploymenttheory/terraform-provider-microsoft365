# Unit Test: Both package_identifier and search_term provided (should fail)
data "microsoft365_utility_microsoft_store_package_manifest_metadata" "test" {
  package_identifier = "Microsoft.PowerToys"
  search_term        = "PowerToys"
}

