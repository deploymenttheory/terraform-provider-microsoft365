---
page_title: "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages Microsoft Entra Global Secure Access Internet Access forwarding policy rules using Microsoft Graph beta /networkAccess/forwardingPolicies/{forwardingPolicyId}/policyRules.
---

# microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule (Resource)

Manages Microsoft Entra Global Secure Access Internet Access forwarding policy rules using Microsoft Graph beta `/networkAccess/forwardingPolicies/{forwardingPolicyId}/policyRules`.

This resource manages Microsoft Entra Global Secure Access Internet Access forwarding policy rules.

## Microsoft Documentation

- [What is Global Secure Access?](https://learn.microsoft.com/en-us/entra/global-secure-access/overview-what-is-global-secure-access)
- [Internet Access](https://learn.microsoft.com/en-us/entra/global-secure-access/concept-internet-access)
- [List policyRules](https://learn.microsoft.com/en-us/graph/api/networkaccess-policy-list-policyrules?view=graph-rest-beta)
- [Update policyRule](https://learn.microsoft.com/en-us/graph/api/networkaccess-policyrule-update?view=graph-rest-beta)

## Microsoft Graph API

This resource uses Microsoft Graph beta `/networkAccess/forwardingPolicies/{forwardingPolicyId}/policyRules`.

The initial supported shape is `#microsoft.graph.networkaccess.internetAccessForwardingRule` with `rule_type` values for FQDN, IP address, IP range, and CIDR/IP subnet destinations. The provider sends explicit Graph `@odata.type` values for the rule and destinations.

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `NetworkAccess.Read.All`
- `NetworkAccess.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Behavior Notes

- Use the forwarding policy `policy_id` from the forwarding profile data source as `forwarding_policy_id`.
- `action` supports `forward` and `bypass`. Graph requires `forward` rules on acquire policies such as `Custom Acquire`, and `bypass` rules on bypass policies such as `Custom bypass`.
- `name`, `action` and `client_fallback_action` are create-time values for this Graph resource. Changing any of them requires Terraform to replace the rule. `ports`, `protocol` and `destinations` are updated in place.
- `protocol` supports `tcp` and `udp`. To match both, create one rule per protocol.
- `ports` entries are single ports (`443`) or inclusive ranges (`1000-2000`) between 0 and 65535.
- Destination `type` supports `fqdn`, `ip_address`, `ip_range`, and `ip_subnet`, and must equal `rule_type`. Use CIDR notation in `value` for `ip_subnet`. Graph accepts only IPv4 addresses, ranges and subnets for Internet Access rules.
- The generated Microsoft Graph beta Go SDK exposes policy rule builders, but this resource uses Kiota `RequestInformation` and custom `Parsable` request/response types for the observed polymorphic Internet Access payload.

## Example Usage

```terraform
data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link" "custom_acquire" {
  traffic_forwarding_type = "internet"
  policy_name             = "Custom Acquire"
}

data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link" "custom_bypass" {
  traffic_forwarding_type = "internet"
  policy_name             = "Custom bypass"
}

resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "fqdn" {
  forwarding_policy_id = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.custom_acquire.policy_id

  name      = "Example Internet Access FQDN rule"
  action    = "forward"
  rule_type = "fqdn"
  ports     = ["80", "443", "8000-8100"]
  protocol  = "tcp"

  # Set only on create; changing it replaces the rule.
  client_fallback_action = "block"

  destinations = [
    {
      type  = "fqdn"
      value = "example.com"
    }
  ]
}

resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "cidr" {
  forwarding_policy_id = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.custom_bypass.policy_id

  name      = "Example Internet Access CIDR bypass rule"
  action    = "bypass"
  rule_type = "ip_subnet"
  ports     = ["443"]
  protocol  = "tcp"

  destinations = [
    {
      type  = "ip_subnet"
      value = "192.0.2.0/24"
    }
  ]
}

resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "ip_range" {
  forwarding_policy_id = data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link.custom_bypass.policy_id

  name      = "Example Internet Access IP range bypass rule"
  action    = "bypass"
  rule_type = "ip_range"
  ports     = ["443"]
  protocol  = "udp"

  destinations = [
    {
      type          = "ip_range"
      begin_address = "192.0.2.10"
      end_address   = "192.0.2.20"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (String) The forwarding action. Possible values are `forward` and `bypass`. Graph requires `forward` for acquire policies and `bypass` for bypass policies; changing this value requires replacement.
- `destinations` (Attributes List) Destinations for the rule. Every destination must have the same `type` as `rule_type`. IP addresses, ranges and subnets must be IPv4; Graph rejects IPv6 for Internet Access rules. (see [below for nested schema](#nestedatt--destinations))
- `forwarding_policy_id` (String) The forwarding policy ID that owns this rule.
- `name` (String) The rule name. Graph does not allow updating this value after creation; changing it requires replacement.
- `ports` (Set of String) Network ports matched by the rule. Each entry is a single port such as `443` or an inclusive range such as `1000-2000`, between 0 and 65535.
- `protocol` (String) Network protocol. Possible values are `tcp` and `udp`; a rule matches one protocol.
- `rule_type` (String) The destination rule type. Possible values are `fqdn`, `ip_address`, `ip_range`, and `ip_subnet`.

### Optional

- `client_fallback_action` (String) What the client does with matching traffic when it cannot reach Global Secure Access. Possible values are `default`, `bypass`, and `block`. When omitted on create, Graph sets `default`; removing it later keeps the existing value. Graph does not allow updating it after creation, so changing it requires replacement.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The forwarding policy rule ID.

<a id="nestedatt--destinations"></a>
### Nested Schema for `destinations`

Required:

- `type` (String) Destination type. Possible values are `fqdn`, `ip_address`, `ip_range`, and `ip_subnet`.

Optional:

- `begin_address` (String) Beginning IPv4 address for `ip_range`.
- `end_address` (String) Ending IPv4 address for `ip_range`. Must not be lower than `begin_address`.
- `value` (String) Destination value for `fqdn`, `ip_address`, and `ip_subnet`. Use CIDR notation for `ip_subnet`. Not used for `ip_range`.


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
# Import using composite ID format: {forwarding_policy_id}/{rule_id}
terraform import microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule.example dad2a411-e330-440d-a7c7-2c830dce5991/66f1c2ee-9f17-4681-a7bc-c324e7dff554
```
