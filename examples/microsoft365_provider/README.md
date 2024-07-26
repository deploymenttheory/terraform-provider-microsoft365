<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

No providers.

## Modules

No modules.

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_auth_method"></a> [auth\_method](#input\_auth\_method) | The authentication method to use for the Entra ID application to authenticate the provider. Options: 'device\_code', 'client\_secret', 'client\_certificate', 'interactive\_browser', 'username\_password'. Can also be set using the `M365_AUTH_METHOD` environment variable. | `string` | `"client_secret"` | no |
| <a name="input_client_certificate_base64"></a> [client\_certificate\_base64](#input\_client\_certificate\_base64) | Base64 encoded PKCS#12 certificate bundle. For use when authenticating as a Service Principal using a Client Certificate. Can also be set using the `M365_CLIENT_CERTIFICATE_BASE64` environment variable. | `string` | `""` | no |
| <a name="input_client_certificate_file_path"></a> [client\_certificate\_file\_path](#input\_client\_certificate\_file\_path) | The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate. Can also be set using the `M365_CLIENT_CERTIFICATE_FILE_PATH` environment variable. | `string` | `""` | no |
| <a name="input_client_certificate_password"></a> [client\_certificate\_password](#input\_client\_certificate\_password) | The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate. Can also be set using the `M365_CLIENT_CERTIFICATE_PASSWORD` environment variable. | `string` | `""` | no |
| <a name="input_client_id"></a> [client\_id](#input\_client\_id) | The client ID for the Entra ID application. This ID is generated when you register an application in the Entra ID (Azure AD) and can be found under App registrations > YourApp > Overview. Can also be set using the `M365_CLIENT_ID` environment variable. | `string` | `""` | no |
| <a name="input_client_secret"></a> [client\_secret](#input\_client\_secret) | The client secret for the Entra ID application. This secret is generated in the Entra ID (Azure AD) and is required for authentication flows such as client credentials and on-behalf-of flows. It can be found under App registrations > YourApp > Certificates & secrets. Required for client credentials and on-behalf-of flows. Can also be set using the `M365_CLIENT_SECRET` environment variable. | `string` | `""` | no |
| <a name="input_cloud"></a> [cloud](#input\_cloud) | The cloud to use for authentication and Graph / Graph Beta API requests. Default is `public`. Valid values are `public`, `gcc`, `gcchigh`, `china`, `dod`, `ex`, `rx`. Can also be set using the `M365_CLOUD` environment variable. | `string` | `"public"` | no |
| <a name="input_debug_mode"></a> [debug\_mode](#input\_debug\_mode) | Flag to enable debug mode for the provider. When enabled, the provider will output additional debug information to the console to help diagnose issues. Can also be set using the `M365_DEBUG_MODE` environment variable. | `bool` | `false` | no |
| <a name="input_enable_chaos"></a> [enable\_chaos](#input\_enable\_chaos) | Enable the chaos handler for testing purposes. When enabled, the chaos handler can simulate specific failure scenarios and random errors in API responses to help test the robustness and resilience of the terraform provider against intermittent issues. This is particularly useful for testing how the provider handles various error conditions and ensures it can recover gracefully. Use with caution in production environments. Can also be set using the `M365_ENABLE_CHAOS` environment variable. | `bool` | `false` | no |
| <a name="input_password"></a> [password](#input\_password) | The password for username/password authentication. Can also be set using the `M365_PASSWORD` environment variable. | `string` | `""` | no |
| <a name="input_proxy_url"></a> [proxy\_url](#input\_proxy\_url) | Specifies the URL of the HTTP proxy server. This URL should be in a valid URL format (e.g., 'http://proxy.example.com:8080'). When 'use\_proxy' is enabled, this URL is used to configure the HTTP client to route requests through the proxy. Ensure the proxy server is reachable and correctly configured to handle the network traffic. Can also be set using the `M365_PROXY_URL` environment variable. | `string` | `""` | no |
| <a name="input_redirect_url"></a> [redirect\_url](#input\_redirect\_url) | The redirect URL for interactive browser authentication. Can also be set using the `M365_REDIRECT_URL` environment variable. | `string` | `""` | no |
| <a name="input_telemetry_optout"></a> [telemetry\_optout](#input\_telemetry\_optout) | Flag to indicate whether to opt out of telemetry. Default is `false`. Can also be set using the `M365_TELEMETRY_OPTOUT` environment variable. | `bool` | `false` | no |
| <a name="input_tenant_id"></a> [tenant\_id](#input\_tenant\_id) | The M365 tenant ID for the Entra ID application. This ID uniquely identifies your Entra ID (EID) instance. It can be found in the Azure portal under Entra ID > Properties. Can also be set using the `M365_TENANT_ID` environment variable. | `string` | `""` | no |
| <a name="input_use_proxy"></a> [use\_proxy](#input\_use\_proxy) | Enables the use of an HTTP proxy for network requests. When set to true, the provider will route all HTTP requests through the specified proxy server. This can be useful for environments that require proxy access for internet connectivity or for monitoring and logging HTTP traffic. Can also be set using the `M365_USE_PROXY` environment variable. | `bool` | `false` | no |
| <a name="input_username"></a> [username](#input\_username) | The username for username/password authentication. Can also be set using the `M365_USERNAME` environment variable. | `string` | `""` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->