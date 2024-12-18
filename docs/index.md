---
page_title: "Provider: Microsoft 365"
description: |-
  
---

# terraform-provider-microsoft365 Provider

The community Microsoft 365 provider allows managing environments and other resources within [Microsoft 365](https://www.microsoft.com/en-gb/microsoft-365/products-apps-services).

!> This code is made available as a public preview. Features are being actively developed and may have restricted or limited functionality. Future updates may introduce breaking changes, but we follow [Semantic Versioning](https://semver.org/) to help mitigate this. The software may contain bugs, errors, or other issues that could cause service interruption or data loss. We recommend backing up your data and testing in non-production environments. Your feedback is valuable to us, so please share any issues or suggestions you encounter via GitHub issues.

## Requirements

This provider requires **Terraform >= 0.12**.  For more information on provider installation and constraining provider versions, see the [Provider Requirements documentation](https://developer.hashicorp.com/terraform/language/providers/requirements).

## Installation

To use this provider, add the following to your Terraform configuration:

```terraform
terraform {
  required_providers {
    microsoft365 = {
      source  = "deploymenttheory/microsoft365"
      version = "~> 0.4.0" # Replace with the latest version
    }
  }
}
```

See the official Terraform documentation for more information about [requiring providers](https://developer.hashicorp.com/terraform/language/providers/requirements).

# Authenticating to Microsoft 365

This Terraform provider supports multiple authentication methods for accessing Microsoft 365 services:

* [Authenticating using Client Secret](#authenticating-using-client-secret)
* [Authenticating using Client Certificate](#authenticating-using-client-certificate)
* [Authenticating using Username and Password](#authenticating-using-username-and-password)
* [Authenticating using Device Code](#authenticating-using-device-code)
* [Authenticating using Interactive Browser](#authenticating-using-interactive-browser)

## Authenticating using Client Secret

The Microsoft 365 provider can use a Service Principal with Client Secret to authenticate to Microsoft 365 services.

1. [Create an app registration in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app)
2. Add the required API permissions for Microsoft Graph
3. Create a client secret in the app registration
4. Configure the provider to use a Service Principal with a Client Secret:

```terraform
provider "microsoft365" {
  auth_method = "client_secret"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id     = "00000000-0000-0000-0000-000000000000"
    client_secret = "your-client-secret"
  }
}
```

## Authenticating using Client Certificate

The Microsoft 365 provider can use certificate-based authentication for enhanced security.

1. [Create an app registration in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app)
2. Add the required API permissions for Microsoft Graph
3. Generate a certificate using openssl or other tools:

```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365
```

4. Merge public and private parts of the certificate files:

```bash
# Using Linux shell
cat *.pem > cert+key.pem

# Using PowerShell
Get-Content .\cert.pem, .\key.pem | Set-Content cert+key.pem
```

5. Generate pkcs12 file:

```bash
openssl pkcs12 -export -out cert.pkcs12 -in cert+key.pem
```

6. Upload the public certificate (`cert.pem`) to your app registration
7. Configure the provider:

```terraform
provider "microsoft365" {
  auth_method = "client_certificate"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id                    = "00000000-0000-0000-0000-000000000000"
    client_certificate          = "${path.cwd}/cert.pkcs12"
    client_certificate_password = "your-certificate-password"
  }
}
```

## Authenticating using Username and Password

The Microsoft 365 provider can authenticate using standard username and password credentials.

1. [Create an app registration in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app)
2. Add the required API permissions for Microsoft Graph
3. Configure the provider:

```terraform
provider "microsoft365" {
  auth_method = "username_password"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id = "00000000-0000-0000-0000-000000000000"
    username  = "user@domain.com"
    password  = "your-password"
  }
}
```

## Authenticating using Device Code

The Microsoft 365 provider can use device code authentication when interactive login isn't possible.

1. [Create an app registration in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app)
2. Add the required API permissions for Microsoft Graph
3. Configure the provider:

```terraform
provider "microsoft365" {
  auth_method = "device_code"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id = "00000000-0000-0000-0000-000000000000"
  }
}
```

## Authenticating using Interactive Browser

The Microsoft 365 provider can authenticate using an interactive browser session.

1. [Create an app registration in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app)
2. Add the required API permissions for Microsoft Graph
3. Configure a redirect URI in your app registration
4. Configure the provider:

```terraform
provider "microsoft365" {
  auth_method = "interactive_browser"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id     = "00000000-0000-0000-0000-000000000000"
    redirect_url = "http://localhost:8888"
    username     = "user@domain.com"  # Optional login hint
  }
}
```

## Using Environment Variables

We recommend using Environment Variables to pass the credentials to the provider.

| Name | Description |
|------|-------------|
| `M365_TENANT_ID` | The Microsoft Entra ID tenant ID |
| `M365_AUTH_METHOD` | The authentication method to use |
| `M365_CLIENT_ID` | The application (client) ID |
| `M365_CLIENT_SECRET` | The client secret value |
| `M365_CLIENT_CERTIFICATE_FILE_PATH` | Path to the certificate file |
| `M365_CLIENT_CERTIFICATE_PASSWORD` | Password for the certificate |
| `M365_USERNAME` | Username for password or browser authentication |
| `M365_PASSWORD` | Password for password authentication |
| `M365_REDIRECT_URI` | Redirect URI for interactive browser authentication |
| `M365_CLOUD` | Cloud environment (defaults to global) |
| `M365_DISABLE_INSTANCE_DISCOVERY` | Disable instance discovery |
| `M365_ADDITIONALLY_ALLOWED_TENANTS` | List of additionally allowed tenant IDs |

-> Variables passed into the provider block will override the environment variables.

## Additional Provider Configuration

The provider supports additional configuration options for client behavior, telemetry, and debugging:

```terraform
provider "microsoft365" {
  # ... authentication configuration ...
  
  debug_mode = false             # ENV: M365_DEBUG_MODE
  telemetry_optout = false       # ENV: M365_TELEMETRY_OPTOUT
  
  client_options = {
    # ... client configuration options ...
  }
}
```

For a complete list of client options, refer to the provider documentation.

> **Security Note**: Store sensitive values like client secrets, certificates, and passwords using environment variables or Terraform's encrypted state management features. Never commit these values directly in your configuration files.

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
    client_id                    = var.client_id
    client_secret                = var.client_secret
    client_certificate           = var.client_certificate
    client_certificate_password  = var.client_certificate_password
    send_certificate_chain       = var.send_certificate_chain
    username                     = var.username
    password                     = var.password
    disable_instance_discovery   = var.disable_instance_discovery
    additionally_allowed_tenants = var.additionally_allowed_tenants
    redirect_url                 = var.redirect_url
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
  default     = ""
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

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `auth_method` (String) The authentication method to use for the Entra ID application to authenticate the provider. Options:
- `device_code`: Uses a device code flow for authentication.
- `client_secret`: Uses a client ID and secret for authentication.
- `client_certificate`: Uses a client certificate (.pfx) for authentication.
- `interactive_browser`: Opens a browser for interactive login.
- `username_password`: Uses username and password for authentication (not recommended for production).
Each method requires different credentials to be provided.
Can also be set using the `M365_AUTH_METHOD` environment variable.
- `cloud` (String) Specifies the Microsoft cloud environment for authentication and API requests. This setting determines the endpoints used for Microsoft Graph and Graph Beta APIs. Valid values:
- `public`: Microsoft Azure Public Cloud (default)
- `dod`: US Department of Defense (DoD) Cloud
- `gcc`: US Government Cloud
- `gcchigh`: US Government High Cloud
- `china`: China Cloud
- `ex`: EagleX Cloud
- `rx`: Secure Cloud (RX)

Can be set using the `M365_CLOUD` environment variable.
- `tenant_id` (String, Sensitive) The Microsoft 365 tenant ID for the Entra ID (formerly Azure AD) application. This GUID uniquely identifies your Entra ID instance.Can be set using the `M365_TENANT_ID` environment variable.

To find your tenant ID:
1. Log in to the [Azure portal](https://portal.azure.com)
2. Navigate to 'Microsoft Entra ID' (formerly Azure Active Directory)
3. In the Overview page, look for 'Tenant ID'

Alternatively, you can use PowerShell:
```powershell
Connect-AzAccount
(Get-AzContext).Tenant.Id
```

Can also be set using the `M365_TENANT_ID` environment variable.

### Optional

- `client_options` (Attributes) Configuration options for the Microsoft Graph client. (see [below for nested schema](#nestedatt--client_options))
- `debug_mode` (Boolean) Flag to enable debug mode for the provider.

This setting enables additional logging and diagnostics for the provider.

Can also be set using the `M365_DEBUG_MODE` environment variable.
- `entra_id_options` (Attributes) Configuration options for Entra ID authentication. (see [below for nested schema](#nestedatt--entra_id_options))
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

<a id="nestedatt--client_options"></a>
### Nested Schema for `client_options`

Optional:

- `chaos_percentage` (Number) Specifies the percentage of requests that should be affected by the chaos handler.

**Key points:**
- Value range: 0 to 100
- Default: 10% if not specified
- Determines the probability of injecting chaos into each request
- Higher values increase the frequency of simulated errors

**Example usage:**
```hcl
provider "microsoft365" {
  enable_chaos = true
  chaos_percentage = 30  # 30% of requests will be affected
}
```

Can be set using the `M365_CHAOS_PERCENTAGE` environment variable.
- `chaos_status_code` (Number) Specifies the HTTP status code to be returned for requests affected by the chaos handler.

**Key points:**
- If not set, a random error status code will be used
- Common error codes: 429 (Too Many Requests), 500 (Internal Server Error), 503 (Service Unavailable)
- Used only when `enable_chaos` is true

**Example usage:**
```hcl
provider "microsoft365" {
  enable_chaos = true
  chaos_status_code = 503  # Simulate a 'Service Unavailable' error
}
```

Can be set using the `M365_CHAOS_STATUS_CODE` environment variable.
- `chaos_status_message` (String) Defines a custom status message to be returned for requests affected by the chaos handler.

**Key points:**
- If not set, a default message will be used
- Allows simulation of specific error messages or conditions
- Used only when `enable_chaos` is true

**Example usage:**
```hcl
provider "microsoft365" {
  enable_chaos = true
  chaos_status_message = "Simulated server overload"
}
```

Can be set using the `M365_CHAOS_STATUS_MESSAGE` environment variable.
- `custom_user_agent` (String) Custom User-Agent string to be sent with requests.
- `enable_chaos` (Boolean) Enable the chaos handler for testing purposes. When enabled, the chaos handler simulates specific failure scenarios and random errors in API responses to help test the robustness and resilience of the terraform provider against intermittent issues. This is particularly useful for testing how the provider handles various error conditions and ensures it can recover gracefully.

**Key points:**
- Default: `false`
- When `true`, adds a chaos handler to the middleware
- Injects an 'X-Chaos-Injected: true' header in affected responses
- Use with caution, especially in production environments

**Example usage:**
```hcl
provider "microsoft365" {
  enable_chaos = true
  chaos_percentage = 20
}
```

Can also be set using the `M365_ENABLE_CHAOS` environment variable.
- `enable_compression` (Boolean) Enable compression for HTTP requests and responses.
- `enable_headers_inspection` (Boolean) Enable inspection of HTTP headers.
- `enable_redirect` (Boolean) Enable automatic following of redirects.
- `enable_retry` (Boolean) Enable automatic retries for failed requests.
- `max_redirects` (Number) Maximum number of redirects to follow.
- `max_retries` (Number) Maximum number of retries for failed requests.
- `proxy_password` (String, Sensitive) Specifies the password for authentication with the proxy server if required.

**Key points:**
- Optional: Only needed if your proxy server requires authentication
- Used in conjunction with `proxy_username`
- Treated as sensitive information and will be masked in logs
- Ignored if `use_proxy` is `false` or if `proxy_url` is not set

**Security note:** It's recommended to set this using an environment variable rather than in the configuration file.

**Example usage:**
```hcl
provider "microsoft365" {
  use_proxy      = true
  proxy_url      = "http://proxy.example.com:8080"
  proxy_username = "proxyuser"
  proxy_password = "proxypass"
}
```

Can be set using the `M365_PROXY_PASSWORD` environment variable.
- `proxy_url` (String) Specifies the URL of the proxy server to be used when `use_proxy` is set to `true`.

**Key points:**
- Must be a valid URL including the scheme (http:// or https://)
- Can include a port number
- Required when `use_proxy` is `true`
- Ignored if `use_proxy` is `false`

**Example usage:**
```hcl
provider "microsoft365" {
  use_proxy = true
  proxy_url = "http://proxy.example.com:8080"
}
```

Can be set using the `M365_PROXY_URL` environment variable.
- `proxy_username` (String) Specifies the username for authentication with the proxy server if required.

**Key points:**
- Optional: Only needed if your proxy server requires authentication
- Used in conjunction with `proxy_password`
- Ignored if `use_proxy` is `false` or if `proxy_url` is not set

**Example usage:**
```hcl
provider "microsoft365" {
  use_proxy      = true
  proxy_url      = "http://proxy.example.com:8080"
  proxy_username = "proxyuser"
  proxy_password = "proxypass"
}
```

Can be set using the `M365_PROXY_USERNAME` environment variable.
- `retry_delay_seconds` (Number) Delay between retry attempts in seconds.
- `timeout_seconds` (Number) Override value for authentication request timeouts in seconds.
- `use_proxy` (Boolean) Enables the use of a proxy server for all network requests made by the provider.

**Key points:**
- Default: `false`
- When `true`, the provider will route all HTTP requests through the specified proxy server
- Requires `proxy_url` to be set
- Useful for environments that require proxy access for internet connectivity

**Example usage:**
```hcl
provider "microsoft365" {
  use_proxy = true
  proxy_url = "http://proxy.example.com:8080"
}
```

Can be set using the `M365_USE_PROXY` environment variable.


<a id="nestedatt--entra_id_options"></a>
### Nested Schema for `entra_id_options`

Optional:

- `additionally_allowed_tenants` (List of String) Specifies additional tenants for which the credential may acquire tokens.Add the wildcard value '*' to allow the credential to acquire tokens for any tenant.

Can be set using the `M365_ADDITIONALLY_ALLOWED_TENANTS` environment variable.
- `client_certificate` (String, Sensitive) Used for the 'client_certificate' authentication method.

The path to the Client Certificate file associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate. Supports PKCS#12 (.pfx or .p12) file format. The file should contain the certificate, private key with an RSA type, and optionally a password which can be defined in client_certificate_password.

The path to the client certificate file for certificate-based authentication with Entra ID (formerly Azure AD). This method is more secure than client secret-based authentication.

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

**Example usage:**
```hcl
provider "microsoft365" {
  client_certificate        = "/path/to/cert.pfx"
}
```

Can be set using the `M365_CLIENT_CERTIFICATE_FILE_PATH` environment variable.
- `client_certificate_password` (String, Sensitive) Used for the 'client_certificate' authentication method.

The password to decrypt the PKCS#12 (.pfx or .p12) client certificate file. Required only if the certificate file is password-protected.

Important notes:
- This password is used to encrypt the private key in the certificate file
- It's not related to any Entra ID settings, but to the certificate file itself
- If your PKCS#12 file was created without a password, leave this field empty
- Treat this password with the same level of security as the certificate itself

When creating a PKCS#12 file with OpenSSL, you'll be prompted for this password:
```bash
openssl pkcs12 -export -out certificate.pfx -inkey key.pem -in cert.pem
```

**Example usage:**
```hcl
provider "microsoft365" {
  client_certificate_password = "certpassword"
}
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

**Example usage:**
```hcl
provider "microsoft365" {
  client_id = "my_client_id"
}
```

Can be set using the `M365_CLIENT_ID` environment variable.
- `client_secret` (String, Sensitive) Used for the 'client_secret' authentication method.

The client secret for the Entra ID application. Required for client credentials authentication. This secret is generated in Entra ID and has an expiration date.

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

**Example usage:**
```hcl
provider "microsoft365" {
  client_secret = "my_client_secret"
}
```

Can be set using the `M365_CLIENT_SECRET` environment variable.
- `disable_instance_discovery` (Boolean) DisableInstanceDiscovery should be set true only by terraform runsauthenticating in disconnected clouds, or private clouds such as Azure Stack.It determines whether the credential requests Microsoft Entra instance metadatafrom https://login.microsoft.com before authenticating. Setting this to true willskip this request, making the application responsible for ensuring the configuredauthority is valid and trustworthy.

Can be set using the `M365_DISABLE_INSTANCE_DISCOVERY` environment variable.
- `password` (String, Sensitive) Used for the 'username_password' authentication method.

The password for resource owner password credentials (ROPC) flow authentication.

**Critical Security Warning:**
- Storing passwords in plain text is a significant security risk
- Use environment variables or secure vaults to manage this sensitive information
- Regularly rotate passwords and monitor for unauthorized access
- Consider using more secure authentication methods when possible

Can be set using the `M365_PASSWORD` environment variable.
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
- `send_certificate_chain` (Boolean) Used for the 'client_certificate' authentication method.

Controls whether the credential sends the public certificate chain in the x5c headerof each token request's JWT. This is required for Subject Name/Issuer (SNI) authenticationand can be used in certain advanced scenarios. Defaults to false. Enable this if yourapplication uses certificate chain validation or if specifically instructed by Azure support.

**Key points:**
- Default value: `false`
- Set to `true` to enable sending the certificate chain

**Use cases:**
1. **Subject Name/Issuer (SNI) Authentication:** Required for SNI authentication scenarios.
2. **Enhanced Security:** Provides additional validation capabilities on the Entra ID side.
3. **Compatibility:** May be necessary for certain Azure services or configurations.

**How it works:**
- When enabled, the full X.509 certificate chain is included in the token request.
- This allows the Microsoft Entra ID (Azure AD) to perform additional validation on the certificate.
- It can help in scenarios where intermediate certificates need to be verified.

**Considerations:**
- Enabling this option increases the size of each token request.
- Only enable if you're certain your scenario requires it.
- Consult Azure documentation or support if you're unsure about enabling this option.

**Example usage:**
```hcl
provider "microsoft365" {
  client_certificate        = "/path/to/cert.pfx"
  client_certificate_password = "certpassword"
  send_certificate_chain    = true
}
```

Only enable this option if you understand its implications or if specifically instructed by Azure support.
- `username` (String) Used for the 'username_password' authentication method.

The username for resource owner password credentials (ROPC) flow authentication.

**Important Security Notice:**
- Resource Owner Password Credentials (ROPC) is considered less secure than other authentication methods
- It should only be used when other, more secure methods are not possible
- Not recommended for production environments
- Does not support multi-factor authentication

Usage:
- Typically, this is the user's email address or User Principal Name (UPN)
- Ensure the user has appropriate permissions for the required operations

**Example usage:**
```hcl
provider "microsoft365" {
  username        = "user_name
}
```

Can be set using the `M365_USERNAME` environment variable.


# Resources and Data Sources

Use the navigation to the left to read about the available resources and data sources.

!> By calling `terraform destroy` all the resources that you've created will be permanently deleted. Please be careful with this command when working with production environments. You can use [prevent-destroy](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#prevent_destroy) lifecycle argument in your resources to prevent accidental deletion.  

## Examples

You can find practical examples of using this provider in our examples directory. These examples demonstrate:
- Basic resource management
- Complex configurations
- Integration with Microsoft Graph API
- Best practices for production deployments

For more examples and use cases, visit our [Microsoft 365 Provider Examples](https://github.com/deploymenttheory/terraform-provider-microsoft365/tree/main/examples) directory.

## Releases

A full list of released versions of the Microsoft 365 Terraform Provider can be found in our [GitHub Releases](https://github.com/deploymenttheory/terraform-provider-microsoft365/releases).

Starting from the initial release, changes to the provider in each version are documented in our [CHANGELOG.md](https://github.com/deploymenttheory/terraform-provider-microsoft365/blob/main/CHANGELOG.md). This provider follows Semantic Versioning for releases, where the version number (MAJOR.MINOR.PATCH) indicates:

- MAJOR version for incompatible API changes
- MINOR version for functionality added in a backwards compatible manner
- PATCH version for backwards compatible bug fixes

## Contributing

We welcome contributions to the Microsoft 365 Provider! Whether it's:
- Bug reports and fixes
- Feature requests and implementations
- Documentation improvements
- Example contributions

Please visit our [GitHub repository](https://github.com/deploymenttheory/terraform-provider-microsoft365) to:
- Open issues
- Submit pull requests
- View contribution guidelines
- Join the community discussions

The provider leverages the Microsoft Graph API through the Kiota-generated SDKs, making it a powerful tool for managing Microsoft 365 resources through Terraform.