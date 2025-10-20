# ADR-001: Design Principles for Terraform Actions

## Status

Accepted

## Date

2025-01-20

## Context

The Microsoft 365 Terraform provider is implementing Terraform Actions to handle non-CRUD operations like device management actions (eSIM activation, device wipe, etc.). These actions differ from traditional resources in that they perform one-time operations rather than managing state.

As we develop multiple actions, we need consistent design principles to ensure:
- Actions are predictable and reliable
- Error handling is appropriate and consistent
- User experience is professional and actionable
- Actions follow Terraform's design philosophy
- Actions can be effectively tested and debugged

The current eSIM activation action implementation revealed several design questions around retry logic, error handling, concurrency control, and user messaging that need principled answers.

## Decision Drivers

* Terraform Actions design philosophy emphasizing specificity and simplicity
* Microsoft Graph API characteristics and error patterns
* Enterprise user expectations for professional tooling
* CLI debuggability and operational transparency
* Consistency across multiple device management actions
* Separation of concerns between validation, execution, and error handling

## Considered Options

* **Complex orchestration approach**: Application-level retry logic, sophisticated concurrency control, stateful operations
* **Simple atomic approach**: Single-attempt operations with comprehensive upfront validation
* **Mixed approach**: Selective retry for specific scenarios with fallback to simple execution

## Decision

Chosen option: "Simple atomic approach", because it aligns with Terraform Actions philosophy and provides better reliability through validation rather than recovery.

## Rationale

Terraform Actions are designed to be simple, specific building blocks rather than general-purpose orchestration systems. The official design decisions emphasize:
- Specificity over generality
- Simplicity and predictability
- Stateless operations
- CLI-first design

Microsoft Graph API operations like eSIM activation are typically atomic - they either succeed immediately or fail for deterministic reasons (invalid device state, wrong permissions, unsupported device). Network-level retries should be handled by the HTTP transport layer, not application logic.

## Consequences

### Positive

* Actions are predictable and deterministic
* Easier to test and debug via CLI
* Reduced complexity in error handling
* Better separation of concerns
* Consistent with Terraform's design philosophy
* Professional, actionable user messaging
* Faster feedback loops through comprehensive validation

### Negative

* No application-level retry for edge cases
* Requires more comprehensive upfront validation
* May need manual retry for rare transient failures

### Neutral

* Shifts error handling emphasis from recovery to prevention
* Requires clear documentation of validation behavior

## Implementation

### Core Principles

1. **Specificity**: Each action performs one specific task extremely well
2. **Atomic Operations**: Single attempt per device/operation with clear success/failure
3. **Validation First**: Comprehensive upfront validation to catch issues before execution
4. **Professional Messaging**: Concise, actionable feedback without emojis or casual language
5. **CLI Debuggability**: All actions must be testable via `terraform action invoke`
6. **Stateless Design**: No internal state tracking or complex orchestration

### Error Handling Strategy

1. **Validation Phase**: Check device existence, permissions, configuration validity
2. **Execution Phase**: Single API call with clear success/failure
3. **Network Retries**: Handled by HTTP client and centralized error system
4. **User Feedback**: Use `resp.SendProgress` for status, `resp.Diagnostics` only for actionable warnings/errors

#### Code Example: Validation First Approach

```go
// GOOD: Comprehensive upfront validation
func (a *ActivateDeviceEsimAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
    var data ActivateDeviceEsimActionModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    // Basic configuration validation
    if len(data.ManagedDevices) == 0 && len(data.ComanagedDevices) == 0 {
        resp.Diagnostics.AddError(
            "No Devices Specified",
            "At least one of 'managed_devices' or 'comanaged_devices' must be provided.")
        return
    }

    // Online validation when enabled and client available
    if validateExists && a.client != nil {
        for _, device := range data.ManagedDevices {
            deviceID := device.DeviceID.ValueString()
            _, err := a.client.DeviceManagement().ManagedDevices().ByManagedDeviceId(deviceID).Get(ctx, nil)
            if err != nil {
                resp.Diagnostics.AddAttributeError(
                    path.Root("managed_devices"),
                    "Device Not Found",
                    fmt.Sprintf("Device %s does not exist or is not enrolled", deviceID))
            }
        }
    }
}
```

#### Code Example: Atomic Execution

```go
// GOOD: Single attempt, clear result
func (a *ActivateDeviceEsimAction) activateEsimManagedDevice(ctx context.Context, device ManagedDeviceActivateEsim) error {
    deviceID := device.DeviceID.ValueString()
    tflog.Debug(ctx, "Activating eSIM for managed device", map[string]any{"device_id": deviceID})

    requestBody := constructManagedDeviceRequest(ctx, device)
    err := a.client.
        DeviceManagement().
        ManagedDevices().
        ByManagedDeviceId(deviceID).
        ActivateDeviceEsim().
        Post(ctx, requestBody, nil)

    // Let centralized error handler deal with the error
    return err
}

// BAD: Application-level retry logic
func (a *ActivateDeviceEsimAction) activateEsimManagedDeviceWithRetry(ctx context.Context, device ManagedDeviceActivateEsim, maxRetries int) error {
    var lastErr error
    for attempt := 0; attempt <= maxRetries; attempt++ {
        if attempt > 0 {
            time.Sleep(time.Duration(attempt) * time.Second) // Complex state tracking
        }
        err := a.activateEsimManagedDevice(ctx, device)
        if err == nil {
            return nil
        }
        lastErr = err
        if isNonRetryableError(ctx, err) { // Complex decision logic
            break
        }
    }
    return lastErr
}
```

### Messaging Standards

#### Code Example: Professional Progress Messaging

```go
// GOOD: Professional, actionable messaging
resp.SendProgress(action.InvokeProgressEvent{
    Message: fmt.Sprintf("eSIM activation completed successfully for %d of %d devices.",
        successCount, totalDevices),
})

resp.SendProgress(action.InvokeProgressEvent{
    Message: fmt.Sprintf("Partial success: %d of %d devices activated. Failed devices: %v",
        successCount, totalDevices, failedDevices),
})

// BAD: Casual messaging with emojis
resp.SendProgress(action.InvokeProgressEvent{
    Message: fmt.Sprintf("✅ Success: All %d devices activated successfully\n"+
        "📱 All devices will have their eSIM enabled with the provided carrier profiles",
        successCount),
})
```

#### Code Example: Appropriate Use of Diagnostics

```go
// GOOD: Use diagnostics for actionable warnings/errors
if successCount > 0 && len(failedDevices) > 0 && !ignorePartialFailures {
    resp.Diagnostics.AddWarning(
        "Partial Success",
        fmt.Sprintf("eSIM activation partially completed. %d of %d devices succeeded. Failed devices: %s",
            successCount, totalDevices, strings.Join(failedDevices, ", ")))
}

if len(failedDevices) == totalDevices && !ignorePartialFailures {
    // Complete failure - use centralized error handler
    errors.HandleKiotaGraphError(ctx, lastError, resp, "Action", a.WritePermissions)
    return
}

// BAD: Using diagnostics for normal progress updates
resp.Diagnostics.AddWarning(
    "Operation Status",
    fmt.Sprintf("Processed %d devices so far", processedCount))
```

#### Code Example: Stateless Design

```go
// GOOD: Stateless action execution
func (a *ActivateDeviceEsimAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
    var data ActivateDeviceEsimActionModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

    // Process each device independently
    for _, device := range data.ManagedDevices {
        err := a.activateEsimManagedDevice(ctx, device)
        // Handle result immediately, no state tracking
        if err != nil {
            failedDevices = append(failedDevices, device.DeviceID.ValueString())
        } else {
            successCount++
        }
    }
    
    // Report final results
    a.reportResults(ctx, resp, successCount, failedDevices, ignorePartialFailures)
}

// BAD: Stateful orchestration
type ActionState struct {
    AttemptCounts    map[string]int
    RetryQueue      []Device
    BackoffTimers   map[string]*time.Timer
}

func (a *ActivateDeviceEsimAction) manageRetries(state *ActionState) {
    // Complex state management violates stateless principle
}
```

### Action Items

* [x] Document design principles in ADR
* [x] Ensure eSIM activation action follows simple atomic approach (retry logic not introduced)
* [x] Verify max_retries configuration option was not added
* [x] Implement validation for device existence with toggle option
* [x] Develop comprehensive action tests
* [x] Implement ignore_partial_failures configuration option
* [ ] Create action implementation guidelines for future actions
* [ ] Review existing actions for consistency with these principles
* [ ] Develop eSIM capability validation if Graph API supports it

### Timeline

- ADR completion: 2025-01-20
- eSIM action implementation: 2025-01-21
- Implementation validation: 2025-10-20
- Guidelines documentation: TBD

## Validation

Success will be measured by:
- All actions can be successfully invoked via CLI
- Error messages provide clear, actionable guidance
- Validation catches 90%+ of issues before execution
- Consistent user experience across all actions
- Reduced support requests due to unclear error messages

## References

* [Terraform Actions Design Decisions](https://danielmschmidt.de/posts/2025-09-26-terraform-actions-design-decisions/)
* [Terraform Action Patterns and Guidelines](https://danielmschmidt.de/posts/2025-09-26-terraform-action-patterns-and-guidelines/)
* [Writing a Terraform Action](https://danielmschmidt.de/posts/2025-09-26-writing-a-terraform-action/)
* [Ansible Terraform Provider Action Example](https://github.com/ansible/terraform-provider-ansible) - Reference implementation

## Notes

This ADR establishes the foundation for all future action implementations in the Microsoft 365 provider. Any deviation from these principles should be documented with clear justification.

The emphasis on validation over retry aligns with the "shift-left" philosophy of catching issues as early as possible in the development/deployment pipeline.

### Implementation Validation (2025-10-20)

The eSIM activation action implementation has been verified to accurately reflect all design principles outlined in this ADR:

**Confirmed Implementation Details:**
- ✅ **Specificity**: Single, well-defined task - eSIM activation only
- ✅ **Atomic Operations**: Single API call per device without application-level retry logic
- ✅ **Validation First**: Comprehensive upfront validation including device existence checks
- ✅ **Professional Messaging**: Clear, actionable progress messages without casual language
- ✅ **CLI Debuggability**: Full support via `terraform action invoke` with detailed logging
- ✅ **Stateless Design**: No state tracking, request-scoped variables only
- ✅ **Configuration Options**: 
  - `validate_device_exists`: Toggle for online device existence validation
  - `ignore_partial_failures`: Handle partial success scenarios gracefully
- ✅ **Error Handling**: Proper use of diagnostics for warnings/errors, centralized error handling via `HandleKiotaGraphError`
- ✅ **Comprehensive Testing**: Unit tests and mock responders present
- ✅ **Duplicate Detection**: Validates and warns about duplicate device IDs
- ✅ **No Retry Logic**: Application-level retry logic was not introduced; HTTP transport layer handles network retries

File structure follows established patterns with separation of concerns across `action.go`, `validate.go`, `invoke.go`, and `construct.go`.