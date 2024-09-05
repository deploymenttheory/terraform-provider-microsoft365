---
page_title: "Provider: m365"
description: |-
  The m365 provider is used to manage m365 resources.  
---

# m365 Provider

The Terraform m365 msgraph provider is a plugin for Terraform that allows for the
management of [m365](https://github.com/microsoftgraph/msgraph-metadata) resources.

## Example Usage

```terraform
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
  tenant_id                   = var.tenant_id
  auth_method                 = var.auth_method
  client_id                   = var.client_id
  client_secret               = var.client_secret
  client_certificate          = var.client_certificate
  client_certificate_password = var.client_certificate_password
  username                    = var.username
  password                    = var.password
  redirect_url                = var.redirect_url
  use_proxy                   = var.use_proxy
  proxy_url                   = var.proxy_url
  cloud                       = var.cloud
  enable_chaos                = var.enable_chaos
  telemetry_optout            = var.telemetry_optout
  debug_mode                  = var.debug_mode
}

variable "cloud" {
  description = "The cloud to use for authentication and Graph / Graph Beta API requests. Default is `public`. Valid values are `public`, `gcc`, `gcchigh`, `china`, `dod`, `ex`, `rx`. Can also be set using the `M365_CLOUD` environment variable."
  type        = string
  default     = "public"
}

variable "tenant_id" {
  description = "The M365 tenant ID for the Entra ID application. This ID uniquely identifies your Entra ID (EID) instance. It can be found in the Azure portal under Entra ID > Properties. Can also be set using the `M365_TENANT_ID` environment variable."
  type        = string
  default     = ""
}

variable "auth_method" {
  description = "The authentication method to use for the Entra ID application to authenticate the provider. Options: 'device_code', 'client_secret', 'client_certificate', 'interactive_browser', 'username_password'. Can also be set using the `M365_AUTH_METHOD` environment variable."
  type        = string
  default     = "client_secret"
}

variable "client_id" {
  description = "The client ID for the Entra ID application. This ID is generated when you register an application in the Entra ID (Azure AD) and can be found under App registrations > YourApp > Overview. Can also be set using the `M365_CLIENT_ID` environment variable."
  type        = string
  default     = ""
}

variable "client_secret" {
  description = "The client secret for the Entra ID application. This secret is generated in the Entra ID (Azure AD) and is required for authentication flows such as client credentials and on-behalf-of flows. It can be found under App registrations > YourApp > Certificates & secrets. Required for client credentials and on-behalf-of flows. Can also be set using the `M365_CLIENT_SECRET` environment variable."
  type        = string
  sensitive   = true
  default     = ""
}

variable "client_certificate" {
  description = "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate. Can also be set using the `M365_CLIENT_CERTIFICATE_FILE_PATH` environment variable."
  type        = string
  sensitive   = true
  default     = ""
}

variable "client_certificate_password" {
  description = "The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate. Can also be set using the `M365_CLIENT_CERTIFICATE_PASSWORD` environment variable."
  type        = string
  sensitive   = true
  default     = ""
}

variable "username" {
  description = "The username for username/password authentication. Can also be set using the `M365_USERNAME` environment variable."
  type        = string
  default     = ""
}

variable "password" {
  description = "The password for username/password authentication. Can also be set using the `M365_PASSWORD` environment variable."
  type        = string
  sensitive   = true
  default     = ""
}

variable "redirect_url" {
  description = "The redirect URL for interactive browser authentication. Can also be set using the `M365_REDIRECT_URL` environment variable."
  type        = string
  default     = ""
}

variable "use_proxy" {
  description = "Enables the use of an HTTP proxy for network requests. When set to true, the provider will route all HTTP requests through the specified proxy server. This can be useful for environments that require proxy access for internet connectivity or for monitoring and logging HTTP traffic. Can also be set using the `M365_USE_PROXY` environment variable."
  type        = bool
  default     = false
}

variable "proxy_url" {
  description = "Specifies the URL of the HTTP proxy server. This URL should be in a valid URL format (e.g., 'http://proxy.example.com:8080'). When 'use_proxy' is enabled, this URL is used to configure the HTTP client to route requests through the proxy. Ensure the proxy server is reachable and correctly configured to handle the network traffic. Can also be set using the `M365_PROXY_URL` environment variable."
  type        = string
  default     = ""
}

variable "enable_chaos" {
  description = "Enable the chaos handler for testing purposes. When enabled, the chaos handler can simulate specific failure scenarios and random errors in API responses to help test the robustness and resilience of the terraform provider against intermittent issues. This is particularly useful for testing how the provider handles various error conditions and ensures it can recover gracefully. Use with caution in production environments. Can also be set using the `M365_ENABLE_CHAOS` environment variable."
  type        = bool
  default     = false
}

variable "telemetry_optout" {
  description = "Flag to indicate whether to opt out of telemetry. Default is `false`. Can also be set using the `M365_TELEMETRY_OPTOUT` environment variable."
  type        = bool
  default     = false
}

variable "debug_mode" {
  description = "Flag to enable debug mode for the provider. When enabled, the provider will output additional debug information to the console to help diagnose issues. Can also be set using the `M365_DEBUG_MODE` environment variable."
  type        = bool
  default     = false
}
``` 

