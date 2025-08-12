# Test that environment variables take precedence over HCL configuration
# Environment variables should override these values:
# M365_CLOUD=gcc (overrides "public")
# M365_AUTH_METHOD=device_code (overrides "client_secret") 
# M365_DEBUG_MODE=true (overrides false)
# M365_TELEMETRY_OPTOUT=true (overrides false)

provider "microsoft365" {
  cloud            = "public"
  auth_method      = "client_secret"
  debug_mode       = false
  telemetry_optout = false
  
  entra_id_options = {
    client_id     = "hcl-client-id"
    client_secret = "hcl-client-secret"
  }
}