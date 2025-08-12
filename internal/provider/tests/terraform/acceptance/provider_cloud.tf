# Test provider with specific cloud environment  
provider "microsoft365" {
  cloud       = "{{.Cloud}}"
  auth_method = "device_code"
}