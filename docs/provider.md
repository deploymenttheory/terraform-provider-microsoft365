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
  cloud            = var.cloud
  tenant_id        = var.tenant_id
  auth_method      = var.auth_method
  telemetry_optout = var.telemetry_optout
  debug_mode       = var.debug_mode

  entra_id_options = {
    client_id                     = var.client_id
    client_secret                 = var.client_secret
    client_certificate            = var.client_certificate
    client_certificate_password   = var.client_certificate_password
    send_certificate_chain        = var.send_certificate_chain
    username                      = var.username
    password                      = var.password
    disable_instance_discovery    = var.disable_instance_discovery
    additionally_allowed_tenants  = var.additionally_allowed_tenants
    redirect_url                  = var.redirect_url
  }

  client_options = {
    enable_headers_inspection = var.enable_headers_inspection
    enable_retry              = var.enable_retry
    max_retries               = var.max_retries
    retry_delay_seconds       = var.retry_delay_seconds
    enable_redirect           = var.enable_redirect
    max_redirects             = var.max_redirects
    enable_compression        = var.enable_compression
    custom_user_agent         = var.custom_user_agent
    use_proxy                 = var.use_proxy
    proxy_url                 = var.proxy_url
    proxy_username            = var.proxy_username
    proxy_password            = var.proxy_password
    timeout_seconds           = var.timeout_seconds
    enable_chaos              = var.enable_chaos
    chaos_percentage          = var.chaos_percentage
    chaos_status_code         = var.chaos_status_code
    chaos_status_message      = var.chaos_status_message
  }
}

# Existing variables
variable "cloud" {
  description = "The cloud to use for authentication and Graph / Graph Beta API requests. Default is `public`. Valid values are `public`, `gcc`, `gcchigh`, `china`, `dod`, `ex`, `rx`. Can also be set using the `M365_CLOUD` environment variable."
  type        = string
  default     = "public"
}

variable "tenant_id" {
  description = "The M365 tenant ID for the Entra ID application. This ID uniquely identifies your Entra ID (EID) instance. It can be found in the Azure portal under Entra ID > Properties. Can also be set using the `M365_TENANT_ID` environment variable."
  type        = string
  default     = "2fd6bb84-ad40-4ec5-9369-a215b25c9952"
}

variable "auth_method" {
  description = "The authentication method to use for the Entra ID application to authenticate the provider. Options: 'device_code', 'client_secret', 'client_certificate', 'interactive_browser', 'username_password'. Can also be set using the `M365_AUTH_METHOD` environment variable."
  type        = string
  default     = "client_certificate"
}

variable "telemetry_optout" {
  description = "Flag to indicate whether to opt out of telemetry. Default is `false`. Can also be set using the `M365_TELEMETRY_OPTOUT` environment variable."
  type        = bool
  default     = false
}

variable "debug_mode" {
  description = "Flag to enable debug mode for the provider. When enabled, the provider will output additional debug information to the console to help diagnose issues. Can also be set using the `M365_DEBUG_MODE` environment variable."
  type        = bool
  default     = true
}

# EntraIDOptions variables
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

variable "send_certificate_chain" {
  description = "Controls whether the credential sends the public certificate chain in the x5c header of each token request's JWT. Can also be set using the `M365_SEND_CERTIFICATE_CHAIN` environment variable."
  type        = bool
  default     = false
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

variable "disable_instance_discovery" {
  description = "Disables the instance discovery in disconnected clouds or private clouds. Can also be set using the `M365_DISABLE_INSTANCE_DISCOVERY` environment variable."
  type        = bool
  default     = false
}

variable "additionally_allowed_tenants" {
  description = "Specifies additional tenants for which the credential may acquire tokens. Can also be set using the `M365_ADDITIONALLY_ALLOWED_TENANTS` environment variable."
  type        = list(string)
  default     = []
}

variable "redirect_url" {
  description = "The redirect URL for interactive browser authentication. Can also be set using the `M365_REDIRECT_URL` environment variable."
  type        = string
  default     = ""
}

# ClientOptions variables
variable "enable_headers_inspection" {
  description = "Enable inspection of HTTP headers. Can also be set using the `M365_ENABLE_HEADERS_INSPECTION` environment variable."
  type        = bool
  default     = false
}

variable "enable_retry" {
  description = "Enable automatic retries for failed requests. Can also be set using the `M365_ENABLE_RETRY` environment variable."
  type        = bool
  default     = true
}

variable "max_retries" {
  description = "Maximum number of retries for failed requests. Can also be set using the `M365_MAX_RETRIES` environment variable."
  type        = number
  default     = 3
}

variable "retry_delay_seconds" {
  description = "Delay between retry attempts in seconds. Can also be set using the `M365_RETRY_DELAY_SECONDS` environment variable."
  type        = number
  default     = 5
}

variable "enable_redirect" {
  description = "Enable automatic following of redirects. Can also be set using the `M365_ENABLE_REDIRECT` environment variable."
  type        = bool
  default     = true
}

variable "max_redirects" {
  description = "Maximum number of redirects to follow. Can also be set using the `M365_MAX_REDIRECTS` environment variable."
  type        = number
  default     = 5
}

variable "enable_compression" {
  description = "Enable compression for HTTP requests and responses. Can also be set using the `M365_ENABLE_COMPRESSION` environment variable."
  type        = bool
  default     = true
}

variable "custom_user_agent" {
  description = "Custom User-Agent string to be sent with requests. Can also be set using the `M365_CUSTOM_USER_AGENT` environment variable."
  type        = string
  default     = ""
}

variable "use_proxy" {
  description = "Enables the use of an HTTP proxy for network requests. Can also be set using the `M365_USE_PROXY` environment variable."
  type        = bool
  default     = false
}

variable "proxy_url" {
  description = "Specifies the URL of the HTTP proxy server. Can also be set using the `M365_PROXY_URL` environment variable."
  type        = string
  default     = ""
}

variable "proxy_username" {
  description = "Username for proxy authentication. Can also be set using the `M365_PROXY_USERNAME` environment variable."
  type        = string
  default     = ""
}

variable "proxy_password" {
  description = "Password for proxy authentication. Can also be set using the `M365_PROXY_PASSWORD` environment variable."
  type        = string
  sensitive   = true
  default     = ""
}

variable "timeout_seconds" {
  description = "Timeout for requests in seconds. Can also be set using the `M365_TIMEOUT_SECONDS` environment variable."
  type        = number
  default     = 300
}

variable "enable_chaos" {
  description = "Enable the chaos handler for testing purposes. Can also be set using the `M365_ENABLE_CHAOS` environment variable."
  type        = bool
  default     = false
}

variable "chaos_percentage" {
  description = "Percentage of requests to apply chaos testing to. Must be between 0 and 100. Can also be set using the `M365_CHAOS_PERCENTAGE` environment variable."
  type        = number
  default     = 10
}

variable "chaos_status_code" {
  description = "HTTP status code to return for chaos-affected requests. If not set, a random error status code will be used. Can also be set using the `M365_CHAOS_STATUS_CODE` environment variable."
  type        = number
  default     = 0
}

variable "chaos_status_message" {
  description = "Custom status message to return for chaos-affected requests. If not set, a default message will be used. Can also be set using the `M365_CHAOS_STATUS_MESSAGE` environment variable."
  type        = string
  default     = ""
}
```

## Schema

Please refer to the provider documentation for details on the configuration schema.