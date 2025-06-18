# Azure SDK for Go azidentity - JSON Responses for Authentication Flows

This document provides comprehensive JSON response examples for each authentication method supported by the Azure SDK for Go's azidentity package.

## Primary References

**Official Azure SDK for Go Documentation:**
- [Azure SDK for Go GitHub Repository](https://github.com/Azure/azure-sdk-for-go)
- [azidentity Package Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity)
- [Azure SDK for Go README](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/README.md)

**Microsoft Identity Platform Documentation:**
- [Microsoft Identity Platform Overview](https://learn.microsoft.com/en-us/entra/identity-platform/)
- [OAuth 2.0 and OpenID Connect Protocols](https://learn.microsoft.com/en-us/entra/identity-platform/v2-protocols)
- [Azure Developer Authentication Guide](https://learn.microsoft.com/en-us/azure/developer/go/sdk/authentication/authentication-overview)

**Community Resources:**
- [Azure SDK Design Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html)
- [Ravikanth Chaganti's Authentication Guide](https://ravichaganti.com/blog/azure-sdk-for-go-authentication-methods-chained-credentials/)

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

## 1. DefaultAzureCredential

**References:**
- [Azure SDK for Go azidentity package](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity)
- [DefaultAzureCredential Overview](https://learn.microsoft.com/en-us/azure/developer/go/sdk/authentication/authentication-overview)
- [Azure SDK for Go Authentication](https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication-managed-identity)

DefaultAzureCredential attempts authentication in this order: Environment → Workload Identity → Managed Identity → Azure CLI → Azure Developer CLI

### Step 1: Environment Credential Attempt (Failed)
```json
{
  "error": "credential_unavailable",
  "error_description": "Environment credential not available. Environment variables are not fully configured.",
  "error_type": "CredentialUnavailableError"
}
```

### Step 2: Workload Identity Credential Attempt (Failed)
```json
{
  "error": "credential_unavailable", 
  "error_description": "Workload identity credential not available. Required environment variables not found.",
  "error_type": "CredentialUnavailableError"
}
```

### Step 3: Managed Identity Success
```json
{
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q...",
  "expires_in": "3599",
  "token_type": "Bearer",
  "resource": "https://management.azure.com/"
}
```

## 2. ManagedIdentityCredential

**References:**
- [Managed Identity Documentation](https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication-managed-identity)
- [Azure Instance Metadata Service (IMDS)](https://learn.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service)
- [Azure SDK for Go Managed Identity](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#ManagedIdentityCredential)

### Step 1: IMDS Probe Request
**Request:** `GET http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.azure.com/`

**Headers:** `Metadata: true`

### Step 2: Successful IMDS Response
```json
{
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q...",
  "expires_in": "3599",
  "token_type": "Bearer",
  "resource": "https://management.azure.com/"
}
```

### Step 3: User-Assigned Managed Identity (with Client ID)
**Request:** `GET http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.azure.com/&client_id=12345678-1234-1234-1234-123456789012`

```json
{
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q...",
  "expires_in": "3599", 
  "token_type": "Bearer",
  "resource": "https://management.azure.com/",
  "client_id": "12345678-1234-1234-1234-123456789012"
}
```

### Error: IMDS Not Available
```json
{
  "error": "imds_unavailable",
  "error_description": "IMDS endpoint is not reachable. This typically indicates the application is not running in an Azure hosting environment.",
  "error_type": "CredentialUnavailableError"
}
```

## 3. ClientSecretCredential (Client Credentials Flow)

**References:**
- [OAuth 2.0 Client Credentials Flow](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-client-creds-grant-flow)
- [Microsoft Identity Platform OAuth 2.0](https://learn.microsoft.com/en-us/entra/identity-platform/v2-protocols)
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

## 4. ClientCertificateCredential

**References:**
- [OAuth 2.0 Client Credentials with Certificates](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-client-creds-grant-flow)
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

## 5. DeviceCodeCredential

**References:**
- [OAuth 2.0 Device Authorization Grant](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-device-code)
- [Device Code Flow in MSAL.NET](https://learn.microsoft.com/en-us/entra/msal/dotnet/acquiring-tokens/desktop-mobile/device-code-flow)
- [DeviceCodeCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#DeviceCodeCredential)
- [Device Code Flow Implementation](https://joonasw.net/view/device-code-flow)
- [Authentication Flows in MSAL](https://learn.microsoft.com/en-us/entra/identity-platform/msal-authentication-flows)

### Step 1: Device Code Initiation Request
**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/devicecode`

**Body:**
```
client_id=12345678-1234-1234-1234-123456789012
&scope=https://graph.microsoft.com/.default
```

### Step 2: Device Code Response
```json
{
  "user_code": "FKDL7G9M8",
  "device_code": "GMMhmHCXhWEzkobqIHGG_EnNYYsAkukHspeYUk9E8...",
  "verification_uri": "https://microsoft.com/devicelogin",
  "expires_in": 900,
  "interval": 5,
  "message": "To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code FKDL7G9M8 to authenticate."
}
```

### Step 3: Token Polling Request (Authorization Pending)
**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token`

**Body:**
```
grant_type=urn:ietf:params:oauth:grant-type:device_code
&client_id=12345678-1234-1234-1234-123456789012
&device_code=GMMhmHCXhWEzkobqIHGG_EnNYYsAkukHspeYUk9E8...
```

### Step 4: Polling Response (Still Pending)
```json
{
  "error": "authorization_pending",
  "error_description": "AADSTS70016: Pending end-user authorization.\r\nTrace ID: 12345678-1234-1234-1234-123456789012\r\nCorrelation ID: 12345678-1234-1234-1234-123456789012\r\nTimestamp: 2025-06-17 10:30:00Z",
  "error_codes": [70016],
  "timestamp": "2025-06-17 10:30:00Z",
  "trace_id": "12345678-1234-1234-1234-123456789012",
  "correlation_id": "12345678-1234-1234-1234-123456789012"
}
```

### Step 5: Successful Authentication Response
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

### Error: Device Code Expired
```json
{
  "error": "expired_token",
  "error_description": "AADSTS70019: Verification code expired.\r\nTrace ID: 12345678-1234-1234-1234-123456789012\r\nCorrelation ID: 12345678-1234-1234-1234-123456789012\r\nTimestamp: 2025-06-17 10:30:00Z",
  "error_codes": [70019],
  "timestamp": "2025-06-17 10:30:00Z",
  "trace_id": "12345678-1234-1234-1234-123456789012",
  "correlation_id": "12345678-1234-1234-1234-123456789012"
}
```

## 6. InteractiveBrowserCredential

**References:**
- [OAuth 2.0 Authorization Code Flow](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-auth-code-flow)
- [Microsoft Graph Authorization](https://learn.microsoft.com/en-us/graph/auth-v2-user)
- [InteractiveBrowserCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#InteractiveBrowserCredential)
- [PKCE for OAuth 2.0](https://datatracker.ietf.org/doc/html/rfc7636)

### Step 1: Authorization Request (Browser Redirect)
**URL:** `https://login.microsoftonline.com/{tenant}/oauth2/v2.0/authorize?client_id=12345678-1234-1234-1234-123456789012&response_type=code&redirect_uri=http://localhost:8080&scope=https://graph.microsoft.com/.default&state=state123&code_challenge=E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM&code_challenge_method=S256`

### Step 2: Authorization Code Response (Redirect to localhost)
**URL:** `http://localhost:8080?code=M0ab92efe-b6fd-df08-87dc-2c6500a7f84d&state=state123&session_state=fe1540c3-a69a-469a-9fa3-8a2470936421`

### Step 3: Token Exchange Request
**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token`

**Body:**
```
grant_type=authorization_code
&client_id=12345678-1234-1234-1234-123456789012
&code=M0ab92efe-b6fd-df08-87dc-2c6500a7f84d
&redirect_uri=http://localhost:8080
&code_verifier=dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk
```

### Step 4: Successful Token Response
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

## 7. AzureCLICredential

**References:**
- [Azure CLI Authentication](https://learn.microsoft.com/en-us/azure/developer/go/sdk/authentication/authentication-overview)
- [AzureCLICredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#AzureCLICredential)
- [Azure SDK for Go Authentication Methods](https://ravichaganti.com/blog/azure-sdk-for-go-authentication-methods-chained-credentials/)

### Step 1: Azure CLI Token Request (Shell Command)
**Command:** `az account get-access-token --resource https://management.azure.com/ --output json`

### Step 2: Azure CLI Success Response
```json
{
  "accessToken": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q...",
  "expiresOn": "2025-06-17 14:30:00.000000",
  "subscription": "12345678-1234-1234-1234-123456789012",
  "tenant": "87654321-4321-4321-4321-210987654321",
  "tokenType": "Bearer"
}
```

### Error: CLI Not Logged In
```json
{
  "error": "credential_unavailable",
  "error_description": "Azure CLI not found or not logged in. Please run 'az login' first.",
  "error_type": "CredentialUnavailableError"
}
```

## 8. EnvironmentCredential

**References:**
- [Environment Variables Authentication](https://learn.microsoft.com/en-us/azure/developer/go/sdk/authentication/authentication-overview)
- [EnvironmentCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#EnvironmentCredential)
- [Azure SDK Environment Configuration](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/README.md)

### Configuration via Environment Variables
```bash
AZURE_CLIENT_ID=12345678-1234-1234-1234-123456789012
AZURE_CLIENT_SECRET=secretValue123
AZURE_TENANT_ID=87654321-4321-4321-4321-210987654321
```

### Token Request
**Request:** `POST https://login.microsoftonline.com/87654321-4321-4321-4321-210987654321/oauth2/v2.0/token`

**Body:**
```
grant_type=client_credentials
&client_id=12345678-1234-1234-1234-123456789012
&client_secret=secretValue123
&scope=https://graph.microsoft.com/.default
```

### Successful Response
```json
{
  "token_type": "Bearer",
  "expires_in": 3599,
  "ext_expires_in": 3599,
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q..."
}
```

## 9. WorkloadIdentityCredential

**References:**
- [Workload Identity Federation](https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication-managed-identity)
- [Azure Workload Identity](https://azure.github.io/azure-workload-identity/docs/)
- [WorkloadIdentityCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#WorkloadIdentityCredential)
- [Token Exchange (RFC 8693)](https://datatracker.ietf.org/doc/html/rfc8693)

### Step 1: Reading Service Account Token
**File:** `/var/run/secrets/azure/tokens/azure-identity-token`

### Step 2: Token Exchange Request
**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token`

**Body:**
```
grant_type=urn:ietf:params:oauth:grant-type:token-exchange
&client_id=12345678-1234-1234-1234-123456789012
&subject_token_type=urn:ietf:params:oauth:token-type:jwt
&subject_token=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsImtpZCI6...
&scope=https://graph.microsoft.com/.default
&client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer
&client_assertion=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6...
```

### Step 3: Successful Response
```json
{
  "token_type": "Bearer",
  "expires_in": 3599,
  "ext_expires_in": 3599,
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q..."
}
```

## 10. AzureDeveloperCLICredential

**References:**
- [Azure Developer CLI](https://learn.microsoft.com/en-us/azure/developer/azure-developer-cli/)
- [AzureDeveloperCLICredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#AzureDeveloperCLICredential)
- [Azure SDK Changelog](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/CHANGELOG.md)

### Step 1: Azure Developer CLI Token Request
**Command:** `azd auth token --output json --scope https://management.azure.com/.default`

### Step 2: AZD Success Response
```json
{
  "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q...",
  "expiresOn": "2025-06-17T14:30:00.000Z"
}
```

### Error: AZD Not Available
```json
{
  "error": "credential_unavailable",
  "error_description": "Azure Developer CLI not found or not logged in. Please run 'azd auth login' first.",
  "error_type": "CredentialUnavailableError"
}
```

## 11. UsernamePasswordCredential (ROPC Flow)

**References:**
- [Resource Owner Password Credentials](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-ropc)
- [ROPC Authentication Flows](https://learn.microsoft.com/en-us/entra/identity-platform/msal-authentication-flows)
- [UsernamePasswordCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#UsernamePasswordCredential)

### Step 1: Resource Owner Password Credentials Request
**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token`

**Body:**
```
grant_type=password
&client_id=12345678-1234-1234-1234-123456789012
&username=user@example.com
&password=userPassword123
&scope=https://graph.microsoft.com/.default
```

### Step 2: Successful Response
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

### Error: MFA Required
```json
{
  "error": "invalid_grant",
  "error_description": "AADSTS50076: Due to a configuration change made by your administrator, or because you moved to a new location, you must use multi-factor authentication to access 'https://graph.microsoft.com'.\r\nTrace ID: 12345678-1234-1234-1234-123456789012\r\nCorrelation ID: 12345678-1234-1234-1234-123456789012\r\nTimestamp: 2025-06-17 10:30:00Z",
  "error_codes": [50076],
  "timestamp": "2025-06-17 10:30:00Z",
  "trace_id": "12345678-1234-1234-1234-123456789012",
  "correlation_id": "12345678-1234-1234-1234-123456789012"
}
```

## 12. OnBehalfOfCredential

**References:**
- [OAuth 2.0 On-Behalf-Of Flow](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-on-behalf-of-flow)
- [OnBehalfOfCredential Documentation](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#OnBehalfOfCredential)
- [Microsoft Identity Platform OBO](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-on-behalf-of-flow)

### Step 1: On-Behalf-Of Token Request
**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token`

**Body:**
```
grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer
&client_id=12345678-1234-1234-1234-123456789012
&client_secret=secretValue123
&assertion=eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6...
&scope=https://graph.microsoft.com/.default
&requested_token_use=on_behalf_of
```

### Step 2: Successful Response
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

## 13. Token Refresh Flow

**References:**
- [OAuth 2.0 Token Refresh](https://learn.microsoft.com/en-us/graph/auth-v2-user)
- [Access Token Refresh](https://learn.microsoft.com/en-us/entra/identity-platform/v2-oauth2-auth-code-flow)
- [Token Caching in azidentity](https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/azidentity/TOKEN_CACHING.md)

### Step 1: Refresh Token Request
**Request:** `POST https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token`

**Body:**
```
grant_type=refresh_token
&client_id=12345678-1234-1234-1234-123456789012
&client_secret=secretValue123
&refresh_token=AwABAAAAvPM1KaPlrEqdFSBzjqfTGAMxZGUTdM0t4B4...
&scope=https://graph.microsoft.com/.default
```

### Step 2: Successful Refresh Response
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

### Error: Refresh Token Expired
```json
{
  "error": "invalid_grant",
  "error_description": "AADSTS70002: Error validating credentials. AADSTS54005: OAuth2 Authorization code was already redeemed, please retry with a new valid code or use an existing refresh token.\r\nTrace ID: 12345678-1234-1234-1234-123456789012\r\nCorrelation ID: 12345678-1234-1234-1234-123456789012\r\nTimestamp: 2025-06-17 10:30:00Z",
  "error_codes": [70002, 54005],
  "timestamp": "2025-06-17 10:30:00Z",
  "trace_id": "12345678-1234-1234-1234-123456789012",
  "correlation_id": "12345678-1234-1234-1234-123456789012"
}
```

## 14. OIDC Discovery Document

**References:**
- [OpenID Connect Discovery](https://learn.microsoft.com/en-us/entra/identity-platform/v2-protocols-oidc)
- [Microsoft Identity Platform Endpoints](https://learn.microsoft.com/en-us/entra/identity-platform/v2-protocols)
- [OpenID Connect Specification](https://openid.net/specs/openid-connect-discovery-1_0.html)
- [Azure AD Token Validation](https://www.voitanos.io/blog/validating-entra-id-generated-oauth-tokens/)

### Discovery Endpoint Response
**URL:** `https://login.microsoftonline.com/{tenant}/v2.0/.well-known/openid-configuration`

```json
{
  "authorization_endpoint": "https://login.microsoftonline.com/{tenant}/oauth2/v2.0/authorize",
  "token_endpoint": "https://login.microsoftonline.com/{tenant}/oauth2/v2.0/token",
  "token_endpoint_auth_methods_supported": [
    "client_secret_post",
    "private_key_jwt",
    "client_secret_basic"
  ],
  "jwks_uri": "https://login.microsoftonline.com/{tenant}/discovery/v2.0/keys",
  "response_modes_supported": [
    "query",
    "fragment",
    "form_post"
  ],
  "subject_types_supported": [
    "pairwise"
  ],
  "id_token_signing_alg_values_supported": [
    "RS256"
  ],
  "response_types_supported": [
    "code",
    "id_token",
    "code id_token",
    "token id_token",
    "token"
  ],
  "scopes_supported": [
    "openid",
    "profile",
    "email",
    "offline_access"
  ],
  "issuer": "https://login.microsoftonline.com/{tenant}/v2.0",
  "microsoft_multi_refresh_token": true,
  "device_authorization_endpoint": "https://login.microsoftonline.com/{tenant}/oauth2/v2.0/devicecode",
  "http_logout_supported": true,
  "frontchannel_logout_supported": true,
  "end_session_endpoint": "https://login.microsoftonline.com/{tenant}/oauth2/v2.0/logout",
  "claims_supported": [
    "sub",
    "iss",
    "cloud_instance_name",
    "cloud_instance_host_name",
    "cloud_graph_host_name",
    "msgraph_host",
    "aud",
    "exp",
    "iat",
    "auth_time",
    "acr",
    "amr",
    "nonce",
    "email",
    "given_name",
    "family_name",
    "nickname"
  ],
  "check_session_iframe": "https://login.microsoftonline.com/{tenant}/oauth2/v2.0/checksession",
  "userinfo_endpoint": "https://graph.microsoft.com/oidc/userinfo",
  "kerberos_endpoint": "https://login.microsoftonline.com/{tenant}/kerberos",
  "tenant_region_scope": null,
  "cloud_instance_name": "microsoftonline.com",
  "cloud_graph_host_name": "graph.windows.net",
  "msgraph_host": "graph.microsoft.com",
  "rbac_url": "https://pas.windows.net"
}
```

## Notes

- All tokens (access_token, refresh_token, id_token) are JWT tokens when returned from Microsoft Entra ID
- Timestamps are in UTC format
- Error codes are specific to Microsoft Entra ID and can be used for programmatic error handling
- The `ext_expires_in` field provides extended token lifetime for resilience during service outages
- Scope values determine what resources and permissions the token can access
- Trace IDs and Correlation IDs are useful for debugging and support requests

These JSON responses represent the actual data structures you'll encounter when implementing authentication flows with the Azure SDK for Go's azidentity package.