data "microsoft365_utility_windows_remediation_script_registry_key_generator" "test" {
  context           = "current_user"
  registry_key_path = "Software\\Test\\"
  value_name        = "TestBinary"
  value_type        = "REG_BINARY"
  value_data        = "01AF3C"
}

