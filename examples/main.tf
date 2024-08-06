# Example backend
terraform {
  required_providers {
    microsoft365 = {
      source  = "deploymenttheory/terraform-provider-microsoft365"
      version = "~> 1.0.0"  
    }
  }
}

# Example provider
provider "microsoft365" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
  auth_method   = "client_secret"
  cloud         = "public"
}

# Example resource
resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "example" {
  display_name = "Example Filter"
  description  = "This is an example filter"
  platform     = "windows10"
  rule         = "(device.manufacturer -eq \"Microsoft\")"
}

# Variables
variable "tenant_id" {
  description = "The Microsoft 365 tenant ID"
  type        = string
}

variable "client_id" {
  description = "The client ID for the Entra ID application"
  type        = string
}

variable "client_secret" {
  description = "The client secret for the Entra ID application"
  type        = string
  sensitive   = true
}