<!-- BEGIN_TF_DOCS -->
### Providers

No providers.

### Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_additionally_allowed_tenants"></a> [additionally_allowed_tenants](#input_additionally_allowed_tenants) | Specifies additional tenants for which the credential may acquire tokens. Can also be set using the `M365_ADDITIONALLY_ALLOWED_TENANTS` environment variable. | `list(string)` | `[]` | no |
| <a name="input_auth_method"></a> [auth_method](#input_auth_method) | The authentication method to use for the Entra ID application to authenticate the provider. Options: 'device_code', 'client_secret', 'client_certificate', 'interactive_browser', 'username_password'. Can also be set using the `M365_AUTH_METHOD` environment variable. | `string` | `"client_certificate"` | no |
| <a name="input_chaos_percentage"></a> [chaos_percentage](#input_chaos_percentage) | Percentage of requests to apply chaos testing to. Must be between 0 and 100. Can also be set using the `M365_CHAOS_PERCENTAGE` environment variable. | `number` | `10` | no |
| <a name="input_chaos_status_code"></a> [chaos_status_code](#input_chaos_status_code) | HTTP status code to return for chaos-affected requests. If not set, a random error status code will be used. Can also be set using the `M365_CHAOS_STATUS_CODE` environment variable. | `number` | `0` | no |
| <a name="input_chaos_status_message"></a> [chaos_status_message](#input_chaos_status_message) | Custom status message to return for chaos-affected requests. If not set, a default message will be used. Can also be set using the `M365_CHAOS_STATUS_MESSAGE` environment variable. | `string` | `""` | no |
| <a name="input_client_certificate"></a> [client_certificate](#input_client_certificate) | The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate. Can also be set using the `M365_CLIENT_CERTIFICATE_FILE_PATH` environment variable. | `string` | `""` | no |
| <a name="input_client_certificate_password"></a> [client_certificate_password](#input_client_certificate_password) | The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate. Can also be set using the `M365_CLIENT_CERTIFICATE_PASSWORD` environment variable. | `string` | `""` | no |
| <a name="input_client_id"></a> [client_id](#input_client_id) | The client ID for the Entra ID application. This ID is generated when you register an application in the Entra ID (Azure AD) and can be found under App registrations > YourApp > Overview. Can also be set using the `M365_CLIENT_ID` environment variable. | `string` | `""` | no |
| <a name="input_client_secret"></a> [client_secret](#input_client_secret) | The client secret for the Entra ID application. This secret is generated in the Entra ID (Azure AD) and is required for authentication flows such as client credentials and on-behalf-of flows. It can be found under App registrations > YourApp > Certificates & secrets. Required for client credentials and on-behalf-of flows. Can also be set using the `M365_CLIENT_SECRET` environment variable. | `string` | `""` | no |
| <a name="input_cloud"></a> [cloud](#input_cloud) | The cloud to use for authentication and Graph / Graph Beta API requests. Default is `public`. Valid values are `public`, `gcc`, `gcchigh`, `china`, `dod`, `ex`, `rx`. Can also be set using the `M365_CLOUD` environment variable. | `string` | `"public"` | no |
| <a name="input_custom_user_agent"></a> [custom_user_agent](#input_custom_user_agent) | Custom User-Agent string to be sent with requests. Can also be set using the `M365_CUSTOM_USER_AGENT` environment variable. | `string` | `""` | no |
| <a name="input_debug_mode"></a> [debug_mode](#input_debug_mode) | Flag to enable debug mode for the provider. When enabled, the provider will output additional debug information to the console to help diagnose issues. Can also be set using the `M365_DEBUG_MODE` environment variable. | `bool` | `true` | no |
| <a name="input_disable_instance_discovery"></a> [disable_instance_discovery](#input_disable_instance_discovery) | Disables the instance discovery in disconnected clouds or private clouds. Can also be set using the `M365_DISABLE_INSTANCE_DISCOVERY` environment variable. | `bool` | `false` | no |
| <a name="input_enable_chaos"></a> [enable_chaos](#input_enable_chaos) | Enable the chaos handler for testing purposes. Can also be set using the `M365_ENABLE_CHAOS` environment variable. | `bool` | `false` | no |
| <a name="input_enable_compression"></a> [enable_compression](#input_enable_compression) | Enable compression for HTTP requests and responses. Can also be set using the `M365_ENABLE_COMPRESSION` environment variable. | `bool` | `true` | no |
| <a name="input_enable_headers_inspection"></a> [enable_headers_inspection](#input_enable_headers_inspection) | Enable inspection of HTTP headers. Can also be set using the `M365_ENABLE_HEADERS_INSPECTION` environment variable. | `bool` | `false` | no |
| <a name="input_enable_redirect"></a> [enable_redirect](#input_enable_redirect) | Enable automatic following of redirects. Can also be set using the `M365_ENABLE_REDIRECT` environment variable. | `bool` | `true` | no |
| <a name="input_enable_retry"></a> [enable_retry](#input_enable_retry) | Enable automatic retries for failed requests. Can also be set using the `M365_ENABLE_RETRY` environment variable. | `bool` | `true` | no |
| <a name="input_max_redirects"></a> [max_redirects](#input_max_redirects) | Maximum number of redirects to follow. Can also be set using the `M365_MAX_REDIRECTS` environment variable. | `number` | `5` | no |
| <a name="input_max_retries"></a> [max_retries](#input_max_retries) | Maximum number of retries for failed requests. Can also be set using the `M365_MAX_RETRIES` environment variable. | `number` | `3` | no |
| <a name="input_password"></a> [password](#input_password) | The password for username/password authentication. Can also be set using the `M365_PASSWORD` environment variable. | `string` | `""` | no |
| <a name="input_proxy_password"></a> [proxy_password](#input_proxy_password) | Password for proxy authentication. Can also be set using the `M365_PROXY_PASSWORD` environment variable. | `string` | `""` | no |
| <a name="input_proxy_url"></a> [proxy_url](#input_proxy_url) | Specifies the URL of the HTTP proxy server. Can also be set using the `M365_PROXY_URL` environment variable. | `string` | `""` | no |
| <a name="input_proxy_username"></a> [proxy_username](#input_proxy_username) | Username for proxy authentication. Can also be set using the `M365_PROXY_USERNAME` environment variable. | `string` | `""` | no |
| <a name="input_redirect_url"></a> [redirect_url](#input_redirect_url) | The redirect URL for interactive browser authentication. Can also be set using the `M365_REDIRECT_URL` environment variable. | `string` | `""` | no |
| <a name="input_retry_delay_seconds"></a> [retry_delay_seconds](#input_retry_delay_seconds) | Delay between retry attempts in seconds. Can also be set using the `M365_RETRY_DELAY_SECONDS` environment variable. | `number` | `5` | no |
| <a name="input_send_certificate_chain"></a> [send_certificate_chain](#input_send_certificate_chain) | Controls whether the credential sends the public certificate chain in the x5c header of each token request's JWT. Can also be set using the `M365_SEND_CERTIFICATE_CHAIN` environment variable. | `bool` | `false` | no |
| <a name="input_telemetry_optout"></a> [telemetry_optout](#input_telemetry_optout) | Flag to indicate whether to opt out of telemetry. Default is `false`. Can also be set using the `M365_TELEMETRY_OPTOUT` environment variable. | `bool` | `false` | no |
| <a name="input_tenant_id"></a> [tenant_id](#input_tenant_id) | The M365 tenant ID for the Entra ID application. This ID uniquely identifies your Entra ID (EID) instance. It can be found in the Azure portal under Entra ID > Properties. Can also be set using the `M365_TENANT_ID` environment variable. | `string` | `"2fd6bb84-1234-abcd-9369-1235b25c1234"` | no |
| <a name="input_timeout_seconds"></a> [timeout_seconds](#input_timeout_seconds) | Timeout for requests in seconds. Can also be set using the `M365_TIMEOUT_SECONDS` environment variable. | `number` | `300` | no |
| <a name="input_use_proxy"></a> [use_proxy](#input_use_proxy) | Enables the use of an HTTP proxy for network requests. Can also be set using the `M365_USE_PROXY` environment variable. | `bool` | `false` | no |
| <a name="input_username"></a> [username](#input_username) | The username for username/password authentication. Can also be set using the `M365_USERNAME` environment variable. | `string` | `""` | no |

### Outputs

No outputs.
<!-- END_TF_DOCS -->