provider "microsoft365" {
  cloud       = "public"
  auth_method = "client_secret"

  entra_id_options = {
    client_id     = "00000000-0000-0000-0000-000000000001"
    client_secret = "test-secret"
  }

  client_options = {
    enable_retry    = true
    max_retries     = 3
    enable_redirect = true
    use_proxy       = false
    enable_chaos    = false
    timeout_seconds = 60
  }

  telemetry_optout = false
  debug_mode       = false
}