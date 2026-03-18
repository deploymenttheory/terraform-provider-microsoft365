# Developer Documentation: autopatch_update_policy Resource

## Overview
This document records the peculiarities, API quirks, and implementation challenges discovered during the development of the `microsoft365_graph_beta_windows_updates_update_policy` resource.

## Table of Contents
1. [API Endpoint Behavior](#api-endpoint-behavior)
2. [Field Validation Rules](#field-validation-rules)
3. [SDK Normalization Issues](#sdk-normalization-issues)
4. [CREATE vs UPDATE Differences](#create-vs-update-differences)
5. [Write-Only Fields](#write-only-fields)
6. [Testing Methodology](#testing-methodology)

---

## API Endpoint Behavior

### Base Endpoint
- **CREATE**: `POST /admin/windows/updates/updatePolicies`
- **READ**: `GET /admin/windows/updates/updatePolicies/{id}`
- **UPDATE**: `PATCH /admin/windows/updates/updatePolicies/{id}`
- **DELETE**: `DELETE /admin/windows/updates/updatePolicies/{id}`

### Misleading Error Messages
The API returns confusing error messages that don't always match the actual problem:
- Error message: `"Schema validation failed for resource 'deploymentPolicy'"`
- Actual resource type: `updatePolicy`
- This occurs when invalid filter types or UPDATE operations include restricted fields

---

## Field Validation Rules

### 1. Filter Type (`compliance_change_rules.content_filter.filter_type`)

**Documentation Says**: Both `driverUpdateFilter` and `windowsUpdateFilter` are valid

**Reality**:
- ✅ **`driverUpdateFilter`**: ONLY valid filter type
- ❌ **`windowsUpdateFilter`**: COMPLETELY INVALID - causes 400 error for ALL duration values (P1D through P180D tested)

**Testing Evidence**:
- Tested all duration values P1D to P180D with `windowsUpdateFilter`
- Result: 100% failure rate (180/180 failed)
- Error: `"Schema validation failed for resource 'updatePolicy'"`

**Implementation**:
```go
Validators: []validator.String{
    stringvalidator.OneOf("driverUpdateFilter"),
}
```

---

### 2. Duration Before Deployment Start (`compliance_change_rules.duration_before_deployment_start`)

**Documentation Says**: ISO 8601 duration format (no specific range mentioned)

**Reality**:
- ✅ **Valid Range**: P1D to P30D (1 to 30 days only)
- ❌ **Invalid**: P31D and above
- Error for invalid values: `"Parameter 'durationBeforeDeploymentStart' in payload has a value that does not match schema"`

**Testing Evidence**:
- Tested P1D through P180D with `driverUpdateFilter`
- Valid: 30 values (P1D to P30D)
- Invalid: 150 values (P31D to P180D)

**Implementation**:
```go
Validators: []validator.String{
    stringvalidator.RegexMatches(
        regexp.MustCompile(`^P([1-9]|[12][0-9]|30)D$`),
        "must be in ISO 8601 duration format P#D where # is between 1 and 30 (e.g., P7D, P14D, P30D)",
    ),
}
PlanModifiers: []planmodifier.String{
    stringplanmodifier.RequiresReplace(),
}
```

**Additional Notes**:
- This field is **immutable** after creation
- Changes require resource replacement (destroy + create)
- Must use `additionalData` in construct.go to bypass SDK normalization (see SDK Issues section)

---

### 3. Duration Between Offers (`deployment_settings.schedule.gradual_rollout.duration_between_offers`)

**Documentation Says**: ISO 8601 duration format (no specific range mentioned)

**Reality**:
- ✅ **Valid Range**: P1D to P30D (1 to 30 days only)
- ❌ **Invalid**: P31D and above

**Testing Evidence**:
- Same validation pattern as `durationBeforeDeploymentStart`
- Maximum value is P30D (30 days)

**Implementation**:
```go
Validators: []validator.String{
    stringvalidator.RegexMatches(
        regexp.MustCompile(`^P([1-9]|[12][0-9]|30)D$`),
        "must be in ISO 8601 duration format P#D where # is between 1 and 30 (e.g., P1D, P7D, P30D)",
    ),
}
```

**Additional Notes**:
- This field IS mutable (can be updated)
- Must use `additionalData` in construct.go to bypass SDK normalization (see SDK Issues section)

---

### 4. Devices Per Offer (`deployment_settings.schedule.gradual_rollout.devices_per_offer`)

**Documentation Says**: Integer representing number of devices

**Reality**:
- ✅ **Valid Range**: 1 to at least 1,000,000 (possibly higher, not tested beyond)
- ❌ **Invalid**: 0 or negative values

**Testing Evidence**:
- Tested values: 1, 10, 50, 100, 500, 1000, 5000, 10000, 50000, 100000, 150000, 200000, 250000, 500000, 1000000
- Result: All valid (100% success rate)

**Implementation**:
```go
Validators: []validator.Int32{
    int32validator.AtLeast(1),
}
```

**Additional Notes**:
- This field IS mutable (can be updated)
- API key is `devicesPerOffer` (plural), not `devicePerOffer` (singular as shown in some MS docs)
- Must use `additionalData` in construct.go with correct plural key name

---

### 5. Start Date Time (`deployment_settings.schedule.start_date_time`)

**Documentation Says**: ISO 8601 datetime format

**Reality**:
- **Optional**: Can be completely omitted
- **NOT Computed**: When omitted, API returns null (not auto-generated)
- **Accepts null**: Explicitly setting to null is valid
- **Accepts future dates**: Can set dates in the future

**Testing Evidence**:
- Test 1 (omitted): ✅ Success - returns null
- Test 2 (explicit future date): ✅ Success - returns provided value
- Test 3 (explicit null): ✅ Success - returns null

**Implementation**:
```go
Optional: true,
// NOT Computed: true
```

---

## SDK Normalization Issues

### Problem: ISO Duration Normalization

The Microsoft Kiota SDK automatically normalizes ISO 8601 durations during deserialization:
- User sends: `P7D` (7 days)
- API returns: `P7D` (7 days)
- SDK deserializes to: `P1W` (1 week)
- Terraform sees: `P1W` in state vs `P7D` in config → diff detected

**Why This Happens**:
The SDK's `ISODuration` type normalizes day-based durations to week-based equivalents:
- P7D → P1W
- P14D → P2W
- P21D → P3W
- P28D → P4W

### Solution 1: Bypass SDK Normalization on SEND (construct.go)

Use `additionalData` to send raw string values instead of using SDK setters:

```go
// WRONG: This causes normalization
duration := abstractions.NewISODuration()
duration.Parse("P7D")  // Gets normalized to P1W internally
complianceRule.SetDurationBeforeDeploymentStart(duration)

// CORRECT: Use additionalData to send raw string
additionalData := complianceRule.GetAdditionalData()
if additionalData == nil {
    additionalData = make(map[string]any)
}
additionalData["durationBeforeDeploymentStart"] = "P7D"  // Raw string, no normalization
complianceRule.SetAdditionalData(additionalData)
```

**Applied to**:
- `durationBeforeDeploymentStart` in `complianceChangeRules`
- `durationBetweenOffers` in `gradualRollout`

### Solution 2: Denormalize on READ (state.go)

Convert SDK-normalized values back to day-based format to match user config:

```go
func denormalizeISODuration(duration string) string {
    // Pattern for week-based durations (e.g., P1W, P2W)
    weekPattern := regexp.MustCompile(`^P(\d+)W$`)
    if matches := weekPattern.FindStringSubmatch(duration); len(matches) == 2 {
        weeks, err := strconv.Atoi(matches[1])
        if err != nil {
            return duration
        }
        days := weeks * 7
        return "P" + strconv.Itoa(days) + "D"
    }
    return duration
}
```

**Conversions**:
- P1W → P7D
- P2W → P14D
- P3W → P21D
- P4W → P28D

**Applied to**:
- `durationBeforeDeploymentStart` when reading from API
- `durationBetweenOffers` when reading from API

---

## CREATE vs UPDATE Differences

### Fields Allowed in CREATE Only

The following fields MUST be included in CREATE but MUST NOT be included in UPDATE:

1. **`audience`** (with `id` field)
   - Required for CREATE
   - Rejected in UPDATE
   - Immutable after creation

2. **`complianceChanges`** (array with ContentApproval object)
   - Required for CREATE (must be `[{ "@odata.type": "#microsoft.graph.windowsUpdates.contentApproval" }]`)
   - Rejected in UPDATE
   - Write-only field (see Write-Only Fields section)

3. **`complianceChangeRules`** (array of ContentApprovalRule objects)
   - Optional for CREATE
   - **REJECTED in UPDATE** (despite documentation saying it's updatable)
   - Immutable after creation

**Testing Evidence**:
- UPDATE with `complianceChangeRules`: ❌ 400 error "Schema validation failed for resource 'deploymentPolicy'"
- UPDATE with only `deploymentSettings`: ✅ Success

### Fields Allowed in UPDATE

Only `deploymentSettings` can be updated:
- ✅ `deploymentSettings.schedule.startDateTime`
- ✅ `deploymentSettings.schedule.gradualRollout.durationBetweenOffers`
- ✅ `deploymentSettings.schedule.gradualRollout.devicesPerOffer`

### Implementation

```go
func constructResource(ctx context.Context, data *WindowsUpdatesAutopatchUpdatePolicyResourceModel, isUpdate bool) (graphmodelswindowsupdates.UpdatePolicyable, error) {
    requestBody := graphmodelswindowsupdates.NewUpdatePolicy()
    
    // For UPDATE operations, skip these fields
    if !isUpdate {
        // Set audience
        // Set complianceChanges
        // Set complianceChangeRules
    }
    
    // Always set deploymentSettings (works for both CREATE and UPDATE)
    if !data.DeploymentSettings.IsNull() {
        // ... set deployment settings
    }
    
    return requestBody, nil
}
```

**CRUD Operations**:
```go
// Create
requestBody, err := constructResource(ctx, &object, false)

// Update
requestBody, err := constructResource(ctx, &plan, true)
```

---

## Write-Only Fields

### `compliance_changes` Field

**Behavior**:
- ✅ **Required in CREATE request**: Must be set to `true` and include ContentApproval object in payload
- ❌ **NOT returned in GET response**: API does not include this field when reading the resource
- ⚠️ **Causes ImportStateVerify failure**: If not ignored, import tests fail

**Implementation**:

1. **Schema Definition**:
```go
"compliance_changes": schema.BoolAttribute{
    Required:            true,
    MarkdownDescription: "Enable compliance changes (content approvals) for this policy. Must be set to `true` to create content approvals.",
}
```

2. **State Mapping** (state.go):
```go
// Note: compliance_changes is write-only and not returned by the API
// The value from config/state is preserved automatically
// DO NOT attempt to map this field from remoteResource
```

3. **Import Test Configuration**:
```go
ImportStateVerifyIgnore: []string{"timeouts", "compliance_changes"},
```

---

## @odata.type Annotations

### Required Annotations

Different objects require different `@odata.type` annotation formats:

1. **Root updatePolicy**: `#microsoft.graph.windowsUpdates.updatePolicy` (with `#` prefix)
2. **ContentApproval**: `#microsoft.graph.windowsUpdates.contentApproval` (with `#` prefix)
3. **ContentApprovalRule**: `#microsoft.graph.windowsUpdates.contentApprovalRule` (with `#` prefix)
4. **DriverUpdateFilter**: `#microsoft.graph.windowsUpdates.driverUpdateFilter` (with `#` prefix)
5. **WindowsUpdateFilter**: `#microsoft.graph.windowsUpdates.windowsUpdateFilter` (with `#` prefix, but INVALID - see validation rules)
6. **RateDrivenRolloutSettings**: `#microsoft.graph.windowsUpdates.rateDrivenRolloutSettings` (with `#` prefix)
7. **DeploymentSettings**: `microsoft.graph.windowsUpdates.deploymentSettings` (WITHOUT `#` prefix) - **OPTIONAL/NOT REQUIRED**

### Implementation Note
The SDK automatically adds most `@odata.type` annotations. The root `updatePolicy` type is set explicitly:

```go
odataType := "#microsoft.graph.windowsUpdates.updatePolicy"
requestBody.SetOdataType(&odataType)
```

---

## Field Naming Discrepancies

### devicesPerOffer vs devicePerOffer

**Microsoft Documentation Shows**: `devicePerOffer` (singular)

**Actual API Expects**: `devicesPerOffer` (plural)

**Testing Evidence**:
- Payload with `devicePerOffer` (singular): ❌ Field ignored by API
- Payload with `devicesPerOffer` (plural): ✅ Success

**Implementation**:
```go
// Use additionalData with correct plural key
additionalData["devicesPerOffer"] = rolloutData.DevicesPerOffer.ValueInt32()
```

**SDK Getter**:
```go
// SDK getter works correctly with plural form
if devices := rateDriven.GetDevicesPerOffer(); devices != nil {
    gradualRolloutData.DevicesPerOffer = types.Int32Value(*devices)
}
```

---

## Testing Methodology

### PowerShell Brute-Force Testing

Due to the numerous discrepancies between documentation and actual API behavior, we created PowerShell scripts to directly test the API:

#### Scripts Created:
1. **`Create-AutopatchUpdatePolicy.ps1`**: Test CREATE operations with various JSON payloads
2. **`Update-AutopatchUpdatePolicy.ps1`**: Test UPDATE operations with various JSON payloads
3. **`Create-DeploymentAudience.ps1`**: Helper to create test audiences
4. **`Delete-DeploymentAudience.ps1`**: Helper to cleanup test audiences
5. **`Run-AllPayloadTests.ps1`**: Orchestrate multiple CREATE payload tests
6. **`Run-AllUpdatePayloadTests.ps1`**: Orchestrate multiple UPDATE payload tests
7. **`Test-DurationValues.ps1`**: Test all P1D-P180D values for both filter types
8. **`Test-DurationBetweenOffers.ps1`**: Test all P1D-P180D values for durationBetweenOffers
9. **`Test-DevicesPerOffer.ps1`**: Test device count ranges

#### Test Payloads Created:
- **CREATE payloads**: `test_payloads/payload_01_*.json` through `payload_09_*.json`
- **UPDATE payloads**: `update_test_payloads/update_01_*.json` through `update_06_*.json`

### Key Discoveries from Testing:

1. **devicesPerOffer naming**: Discovered through payload_07_devicesPerOffer.json
2. **windowsUpdateFilter invalid**: Discovered by testing payload_09_windowsUpdateFilter.json
3. **UPDATE restrictions**: Discovered by testing 6 different UPDATE payloads
4. **Duration limits**: Discovered by testing all values P1D through P180D
5. **Device count range**: Discovered by testing 1 to 1,000,000 devices

---

## SDK Normalization Issues (Detailed)

### The Problem

The Microsoft Graph SDK for Go uses the `abstractions.ISODuration` type which automatically normalizes durations:

```go
// When SDK deserializes API response:
API returns:     "durationBeforeDeploymentStart": "P7D"
SDK parses:      ISODuration.Parse("P7D")
SDK normalizes:  P7D → P1W (because 7 days = 1 week)
SDK serializes:  duration.String() returns "P1W"
Terraform sees:  "P1W" in state vs "P7D" in config → DIFF!
```

### Why We Can't Use Plan Modifiers

**Attempted Solution 1: NormalizeISODurationModifier**
- Tried to normalize P7D → P1W in the plan
- Result: ❌ Failed - Terraform rejects plan modifiers that change config values
- Error: `"planned value ... does not match config value"`

**Attempted Solution 2: SuppressISODurationDiffModifier**
- Tried to suppress diff between P7D and P1W
- Result: ❌ Failed - Doesn't handle computed timestamp fields in sets
- Error: `"planned set element ... does not correlate with any element in actual"`

### The Working Solution: Two-Pronged Approach

#### 1. Send Raw Strings (construct.go)

Bypass SDK normalization when SENDING to API:

```go
// For durationBeforeDeploymentStart
if !rule.DurationBeforeDeploymentStart.IsNull() {
    rawDuration := rule.DurationBeforeDeploymentStart.ValueString()
    if rawDuration != "" {
        additionalData := complianceRule.GetAdditionalData()
        if additionalData == nil {
            additionalData = make(map[string]any)
        }
        additionalData["durationBeforeDeploymentStart"] = rawDuration
        complianceRule.SetAdditionalData(additionalData)
    }
}

// For durationBetweenOffers
if !rolloutData.DurationBetweenOffers.IsNull() {
    rawDuration := rolloutData.DurationBetweenOffers.ValueString()
    if rawDuration != "" {
        additionalData["durationBetweenOffers"] = rawDuration
    }
}
```

#### 2. Denormalize on Read (state.go)

Convert SDK-normalized values back to day-based format when READING from API:

```go
func denormalizeISODuration(duration string) string {
    weekPattern := regexp.MustCompile(`^P(\d+)W$`)
    if matches := weekPattern.FindStringSubmatch(duration); len(matches) == 2 {
        weeks, err := strconv.Atoi(matches[1])
        if err != nil {
            return duration
        }
        days := weeks * 7
        return "P" + strconv.Itoa(days) + "D"
    }
    return duration
}

// Usage when mapping from API response
if duration := contentApprovalRule.GetDurationBeforeDeploymentStart(); duration != nil {
    durationStr := denormalizeISODuration(duration.String())
    ruleModel.DurationBeforeDeploymentStart = types.StringValue(durationStr)
}
```

### Why This Works

1. **On CREATE/UPDATE**: We send `P7D` as raw string → API accepts and stores `P7D`
2. **API Response**: API returns `P7D` in JSON
3. **SDK Deserialization**: SDK normalizes to `P1W` internally
4. **Our Denormalization**: We convert `P1W` back to `P7D` before storing in Terraform state
5. **Result**: Terraform state has `P7D` matching user config → No diff!

---

## CREATE vs UPDATE Differences

### CREATE Operation

**Required Fields**:
- `@odata.type`: `#microsoft.graph.windowsUpdates.updatePolicy`
- `audience.id`: Reference to deployment audience
- `complianceChanges`: Array with ContentApproval object

**Optional Fields**:
- `complianceChangeRules`: Array of ContentApprovalRule objects
- `deploymentSettings`: Deployment configuration

**Example Payload**:
```json
{
    "@odata.type": "#microsoft.graph.windowsUpdates.updatePolicy",
    "audience": {
        "id": "audience-guid"
    },
    "complianceChanges": [
        {
            "@odata.type": "#microsoft.graph.windowsUpdates.contentApproval"
        }
    ],
    "complianceChangeRules": [
        {
            "@odata.type": "#microsoft.graph.windowsUpdates.contentApprovalRule",
            "contentFilter": {
                "@odata.type": "#microsoft.graph.windowsUpdates.driverUpdateFilter"
            },
            "durationBeforeDeploymentStart": "P7D"
        }
    ],
    "deploymentSettings": {
        "schedule": {
            "gradualRollout": {
                "@odata.type": "#microsoft.graph.windowsUpdates.rateDrivenRolloutSettings",
                "devicesPerOffer": 1000,
                "durationBetweenOffers": "P1D"
            }
        }
    }
}
```

### UPDATE Operation

**Allowed Fields**:
- `deploymentSettings`: ONLY field that can be updated

**Rejected Fields** (despite docs saying they're updatable):
- ❌ `audience`: Cannot be changed
- ❌ `complianceChanges`: Cannot be changed
- ❌ `complianceChangeRules`: **REJECTED BY API** (causes 400 error)

**Testing Evidence**:
Created 6 UPDATE test payloads:
1. `update_01_full_with_odata.json`: ❌ Failed (includes complianceChangeRules)
2. `update_02_no_root_odata.json`: ❌ Failed (includes complianceChangeRules)
3. `update_03_no_deployment_odata.json`: ❌ Failed (includes complianceChangeRules)
4. `update_04_minimal_no_odata.json`: ❌ Failed (includes complianceChangeRules)
5. `update_05_only_compliance_rules.json`: ❌ Failed (only complianceChangeRules)
6. `update_06_only_deployment_settings.json`: ✅ **SUCCESS** (only deploymentSettings)

**Example Valid UPDATE Payload**:
```json
{
    "deploymentSettings": {
        "schedule": {
            "gradualRollout": {
                "@odata.type": "#microsoft.graph.windowsUpdates.rateDrivenRolloutSettings",
                "devicesPerOffer": 2000,
                "durationBetweenOffers": "P2D"
            }
        }
    }
}
```

### Implementation Strategy

```go
func constructResource(ctx context.Context, data *WindowsUpdatesAutopatchUpdatePolicyResourceModel, isUpdate bool) (graphmodelswindowsupdates.UpdatePolicyable, error) {
    requestBody := graphmodelswindowsupdates.NewUpdatePolicy()
    
    odataType := "#microsoft.graph.windowsUpdates.updatePolicy"
    requestBody.SetOdataType(&odataType)
    
    // For UPDATE operations, don't send audience, complianceChanges, or complianceChangeRules
    // API testing shows that UPDATE only accepts deploymentSettings, despite docs saying otherwise
    if !isUpdate {
        // Set audience (CREATE only)
        // Set complianceChanges (CREATE only)
    }
    
    // complianceChangeRules can only be set during CREATE, not UPDATE
    if !isUpdate && !data.ComplianceChangeRules.IsNull() {
        // Set complianceChangeRules (CREATE only)
    }
    
    // deploymentSettings works for both CREATE and UPDATE
    if !data.DeploymentSettings.IsNull() {
        // Set deployment settings
    }
    
    return requestBody, nil
}
```

---

## Write-Only Fields

### `compliance_changes` Field

**Problem**: This field is write-only (request-only), not returned in API responses.

**Symptoms**:
- Field is required in CREATE request
- Field is NOT present in GET response
- Causes `ImportStateVerify` to fail with: `"compliance_changes": "true",` difference

**Solution**:

1. **Do NOT map from API response** (state.go):
```go
func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchUpdatePolicyResourceModel, remoteResource graphmodelswindowsupdates.UpdatePolicyable) {
    // ... other mappings ...
    
    // Note: compliance_changes is write-only and not returned by the API
    // The value from config/state is preserved automatically
    // DO NOT add: data.ComplianceChanges = ...
}
```

2. **Ignore in ImportStateVerify** (resource_acceptance_test.go):
```go
ImportStateVerifyIgnore: []string{"timeouts", "compliance_changes"},
```

**Why This Works**:
- Terraform automatically preserves the config value in state when the field isn't set during Read
- The value from the user's configuration is maintained across refresh operations

---

## Nested Object Handling

### Problem: Unknown Values for Optional Nested Objects

When nested objects are optional and not provided in config, they need to handle "unknown" values during plan phase.

**Original Implementation** (WRONG):
```go
type WindowsUpdatesAutopatchUpdatePolicyResourceModel struct {
    DeploymentSettings    *DeploymentSettingsModel  // Pointer type
}
```

**Problem**: Caused `"Received unknown value"` errors for minimal configs without deployment_settings.

**Solution**: Use `types.Object` instead of struct pointers:
```go
type WindowsUpdatesAutopatchUpdatePolicyResourceModel struct {
    DeploymentSettings    types.Object  // Can handle Unknown, Null, and Known values
}
```

**Extraction** (construct.go):
```go
if !data.DeploymentSettings.IsNull() && !data.DeploymentSettings.IsUnknown() {
    var deploymentSettingsData DeploymentSettingsModel
    diags := data.DeploymentSettings.As(ctx, &deploymentSettingsData, basetypes.ObjectAsOptions{})
    if diags.HasError() {
        return nil, fmt.Errorf("failed to extract deployment_settings data: %s", diags.Errors()[0].Detail())
    }
    // ... use deploymentSettingsData
}
```

**State Mapping** (state.go):
```go
if settings := remoteResource.GetDeploymentSettings(); settings != nil {
    deploymentSettingsData := DeploymentSettingsModel{}
    // ... populate deploymentSettingsData
    deploymentSettingsObj, diags := types.ObjectValueFrom(ctx, DeploymentSettingsAttrTypes, deploymentSettingsData)
    if !diags.HasError() {
        data.DeploymentSettings = deploymentSettingsObj
    } else {
        data.DeploymentSettings = types.ObjectNull(DeploymentSettingsAttrTypes)
    }
} else {
    data.DeploymentSettings = types.ObjectNull(DeploymentSettingsAttrTypes)
}
```

---

## Set vs List for compliance_change_rules

### Evolution of Implementation

**Attempt 1**: `schema.SetNestedAttribute` + `types.Set`
- Problem: Set correlation errors due to computed timestamp fields with UnknownVal in plan vs concrete values in state

**Attempt 2**: `schema.ListNestedAttribute` + `types.List`
- Problem: Fixed set correlation but duration normalization (P7D vs P1W) still caused diffs

**Final Solution**: Back to `schema.SetNestedAttribute` + `types.Set` with denormalization
- Denormalization ensures P1W → P7D conversion
- Set elements now match exactly between plan and state
- All tests pass

---

## Computed Timestamp Fields

### Fields in compliance_change_rules

Three timestamp fields are computed (read-only):
- `created_date_time`
- `last_evaluated_date_time`
- `last_modified_date_time`

**Implementation**:
```go
"created_date_time": schema.StringAttribute{
    Computed:            true,
    MarkdownDescription: "The date and time when the rule was created. Read-only.",
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.UseStateForUnknown(),
    },
},
```

**Why UseStateForUnknown()**:
- During CREATE, these fields are Unknown in plan
- After CREATE, they have concrete values from API
- `UseStateForUnknown()` tells Terraform to use state value when plan value is Unknown
- Prevents unnecessary diffs on refresh

---

## Common Pitfalls and Gotchas

### 1. Don't Trust the Documentation
- Microsoft's API documentation is incomplete and sometimes incorrect
- Always validate with actual API testing
- Document discrepancies for future reference

### 2. SDK Getters vs additionalData
- **For SEND (construct.go)**: Use `additionalData` for duration fields to avoid normalization
- **For READ (state.go)**: Use SDK getters (they work correctly), then denormalize
- **Don't use additionalData fallbacks**: SDK getters are reliable, fallbacks are unnecessary

### 3. UPDATE Operation Limitations
- Despite docs saying `complianceChangeRules` is updatable, it's NOT
- Only `deploymentSettings` can be updated
- Mark immutable fields with `RequiresReplace()` to trigger destroy+create

### 4. Filter Type Validation
- Only `driverUpdateFilter` works
- `windowsUpdateFilter` fails 100% of the time
- Add schema validation to prevent user errors

### 5. Duration Validation
- Maximum is P30D for both duration fields
- Add regex validation to catch invalid values early
- Provide clear error messages

### 6. Write-Only Fields
- `compliance_changes` is write-only
- Don't attempt to read it from API response
- Add to `ImportStateVerifyIgnore`

---

## Testing Commands

### Run Acceptance Tests
```bash
cd /Users/dafyddwatkins/GitHub/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_update_policy
export M365_CLOUD="public"
export M365_TENANT_ID="your-tenant-id"
export M365_AUTH_METHOD="client_secret"
export M365_CLIENT_ID="your-client-id"
export M365_CLIENT_SECRET="your-client-secret"
TF_ACC=1 go test -v -run TestAccResourceWindowsUpdatesUpdatePolicy -timeout 120m
```

### Run PowerShell Validation Tests
```powershell
# Test duration values for durationBeforeDeploymentStart
./Test-DurationValues.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx"

# Test duration values for durationBetweenOffers
./Test-DurationBetweenOffers.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx"

# Test device count values
./Test-DevicesPerOffer.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx"

# Test UPDATE payloads
./Run-AllUpdatePayloadTests.ps1 -TenantId "xxx" -ClientId "xxx" -ClientSecret "xxx"
```

---

## Summary of Validation Rules

| Field | Valid Values | Documentation | Reality | Mutable |
|-------|-------------|---------------|---------|---------|
| `filter_type` | `driverUpdateFilter`, `windowsUpdateFilter` | Both valid | Only `driverUpdateFilter` | No (RequiresReplace) |
| `duration_before_deployment_start` | ISO 8601 duration | No range specified | P1D to P30D only | No (RequiresReplace) |
| `duration_between_offers` | ISO 8601 duration | No range specified | P1D to P30D only | Yes |
| `devices_per_offer` | Integer | No range specified | 1 to 1,000,000+ | Yes |
| `start_date_time` | ISO 8601 datetime | Optional | Optional, not computed | Yes |
| `compliance_changes` | Boolean | Required | Write-only (not in GET) | No |
| `compliance_change_rules` | Array | Updatable | CREATE only, not updatable | No (RequiresReplace) |
| `audience` | Object with id | Required | CREATE only, not updatable | No (RequiresReplace) |

---

## Lessons Learned

1. **Beta APIs are unpredictable**: Microsoft Graph Beta endpoints often have undocumented restrictions
2. **Test everything**: Don't assume documentation is correct
3. **PowerShell testing is invaluable**: Direct API testing reveals ground truth
4. **SDK abstractions can cause issues**: Be aware of type conversions and normalizations
5. **Document everything**: Future developers will thank you

---

## Future Considerations

### If Microsoft Fixes These Issues

If Microsoft updates the API to:
1. **Support `windowsUpdateFilter`**: Remove `OneOf` validator, allow both filter types
2. **Support updating `complianceChangeRules`**: Remove `RequiresReplace()`, update construct logic
3. **Extend duration limits**: Update regex validators to new ranges
4. **Fix SDK normalization**: Remove `additionalData` workarounds and denormalization logic

### Monitoring for Changes

Periodically re-run the PowerShell test scripts to detect if API behavior changes:
- `Test-DurationValues.ps1`: Check if duration limits change
- `Test-DevicesPerOffer.ps1`: Check if device limits change
- `Run-AllUpdatePayloadTests.ps1`: Check if UPDATE restrictions are lifted

---

## References

- [Create updatePolicy API Docs](https://learn.microsoft.com/en-us/graph/api/adminwindowsupdates-post-updatepolicies?view=graph-rest-beta)
- [Update updatePolicy API Docs](https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatepolicy-update?view=graph-rest-beta)
- [updatePolicy Resource Docs](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-updatepolicy?view=graph-rest-beta)

---

**Document Version**: 1.0  
**Last Updated**: 2026-03-18  
**Tested Against**: Microsoft Graph API Beta
