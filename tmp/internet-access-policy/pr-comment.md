# Internet Access Policy Terraform Support

## What Changed

- Added data source `microsoft365_graph_beta_identity_and_access_network_forwarding_profile`.
- Added resource `microsoft365_graph_beta_identity_and_access_network_forwarding_profile_policy_link`.
- Added resource `microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule`.
- Registered all three in the provider.
- Added docs/templates/examples and unit tests for serialization and state conversion.

## Design Notes

- Forwarding profile reads use the generated Microsoft Graph beta SDK typed builders for `/networkAccess/forwardingProfiles` with `$expand=policies($expand=policy)`.
- Forwarding profile policy link reads use the SDK typed builder. State updates use Kiota `RequestInformation` with a small custom `Parsable` PATCH body so the provider can send the observed state-only payload.
- Internet Access forwarding policy rules use Kiota `RequestInformation` and custom `Parsable` request/response types. The SDK has policyRules builders, but the observed `internetAccessForwardingRule` and destination polymorphism are safer to model explicitly for now.
- Policy link destroy is a remote no-op and only removes Terraform state, because Microsoft manages these links.
- Internet Access rule destroy calls DELETE and expects 204.
- `name` and `action` are create-time values for Internet Access rules. Live Graph probing rejected PATCHing `name`, and Graph constrains `action` by policy type (`forward` for acquire policies, `bypass` for bypass policies).
- Internet Access rule create/update/delete retry Graph `412 PreconditionFailed`, matching the existing web filtering policy rule pattern for parallel writes within one policy.

## Example HCL

```hcl
data "microsoft365_graph_beta_identity_and_access_network_forwarding_profile" "internet" {
  traffic_forwarding_type = "internet"
}

locals {
  internet_profile = one(data.microsoft365_graph_beta_identity_and_access_network_forwarding_profile.internet.items)
  custom_acquire_policy = one([
    for link in local.internet_profile.policies : link
    if link.policy_name == "Custom Acquire"
  ])
  custom_bypass_policy = one([
    for link in local.internet_profile.policies : link
    if link.policy_name == "Custom bypass"
  ])
}

resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "fqdn" {
  forwarding_policy_id = local.custom_acquire_policy.policy_id

  name      = "Example Internet Access FQDN rule"
  action    = "forward"
  rule_type = "fqdn"
  ports     = ["80", "443"]
  protocol  = "tcp"

  destinations = [
    {
      type  = "fqdn"
      value = "example.com"
    }
  ]
}

resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "cidr_bypass" {
  forwarding_policy_id = local.custom_bypass_policy.policy_id

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
```

## Graph Evidence Summary

Observed from user-provided Graph/DevTools traffic and live service principal probing:

- `GET /beta/networkaccess/forwardingProfiles/{id}?$expand=policies($expand=policy)` returns forwarding profile fields and policy links; `policy_link_id` and `policy_id` are distinct.
- `PATCH /beta/networkaccess/forwardingProfiles/{forwardingProfileId}/policies/{policyLinkId}` accepts `{"state":"enabled"}`.
- `PATCH /beta/networkaccess/forwardingProfiles/{forwardingProfileId}/policies/{policyLinkId}` also works for Microsoft 365 traffic profile policy links; live service principal probing changed Exchange Online `enabled -> disabled -> enabled` with 204 responses and GET verification after each PATCH.
- `POST /beta/networkaccess/forwardingPolicies/{forwardingPolicyId}/policyRules` accepts `#microsoft.graph.networkaccess.internetAccessForwardingRule` with FQDN destinations, ports, protocol, and action.
- Rule item `GET` returns `clientFallbackAction`.
- Rule item `PATCH` returned 204.
- Rule item `DELETE` returned 204.
- Live create/update/delete confirmed destination shapes for `fqdn`, `ipAddress`, `ipRange`, and `ipSubnet`.
- Live probing confirmed `action = "forward"` belongs on acquire policies and `action = "bypass"` belongs on bypass policies.
- Live probing confirmed PATCH must omit `name`; changing `name` now requires replacement.
- Live destroy confirmed temporary rules return 404 after deletion, while the Microsoft-managed forwarding profile policy link still returns 200.

## Verification

```text
go test ./internal/services/datasources/identity_and_access/graph_beta/network_forwarding_profile ./internal/services/resources/identity_and_access/graph_beta/network_forwarding_profile_policy_link ./internal/services/resources/identity_and_access/graph_beta/network_internet_access_forwarding_policy_rule
ok   github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/network_forwarding_profile
ok   github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_forwarding_profile_policy_link
ok   github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_internet_access_forwarding_policy_rule
```

```text
make userdocs
Generated docs, then failed validation on pre-existing oversized docs:
- docs/resources/graph_beta_device_and_app_management_targeted_managed_app_configuration.md
- docs/resources/graph_beta_device_management_settings_catalog_configuration_policy.md
- docs/resources/graph_beta_device_management_settings_catalog_inventory_policy.md
```

## Terraform Plan/Apply/Update/Destroy Logs

Live Terraform was run with a temporary Azure CLI-created service principal granted `NetworkAccess.Read.All` and `NetworkAccess.ReadWrite.All`.

Key log files:

- `tmp/internet-access-policy/live/terraform-phase1-plan.log`
- `tmp/internet-access-policy/live/terraform-phase1b-apply.log`
- `tmp/internet-access-policy/live/terraform-phase2b-plan.log`
- `tmp/internet-access-policy/live/terraform-phase2b-apply.log`
- `tmp/internet-access-policy/live/terraform-destroy-plan.log`
- `tmp/internet-access-policy/live/terraform-destroy-apply.log`
- `tmp/internet-access-policy/live/rule-get-after-update.summary.json`
- `tmp/internet-access-policy/live/rule-get-after-destroy-status.tsv`
- `tmp/internet-access-policy/live/policy-link-after-destroy.summary.json`

Summary:

```text
phase1 apply: partial live probe exposed Graph constraints:
- action=bypass on Custom Acquire returned 400 Only Forward action is allowed for acquire policies
- parallel writes returned 412 PreconditionFailed

phase1b apply: 3 added, 0 changed, 0 destroyed
phase2b apply: 0 added, 5 changed, 0 destroyed
destroy apply: 0 added, 0 changed, 6 destroyed

post-destroy rule GET statuses:
fqdn_bypass_udp  404
fqdn_forward_tcp 404
ip_address       404
ip_range         404
ip_subnet        404

post-destroy policy link GET:
status 200, id f576d498-0067-4cc8-960b-b6e3ebf571ea, state enabled, policy Custom Acquire
```
