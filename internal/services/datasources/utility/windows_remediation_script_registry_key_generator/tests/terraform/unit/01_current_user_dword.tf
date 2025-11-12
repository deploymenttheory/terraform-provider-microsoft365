data "microsoft365_utility_windows_remediation_script_registry_key_generator" "test" {
  context           = "current_user"
  registry_key_path = "Software\\Policies\\Microsoft\\WindowsStore\\"
  value_name        = "RequirePrivateStoreOnly"
  value_type        = "REG_DWORD"
  value_data        = "1"
}

