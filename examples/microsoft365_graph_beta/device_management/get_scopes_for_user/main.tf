terraform {
  required_providers {
    microsoft365 = {
      source  = "deploymenttheory/microsoft365"
      version = ">= 1.0.0"
    }
  }
  required_version = ">= 1.0.0"
}

provider "microsoft365" {
  # Configuration options
}

# Get scopes for a specific user for a resource operation
data "microsoft365_graph_beta_device_management_get_scopes_for_user" "example" {
  resource_operation_id = "00000000-0000-0000-0000-000000000000" # Replace with an actual resource operation ID
  user_id               = "11111111-1111-1111-1111-111111111111" # Replace with an actual user ID
}

# Output the scopes
output "user_scopes" {
  value = data.microsoft365_graph_beta_device_management_get_scopes_for_user.example.scopes
} 