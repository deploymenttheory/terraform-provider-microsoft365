---
page_title: "Using HashiCorp Vault with Microsoft 365 Provider"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to securely manage Microsoft 365 credentials using HashiCorp Vault.
---

# Using HashiCorp Vault with Microsoft 365 Provider

This guide demonstrates how to securely manage Microsoft 365 provider credentials using HashiCorp Vault. Vault provides a secure way to store, access, and rotate sensitive credentials such as client secrets, certificates, and tokens.

## Prerequisites

- [HashiCorp Vault](https://www.vaultproject.io/) installed and configured
- Vault CLI access with appropriate permissions
- Basic familiarity with Vault concepts (secrets engines, policies, authentication)
- Microsoft 365 credentials already created (app registrations, client secrets, etc.)

## Setting Up Vault for Microsoft 365 Credentials

### 1. Enable the KV Secrets Engine

If not already enabled, set up a Key-Value secrets engine:

```bash
# Enable KV version 2 secrets engine
vault secrets enable -version=2 -path=secret kv
```

### 2. Store Microsoft 365 Credentials

#### Storing Client Secret Credentials

```bash
# Store client secret credentials
vault kv put secret/microsoft365/client-secret \
  tenant_id="00000000-0000-0000-0000-000000000000" \
  client_id="00000000-0000-0000-0000-000000000000" \
  client_secret="your-client-secret"
```

#### Storing Client Certificate Credentials

```bash
# Store certificate path and password
vault kv put secret/microsoft365/client-certificate \
  tenant_id="00000000-0000-0000-0000-000000000000" \
  client_id="00000000-0000-0000-0000-000000000000" \
  client_certificate_path="/path/to/cert.pfx" \
  client_certificate_password="certificate-password"

# Alternatively, store the entire certificate content
# Not recommended for large certificates
vault kv put secret/microsoft365/client-certificate-content \
  tenant_id="00000000-0000-0000-0000-000000000000" \
  client_id="00000000-0000-0000-0000-000000000000" \
  certificate_content="$(base64 /path/to/cert.pfx)" \
  certificate_password="certificate-password"
```

#### Storing OIDC Credentials

```bash
# Store OIDC configuration
vault kv put secret/microsoft365/oidc \
  tenant_id="00000000-0000-0000-0000-000000000000" \
  client_id="00000000-0000-0000-0000-000000000000"
```

### 3. Create a Vault Policy

Create a policy that grants read-only access to the Microsoft 365 secrets:

```hcl
# m365-policy.hcl
path "secret/data/microsoft365/*" {
  capabilities = ["read"]
}
```

Apply the policy:

```bash
vault policy write m365-policy m365-policy.hcl
```

## Retrieving Vault Secrets for Terraform

There are several ways to securely retrieve and use secrets from Vault in Terraform.

### Method 1: Vault Provider

Use the Vault provider to retrieve secrets directly in your Terraform configuration:

```terraform
terraform {
  required_providers {
    vault = {
      source  = "hashicorp/vault"
      version = "~> 3.20.0"
    }
    microsoft365 = {
      source  = "deploymenttheory/microsoft365"
      version = "~> 1.0.0"
    }
  }
}

provider "vault" {
  # Authentication can be configured through environment variables:
  # VAULT_ADDR, VAULT_TOKEN, etc.
}

# Retrieve client secret credentials
data "vault_kv_secret_v2" "m365_client_secret" {
  mount = "secret"
  name  = "microsoft365/client-secret"
}

# Configure Microsoft 365 provider with client secret
provider "microsoft365" {
  auth_method = "client_secret"
  tenant_id   = data.vault_kv_secret_v2.m365_client_secret.data["tenant_id"]
  entra_id_options = {
    client_id     = data.vault_kv_secret_v2.m365_client_secret.data["client_id"]
    client_secret = data.vault_kv_secret_v2.m365_client_secret.data["client_secret"]
  }
}

# ... rest of your Terraform configuration ...
```

### Method 2: Environment Variables with Vault CLI

Create a wrapper script that fetches secrets and runs Terraform:

```bash
#!/bin/bash
# run-terraform.sh

# Fetch client secret credentials
VAULT_DATA=$(vault kv get -format=json secret/microsoft365/client-secret)

# Extract and export as environment variables
export M365_TENANT_ID=$(echo $VAULT_DATA | jq -r '.data.data.tenant_id')
export M365_AUTH_METHOD="client_secret"
export M365_CLIENT_ID=$(echo $VAULT_DATA | jq -r '.data.data.client_id')
export M365_CLIENT_SECRET=$(echo $VAULT_DATA | jq -r '.data.data.client_secret')

# Run Terraform with all arguments passed to this script
terraform "$@"
```

Make the script executable and use it:

```bash
chmod +x run-terraform.sh
./run-terraform.sh apply
```

Your Terraform configuration can then be simplified:

```terraform
provider "microsoft365" {
  auth_method = "client_secret"
  # Credentials will be read from environment variables
}
```

### Method 3: Vault Agent Templates

Vault Agent can automatically fetch and refresh secrets:

1. Configure Vault Agent template:

```hcl
# agent-config.hcl
template {
  destination = "/path/to/m365-credentials.env"
  contents = <<EOT
export M365_TENANT_ID="{{with secret "secret/microsoft365/client-secret"}}{{.Data.data.tenant_id}}{{end}}"
export M365_AUTH_METHOD="client_secret"
export M365_CLIENT_ID="{{with secret "secret/microsoft365/client-secret"}}{{.Data.data.client_id}}{{end}}"
export M365_CLIENT_SECRET="{{with secret "secret/microsoft365/client-secret"}}{{.Data.data.client_secret}}{{end}}"
EOT
}
```

2. Run Vault Agent:

```bash
vault agent -config=agent-config.hcl
```

3. Source the environment file before running Terraform:

```bash
source /path/to/m365-credentials.env
terraform apply
```

## Advanced Configuration

### Certificate-Based Authentication

For certificate-based authentication, you'll need to handle the certificate file:

```terraform
data "vault_kv_secret_v2" "m365_certificate" {
  mount = "secret"
  name  = "microsoft365/client-certificate"
}

# If certificate content is stored in Vault
resource "local_file" "certificate" {
  content_base64 = data.vault_kv_secret_v2.m365_certificate.data["certificate_content"]
  filename       = "${path.module}/cert.pfx"
  file_permission = "0600"  # Restrict access to the file
}

provider "microsoft365" {
  auth_method = "client_certificate"
  tenant_id   = data.vault_kv_secret_v2.m365_certificate.data["tenant_id"]
  entra_id_options = {
    client_id                  = data.vault_kv_secret_v2.m365_certificate.data["client_id"]
    client_certificate         = local_file.certificate.filename
    client_certificate_password = data.vault_kv_secret_v2.m365_certificate.data["certificate_password"]
  }
}
```

### Dynamic Secrets

For more advanced setups, consider using Vault's Azure secrets engine to dynamically generate credentials:

```bash
# Enable Azure secrets engine
vault secrets enable azure

# Configure Azure secrets engine
vault write azure/config \
  subscription_id=your-subscription-id \
  tenant_id=your-tenant-id \
  client_id=your-client-id \
  client_secret=your-client-secret

# Create a role for dynamic credentials
vault write azure/roles/my-role \
  application_object_id=your-app-object-id \
  ttl=1h
```

## Security Best Practices

1. **Secrets Rotation**
   - Regularly rotate client secrets in both Azure and Vault
   - Use Vault's built-in TTL functionality to enforce rotation

2. **Limit Access**
   - Use specific Vault policies to limit which users/systems can access Microsoft 365 credentials
   - Implement Vault response wrapping for additional security

3. **Audit and Monitoring**
   - Enable Vault audit logging
   - Monitor access to Microsoft 365 credentials
   - Set up alerts for unusual access patterns

4. **CI/CD Integration**
   - When using CI/CD pipelines, use Vault's JWT/OIDC auth method to authenticate CI/CD systems
   - Avoid storing tokens as persistent environment variables in CI/CD systems

## Example: GitHub Actions Workflow with Vault

```yaml
name: Terraform with Vault

on:
  push:
    branches: [ main ]

jobs:
  terraform:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v3

      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v2
        with:
          terraform_version: "1.5.0"

      - name: Vault Login
        uses: hashicorp/vault-action@v2
        with:
          url: ${{ secrets.VAULT_ADDR }}
          method: approle
          roleId: ${{ secrets.VAULT_ROLE_ID }}
          secretId: ${{ secrets.VAULT_SECRET_ID }}
          secrets: |
            secret/data/microsoft365/client-secret tenant_id | M365_TENANT_ID ;
            secret/data/microsoft365/client-secret client_id | M365_CLIENT_ID ;
            secret/data/microsoft365/client-secret client_secret | M365_CLIENT_SECRET

      - name: Terraform Init
        run: terraform init

      - name: Terraform Apply
        run: |
          export M365_AUTH_METHOD="client_secret"
          terraform apply -auto-approve
```

## Troubleshooting

- **Vault authentication issues**: Check that your Vault token has the correct permissions and hasn't expired
- **Secret not found**: Verify the path to your secret and that it exists in Vault
- **Environment variables not set**: Ensure environment variables are correctly exported and available in the Terraform process
- **Permission denied**: Confirm your Vault policy grants read access to the secrets
- **Certificate issues**: When using certificates, check file permissions and path validity

## Additional Resources

- [HashiCorp Vault Documentation](https://developer.hashicorp.com/vault/docs)
- [Vault Provider for Terraform](https://registry.terraform.io/providers/hashicorp/vault/latest/docs)
- [Securing Credentials with Environment Variables](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables)
- [Microsoft 365 Provider Documentation](https://registry.terraform.io/providers/deploymenttheory/microsoft365/latest/docs)