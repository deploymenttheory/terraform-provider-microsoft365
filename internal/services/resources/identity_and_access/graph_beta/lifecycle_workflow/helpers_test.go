package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"
	"testing"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newGraphError(status int, message string) error {
	mainError := odataerrors.NewMainError()
	code := "BAD_REQUEST"
	mainError.SetCode(&code)
	mainError.SetMessage(&message)

	err := odataerrors.NewODataError()
	err.SetStatusCode(status)
	err.SetErrorEscaped(mainError)
	return err
}

func TestUnitWriteWithTaskArgumentRetryRetriesUntilAccepted(t *testing.T) {
	calls := 0
	err := writeWithTaskArgumentRetry(context.Background(), func(ctx context.Context) error {
		calls++
		if calls == 1 {
			return newGraphError(400, "Task arguments missing or invalid")
		}
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, 2, calls)
}

func TestUnitWriteWithTaskArgumentRetryReturnsOtherErrorsImmediately(t *testing.T) {
	calls := 0
	err := writeWithTaskArgumentRetry(context.Background(), func(ctx context.Context) error {
		calls++
		return newGraphError(400, "Missing or invalid parameters (Parameter 'Category')")
	})

	require.Error(t, err)
	assert.Equal(t, 1, calls)
	assert.Contains(t, err.Error(), "Category")
}
