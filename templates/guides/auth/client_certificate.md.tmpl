---
page_title: "Authentication with Client Certificate"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using a client certificate.
---

# Authentication with Client Certificate

The Microsoft 365 provider can use a Service Principal with a Client Certificate to authenticate to Microsoft 365 services. This is a more secure authentication method compared to client secrets, as certificates are less susceptible to certain types of attacks and can be managed with stronger security controls.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Setup](#setup)
  - [Manual Setup](#manual-setup)
  - [Generate a Certificate](#generate-a-certificate)
  - [Setup Using Terraform](#setup-using-terraform)
  - [Setup Using PowerShell](#setup-using-powershell)
- [Configuration](#configuration)
  - [Using Terraform Configuration](#using-terraform-configuration)
  - [Using Environment Variables](#using-environment-variables-recommended)
- [Certificate Management](#certificate-management)
  - [Certificate Rotation](#certificate-rotation)
  - [Certificate Storage](#certificate-storage)
- [Integration with HashiCorp Vault](#integration-with-hashicorp-vault)
  - [Storing Certificates in Vault](#storing-certificates-in-vault)
  - [Retrieving Certificates During Terraform Runs](#retrieving-certificates-during-terraform-runs)
  - [Vault Agent Integration](#vault-agent-integration)
  - [Security Considerations for Vault Integration](#security-considerations-for-vault-integration) 
- [Security Considerations](#security-considerations)
- [Troubleshooting](#troubleshooting)

## Prerequisites

- A Microsoft Entra ID tenant
- Permissions to create an app registration in your tenant
- Tools for generating certificates (OpenSSL, PowerShell, or other certificate management tools)

## Setup

### Manual Setup

1. [Create an app registration in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app)
2. Add the required API permissions for Microsoft Graph
   - Navigate to "API permissions" in your app registration
   - Click "Add a permission" and select "Microsoft Graph"
   - Choose "Application permissions" for automation scenarios or "Delegated permissions" for user context
   - **Apply least privilege principles**: Only add permissions specific to the resources you need to manage
   - Click "Grant admin consent" for these permissions
3. Upload the client certificate to the app registration
   - Navigate to "Certificates & secrets" in your app registration
   - In the "Certificates" section, click "Upload certificate"
   - Upload your public certificate (.cer or .pem) file
   - Add a description and expiration date (typically matches the certificate's validity)

### Generate a Certificate

Before uploading a certificate to your app registration, you'll need to generate one. There are several ways to do this:

#### Using OpenSSL (Cross-platform)

```bash
# Generate a private key
openssl genrsa -out key.pem 4096

# Generate a certificate signing request (CSR)
openssl req -new -key key.pem -out cert.csr

# Generate a self-signed certificate (valid for 365 days)
openssl x509 -req -days 365 -in cert.csr -signkey key.pem -out cert.pem

# Create a PKCS#12 file that contains both the certificate and private key
openssl pkcs12 -export -out cert.pfx -inkey key.pem -in cert.pem -password pass:YourSecurePassword
```

#### Using PowerShell (Windows)

```powershell
# Create a self-signed certificate
$cert = New-SelfSignedCertificate -Subject "CN=Microsoft365TerraformAuth" `
  -CertStoreLocation "Cert:\CurrentUser\My" `
  -KeyExportPolicy Exportable `
  -KeySpec Signature `
  -KeyLength 2048 `
  -KeyAlgorithm RSA `
  -HashAlgorithm SHA256 `
  -NotAfter (Get-Date).AddYears(1)

# Export the public certificate (for uploading to Azure)
Export-Certificate -Cert $cert -FilePath "cert.cer" -Type CERT

# Export the certificate with private key (for Terraform)
$pwd = ConvertTo-SecureString -String "YourSecurePassword" -Force -AsPlainText
Export-PfxCertificate -Cert $cert -FilePath "cert.pfx" -Password $pwd
```

### Setup Using Terraform

You can create the app registration and configure certificate authentication using Terraform with the AzureAD provider:

```terraform
terraform {
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.47.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.4.0"
    }
  }
}

provider "azuread" {
  # Configuration options
}

# Generate a private key
resource "tls_private_key" "cert_key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

# Generate a self-signed certificate
resource "tls_self_signed_cert" "cert" {
  private_key_pem = tls_private_key.cert_key.private_key_pem

  subject {
    common_name = "Microsoft365TerraformAuth"
  }

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "client_auth",
  ]
}

# Save the certificate to a file (for debugging or reference)
resource "local_file" "cert_file" {
  content  = tls_self_signed_cert.cert.cert_pem
  filename = "${path.module}/cert.pem"
}

# Save the private key to a file (for debugging or reference - secure this file!)
resource "local_file" "key_file" {
  content  = tls_private_key.cert_key.private_key_pem
  filename = "${path.module}/key.pem"
  file_permission = "0600" # Restrict permissions
}

# Create an application registration
resource "azuread_application" "microsoft365_app" {
  display_name = "Microsoft365 Terraform Provider"
}

# Create a service principal for the application
resource "azuread_service_principal" "microsoft365_sp" {
  application_id = azuread_application.microsoft365_app.application_id
}

# Add Microsoft Graph permissions for Microsoft 365 management
resource "azuread_application_api_permission" "graph_permissions" {
  application_id = azuread_application.microsoft365_app.id
  api_id         = "00000003-0000-0000-c000-000000000000" # Microsoft Graph ID

  # Example: Microsoft Intune device management permissions
  api_permissions {
    id   = "78145de6-330d-4800-a6ce-494ff2d33d07" # DeviceManagementApps.ReadWrite.All
    type = "Role"                                 # Application permission
  }
  
  api_permissions {
    id   = "dc377aa6-52d8-4e23-b271-2a7ae04cedf3" # DeviceManagementConfiguration.Read.All
    type = "Role"                                 # Application permission
  }
}

# Grant admin consent (requires an authenticated provider with sufficient permissions)
resource "azuread_service_principal_api_permission_grant" "grant_admin_consent" {
  service_principal_id = azuread_service_principal.microsoft365_sp.id
  api_id               = "00000003-0000-0000-c000-000000000000" # Microsoft Graph ID
  scope                = each.value

  for_each = toset([
    "DeviceManagementApps.ReadWrite.All",
    "DeviceManagementConfiguration.Read.All",
  ])
  
  depends_on = [
    azuread_application_api_permission.graph_permissions
  ]
}

# Add the certificate to the application
resource "azuread_application_certificate" "microsoft365_cert" {
  application_id = azuread_application.microsoft365_app.id
  type           = "AsymmetricX509Cert"
  value          = tls_self_signed_cert.cert.cert_pem
  end_date       = timeadd(timestamp(), "8760h") # 1 year
}

# Output important values
output "tenant_id" {
  value     = data.azuread_client_config.current.tenant_id
  sensitive = false
}

output "client_id" {
  value     = azuread_application.microsoft365_app.application_id
  sensitive = false
}

output "certificate_path" {
  value     = local_file.key_file.filename
  sensitive = false
}
```

After running this Terraform configuration, you'll need to create a PKCS#12 (.pfx) file manually or add additional resources to do this. 
This is because Terraform doesn't have native support for generating PKCS#12 files.

### Setup Using PowerShell

You can create an app registration and configure certificate authentication using PowerShell:

```powershell
# Connect to Microsoft Graph with the ability to create applications
Connect-MgGraph -Scopes "Application.ReadWrite.All", "Directory.ReadWrite.All"

# Create a self-signed certificate
$certName = "Microsoft365TerraformAuth"
$cert = New-SelfSignedCertificate -Subject "CN=$certName" `
    -CertStoreLocation "Cert:\CurrentUser\My" `
    -KeyExportPolicy Exportable `
    -KeySpec Signature `
    -KeyLength 2048 `
    -KeyAlgorithm RSA `
    -HashAlgorithm SHA256 `
    -NotAfter (Get-Date).AddYears(1)

# Export the public certificate for uploading to Azure
$certOutputPath = "$env:USERPROFILE\Documents\$certName.cer"
Export-Certificate -Cert $cert -FilePath $certOutputPath -Type CERT

# Export the certificate with private key for use with Terraform
$pfxOutputPath = "$env:USERPROFILE\Documents\$certName.pfx"
$pfxPassword = "YourSecurePassword" # Change to a secure password
$secPfxPassword = ConvertTo-SecureString -String $pfxPassword -Force -AsPlainText
Export-PfxCertificate -Cert $cert -FilePath $pfxOutputPath -Password $secPfxPassword

# Create a new app registration
$appParams = @{
    DisplayName = "Microsoft365 Terraform Provider"
}
$app = New-MgApplication @appParams

# Create a service principal for the application
$spParams = @{
    AppId = $app.AppId
    DisplayName = $app.DisplayName
}
$sp = New-MgServicePrincipal @spParams

# Read the certificate data for upload to the application
$certData = [System.Convert]::ToBase64String((Get-Content $certOutputPath -Encoding Byte))

# Add the certificate to the application
$certParams = @{
    Type = "AsymmetricX509Cert"
    Usage = "Verify"
    Value = $certData
}
$appCert = New-MgApplicationKey -ApplicationId $app.Id -KeyCredential $certParams

# Define the Microsoft Graph API permissions we need
$graphResourceId = "00000003-0000-0000-c000-000000000000" # Microsoft Graph ID

# Get the app role IDs for the permissions we want
$graphServicePrincipal = Get-MgServicePrincipal -Filter "appId eq '$graphResourceId'"
$permission1 = $graphServicePrincipal.AppRoles | Where-Object { $_.Value -eq "DeviceManagementApps.ReadWrite.All" }
$permission2 = $graphServicePrincipal.AppRoles | Where-Object { $_.Value -eq "DeviceManagementConfiguration.Read.All" }

# Add permissions to the application
$reqResourceAccess = @(
    @{
        ResourceAppId = $graphResourceId
        ResourceAccess = @(
            @{
                Id = $permission1.Id
                Type = "Role"
            },
            @{
                Id = $permission2.Id
                Type = "Role"
            }
        )
    }
)

Update-MgApplication -ApplicationId $app.Id -RequiredResourceAccess $reqResourceAccess

# Grant admin consent (requires admin role)
# Note: This requires additional admin permissions
foreach ($permission in @($permission1.Id, $permission2.Id)) {
    $grantParams = @{
        ClientId    = $sp.Id
        ResourceId  = $graphServicePrincipal.Id
        Scope       = $permission.Value
    }
    # This command requires admin consent rights
    # New-MgOAuth2PermissionGrant @grantParams
    Write-Output "Manual admin consent is required. Please grant consent in the Azure portal."
}

# Output the important information for configuring the Terraform provider
Write-Output "App Registration created successfully"
Write-Output "Tenant ID: $((Get-MgContext).TenantId)"
Write-Output "Client ID: $($app.AppId)"
Write-Output "Certificate Thumbprint: $($cert.Thumbprint)"
Write-Output "Certificate Path: $pfxOutputPath"
Write-Output "Certificate Password: $pfxPassword"
Write-Output "Remember to grant admin consent for the required permissions in the Azure portal"
```

**Note:** The PowerShell script requires the Microsoft Graph PowerShell SDK, which you can install with:

```powershell
Install-Module Microsoft.Graph -Force
```

After running the script, you'll need to manually grant admin consent for the API permissions in the Azure portal, unless your account has sufficient privileges to grant consent programmatically.

## Configuration

After creating the app registration and configuring certificate authentication, configure the Microsoft 365 provider to use them.

### Using Terraform Configuration

```terraform
provider "microsoft365" {
  auth_method = "client_certificate"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id                   = "00000000-0000-0000-0000-000000000000"
    client_certificate          = "/path/to/certificate.pfx"
    client_certificate_password = "YourSecurePassword"
    send_certificate_chain      = true # Optional, set to true if needed
  }
}
```

### Using Environment Variables

```bash
# Set these environment variables before running Terraform
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="client_certificate"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export M365_CLIENT_CERTIFICATE_FILE_PATH="/path/to/certificate.pfx"
export M365_CLIENT_CERTIFICATE_PASSWORD="YourSecurePassword"
export M365_SEND_CERTIFICATE_CHAIN="true" # Optional
```

With environment variables set, your Terraform configuration would still need the provider block:

```terraform
provider "microsoft365" {
  auth_method = "client_certificate"
  # The credentials will be read from environment variables
}
```

**Note:** The provider automatically reads values from environment variables if they're not specified in the provider configuration. When you set the environment variables, the provider will use those values even if they're not explicitly mapped in your Terraform configuration.

## Certificate Management

### Certificate Rotation

Certificates have expiration dates and should be rotated before they expire to avoid disruption to your infrastructure:

1. Generate a new certificate following the methods above
2. Upload the new certificate to your Entra ID app registration
  - Both the old and new certificates can be active simultaneously, allowing for a smooth transition
3. Update your Terraform configuration or environment variables to use the new certificate
4. After confirming the new certificate works, you can optionally remove the old certificate from the app registration

**Best Practice:** Set up a pipeline job to scan regularly to identify certificate expiration to ensure you have sufficient time to perform the rotation.

### Certificate Storage

For secure certificate storage, consider the following options:

1. **Azure Key Vault**: Store your certificate in Azure Key Vault and retrieve it during Terraform execution
2. **HashiCorp Vault**: Store your certificate in HashiCorp Vault's PKI engine
3. **Secure file storage**: Ensure your certificate files have restrictive permissions (e.g., `chmod 600`)
4. **CI/CD Secrets**: Store certificates in your CI/CD pipeline's secure secret storage

Example using Azure Key Vault:

```terraform
# Retrieve certificate from Azure Key Vault
data "azurerm_key_vault" "example" {
  name                = "my-keyvault"
  resource_group_name = "my-resource-group"
}

data "azurerm_key_vault_certificate" "example" {
  name         = "my-certificate"
  key_vault_id = data.azurerm_key_vault.example.id
}

# Configure provider with certificate from Key Vault
provider "microsoft365" {
  auth_method = "client_certificate"
  tenant_id   = var.tenant_id
  entra_id_options = {
    client_id                   = var.client_id
    client_certificate          = data.azurerm_key_vault_certificate.example.certificate_data_base64
    client_certificate_password = "" # If your certificate is password-protected
  }
}
```

## Integration with HashiCorp Vault

HashiCorp Vault provides secure storage for sensitive data, including certificates. This section covers how to store and retrieve certificate-based credentials using Vault.

### Storing Certificates in Vault

You can store certificates in Vault using several methods, depending on your needs:

#### Method 1: Using the Key/Value Secrets Engine

Store certificate data and related credentials:

```bash
# Store tenant_id and client_id as regular KV pairs
vault kv put secret/microsoft365/credentials \
  tenant_id="00000000-0000-0000-0000-000000000000" \
  client_id="00000000-0000-0000-0000-000000000000"

# Store the certificate and password separately
# Note: Base64 encode the .pfx file to store it as a string
vault kv put secret/microsoft365/certificate \
  certificate="$(base64 -i /path/to/certificate.pfx)" \
  password="YourSecurePassword"
```

#### Method 2: Using Vault's PKI Secrets Engine

For more advanced certificate management, you can use Vault's PKI secrets engine:

```bash
# Enable the PKI secrets engine
vault secrets enable pki

# Configure the PKI secrets engine
vault write pki/root/generate/internal \
  common_name="example.com" \
  ttl=8760h

# Create a role for generating certificates
vault write pki/roles/microsoft365 \
  allowed_domains="example.com" \
  allow_subdomains=true \
  max_ttl=8760h

# Generate a certificate
vault write pki/issue/microsoft365 \
  common_name="microsoft365.example.com" \
  ttl=8760h
```

### Retrieving Certificates During Terraform Runs

#### Method 1: Using the Vault Provider with KV Secrets Engine

```terraform
provider "vault" {
  # Vault provider configuration
}

# Get tenant ID and client ID
data "vault_kv_secret_v2" "microsoft365_creds" {
  mount = "secret"
  name  = "microsoft365/credentials"
}

# Get certificate and password
data "vault_kv_secret_v2" "microsoft365_cert" {
  mount = "secret"
  name  = "microsoft365/certificate"
}

# Save the certificate to a temporary file for use
resource "local_file" "certificate" {
  content_base64 = data.vault_kv_secret_v2.microsoft365_cert.data["certificate"]
  filename       = "${path.module}/temp_certificate.pfx"
  file_permission = "0600" # Restrict permissions
}

provider "microsoft365" {
  auth_method = "client_certificate"
  tenant_id   = data.vault_kv_secret_v2.microsoft365_creds.data["tenant_id"]
  entra_id_options = {
    client_id                   = data.vault_kv_secret_v2.microsoft365_creds.data["client_id"]
    client_certificate          = local_file.certificate.filename
    client_certificate_password = data.vault_kv_secret_v2.microsoft365_cert.data["password"]
  }
}
```

#### Method 2: Using External Data Source

```terraform
# Get credentials from Vault
data "external" "vault_creds" {
  program = [
    "/bin/bash", "-c",
    "vault kv get -format=json secret/microsoft365/credentials | jq '{tenant_id: .data.data.tenant_id, client_id: .data.data.client_id}'"
  ]
}

# Get certificate from Vault and save it to a temporary file
data "external" "vault_cert" {
  program = [
    "/bin/bash", "-c",
    "CERT_DATA=$(vault kv get -format=json secret/microsoft365/certificate); echo $CERT_DATA | jq -r '.data.data.certificate' | base64 -d > /tmp/cert.pfx; echo $CERT_DATA | jq '{password: .data.data.password, path: \"/tmp/cert.pfx\"}'"
  ]
}

provider "microsoft365" {
  auth_method = "client_certificate"
  tenant_id   = data.external.vault_creds.result.tenant_id
  entra_id_options = {
    client_id                   = data.external.vault_creds.result.client_id
    client_certificate          = data.external.vault_cert.result.path
    client_certificate_password = data.external.vault_cert.result.password
  }
}
```

#### Method 3: Using Environment Variables with Vault CLI

Create a wrapper script to fetch from Vault and set environment variables:

```bash
#!/bin/bash
# fetch-cert-secrets.sh

# Get credentials from Vault
VAULT_CREDS=$(vault kv get -format=json secret/microsoft365/credentials)
VAULT_CERT=$(vault kv get -format=json secret/microsoft365/certificate)

# Extract values
export M365_TENANT_ID=$(echo $VAULT_CREDS | jq -r '.data.data.tenant_id')
export M365_AUTH_METHOD="client_certificate"
export M365_CLIENT_ID=$(echo $VAULT_CREDS | jq -r '.data.data.client_id')

# Extract and save certificate to a temporary file
CERT_PATH="/tmp/m365_cert_$(date +%s).pfx"
echo $VAULT_CERT | jq -r '.data.data.certificate' | base64 -d > $CERT_PATH
chmod 600 $CERT_PATH
export M365_CLIENT_CERTIFICATE_FILE_PATH=$CERT_PATH
export M365_CLIENT_CERTIFICATE_PASSWORD=$(echo $VAULT_CERT | jq -r '.data.data.password')

# Run terraform with the environment variables set
terraform "$@"

# Clean up the temporary certificate file
rm -f $CERT_PATH
```

Then run:

```bash
chmod +x fetch-cert-secrets.sh
./fetch-cert-secrets.sh apply
```

### Vault Agent Integration

For more advanced scenarios, you can use Vault Agent to automatically retrieve and refresh secrets:

1. Configure Vault Agent template for credentials:

```hcl
template {
  destination = "/path/to/terraform_creds.env"
  contents = <<EOT
  export M365_TENANT_ID={{ "{{" }}with secret "secret/microsoft365/credentials"{{ "}}" }}{{ "{{" }}.Data.data.tenant_id{{ "}}" }}{{ "{{" }}end{{ "}}" }}
  export M365_AUTH_METHOD="client_certificate"
  export M365_CLIENT_ID={{ "{{" }}with secret "secret/microsoft365/credentials"{{ "}}" }}{{ "{{" }}.Data.data.client_id{{ "}}" }}{{ "{{" }}end{{ "}}" }}
  export M365_CLIENT_CERTIFICATE_PASSWORD={{ "{{" }}with secret "secret/microsoft365/certificate"{{ "}}" }}{{ "{{" }}.Data.data.password{{ "}}" }}{{ "{{" }}end{{ "}}" }}
  EOT
}
```

2. Configure Vault Agent to save the certificate:

```hcl
template {
  destination = "/path/to/terraform_cert.pfx"
  contents = "{{ "{{" }}with secret \"secret/microsoft365/certificate\"{{ "}}" }}{{ "{{" }}.Data.data.certificate | base64Decode{{ "}}" }}{{ "{{" }}end{{ "}}" }}"
  perms = 0600
}
```

3. Source the environment file and run Terraform:

```bash
source /path/to/terraform_creds.env
export M365_CLIENT_CERTIFICATE_FILE_PATH=/path/to/terraform_cert.pfx
terraform apply
```

### Security Considerations for Vault Integration

- Use Vault's built-in access controls to restrict who can access the stored certificate
- Consider using response wrapping for sensitive certificate data
- Implement regular rotation of certificates stored in Vault
- Use Vault's audit logging to track access to certificate secrets
- For automation, use Vault's dynamic credentials capabilities where possible
- Clean up temporary certificate files after use

## Security Considerations

- Certificate-based authentication is generally more secure than client secrets
- Use strong passwords for your PKCS#12 (.pfx) files
- Protect private keys and certificates with proper file permissions and secure storage
- Never commit certificates or private keys to version control
- Use certificates with reasonable expiration periods (typically 1-2 years)
- Follow the principle of least privilege for API permissions
- Consider using a hardware security module (HSM) for storing certificate private keys in high-security environments
- For maximum security, consider using certificate chains or certificates issued by a trusted Certificate Authority rather than self-signed certificates

## Troubleshooting

- **Authentication failed**: Verify the tenant ID, client ID, and certificate details are correct
- **Certificate not recognized**: Ensure the certificate is properly formatted and uploaded to the app registration
- **Certificate expired**: Check if your certificate has expired and create a new one if necessary
- **Permission denied**: Ensure you've granted admin consent for the required permissions
- **Certificate format issues**: Ensure your certificate is in the correct format (PKCS#12 for the provider)
- **SendCertificateChain issues**: If you're having problems with authentication, try setting `send_certificate_chain` to `true`

If you're still having issues, examine the logs by setting the debug mode:

```bash
export M365_DEBUG_MODE="true"
```

Or in your Terraform configuration:

```terraform
provider "microsoft365" {
  debug_mode = true
  # Other configuration...
}
```