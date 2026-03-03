---
page_title: "Authentication with Azure CLI"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using the Azure CLI (az).
---

# Authentication with Azure CLI

The Microsoft 365 provider can leverage the Azure CLI (az) authentication to simplify the authentication experience. This method uses the existing authenticated session from the Azure CLI, making it ideal for local development and service account scenarios.

## Table of Contents

- [Prerequisites](#prerequisites)
- [How It Works](#how-it-works)
- [Use Cases](#use-cases)
- [Setup](#setup)
  - [Installing the Azure CLI](#installing-the-azure-cli)
  - [Authentication Steps](#authentication-steps)
- [Terraform Configuration](#terraform-configuration)
  - [Using Provider](#using-the-provider)
  - [Using Environment Variables](#using-environment-variables-recommended)
- [Integration with Development Workflows](#integration-with-development-workflows)
  - [Visual Studio Code Integration](#visual-studio-code-integration)
  - [Switching Between Authentication Methods](#switching-between-authentication-methods)
- [Limitations](#limitations)
- [Security Considerations](#security-considerations)
- [Troubleshooting](#troubleshooting)

## Prerequisites

- [Azure CLI (az)](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli) installed
- Successfully authenticated with `az login`
- A Microsoft Entra ID tenant

## How It Works

This authentication method leverages the Azure CLI authentication, which stores tokens in a local credential cache. When you use this method:

1. The provider checks if the Azure CLI is installed and authenticated
2. It uses the existing credential to acquire tokens for Microsoft Graph
3. No additional app registrations or secrets are required

The Azure CLI authentication method simplifies development by:

- Eliminating the need to create and manage separate app registrations
- Removing the need to handle sensitive client secrets or certificates
- Using the same authentication context as your other Azure CLI operations
- Supporting automatic token renewal when tokens expire

## Use Cases

Azure CLI authentication is ideal for:

- **Local development**: Quick setup for development environments
- **Service account usage**: Run Terraform with a dedicated service account logged into the Azure CLI
- **Prototyping**: Rapidly test Terraform configurations without configuration overhead
- **Cross-service development**: Maintain consistent authentication when working with both Azure and Microsoft 365 resources
- **Personal automation**: Scripts and tools that run in the context of your own user account
- **Testing and troubleshooting**: Simplified authentication for debugging issues

This approach is especially valuable for developers who:

- Already use the Azure CLI in their daily workflow
- Need to quickly switch between multiple projects or tenants
- Want to minimize credential management during development
- Prefer to use their own user account permissions for development

## Setup

### Installing the Azure CLI

Before you can use this authentication method, you need to install the Azure CLI:

You can follow the installation steps based on your OS type here: [Azure CLI (az) installation guide](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli)

You can verify the installation with:

```bash
az version
```

### Authentication Steps

1. Authenticate with the Azure CLI:

   ```bash
   # Basic authentication to your default tenant
   az login

   # Or specify a tenant ID
   az login --tenant 00000000-0000-0000-0000-000000000000
   ```

2. No additional app registration setup is required for this authentication method

## Terraform Configuration

### Using the Provider

The minimal configuration required:

```terraform
provider "microsoft365" {
  auth_method = "azure_cli"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
}
```

You can also specify additional options:

```terraform
provider "microsoft365" {
  auth_method = "azure_cli"
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
export M365_AUTH_METHOD="azure_cli"

# Optional variables
export M365_ADDITIONALLY_ALLOWED_TENANTS="tenant1,tenant2,*"
export M365_DEBUG_MODE="true"
```

With environment variables set, your Terraform configuration would still need the provider block:

```terraform
provider "microsoft365" {
  auth_method = "azure_cli"
  # The credentials will be read from environment variables
}
```

## Integration with Development Workflows

### VS Code Task Configuration

To streamline your Terraform workflow using Azure CLI authentication in VS Code:

1. Create a `.vscode` directory at the root of your vscode workspace.
2. Create or edit a file called `tasks.json` inside this directory
3. Add the task configuration shown below

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
          "M365_AUTH_METHOD": "azure_cli",
          "M365_TENANT_ID": "your-tenant-id-here"
        }
      },
      "problemMatcher": []
    },
    {
      "label": "Terraform Plan",
      "type": "shell",
      "command": "terraform plan",
      "options": {
        "env": {
          "M365_AUTH_METHOD": "azure_cli",
          "M365_TENANT_ID": "your-tenant-id-here"
        }
      },
      "problemMatcher": []
    },
    {
      "label": "Terraform Destroy",
      "type": "shell",
      "command": "terraform destroy",
      "options": {
        "env": {
          "M365_AUTH_METHOD": "azure_cli",
          "M365_TENANT_ID": "your-tenant-id-here"
        }
      },
      "problemMatcher": []
    }
  ]
}
```

4. Replace "your-tenant-id-here" with your actual Microsoft Entra ID tenant ID
5. Save the file
6. Access these tasks in VS Code by:

Pressing Ctrl+Shift+P (or Cmd+Shift+P on macOS)
Typing "Tasks: Run Task"
Selecting one of your defined tasks



You can extend this with additional customized tasks for your specific workflow needs.

This configuration creates three common Terraform tasks (init+apply, plan, and destroy) that all use the Azure CLI authentication method automatically, without requiring you to set those environment variables manually each time.

### Switching Between Authentication Methods

During the development lifecycle, you might need to switch between authentication methods:

```bash
# For local development using Azure CLI authentication
export M365_AUTH_METHOD="azure_cli"

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
    export M365_AUTH_METHOD="azure_cli"
    echo "Switched to developer mode using Azure CLI authentication"
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

- **Installation requirement**: Requires the Azure CLI to be installed and available in PATH
- **User permissions**: The authentication uses the permissions of the currently logged-in user or service principal
- **Session expiration**: Sessions might expire, requiring re-authentication
- **Limited control**: Less control over token lifetimes and other authentication parameters
- **Not for production**: Not recommended for production deployment scenarios; consider using dedicated service principal authentication methods instead

## Security Considerations

- Azure CLI stores tokens in a local credential cache, which is encrypted but still exists on the developer's machine
- The authentication uses the logged-in identity (user or service principal), so actions performed will be attributable to that identity
- Token refresh is handled automatically, but the initial authentication requires interactive login (or service principal credentials)
- Ensure your user account follows the principle of least privilege
- For shared or public computers, be cautious of leaving authenticated sessions active
- Log out when finished by running `az logout`
- For production environments, consider using a different authentication method with dedicated service principals

## Troubleshooting

- **Azure CLI not found**:
  ```
  Error: Failed to create credential strategy: azure cli not found in PATH
  ```
  Ensure that `az` is installed and available in your system PATH. Verify with `az version`.

- **Authentication expired**:
  ```
  Error: Failed to get token: azure cli not authenticated
  ```
  Run `az login` to authenticate before using Terraform.

- **Multiple tenant scenarios**:
  ```
  Error: Failed to get token: tenant ID mismatch
  ```
  If you work with multiple tenants, authenticate to the specific tenant with:
  ```bash
  az login --tenant <tenant-id>
  ```

- **Permission errors**:

  ```bash
  Error: Insufficient privileges to complete the operation
  ```

  The authenticated identity must have the necessary permissions for Microsoft Graph operations. Check your user or service principal permissions in the Microsoft Entra admin center.

- **Debug mode**: Enable debug logging for more detailed information:

  ```terraform
  provider "microsoft365" {
    auth_method = "azure_cli"
    tenant_id   = "00000000-0000-0000-0000-000000000000"
    debug_mode  = true
  }
  ```

- **Version compatibility**: Ensure you're using the latest versions of:
  - Azure CLI: `az version`
  - Terraform: `terraform version`
  - Microsoft 365 Provider: Check your provider version constraints

  ## Additional Resources

- [Microsoft Graph permissions reference](https://learn.microsoft.com/en-us/graph/permissions-reference)
- [Azure CLI documentation](https://learn.microsoft.com/en-us/cli/azure/)
- [Azure CLI GitHub repository](https://github.com/Azure/azure-cli)
- [Terraform Microsoft 365 Provider examples](https://github.com/hashicorp/terraform-provider-azuread/tree/main/examples)
- [Microsoft Learn: Authenticate to Azure using Azure CLI](https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli)
