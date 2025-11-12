data "microsoft365_utility_windows_remediation_script_registry_key_generator" "test" {
  context           = "all_users"
  registry_key_path = "Software\\MyApp\\Settings\\"
  value_name        = "EnableFeature"
  value_type        = "REG_SZ"
  value_data        = "Enabled"
}

