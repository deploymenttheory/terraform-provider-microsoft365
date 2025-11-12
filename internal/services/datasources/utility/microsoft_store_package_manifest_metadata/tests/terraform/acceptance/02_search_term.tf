# Acceptance Test: Search for packages by search term
# Using "PC Manager" which returns Microsoft PC Manager package
data "microsoft365_utility_microsoft_store_package_manifest_metadata" "test" {
  search_term = "PC Manager"
}

