package crud

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReadWithRetryOptions_DefaultsMatchDefaultReadWithRetryOptions(t *testing.T) {
	defaults := DefaultReadWithRetryOptions()

	opts := NewReadWithRetryOptions("create", "microsoft365_test_resource")

	assert.Equal(t, defaults.MaxRetries, opts.MaxRetries)
	assert.Equal(t, defaults.RetryInterval, opts.RetryInterval)
	assert.Equal(t, "create", opts.Operation)
	assert.Equal(t, "microsoft365_test_resource", opts.ResourceTypeName)
	assert.Nil(t, opts.ConsistencyPredicate)
}

func TestNewReadWithRetryOptions_AppliesOptions(t *testing.T) {
	predicateCalls := 0
	predicate := func(ctx context.Context, state tfsdk.State) bool {
		predicateCalls++
		return true
	}

	tests := []struct {
		name             string
		options          []ReadWithRetryOption
		wantMaxRetries   int
		wantInterval     time.Duration
		wantPredicateSet bool
	}{
		{name: "predicate only", options: []ReadWithRetryOption{ReadWithConsistencyPredicate(predicate)}, wantMaxRetries: 30, wantInterval: 2 * time.Second, wantPredicateSet: true},
		{name: "predicate and max retries", options: []ReadWithRetryOption{ReadWithConsistencyPredicate(predicate), ReadWithMaxRetries(1000)}, wantMaxRetries: 1000, wantInterval: 2 * time.Second, wantPredicateSet: true},
		{name: "predicate, interval and max retries", options: []ReadWithRetryOption{ReadWithConsistencyPredicate(predicate), ReadWithRetryInterval(10 * time.Millisecond), ReadWithMaxRetries(60)}, wantMaxRetries: 60, wantInterval: 10 * time.Millisecond, wantPredicateSet: true},
		{name: "max retries only", options: []ReadWithRetryOption{ReadWithMaxRetries(5)}, wantMaxRetries: 5, wantInterval: 2 * time.Second},
		{name: "nil predicate keeps fast exit", options: []ReadWithRetryOption{ReadWithConsistencyPredicate(nil)}, wantMaxRetries: 30, wantInterval: 2 * time.Second},
		{name: "nil option is ignored", options: []ReadWithRetryOption{nil, ReadWithMaxRetries(7)}, wantMaxRetries: 7, wantInterval: 2 * time.Second},
		{name: "later option wins", options: []ReadWithRetryOption{ReadWithMaxRetries(7), ReadWithMaxRetries(9)}, wantMaxRetries: 9, wantInterval: 2 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewReadWithRetryOptions("update", "microsoft365_test_resource", tt.options...)

			assert.Equal(t, "update", opts.Operation)
			assert.Equal(t, "microsoft365_test_resource", opts.ResourceTypeName)
			assert.Equal(t, tt.wantMaxRetries, opts.MaxRetries)
			assert.Equal(t, tt.wantInterval, opts.RetryInterval)
			if !tt.wantPredicateSet {
				assert.Nil(t, opts.ConsistencyPredicate)
				return
			}
			require.NotNil(t, opts.ConsistencyPredicate)
			before := predicateCalls
			assert.True(t, opts.ConsistencyPredicate(context.Background(), tfsdk.State{}))
			assert.Equal(t, before+1, predicateCalls)
		})
	}
}

func TestNewReadWithRetryOptions_UsedByReadWithRetry(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reads := 0
	readFunc := func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) { reads++ }
	consistentAfter := 3
	predicate := func(ctx context.Context, state tfsdk.State) bool { return reads >= consistentAfter }

	opts := NewReadWithRetryOptions("create", "microsoft365_test_resource", ReadWithConsistencyPredicate(predicate), ReadWithRetryInterval(5*time.Millisecond), ReadWithMaxRetries(10))

	err := ReadWithRetry(ctx, readFunc, resource.ReadRequest{}, &mockStateContainer{}, opts)

	require.NoError(t, err)
	assert.Equal(t, consistentAfter, reads)
}
