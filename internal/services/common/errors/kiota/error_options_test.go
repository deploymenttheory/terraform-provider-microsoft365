package errors

import (
	"context"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnitGraphErrorOptions_ReadBadRequest(t *testing.T) {
	for _, preserve := range []bool{false, true} {
		state := tfsdk.State{Schema: schema.Schema{Attributes: map[string]schema.Attribute{"id": schema.StringAttribute{Computed: true}}}}
		require.False(t, state.Set(context.Background(), map[string]string{"id": "existing"}).HasError())
		resp := resource.ReadResponse{State: state}
		err := abstractions.NewApiError()
		err.SetStatusCode(400)
		if preserve {
			HandleKiotaGraphErrorWithOptions(context.Background(), err, &resp, constants.TfOperationRead, nil, GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		} else {
			HandleKiotaGraphError(context.Background(), err, &resp, constants.TfOperationRead, nil)
		}
		require.Equal(t, preserve, resp.Diagnostics.HasError())
		require.Equal(t, !preserve, resp.State.Raw.IsNull())
		if preserve {
			require.True(t, state.Raw.Equal(resp.State.Raw))
		}
	}
}
