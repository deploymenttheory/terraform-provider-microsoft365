---
page_title: "Authentication with GitHub OIDC"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using GitHub Actions OIDC tokens.
---

# Authentication with GitHub OIDC

The Microsoft 365 provider supports authentication using GitHub Actions' OpenID Connect (OIDC) tokens. This approach allows Terraform to authenticate to Microsoft 365 services directly from GitHub Actions workflows without storing long-lived credentials as GitHub secrets.

## How GitHub OIDC Authentication Works

1. GitHub Actions generates a short-lived OIDC token for each workflow run
2. The Microsoft 365 provider requests this token during Terraform execution
3. The provider presents the token to Microsoft Entra ID
4. Based on a pre-configured trust relationship, Entra ID issues a Microsoft Graph access token
5. The provider uses this token to authenticate API requests

## Prerequisites

- A GitHub repository where you'll run Terraform
- Permissions to create and configure app registrations in Microsoft Entra ID
- Ability to modify GitHub Actions workflows

## Setup

### 1. Create an App Registration in Microsoft Entra ID

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

## Additional Resources

- [GitHub OIDC Documentation](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect)
- [Microsoft Entra ID Workload Identity Federation](https://learn.microsoft.com/en-us/azure/active-directory/develop/workload-identity-federation)
- [Securing GitHub Actions with OpenID Connect](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect)