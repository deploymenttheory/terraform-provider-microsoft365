# Internet Access Policy Graph Contract Notes

Date: 2026-07-10

## Authentication

Live probing with service principal authentication was not run in this workspace because the required environment variables were not set:

- `AZURE_TENANT_ID`: unset
- `AZURE_CLIENT_ID`: unset
- `AZURE_CLIENT_SECRET`: unset

`az` is installed (`azure-cli 2.87.0`). To run live probing:

```bash
az login --service-principal --tenant "$AZURE_TENANT_ID" -u "$AZURE_CLIENT_ID" -p "$AZURE_CLIENT_SECRET"
TOKEN="$(az account get-access-token --resource https://graph.microsoft.com --query accessToken -o tsv)"
```

The service principal needs Microsoft Graph application permissions:

- `NetworkAccess.Read.All`
- `NetworkAccess.ReadWrite.All`

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
- `Default Acquire` policy link ID observed as `09837256-2cba-4dde-a121-4d6a129f13db`; linked policy ID observed as `f0474b3e-307a-4230-bc1c-cd8ac2f1a2cf`.

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

Provider implementation sends the full writable body on update (`id`, `name`, `action`, `ruleType`, `ports`, `protocol`, `destinations`, and `@odata.type`) and refreshes state with GET after 204.

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

## Probe Matrix Still Needed

Run these against a disposable rule in the Custom Acquire forwarding policy:

- POST/GET/PATCH/DELETE FQDN rule with `protocol = tcp`.
- POST/GET/PATCH/DELETE FQDN rule with `protocol = udp`.
- POST/GET/PATCH/DELETE IP address rule.
- POST/GET/PATCH/DELETE IP range rule.
- POST/GET/PATCH/DELETE CIDR/IP subnet rule.
- POST with `action = forward`.
- POST with `action = bypass`.
- PATCH `name`, `action`, `protocol`, `ports`, and `destinations` independently.
- Invalid probes:
  - missing rule `@odata.type`
  - invalid `ruleType`
  - unsupported destination type
  - empty `ports`
  - invalid `protocol`
  - mismatched `ruleType` and destination `@odata.type`

