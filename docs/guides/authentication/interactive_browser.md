---
page_title: "Authentication with Interactive Browser"
subcategory: "Authentication"
description: |-
  This guide demonstrates how to authenticate with Microsoft 365 using interactive browser authentication.
---

# Authentication with Interactive Browser

The Microsoft 365 provider can authenticate using an interactive browser session. This method automatically opens a web browser to authenticate the user and is ideal for local development scenarios.

## Prerequisites

- A Microsoft Entra ID tenant
- Permissions to create an app registration in your tenant
- Access to a web browser on the machine running Terraform

## Setup

1. [Create an app registration in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity-platform/quickstart-register-app)
2. Configure the app registration:
   - Under "Authentication", add a platform configuration for "Web"
   - Add a redirect URI (e.g., `http://localhost:8888`)
   - Ensure "Access tokens" and "ID tokens" are checked under "Implicit grant and hybrid flows"
3. Add the required API permissions for Microsoft Graph
   - Navigate to "API permissions" in your app registration
   - Click "Add a permission" and select "Microsoft Graph"
   - Choose "Delegated permissions" (interactive browser uses delegated permissions)
   - Add the necessary permissions depending on your intended use
   - Click "Grant admin consent" for these permissions

## Configuration

### Using Terraform Configuration

```terraform
provider "microsoft365" {
  auth_method = "interactive_browser"
  tenant_id   = "00000000-0000-0000-0000-000000000000"
  entra_id_options = {
    client_id    = "00000000-0000-0000-0000-000000000000"
    redirect_url = "http://localhost:8888"
    username     = "user@domain.com"  # Optional login hint
  }
}
```

### Using Environment Variables (Recommended)

```bash
# Set these environment variables before running Terraform
export M365_TENANT_ID="00000000-0000-0000-0000-000000000000"
export M365_AUTH_METHOD="interactive_browser"
export M365_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export M365_REDIRECT_URI="http://localhost:8888"
export M365_USERNAME="user@domain.com"  # Optional
```


## Usage Workflow

When you run Terraform with interactive browser authentication:

1. The provider will automatically open your default web browser
2. You'll be directed to the Microsoft login page
3. After successful authentication, you'll be redirected to the configured redirect URL
4. The browser may display a success message or a blank page
5. Terraform will continue once authentication is complete

## Configuration Options

### Login Hint

You can provide a `username` to pre-populate the login page:

```terraform
entra_id_options = {
  # ... other settings ...
  username = "user@domain.com"
}
```

### Redirect URL

The `redirect_url` must exactly match one of the redirect URIs configured in your app registration:

```terraform
entra_id_options = {
  # ... other settings ...
  redirect_url = "http://localhost:8888"
}
```

## Use Cases

Interactive browser authentication is ideal for:

- Local development
- Testing and troubleshooting
- First-time setup and configuration
- Scenarios where you need to work with user-specific permissions

## Security Considerations

- This method grants permissions based on the authenticated user's privileges
- For automated processes, consider using client secret, certificate, or OIDC authentication instead
- The default token lifetime is one hour
- For shared machines, be cautious as the browser may retain cookies

## Troubleshooting

- **Browser doesn't open**: The provider may not be able to launch a browser automatically. You can manually open the browser and navigate to the URL displayed in the logs.
- **Redirect error**: Ensure the redirect URL in your configuration exactly matches the one in your app registration, including protocol (http/https) and any trailing slashes.
- **Permission denied**: Ensure you've granted admin consent for the required permissions.
- **Browser automation blocked**: Some security tools may block automated browser launching. In these cases, use device code authentication instead.