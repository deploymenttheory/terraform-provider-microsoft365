# Unit Test: Get package manifest by package identifier
# Using Microsoft PC Manager (9PM860492SZD) - a known stable package in Microsoft Store
data "microsoft365_utility_microsoft_store_package_manifest_metadata" "test" {
  package_identifier = "9PM860492SZD"
}

