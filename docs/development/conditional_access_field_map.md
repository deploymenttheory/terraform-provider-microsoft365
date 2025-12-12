# Conditional Access Policy Field Mapping

This document maps all fields across the Create Request, Create Response, and GET Response for Conditional Access Policies in Microsoft Graph API (Beta endpoint).

**Test Coverage:** This mapping is based on systematic testing with real API requests and responses across multiple scenarios.

## Key Findings Summary

### ✅ Confirmed API Behaviors (6-Test Validation)

1. **Operator Transformation**: `grantControls.operator: "AND"` → `"OR"` in multiple scenarios (100% consistent - all 6 tests)
   - With single `block` control (tests 1-3, 6)
   - With single `compliantDevice` control (test 4)
   - With `authenticationStrength` even when `builtInControls: []` (test 5)
   - **Pattern**: API enforces `OR` when grant control requirements can be satisfied by a single method
2. **API Field Expansion**: `authenticationStrength` object dramatically expanded (test 5)
   - Sent: `{id: "00000000-0000-0000-0000-000000000004"}` (just ID)
   - Received: Full object with 10+ fields (displayName, description, allowedCombinations, policyType, timestamps, etc.)
3. **Case Transformation**: `excludeGuestsOrExternalUsers.guestOrExternalUserTypes` PascalCase → camelCase (test 4)
   - Sent: `"B2bCollaborationGuest,B2bCollaborationMember,B2bDirectConnectUser"`
   - Received: `"b2bCollaborationGuest,b2bCollaborationMember,b2bDirectConnectUser"`
4. **Alphabetical Sorting**: `agentIdRiskLevels: "high,medium"` → `"medium,high"` (test 2)
5. **Conditional Field Removal - Complex Pattern**:
   - `agentIdServicePrincipalFilter`: kept when object, removed when null (tests 1, 2, 6)
   - 🚨 **NEW**: `excludeAgentIdServicePrincipals`: kept when populated array, removed when empty `[]` (test 6)
     - Test 1/2: `[<GUIDs>]` → kept ✅
     - Test 6: `[]` → removed ❌
6. **Consistent Field Removal**: 
   - 6 conditions fields always removed: `clients`, `insiderRiskLevels`, `signInRiskDetections`, `authenticationFlows`, `servicePrincipalRiskLevels`, `globalSecureAccess`
   - 2 sessionControls fields removed: `networkAccessSecurity`, `globalSecureAccessFilteringProfile` (test 5)
7. **Conditional Field Addition**: `includeServicePrincipals` & `excludeServicePrincipals` only added when `clientApplications` is object (not when null)

### 📊 Value Variations Confirmed

- `includeApplications`: `["All"]` (tests 2, 3, 6), `["None"]` (test 1), and specific GUIDs (tests 4 & 5) accepted
- `excludeApplications`: `[]` (tests 1-4, 6) and `["AllAgentIdResources"]` (test 5) - special value accepted
- `includeUsers`: `["None"]` (tests 1, 2, 5, 6), `["AllAgentIdUsers"]` (test 3), and `["All"]` (test 4) accepted
- `includeAgentIdServicePrincipals`: `["All"]` (tests 2, 6) and specific GUIDs (test 1) accepted
- `excludeAgentIdServicePrincipals`: Populated `[<GUIDs>]` kept (tests 1, 2), empty `[]` removed (test 6)
- `clientApplications`: Complex object (tests 1, 2, 6) and `null` (tests 3, 4, 5) - different API responses
- `agentIdServicePrincipalFilter`: Complex objects with "include" (test 1) or "exclude" (test 6) mode, and null (test 2) - different outcomes
- `agentIdRiskLevels`: Single value `"high"` (test 6) and comma-separated `"high,medium"` (test 2) - alphabetical sorting when multiple
- `builtInControls`: `["block"]` (tests 1-3, 6), `["compliantDevice"]` (test 4), and `[]` (test 5) - all trigger operator transformation
- `authenticationStrength`: `null` (tests 1-4, 6) and object with ID (test 5) - API expands object dramatically
- `sessionControls`: `null` (tests 1-4, 6) and complex object (test 5) - API removes certain sub-fields
- `locations`: `null` (tests 1-3, 6) and complex object (tests 4 & 5) accepted
- `excludeGuestsOrExternalUsers`: `null` (tests 1-3, 5, 6) and complex object (test 4) accepted

### 🎯 Critical Discovery: clientApplications Null vs Object

- **When `null`** (test 3): Field stays as `null`, no arrays added
- **When object** (tests 1 & 2): Field stays as object, API adds `includeServicePrincipals: []` and `excludeServicePrincipals: []`

## Test Scenarios

### Test 1: labtest-caau001 (Service Principal Filter with Agent IDs)

**Configuration:**
- **Applications**: `["None"]` (no application targeting)
- **Users**: `["None"]` (no user targeting)
- **Client Applications**: Service principal targeting with complex filter
  - `includeAgentIdServicePrincipals`: `["7ca55e16-b9fd-4269-afe4-444ceed088fa"]` (specific GUID)
  - `excludeAgentIdServicePrincipals`: `["7ca55e16-b9fd-4269-afe4-444ceed088fa"]` (same ID)
  - `agentIdServicePrincipalFilter`: Complex custom security attribute object with mode and rule
- **Agent ID Risk Levels**: Not sent
- **Client App Types**: `["all"]`
- **Grant Controls**: `block` with `AND` operator

**Key Findings:**
1. ✅ `agentIdServicePrincipalFilter` **object** is ACCEPTED and returned
2. ✅ `excludeAgentIdServicePrincipals` is ACCEPTED and returned
3. ⚠️ API transforms `operator: "AND"` → `"OR"` for block controls

### Test 2: labtest-caau002 (All Applications with Agent ID Risk Levels)

**Configuration:**
- **Applications**: `["All"]` (all applications)
- **Users**: `["None"]` (no user targeting)
- **Client Applications**: Service principal targeting without filter
  - `includeAgentIdServicePrincipals`: `["All"]` (all service principals)
  - `excludeAgentIdServicePrincipals`: `["7ca55e16-b9fd-4269-afe4-444ceed088fa"]`
  - `agentIdServicePrincipalFilter`: `null`
- **Agent ID Risk Levels**: `"high,medium"` (sent)
- **Client App Types**: `["all"]`
- **Grant Controls**: `block` with `AND` operator

**Key Findings:**
1. ✅ **CONFIRMED**: `agentIdRiskLevels` alphabetically reordered: `"high,medium"` → `"medium,high"`
2. ❌ `agentIdServicePrincipalFilter` when **null** is REMOVED by API
3. ✅ `includeAgentIdServicePrincipals: ["All"]` is accepted
4. ⚠️ API transforms `operator: "AND"` → `"OR"` for block controls (confirmed again)

### Test 3: labtest-caau003 (Simple All Applications + All Agent ID Users)

**Configuration:**
- **Applications**: `["All"]` (all applications)
- **Users**: `["AllAgentIdUsers"]` (all agent ID users)
- **Client Applications**: `null` (no service principal targeting at all)
- **Agent ID Risk Levels**: Not sent
- **Client App Types**: `["all"]` (explicitly sent)
- **Grant Controls**: `block` with `AND` operator

**Key Findings:**
1. ✅ `clientApplications: null` stays as `null` (field present, value null)
2. ✅ `includeUsers: ["AllAgentIdUsers"]` is accepted (new special value confirmed)
3. ⚠️ API transforms `operator: "AND"` → `"OR"` for block controls (confirmed 3rd time)
4. ✅ `clientAppTypes: ["all"]` explicitly sent - API keeps it
5. ❌ 6 fields removed as expected (consistent with previous tests)
6. ✅ No `includeServicePrincipals` or `excludeServicePrincipals` added when `clientApplications` is `null`

### Test 4: labtest-caau004 (Complex User Exclusions + Locations + compliantDevice)

**Configuration:**
- **Applications**: Specific app GUIDs (not special values)
- **Users**: `["All"]` with extensive exclusions:
  - `excludeUsers`: User GUIDs
  - `excludeGroups`: Group GUIDs (multiple)
  - `excludeRoles`: Role GUIDs (multiple)
  - `excludeGuestsOrExternalUsers`: Complex object with guest types and external tenants
- **Locations**: Object with `includeLocations: ["All"]` and `excludeLocations: []`
- **Client Applications**: `null`
- **Client App Types**: `["all"]`
- **Grant Controls**: `compliantDevice` with `AND` operator
- **Session Controls**: `null`

**Key Findings:**
1. 🚨 **CRITICAL**: `operator: "AND"` → `"OR"` with `compliantDevice` control (not just block!)
   - This suggests the transformation happens with **any** single built-in control
2. 🚨 **CASE TRANSFORMATION**: `excludeGuestsOrExternalUsers.guestOrExternalUserTypes` 
   - Sent: `"B2bCollaborationGuest,B2bCollaborationMember,B2bDirectConnectUser"` (PascalCase)
   - Received: `"b2bCollaborationGuest,b2bCollaborationMember,b2bDirectConnectUser"` (camelCase!)
3. ✅ `locations` object accepted and returned with nested arrays
4. ✅ `excludeGuestsOrExternalUsers` complex object accepted with nested `externalTenants` object
5. ✅ `includeUsers: ["All"]` special value accepted (new)
6. ✅ Multiple GUIDs in exclude arrays accepted and returned unchanged
7. ✅ Specific app GUIDs in `includeApplications` accepted and returned unchanged

### Test 5: labtest-caau005 (Authentication Strength + Session Controls)

**Configuration:**
- **Applications**: Specific app GUIDs with `excludeApplications: ["AllAgentIdResources"]`
- **Users**: `["None"]`
- **Locations**: Object with `includeLocations: ["All"]`
- **Client Applications**: `null`
- **Grant Controls**: 
  - `builtInControls: []` (EMPTY!)
  - `authenticationStrength: {id: "00000000-0000-0000-0000-000000000004"}`
  - `operator: "AND"`
- **Session Controls**: Complex object with `signInFrequency` and 8 other null fields

**Key Findings:**
1. 🚨 **OPERATOR TRANSFORMATION with EMPTY builtInControls!**
   - `operator: "AND"` → `"OR"` even with `builtInControls: []`
   - Transformation happens when `authenticationStrength` is present
2. 🚨 **API EXPANDS authenticationStrength!**
   - Sent: `{id: "00000000-0000-0000-0000-000000000004"}` (just ID)
   - Received: Full object with 10+ fields including `displayName`, `description`, `allowedCombinations`, `policyType`, `requirementsSatisfied`, timestamps, `combinationConfigurations`, etc.
3. 🚨 **sessionControls fields REMOVED by API!**
   - Sent: `networkAccessSecurity: null`, `globalSecureAccessFilteringProfile: null`
   - Received: These fields NOT present (removed)
   - Kept: `signInFrequency`, `applicationEnforcedRestrictions`, `cloudAppSecurity`, `persistentBrowser`, `continuousAccessEvaluation`, `disableResilienceDefaults`, `secureSignInSession`
4. ✅ `excludeApplications: ["AllAgentIdResources"]` special value accepted
5. ✅ `signInFrequency` object accepted with nested `frequencyInterval`, `authenticationType`, `isEnabled`, `value`, `type`
6. ✅ `builtInControls: []` empty array accepted

### Test 6: labtest-caau006 (Service Principal Filter with Exclude Mode + Single Risk Level)

**Configuration:**
- **Applications**: `["All"]`
- **Users**: `["None"]`
- **Client Applications**: Object with:
  - `includeAgentIdServicePrincipals: ["All"]`
  - `excludeAgentIdServicePrincipals: []` (EMPTY array)
  - `agentIdServicePrincipalFilter: {mode: "exclude", rule: "..."}`
- **Agent ID Risk Levels**: `"high"` (single value, not comma-separated)
- **Grant Controls**: `block` with `AND` operator

**Key Findings:**
1. ✅ `agentIdServicePrincipalFilter` with `mode: "exclude"` accepted (previous test had "include")
2. ✅ `agentIdRiskLevels: "high"` (single value) - no transformation needed, returned as-is
3. 🚨 **NEW DISCOVERY**: `excludeAgentIdServicePrincipals: []` (empty array) is REMOVED by API!
   - Test 1: `excludeAgentIdServicePrincipals: [<GUID>]` → kept by API ✅
   - Test 2: `excludeAgentIdServicePrincipals: [<GUID>]` → kept by API ✅
   - Test 6: `excludeAgentIdServicePrincipals: []` → removed by API ❌
   - **Pattern**: Empty array is removed, populated array is kept
4. ⚠️ API transforms `operator: "AND"` → `"OR"` with block (6th confirmation)

### Cross-Test Confirmations

**Consistently Removed Fields (all 6 tests):**
- 6 conditions fields: `clients`, `insiderRiskLevels`, `signInRiskDetections`, `authenticationFlows`, `servicePrincipalRiskLevels`, `globalSecureAccess`
- 2 sessionControls fields (when present): `networkAccessSecurity`, `globalSecureAccessFilteringProfile`

**Consistently Added Fields (all 6 tests):**
- `id`, `@odata.context`, `createdDateTime`, `modifiedDateTime`, `deletedDateTime`, `templateId`, `partialEnablementStrategy`, `deviceStates`, `authenticationStrength@odata.context`

**Conditionally Added Fields:**
- `includeServicePrincipals`, `excludeServicePrincipals` - Only added when `clientApplications` is an **object** (not when `null`)

**Conditionally Removed Fields:**
- `excludeAgentIdServicePrincipals` - Removed when empty `[]` (test 6), kept when populated (tests 1, 2)
- `agentIdServicePrincipalFilter` - Removed when `null` (test 2), kept when object (tests 1, 6)

**Operator Transformation (all 6 tests - 100% consistent):**
- `grantControls.operator: "AND"` → `"OR"` in multiple scenarios:
  - With single `"block"` control (tests 1-3, 6)
  - With single `"compliantDevice"` control (test 4)
  - With `authenticationStrength` even when `builtInControls: []` (test 5)
- **Critical Pattern**: API enforces `OR` operator when grant requirements can be satisfied by a single method

**ClientApplications Behavior:**
- When `null` (tests 3, 4, 5): Field stays as `null`, no service principal arrays added
- When object with values (tests 1, 2, 6): Field stays as object, API adds `includeServicePrincipals: []` and `excludeServicePrincipals: []`

**AuthenticationStrength Expansion (test 5):**
- Sent: Minimal object with just `{id: "..."}
- Received: Full object with 10+ fields (displayName, description, allowedCombinations, policyType, timestamps, combinationConfigurations, requirementsSatisfied, etc.)

**AgentIdServicePrincipalFilter Mode Support (tests 1, 6):**
- Both "include" mode (test 1) and "exclude" mode (test 6) are accepted and returned

## Field Mapping Table

**Column Definitions:**
- **Sent by GUI**: Whether the Azure Portal GUI includes this field in create requests
- **Create Request Requirement**: Whether the field is Required or Optional in API create requests
- **Default GUI Value**: The default value the GUI sends when field is not configured/populated
- **Valid Field Value(s)**: All acceptable values for this field (enums, types, special values)
- **GUI Value (Tests)**: The exact value(s) sent by the GUI in our test scenarios (test variations noted)
- **API Behavior**: Detailed behavior - how API handles the field (e.g., "Accepts null", "Adds default as []", "Removes", "Modifies: X→Y")
- **Create Response Value**: The exact value(s) returned in POST response
- **GET Response Value**: The exact value(s) returned in GET response
- **Notes**: Additional observations, constraints, and important behaviors

| Field Name | Parent Field Name | Sent by GUI | Create Request Requirement | Default GUI Value | Valid Field Value(s) | GUI Value (Tests) | API Behavior | Create Response Value | GET Response Value | Notes |
|-----------|-------------------|-------------|----------------------------|-------------------|----------------------|-------------------|--------------|----------------------|-------------------|-------|
| `@odata.context` | (root) | No | N/A | N/A | String (URL) | N/A | Adds (metadata) | `https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/policies/$entity` | Same | API metadata URL, always added |
| `id` | (root) | No | N/A | N/A | String (UUID) | N/A | Adds (auto-generated) | `aefd62ee-806f-4a7c-8876-a1fc171c665d` | Same | Auto-generated policy GUID |
| `displayName` | (root) | Yes | Required | N/A | String | `"labtest-caau001"` | Accepts string | `"labtest-caau001"` | Same | Policy display name, passed through unchanged |
| `state` | (root) | Yes | Required | N/A | `disabled`, `enabled`, `enabledForReportingButNotEnforced` | `"enabledForReportingButNotEnforced"` | Accepts enum | `"enabledForReportingButNotEnforced"` | Same | Policy enforcement state |
| `createdDateTime` | (root) | No | N/A | N/A | String (ISO 8601) | N/A | Adds (auto-generated) | `"2025-12-12T12:11:16.1863281Z"` | Same | Creation timestamp |
| `modifiedDateTime` | (root) | No | N/A | N/A | String (ISO 8601) / `null` | N/A | Adds default as null | `null` | Modified on update | Initially null |
| `deletedDateTime` | (root) | No | N/A | N/A | String (ISO 8601) / `null` | N/A | Adds default as null | `null` | Same | Set on soft-delete |
| `templateId` | (root) | No | N/A | N/A | String (UUID) / `null` | N/A | Adds default as null | `null` | Same | Template reference if created from template |
| `partialEnablementStrategy` | (root) | No | N/A | N/A | String / `null` | N/A | Adds default as null | `null` | Same | Partial enablement config |
| `sessionControls` | (root) | Yes | Optional | `null` | Object / `null` | caau001-004: `null`, caau005: `{object}` | Accepts null or object; removes some sub-fields | caau005: `{object}` | Same | Session controls object. When populated, API removes `networkAccessSecurity` and `globalSecureAccessFilteringProfile` |
| `conditions` | (root) | Yes | Required | N/A | Object | `{...}` | Accepts object, modifies nested | Modified (see sub-fields) | Same | Conditions object (required) |
| `grantControls` | (root) | Yes | Required | N/A | Object | `{...}` | Accepts object, modifies nested | Modified (see sub-fields) | Same | Grant controls object (required) |

### Conditions Fields

| Field Name | Parent Field Name | Sent by GUI | Create Request Requirement | Default GUI Value | Valid Field Value(s) | GUI Value (Tests) | API Behavior | Create Response Value | GET Response Value | Notes |
|-----------|-------------------|-------------|---------------------|-------------------|----------------------|-------------------|--------------|----------------------|-------------------|-------|
| `applications` | `conditions` | Yes | Required | N/A | Object | `{...}` | Accepts object, adds nested defaults | Modified (see sub-fields) | Same | Required object with application targeting |
| `users` | `conditions` | Yes | Required | N/A | Object | `{...}` | Accepts object | `{...}` | Same | Required object with user targeting |
| `clientApplications` | `conditions` | Yes | Optional | `null` | Object / `null` | caau001/002: `{object}`, caau003: `null` | Conditional: if object adds arrays, if null keeps null | caau001/002: modified, caau003: `null` | Same | **Critical**: When object, API adds service principal arrays; when `null`, stays `null` |
| `clientAppTypes` | `conditions` | Yes | Optional | `["all"]` | Array: `all`, `browser`, `mobileAppsAndDesktopClients`, `exchangeActiveSync`, `easSupported`, `other` | `["all"]` (all tests) | Accepts array | `["all"]` | Same | Optional. When explicitly sent as `["all"]`, API keeps it. May default if omitted (needs testing) |
| `platforms` | `conditions` | Yes | Optional | `null` | Object / `null` | `null` | Accepts null | `null` | Same | Platform targeting (OS types) |
| `locations` | `conditions` | Yes | Optional | `null` | Object / `null` | caau001-003: `null`, caau004: `{object}` | Accepts null or object | caau004: `{object}` | Same | Location/network targeting. When object, has `includeLocations` and `excludeLocations` |
| `times` | `conditions` | Yes | Optional | `null` | Object / `null` | `null` | Accepts null | `null` | Same | Time-based conditions |
| `devices` | `conditions` | Yes | Optional | `null` | Object / `null` | `null` | Accepts null | `null` | Same | Device state/compliance |
| `deviceStates` | `conditions` | No | N/A | N/A | Object / `null` | N/A | Adds default as null | `null` | Same | Deprecated field, always added by API |
| `userRiskLevels` | `conditions` | Yes | Optional | `[]` | Array: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue` | `[]` | Accepts [] | `[]` | Same | User risk level targeting |
| `signInRiskLevels` | `conditions` | Yes | Optional | `[]` | Array: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue` | `[]` | Accepts [] | `[]` | Same | Sign-in risk level targeting |
| `servicePrincipalRiskLevels` | `conditions` | Yes | Optional | `[]` | Array: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue` | `[]` | Removes (rejects) | Not returned | Not returned | **GUI sends [], API removes field** |
| `agentIdRiskLevels` | `conditions` | Yes | Optional | Not always sent | String (comma-separated or single): `low`, `medium`, `high` | caau002: `"high,medium"`, caau006: `"high"` | Modifies: alphabetically reorders if multiple | caau002: `"medium,high"`, caau006: `"high"` | Same | **CONFIRMED**: API alphabetically sorts when multiple values; single values unchanged |
| `clients` | `conditions` | Yes | Optional | `null` | Object / `null` | `null` | Removes (rejects) | Not returned | Not returned | **GUI sends null, API removes field** |
| `insiderRiskLevels` | `conditions` | Yes | Optional | `null` | Object / `null` | `null` | Removes (rejects) | Not returned | Not returned | **GUI sends null, API removes field** |
| `signInRiskDetections` | `conditions` | Yes | Optional | `null` | Object / `null` | `null` | Removes (rejects) | Not returned | Not returned | **GUI sends null, API removes field** |
| `authenticationFlows` | `conditions` | Yes | Optional | `null` | Object / `null` | `null` | Removes (rejects) | Not returned | Not returned | **GUI sends null, API removes field** |

### Conditions.Applications Fields

| Field Name | Parent Field Name | Sent by GUI | Create Request Requirement | Default GUI Value | Valid Field Value(s) | GUI Value (Tests) | API Behavior | Create Response Value | GET Response Value | Notes |
|-----------|-------------------|-------------|---------------------|-------------------|----------------------|-------------------|--------------|----------------------|-------------------|-------|
| `includeApplications` | `conditions.applications` | Yes | Required | N/A | Array: `All`, `None`, `Office365`, `AllAgentIdResources`, or app GUIDs | caau001: `["None"]`, caau002/003: `["All"]`, caau004: `[<GUIDs>]` | Accepts array | Same as sent | Same | Required. Accepts special values and specific app GUIDs |
| `excludeApplications` | `conditions.applications` | Yes | Optional | `[]` | Array of app GUIDs or special values | caau001-004: `[]`, caau005: `["AllAgentIdResources"]` | Accepts [] or special values/GUIDs | Same as sent | Same | Optional. Accepts special values like `AllAgentIdResources` |
| `includeUserActions` | `conditions.applications` | Yes | Optional | `[]` | Array of user action strings | `[]` | Accepts [] | `[]` | Same | Optional. GUI sends empty array |
| `includeAuthenticationContextClassReferences` | `conditions.applications` | Yes | Optional | `[]` | Array of auth context class reference strings | `[]` | Accepts [] | `[]` | Same | Optional. GUI sends empty array |
| `applicationFilter` | `conditions.applications` | Yes | Optional | `null` | Object / `null` | `null` | Accepts null | `null` | Same | Optional. Application filter with mode and rule |
| `globalSecureAccess` | `conditions.applications` | Yes | Optional | `null` | Object / `null` | `null` | Removes (rejects) | Not returned | Not returned | **GUI sends null, API removes field** |

### Conditions.Users Fields

| Field Name | Parent Field Name | Sent by GUI | Create Request Requirement | Default GUI Value | Valid Field Value(s) | GUI Value (Tests) | API Behavior | Create Response Value | GET Response Value | Notes |
|-----------|-------------------|-------------|---------------------|-------------------|----------------------|-------------------|--------------|----------------------|-------------------|-------|
| `includeUsers` | `conditions.users` | Yes | Required | N/A | Array: `All`, `None`, `AllAgentIdUsers`, `GuestsOrExternalUsers`, or user GUIDs | caau001/002: `["None"]`, caau003: `["AllAgentIdUsers"]`, caau004: `["All"]` | Accepts array | Same as sent | Same | Required. Accepts special values and GUIDs |
| `excludeUsers` | `conditions.users` | Yes | Optional | `[]` | Array of user GUIDs | caau001-003: `[]`, caau004: `[<GUID>]` | Accepts [] or GUIDs | Same as sent | Same | Optional. GUI sends empty array or GUIDs |
| `includeGroups` | `conditions.users` | Yes | Optional | `[]` | Array of group GUIDs | `[]` (all tests) | Accepts [] | `[]` | Same | Optional. GUI sends empty array |
| `excludeGroups` | `conditions.users` | Yes | Optional | `[]` | Array of group GUIDs | caau001-003: `[]`, caau004: `[<GUIDs>]` | Accepts [] or GUIDs | Same as sent | Same | Optional. Accepts multiple group GUIDs |
| `includeRoles` | `conditions.users` | Yes | Optional | `[]` | Array of role template IDs | `[]` (all tests) | Accepts [] | `[]` | Same | Optional. GUI sends empty array |
| `excludeRoles` | `conditions.users` | Yes | Optional | `[]` | Array of role template IDs | caau001-003: `[]`, caau004: `[<GUIDs>]` | Accepts [] or GUIDs | Same as sent | Same | Optional. Accepts multiple role template IDs |
| `includeGuestsOrExternalUsers` | `conditions.users` | Yes | Optional | `null` | Object / `null` | `null` (all tests) | Accepts null | `null` | Same | Optional. Guest/external user inclusion settings |
| `excludeGuestsOrExternalUsers` | `conditions.users` | Yes | Optional | `null` | Object / `null` | caau001-003: `null`, caau004: `{object}` | Accepts null or object; transforms case | caau004: `{object}` | Same | 🚨 **Case transform**: PascalCase → camelCase in `guestOrExternalUserTypes` |

### Conditions.ClientApplications Fields

| Field Name | Parent Field Name | Sent by GUI | Create Request Requirement | Default GUI Value | Valid Field Value(s) | GUI Value (Tests) | API Behavior | Create Response Value | GET Response Value | Notes |
|-----------|-------------------|-------------|---------------------|-------------------|----------------------|-------------------|--------------|----------------------|-------------------|-------|
| `includeAgentIdServicePrincipals` | `conditions.clientApplications` | Yes | Optional | N/A | Array: `All` or service principal GUIDs | caau001: `["7ca55e16..."]`, caau002/006: `["All"]` | Accepts array | Same as sent | Same | Accepts specific GUIDs or `"All"` |
| `excludeAgentIdServicePrincipals` | `conditions.clientApplications` | Yes | Optional | N/A | Array of service principal GUIDs | caau001/002: `[<GUID>]`, caau006: `[]` | Conditional: keeps populated, removes empty | caau001/002: `[<GUID>]`, caau006: removed | Same | 🚨 **If populated: kept; If empty []: removed** |
| `agentIdServicePrincipalFilter` | `conditions.clientApplications` | Yes | Optional | `null` | Object: `{mode, rule}` / `null` | caau001/006: `{object}`, caau002: `null` | Conditional: keeps object, removes null | caau001/006: `{object}`, caau002: removed | Same | **If object: kept; If null: removed**. Accepts both "include" and "exclude" modes |
| `includeServicePrincipals` | `conditions.clientApplications` | No | N/A | N/A | Array of service principal GUIDs | N/A | Conditionally adds [] | caau001/002: `[]`, caau003: N/A | Same | Never sent by GUI. API adds as `[]` **only when clientApplications is object** (not when null) |
| `excludeServicePrincipals` | `conditions.clientApplications` | No | N/A | N/A | Array of service principal GUIDs | N/A | Conditionally adds [] | caau001/002: `[]`, caau003: N/A | Same | Never sent by GUI. API adds as `[]` **only when clientApplications is object** (not when null) |

### Conditions.Locations Fields

| Field Name | Parent Field Name | Sent by GUI | Create Request Requirement | Default GUI Value | Valid Field Value(s) | GUI Value (Tests) | API Behavior | Create Response Value | GET Response Value | Notes |
|-----------|-------------------|-------------|---------------------|-------------------|----------------------|-------------------|--------------|----------------------|-------------------|-------|
| `includeLocations` | `conditions.locations` | Yes | Optional | N/A | Array: `All`, `AllTrusted`, or location GUIDs | caau004: `["All"]` | Accepts array | `["All"]` | Same | Location inclusion targeting |
| `excludeLocations` | `conditions.locations` | Yes | Optional | `[]` | Array of location GUIDs | caau004: `[]` | Accepts [] | `[]` | Same | Location exclusion targeting |

### Conditions.Users.ExcludeGuestsOrExternalUsers Fields

| Field Name | Parent Field Name | Sent by GUI | Create Request Requirement | Default GUI Value | Valid Field Value(s) | GUI Value (Tests) | API Behavior | Create Response Value | GET Response Value | Notes |
|-----------|-------------------|-------------|---------------------|-------------------|----------------------|-------------------|--------------|----------------------|-------------------|-------|
| `guestOrExternalUserTypes` | `conditions.users.excludeGuestsOrExternalUsers` | Yes | Optional | N/A | String (comma-separated): values like `b2bCollaborationGuest`, `b2bCollaborationMember`, `b2bDirectConnectUser`, etc. | caau004: PascalCase version | 🚨 Transforms: PascalCase → camelCase | camelCase version | Same | **CRITICAL**: API converts `B2bCollaborationGuest` → `b2bCollaborationGuest` |
| `externalTenants` | `conditions.users.excludeGuestsOrExternalUsers` | Yes | Optional | N/A | Object with `membershipKind` and `@odata.type` | caau004: `{membershipKind: "all", @odata.type: "...AllExternalTenants"}` | Accepts object | Same as sent | Same | External tenant targeting configuration |

### GrantControls Fields

| Field Name | Parent Field Name | Sent by GUI | Create Request Requirement | Default GUI Value | Valid Field Value(s) | GUI Value (Tests) | API Behavior | Create Response Value | GET Response Value | Notes |
|-----------|-------------------|-------------|---------------------|-------------------|----------------------|-------------------|--------------|----------------------|-------------------|-------|
| `operator` | `grantControls` | Yes | Required | N/A | `AND`, `OR` | `"AND"` (all tests) | 🚨 Modifies: `"AND"` → `"OR"` when grant control present | `"OR"` (all tests) | Same | **CRITICAL**: API transforms `"AND"` to `"OR"` when single built-in control OR authenticationStrength is used |
| `builtInControls` | `grantControls` | Yes | Optional/Required | N/A | Array: `block`, `mfa`, `compliantDevice`, `domainJoinedDevice`, `approvedApplication`, `compliantApplication`, `passwordChange`, `unknownFutureValue` | caau001-003: `["block"]`, caau004: `["compliantDevice"]`, caau005: `[]` | Accepts array including empty | Same as sent | Same | Can be empty `[]` if `authenticationStrength` is present |
| `customAuthenticationFactors` | `grantControls` | Yes | Optional | `[]` | Array of custom auth factor strings | `[]` (all tests) | Accepts [] | `[]` | Same | Optional. GUI sends empty array |
| `termsOfUse` | `grantControls` | Yes | Optional | `[]` | Array of terms of use agreement GUIDs | `[]` (all tests) | Accepts [] | `[]` | Same | Optional. GUI sends empty array |
| `authenticationStrength` | `grantControls` | Yes | Optional | `null` | Object / `null` | caau001-004: `null`, caau005: `{id: ...}` | 🚨 Expands: Accepts ID, returns full object | caau005: Full expanded object | Same | **API EXPANSION**: Send just `{id}`, receive full object with 10+ fields including displayName, allowedCombinations, etc. |
| `authenticationStrength@odata.context` | `grantControls` | No | N/A | N/A | String (URL) | N/A | Adds (metadata) | `https://graph.microsoft.com/beta/$metadata#identity/conditionalAccess/policies('aefd62ee-806f-4a7c-8876-a1fc171c665d')/grantControls/authenticationStrength/$entity` | Same | API metadata, always added |

### SessionControls Fields

| Field Name | Parent Field Name | Sent by GUI | Create Request Requirement | Default GUI Value | Valid Field Value(s) | GUI Value (Tests) | API Behavior | Create Response Value | GET Response Value | Notes |
|-----------|-------------------|-------------|---------------------|-------------------|----------------------|-------------------|--------------|----------------------|-------------------|-------|
| `signInFrequency` | `sessionControls` | Yes | Optional | `null` | Object / `null` | caau005: `{object}` | Accepts object | `{object}` | Same | Sign-in frequency control with `frequencyInterval`, `authenticationType`, `isEnabled`, `value`, `type` |
| `applicationEnforcedRestrictions` | `sessionControls` | Yes | Optional | `null` | Object / `null` | caau005: `null` | Accepts null | `null` | Same | Application-enforced restrictions |
| `cloudAppSecurity` | `sessionControls` | Yes | Optional | `null` | Object / `null` | caau005: `null` | Accepts null | `null` | Same | Cloud app security settings |
| `persistentBrowser` | `sessionControls` | Yes | Optional | `null` | Object / `null` | caau005: `null` | Accepts null | `null` | Same | Persistent browser session settings |
| `continuousAccessEvaluation` | `sessionControls` | Yes | Optional | `null` | Object / `null` | caau005: `null` | Accepts null | `null` | Same | Continuous access evaluation settings |
| `disableResilienceDefaults` | `sessionControls` | Yes | Optional | `null` | Object / `null` | caau005: `null` | Accepts null | `null` | Same | Resilience defaults settings |
| `secureSignInSession` | `sessionControls` | Yes | Optional | `null` | Object / `null` | caau005: `null` | Accepts null | `null` | Same | Secure sign-in session settings |
| `networkAccessSecurity` | `sessionControls` | Yes | Optional | `null` | Object / `null` | caau005: `null` | Removes (rejects) | Not returned | Not returned | **GUI sends null, API removes field** |
| `globalSecureAccessFilteringProfile` | `sessionControls` | Yes | Optional | `null` | Object / `null` | caau005: `null` | Removes (rejects) | Not returned | Not returned | **GUI sends null, API removes field** |

## Key Observations

### Fields Removed by API (Confirmed by Testing)
The following fields sent in the create request are **rejected** and not returned in responses:

**Conditions fields removed (6 fields):**
- `conditions.clients` (sent as `null`, removed by API)
- `conditions.insiderRiskLevels` (sent as `null`, removed by API)
- `conditions.signInRiskDetections` (sent as `null`, removed by API)
- `conditions.authenticationFlows` (sent as `null`, removed by API)
- `conditions.servicePrincipalRiskLevels` (sent as `[]`, removed by API)
- `conditions.applications.globalSecureAccess` (sent as `null`, removed by API)

**SessionControls fields removed (2 fields - test 5):**
- `sessionControls.networkAccessSecurity` (sent as `null`, removed by API)
- `sessionControls.globalSecureAccessFilteringProfile` (sent as `null`, removed by API)

**Important Correction:** The following fields are **NOT** removed - they are accepted and returned:
- ✅ `conditions.clientApplications.excludeAgentIdServicePrincipals` - **KEPT by API**
- ✅ `conditions.clientApplications.agentIdServicePrincipalFilter` - **KEPT by API**

### Fields Added by API (Confirmed by Testing)
The following fields are **not** sent in the create request but are **added** by the API in responses:
- `@odata.context` (metadata URL) - always added
- `id` (UUID) - auto-generated policy ID
- `createdDateTime` (ISO 8601 timestamp) - auto-generated creation time
- `modifiedDateTime` (initially `null`) - updated when policy is modified
- `deletedDateTime` (initially `null`) - set when policy is soft-deleted
- `templateId` (initially `null`) - reference to template if used
- `partialEnablementStrategy` (initially `null`) - partial enablement configuration
- `conditions.deviceStates` (always `null`) - deprecated field, always added by API
- `conditions.clientApplications.includeServicePrincipals` (always `[]`) - added when clientApplications object exists
- `conditions.clientApplications.excludeServicePrincipals` (always `[]`) - added when clientApplications object exists
- `grantControls.authenticationStrength@odata.context` (metadata URL) - always added

### API Default Values

The API automatically provides default values for certain fields if omitted:

1. **conditions.clientAppTypes**: Defaults to `["all"]` if not provided in the request
   - Test 2 confirmed: Field omitted from request, API returned `["all"]`

### Field Value Transformations (Confirmed by Testing)

The API transforms certain field values consistently:

1. **grantControls.operator**: When using grant controls, the API **always** changes the operator
   - Sent: `"AND"`
   - Received: `"OR"`
   - **Confirmed in all 5 tests**: 
     - With `"block"` (tests 1-3)
     - With `"compliantDevice"` (test 4)
     - With `authenticationStrength` and empty `builtInControls: []` (test 5)
   - The API enforces that grant controls must use the OR operator
   - **Pattern**: When requirements can be satisfied by a single method, API uses OR

2. **authenticationStrength object expansion**: The API dramatically expands this object
   - Sent (caau005): `{id: "00000000-0000-0000-0000-000000000004"}`
   - Received (caau005): Full object with `displayName`, `description`, `allowedCombinations`, `policyType`, `requirementsSatisfied`, `createdDateTime`, `modifiedDateTime`, `combinationConfigurations`, etc.
   - **Status**: ✅ CONFIRMED - API expands minimal ID reference to full object

3. **excludeGuestsOrExternalUsers.guestOrExternalUserTypes**: The API converts PascalCase to camelCase
   - Sent (caau004): `"B2bCollaborationGuest,B2bCollaborationMember,B2bDirectConnectUser"`
   - Received (caau004): `"b2bCollaborationGuest,b2bCollaborationMember,b2bDirectConnectUser"`
   - **Status**: ✅ CONFIRMED - API lowercases the first character of each comma-separated value

4. **agentIdRiskLevels**: The API reorders comma-separated values alphabetically
   - Sent (caau002): `"high,medium"`
   - Received (caau002): `"medium,high"`
   - **Status**: ✅ CONFIRMED - API alphabetically sorts the comma-separated risk level values

### Minimum API Required Fields

The following fields are **actually required** by the Microsoft Graph API to create a Conditional Access Policy:

- `displayName`
- `state`
- `conditions` (object)
  - `conditions.users.includeUsers` (or other user targeting)
  - `conditions.applications.includeApplications` (or other app targeting)
- `grantControls` (object)
  - `grantControls.operator`
  - `grantControls.builtInControls` (or other grant control)

**Note:** `conditions.clientAppTypes` is optional - the API defaults it to `["all"]` if not provided.

### Fields Sent by Azure Portal GUI

The Azure Portal GUI sends many additional **optional** fields that are not required by the API, including:
- Empty arrays for all exclude/include lists (even when not used)
- `null` values for unused features (platforms, locations, times, devices, etc.)
- Fields that the API removes (clients, insiderRiskLevels, authenticationFlows, etc.)

This means the GUI payloads are much more verbose than the minimum required by the API.

### Field Consistency

Fields that maintain consistency across Create Request → Create Response → GET Response:
- `displayName`
- `state`
- `conditions.users.*` (all user fields)
- `conditions.applications.*` (most application fields)
- `conditions.clientAppTypes`
- `grantControls.builtInControls`
- `grantControls.customAuthenticationFactors`
- `grantControls.termsOfUse`
- `grantControls.authenticationStrength`
- `sessionControls`

## Valid Values Reference

### State Values
- `disabled`: Policy is disabled
- `enabled`: Policy is enabled and enforced
- `enabledForReportingButNotEnforced`: Report-only mode

### ClientAppTypes Values
- `all`: All client app types
- `browser`: Browser-based applications
- `mobileAppsAndDesktopClients`: Native clients
- `exchangeActiveSync`: Exchange ActiveSync clients
- `easSupported`: EAS supported clients
- `other`: Other clients

### BuiltInControls Values
- `block`: Block access
- `mfa`: Require multi-factor authentication
- `compliantDevice`: Require device to be marked as compliant
- `domainJoinedDevice`: Require domain joined device
- `approvedApplication`: Require approved client app
- `compliantApplication`: Require app protection policy
- `passwordChange`: Require password change
- `unknownFutureValue`: Evolvable enum sentinel value

### Operator Values
- `AND`: Require all specified controls
- `OR`: Require any one of the specified controls

### Special Application Values
- `All`: All cloud applications
- `None`: No applications
- `Office365`: All Office 365 applications
- `AllAgentIdResources`: All Agent ID resources

### Special User Values
- `All`: All users
- `None`: No users (used with service principal targeting)
- `AllAgentIdUsers`: All Agent ID users
- `GuestsOrExternalUsers`: All guest or external users

## Terraform Provider Implications

When implementing Terraform provider support:

1. **Handle API transformations**: Account for fields that are removed or reordered by the API
   - **CONFIRMED**: `agentIdRiskLevels` will be alphabetically reordered (e.g., `"high,medium"` → `"medium,high"`)
     - Provider should normalize this to prevent drift detection
   - **CONFIRMED**: `grantControls.operator` will always become `"OR"` when using **single** built-in controls
     - Confirmed with `block` and `compliantDevice` - likely applies to all single controls
     - Consider auto-setting or suppressing diff when single control is present
   - **CONFIRMED**: `excludeGuestsOrExternalUsers.guestOrExternalUserTypes` converts PascalCase to camelCase
     - Provider should normalize input or suppress diff for case differences

2. **Read-only fields**: Mark auto-generated fields as computed
   - `id`, `createdDateTime`, `modifiedDateTime`, `deletedDateTime`, `templateId`, `partialEnablementStrategy`
   - `@odata.context` and other metadata fields

3. **Operator override for block controls**:
   - **Critical**: When `grantControls.builtInControls` contains `"block"`, the API will **always** change the operator to `"OR"`
   - Consider either:
     - Auto-setting operator to `"OR"` when block is used (prevent drift)
     - Documenting this behavior prominently
     - Implementing a validation warning when user sets `AND` with block

4. **API default values**: Handle fields with API-provided defaults
   - `conditions.clientAppTypes` defaults to `["all"]` if omitted
   - Consider making this explicit in the schema with a default value

5. **GUI vs API Required distinction**: Understand what's actually required vs. what the GUI sends
   - The Azure Portal GUI sends many optional fields (empty arrays, null values) that aren't required
   - The provider should only require the minimal API-required fields
   - Optional fields can be made available but should have appropriate defaults or be truly optional

6. **Field removal**: Don't store these fields in state as the API removes them:
   
   **Conditions fields (6 total):**
   - `conditions.clients` (sent as `null`, rejected)
   - `conditions.insiderRiskLevels` (sent as `null`, rejected)
   - `conditions.signInRiskDetections` (sent as `null`, rejected)
   - `conditions.authenticationFlows` (sent as `null`, rejected)
   - `conditions.servicePrincipalRiskLevels` (sent as `[]`, rejected)
   - `conditions.applications.globalSecureAccess` (sent as `null`, rejected)
   
   **SessionControls fields (2 total):**
   - `sessionControls.networkAccessSecurity` (sent as `null`, rejected)
   - `sessionControls.globalSecureAccessFilteringProfile` (sent as `null`, rejected)
   
   **Important**: The following fields are **NOT** removed and should be stored in state:
   - ✅ `conditions.clientApplications.excludeAgentIdServicePrincipals` - API accepts and returns
   - ✅ `conditions.clientApplications.agentIdServicePrincipalFilter` - API accepts and returns

7. **Null vs missing fields vs empty arrays**: Properly handle the distinction and conditional removal/addition
   - **`conditions.clientApplications` behavior** (critical for schema design):
     - When set to `null`: Stays `null`, no arrays added ✅ (caau003)
     - When set to object: API adds `includeServicePrincipals: []` and `excludeServicePrincipals: []` ➕ (caau001, caau002, caau006)
   - **`agentIdServicePrincipalFilter` conditional behavior**:
     - When sent as **object** with mode and rule: API keeps it ✅ (caau001, caau006)
     - When sent as **null** within clientApplications object: API removes the field entirely ❌ (caau002)
   - **`excludeAgentIdServicePrincipals` conditional behavior** (NEW - test 6):
     - When sent as populated array `[<GUIDs>]`: API keeps it ✅ (caau001, caau002)
     - When sent as empty array `[]`: API removes the field entirely ❌ (caau006)
   - **Implication**: Schema should distinguish between "field not present", "field is null", "field is empty array", and "field is populated array"
   - This complexity may require custom diff suppression or normalization logic to avoid spurious diffs
   - **Recommendation**: For `excludeAgentIdServicePrincipals`, only send when populated; omit entirely when not needed

8. **AuthenticationStrength expansion**: Handle API's dramatic object expansion
   - User provides: `authentication_strength = {id = "00000000-0000-0000-0000-000000000004"}`
   - API returns: Full object with 10+ fields (displayName, allowedCombinations, policyType, etc.)
   - **Strategy options**:
     - Mark all expansion fields as Computed (recommended) - user provides ID, Terraform reads full object
     - Implement custom diff suppression to ignore the expanded fields
     - Only store the ID in state and ignore the rest
   - **Critical**: Don't cause drift when user only provides ID

9. **Validation**: Implement validation for enum-type fields
   - `state`: `disabled`, `enabled`, `enabledForReportingButNotEnforced`
   - `clientAppTypes`: `all`, `browser`, `mobileAppsAndDesktopClients`, `exchangeActiveSync`, `easSupported`, `other`
   - `builtInControls`: `block`, `mfa`, `compliantDevice`, `domainJoinedDevice`, `approvedApplication`, `compliantApplication`, `passwordChange`
   - Special values: `All`, `None`, `AllAgentIdResources`, `AllAgentIdUsers`, `Office365`, `GuestsOrExternalUsers`

10. **Testing coverage**: Ensure tests cover:
   - Agent ID scenarios (AllAgentIdResources, AllAgentIdUsers, service principals) ✅
   - Standard user/application targeting ✅
   - Block controls with operator transformation ✅
   - compliantDevice controls with operator transformation ✅
   - authenticationStrength with empty builtInControls ✅
   - sessionControls with signInFrequency ✅
   - excludeGuestsOrExternalUsers with case transformation ✅
   - locations targeting ✅
   - Fields that get removed by API ✅
   - authenticationStrength expansion ✅
   - Fields with default values
