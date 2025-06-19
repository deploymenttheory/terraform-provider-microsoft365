# Microsoft Entra ID JSON Responses for Authentication

This document provides JSON response examples for authentication methods used by this Terraform provider.

## Primary References

**Official Azure SDK for Go Documentation:**

- [Azure SDK for Go GitHub Repository](https://github.com/Azure/azure-sdk-for-go)
- [azidentity Package Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity)

**Microsoft Identity Platform Documentation:**

- [Microsoft Identity Platform Overview](https://learn.microsoft.com/en-us/entra/identity-platform/)
- [OAuth 2.0 and OpenID Connect Protocols](https://learn.microsoft.com/en-us/entra/identity-platform/v2-protocols)

## Common Data Structures

### AccessToken Response (Successful)

```json
{
  "token_type": "Bearer",
  "scope": "https://graph.microsoft.com/.default",
  "expires_in": 3599,
  "ext_expires_in": 3599,
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q...",
  "refresh_token": "AwABAAAAvPM1KaPlrEqdFSBzjqfTGAMxZGUTdM0t4B4...",
  "id_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJhdWQiOiIyZDRkMTFhMi1mODE0LTQ2YTctOD..."
}
```

### Error Response (Common Structure)

```json
{
  "error": "invalid_client",
  "error_description": "AADSTS70002: Error validating credentials. AADSTS50012: Invalid client secret is provided.\r\nTrace ID: 12345678-1234-1234-1234-123456789012\r\nCorrelation ID: 12345678-1234-1234-1234-123456789012\r\nTimestamp: 2025-06-17 10:30:00Z",
  "error_codes": [70002, 50012],
  "timestamp": "2025-06-17 10:30:00Z",
  "trace_id": "12345678-1234-1234-1234-123456789012",
  "correlation_id": "12345678-1234-1234-1234-123456789012"
}
```

---

## Supported Authentication Methods

The provider supports the following authentication methods. The sections below provide details and example JSON responses for each.

1.  [`azure_developer_cli`](#1-azure_developer_cli)
2.  [`client_secret`](#2-client_secret)
3.  [`client_certificate`](#3-client_certificate)
4.  [`device_code`](#4-device_code)
5.  [`interactive_browser`](#5-interactive_browser)
6.  [`workload_identity`](#6-workload_identity)
7.  [`managed_identity`](#7-managed_identity)
8.  [`oidc`](#8-oidc)
9.  [`oidc_github`](#9-oidc_github)
10. [`oidc_azure_devops`](#10-oidc_azure_devops)


---

## 1. `azure_developer_cli`

Authenticates as the user logged into the Azure Developer CLI (`azd`). The user must first run `azd auth login` to authenticate.

**References:**
- [AzureDeveloperCLICredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#AzureDeveloperCLICredential)

### Step 1: `azd` Token Request

The provider invokes the Azure Developer CLI to obtain a token. This process is opaque to the provider.

### Step 2: Successful Token Response

If `azd` has a valid, cached token, it will be returned. The JSON response from Microsoft Entra ID has the standard `AccessToken` structure.

```json
{
  "token_type": "Bearer",
  "expires_in": 3599,
  "ext_expires_in": 3599,
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q..."
}
```

### Error: Azure Developer CLI Not Available

If `azd` is not installed or not in the system's PATH, the SDK will return an error.

```json
{
  "error": "credential_unavailable",
  "error_description": "Azure Developer CLI is not installed or not in the PATH. Please install it from https://aka.ms/azure-dev/install.",
  "error_type": "CredentialUnavailableError"
}
```

### Error: Not Logged In

If the user is not logged into `azd`, the tool will fail to provide a token.

```json
{
  "error": "authentication_failed",
  "error_description": "Please run 'azd auth login' to setup account.",
  "error_type": "AuthenticationFailedError"
}
```

---

## 2. `client_secret` (Client Credentials Flow)

**References:**

- [OAuth 2.0 Client Credentials Flow](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-client-creds-grant-flow)
- [ClientSecretCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#ClientSecretCredential)

### Step 1: Token Request to Microsoft Entra ID

**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token`

**Body:**

```
grant_type=client_credentials
&client_id=12345678-1234-1234-1234-123456789012
&client_secret=secretValue123
&scope=https://graph.microsoft.com/.default
```

### Step 2: Successful Token Response

```json
{
  "token_type": "Bearer",
  "expires_in": 3599,
  "ext_expires_in": 3599,
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q..."
}
```

### Error: Invalid Client Secret

```json
{
  "error": "invalid_client",
  "error_description": "AADSTS70002: Error validating credentials. AADSTS50012: Invalid client secret is provided.\r\nTrace ID: 12345678-1234-1234-1234-123456789012\r\nCorrelation ID: 12345678-1234-1234-1234-123456789012\r\nTimestamp: 2025-06-17 10:30:00Z",
  "error_codes": [70002, 50012],
  "timestamp": "2025-06-17 10:30:00Z",
  "trace_id": "12345678-1234-1234-1234-123456789012",
  "correlation_id": "12345678-1234-1234-1234-123456789012"
}
```

---

## 3. `client_certificate`

**References:**

- [Certificate Credentials](https://learn.microsoft.com/en-us/entra/identity-platform/certificate-credentials)
- [ClientCertificateCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#ClientCertificateCredential)

### Step 1: Certificate-based Token Request

**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token`

**Body:**

```
grant_type=client_credentials
&client_id=12345678-1234-1234-1234-123456789012
&client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer
&client_assertion=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6...
&scope=https://graph.microsoft.com/.default
```

### Step 2: Successful Response

```json
{
  "token_type": "Bearer",
  "expires_in": 3599,
  "ext_expires_in": 3599,
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q..."
}
```

### Error: Invalid Certificate

```json
{
  "error": "invalid_client",
  "error_description": "AADSTS70002: Error validating credentials. AADSTS50013: Invalid client certificate.\r\nTrace ID: 12345678-1234-1234-1234-123456789012\r\nCorrelation ID: 12345678-1234-1234-1234-123456789012\r\nTimestamp: 2025-06-17 10:30:00Z",
  "error_codes": [70002, 50013],
  "timestamp": "2025-06-17 10:30:00Z",
  "trace_id": "12345678-1234-1234-1234-123456789012",
  "correlation_id": "12345678-1234-1234-1234-123456789012"
}
```

---

## 4. `device_code`

**References:**

- [OAuth 2.0 Device Authorization Grant](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-device-code)
- [DeviceCodeCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#DeviceCodeCredential)

### Step 1: Device Code Initiation Request

**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/devicecode`

**Body:** `client_id=12345678-1234-1234-1234-123456789012&scope=https://graph.microsoft.com/.default`

### Step 2: Device Code Response

```json
{
  "user_code": "A1B2C3D4E",
  "device_code": "DAQABAAEAAAD...Zz0",
  "verification_uri": "https://microsoft.com/devicelogin",
  "expires_in": 900,
  "interval": 5,
  "message": "To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code A1B2C3D4E to authenticate."
}
```

### Step 3: Token Polling Request

The provider will poll the token endpoint until the user authenticates.

**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token`

**Body:** `grant_type=urn:ietf:params:oauth:grant-type:device_code&client_id=...&device_code=...`

### Step 4: Successful Token Response (After User Authentication)

This response is returned once the user completes the device login flow in their browser.

```json
{
  "token_type": "Bearer",
  "scope": "https://graph.microsoft.com/.default",
  "expires_in": 3599,
  "ext_expires_in": 3599,
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q...",
  "refresh_token": "AwABAAAAvPM1KaPlrEqdFSBzjqfTGAMxZGUTdM0t4B4..."
}
```

### Error: Pending User Action

While polling before the user has authenticated, the service returns this error.

```json
{
  "error": "authorization_pending",
  "error_description": "AADSTS70016: The request is pending and the device is not yet authorized. Please try again later.\r\n[...]"
}
```

---

## 5. `interactive_browser`

**References:**

- [InteractiveBrowserCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#InteractiveBrowserCredential)

This flow opens the default system browser for the user to authenticate. The JSON responses are handled by the SDK and are not directly exposed. The end result is either a successful token acquisition or an error.

### Successful Authentication

A successful authentication results in a standard `AccessToken` being available to the provider.

### Error: Browser Unavailable or User Cancellation

If a browser cannot be opened or the user cancels the authentication, an error is returned by the SDK.

```json
{
  "error": "authentication_failed",
  "error_description": "Failed to open browser: ... or user canceled authentication.",
  "error_type": "AuthenticationFailedError"
}
```
---

## 6. `workload_identity`

**References:**

- [Microsoft Entra Workload ID](https://learn.microsoft.com/en-us/azure/aks/workload-identity-overview)
- [WorkloadIdentityCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#WorkloadIdentityCredential)

This method uses a service account token from the environment (e.g., Kubernetes) to federate with Microsoft Entra ID.

### Step 1: Read Federated Token

The provider reads a federated token from a file path specified by environment variables (e.g., `AZURE_FEDERATED_TOKEN_FILE`).

### Step 2: Token Request with Federated Token

The federated token is used as a `client_assertion`.

**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token`

**Body:**
```
grant_type=client_credentials
&client_id=12345678-1234-1234-1234-123456789012
&client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer
&client_assertion=eyJhbGciOiJSUzI1Ni...
&scope=https://graph.microsoft.com/.default
```

### Step 3: Successful Token Response

```json
{
  "token_type": "Bearer",
  "expires_in": 3599,
  "ext_expires_in": 3599,
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q..."
}
```

### Error: Federated Token Invalid

```json
{
  "error": "invalid_request",
  "error_description": "AADSTS70021: No matching federated identity record found for presented assertion.",
  "error_codes": [70021]
}
```

---

## 7. `managed_identity`

**References:**

- [Managed Identity Documentation](https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication-managed-identity)
- [ManagedIdentityCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#ManagedIdentityCredential)

### Step 1: IMDS Probe Request

The provider requests a token from the Instance Metadata Service (IMDS) endpoint within an Azure host.

**Request:** `GET http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://graph.microsoft.com/`

**Headers:** `Metadata: true`

### Step 2: Successful IMDS Response

```json
{
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q...",
  "expires_in": "86399",
  "token_type": "Bearer",
  "resource": "https://graph.microsoft.com/"
}
```

### Error: IMDS Not Available

```json
{
  "error": "imds_unavailable",
  "error_description": "IMDS endpoint is not reachable. This typically indicates the application is not running in an Azure hosting environment with Managed Identity enabled.",
  "error_type": "CredentialUnavailableError"
}
```

---

## 8. `oidc` (Generic OIDC / Workload Identity Federation)

This method uses an OIDC token from an external identity provider (IdP) to authenticate, following the [Workload Identity Federation](https://learn.microsoft.com/en-us/entra/workload-id/workload-identity-federation) flow. The external token is exchanged for a Microsoft Entra ID access token.

This is the underlying mechanism for provider-specific OIDC authentication like `oidc_github` and `oidc_azure_devops`.

### Step 1: Obtain OIDC Token from External IdP

The CI/CD environment or external workload is responsible for obtaining a JWT from the OIDC provider. This token is typically made available to the provider via an environment variable (e.g., `AZURE_FEDERATED_TOKEN_FILE`).

### Step 2: Token Exchange Request

The provider uses the obtained OIDC token as a `client_assertion` in a client credentials grant request to Microsoft Entra ID.

**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token`

**Body:**
```
grant_type=client_credentials
&client_id={client_id}
&client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer
&client_assertion={oidc_token_from_idp}
&scope=https://graph.microsoft.com/.default
```

### Step 3: Successful Token Response

If the `client_assertion` is valid and the federated credential is configured correctly in Microsoft Entra ID, a standard access token is returned.

```json
{
  "token_type": "Bearer",
  "expires_in": 3599,
  "ext_expires_in": 3599,
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q..."
}
```

### Error: Invalid Federated Credential

If the trust relationship is not configured correctly (e.g., mismatched issuer, subject, or audience), Entra ID will reject the assertion.

```json
{
  "error": "invalid_request",
  "error_description": "AADSTS70021: No matching federated identity record found for presented assertion.",
  "error_codes": [70021]
}
```

---

## 9. `oidc_github` (GitHub OIDC Provider)

This method authenticates using an OIDC token provided by a GitHub Actions workflow. A [federated identity credential](https://learn.microsoft.com/en-us/entra/workload-id/workload-identity-federation-create-trust?pivots=identity-wif-apps-methods-azcli#github-actions-example) must be configured on the service principal in Microsoft Entra ID.

**References:**
- [Configuring a federated credential for a GitHub repo](https://learn.microsoft.com/en-us/entra/workload-id/workload-identity-federation-create-trust?pivots=identity-wif-apps-methods-azcli#github-actions-example)

### Federated Credential Configuration (in Entra ID)
- **Issuer**: `https://token.actions.githubusercontent.com`
- **Audience**: `api://AzureADTokenExchange`
- **Subject**: Varies based on the trigger (e.g., `repo:my-org/my-repo:ref:refs/heads/main`)

The token exchange request and responses are the same as the [generic `oidc` flow](#8-oidc-generic-oidc--workload-identity-federation).

---

## 10. `oidc_azure_devops` (Azure DevOps OIDC Provider)

This method authenticates using an OIDC token from an Azure DevOps pipeline via a [Workload Identity federation service connection](https://learn.microsoft.com/en-us/azure/devops/pipelines/release/configure-workload-identity?view=azure-devops).

**References:**
- [Introduction to Azure DevOps Workload identity federation](https://devblogs.microsoft.com/devops/introduction-to-azure-devops-workload-identity-federation-oidc-with-terraform/)

### Federated Credential Configuration (in Entra ID)

- **Issuer**: `https://vstoken.dev.azure.com/{organization_id}`
- **Audience**: `api://AzureADTokenExchange`
- **Subject**: `sc://{organization_name}/{project_name}/{service_connection_name}`

The token exchange request and responses are the same as the [generic `oidc` flow](#8-oidc-generic-oidc--workload-identity-federation).