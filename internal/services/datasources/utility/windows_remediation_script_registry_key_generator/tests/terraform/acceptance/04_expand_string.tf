data "microsoft365_utility_windows_remediation_script_registry_key_generator" "test" {
  context           = "all_users"
  registry_key_path = "Software\\MyApp\\"
  value_name        = "InstallPath"
  value_type        = "REG_EXPAND_SZ"
  value_data        = "%ProgramFiles%\\MyApp"
}

