---
page_title: "Authentication with Managed Identity"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using managed identities for Azure resources.
---

# Authentication with Managed Identity

The Microsoft 365 provider can use Azure managed identities to authenticate to Microsoft 365 services. This authentication method eliminates the need for secrets when Terraform is running on Azure resources, such as Virtual Machines, App Services, Azure Functions, or Azure Container Instances.

## How Managed Identity Authentication Works

Managed identities provide Azure resources with an automatically managed identity in Microsoft Entra ID. When enabled:

1. Azure automatically creates an identity for the resource
2. The resource can request tokens without handling credentials
3. The Microsoft 365 provider can use this identity to authenticate

Managed identities come in two forms:

- **System-assigned**: Tied to the lifecycle of the Azure resource
- **User-assigned**: Created as standalone resources and assigned to one or more Azure resources

## Prerequisites

- An Azure resource with managed identity enabled
- Permissions to manage identities and role assignments in your Azure environment

## Setup

### 1. Enable Managed Identity

#### For Virtual Machine

```bash
# Enable system-assigned managed identity
az vm identity assign --name myVM --resource-group myResourceGroup

# Or, assign a user-assigned managed identity
az vm identity assign --name myVM --resource-group myResourceGroup --identities /subscriptions/SUBSCRIPTION_ID/resourcegroups/RESOURCE_GROUP/providers/Microsoft.ManagedIdentity/userAssignedIdentities/IDENTITY_NAME
```

#### For App Service

```bash
# Enable system-assigned managed identity
az webapp identity assign --name myWebApp --resource-group myResourceGroup

# Or, assign a user-assigned managed identity
az webapp identity assign --name myWebApp --resource-group myResourceGroup --identities /subscriptions/SUBSCRIPTION_ID/resourcegroups/RESOURCE_GROUP/providers/Microsoft.ManagedIdentity/userAssignedIdentities/IDENTITY_NAME
```

### 2. Create an App Registration

```bash
# Create an app registration
APP_ID=$(az ad app create --display-name "Microsoft365Provider" --query appId -o tsv)

# Create service principal
az ad sp create --id $APP_ID

# Add API permissions (Graph API)
az ad app permission add \
  --id $APP_ID \
  --api 00000003-0000-0000-c000-000000000000 \
  --api-permissions PERMISSIONS_LIST_HERE

# Grant admin consent
az ad app permission admin-consent --id $APP_ID
```

### 3. Assign Permissions to Managed Identity

For the managed identity to use the app registration, you need to configure a federated credential.

```bash
# If using a system-assigned managed identity
PRINCIPAL_ID=$(az vm identity show --name myVM --resource-group myResourceGroup --query principalId -o tsv)

# If using a user-assigned managed identity
PRINCIPAL_ID=$(az identity show --name myIdentity --resource-group myResourceGroup --query principalId -o tsv)

# Configure the managed identity to have access to the app registration
az ad app owner add --id $APP_ID --owner-object-id $PRINCIPAL_ID
```

## Configuration

### Using Terraform Configuration

#### System-assigned Managed Identity

```terraform
provider "microsoft365" {
  auth_method = "managed_identity"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  # No additional entra_id_options required for system-assigned identity
}
```

#### User-assigned Managed Identity

```terraform
provider "microsoft365" {
  auth_method = "managed_identity"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    # Specify either the client ID or resource ID of the user-assigned managed identity
    managed_identity_id = "00000000-0000-0000-0000-000000000000" # Client ID
    # OR
    managed_identity_id = "/subscriptions/SUBSCRIPTION_ID/resourceGroups/RESOURCE_GROUP/providers/Microsoft.ManagedIdentity/userAssignedIdentities/IDENTITY_NAME" # Resource ID
  }
}
```

### Using Environment Variables (Recommended)

```bash
# Set these environment variables before running Terraform
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="managed_identity"

# For user-assigned managed identity (optional)
export M365_MANAGED_IDENTITY_ID="00000000-0000-0000-0000-000000000000"
# Or using the Azure SDK environment variable
export AZURE_CLIENT_ID="00000000-0000-0000-0000-000000000000"
```

Then your Terraform configuration can be simplified:

```terraform
provider "microsoft365" {}
```

## Use Cases

Managed identity authentication is ideal for:

- Terraform running on Azure VMs
- Azure DevOps self-hosted agents on Azure resources
- Azure Automation runbooks
- Azure Functions executing Terraform
- Any scenario where Terraform runs on an Azure-hosted resource

## Advantages

- No secrets to manage or rotate
- Azure automatically manages the identity lifecycle
- Reduced risk of credential exposure
- Simplified operations
- Seamless integration with Azure RBAC

## Limitations

- Only works when running Terraform on Azure resources with managed identity support
- Not available for local development unless using Azure CLI login as a fallback
- System-assigned identities are tied to the lifecycle of the Azure resource

## Troubleshooting

- **Managed identity not available**: Verify the managed identity is properly enabled on your Azure resource
- **Identity not authorized**: Ensure the managed identity has the necessary permissions assigned
- **Token acquisition failed**: Check if your resource has network restrictions preventing token acquisition
- **User-assigned identity not found**: Verify the ID format and that the identity exists and is attached to your resource

## Additional Resources

- [What are managed identities for Azure resources?](https://learn.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview)
- [How to use managed identities with Azure VMs](https://learn.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/qs-configure-portal-windows-vm)
- [How to use managed identities with App Service](https://learn.microsoft.com/en-us/azure/app-service/overview-managed-identity)