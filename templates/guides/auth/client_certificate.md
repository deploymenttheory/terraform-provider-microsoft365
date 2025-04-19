---
page_title: "Authentication with Client Certificate"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using a client certificate.
---

# Authentication with Client Certificate

The Microsoft 365 provider can use certificate-based authentication for enhanced security when connecting to Microsoft 365 services. This authentication method is more secure than client secret-based authentication as it eliminates the need to manage and rotate secrets.

## Prerequisites

- A Microsoft Entra ID tenant
- Permissions to create an app registration in your tenant
- OpenSSL or similar tool to generate certificates

## Setup

1. [Create an app registration in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app)
2. Add the required API permissions for Microsoft Graph
   - Navigate to "API permissions" in your app registration
   - Click "Add a permission" and select "Microsoft Graph"
   - Choose "Application permissions" for automation scenarios or "Delegated permissions" for user context
   - Add the necessary permissions depending on your intended use
   - Click "Grant admin consent" for these permissions
3. Generate a certificate using OpenSSL or other tools:

   ```bash
   # Generate a private key and certificate
   openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365 -nodes
   
   # Merge public and private parts of the certificate files
   # Using Linux/macOS
   cat cert.pem key.pem > cert+key.pem
   
   # Using PowerShell
   # Get-Content .\cert.pem, .\key.pem | Set-Content cert+key.pem
   
   # Generate PKCS#12 file
   openssl pkcs12 -export -out cert.pfx -inkey key.pem -in cert.pem -passout pass:YourSecurePassword
   ```

4. Upload the public certificate (`cert.pem`) to your app registration
   - Navigate to "Certificates & secrets" in your app registration
   - Click "Upload certificate"
   - Browse and select your `cert.pem` file
   - Add a description and click "Add"

## Configuration

### Using Terraform Configuration

```terraform
provider "microsoft365" {
  auth_method = "client_certificate"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id                   = "00000000-0000-0000-0000-000000000000"
    client_certificate          = "${path.cwd}/cert.pfx"
    client_certificate_password = "YourSecurePassword"
    send_certificate_chain      = false # Set to true for Subject Name/Issuer (SNI) authentication
  }
}
```

### Using Environment Variables (Recommended)

```bash
# Set these environment variables before running Terraform
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="client_certificate"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export M365_CLIENT_CERTIFICATE_FILE_PATH="/path/to/cert.pfx"
export M365_CLIENT_CERTIFICATE_PASSWORD="YourSecurePassword"
export M365_SEND_CERTIFICATE_CHAIN="false" # Optional
```

## Advanced Options

### Certificate Chain

Set `send_certificate_chain` to `true` if you need to send the certificate chain for:

- Subject Name/Issuer (SNI) authentication
- Scenarios where intermediate certificates need to be validated
- When specifically required by your Azure configuration

```terraform
entra_id_options = {
  # ... other settings ...
  send_certificate_chain = true
}
```

## Security Considerations

- Store certificates securely and restrict access to the private key
- Use strong certificate passwords
- Set appropriate certificate validity periods (typically 1-2 years)
- Implement a certificate rotation process before expiration
- Consider using a hardware security module (HSM) for additional security

## Troubleshooting

- **Certificate not found**: Verify the path to your certificate file is correct
- **Invalid certificate format**: Ensure you're using a valid PKCS#12 (.pfx or .p12) file
- **Authentication failed**: Verify the certificate has been properly uploaded to your app registration
- **Permission denied**: Ensure you've granted admin consent for the required permissions