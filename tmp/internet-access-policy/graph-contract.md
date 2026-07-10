# Internet Access Policy Graph Contract Notes

Date: 2026-07-10

## Authentication

Live probing was run with a temporary service principal created by Azure CLI:

- Tenant ID: `2cbe968c-9683-4d8a-af06-dab1bb350a04`
- App display name: `codex-iap-terraform-20260710073944`
- App ID: `ea74fefa-eecb-4b80-a543-a750177d1a67`

The service principal was granted Microsoft Graph application permissions:

- `NetworkAccess.Read.All`
- `NetworkAccess.ReadWrite.All`

Authentication command shape:

```bash
az login --service-principal --tenant "$AZURE_TENANT_ID" -u "$AZURE_CLIENT_ID" -p "$AZURE_CLIENT_SECRET"
TOKEN="$(az account get-access-token --resource https://graph.microsoft.com --query accessToken -o tsv)"
```

The secret is intentionally not recorded.

## Observed Forwarding Profile GET

Endpoint:

```http
GET https://graph.microsoft.com/beta/networkaccess/forwardingProfiles/72661c0d-027e-4dff-8c76-af103f200903?$expand=policies($expand=policy)
```

Key observed contract:

- Forwarding profile ID: `72661c0d-027e-4dff-8c76-af103f200903`
- Name: `Internet traffic forwarding profile`
- `trafficForwardingType`: `internet`
- `clientFallbackAction`: `bypass`
- `policies` contains `#microsoft.graph.networkaccess.forwardingPolicyLink` items.
- Link ID and policy ID are distinct.
- `Custom Acquire` policy ID observed as `dad2a411-e330-440d-a7c7-2c830dce5991`.
- `Custom bypass` policy ID observed as `20b11851-d497-4a6e-8c98-a9bd9461eeb2`.
- `Default Acquire` policy link ID observed as `09837256-2cba-4dde-a121-4d6a129f13db`; linked policy ID observed as `f0474b3e-307a-4230-bc1c-cd8ac2f1a2cf`.
- Microsoft 365 traffic profile ID observed as `233d4bc3-a943-44f1-8f7a-62852fdb79d5`.
- Microsoft 365 `Exchange Online` policy link ID observed as `e4dad1ae-c5c7-4ad6-aa76-3e6fb6e734f2`; linked policy ID observed as `1d6ed541-47fe-4899-9868-4edde0465e2e`.

## Observed Policy Link State PATCH

Endpoint:

```http
PATCH https://graph.microsoft.com/beta/networkaccess/forwardingProfiles/72661c0d-027e-4dff-8c76-af103f200903/policies/09837256-2cba-4dde-a121-4d6a129f13db
```

Body:

```json
{"state":"enabled"}
```

Provider implementation sends:

```json
{
  "@odata.type": "#microsoft.graph.networkaccess.forwardingPolicyLink",
  "state": "enabled"
}
```

Live service principal observation for a Microsoft 365 traffic profile policy link:

```http
PATCH https://graph.microsoft.com/beta/networkaccess/forwardingProfiles/233d4bc3-a943-44f1-8f7a-62852fdb79d5/policies/e4dad1ae-c5c7-4ad6-aa76-3e6fb6e734f2
```

Provider-shaped body:

```json
{"@odata.type":"#microsoft.graph.networkaccess.forwardingPolicyLink","state":"disabled"}
```

Observed status: `204`.

Verification:

- Initial GET returned `state = "enabled"` for policy link `e4dad1ae-c5c7-4ad6-aa76-3e6fb6e734f2`, policy `Exchange Online`, `trafficForwardingType = "m365"`.
- PATCH to `disabled` returned `204`; follow-up GET returned `state = "disabled"`.
- PATCH restore to `enabled` returned `204`; follow-up GET returned `state = "enabled"`.

The provider resource is generic for forwarding profile policy links and uses the same PATCH endpoint for Internet Access and Microsoft 365 traffic profile policy links.

## Observed Internet Access FQDN Rule Create

Endpoint:

```http
POST https://graph.microsoft.com/beta/networkaccess/forwardingPolicies/dad2a411-e330-440d-a7c7-2c830dce5991/policyRules
```

Request body:

```json
{
  "name": "Custom Acquire policy internet rule",
  "action": "forward",
  "destinations": [
    {
      "@odata.type": "#microsoft.graph.networkaccess.fqdn",
      "value": "example.com"
    }
  ],
  "ruleType": "fqdn",
  "ports": ["80", "443"],
  "protocol": "udp",
  "@odata.type": "#microsoft.graph.networkaccess.internetAccessForwardingRule"
}
```

Response body:

```json
{
  "@odata.context": "https://graph.microsoft.com/beta/$metadata#networkAccess/forwardingPolicies('dad2a411-e330-440d-a7c7-2c830dce5991')/policyRules/$entity",
  "@odata.type": "#microsoft.graph.networkaccess.internetAccessForwardingRule",
  "id": "66f1c2ee-9f17-4681-a7bc-c324e7dff554",
  "name": "Custom Acquire policy internet rule",
  "ruleType": "fqdn",
  "action": "forward",
  "clientFallbackAction": "block",
  "ports": ["80", "443"],
  "protocol": "udp",
  "destinations": [
    {
      "@odata.type": "#microsoft.graph.networkaccess.fqdn",
      "value": "example.com"
    }
  ]
}
```

## Observed Internet Access Rule GET

Endpoint:

```http
GET https://graph.microsoft.com/beta/networkaccess/forwardingPolicies/dad2a411-e330-440d-a7c7-2c830dce5991/policyRules/66f1c2ee-9f17-4681-a7bc-c324e7dff554
```

Returned the same `internetAccessForwardingRule` shape as create response.

## Observed Internet Access Rule PATCH

Endpoint:

```http
PATCH https://graph.microsoft.com/beta/networkaccess/forwardingPolicies/dad2a411-e330-440d-a7c7-2c830dce5991/policyRules/66f1c2ee-9f17-4681-a7bc-c324e7dff554
```

Request body:

```json
{
  "id": "66f1c2ee-9f17-4681-a7bc-c324e7dff554",
  "ruleType": "fqdn",
  "destinations": [
    {
      "@odata.type": "#microsoft.graph.networkaccess.fqdn",
      "value": "example.com"
    }
  ],
  "ports": ["80", "443"],
  "protocol": "tcp",
  "@odata.type": "#microsoft.graph.networkaccess.internetAccessForwardingRule"
}
```

Observed status: `204`.

Provider implementation sends `id`, `ruleType`, `ports`, `protocol`, `destinations`, and `@odata.type` on update, and refreshes state with GET after 204. Live probing showed `name` is not patchable.

## Observed Internet Access Rule DELETE

Endpoint:

```http
DELETE https://graph.microsoft.com/beta/networkaccess/forwardingPolicies/dad2a411-e330-440d-a7c7-2c830dce5991/policyRules/66f1c2ee-9f17-4681-a7bc-c324e7dff554
```

Observed status: `204`.

## Destination Shapes Implemented

The user noted that Internet Access destinations include FQDN, IP, IP range, and CIDR. The provider currently serializes:

| Terraform type | Graph `ruleType` | Graph destination `@odata.type` | Fields |
| --- | --- | --- | --- |
| `fqdn` | `fqdn` | `#microsoft.graph.networkaccess.fqdn` | `value` |
| `ip_address` | `ipAddress` | `#microsoft.graph.networkaccess.ipAddress` | `value` |
| `ip_range` | `ipRange` | `#microsoft.graph.networkaccess.ipRange` | `beginAddress`, `endAddress` |
| `ip_subnet` | `ipSubnet` | `#microsoft.graph.networkaccess.ipSubnet` | `value` with CIDR notation |

The non-FQDN shapes are implemented from the observed Graph type naming pattern and should be live-probed before merge in a tenant with safe test data.

## Live Terraform Probe Results

Terraform was run against a local provider build through `TF_CLI_CONFIG_FILE` development overrides.

Successful create/apply coverage:

- `fqdn` + `forward` + `tcp` on `Custom Acquire`.
- `fqdn` + `bypass` + `udp` on `Custom bypass`.
- `ipAddress` + `forward` + `tcp` on `Custom Acquire`.
- `ipRange` + `bypass` + `udp` on `Custom bypass`.
- `ipSubnet` + `forward` + `tcp` on `Custom Acquire`.
- Forwarding profile policy link adopt/apply for `Custom Acquire`; destroy removed Terraform state only.

Successful update coverage:

- FQDN destination value, ports, and protocol.
- IP address value, ports, and protocol.
- IP range begin/end address, ports, and protocol.
- IP subnet CIDR value, ports, and protocol.

Successful destroy coverage:

- All five temporary rule GETs returned `404` after Terraform destroy.
- The Microsoft-managed forwarding policy link GET returned `200` after Terraform destroy.

Observed Graph constraints:

- `action = "bypass"` on the `Custom Acquire` policy returned `400` with `Only Forward action is allowed for acquire policies`.
- Patching `name` returned `400` with `Updating the property Name for entity InternetAccessForwardingRule is not allowed`.
- Parallel rule writes can return `412 PreconditionFailed`; the resource now retries precondition failures like the existing web filtering rule resource.

Probe matrix still not covered:

- Missing rule `@odata.type`.
- Invalid `ruleType`.
- Unsupported destination type.
- Empty `ports`.
- Invalid `protocol`.
- Mismatched `ruleType` and destination `@odata.type`.
