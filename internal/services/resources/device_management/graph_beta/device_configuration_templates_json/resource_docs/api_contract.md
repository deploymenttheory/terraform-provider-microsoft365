# Device configuration template JSON API contract

This resource manages the flat `/beta/deviceManagement/deviceConfigurations` profile object. It is not a Settings Catalog `/configurationPolicies` resource. The public interface follows the metadata/settings/assignments separation: `display_name`, `description`, `role_scope_tag_ids`, complete writable `settings` JSON including `@odata.type`, and standard `assignments`.

## Evidence and boundaries

Baseline fixtures in `../tests/responses` and `../tests/terraform` were derived from disposable-profile curl captures made on 2026-09-28. Additional unit fixtures model response variants and failures; acceptance fixtures exercise those request shapes against Graph. Tenant IDs, resource IDs, timestamps, and encrypted secret references were replaced with synthetic values. The certificate fixture is a public test certificate. No existing managed profile was modified during the study.

A follow-up curl check on 2026-09-29 verified that Base64-encoded XML writes return exactly the original multiline Unicode text through the OMA plaintext endpoint. The disposable profile was deleted and a subsequent GET returned 404. The XML unit and acceptance fixtures now configure cleartext values.

The Graph beta metadata inventory contained 140 device-configuration entity types, including base and abstract types. The relationship map includes profile-to-profile certificate navigation properties, including inherited relationships. Live tests exclude the inherited rootCertificate routes on Android owner/work imported-PFX types because Graph rejects those routes. This inventory establishes API shape, not successful live testing of every type. The resource does not maintain an allowlist of root profile types. Graph validates their writable properties.

## Verified request behavior

| Operation | Request | Observed result |
| --- | --- | --- |
| Create | `POST /deviceConfigurations` | `201`, profile ID returned |
| Read | `GET /deviceConfigurations/{id}` | `200`, complete flat profile with defaults and read-only metadata |
| Update | `PATCH /deviceConfigurations/{id}` | `204`; omitted properties retain their previous values |
| Replace profile body | `PUT /deviceConfigurations/{id}` | `400`; unsupported |
| Assign or unassign | `POST /deviceConfigurations/{id}/assign` | `200` with empty body; an empty assignments array clears assignments |
| Read assignments | `GET /deviceConfigurations/{id}/assignments` | Collection of targets |
| Delete | `DELETE /deviceConfigurations/{id}` | `200` with empty body, followed by `404` on read |
| Read certificate reference | `GET /deviceConfigurations/{id}/microsoft.graph.{type}/{relationship}` | Singleton object or collection; unbound singleton returns `404` |
| Replace singleton reference | `PUT .../{relationship}/$ref` with `@odata.id` | `204` |
| Remove singleton reference | `DELETE .../{relationship}/$ref` | `204` |
| Add collection reference | `POST .../{relationship}/$ref` with `@odata.id` | `204` |
| Remove collection reference | `DELETE .../{relationship}/{relatedId}/$ref` | `204` |

A singleton `@odata.bind: null` and collection `@odata.bind: []` were rejected by Graph. PATCH of a subset does not remove previous references. `$expand` and uncast SCEP navigation reads were also rejected. The implementation uses the verified typed routes and explicit reference operations; omission of a binding removes that relationship during update.

## State behavior

JSON processing occurs in construction and state mapping. There are no resource or attribute plan modifiers in this resource. State uses the existing `normalize.JSONAlphabetically` helper. If the entire returned writable JSON is semantically equivalent to the configured JSON, the exact configured string is retained. This preserves formatting without rewriting a known planned value. Different remote values remain visible as drift. Import reads the full remote payload without prior configuration.

The state mapper excludes root metadata owned by separate attributes, Graph response annotations, and verified read-only fields. It does not project the response onto configured keys. Configure complete writable settings, including returned default/null/empty values. To reset an ordinary setting, provide an API-supported explicit value; omitting it from a PATCH does not reset it.

Navigation collections are unordered. State preserves configured URL spelling and order only when the actual related IDs match. Import emits deterministic sorted URLs. Tests compare every imported JSON field while allowing ordering differences only in reference collections. Ordinary settings arrays retain their order.

## Nontrivial JSON cases

- **Windows OMA:** integer, boolean, floating-point, date/time, string, Base64, and XML values have fixtures. Encrypted string/Base64/XML values require `getOmaSettingPlainTextValue`. The endpoint returns encoded Base64 data unchanged, but returns decoded XML text for `omaSettingStringXml`; configure XML in cleartext. The constructor encodes XML for the API, and state retains the plaintext response. Unencrypted XML responses are decoded from the wire encoding. Binary OMA values retain their Base64 representation. Encryption flags and secret-reference IDs are response-only.
- **Wi-Fi:** iOS, macOS, Windows, Android Enterprise, Android work profile, and AOSP return null for `preSharedKey`. State retains the configured key when Graph masks it. Where Graph returns `preSharedKeyIsSet: false`, the clear is authoritative and is exposed as drift. Import cannot recover this key. Acceptance tests check all other imported settings against captured responses.
- **Windows Update:** six derived pause/rollback timestamps plus `qualityUpdatesWillBeRolledBack` and `featureUpdatesWillBeRolledBack` cannot be echoed in a settings PATCH. Removing all eight allowed the full writable payload to be updated successfully. These are omitted from settings state.
- **Nested OData:** concrete home-screen page/folder-page discriminators can be absent in GET responses while abstract icon discriminators remain. The examples follow the actual readable/writable shape rather than inventing discriminator insertion rules.
- **Android Enterprise Wi-Fi:** embedded root bindings can be accepted without being stored. The plural root relationship returned an empty collection even after a successful PUT. The singular `rootCertificateForServerValidation/$ref` PUT persisted the reference and GET returned it. The example uses this singular binding, and creation verifies actual relationships and writes missing references explicitly.
- **Certificates:** SCEP, enterprise Wi-Fi, VPN, and email can hold references absent from the ordinary profile GET. The mapper reads these references from their typed navigation endpoints so import and drift do not depend on previously configured IDs.

Secrets remain in Terraform state. Root profile type changes require an explicit Terraform replacement; update rejects an in-place type change before making API writes.

## Assignment and encoding regression coverage

The assignment lifecycle tests cover one group target expanding to all four target types, and the reverse transition back to one group. Acceptance tests deploy the inclusion and exclusion security groups in a separate first step, wait for directory consistency, and then create the profile. The maximal set contains all devices, all licensed users, the inclusion group, and the exclusion group. The custom profile carries an empty payload, so these broad assignment targets do not configure device settings. Both directions verify target membership, the unchanged profile ID, import, an empty plan, and destruction. Existing tests retain coverage for filters, absent assignments, and explicit empty assignments.

XML lifecycle tests configure cleartext, update it to multiline Unicode text, import it, and verify an empty plan. Mixed OMA tests combine XML, ordinary strings that resemble Base64, binary bytes, integers, booleans, and floating-point values, then update and reorder the array. A separate unit fixture verifies decoding an unencrypted XML response. Encoding is selected by the OMA type, never by guessing whether an arbitrary string resembles Base64.

## Existing endpoint implementations consulted

The endpoint already had implementations in macOS device configuration templates, macOS software update configuration, Windows update rings, Windows custom configuration, Windows trusted root certificates, and the v1 device configuration assignment resource. Windows update ring actions/data sources, policy set validation, and the configuration-policy import PowerShell tooling also use this endpoint. These uses demonstrate that `/deviceConfigurations` is cross-platform.

The new resource reuses the provider's SDK adapter, timeout/retry/error handling, conversion helpers, standard assignment schema/model, test factories, Graph existence checks, destruction checks, and HCL/JSON fixture readers. Settings Catalog JSON provides the metadata/settings separation and common state normalization pattern; its collection-shaped payload and update semantics are not copied onto this flat endpoint.

### Prior endpoint references

Source inventory before this resource was added (excluding tests and mocks):

- `internal/services/datasources/device_management/graph_beta/windows_update_ring/datasource.go`
- `internal/services/datasources/device_management/graph_beta/windows_update_ring/read.go`
- `internal/services/resources/device_management/graph_beta/macos_device_configuration_templates/construct_resource.go`
- `internal/services/resources/device_management/graph_beta/macos_device_configuration_templates/crud.go`
- `internal/services/resources/device_management/graph_beta/macos_device_configuration_templates/resource.go`
- `internal/services/resources/device_management/graph_beta/macos_device_configuration_templates/resource_docs/create_req.md`
- `internal/services/resources/device_management/graph_beta/macos_device_configuration_templates/state.go`
- `internal/services/resources/device_management/graph_beta/macos_software_update_configuration/crud.go`
- `internal/services/resources/device_management/graph_beta/macos_software_update_configuration/resource.go`
- `internal/services/resources/device_management/graph_beta/policy_set/reqs/reqs.md`
- `internal/services/resources/device_management/graph_beta/policy_set/validate.go`
- `internal/services/resources/device_management/graph_beta/windows_custom_configuration/crud.go`
- `internal/services/resources/device_management/graph_beta/windows_custom_configuration/resource.go`
- `internal/services/resources/device_management/graph_beta/windows_trusted_root_certificate/crud.go`
- `internal/services/resources/device_management/graph_beta/windows_trusted_root_certificate/resource.go`
- `internal/services/resources/device_management/graph_beta/windows_update_ring/crud.go`
- `internal/services/resources/device_management/graph_beta/windows_update_ring/resource.go`
- `internal/services/resources/device_management/graph_beta/windows_update_ring/resource_docs/update_req.md`
- `internal/services/resources/device_management/graph_beta/windows_update_ring_action/crud.go`
- `internal/services/resources/device_management/graph_beta/windows_update_ring_action/resource.go`
- `internal/services/resources/device_management/graph_v1.0/device_configuration_assignment/crud.go`
- `internal/services/resources/device_management/graph_v1.0/device_configuration_assignment/resource.go`
- `scripts/powershell/device_management/Test-ConfigurationPolicyForTFImport.ps1`

## UI template coverage

See [the per-platform checklist](template_coverage.md) for the 72 requested UI entries, lifecycle evidence, and separate API families. The current Delivery Optimization wizard uses a versioned configuration-policy template; its older device-configuration type remains available separately.
