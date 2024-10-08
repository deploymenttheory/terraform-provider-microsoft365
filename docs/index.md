---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "microsoft365 Provider"
subcategory: ""
description: |-
  
---

# microsoft365 Provider





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
