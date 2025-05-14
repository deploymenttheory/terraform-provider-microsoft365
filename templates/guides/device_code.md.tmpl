---
page_title: "Authentication with Device Code"
subcategory: "Guides/Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using device code authentication.
---

# Authentication with Device Code

The Microsoft 365 provider can use device code authentication when interactive login isn't possible directly within the application. This is useful for environments without a web browser or where the user cannot directly interact with a login prompt.

## How Device Code Authentication Works

1. The provider requests a device code from Microsoft Entra ID
2. The provider displays a message with:
   - A unique code
   - A URL where the code should be entered
3. The user visits the URL on any device with a browser and enters the code
4. The user authenticates in the browser
5. The provider receives an access token once authentication is complete

## Prerequisites

- A Microsoft Entra ID tenant
- Permissions to create an app registration in your tenant

## Setup

1. [Create an app registration in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app)
2. Configure the app registration:
   - Set "Mobile and desktop applications" as platform type
   - Add `https://login.microsoftonline.com/common/oauth2/nativeclient` as a redirect URI
3. Add the required API permissions for Microsoft Graph
   - Navigate to "API permissions" in your app registration
   - Click "Add a permission" and select "Microsoft Graph"
   - Choose "Delegated permissions" (device code flow requires delegated permissions)
   - Add the necessary permissions depending on your intended use
   - Click "Grant admin consent" for these permissions

## Configuration

### Using Terraform Configuration

```terraform
provider "microsoft365" {
  auth_method = "device_code"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id = "00000000-0000-0000-0000-000000000000"
  }
}
```

### Using Environment Variables (Recommended)

```bash
# Set these environment variables before running Terraform
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="device_code"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
```


## Usage Workflow

When you run Terraform with device code authentication, you'll see a message similar to:

```
To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code ABC123XYZ to authenticate.
```

1. Open the URL in any web browser
2. Enter the code displayed
3. Sign in with your Microsoft credentials
4. Terraform will continue once authentication is complete

## Use Cases

Device code authentication is ideal for:

- CI/CD pipelines with manual intervention
- Environments without a web browser
- Remote terminals or SSH sessions
- Scenarios where redirection to a local web server isn't possible

## Security Considerations

- Device code authentication requires user interaction for each token acquisition
- The default token lifetime is one hour
- For automated processes, consider using client secret, certificate, or OIDC authentication instead
- This authentication method grants permissions based on the authenticated user's privileges

## Troubleshooting

- **Code expired**: If you don't authenticate within the time limit (typically 15 minutes), you'll need to restart the process
- **Permission denied**: Ensure you've granted admin consent for the required permissions
- **No code displayed**: Verify your terminal can display output from the provider