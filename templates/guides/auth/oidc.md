---
page_title: "Authentication with Generic OIDC"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using generic OpenID Connect (OIDC) tokens.
---

# Authentication with Generic OIDC

The Microsoft 365 provider supports authentication using generic OpenID Connect (OIDC) tokens. This approach allows for federated authentication from any OIDC-compatible identity provider, enabling secure authentication without managing long-lived secrets.

## How Generic OIDC Authentication Works

1. A trusted OIDC provider generates a token (JWT)
2. The token is provided to the Microsoft 365 provider through a file or environment variable
3. The provider exchanges this token for a Microsoft Graph access token
4. This exchange is based on a trust relationship established in Microsoft Entra ID

## Prerequisites

- An OIDC token provider capable of generating valid JWTs
- Permissions to create and configure app registrations in Microsoft Entra ID
- Ability to configure federated identity credentials

## Setup

### 1. Create an App Registration

```bash
# Set variables
TENANT_ID="00000000-0000-0000-0000-000000000000"
APP_NAME="terraform-m365-provider"

# Create the app registration
APP_ID=$(az ad app create --display-name $APP_NAME --query appId -o tsv)

# Create service principal for the application
az ad sp create --id $APP_ID

# Grant API permissions
az ad app permission add \
  --id $APP_ID \
  --api 00000003-0000-0000-c000-000000000000 \
  --api-permissions PERMISSIONS_LIST_HERE

# Grant admin consent
az ad app permission admin-consent --id $APP_ID
```

### 2. Configure Federated Identity Credential

Configure a federated credential in your Entra ID application. The configuration will depend on your specific OIDC provider.

```bash
# Example for a generic OIDC provider
az ad app federated-credential create \
  --id $APP_ID \
  --parameters "{\"name\":\"generic-oidc-credential\",\"issuer\":\"https://token.issuer.example.com\",\"subject\":\"specific-subject-claim\",\"description\":\"Generic OIDC federated credential\",\"audiences\":[\"api://AzureADTokenExchange\"]}"
```

The key parameters to configure are:
- `issuer`: The OIDC issuer URL of your identity provider
- `subject`: The subject identity you want to trust (varies by provider)
- `audiences`: The intended audience of the token (usually "api://AzureADTokenExchange")

## Configuration

### Using Terraform Configuration

```terraform
provider "microsoft365" {
  auth_method = "oidc"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id          = "00000000-0000-0000-0000-000000000000"
    oidc_token_file_path = "/path/to/oidc-token.jwt"
  }
}
```

### Using Environment Variables (Recommended)

```bash
# Set these environment variables before running Terraform
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="oidc"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export M365_OIDC_TOKEN_FILE_PATH="/path/to/oidc-token.jwt"
```

Then your Terraform configuration can be simplified:

```terraform
provider "microsoft365" {}
```

## Token File Format

The OIDC token file should contain a valid JWT token as plain text. For example:

```
eyJhbGciOiJSUzI1NiIsImtpZCI6IkMyRjU2RDU1MkYyQzNCQzg2MDI4MjRCNjA2QkM3NzdDIiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwczovL3Rva2VuLmlzc3Vlci5leGFtcGxlLmNvbSIsInN1YiI6InNwZWNpZmljLXN1YmplY3QtY2xhaW0iLCJhdWQiOiJhcGk6Ly9BenVyZUFEVG9rZW5FeGNoYW5nZSIsImV4cCI6MTY5OTEyMzQ1NiwiaWF0IjoxNjk5MTIzMTU2fQ.signature
```

The JWT should include:
- `iss` (issuer): Must match the issuer configured in the federated credential
- `sub` (subject): Must match the subject configured in the federated credential
- `aud` (audience): Typically "api://AzureADTokenExchange"

## Use Cases

Generic OIDC authentication is ideal for:

- Custom CI/CD systems that support OIDC
- Self-hosted automation platforms
- Integration with non-standard identity providers
- Any scenario where you can generate a compliant OIDC token

## Security Considerations

- The OIDC token should have a short lifetime (typically under 1 hour)
- Protect access to the token file using appropriate file system permissions
- Ensure your OIDC issuer is properly secured
- Configure the federated credential with specific claims to limit trust
- Consider adding conditional access policies for additional security

## Troubleshooting

- **Invalid token**: Ensure the token is valid and not expired
- **Token not found**: Verify the path to the token file is correct
- **Authentication failed**: Check that the issuer, subject, and audience in the token match the federated credential configuration
- **Permission denied**: Ensure you've granted admin consent for the required Microsoft Graph permissions
- **Missing claims**: Verify your token includes all required claims (iss, sub, aud)

## Examples for Specific Providers

### HashiCorp Vault

If using HashiCorp Vault as your OIDC provider:

```bash
# Generate a token from Vault
AZURE_OIDC_TOKEN=$(vault read -field=token identity/oidc/token/azure-role)

# Save to file
echo $AZURE_OIDC_TOKEN > /path/to/oidc-token.jwt
```

### Custom Identity Server

For a custom identity server, ensure it supports the following:

1. JWT token generation with RS256 signing
2. Configuration of issuer, subject, and audience claims
3. Proper key rotation and token validation

## Additional Resources

- [OpenID Connect Core 1.0 specification](https://openid.net/specs/openid-connect-core-1_0.html)
- [Microsoft Entra ID Workload Identity Federation](https://learn.microsoft.com/en-us/azure/active-directory/develop/workload-identity-federation)