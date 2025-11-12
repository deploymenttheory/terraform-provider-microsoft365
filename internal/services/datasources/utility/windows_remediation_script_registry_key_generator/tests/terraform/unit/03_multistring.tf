data "microsoft365_utility_windows_remediation_script_registry_key_generator" "test" {
  context           = "current_user"
  registry_key_path = "Software\\Test\\"
  value_name        = "TestMulti"
  value_type        = "REG_MULTI_SZ"
  value_data        = "Line1\nLine2\nLine3"
}

