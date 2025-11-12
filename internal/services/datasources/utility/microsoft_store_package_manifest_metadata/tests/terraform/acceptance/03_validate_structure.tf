# Acceptance Test: Validate full package manifest structure
# Using WPS Office 2022 (XP8M1ZJCZ99QJW) - has rich metadata including locales, installers, agreements
data "microsoft365_utility_microsoft_store_package_manifest_metadata" "test" {
  package_identifier = "XP8M1ZJCZ99QJW"
}

