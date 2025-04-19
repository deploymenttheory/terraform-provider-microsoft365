---
page_title: "Authentication with Azure Developer CLI"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using the Azure Developer CLI (azd).
---

# Authentication with Azure Developer CLI

The Microsoft 365 provider can leverage the Azure Developer CLI (azd) authentication to simplify the development experience. This method uses the existing authenticated session from the Azure Developer CLI, making it ideal for local development scenarios.

## Prerequisites

- [Azure Developer CLI (azd)](https://learn.microsoft.com/en-us/azure/developer/azure-developer-cli/install-azd) installed
- Successfully authenticated with `azd auth login`
- A Microsoft Entra ID tenant

## How It Works

This authentication method leverages the existing Azure Developer CLI authentication, which stores tokens in a local credential cache. When you use this method:

1. The provider checks if the Azure Developer CLI is installed and authenticated
2. It uses the existing credential to acquire tokens for Microsoft Graph
3. No additional app registrations or secrets are required

## Setup

1. Install the Azure Developer CLI:
   ```bash
   # For Windows (using winget)
   winget install Microsoft.Azd
   
   # For macOS (using Homebrew)
   brew install azure-developer-cli
   
   # For Linux
   curl -fsSL https://aka.ms/install-azd.sh | bash
   ```

2. Authenticate with the Azure Developer CLI:
   ```bash
   azd auth login
   ```

3. No additional app registration setup is required

## Configuration

### Using Terraform Configuration

```terraform
provider "microsoft365" {
  auth_method = "azure_developer_cli"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  # No additional entra_id_options required
}
```

### Using Environment Variables (Recommended)

```bash
# Set these environment variables before running Terraform
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="azure_developer_cli"
```

Then your Terraform configuration can be simplified:

```terraform
provider "microsoft365" {}
```

## Use Cases

Azure Developer CLI authentication is ideal for:

- Local development
- Testing and troubleshooting
- Developer environments where `azd` is already in use
- Streamlined workflows without managing separate credentials

## Limitations

- Requires the Azure Developer CLI to be installed and available in PATH
- Only works for interactive development scenarios
- The authentication uses the permissions of the currently logged-in user
- Not suitable for automated workflows or CI/CD pipelines

## Troubleshooting

- **Azure Developer CLI not found**: Ensure that `azd` is installed and available in your PATH
- **Authentication error**: Run `azd auth login` to authenticate before using Terraform
- **Multiple tenant scenarios**: If you work with multiple tenants, make sure you've authenticated to the correct tenant with `azd auth login --tenant-id <tenant-id>`
- **Permission errors**: The authenticated user must have the necessary permissions for Microsoft Graph operations