data "microsoft365_utility_windows_remediation_script_registry_key_generator" "test" {
  context           = "current_user"
  registry_key_path = "Software\\MyCompany\\Settings\\"
  value_name        = "LargeValue"
  value_type        = "REG_QWORD"
  value_data        = "9223372036854775807"
}

