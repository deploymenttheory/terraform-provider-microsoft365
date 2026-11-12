provider "microsoft365" {
  cloud       = "public"
  auth_method = "client_secret"
  client_options = {
    enable_retry        = true
    max_retries         = 5
    retry_delay_seconds = 10
  }
}

data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type  = "display_name"
  filter_value = "NonExistentScript"
}
