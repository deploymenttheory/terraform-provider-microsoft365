package crud

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// mockStateContainer implements StateContainer for testing
type mockStateContainer struct {
	state tfsdk.State
}

func (m *mockStateContainer) GetState() tfsdk.State {
	return m.state
}

func (m *mockStateContainer) SetState(state tfsdk.State) {
	m.state = state
}

func TestReadWithRetry_SuccessOnFirstAttempt(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	callCount := 0
	mockReadFunc := func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
		callCount++
	}

	readReq := resource.ReadRequest{}
	stateContainer := &mockStateContainer{}
	opts := DefaultReadWithRetryOptions()

	err := ReadWithRetry(ctx, mockReadFunc, readReq, stateContainer, opts)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got: %d", callCount)
	}
}

func TestReadWithRetry_SuccessAfterRetries(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	callCount := 0
	mockReadFunc := func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
		callCount++
		if callCount < 3 {
			resp.Diagnostics.AddError("test error", "test error message")
		}
	}

	readReq := resource.ReadRequest{}
	stateContainer := &mockStateContainer{}
	opts := ReadWithRetryOptions{
		MaxRetries:    5,
		RetryInterval: 100 * time.Millisecond,
		Operation:     "Test",
	}

	err := ReadWithRetry(ctx, mockReadFunc, readReq, stateContainer, opts)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if callCount != 3 {
		t.Errorf("Expected 3 calls, got: %d", callCount)
	}
}

func TestReadWithRetry_AllRetriesExhausted(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	callCount := 0
	mockReadFunc := func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
		callCount++
		resp.Diagnostics.AddError("persistent error", "this error persists")
	}

	readReq := resource.ReadRequest{}
	stateContainer := &mockStateContainer{}
	opts := ReadWithRetryOptions{
		MaxRetries:    3,
		RetryInterval: 100 * time.Millisecond,
		Operation:     "Test",
	}

	err := ReadWithRetry(ctx, mockReadFunc, readReq, stateContainer, opts)

	if err == nil {
		t.Error("Expected error but got none")
	}
	if callCount != 4 { // MaxRetries + 1
		t.Errorf("Expected 4 calls, got: %d", callCount)
	}
	if !(contains(err.Error(), "failed to read resource state for") && contains(err.Error(), "after 4 attempts")) {
		t.Errorf("Error message doesn't contain expected text: %v", err)
	}
}

func TestReadWithRetry_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	callCount := 0
	mockReadFunc := func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
		callCount++
		resp.Diagnostics.AddError("error", "error message")
	}

	readReq := resource.ReadRequest{}
	stateContainer := &mockStateContainer{}
	opts := ReadWithRetryOptions{
		MaxRetries:    10,
		RetryInterval: 100 * time.Millisecond,
		Operation:     "Test",
	}

	err := ReadWithRetry(ctx, mockReadFunc, readReq, stateContainer, opts)

	if err == nil {
		t.Error("Expected timeout error but got none")
	}
	if callCount >= 10 {
		t.Errorf("Expected fewer calls due to timeout, got: %d", callCount)
	}
}

func TestReadWithRetry_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	callCount := 0
	mockReadFunc := func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
		callCount++
		if callCount == 2 {
			cancel()
		}
		resp.Diagnostics.AddError("error", "error message")
	}

	readReq := resource.ReadRequest{}
	stateContainer := &mockStateContainer{}
	opts := ReadWithRetryOptions{
		MaxRetries:    5,
		RetryInterval: 100 * time.Millisecond,
		Operation:     "Test",
	}

	err := ReadWithRetry(ctx, mockReadFunc, readReq, stateContainer, opts)

	if err == nil {
		t.Error("Expected cancellation error but got none")
	}
	if !contains(err.Error(), "context cancelled") {
		t.Errorf("Expected context cancellation error, got: %v", err)
	}
}

func TestReadWithRetry_NoDeadline(t *testing.T) {
	ctx := context.Background()

	mockReadFunc := func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {}

	readReq := resource.ReadRequest{}
	stateContainer := &mockStateContainer{}
	opts := DefaultReadWithRetryOptions()

	err := ReadWithRetry(ctx, mockReadFunc, readReq, stateContainer, opts)

	if err == nil {
		t.Error("Expected error for context without deadline")
	}
	if !contains(err.Error(), "context must have a deadline") {
		t.Errorf("Expected deadline error, got: %v", err)
	}
}

func TestReadWithRetry_InsufficientTime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	mockReadFunc := func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {}

	readReq := resource.ReadRequest{}
	stateContainer := &mockStateContainer{}
	opts := ReadWithRetryOptions{
		MaxRetries:    10,
		RetryInterval: 2 * time.Second,
		Operation:     "Test",
	}

	err := ReadWithRetry(ctx, mockReadFunc, readReq, stateContainer, opts)

	if err == nil {
		t.Error("Expected insufficient time error")
	}
	if !contains(err.Error(), "insufficient time remaining") {
		t.Errorf("Expected insufficient time error, got: %v", err)
	}
}

func TestReadWithRetry_DefaultOptions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mockReadFunc := func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {}

	readReq := resource.ReadRequest{}
	stateContainer := &mockStateContainer{}
	opts := ReadWithRetryOptions{}

	err := ReadWithRetry(ctx, mockReadFunc, readReq, stateContainer, opts)

	if err != nil {
		t.Errorf("Expected no error with default options, got: %v", err)
	}
}

func TestReadWithRetry_StateContainerUpdated(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	expectedState := tfsdk.State{}
	mockReadFunc := func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
		resp.State = expectedState
	}

	readReq := resource.ReadRequest{}
	stateContainer := &mockStateContainer{}
	opts := DefaultReadWithRetryOptions()

	err := ReadWithRetry(ctx, mockReadFunc, readReq, stateContainer, opts)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if stateContainer.state != expectedState {
		t.Error("State container was not updated with the expected state")
	}
}

func TestCreateResponseContainer(t *testing.T) {
	resp := &resource.CreateResponse{}
	container := &CreateResponseContainer{CreateResponse: resp}

	state := tfsdk.State{}
	container.SetState(state)

	if container.GetState() != state {
		t.Error("CreateResponseContainer state management failed")
	}
	if resp.State != state {
		t.Error("CreateResponseContainer didn't update underlying response state")
	}
}

func TestUpdateResponseContainer(t *testing.T) {
	resp := &resource.UpdateResponse{}
	container := &UpdateResponseContainer{UpdateResponse: resp}

	state := tfsdk.State{}
	container.SetState(state)

	if container.GetState() != state {
		t.Error("UpdateResponseContainer state management failed")
	}
	if resp.State != state {
		t.Error("UpdateResponseContainer didn't update underlying response state")
	}
}

func TestDefaultReadWithRetryOptions(t *testing.T) {
	opts := DefaultReadWithRetryOptions()

	if opts.MaxRetries != 30 {
		t.Errorf("Expected MaxRetries to be 30, got: %d", opts.MaxRetries)
	}
	if opts.RetryInterval != 2*time.Second {
		t.Errorf("Expected RetryInterval to be 2s, got: %v", opts.RetryInterval)
	}
	if opts.Operation != "Operation" {
		t.Errorf("Expected Operation to be 'Operation', got: %s", opts.Operation)
	}
}

func TestReadWithRetry_MaxRetriesAdjustedForTime(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	callCount := 0
	mockReadFunc := func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
		callCount++
		resp.Diagnostics.AddError("error", "error message")
	}

	readReq := resource.ReadRequest{}
	stateContainer := &mockStateContainer{}
	opts := ReadWithRetryOptions{
		MaxRetries:    100,
		RetryInterval: 500 * time.Millisecond,
		Operation:     "Test",
	}

	err := ReadWithRetry(ctx, mockReadFunc, readReq, stateContainer, opts)

	if err == nil {
		t.Error("Expected error but got none")
	}
	if callCount > 10 {
		t.Errorf("Expected limited calls due to time constraint, got: %d", callCount)
	}
}

func TestExtractResourceID(t *testing.T) {
	ctx := context.Background()

	// Test with empty state
	emptyState := tfsdk.State{}
	id := extractResourceID(ctx, emptyState)
	if id != "unknown" {
		t.Errorf("Expected 'unknown' for empty state, got: %s", id)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || s[0:len(substr)] == substr || contains(s[1:], substr))
}
