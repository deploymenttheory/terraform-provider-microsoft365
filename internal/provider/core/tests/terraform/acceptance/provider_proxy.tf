# Test provider with proxy configuration
provider "microsoft365" {
  auth_method = "device_code"

  client_options = {
    use_proxy      = true
    proxy_url      = "http://proxy.example.com:8080"
    proxy_username = "proxy-user"
    proxy_password = "proxy-password"
  }
}