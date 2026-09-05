package errors

import (
	"context"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/require"
)

func TestUnitGraphReadErrorOptions(t *testing.T) {
	for _, strict := range []bool{false, true} {
		state := tfsdk.State{Schema: schema.Schema{Attributes: map[string]schema.Attribute{"id": schema.StringAttribute{Computed: true}}}}
		require.False(t, state.Set(context.Background(), struct {
			ID string `tfsdk:"id"`
		}{ID: "existing"}).HasError())
		resp := resource.ReadResponse{State: state}
		err := abstractions.NewApiError()
		err.SetStatusCode(400)
		if strict {
			HandleKiotaGraphErrorWithOptions(context.Background(), err, &resp, constants.TfOperationRead, nil, GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		} else {
			HandleKiotaGraphError(context.Background(), err, &resp, constants.TfOperationRead, nil)
		}
		require.Equal(t, strict, resp.Diagnostics.HasError())
		require.Equal(t, !strict, resp.State.Raw.IsNull())
	}
}
