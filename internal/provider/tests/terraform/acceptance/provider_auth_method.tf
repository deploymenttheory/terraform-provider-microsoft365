# Test provider with specific authentication method
provider "microsoft365" {
  auth_method = "{{.AuthMethod}}"
}