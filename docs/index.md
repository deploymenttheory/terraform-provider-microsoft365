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


# Example resource
resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "example" {
  display_name = "Example Filter"
  description  = "This is an example filter"
  platform     = "windows10"
  rule         = "(device.manufacturer -eq \"Microsoft\")"
}

# Variables
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

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `auth_method` (String) The authentication method to use for the Entra ID application to authenticate the provider. Options:
- `device_code`: Uses a device code flow for authentication.
- `client_secret`: Uses a client ID and secret for authentication.
- `client_certificate`: Uses a client certificate (.pfx) for authentication.
- `interactive_browser`: Opens a browser for interactive login.
- `username_password`: Uses username and password for authentication (not recommended for production).
- `client_assertion`: Uses a client assertion (OIDC token) for authentication, suitable for CI/CD and server-to-server scenarios.

Each method requires different credentials to be provided. Can also be set using the `M365_AUTH_METHOD` environment variable.
- `cloud` (String) Specifies the Microsoft cloud environment for authentication and API requests. This setting determines the endpoints used for Microsoft Graph and Graph Beta APIs.

Valid values:
- `public`: Microsoft Azure Public Cloud (default)
- `dod`: US Department of Defense (DoD) Cloud
- `gcc`: US Government Cloud
- `gcchigh`: US Government High Cloud
- `china`: China Cloud
- `ex`: EagleX Cloud
- `rx`: Secure Cloud (RX)

Can be set using the `M365_CLOUD` environment variable.
- `tenant_id` (String, Sensitive) The Microsoft 365 tenant ID for the Entra ID (formerly Azure AD) application. This GUID uniquely identifies your Entra ID instance.

To find your tenant ID:
1. Log in to the [Azure portal](https://portal.azure.com)
2. Navigate to 'Microsoft Entra ID' (formerly Azure Active Directory)
3. In the Overview page, look for 'Tenant ID'

Alternatively, you can use PowerShell:
```powershell
(Get-AzureADTenantDetails).ObjectId
```

Or Azure CLI:
```bash
az account show --query tenantId -o tsv
```

Can be set using the `M365_TENANT_ID` environment variable.

### Optional

- `client_certificate` (String, Sensitive) The path to the client certificate file for certificate-based authentication with Entra ID (formerly Azure AD). This method is more secure than client secret-based authentication.

Requirements:
- File format: PKCS#12 (.pfx or .p12)
- Contents: Certificate, private key, and optionally a certificate chain
- The public key certificate must be uploaded to Entra ID

Steps to set up certificate authentication:
1. Generate a self-signed certificate or obtain one from a trusted Certificate Authority
2. Convert the certificate to PKCS#12 format if necessary
3. Upload the public key to Entra ID:
   - Go to Azure Portal > 'Microsoft Entra ID' > 'App registrations' > [Your App] > 'Certificates & secrets'
   - Upload the public key portion of your certificate
4. Provide the path to the PKCS#12 file in this attribute

Using OpenSSL to create a self-signed certificate:
```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365
openssl pkcs12 -export -out certificate.pfx -inkey key.pem -in cert.pem
```

Can be set using the `M365_CLIENT_CERTIFICATE_FILE_PATH` environment variable.
- `client_certificate_password` (String, Sensitive) The password to decrypt the PKCS#12 (.pfx or .p12) client certificate file. This is required only if the certificate file is password-protected.

Important notes:
- This password is used to encrypt the private key in the certificate file
- It's not related to any Entra ID settings, but to the certificate file itself
- If your PKCS#12 file was created without a password, leave this field empty
- Treat this password with the same level of security as the certificate itself

When creating a PKCS#12 file with OpenSSL, you'll be prompted for this password:
```bash
openssl pkcs12 -export -out certificate.pfx -inkey key.pem -in cert.pem
```

Can be set using the `M365_CLIENT_CERTIFICATE_PASSWORD` environment variable.
- `client_id` (String, Sensitive) The client ID (application ID) for the Entra ID (formerly Azure AD) application. This GUID is generated when you register an application in Entra ID.

To find or create a client ID:
1. Log in to the [Azure portal](https://portal.azure.com)
2. Navigate to 'Microsoft Entra ID' > 'App registrations'
3. Select your application or create a new one
4. The client ID is listed as 'Application (client) ID' in the Overview page

Using Azure CLI:
```bash
az ad app list --query "[].{appId:appId, displayName:displayName}"
```

Using Microsoft Graph PowerShell:
```powershell
Get-MgApplication -Filter "displayName eq 'Your App Name'" | Select-Object AppId, DisplayName
```

Can be set using the `M365_CLIENT_ID` environment variable.
- `client_secret` (String, Sensitive) The client secret for the Entra ID (formerly Azure AD) application. This secret is required for client credentials authentication flow.

Important notes:
- Client secrets are sensitive and should be handled securely
- Secrets have an expiration date and need to be rotated periodically
- Use managed identities or certificate-based authentication when possible for improved security

To create a client secret:
1. Log in to the [Azure portal](https://portal.azure.com)
2. Navigate to 'Microsoft Entra ID' > 'App registrations'
3. Select your application
4. Go to 'Certificates & secrets' > 'Client secrets'
5. Click 'New client secret' and set a description and expiration
6. Copy the secret value immediately (it won't be shown again)

Using Azure CLI:
```bash
az ad app credential reset --id <app-id> --append
```

Using Microsoft Graph PowerShell:
```powershell
$credential = @{
    displayName = 'My Secret'
    endDateTime = (Get-Date).AddMonths(6)
}
New-MgApplicationPassword -ApplicationId <app-id> -PasswordCredential $credential
```

Can be set using the `M365_CLIENT_SECRET` environment variable.
- `debug_mode` (Boolean) Flag to enable debug mode for the provider.Can also be set using the `M365_DEBUG_MODE` environment variable.
- `enable_chaos` (Boolean) Enable the chaos handler for testing purposes. When enabled, the chaos handler can simulate specific failure scenarios and random errors in API responses to help test the robustness and resilience of the terraform provider against intermittent issues. This is particularly useful for testing how the provider handles various error conditions and ensures it can recover gracefully. Use with caution in production environments. Can also be set using the `M365_ENABLE_CHAOS` environment variable.
- `password` (String, Sensitive) The password for resource owner password credentials (ROPC) flow authentication.

**Critical Security Warning:**
- Storing passwords in plain text is a significant security risk
- Use environment variables or secure vaults to manage this sensitive information
- Regularly rotate passwords and monitor for unauthorized access
- Consider using more secure authentication methods when possible

Can be set using the `M365_PASSWORD` environment variable.
- `proxy_url` (String) Specifies the URL of the HTTP proxy server for routing requests when `use_proxy` is enabled.

**Format:**
- Must be a valid URL (e.g., `http://proxy.example.com:8080`)
- Supports HTTP and HTTPS protocols

**Usage:**
- When `use_proxy` is set to `true`, all HTTP(S) requests will be routed through this proxy
- Ensure the proxy server is reachable and correctly configured to handle the traffic

**Examples:**
- HTTP proxy: `http://proxy.example.com:8080`
- HTTPS proxy: `https://secure-proxy.example.com:443`
- Authenticated proxy: `http://username:password@proxy.example.com:8080`

**Security Considerations:**
- Use HTTPS for the proxy URL when possible to encrypt proxy communications
- If using an authenticated proxy, consider setting the URL via the environment variable to avoid exposing credentials in configuration files
- Ensure the proxy server is trusted and secure

Can be set using the `M365_PROXY_URL` environment variable.
- `redirect_url` (String) The redirect URL (also known as reply URL or callback URL) for OAuth 2.0 authentication flows that require a callback, such as the Authorization Code flow or interactive browser authentication.

**Important:**
- This URL must be registered in your Entra ID (formerly Azure AD) application settings
- For local development, typically use `http://localhost:port`
- For production, use a secure HTTPS URL

To configure in Entra ID:
1. Go to Azure Portal > 'Microsoft Entra ID' > 'App registrations'
2. Select your application
3. Go to 'Authentication' > 'Platform configurations'
4. Add or update the redirect URI

Security considerations:
- Use a unique path for your redirect URL to prevent potential conflicts
- Avoid using wildcard URLs in production environments
- Regularly audit and remove any unused redirect URLs

Example values:
- Local development: `http://localhost:8000/auth/callback`
- Production: `https://yourdomain.com/auth/microsoft365/callback`

Can be set using the `M365_REDIRECT_URL` environment variable.
- `telemetry_optout` (Boolean) Controls the collection of telemetry data for the Microsoft 365 provider by Microsoft Services.

**Usage:**
- Set to `true` to disable all telemetry collection
- Set to `false` (default) to allow telemetry collection

**Behavior:**
- When set to `true`, it prevents the addition of any telemetry data to API requests
- This affects the User-Agent string and other potential telemetry mechanisms

**Privacy:**
- Telemetry, when enabled, may include provider version, Terraform version, and general usage patterns
- No personally identifiable information (PII) or sensitive data is collected

**Recommendations:**
- For development or non-sensitive environments, consider leaving telemetry enabled to support product improvement
- For production or sensitive environments, you may choose to opt out

Can be set using the `M365_TELEMETRY_OPTOUT` environment variable.
- `use_proxy` (Boolean) Enables the use of an HTTP proxy for network requests. When set to true, the provider will route all HTTP requests through the specified proxy server. This can be useful for environments that require proxy access for internet connectivity or for monitoring and logging HTTP traffic. Can also be set using the `M365_USE_PROXY` environment variable.
- `username` (String) The username for resource owner password credentials (ROPC) flow authentication.

**Important Security Notice:**
- Resource Owner Password Credentials (ROPC) is considered less secure than other authentication methods
- It should only be used when other, more secure methods are not possible
- Not recommended for production environments
- Does not support multi-factor authentication

Usage:
- Typically, this is the user's email address or User Principal Name (UPN)
- Ensure the user has appropriate permissions for the required operations

Can be set using the `M365_USERNAME` environment variable.