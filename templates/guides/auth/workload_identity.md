---
page_title: "Authentication with Workload Identity"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using workload identity federation for Kubernetes.
---

# Authentication with Workload Identity

The Microsoft 365 provider supports workload identity federation, which allows Kubernetes pods to authenticate to Microsoft 365 services without storing client secrets. This provides a more secure approach for containerized applications running in Kubernetes.

## How Workload Identity Works

Workload identity federation creates a trust relationship between:
1. A Kubernetes service account
2. A Microsoft Entra ID application

Instead of using long-lived secrets, the Kubernetes service account token is exchanged for an Azure access token through the OIDC protocol. This approach:
- Eliminates the need to manage secrets in Kubernetes
- Reduces the risk of credential leakage
- Aligns with zero trust security principles

## Prerequisites

- A Kubernetes cluster with workload identity configured
  - For AKS: [Configure workload identity on AKS](https://learn.microsoft.com/en-us/azure/aks/workload-identity-overview)
  - For other Kubernetes: [Configure Kubernetes federation with Microsoft Entra ID](https://learn.microsoft.com/en-us/azure/active-directory/develop/workload-identity-federation-create-trust-kubernetes)
- Permissions to create/modify:
  - Kubernetes service accounts
  - Microsoft Entra ID applications

## Setup

### 1. Create a Microsoft Entra ID Application

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

### 2. Configure Kubernetes Service Account

```bash
# Create a namespace if it doesn't exist
kubectl create namespace terraform-m365

# Create a service account
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: terraform-m365-sa
  namespace: terraform-m365
  annotations:
    azure.workload.identity/client-id: $APP_ID
EOF
```

### 3. Configure Federated Identity

```bash
# Set Kubernetes variables
SERVICE_ACCOUNT_NAME="terraform-m365-sa"
SERVICE_ACCOUNT_NAMESPACE="terraform-m365"
SERVICE_ACCOUNT_ISSUER="$(kubectl get --raw /.well-known/openid-configuration | jq -r '.issuer')"

# Configure federated identity
az ad app federated-credential create \
  --id $APP_ID \
  --parameters "{\"name\":\"kubernetes-federated-credential\",\"issuer\":\"$SERVICE_ACCOUNT_ISSUER\",\"subject\":\"system:serviceaccount:$SERVICE_ACCOUNT_NAMESPACE:$SERVICE_ACCOUNT_NAME\",\"description\":\"Kubernetes service account federated credential\",\"audiences\":[\"api://AzureADTokenExchange\"]}"
```

## Configuration

### Using Terraform Configuration

```terraform
provider "microsoft365" {
  auth_method = "workload_identity"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id = "00000000-0000-0000-0000-000000000000"
    # Optionally specify the token file path if not using the default location
    federated_token_file_path = "/var/run/secrets/azure/tokens/azure-identity-token"
  }
}
```

### Using Environment Variables (Recommended)

```bash
# Set these environment variables before running Terraform
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="workload_identity"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export AZURE_FEDERATED_TOKEN_FILE="/var/run/secrets/azure/tokens/azure-identity-token"
```

Then your Terraform configuration can be simplified:

```terraform
provider "microsoft365" {}
```

## Use Cases

Workload identity authentication is ideal for:

- Running Terraform in Kubernetes pods
- CI/CD pipelines in Kubernetes-based environments like Tekton or Argo
- Secure, automated workflows without managing secrets
- Production environments following zero trust security principles

## Security Considerations

- The service account token file is automatically mounted in pods associated with the configured service account
- The token is short-lived and automatically rotated by Kubernetes
- Configure RBAC in Kubernetes to limit which pods can use the service account
- Configure conditional access policies in Azure to further restrict access

## Troubleshooting

- **File not found**: Ensure the pod has the correct service account and the token file is mounted correctly
- **Authentication failed**: Verify the federated credential is configured correctly
- **Permission denied**: Ensure you've granted admin consent for the required Microsoft Graph permissions
- **OIDC issuer mismatch**: Verify the OIDC issuer URL in your federated credential configuration matches your Kubernetes cluster's issuer