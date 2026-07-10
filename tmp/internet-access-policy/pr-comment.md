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
```

## Graph Evidence Summary

Observed from user-provided Graph/DevTools traffic:

- `GET /beta/networkaccess/forwardingProfiles/{id}?$expand=policies($expand=policy)` returns forwarding profile fields and policy links; `policy_link_id` and `policy_id` are distinct.
- `PATCH /beta/networkaccess/forwardingProfiles/{forwardingProfileId}/policies/{policyLinkId}` accepts `{"state":"enabled"}`.
- `POST /beta/networkaccess/forwardingPolicies/{forwardingPolicyId}/policyRules` accepts `#microsoft.graph.networkaccess.internetAccessForwardingRule` with FQDN destinations, ports, protocol, and action.
- Rule item `GET` returns `clientFallbackAction`.
- Rule item `PATCH` returned 204.
- Rule item `DELETE` returned 204.

Additional live probing with service principal credentials was not run in this workspace because `AZURE_TENANT_ID`, `AZURE_CLIENT_ID`, and `AZURE_CLIENT_SECRET` were unset. See `tmp/internet-access-policy/graph-contract.md` for the probe matrix.

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

Not run. Live Terraform apply requires a tenant with Global Secure Access enabled and service principal credentials with `NetworkAccess.ReadWrite.All`.
