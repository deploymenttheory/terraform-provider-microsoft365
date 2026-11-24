package crud

import (
	"context"
	stderrors "errors"
	"fmt"
	"testing"
	"time"

	betaodataerrors "github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// newGraphError returns an ODataError carrying an HTTP status code and an OData error code,
// as the Kiota adapter does for a failed Graph request. It uses the v1.0 SDK type because
// that is the ODataError errors.GraphError reads the error code from.
func newGraphError(statusCode int, code, message string) error {
	mainError := odataerrors.NewMainError()
	mainError.SetCode(&code)
	mainError.SetMessage(&message)
	odataError := odataerrors.NewODataError()
	odataError.SetStatusCode(statusCode)
	odataError.SetErrorEscaped(mainError)
	return odataError
}

// fastPreconditionRetryOptions keeps the default codes and retry count with a short interval.
func fastPreconditionRetryOptions() PreconditionRetryOptions {
	opts := DefaultPreconditionRetryOptions()
	opts.RetryInterval = time.Millisecond
	return opts
}

func TestDefaultPreconditionRetryOptions_MatchExistingWrappers(t *testing.T) {
	opts := DefaultPreconditionRetryOptions()

	assert.Equal(t, 3, opts.MaxRetries)
	assert.Equal(t, 5*time.Second, opts.RetryInterval)
	assert.Equal(t, []int{412}, opts.StatusCodes)
	assert.Equal(t, []string{"PreconditionFailed"}, opts.ErrorCodes)
}

func TestWriteWithPreconditionRetry_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := WriteWithPreconditionRetry(context.Background(), func(ctx context.Context) error {
		calls++
		return nil
	}, fastPreconditionRetryOptions())

	require.NoError(t, err)
	assert.Equal(t, 1, calls)
}

func TestWriteWithPreconditionRetry_SucceedsAfterPreconditionFailures(t *testing.T) {
	tests := []struct {
		name  string
		error error
	}{
		{name: "status and code match", error: newGraphError(412, "PreconditionFailed", "policy changed")},
		{name: "status matches only", error: newGraphError(412, "SomethingElse", "policy changed")},
		{name: "code matches only", error: newGraphError(409, "PreconditionFailed", "policy changed")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			err := WriteWithPreconditionRetry(context.Background(), func(ctx context.Context) error {
				calls++
				if calls < 3 {
					return tt.error
				}
				return nil
			}, fastPreconditionRetryOptions())

			require.NoError(t, err)
			assert.Equal(t, 3, calls)
		})
	}
}

func TestWriteWithPreconditionRetry_NonMatchingErrorReturnsAtOnce(t *testing.T) {
	tests := []struct {
		name  string
		error error
	}{
		{name: "bad request", error: newGraphError(400, "BadRequest", "invalid rule")},
		{name: "conflict", error: newGraphError(409, "Conflict", "duplicate rule")},
		{name: "plain error", error: stderrors.New("connection reset")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			err := WriteWithPreconditionRetry(context.Background(), func(ctx context.Context) error {
				calls++
				return tt.error
			}, fastPreconditionRetryOptions())

			assert.Same(t, tt.error, err)
			assert.Equal(t, 1, calls)
		})
	}
}

func TestWriteWithPreconditionRetry_ReturnsLastErrorWhenRetriesRunOut(t *testing.T) {
	var returned []error
	err := WriteWithPreconditionRetry(context.Background(), func(ctx context.Context) error {
		writeErr := newGraphError(412, "PreconditionFailed", fmt.Sprintf("attempt %d", len(returned)+1))
		returned = append(returned, writeErr)
		return writeErr
	}, fastPreconditionRetryOptions())

	require.Len(t, returned, 4, "the first attempt plus 3 retries")
	assert.Same(t, returned[3], err)
}

func TestWriteWithPreconditionRetry_ReturnsLastErrorAtDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	opts := DefaultPreconditionRetryOptions()
	opts.RetryInterval = time.Minute
	writeErr := newGraphError(412, "PreconditionFailed", "policy changed")
	calls := 0

	start := time.Now()
	err := WriteWithPreconditionRetry(ctx, func(ctx context.Context) error {
		calls++
		return writeErr
	}, opts)

	require.Error(t, err)
	assert.Less(t, time.Since(start), 10*time.Second)
	assert.Equal(t, 1, calls)
	assert.ErrorIs(t, err, writeErr)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, 412, errors.GraphError(context.Background(), err).StatusCode)
}

func TestWriteWithPreconditionRetry_ZeroOptionsUseDefaults(t *testing.T) {
	opts := PreconditionRetryOptions{RetryInterval: time.Millisecond}

	calls := 0
	err := WriteWithPreconditionRetry(context.Background(), func(ctx context.Context) error {
		calls++
		return newGraphError(412, "PreconditionFailed", "policy changed")
	}, opts)
	require.Error(t, err)
	assert.Equal(t, 4, calls, "default MaxRetries is 3")

	calls = 0
	err = WriteWithPreconditionRetry(context.Background(), func(ctx context.Context) error {
		calls++
		return newGraphError(400, "BadRequest", "invalid rule")
	}, opts)
	require.Error(t, err)
	assert.Equal(t, 1, calls, "default codes do not include 400")
}

func TestWriteWithPreconditionRetry_CustomCodes(t *testing.T) {
	opts := PreconditionRetryOptions{MaxRetries: 2, RetryInterval: time.Millisecond, StatusCodes: []int{409}, Operation: "update", ResourceTypeName: "microsoft365_test_rule", LogFields: map[string]any{"policy_id": "p1"}}

	calls := 0
	err := WriteWithPreconditionRetry(context.Background(), func(ctx context.Context) error {
		calls++
		return newGraphError(409, "Conflict", "policy busy")
	}, opts)
	require.Error(t, err)
	assert.Equal(t, 3, calls, "409 is retried MaxRetries times")

	calls = 0
	err = WriteWithPreconditionRetry(context.Background(), func(ctx context.Context) error {
		calls++
		return newGraphError(412, "PreconditionFailed", "policy changed")
	}, opts)
	require.Error(t, err)
	assert.Equal(t, 1, calls, "custom codes replace the defaults")
}

// The network policy rule packages map Graph errors to the beta SDK ODataError, from which
// errors.GraphError reads the status code but not the OData error code. For those callers only
// the 412 status triggers a retry, exactly as in their current wrappers.
func TestWriteWithPreconditionRetry_BetaODataErrorMatchesOnStatusOnly(t *testing.T) {
	newBetaGraphError := func(statusCode int, code string) error {
		message := "policy changed"
		mainError := betaodataerrors.NewMainError()
		mainError.SetCode(&code)
		mainError.SetMessage(&message)
		odataError := betaodataerrors.NewODataError()
		odataError.SetStatusCode(statusCode)
		odataError.SetErrorEscaped(mainError)
		return odataError
	}

	calls := 0
	err := WriteWithPreconditionRetry(context.Background(), func(ctx context.Context) error {
		calls++
		return newBetaGraphError(412, "PreconditionFailed")
	}, fastPreconditionRetryOptions())
	require.Error(t, err)
	assert.Equal(t, 4, calls, "412 is retried")

	calls = 0
	err = WriteWithPreconditionRetry(context.Background(), func(ctx context.Context) error {
		calls++
		return newBetaGraphError(409, "PreconditionFailed")
	}, fastPreconditionRetryOptions())
	require.Error(t, err)
	assert.Equal(t, 1, calls, "the beta OData error code is not visible to errors.GraphError")
}
