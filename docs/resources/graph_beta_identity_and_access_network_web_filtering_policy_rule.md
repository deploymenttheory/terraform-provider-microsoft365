---
page_title: "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy_rule Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages rules for Microsoft Entra Global Secure Access web filtering policies using the Microsoft Graph beta /networkaccess/webFilteringPolicies/{id}/policyRules endpoint observed from the Entra portal.
---

# microsoft365_graph_beta_identity_and_access_network_web_filtering_policy_rule (Resource)

Manages rules for Microsoft Entra Global Secure Access web filtering policies using the Microsoft Graph beta `/networkaccess/webFilteringPolicies/{id}/policyRules` endpoint observed from the Entra portal.

This resource manages rules for Microsoft Entra Global Secure Access **Web content filtering policies**.

## Microsoft Documentation

Microsoft Graph beta metadata and Microsoft Learn do not currently expose a documented `webFilteringRule` resource type for this portal surface.

The provider implementation is based on Microsoft Entra admin center network traffic observed from the Global Secure Access Web content filtering blade.

## Microsoft Graph API

This resource uses the beta `/networkaccess/webFilteringPolicies/{webFilteringPolicyId}/policyRules` endpoint.

Observed portal traffic for creating a rule uses:

```http
POST /beta/networkaccess/webFilteringPolicies/{webFilteringPolicyId}/policyRules
```

with `@odata.type = #microsoft.graph.networkaccess.webFilteringRule`, `action.@odata.type`, `settings.status`, and `matchingConditions.destinations.targets`.
When custom headers are configured for an allow rule, the portal serializes them under `action.headerSettings.modifications`.
Custom header values cannot contain CR/LF characters or escaped CR/LF sequences.

Destination targets observed from the portal include:

- `#microsoft.graph.networkaccess.webFilteringUrlDestination`
- `#microsoft.graph.networkaccess.webFilteringWebCategoryDestination`

Web category IDs are passed through unchanged. For example, the portal was observed sending `AlcoholAndTobacco` and `AIAgents`.

## Behavior Notes

- At least one destination must be configured using `urls_or_fqdns`, `web_categories`, or both.
- `urls_or_fqdns` values are URL/FQDN patterns without a URL scheme. The Entra portal displays these as comma-separated text, but Microsoft Graph stores them as a values array and Terraform follows that API shape with one set element per destination, for example `["www.MySite.com", "www.MySite.com/a/b/c", "www.MySite.com/a/*", "*.mysite.com"]`. The provider does not validate URL/FQDN pattern syntax locally; unsupported values are returned by Microsoft Graph.
- `web_categories` values are category IDs sent to Microsoft Graph unchanged. The portal category picker displays friendly names, but the API payload uses IDs such as `AIAgents`.
- If `http_methods` is omitted, the rule matches all HTTP methods supported by the service. If `session_types` is omitted, the rule is not limited to a specific source session type.
- Lower `priority` values are evaluated before higher values. Keep priorities unique within a policy to make rule ordering predictable.
- `action = "allow"` permits matching traffic and can add custom request headers. `action = "block"` blocks matching traffic and cannot be combined with `custom_headers`.
- `custom_headers` inserts custom HTTP request headers into matching traffic only when `action = "allow"`. This can be useful for tagging traffic for downstream services or routing decisions, but it should not be used for secrets or sensitive user data because the header can be forwarded to the destination. Some tenants may reject this setting when the backend header modifications feature is not enabled.
- Custom header values cannot include CR/LF characters or common escaped CR/LF sequences. This mirrors the Entra portal validation and prevents header-injection style inputs before the request reaches Graph.

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `NetworkAccess.Read.All`
- `NetworkAccess.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Example Usage

```terraform
resource "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy" "example" {
  name           = "Web Content Filtering Policy"
  description    = "Global Secure Access web filtering policy managed by Terraform"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy_rule" "example" {
  web_filtering_policy_id = microsoft365_graph_beta_identity_and_access_network_web_filtering_policy.example.id

  name        = "Example Web Content Filtering Rule"
  description = "Allow matching web traffic"
  priority    = 100
  action      = "allow"
  status      = "enabled"

  urls_or_fqdns  = ["*.example.com"]
  web_categories = ["AlcoholAndTobacco"]
  http_methods   = ["get"]
  session_types  = ["user", "agent"]
}

resource "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy_rule" "with_headers" {
  web_filtering_policy_id = microsoft365_graph_beta_identity_and_access_network_web_filtering_policy.example.id

  name        = "Allow With Custom Headers"
  description = "Allow matching web traffic and add custom headers"
  priority    = 200
  action      = "allow"
  status      = "enabled"

  urls_or_fqdns = ["headers.example.com"]

  custom_headers = [
    {
      header_name  = "X-Managed-By"
      header_value = "Terraform"
    }
  ]
}

resource "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy_rule" "category_only" {
  web_filtering_policy_id = microsoft365_graph_beta_identity_and_access_network_web_filtering_policy.example.id

  name        = "Block AI Agents Category"
  description = "Block traffic that matches a selected web category"
  priority    = 300
  action      = "block"
  status      = "enabled"

  web_categories = ["AIAgents"]
  session_types  = ["user", "agent"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (String) The action for matching traffic. Possible values are `allow` and `block`.
- `name` (String) The name of the web filtering rule.
- `priority` (Number) The rule priority. Lower numbers are evaluated before higher numbers. The Entra portal accepts values from 100 to 65000.
- `status` (String) The rule status. Possible values are `enabled` and `disabled`.
- `web_filtering_policy_id` (String) The ID of the web filtering policy that owns this rule.

### Optional

- `custom_headers` (Attributes List) Custom request headers to add for allow rules. Microsoft Graph accepts these only when `action` is `allow`; the Entra portal serializes them as `action.headerSettings.modifications`. Some tenants may reject this setting with a BadRequest response when the backend header modifications feature is not enabled. (see [below for nested schema](#nestedatt--custom_headers))
- `description` (String) Optional description of the web filtering rule. Maximum length is 8192 characters.
- `http_methods` (Set of String) HTTP methods that must match the rule. The Entra portal sends these as comma-separated lowercase values.
- `session_types` (Set of String) Session types that must match the rule. Possible values are `user` and `agent`.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `urls_or_fqdns` (Set of String) URL or FQDN destination patterns for the rule, for example `www.MySite.com`, `www.MySite.com/a/b/c`, `www.MySite.com/a/*`, or `*.mysite.com`. Use `*` to match any URL or FQDN. At least one of `urls_or_fqdns` or `web_categories` must be specified; both can be set on the same rule. If set, this attribute must contain at least one value. The Entra portal shows URL/FQDN destinations as comma-separated text, while Microsoft Graph stores them as a values array; Terraform follows the Graph shape with one set element per destination.
- `web_categories` (Set of String) Web category IDs for the rule, for example `AlcoholAndTobacco`. At least one of `urls_or_fqdns` or `web_categories` must be specified; both can be set on the same rule. If set, this attribute must contain at least one value. Category IDs are passed through to Microsoft Graph unchanged.

### Read-Only

- `id` (String) The unique identifier for the web filtering policy rule.

<a id="nestedatt--custom_headers"></a>
### Nested Schema for `custom_headers`

Required:

- `header_name` (String) The custom header name.
- `header_value` (String) The custom header value.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash
# Import using ID format: {web_filtering_policy_id}/{rule_id}
terraform import microsoft365_graph_beta_identity_and_access_network_web_filtering_policy_rule.example 00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111
```
