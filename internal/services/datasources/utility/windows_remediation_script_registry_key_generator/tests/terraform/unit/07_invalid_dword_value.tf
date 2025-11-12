data "microsoft365_windows_remediation_script_registry_key_generator" "test" {
  context           = "current_user"
  registry_key_path = "Software\\Test\\"
  value_name        = "TestValue"
  value_type        = "REG_DWORD"
  value_data        = "not_a_number"
}

