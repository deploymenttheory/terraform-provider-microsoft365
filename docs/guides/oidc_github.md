---
page_title: "Authentication with GitHub OIDC"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using GitHub Actions OIDC tokens.
---

# Authentication with GitHub OIDC

The Microsoft 365 provider supports authentication using GitHub Actions' OpenID Connect (OIDC) tokens. This approach allows Terraform to authenticate to Microsoft 365 services directly from GitHub Actions workflows without storing long-lived credentials as GitHub secrets.

## Table of Contents

- [Prerequisites](#prerequisites)
- [How GitHub OIDC Authentication Works](#how-github-oidc-authentication-works)
- [Setup](#setup)
  - [Create an App Registration in Microsoft Entra ID](#1-create-an-app-registration-in-microsoft-entra-id)
  - [Configure Federated Identity Credential](#2-configure-federated-identity-credential)
- [Common Subject Patterns](#common-subject-patterns)
- [GitHub Actions Workflow Configuration](#github-actions-workflow-configuration)
- [Provider Configuration](#provider-configuration)
  - [Using Terraform Configuration](#using-terraform-configuration)
  - [Using Environment Variables](#using-environment-variables-recommended)
- [Security Best Practices](#security-best-practices)
- [Troubleshooting](#troubleshooting)
- [Additional Resources](#additional-resources)

## Prerequisites

- A GitHub repository where you'll run Terraform
- Permissions to create and configure app registrations in Microsoft Entra ID
- Ability to modify GitHub Actions workflows
- Azure CLI installed (for setup commands)
- Terraform provider deploymenttheory/microsoft365 version >= v0.11.0-alpha

## How GitHub OIDC Authentication Works

1. The workflow runs with permissions: id-token: write, causing the runner to prepare an OIDC token for this job.
2. GitHub injects ACTIONS_ID_TOKEN_REQUEST_URL and ACTIONS_ID_TOKEN_REQUEST_TOKEN into the job’s environment.
3. GitHubOIDCStrategy.GetCredential() reads those two variables and determines the audience (defaulting to api://AzureADTokenExchange).
4. The provider issues an HTTP GET to the GitHub URL in `ACTIONS_ID_TOKEN_REQUEST_URL`:

   ```bash
   GET https://token.actions.githubusercontent.com/<repo-owner>/<repo-name>/_apis/oidc/token?audience=<your-audience>
   Authorization: Bearer $ACTIONS_ID_TOKEN_REQUEST_TOKEN
   Accept: application/json; api-version=2.0
   ```

   to request the short-lived JWT.
5. GitHub’s OIDC provider returns a JSON payload containing the short-lived JWT in its value field.
6. That JWT is handed to the provider credential factory azidentity.NewClientAssertionCredential() along with the Azure tenant and client IDs.
7. Entra ID’s workload-federation endpoint verifies the JWT’s issuer, audience, expiration, and signature via its JWKS URI.
8. Entra ID issues an OAuth2 access token scoped for Microsoft Graph.
9. The resulting ClientAssertionCredential (implementing azcore.TokenCredential) holds that access token and refreshes it as needed.
10. Terraform calls Microsoft Graph (Intune/M365 APIs), automatically attaching the bearer token from the TokenCredential, without ever storing long-lived Azure secrets in GitHub.

Key benefits of this approach include:

- **No stored secrets**: Long-lived secrets don't need to be stored in GitHub
- **Automatic rotation**: Tokens are short-lived and automatically rotated
- **Conditional access**: Fine-grained control over which workflows can obtain tokens
- **Reduced attack surface**: Eliminates risk of leaked or compromised credentials

```bash
+-------------------+          +--------------------------------------+          +-----------------+          +--------------------+
|                   |          |                                      |          |                 |          |                    |
| GitHub Actions    |          | GitHub OIDC Provider                 |          | Azure AD /      |          | Microsoft Graph    |
| Runner (job)      |          | (token.actions.githubusercontent.com)|          | Microsoft Entra |          | (Intune/M365 APIs) |
|                   |          |                                      |          | ID              |          |                    |
+---------+---------+          +-----------------+--------------------+          +--------+--------+          +--------+-----------+
          |                                      |                                        |                            |
          | (1) permissions:id-token: write      |                                        |                            |
          | triggers prep of OIDC token          |                                        |                            |
          |─────────────────────────────────────►|                                        |                            |
          |                                      |                                        |                            |
          |     (2) inject env vars:             |                                        |                            |
          |     ACTIONS_ID_TOKEN_REQUEST_URL     |                                        |                            |
          |     ACTIONS_ID_TOKEN_REQUEST_TOKEN   |                                        |                            |
          |◄─────────────────────────────────────|                                        |                            |
          |                                      |                                        |                            |
+---------▼---------+                            |                                        |                            |
|                   |                            |                                        |                            |
| Terraform         | (3) GitHubOIDCStrategy.GetCredential()                              |                            |
| Microsoft 365     | reads env vars & determines audience                                |                            |
| Provider          | (default: api://AzureADTokenExchange)                               |                            |
| (plugin)          |                            |                                        |                            |
+---------+---------+                            |                                        |                            |
          |                                      |                                        |                            |
          | (4) HTTP GET request to GitHub OIDC Provider:                                 |                            |
          |     GET $URL?audience=<audience>     |                                        |                            |
          |     Authorization: Bearer $TOKEN     |                                        |                            |
          |     Accept: application/json; api-version=2.0                                 |                            |
          |─────────────────────────────────────►|                                        |                            |
          |                                      |                                        |                            |
          |                                      | (5) GitHub validates request           |                            |
          |                                      | and returns JSON payload with JWT:     |                            |
          |                                      | { "value": "<short-lived JWT>" }       |                            |
          |◄─────────────────────────────────────|                                        |                            |
          |                                      |                                        |                            |
          | (6) azidentity.NewClientAssertionCredential()                                 |                            |
          | passing JWT as client assertion      |                                        |                            |
          | with tenant ID and client ID         |                                        |                            |
          |──────────────────────────────────────────────────────────────────────────────►|                            |
          |                                      |                                        |                            |
          |                                      |      (7) Entra ID's workload-federation|                            |
          |                                      |      endpoint verifies JWT's issuer,   |                            |
          |                                      |      audience, expiration, and         |                            |
          |                                      |      signature via its JWKS URI        |                            |
          |                                      |                                        |                            |
          |                                      |      (8) Entra ID issues an OAuth2     |                            |
          |                                      |      access token scoped for           |                            |
          |                                      |      Microsoft Graph                   |                            |
          |◄──────────────────────────────────────────────────────────────────────────────|                            |
          |                                      |                                        |                            |
          | (9) ClientAssertionCredential        |                                        |                            |
          | (implements azcore.TokenCredential)  |                                        |                            |
          | holds access token and refreshes     |                                        |                            |
          | it as needed                         |                                        |                            |
          |                                      |                                        |                            |
          | (10) Terraform calls Microsoft Graph APIs                                     |                            |
          | automatically attaching bearer token |                                        |                            |
          | Authorization: Bearer <token>        |                                        |                            |
          |───────────────────────────────────────────────────────────────────────────────────────────────────────────►|
          |                                      |                                        |                            |
```

## Setup

### 1. Create an App Registration in Microsoft Entra ID

You can create the app registration using Azure CLI:

```bash
# 1. Variables
export TENANT_ID="00000000-0000-0000-0000-000000000000"
export APP_NAME="terraform-m365-provider"
export GITHUB_ORG="your-github-org"
export GITHUB_REPO="your-github-repo"

# 2. Create the app registration
APP_ID=$(az ad app create \
  --display-name "$APP_NAME" \
  --query appId -o tsv)
# Also grab the object ID (needed for some commands)
APP_OBJECT_ID=$(az ad app show \
  --id "$APP_ID" \
  --query id -o tsv)

echo "✔️ App created. Client (app) ID: $APP_ID, Object ID: $APP_OBJECT_ID"

# 3. Create a service principal for the app
az ad sp create --id "$APP_ID"
echo "✔️ Service principal created"

# 4. Grant Microsoft Graph application-level permissions
#    Example: grant intune device management permissions
az ad app permission add \
  --id "$APP_ID" \
  --api 00000003-0000-0000-c000-000000000000 \
  --api-permissions \
    78145de6-330d-4800-a6ce-494ff2d33d07=Role \
    7a6ee1e7-141e-4cec-ae74-d9db155731ff=Role \
    dc377aa6-52d8-4e23-b271-2a7ae04cedf3=Role \
    9241abd9-d0e6-425a-bd4f-47ba86e767a4=Role \
  --output none
#   Then grant admin consent
az ad app permission admin-consent \
  --id "$APP_ID"
echo "✔️ Graph API permissions granted and consented"
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

When implementing GitHub OIDC authentication with Microsoft 365, you might encounter various errors at different stages of the authentication flow. Here's how to troubleshoot common issues based on where they occur in the process:

### GitHub Actions Environment Issues (Steps 1-2)

- **Missing Environment Variables**
  ```bash
  Error: GetCredential: environment variable ACTIONS_ID_TOKEN_REQUEST_URL not set
  ```
  This occurs when the GitHub Actions environment doesn't have the necessary OIDC environment variables. Ensure your workflow has `permissions: id-token: write` directive properly configured.

- **OIDC Provider Not Available**
  ```bash
  Error: Failed to get token: OIDC provider not available
  ```
  This happens when trying to use GitHub OIDC authentication outside a GitHub Actions workflow. The authentication method only works within GitHub Actions with proper permissions.

### Token Request Issues (Steps 3-5)

- **Token Request Errors**
  ```bash
  Error: Failed to get token: error requesting token from GitHub Actions OIDC provider
  ```
  This can occur when the provider is unable to communicate with GitHub's OIDC endpoint. Check network connectivity and firewall rules that might block outbound connections.

- **Malformed Token Request**
  ```bash
  Error: Bad Request: invalid audience or scope parameter
  ```
  This happens when the audience parameter in the request doesn't match expected formats. Try setting `M365_OIDC_AUDIENCE` environment variable explicitly:
  
  ```bash
  export M365_OIDC_AUDIENCE="api://AzureADTokenExchange"
  ```

### Azure AD/Entra ID Validation Issues (Steps 6-8)

- **Federated Credential Configuration**
  ```bash
  Error: Failed to get token: failed to exchange GitHub Actions OIDC token: oauth2: "invalid_client" "AADSTS7000215: Invalid client secret provided"
  ```
  This occurs when Entra ID attempts to validate the JWT but the federated credential is improperly configured. The error suggests Azure is expecting a client secret rather than accepting the JWT for workload federation. Verify the subject and issuer in your federated credential configuration match what GitHub is providing.
  
  You can check your federated credential configuration using Azure CLI:
  
  ```bash
  az ad app federated-credential list --id $APP_ID -o table
  ```

- **JWT Signature Validation**
  ```bash
  Error: AADSTS700027: Client assertion contains an invalid signature
  ```
  This indicates the JWT signature validation failed during step 7. Check that the Entra ID application has access to the JWKS URI for GitHub's OIDC provider.

- **Invalid Audience in JWT**
  ```bash
  Error: AADSTS700024: Client assertion with invalid audience claim
  ```
  The audience claim in the JWT doesn't match what's expected in the federated credential configuration. Ensure the audience in your federated credential matches the one used in your GitHub OIDC request.
  
  Update your federated credential audience with:
  
  ```bash
  az ad app federated-credential update --id $APP_ID --federated-credential-id $CREDENTIAL_ID --audiences "['api://AzureADTokenExchange']"
  ```

- **Mismatched Subject**
  ```bash
  Error: AADSTS70021: No matching federated identity record found for presented assertion
  ```
  This occurs when the subject claim in the JWT (e.g., `repo:org/repo:ref:refs/heads/main`) doesn't match any federated credential in Entra ID. Double-check that your subject pattern exactly matches your workflow's context.
  
  You can create a new federated credential with the correct subject:
  
  ```bash
  az ad app federated-credential create \
    --id $APP_ID \
    --parameters "{\"name\":\"github-federated-credential\",\"issuer\":\"https://token.actions.githubusercontent.com\",\"subject\":\"repo:${GITHUB_ORG}/${GITHUB_REPO}:ref:refs/heads/main\",\"description\":\"GitHub Actions federated credential\",\"audiences\":[\"api://AzureADTokenExchange\"]}"
  ```

### Authorization Issues (Steps 9-10)

- **Insufficient Permissions**
  ```bash
  Error: Authorization_RequestDenied: Insufficient privileges to complete the operation
  ```
  This happens when the Entra ID app registration doesn't have the necessary Microsoft Graph permissions. Ensure you've added and granted admin consent for all required permissions.
  
  Grant permissions and consent with:
  
  ```bash
  # Add permissions
  az ad app permission add \
    --id $APP_ID \
    --api 00000003-0000-0000-c000-000000000000 \
    --api-permissions 78145de6-330d-4800-a6ce-494ff2d33d07=Role
  
  # Grant admin consent
  az ad app permission admin-consent --id $APP_ID
  ```

- **Token Refresh Issues**
  ```bash
  Error: Failed to refresh token: token has expired and cannot be refreshed
  ```
  The ClientAssertionCredential is unable to refresh the token. This might occur if your workflow runs longer than the token validity period.

### Verification Steps

When troubleshooting, systematically verify each component:

1. **GitHub Actions workflow configuration**
   - Confirm the `permissions: id-token: write` directive is present
   - Verify the workflow is running in the expected context (branch, repository, etc.)

2. **Entra ID app registration**
   - Check that the app ID and tenant ID used in the provider configuration are correct:

     ```bash
     # Get app details
     az ad app show --id $APP_ID --query "{clientId:appId,objectId:id}" -o table
     ```

   - Verify all required Microsoft Graph permissions are granted:

     ```bash
     # List permissions
     az ad app permission list --id $APP_ID -o table
     ```

3. **Federated credential configuration**
   - Ensure the subject pattern exactly matches your workflow's context
   - Verify the issuer is set to `https://token.actions.githubusercontent.com`
   - Confirm the audience matches what's used in your provider configuration

     ```bash
     # List federated credentials
     az ad app federated-credential list --id $APP_ID -o yaml
     ```

4. **Provider configuration**
   - Check that environment variables are correctly set:

     ```bash
     # In your workflow
     env:
      M365_TENANT_ID: "${{ secrets.TENANT_ID }}"
      M365_AUTH_METHOD: "oidc_github"
      M365_CLIENT_ID: "${{ secrets.CLIENT_ID }}"
      M365_OIDC_AUDIENCE: "api://AzureADTokenExchange"
     ```

   - Verify the auth_method is set to "oidc_github"

For persistent issues, enable more detailed logging by setting the environment variable:

```bash
export TF_LOG=DEBUG
```

This will provide more detailed information about each step of the authentication process.
## Additional Resources

- [GitHub OIDC Documentation](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/configuring-openid-connect-in-azure)
- [About security hardening with OpenID Connect](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect#configuring-the-oidc-trust-with-the-cloud)
- [Configure an app to trust an external identity provider](https://learn.microsoft.com/en-us/entra/workload-id/workload-identity-federation-create-trust)
- [Securing GitHub Actions with OpenID Connect](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect)
- [HashiCorp Vault JWT Auth Method](https://developer.hashicorp.com/vault/docs/auth/jwt)
- [Vault with GitHub Actions](https://developer.hashicorp.com/vault/tutorials/app-integration/github-actions)
