# Test provider with client certificate authentication
provider "microsoft365" {
  auth_method = "client_certificate"
  
  entra_id_options = {
    client_id                    = "00000000-0000-0000-0000-000000000001"
    client_certificate          = "/path/to/cert.pfx"
    client_certificate_password = "cert-password"
    send_certificate_chain      = true
  }
}