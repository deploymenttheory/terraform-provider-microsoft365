# Live Terraform validation — 2026-09-07

Executed against a feature-enabled development tenant with Terraform 1.14.6 and a locally built provider. Authentication used environment variables and a temporary application with only `NetworkAccess.ReadWrite.All`; no credentials are included here. Only disposable, unassigned `tf-api-probe-` resources were managed. This verifies configuration persistence, not traffic evaluation.

## Manual lifecycle: HCL → plan → apply → GET

Initial configuration (as executed, including all output blocks):

```hcl
terraform {
  required_providers {
    microsoft365 = {
      source = "deploymenttheory/microsoft365"
    }
  }
}

provider "microsoft365" {
  auth_method = "client_secret"
  cloud = "public"
}

resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy" "test" {
  name           = "tf-api-probe-mcp-manual-20260907"
  description    = "Initial description"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule" "test" {
  mcp_policy_id       = microsoft365_graph_beta_identity_and_access_network_mcp_policy.test.id
  name                = "tf-api-probe-mcp-rule-20260907"
  description         = "Initial description"
  action              = "allow"
  priority            = 1000
  enabled             = true
  matching_conditions = { tool_matching = { names = { values = ["tf_probe_tool"], match_type = "exactMatch" }, methods = "call" } }
}

output "policy_id" {
  value = microsoft365_graph_beta_identity_and_access_network_mcp_policy.test.id
}

output "rule_id" {
  value = microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule.test.id
}

output "rule_status" {
  value = microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule.test.status
}

output "rule_parent_id" {
 value = microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule.test.mcp_policy_id
}
```

`terraform plan -out=01.plan` → `2 to add, 0 to change, 0 to destroy`.
`terraform apply 01.plan` → `2 added, 0 changed, 0 destroyed`.
Policy and rule item GETs returned 200; `rule_status = "enabled"`. Subsequent `terraform plan -detailed-exitcode` returned 0.

Updated resource blocks below replace the initial blocks; provider and outputs remain as above. Descriptions are omitted, tool matching is removed, and the other conditions are added:

```hcl
resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy" "test" {
  name           = "tf-api-probe-mcp-manual-updated-20260907"
  default_action = "block"
}

resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule" "test" {
  mcp_policy_id = microsoft365_graph_beta_identity_and_access_network_mcp_policy.test.id
  name          = "tf-api-probe-mcp-rule-updated-20260907"
  action        = "block"
  priority      = 65001
  enabled       = false
  matching_conditions = {
    server_urls = {
      values     = ["https://Example.invalid/mcp/", "https://example.invalid/MCP"]
      match_type = "notContains"
    }
    protocol_versions = {
      values     = ["2025-11-25", "custom-version"]
      match_type = "exactMatch"
    }
    insecure_connection = "required"
    missing_prm         = "excluded"
    resource_matching = {
      names   = { values = [], match_type = "contains" }
      methods = "read"
    }
    prompt_matching = {
      names   = { values = ["tf_probe_prompt"], match_type = "notExactMatch" }
      methods = "get"
    }
  }
}

```

Plan → `0 to add, 2 to change, 0 to destroy`; apply → `0 added, 2 changed, 0 destroyed`. Both IDs stayed unchanged. Item GETs returned 200 with null descriptions, block actions, priority 65001, disabled status, toolMatching null, the exact URL/order/case/trailing-slash values, custom protocol version, required insecureConnection, excluded missingPrm, and an empty resource names array. No-op plan returned 0.

Next, set both descriptions to `""`, set rule `enabled = true`, and remove the entire `matching_conditions` attribute. Plan → `0 to add, 2 to change, 0 to destroy`; apply → `0 added, 2 changed, 0 destroyed`. IDs stayed unchanged. GET returned empty descriptions, enabled status, and `matchingConditions = {"sources":null,"destinations":null}`. No-op plan returned 0.

## CLI and identity imports

After backing up state, remove the two resource addresses from local state, then run:

```sh
terraform import microsoft365_graph_beta_identity_and_access_network_mcp_policy.test 8ddd80b0-df2e-4475-b800-d550a561f889
terraform import microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule.test 8ddd80b0-df2e-4475-b800-d550a561f889/2dd4815c-11ab-4323-aefa-af43bdf6fcfb
terraform plan -detailed-exitcode
```

Both imports succeeded; plan returned 0. Remove the addresses from local state again, then add these executed identity import blocks:

```hcl
import {
 to = microsoft365_graph_beta_identity_and_access_network_mcp_policy.test
 identity = { id = "8ddd80b0-df2e-4475-b800-d550a561f889" }
}
import {
 to = microsoft365_graph_beta_identity_and_access_network_mcp_policy_rule.test
 identity = { mcp_policy_id = "8ddd80b0-df2e-4475-b800-d550a561f889", id = "2dd4815c-11ab-4323-aefa-af43bdf6fcfb" }
}
```

Plan → `2 to import, 0 to add, 0 to change, 0 to destroy`; apply → `2 imported, 0 added, 0 changed, 0 destroyed`. Remove the import blocks; no-op plan returned 0. Neither import flow created cloud resources.

## Parent replacement

First add the following parent and output, keeping the rule under the original policy:

```hcl
resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy" "other" {
  name           = "tf-api-probe-mcp-other-20260907"
  default_action = "allow"
}

output "other_policy_id" {
  value = microsoft365_graph_beta_identity_and_access_network_mcp_policy.other.id
}
```

Plan/apply: 1 added, 0 changed, 0 destroyed; no-op plan 0. Then change the rule's `mcp_policy_id` to `microsoft365_graph_beta_identity_and_access_network_mcp_policy.other.id`.

Plan excerpt: `mcp_policy_id = "8ddd80b0-df2e-4475-b800-d550a561f889" -> "a4531a13-5db4-455f-b9d9-6834c55508af" # forces replacement`.
Plan/apply: 1 added, 0 changed, 1 destroyed. Rule ID changed from R1 to R2; R1 GET 404, current rule GET 200; no-op plan 0.

Next add this new parent/output and change the rule's `mcp_policy_id` to `microsoft365_graph_beta_identity_and_access_network_mcp_policy.new.id` in the same plan:

```hcl
resource "microsoft365_graph_beta_identity_and_access_network_mcp_policy" "new" {
  name           = "tf-api-probe-mcp-new-20260907"
  default_action = "allow"
}

output "new_policy_id" {
  value = microsoft365_graph_beta_identity_and_access_network_mcp_policy.new.id
}
```

Plan excerpt: `mcp_policy_id = "a4531a13-5db4-455f-b9d9-6834c55508af" -> (known after apply) # forces replacement`.
Plan/apply: 2 added, 0 changed, 1 destroyed. Rule ID changed from R2 to R3; R2 GET 404, current rule GET 200; no-op plan 0.

| Label | Created ID | Parent |
| --- | --- | --- |
| P1 | 8ddd80b0-df2e-4475-b800-d550a561f889 | — |
| P2 | a4531a13-5db4-455f-b9d9-6834c55508af | — |
| P3 | 0c05b71f-f610-4171-b1a1-40b83c366c17 | — |
| R1 | 2dd4815c-11ab-4323-aefa-af43bdf6fcfb | P1 |
| R2 | b523e98b-3641-4d23-81ba-89fd02832acb | P2 |
| R3 | 7c6f1f71-21d1-4f5b-aabc-7349e6729588 | P3 |

## Go acceptance tests (separate execution)

With application authentication environment variables, `TF_ACC=1`, and the manual dev-overrides configuration removed:

```sh
GOFLAGS=-p=2 go test -vet=off -p 1 -v -timeout 30m \
  ./internal/services/resources/identity_and_access/graph_beta/network_mcp_policy \
  ./internal/services/resources/identity_and_access/graph_beta/network_mcp_policy_rule \
  -run '^TestAccResourceNetworkMCPPolicy' -count=1
```

- `TestAccResourceNetworkMCPPolicy_01_Lifecycle`: PASS (22.31s).
- `TestAccResourceNetworkMCPPolicyRule_01_Lifecycle`: PASS (37.22s).

Both include create, import, updates, null/empty descriptions and CheckDestroy. Rule coverage also includes disable/re-enable and condition removal. Acceptance policy ID: `ee68284d-d5dc-4c89-bd46-37d50ef69338`; rule test parent: `fc270598-ba25-4717-a319-0f01f5cb61df`; rule: `49c8161f-e8a3-49d3-a734-4dd864ddf1d2`. These IDs were recorded by the tests and independently checked during cleanup.

## Destroy and credential cleanup

`terraform plan -destroy -out=08-destroy.plan` → `0 to add, 0 to change, 4 to destroy`; `terraform apply 08-destroy.plan` → `0 added, 0 changed, 4 destroyed` (three parents and the final rule).

Independent item GETs returned **404 for all nine recorded policy/rule URLs**, including replaced rules and both acceptance-test parents. The collection returned to the original empty baseline. No existing cloud resources, assignments or links were modified, and no external MCP server was contacted.

The temporary secret was removed with `removePassword` (204), and its key was confirmed absent. The temporary service principal and application were deleted; subsequent GETs returned **404 for both**. The local secret file was deleted. Directory deletion required subsequent reads to observe absence; this is separate from MCP resource consistency.
