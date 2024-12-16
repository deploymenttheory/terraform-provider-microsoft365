# Microsoft365 Provider Configuration with Dummy Values

terraform {
  required_providers {
    microsoft365 = {
      source  = "deploymenttheory/microsoft365"
      version = "~> 1.0" # Replace with the actual version you're using
    }
  }
}

provider "microsoft365" {
  cloud                       = "public"
  tenant_id                   = "11111111-1111-1111-1111-111111111111"
  auth_method                 = "client_secret"
  client_id                   = "22222222-2222-2222-2222-222222222222"
  client_secret               = "dummyClientSecret123!@#"
  client_certificate          = "/path/to/dummy/certificate.pfx"
  client_certificate_password = "dummyCertPassword456!@#"
  username                    = "dummyuser@example.com"
  password                    = "dummyPassword789!@#"
  redirect_url                = "http://localhost:8080/oauth/callback"
  use_proxy                   = true
  proxy_url                   = "http://dummy-proxy.example.com:8080"
  enable_chaos                = false
  telemetry_optout            = false
  debug_mode                  = true
}