provider "microsoft365" {
  cloud       = "public"
  auth_method = "azure_developer_cli"
  client_options = {
    use_proxy = false
  }
}

data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type  = "display_name"
  filter_value = "NonExistentScript"
}
