package errors

import (
	"context"
	"fmt"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/require"
	"io"
	"net/url"
	"testing"
)

func TestUnitGraphErrorOptions_ReadBadRequest(t *testing.T) {
	for _, preserve := range []bool{false, true} {
		state := tfsdk.State{Schema: schema.Schema{Attributes: map[string]schema.Attribute{"id": schema.StringAttribute{Computed: true}}}}
		require.False(t, state.Set(context.Background(), struct {
			ID string `tfsdk:"id"`
		}{ID: "existing"}).HasError())
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

func TestUnitGraphErrorOptions_ReadTransportError(t *testing.T) {
	for _, wrapped := range []bool{false, true} {
		for _, preserve := range []bool{false, true} {
			t.Run(fmt.Sprintf("wrapped=%t/preserve=%t", wrapped, preserve), func(t *testing.T) {
				ctx := context.Background()
				state := tfsdk.State{Schema: schema.Schema{Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{Computed: true},
				}}}
				require.False(t, state.Set(ctx, struct {
					ID string `tfsdk:"id"`
				}{ID: "existing"}).HasError())
				resp := resource.ReadResponse{State: state}
				var err error = &url.Error{Op: "Get", URL: "https://graph.microsoft.com/beta/test", Err: io.EOF}
				if wrapped {
					err = fmt.Errorf("read request failed: %w", err)
				}
				if preserve {
					HandleKiotaGraphErrorWithOptions(ctx, err, &resp, constants.TfOperationRead, nil,
						GraphErrorOptions{PreserveStateOnReadBadRequest: true})
				} else {
					HandleKiotaGraphError(ctx, err, &resp, constants.TfOperationRead, nil)
				}
				require.True(t, resp.Diagnostics.HasError())
				require.True(t, state.Raw.Equal(resp.State.Raw), "transport errors must preserve the existing resource")
			})
		}
	}
}
