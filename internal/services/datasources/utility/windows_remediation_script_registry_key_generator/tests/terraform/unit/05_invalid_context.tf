data "microsoft365_utility_windows_remediation_script_registry_key_generator" "test" {
  context           = "invalid_context"
  registry_key_path = "Software\\Test\\"
  value_name        = "TestValue"
  value_type        = "REG_DWORD"
  value_data        = "1"
}

