# Hardcoded Wait Audit — Eventual Consistency

Audit of all `time.Sleep` calls in resource `crud.go` files.
Generated: 2026-03-26

---

## Summary

| Category | Count |
|---|---|
| Resources with sleeps in Create | 20 |
| Resources with sleeps in Update | 17 |
| Resources with sleeps in Delete | 3 |
| Resources with rate-limit sleeps (loops) | 2 |
| **Total resources affected** | **27** |

---

## Resources with Hardcoded Waits in Create / Update

| Resource | Create | Update | Notes |
|---|---|---|---|
| `agents/graph_beta/agent_identity` | 20s | 15s | |
| `agents/graph_beta/agent_identity_blueprint` | 25s | 15s | |
| `agents/graph_beta/agent_identity_blueprint_certificate_credential` | 15s | 5s + 15s | Two sleeps in Update |
| `agents/graph_beta/agent_identity_blueprint_federated_identity_credential` | 10s | 15s | |
| `agents/graph_beta/agent_identity_blueprint_identifier_uri` | 10s | — | |
| `agents/graph_beta/agent_identity_blueprint_service_principal` | 10s | 10s | |
| `agents/graph_beta/agent_instance` | 5s | 5s | |
| `agents/graph_beta/agent_user` | 5s | 5s | |
| `applications/graph_beta/application` | 25s | 5s (loop) + 15s | |
| `applications/graph_beta/application_certificate_credential` | 30s + 10s (retry) | 5s + 15s | Two sleeps in Create; two in Update |
| `applications/graph_beta/application_federated_identity_credential` | 10s | 15s | |
| `applications/graph_beta/application_identifier_uri` | 10s | — | |
| `applications/graph_beta/application_owner` | 5s | — | |
| `applications/graph_beta/service_principal` | — | 20s | |
| `applications/graph_beta/service_principal_owner` | 5s | — | |
| `groups/graph_beta/license_assignment` | 15s | 15s | |
| `identity_and_access/graph_beta/administrative_unit` | — | 10s | |
| `identity_and_access/graph_beta/administrative_unit_membership` | 20s (conditional) | 20s (conditional) | Only when members change |
| `identity_and_access/graph_beta/conditional_access_policy` | — | 10s | |
| `identity_and_access/graph_beta/cross_tenant_access_default_settings` | 20s | 20s | |
| `identity_and_access/graph_beta/cross_tenant_access_partner_settings` | 20s | 30s | |
| `identity_and_access/graph_beta/cross_tenant_access_policy` | 10s | 10s | |
| `users/graph_beta/user` | — | 15s | |
| `users/graph_beta/user_manager` | — | 2s | |
| `windows_updates/autopatch_deployment_state` | 10s | 10s | |
| `windows_updates/autopatch_updatable_asset_group` | 20s (conditional) | 20s | |

---

## Resources with Hardcoded Waits in Delete

| Resource | Delete | Notes |
|---|---|---|
| `identity_and_access/graph_beta/cross_tenant_access_default_settings` | 20s | |
| `device_management/graph_beta/assignment_filter` | 10s | |
| `windows_365/graph_beta/cloud_pc_provisioning_policy` | 2s | |

---

## Rate-Limit / Throttle Sleeps (not eventual consistency)

These sleeps are inside processing loops to avoid API throttling. They are distinct from eventual consistency waits but are still hardcoded.

| Resource | Function | Duration | Notes |
|---|---|---|---|
| `device_management/graph_beta/device_compliance_notification_template` | Update | 2s | In notification message loop |
| `device_management/graph_beta/device_enrollment_notification` | Create + Update | 1s–2s | Multiple sleeps in item processing loops |

---

## Excluded from Report

The following contain `time.Sleep` via a `retryDelay` variable (not a hardcoded literal) and are therefore dynamic retry patterns rather than fixed waits:

- `identity_and_access/graph_beta/administrative_unit_membership` — Read (retry loop, `retryDelay = 5s`)
- `groups/graph_beta/group` — Read (retry loop, `retryDelay`)
