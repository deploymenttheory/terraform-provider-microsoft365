# Test provider with client secret authentication
provider "microsoft365" {
  auth_method = "client_secret"
  
  entra_id_options = {
    client_id     = "00000000-0000-0000-0000-000000000001"
    client_secret = "test-client-secret"
  }
}