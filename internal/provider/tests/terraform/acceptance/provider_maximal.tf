# Test provider with complete configuration
provider "microsoft365" {
  cloud            = "public"
  auth_method      = "client_secret"
  telemetry_optout = false
  debug_mode       = false

  entra_id_options = {
    client_id                    = "00000000-0000-0000-0000-000000000001"
    client_secret               = "test-client-secret"
    disable_instance_discovery   = false
    additionally_allowed_tenants = ["*"]
    redirect_url                = "http://localhost:8000/auth/callback"
  }

  client_options = {
    enable_headers_inspection = true
    enable_retry             = true
    max_retries             = 3
    retry_delay_seconds     = 1
    enable_redirect         = true
    max_redirects           = 5
    enable_compression      = true
    custom_user_agent       = "terraform-provider-microsoft365/test"
    use_proxy               = false
    timeout_seconds         = 60
    enable_chaos            = false
    chaos_percentage        = 10
    chaos_status_code       = 503
    chaos_status_message    = "Test chaos message"
  }
}