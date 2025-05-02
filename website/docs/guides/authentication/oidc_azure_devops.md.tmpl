---
page_title: "Authentication with Azure DevOps OIDC"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using Azure DevOps Pipelines OIDC tokens.
---

# Authentication with Azure DevOps OIDC

The Microsoft 365 provider supports authentication using Azure DevOps Pipelines' OpenID Connect (OIDC) tokens. This approach allows Terraform to authenticate to Microsoft 365 services directly from Azure DevOps pipelines without storing long-lived credentials as pipeline variables or service connections.

## How Azure DevOps OIDC Authentication Works

1. Azure DevOps Pipelines generates a short-lived OIDC token for each pipeline run
2. The Microsoft 365 provider requests this token during Terraform execution
3. The provider presents the token to Microsoft Entra ID
4. Based on a pre-configured trust relationship, Entra ID issues a Microsoft Graph access token
5. The provider uses this token to authenticate API requests

## Prerequisites

- An Azure DevOps organization and project
- Permissions to create and configure app registrations in Microsoft Entra ID
- Ability to modify Azure DevOps pipelines
- Permissions to create service connections in Azure DevOps

## Setup

### 1. Create an App Registration in Microsoft Entra ID

```bash
# Set variables
TENANT_ID="00000000-0000-0000-0000-000000000000"
APP_NAME="terraform-m365-provider"

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

### 2. Create a Federated Credential Service Connection in Azure DevOps

1. In your Azure DevOps project, go to **Project settings** > **Service connections**
2. Click **New service connection** > **Azure Resource Manager**
3. Select **Workload Identity federation (manual)**
4. Fill in the required information:
   - **Subscription ID**: Your Azure subscription ID
   - **Tenant ID**: Your Microsoft Entra ID tenant ID
   - **Service connection name**: A name for your connection (e.g., "M365Provider")
   - **Service principal client ID**: The app registration client ID from step 1
   - **Service principal name**: The name of your app registration
   - **Issuer**: Typically `<your-azure-devops-org>/<your-azure-devops-project>` (this value is shown in the UI)
   - **Subject**: Typically `sc://<your-azure-devops-org>/<your-azure-devops-project>/<service-connection-name>`

5. Complete the service connection setup

### 3. Configure Federated Identity Credential in Entra ID

After creating the service connection, configure a federated credential in your Entra ID application:

```bash
# Set variables
ADO_ORG="your-azure-devops-org"
ADO_PROJECT="your-azure-devops-project"
SERVICE_CONNECTION_NAME="M365Provider"

# Create federated credential
az ad app federated-credential create \
  --id $APP_ID \
  --parameters "{\"name\":\"ado-federated-credential\",\"issuer\":\"https://vstoken.dev.azure.com/${ADO_ORG}\",\"subject\":\"sc://${ADO_ORG}/${ADO_PROJECT}/${SERVICE_CONNECTION_NAME}\",\"description\":\"Azure DevOps Pipeline federated credential\",\"audiences\":[\"api://AzureADTokenExchange\"]}"
```

The exact issuer URL and subject formats may vary. Check the Azure DevOps UI for the specific values to use.

## Azure DevOps Pipeline Configuration

Configure your Azure DevOps pipeline to use OIDC:

```yaml
trigger:
- main

pool:
  vmImage: ubuntu-latest

variables:
  serviceConnectionId: 'M365Provider' # The name of your service connection
  
steps:
- task: TerraformInstaller@0
  inputs:
    terraformVersion: 'latest'

- task: TerraformTaskV3@3
  displayName: 'Terraform init'
  inputs:
    provider: 'azurerm'
    command: 'init'
    backendServiceArm: '$(serviceConnectionId)'
    backendAzureRmResourceGroupName: 'your-resource-group'
    backendAzureRmStorageAccountName: 'your-storage-account'
    backendAzureRmContainerName: 'tfstate'
    backendAzureRmKey: 'terraform.tfstate'

- task: TerraformTaskV3@3
  displayName: 'Terraform apply'
  inputs:
    provider: 'azurerm'
    command: 'apply'
    environmentServiceNameAzureRM: '$(serviceConnectionId)'
  env:
    M365_TENANT_ID: '$(tenantId)'
    M365_AUTH_METHOD: 'oidc_azure_devops'
    M365_CLIENT_ID: '$(clientId)'
    M365_ADO_SERVICE_CONNECTION_ID: '$(serviceConnectionId)'
```

## Provider Configuration

### Using Terraform Configuration

```terraform
provider "microsoft365" {
  auth_method = "oidc_azure_devops"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id                = "00000000-0000-0000-0000-000000000000"
    ado_service_connection_id = "M365Provider"
  }
}
```

### Using Environment Variables (Recommended)

```bash
# Set these environment variables in your Azure DevOps pipeline
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="oidc_azure_devops"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export M365_ADO_SERVICE_CONNECTION_ID="M365Provider"
# You can also use the Azure ARM environment variables
export ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID="M365Provider"
export ARM_OIDC_AZURE_SERVICE_CONNECTION_ID="M365Provider"
```

Then your Terraform configuration can be simplified:

```terraform
provider "microsoft365" {}
```

## Required Pipeline Environment Variables

Azure DevOps automatically sets the following environment variables that are required for OIDC authentication:

- `SYSTEM_ACCESSTOKEN`: The token used to authenticate to Azure DevOps services
- `SYSTEM_OIDCREQUESTURI`: The URI to request OIDC tokens

These variables need to be available to your Terraform commands. For most tasks, these are automatically available, but if you're using custom scripts, you may need to ensure they're passed through.

## Security Best Practices

1. **Use pipeline conditions**:
   - Restrict which branches or paths can trigger the pipeline
   - Consider using environments with approval gates for sensitive operations

2. **Add conditional access policies**:
   - Configure additional conditions in Microsoft Entra ID
   - Restrict access based on IP ranges or other attributes

3. **Limit API permissions**:
   - Grant only the minimum required permissions to the app registration
   - Use application-level permissions rather than delegated permissions for automation

4. **Secure your pipelines**:
   - Enable pipeline validation
   - Implement branch policies
   - Consider using YAML pipelines for better auditability

## Troubleshooting

- **Authentication failed**: Verify the federated credential is configured correctly in Entra ID
- **Service connection issues**: Check the service connection configuration in Azure DevOps
- **Permission denied**: Ensure you've granted admin consent for the required Microsoft Graph permissions
- **Missing environment variables**: Verify that the required environment variables are set and accessible in your pipeline
- **Subject or issuer mismatch**: Double-check that the subject and issuer values match between Azure DevOps and the federated credential in Entra ID

## Additional Resources

- [Azure Pipelines OIDC Documentation](https://learn.microsoft.com/en-us/azure/devops/pipelines/library/connect-to-azure?view=azure-devops#create-an-azure-resource-manager-service-connection-that-uses-workload-identity-federation)
- [Microsoft Entra ID Workload Identity Federation](https://learn.microsoft.com/en-us/azure/active-directory/develop/workload-identity-federation)
- [Azure DevOps Pipeline Security Best Practices](https://learn.microsoft.com/en-us/azure/devops/pipelines/security/overview?view=azure-devops)