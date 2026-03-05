---
page_title: "Authentication with Username and Password"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using username and password (ROPC) authentication.
---

# Authentication with Username and Password

The Microsoft 365 provider can use username and password authentication via the Resource Owner Password Credentials (ROPC) flow. This method authenticates directly with a user's credentials without requiring browser interaction.

~> **Security Warning:** ROPC is considered less secure than other authentication methods. It does not support multi-factor authentication (MFA) and is not recommended for production environments. Microsoft recommends using more secure alternatives such as client secret, client certificate, or OIDC authentication. Use this method only when other options are not available.

## How Username/Password Authentication Works

1. The provider sends the user's credentials (username and password) directly to Microsoft Entra ID's token endpoint
2. Microsoft Entra ID validates the credentials
3. If validation succeeds, an access token is issued

## Prerequisites

- A Microsoft Entra ID tenant
- Permissions to create an app registration in your tenant
- A user account with the necessary permissions
- MFA must **not** be enabled for the user account
- The app registration must have **"Allow public client flows"** enabled (see [Setup](#setup) below). Without this setting, authentication will fail with error `AADSTS7000218`

## Setup

1. [Create an app registration in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app)
2. Configure the app registration:
   - Set the application as a public client (under "Authentication" > "Advanced settings" > "Allow public client flows" = Yes)
3. Add the required API permissions for Microsoft Graph
   - Navigate to "API permissions" in your app registration
   - Click "Add a permission" and select "Microsoft Graph"
   - Choose "Delegated permissions" (ROPC flow requires delegated permissions)
   - Add the necessary permissions depending on your intended use
   - Click "Grant admin consent" for these permissions

## Configuration

### Using Terraform Configuration

```terraform
provider "microsoft365" {
  auth_method = "username_password"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id = "00000000-0000-0000-0000-000000000000"
    username  = "user@example.com"
    password  = "your-password"
  }
}
```

### Using Environment Variables (Recommended)

```bash
# Set these environment variables before running Terraform
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="username_password"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export M365_USERNAME="user@example.com"
export M365_PASSWORD="your-password"
```

## Limitations

- **No MFA support**: Accounts with multi-factor authentication enabled cannot use this method
- **No federated accounts**: Federated user accounts (e.g., ADFS) are not supported
- **Conditional Access**: Conditional Access policies requiring device compliance or location-based restrictions may block authentication
- **Deprecated by Microsoft**: Microsoft has deprecated the ROPC flow and recommends alternative authentication methods

## Security Considerations

- Store credentials securely using environment variables or a secrets manager
- Use a dedicated service account with minimal required permissions
- Rotate passwords regularly
- Monitor sign-in logs for suspicious activity
- Consider migrating to client secret, client certificate, or OIDC authentication for production workloads

## Troubleshooting

- **Invalid credentials**: Verify the username (typically UPN format: user@domain.com) and password are correct
- **MFA required**: Disable MFA for the account or switch to a different authentication method
- **Permission denied**: Ensure admin consent has been granted for the required API permissions
- **AADSTS7000218 / Public client not enabled**: The ROPC flow requires the app registration to be configured as a public client. Navigate to your app registration in the Azure portal, go to "Authentication" > "Advanced settings", and set "Allow public client flows" to **Yes**. If this is not enabled, authentication will fail with `AADSTS7000218: The request body must contain the following parameter: 'client_assertion' or 'client_secret'`
- **AADSTS65001 / Consent required**: The ROPC flow uses delegated permissions, which require admin consent. Navigate to "API permissions" in your app registration and click "Grant admin consent"
