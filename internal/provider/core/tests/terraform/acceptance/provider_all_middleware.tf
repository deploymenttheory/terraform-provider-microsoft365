provider "microsoft365" {
  cloud       = "public"
  auth_method = "client_secret"
  client_options = {
    enable_retry              = true
    max_retries               = 3
    retry_delay_seconds       = 5
    enable_redirect           = true
    max_redirects             = 10
    enable_compression        = true
    enable_headers_inspection = true
    custom_user_agent         = "TestAgent/1.0"
    timeout_seconds           = 120
  }
}

data "microsoft365_graph_beta_device_management_windows_remediation_script" "test" {
  filter_type  = "display_name"
  filter_value = "NonExistentScript"
}
