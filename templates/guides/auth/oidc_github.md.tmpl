---
page_title: "Authentication with GitHub OIDC"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using GitHub Actions OIDC tokens.
---

# Authentication with GitHub OIDC

The Microsoft 365 provider supports authentication using GitHub Actions' OpenID Connect (OIDC) tokens. This approach allows Terraform to authenticate to Microsoft 365 services directly from GitHub Actions workflows without storing long-lived credentials as GitHub secrets.

## Table of Contents

- [How GitHub OIDC Authentication Works](#how-github-oidc-authentication-works)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
  - [Create an App Registration in Microsoft Entra ID](#1-create-an-app-registration-in-microsoft-entra-id)
  - [Configure Federated Identity Credential](#2-configure-federated-identity-credential)
- [Common Subject Patterns](#common-subject-patterns)
- [GitHub Actions Workflow Configuration](#github-actions-workflow-configuration)
- [Provider Configuration](#provider-configuration)
  - [Using Terraform Configuration](#using-terraform-configuration)
  - [Using Environment Variables](#using-environment-variables-recommended)
- [Integration with HashiCorp Vault](#integration-with-hashicorp-vault)
  - [Vault JWT Auth Method](#vault-jwt-auth-method)
  - [GitHub Actions with Vault](#github-actions-with-vault)
  - [Security Considerations for Vault Integration](#security-considerations-for-vault-integration)
- [Security Best Practices](#security-best-practices)
- [Troubleshooting](#troubleshooting)
- [Additional Resources](#additional-resources)

## How GitHub OIDC Authentication Works

1. GitHub Actions generates a short-lived OIDC token for each workflow run
2. The Microsoft 365 provider requests this token during Terraform execution
3. The provider presents the token to Microsoft Entra ID
4. Based on a pre-configured trust relationship, Entra ID issues a Microsoft Graph access token
5. The provider uses this token to authenticate API requests

Key benefits of this approach include:

- **No stored secrets**: Long-lived secrets don't need to be stored in GitHub
- **Automatic rotation**: Tokens are short-lived and automatically rotated
- **Conditional access**: Fine-grained control over which workflows can obtain tokens
- **Reduced attack surface**: Eliminates risk of leaked or compromised credentials

## Prerequisites

- A GitHub repository where you'll run Terraform
- Permissions to create and configure app registrations in Microsoft Entra ID
- Ability to modify GitHub Actions workflows
- Azure CLI installed (for setup commands)

## Setup

### 1. Create an App Registration in Microsoft Entra ID

You can create the app registration using Azure CLI:

```bash
# Set variables
TENANT_ID="00000000-0000-0000-0000-000000000000"
APP_NAME="terraform-m365-provider"
GITHUB_ORG="your-github-org"
GITHUB_REPO="your-github-repo"

# Create the app registration
APP_ID=$(az ad app create --display-name $APP_NAME --query appId -o tsv)

# Create service principal for the application
az ad sp create --id $APP_ID

# Grant API permissions for Microsoft Graph
az ad app permission add \
  --id $APP_ID \
  --api 00000003-0000-0000-c000-000000000000 \
  --api-permissions PERMISSIONS_LIST_HERE

# Grant admin consent
az ad app permission admin-consent --id $APP_ID
```

Alternatively, you can create the app registration using the Azure portal:

1. Go to Microsoft Entra ID > App registrations > New registration
2. Enter a name for the application (e.g., "terraform-m365-provider")
3. Select "Accounts in this organizational directory only"
4. Click "Register"
5. Navigate to "API permissions"
6. Add the necessary Microsoft Graph permissions
7. Grant admin consent
8. Note the Application (client) ID and Tenant ID for later use

### 2. Configure Federated Identity Credential

Create a federated credential in your Entra ID application that trusts GitHub Actions:

```bash
# For a specific branch
az ad app federated-credential create \
  --id $APP_ID \
  --parameters "{\"name\":\"github-federated-credential\",\"issuer\":\"https://token.actions.githubusercontent.com\",\"subject\":\"repo:${GITHUB_ORG}/${GITHUB_REPO}:ref:refs/heads/main\",\"description\":\"GitHub Actions federated credential\",\"audiences\":[\"api://AzureADTokenExchange\"]}"

# For all branches
az ad app federated-credential create \
  --id $APP_ID \
  --parameters "{\"name\":\"github-federated-credential-all-branches\",\"issuer\":\"https://token.actions.githubusercontent.com\",\"subject\":\"repo:${GITHUB_ORG}/${GITHUB_REPO}:*\",\"description\":\"GitHub Actions federated credential for all branches\",\"audiences\":[\"api://AzureADTokenExchange\"]}"
```

You can create multiple federated credentials with different subject filters to control which workflows can obtain tokens.

## Common Subject Patterns

| Scenario | Subject Format | Example |
|----------|----------------|---------|
| Specific branch | `repo:ORG/REPO:ref:refs/heads/BRANCH` | `repo:octo-org/octo-repo:ref:refs/heads/main` |
| Any branch | `repo:ORG/REPO:*` | `repo:octo-org/octo-repo:*` |
| Pull requests | `repo:ORG/REPO:pull_request` | `repo:octo-org/octo-repo:pull_request` |
| Specific tag | `repo:ORG/REPO:ref:refs/tags/TAG` | `repo:octo-org/octo-repo:ref:refs/tags/v1.0.0` |
| Specific environment | `repo:ORG/REPO:environment:ENV` | `repo:octo-org/octo-repo:environment:production` |

## GitHub Actions Workflow Configuration

Configure your GitHub Actions workflow to request and use the OIDC token:

```yaml
name: Terraform Microsoft 365

on:
  push:
    branches: [ main ]

# Permission to request the OIDC JWT ID token
permissions:
  id-token: write  # Required for OIDC
  contents: read   # Required for checkout

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v2
        with:
          terraform_version: "1.5.0"

      - name: Terraform Init
        run: terraform init

      - name: Terraform Apply
        env:
          M365_TENANT_ID: "00000000-0000-0000-0000-000000000000"
          M365_AUTH_METHOD: "oidc_github"
          M365_CLIENT_ID: "00000000-0000-0000-0000-000000000000"
          # Optional: Set a custom audience if needed
          # M365_OIDC_AUDIENCE: "api://AzureADTokenExchange"
        run: terraform apply -auto-approve
```

## Provider Configuration

### Using Terraform Configuration

```terraform
provider "microsoft365" {
  auth_method = "oidc_github"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id = "00000000-0000-0000-0000-000000000000"
  }
}
```

### Using Environment Variables (Recommended)

```bash
# Set these environment variables in your GitHub Actions workflow
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="oidc_github"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
```

Then your Terraform configuration can be simplified:

```terraform
provider "microsoft365" {}
```

## Integration with HashiCorp Vault

You can enhance your GitHub Actions OIDC authentication by integrating with HashiCorp Vault to securely manage and retrieve sensitive configuration data needed for Microsoft 365 provider.

### Vault JWT Auth Method

HashiCorp Vault's JWT auth method can validate the GitHub Actions OIDC token:

1. Enable the JWT auth method in Vault:

```bash
vault auth enable jwt
```

2. Configure the JWT auth method to trust GitHub Actions tokens:

```bash
vault write auth/jwt/config \
  bound_issuer="https://token.actions.githubusercontent.com" \
  jwks_url="https://token.actions.githubusercontent.com/.well-known/jwks"
```

3. Create a role for your GitHub repository:

```bash
vault write auth/jwt/role/github-oidc \
  role_type="jwt" \
  bound_audiences="api://AzureADTokenExchange" \
  bound_claims_type="glob" \
  bound_claims='{"sub":"repo:your-org/your-repo:*"}' \
  user_claim="sub" \
  policies="microsoft365-policy" \
  ttl="1h"
```

4. Create a policy to allow access to your Microsoft 365 credentials:

```bash
vault policy write microsoft365-policy - <<EOF
path "secret/data/microsoft365/*" {
  capabilities = ["read"]
}
EOF
```

5. Store your Microsoft 365 configuration in Vault:

```bash
vault kv put secret/microsoft365/credentials \
  tenant_id="00000000-0000-0000-0000-000000000000" \
  client_id="00000000-0000-0000-0000-000000000000"
```

### GitHub Actions with Vault

Integrate Vault into your GitHub Actions workflow:

```yaml
name: Terraform Microsoft 365 with Vault

on:
  push:
    branches: [ main ]

permissions:
  id-token: write
  contents: read

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
        
      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v2
        with:
          terraform_version: "1.5.0"
          
      - name: Authenticate to Vault with GitHub Actions OIDC
        id: vault-auth
        run: |
          # Install Vault CLI
          curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
          sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
          sudo apt-get update && sudo apt-get install vault
          
          # Get GitHub Actions OIDC token
          JWT=$(curl -s -H "Authorization: bearer $ACTIONS_ID_TOKEN_REQUEST_TOKEN" "$ACTIONS_ID_TOKEN_REQUEST_URL" | jq -r .value)
          
          # Authenticate to Vault using the JWT
          VAULT_ADDR="https://your-vault-server:8200"
          VAULT_AUTH_RESPONSE=$(vault write -format=json auth/jwt/login role=github-oidc jwt=$JWT)
          
          # Export Vault token
          VAULT_TOKEN=$(echo $VAULT_AUTH_RESPONSE | jq -r .auth.client_token)
          echo "::add-mask::$VAULT_TOKEN"
          echo "VAULT_TOKEN=$VAULT_TOKEN" >> $GITHUB_ENV
          echo "VAULT_ADDR=$VAULT_ADDR" >> $GITHUB_ENV
      
      - name: Get Microsoft 365 credentials from Vault
        run: |
          # Fetch credentials from Vault
          VAULT_DATA=$(vault kv get -format=json secret/microsoft365/credentials)
          
          # Set Microsoft 365 provider environment variables
          echo "M365_TENANT_ID=$(echo $VAULT_DATA | jq -r .data.data.tenant_id)" >> $GITHUB_ENV
          echo "M365_AUTH_METHOD=oidc_github" >> $GITHUB_ENV
          echo "M365_CLIENT_ID=$(echo $VAULT_DATA | jq -r .data.data.client_id)" >> $GITHUB_ENV
          
      - name: Terraform Init
        run: terraform init
        
      - name: Terraform Apply
        run: terraform apply -auto-approve
```

### Security Considerations for Vault Integration

When integrating HashiCorp Vault with GitHub OIDC:

1. **Time-bound tokens**: Configure short TTLs for Vault tokens to limit exposure
2. **Least privilege**: Create specific Vault policies that grant access only to required secrets
3. **Bound claims**: Tightly scope JWT roles to specific repositories and workflows
4. **Audit logging**: Enable comprehensive auditing in Vault to track token usage
5. **Response wrapping**: Use Vault's response wrapping for additional security when retrieving sensitive values
6. **Secure Vault communication**: Ensure TLS is configured properly for all communication to Vault
7. **Namespace isolation**: In Vault Enterprise, use namespaces to isolate different teams and repositories
8. **Dynamic secrets**: Consider using Vault's dynamic secrets capabilities where possible

## Security Best Practices

1. **Narrow the scope of trust**:
   - Use specific branch, environment, or tag conditions instead of wildcard subjects
   - Create separate app registrations for different repositories or workflows

2. **Add conditional access policies**:
   - Configure additional conditions in Microsoft Entra ID
   - Restrict access based on IP ranges, specific actions, or other attributes

3. **Limit API permissions**:
   - Grant only the minimum required permissions to the app registration
   - Use application-level permissions rather than delegated permissions for automation

4. **Enable auditing**:
   - Monitor token issuance and usage in Microsoft Entra ID logs
   - Set up alerts for suspicious authentication patterns

## Troubleshooting

- **Missing ID token**: Ensure the workflow has `id-token: write` permission
- **Authentication failed**: Verify the federated credential is configured correctly
- **Subject mismatch**: Check that the subject pattern in Entra ID matches your GitHub workflow conditions
- **Permission denied**: Ensure you've granted admin consent for the required Microsoft Graph permissions
- **Invalid audience**: Confirm the audience claim matches what's expected in the federated credential

For GitHb OIDC specific errors:

```
Error: Failed to get token: failed to exchange GitHub Actions OIDC token: oauth2: "invalid_client" "AADSTS7000215: Invalid client secret provided"
```
This could indicate that the federated credential is not properly configured. Verify the subject and issuer are correct.

```
Error: Failed to get token: OIDC provider not available
```
This may occur if you're trying to use GitHub OIDC authentication outside of a GitHub Actions workflow. This authentication method only works within GitHub Actions.

## Additional Resources

- [GitHub OIDC Documentation](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/configuring-openid-connect-in-azure)
- [Microsoft Entra ID Workload Identity Federation](https://learn.microsoft.com/en-us/azure/active-directory/develop/workload-identity-federation)
- [Securing GitHub Actions with OpenID Connect](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect)
- [HashiCorp Vault JWT Auth Method](https://developer.hashicorp.com/vault/docs/auth/jwt)
- [Vault with GitHub Actions](https://developer.hashicorp.com/vault/tutorials/app-integration/github-actions)