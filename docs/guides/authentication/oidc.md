---
page_title: "Authentication with Generic OIDC"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using generic OpenID Connect (OIDC) tokens.
---

# Authentication with Generic OIDC

The Microsoft 365 provider supports authentication using generic OpenID Connect (OIDC) tokens. This approach allows for federated authentication from any OIDC-compatible identity provider, enabling secure authentication without managing long-lived secrets.

> [!NOTE]
> Generic OIDC authentication is particularly valuable for custom CI/CD pipelines, enterprise identity systems, and security-focused environments where eliminating static credentials is a priority.

## How Generic OIDC Authentication Works

1. A trusted OIDC provider (identity provider) generates a JWT token
2. This token contains claims about the identity requesting access (issuer, subject, audience)
3. The token is provided to the Microsoft 365 provider via a file or environment variable
4. The provider exchanges this token for a Microsoft Graph access token based on a pre-configured trust relationship in Microsoft Entra ID
5. All API calls use the acquired token for authorization

This process enables a modern, secure authentication flow without the need to manage, rotate, or secure client secrets.

## Prerequisites

- An OIDC token provider capable of generating valid JWTs with the required claims
- Permissions to create and configure app registrations in Microsoft Entra ID
- Ability to configure federated identity credentials
- Terraform provider deploymenttheory/microsoft365 version >= v0.11.0-alpha

## Common Use Cases

Generic OIDC authentication is ideal for:

- **Custom CI/CD Systems**: Self-hosted Jenkins, TeamCity, or other CI/CD platforms that support OIDC
- **Enterprise Identity Systems**: Integration with existing corporate identity providers
- **Security-Focused Environments**: Organizations implementing zero-trust security principles
- **Hybrid Automation**: Scenarios where Terraform runs across diverse environments
- **Custom Workflow Tools**: Internal developer platforms or automation frameworks

## Setup

You can configure the required infrastructure using either the Azure CLI or Terraform.

### Using Azure CLI

#### 1. Create an App Registration

```bash
# Set variables
TENANT_ID="00000000-0000-0000-0000-000000000000"
APP_NAME="terraform-provider-microsoft365"

# Create the app registration
APP_ID=$(az ad app create --display-name $APP_NAME --query appId -o tsv)
APP_OBJECT_ID=$(az ad app show --id $APP_ID --query id -o tsv)

# Create service principal for the application
az ad sp create --id $APP_ID

# Grant API permissions (example for Microsoft Graph)
az ad app permission add \
  --id $APP_ID \
  --api 00000003-0000-0000-c000-000000000000 \
  --api-permissions 9241abd9-d0e6-425a-bd4f-47ba86e767a4=Role

# Grant admin consent
az ad app permission admin-consent --id $APP_ID
```

#### 2. Configure Federated Identity Credential

Configure a federated credential in your Entra ID application. The configuration will depend on your specific OIDC provider.

```bash
# Example for a generic OIDC provider
az ad app federated-credential create \
  --id $APP_OBJECT_ID \
  --parameters "{\"name\":\"generic-oidc-credential\",\"issuer\":\"https://token.issuer.example.com\",\"subject\":\"specific-subject-claim\",\"description\":\"Generic OIDC federated credential\",\"audiences\":[\"api://AzureADTokenExchange\"]}"
```

The key parameters to configure are:

- `issuer`: The OIDC issuer URL of your identity provider
- `subject`: The subject identity you want to trust (varies by provider)
- `audiences`: The intended audience of the token (usually "api://AzureADTokenExchange")

### Using Terraform

```terraform
provider "azurerm" {
  features {}
}

provider "azuread" {}

# Create Microsoft Entra ID application
resource "azuread_application" "terraform_m365" {
  display_name = "terraform-m365-provider"
}

# Create service principal
resource "azuread_service_principal" "terraform_m365" {
  application_id = azuread_application.terraform_m365.application_id
}

# Add API permissions for Microsoft Graph
resource "azuread_application_api_permission" "graph_permissions" {
  application_object_id = azuread_application.terraform_m365.object_id
  
  api_id = "00000003-0000-0000-c000-000000000000" # Microsoft Graph API
  
  # Example: DeviceManagementConfiguration.ReadWrite.All
  api_permissions {
    id   = "9241abd9-d0e6-425a-bd4f-47ba86e767a4"
    type = "Role"
  }
}

# Grant admin consent
resource "azuread_application_api_permission_admin_consent" "graph_permissions" {
  application_object_id = azuread_application.terraform_m365.object_id
}

# Add federated identity credential
resource "azuread_application_federated_identity_credential" "generic_oidc" {
  application_object_id = azuread_application.terraform_m365.object_id
  display_name          = "generic-oidc-credential"
  description           = "Generic OIDC federated credential"
  audiences             = ["api://AzureADTokenExchange"]
  issuer                = "https://token.issuer.example.com"
  subject               = "specific-subject-claim"
}

# Output important values
output "tenant_id" {
  value = data.azuread_client_config.current.tenant_id
}

output "client_id" {
  value = azuread_application.terraform_m365.application_id
}
```

## Generating OIDC Tokens

The specific method for generating OIDC tokens depends on your identity provider. Here are examples for common systems:

### HashiCorp Vault as OIDC Provider

#### 1. Configure Vault as an OIDC Provider

```hcl
# Enable the OIDC identity provider
resource "vault_identity_oidc_provider" "azure" {
  name               = "azure"
  https_enabled      = true
  issuer_host        = "vault.example.com"
  allowed_client_ids = ["*"]
}

# Create a key for signing tokens
resource "vault_identity_oidc_key" "key" {
  name      = "azure-key"
  algorithm = "RS256"
}

# Create a role for issuing tokens
resource "vault_identity_oidc_role" "azure_role" {
  name      = "azure-role"
  key       = vault_identity_oidc_key.key.name
  ttl       = 3600
  
  client_id = "api://AzureADTokenExchange"
  
  template = <<EOF
{
  "iss": "https://vault.example.com",
  "sub": "specific-subject-claim",
  "aud": "api://AzureADTokenExchange"
}
EOF
}

# Assign the role to an entity or group
resource "vault_identity_oidc_key_allowed_client_id" "azure" {
  key_name          = vault_identity_oidc_key.key.name
  allowed_client_id = vault_identity_oidc_role.azure_role.client_id
}
```

#### 2. Generate a Token

```bash
# Generate a token from Vault
AZURE_OIDC_TOKEN=$(vault read -field=token identity/oidc/token/azure-role)

# Save to file
echo $AZURE_OIDC_TOKEN > /path/to/oidc-token.jwt
```

### Custom Identity Server Requirements

For a custom identity server, ensure it supports:

1. JWT token generation with RS256 signing
2. Configuration of issuer, subject, and audience claims
3. Proper key rotation and token validation

The token must include these claims:

- `iss` (issuer): Must match the issuer configured in the federated credential
- `sub` (subject): Must match the subject configured in the federated credential
- `aud` (audience): Typically "api://AzureADTokenExchange"
- `exp` (expiration time): Token expiration timestamp
- `iat` (issued at): Token issuance timestamp

## Microsoft 365 Provider Configuration

### Using Terraform Configuration

```terraform
provider "microsoft365" {
  auth_method = "oidc"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id           = "00000000-0000-0000-0000-000000000000"
    oidc_token_file_path = "/path/to/oidc-token.jwt"
  }
}
```

### Using Environment Variables (Recommended)

```bash
# Set these environment variables before running Terraform
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="oidc"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export M365_OIDC_TOKEN_FILE_PATH="/path/to/oidc-token.jwt"
```

Then your Terraform configuration can be simplified:

```terraform
provider "microsoft365" {
  auth_method = "oidc"
  # No need to specify credentials as they're read from environment variables
}
```

## Token File Format

The OIDC token file should contain a valid JWT token as plain text. For example:

```
eyJhbGciOiJSUzI1NiIsImtpZCI6IkMyRjU2RDU1MkYyQzNCQzg2MDI4MjRCNjA2QkM3NzdDIiwidHlwIjoiSldUIn0.eyJpc3MiOiJodHRwczovL3Rva2VuLmlzc3Vlci5leGFtcGxlLmNvbSIsInN1YiI6InNwZWNpZmljLXN1YmplY3QtY2xhaW0iLCJhdWQiOiJhcGk6Ly9BenVyZUFEVG9rZW5FeGNoYW5nZSIsImV4cCI6MTY5OTEyMzQ1NiwiaWF0IjoxNjk5MTIzMTU2fQ.signature
```

## Using HashiCorp Vault for Secret Management

Even when using OIDC authentication, you might still want to use HashiCorp Vault to manage other sensitive information related to your Microsoft 365 configuration.

### 1. Store Microsoft 365 Configuration in Vault

```terraform
# Store Microsoft 365 configuration in Vault
resource "vault_kv_secret_v2" "m365_config" {
  mount = "secret"
  name  = "microsoft365/config"
  data_json = jsonencode({
    tenant_id = "00000000-0000-0000-0000-000000000000"
    client_id = "00000000-0000-0000-0000-000000000000"
  })
}
```

### 2. Retrieve Configuration and Generate OIDC Token

```bash
#!/bin/bash
# Retrieve configuration from Vault
TENANT_ID=$(vault kv get -field=tenant_id secret/microsoft365/config)
CLIENT_ID=$(vault kv get -field=client_id secret/microsoft365/config)

# Generate OIDC token
OIDC_TOKEN=$(vault read -field=token identity/oidc/token/azure-role)

# Export as environment variables
export M365_TENANT_ID="$TENANT_ID"
export M365_AUTH_METHOD="oidc"
export M365_CLIENT_ID="$CLIENT_ID"

# Either save to file and reference the file path
echo "$OIDC_TOKEN" > /tmp/oidc-token.jwt
export M365_OIDC_TOKEN_FILE_PATH="/tmp/oidc-token.jwt"

# Run Terraform
terraform apply
```

### 3. Using Vault Agent for Automation

Vault Agent can automate the process of fetching secrets and OIDC tokens:

# Vault Agent configuration
auto_auth {
  method "kubernetes" {
    mount_path = "auth/kubernetes"
    config = {
      role = "terraform-role"
    }
  }
}

template {
  destination = "/path/to/env-file.sh"
  contents = <<EOT
  #!/bin/bash
  export M365_TENANT_ID={{ "{{" }}with secret "secret/microsoft365/config"{{ "}}" }}{{ "{{" }}.Data.data.tenant_id{{ "}}" }}{{ "{{" }}end{{ "}}" }}
  export M365_CLIENT_ID={{ "{{" }}with secret "secret/microsoft365/config"{{ "}}" }}{{ "{{" }}.Data.data.client_id{{ "}}" }}{{ "{{" }}end{{ "}}" }}
  export M365_AUTH_METHOD="oidc"

  # Get OIDC token and save to file
  OIDC_TOKEN={{ "{{" }}with secret "identity/oidc/token/azure-role"{{ "}}" }}{{ "{{" }}.Data.token{{ "}}" }}{{ "{{" }}end{{ "}}" }}
  echo "$OIDC_TOKEN" > /tmp/oidc-token.jwt
  export M365_OIDC_TOKEN_FILE_PATH="/tmp/oidc-token.jwt"
EOT
}

## Integration Examples

### Jenkins Pipeline with OIDC

For a Jenkins pipeline that uses OIDC authentication:

```groovy
pipeline {
    agent any
    
    environment {
        M365_TENANT_ID = credentials('m365-tenant-id')
        M365_CLIENT_ID = credentials('m365-client-id')
        M365_AUTH_METHOD = 'oidc'
    }
    
    stages {
        stage('Generate OIDC Token') {
            steps {
                script {
                    // Get OIDC token from your identity provider
                    def token = sh(script: 'curl -s -X POST https://your-identity-provider/token --data "audience=api://AzureADTokenExchange"', returnStdout: true).trim()
                    writeFile file: 'oidc-token.jwt', text: token
                    env.M365_OIDC_TOKEN_FILE_PATH = "${WORKSPACE}/oidc-token.jwt"
                }
            }
        }
        
        stage('Terraform') {
            steps {
                sh 'terraform init'
                sh 'terraform apply -auto-approve'
            }
        }
    }
    
    post {
        always {
            // Clean up the token file
            sh 'rm -f ${M365_OIDC_TOKEN_FILE_PATH}'
        }
    }
}
```

### GitLab CI/CD with OIDC

For GitLab CI/CD pipelines, you can use the GitLab OIDC provider to authenticate with Microsoft 365:

```yaml
variables:
  TF_VAR_tenant_id: ${M365_TENANT_ID}
  TF_VAR_client_id: ${M365_CLIENT_ID}

stages:
  - prepare
  - deploy
  - cleanup

before_script:
  - export M365_AUTH_METHOD="oidc"
  - export M365_TENANT_ID=${TF_VAR_tenant_id}
  - export M365_CLIENT_ID=${TF_VAR_client_id}

generate-oidc-token:
  stage: prepare
  image: alpine:latest
  script:
    # Install dependencies
    - apk add --no-cache curl jq
    
    # Request JWT token from GitLab's OIDC provider
    - >
      TOKEN=$(curl -s --request POST --header "Content-Type:application/json" 
      "${CI_JOB_JWT_URL}" | jq -r .token)
    
    # Save token to file with secure permissions
    - echo "$TOKEN" > oidc-token.jwt
    - chmod 600 oidc-token.jwt
    
    # Export the token file path for later stages
    - echo "M365_OIDC_TOKEN_FILE_PATH=${PWD}/oidc-token.jwt" >> variables.env
  artifacts:
    paths:
      - oidc-token.jwt
    reports:
      dotenv: variables.env

terraform-deploy:
  stage: deploy
  image: hashicorp/terraform:latest
  dependencies:
    - generate-oidc-token
  script:
    - export M365_OIDC_TOKEN_FILE_PATH=${PWD}/oidc-token.jwt
    - terraform init
    - terraform validate
    - terraform plan -out=tfplan
    - terraform apply -auto-approve tfplan

cleanup:
  stage: cleanup
  dependencies:
    - generate-oidc-token
  script:
    # Securely remove the token file
    - rm -f ${PWD}/oidc-token.jwt
  when: always
```

To configure this pipeline:

1. Store your Microsoft 365 tenant ID and client ID as GitLab CI/CD variables (`M365_TENANT_ID` and `M365_CLIENT_ID`)
2. Configure the Microsoft Entra ID federated credential to trust GitLab's OIDC provider
3. Set the subject identifier to match your GitLab project's path or specific job identifiers

For the federated credential in Microsoft Entra ID, use these settings:

- **Issuer**: `https://gitlab.com` (or your self-hosted GitLab instance URL)
- **Subject**: `project_path:<group>/<project>:ref_type:branch:ref:main` (adjust as needed)
- **Audience**: `api://AzureADTokenExchange`

This configuration allows GitLab CI/CD to generate OIDC tokens that Microsoft Entra ID will trust, enabling secure authentication without storing secrets in your GitLab repository or CI/CD variables.

### Custom CI/CD Integration

For custom automation systems, create a wrapper script that:

1. Acquires an OIDC token from your identity provider
2. Sets up the environment variables
3. Invokes Terraform
4. Securely cleans up afterward

```bash
#!/bin/bash
set -e

# Acquire OIDC token
OIDC_TOKEN=$(curl -s -X POST https://your-identity-provider/token --data "audience=api://AzureADTokenExchange")

# Save token to temporary file with secure permissions
TOKEN_FILE=$(mktemp)
echo "$OIDC_TOKEN" > "$TOKEN_FILE"
chmod 600 "$TOKEN_FILE"

# Set environment variables
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="oidc"
export M365_OIDC_TOKEN_FILE_PATH="$TOKEN_FILE"

# Run Terraform
terraform apply -auto-approve

# Clean up
rm -f "$TOKEN_FILE"
```

## Security Considerations

- **Token Lifetime**: OIDC tokens should have a short lifetime (typically under 1 hour)
- **Token Storage**: Protect access to the token file using appropriate file system permissions (chmod 600)
- **Subject Claim Specificity**: Be as specific as possible with the subject claim to limit the scope of trust
- **Conditional Access**: Consider implementing conditional access policies in Microsoft Entra ID
- **Audit Logging**: Enable comprehensive logging for authentication events
- **Claim Validation**: Ensure your federated credential configuration properly validates all required claims
- **Secure Transport**: Always use HTTPS for any API calls that transmit or receive tokens

## Troubleshooting

- **Invalid token**: Ensure the token is valid and not expired
- **Token not found**: Verify the path to the token file is correct
- **Authentication failed**: Check that the issuer, subject, and audience in the token match the federated credential configuration
- **Permission denied**: Ensure you've granted admin consent for the required Microsoft Graph permissions
- **Missing claims**: Verify your token includes all required claims (iss, sub, aud)
- **Claim format issues**: Some claims may need specific formatting; check Microsoft's documentation
- **Token signature validation**: Ensure your token is properly signed with a supported algorithm (RS256)

## Additional Resources

- [Terraform Cloud with workspace OIDC](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/dynamic-provider-credentials)
- [OpenID Connect Core 1.0 specification](https://openid.net/specs/openid-connect-core-1_0.html)
- [HashiCorp Terraform AzureAD Provider documentation](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs)
- [HashiCorp Vault as OIDC Provider](https://developer.hashicorp.com/vault/docs/secrets/identity/oidc-provider)
- [Managing JWT/OIDC tokens with Vault](https://developer.hashicorp.com/vault/docs/secrets/identity/oidc-provider#generate-a-token)
- [Securing CI/CD pipelines with OIDC](https://learn.microsoft.com/en-us/azure/active-directory/develop/workload-identity-federation-create-trust)
