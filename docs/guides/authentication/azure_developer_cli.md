---
page_title: "Authentication with Azure Developer CLI"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using the Azure Developer CLI (azd).
---

# Authentication with Azure Developer CLI

The Microsoft 365 provider can leverage the Azure Developer CLI (azd) authentication to simplify the development experience. This method uses the existing authenticated session from the Azure Developer CLI, making it ideal for local development scenarios.

## Table of Contents

- [Prerequisites](#prerequisites)
- [How It Works](#how-it-works)
- [Setup](#setup)
  - [Installing Azure Developer CLI](#installing-azure-developer-cli)
  - [Authentication Steps](#authentication-steps)
- [Configuration](#configuration)
  - [Using Terraform Configuration](#using-terraform-configuration)
  - [Using Environment Variables](#using-environment-variables-recommended)
- [Use Cases](#use-cases)
- [Integration with Development Workflows](#integration-with-development-workflows)
  - [Visual Studio Code Integration](#visual-studio-code-integration)
  - [Switching Between Authentication Methods](#switching-between-authentication-methods)
- [Limitations](#limitations)
- [Security Considerations](#security-considerations)
- [Troubleshooting](#troubleshooting)

## Prerequisites

- [Azure Developer CLI (azd)](https://learn.microsoft.com/en-us/azure/developer/azure-developer-cli/install-azd) installed
- Successfully authenticated with `azd auth login`
- A Microsoft Entra ID tenant

## How It Works

This authentication method leverages the existing Azure Developer CLI authentication, which stores tokens in a local credential cache. When you use this method:

1. The provider checks if the Azure Developer CLI is installed and authenticated
2. It uses the existing credential to acquire tokens for Microsoft Graph
3. No additional app registrations or secrets are required

The Azure Developer CLI authentication method simplifies development by:
- Eliminating the need to create and manage separate app registrations
- Removing the need to handle sensitive client secrets or certificates
- Using the same authentication context as your other Azure development tools
- Supporting automatic token renewal when tokens expire

## Setup

### Installing Azure Developer CLI

Before you can use this authentication method, you need to install the Azure Developer CLI:

**For Windows:**
```powershell
# Using winget
winget install Microsoft.Azd

# Using PowerShell
powershell -ex AllSigned -c "Invoke-RestMethod 'https://aka.ms/install-azd.ps1' | Invoke-Expression"
```

**For macOS:**
```bash
# Using Homebrew
brew install azure-developer-cli

# Using curl
curl -fsSL https://aka.ms/install-azd.sh | bash
```

**For Linux:**
```bash
curl -fsSL https://aka.ms/install-azd.sh | bash
```

Verify the installation:
```bash
azd version
```

### Authentication Steps

1. Authenticate with the Azure Developer CLI:
   ```bash
   # Basic authentication to your default tenant
   azd auth login
   
   # Or specify a tenant ID
   azd auth login --tenant-id 00000000-0000-0000-0000-000000000000
   ```

2. Verify authentication is successful:
   ```bash
   azd account show
   ```

3. No additional app registration setup is required for this authentication method

## Configuration

### Using Terraform Configuration

The minimal configuration required:

```terraform
provider "microsoft365" {
  auth_method = "azure_developer_cli"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
}
```

You can also specify additional options:

```terraform
provider "microsoft365" {
  auth_method = "azure_developer_cli"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    additionally_allowed_tenants = ["*"] # Allow multi-tenant access
  }
  debug_mode = true # Enable for troubleshooting
}
```

### Using Environment Variables (Recommended)

Set environment variables before running Terraform:

```bash
# Required variables
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="azure_developer_cli"

# Optional variables
export M365_ADDITIONALLY_ALLOWED_TENANTS="tenant1,tenant2,*"
export M365_DEBUG_MODE="true"
```

With environment variables set, your Terraform configuration would still need the provider block:

```terraform
provider "microsoft365" {
  auth_method = "azure_developer_cli"
  # The credentials will be read from environment variables
}
```

## Use Cases

Azure Developer CLI authentication is ideal for:

- **Local development**: Quick setup for development environments
- **Prototyping**: Rapidly test Terraform configurations without configuration overhead
- **Cross-service development**: Maintain consistent authentication when working with both Azure and Microsoft 365 resources
- **Personal automation**: Scripts and tools that run in the context of your own user account
- **Testing and troubleshooting**: Simplified authentication for debugging issues

This approach is especially valuable for developers who:
- Already use Azure Developer CLI in their workflow
- Need to quickly switch between multiple projects or tenants
- Want to minimize credential management during development
- Prefer to use their own user account permissions for development

## Integration with Development Workflows

### Visual Studio Code Integration

The Azure Developer CLI works seamlessly with Visual Studio Code:

1. Install the [Azure Developer CLI extension](https://marketplace.visualstudio.com/items?itemName=ms-azuretools.azure-dev) for VS Code
2. Use the Azure icon in the activity bar to sign in
3. Run Terraform commands from the integrated terminal after authentication

Example VS Code task configuration:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Terraform Init and Apply",
      "type": "shell",
      "command": "terraform init && terraform apply -auto-approve",
      "options": {
        "env": {
          "M365_AUTH_METHOD": "azure_developer_cli",
          "M365_TENANT_ID": "00000000-0000-0000-0000-000000000000"
        }
      },
      "problemMatcher": []
    }
  ]
}
```

### Switching Between Authentication Methods

During the development lifecycle, you might need to switch between authentication methods:

```bash
# For local development using Azure Developer CLI authentication
export M365_AUTH_METHOD="azure_developer_cli"

# For CI/CD pipelines or production deployments, switch to service principal
export M365_AUTH_METHOD="client_secret"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export M365_CLIENT_SECRET="your-client-secret"
```

Create a shell script to easily switch between profiles:

```bash
#!/bin/bash
# switch-auth-method.sh

case "$1" in
  dev)
    export M365_AUTH_METHOD="azure_developer_cli"
    echo "Switched to developer mode using Azure Developer CLI authentication"
    ;;
  prod)
    export M365_AUTH_METHOD="client_secret"
    export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
    export M365_CLIENT_SECRET="$(cat ~/.secrets/m365_client_secret)"
    echo "Switched to production mode using service principal authentication"
    ;;
  *)
    echo "Usage: $0 {dev|prod}"
    exit 1
    ;;
esac
```

## Limitations

- **Installation requirement**: Requires the Azure Developer CLI to be installed and available in PATH
- **Interactive only**: Only works for interactive development scenarios
- **User permissions**: The authentication uses the permissions of the currently logged-in user
- **Automated environments**: Not suitable for automated workflows or CI/CD pipelines
- **Session expiration**: Sessions might expire, requiring re-authentication
- **Limited control**: Less control over token lifetimes and other authentication parameters
- **Not for production**: Not recommended for production deployment scenarios

## Security Considerations

- Azure Developer CLI stores tokens in a local credential cache, which is encrypted but still exists on the developer's machine
- The authentication uses the developer's own user account, so actions performed will be attributable to that user
- Token refresh is handled automatically, but the initial authentication requires interactive login
- Ensure your user account follows the principle of least privilege
- For shared or public computers, be cautious of leaving authenticated sessions active
- Log out when finished by running `azd auth logout`
- For production environments, consider using a different authentication method with dedicated service principals

## Troubleshooting

- **Azure Developer CLI not found**: 
  ```
  Error: Failed to create credential strategy: azure developer cli not found in PATH
  ```
  Ensure that `azd` is installed and available in your system PATH. Verify with `azd version`.

- **Authentication expired**: 
  ```
  Error: Failed to get token: azure developer cli not authenticated
  ```
  Run `azd auth login` to authenticate before using Terraform.

- **Multiple tenant scenarios**: 
  ```
  Error: Failed to get token: tenant ID mismatch
  ```
  If you work with multiple tenants, authenticate to the specific tenant with:
  ```bash
  azd auth login --tenant-id <tenant-id>
  ```

- **Permission errors**: 
  ```
  Error: Insufficient privileges to complete the operation
  ```
  The authenticated user must have the necessary permissions for Microsoft Graph operations. Check your user permissions in the Microsoft Entra admin center.

- **Debug mode**: Enable debug logging for more detailed information:
  ```terraform
  provider "microsoft365" {
    auth_method = "azure_developer_cli"
    tenant_id   = "00000000-0000-0000-0000-000000000000"
    debug_mode  = true
  }
  ```

- **Version compatibility**: Ensure you're using the latest versions of:
  - Azure Developer CLI: `azd version`
  - Terraform: `terraform version`
  - Microsoft 365 Provider: Check your provider version constraints