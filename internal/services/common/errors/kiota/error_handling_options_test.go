package errors_test

import (
	"context"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnitResourceNetworkCloudFirewall_ReadErrorOptions(t *testing.T) {
	for _, tc := range []struct {
		name                     string
		status                   int
		opt, removed, diagnostic bool
	}{
		{"legacy400", 400, false, true, false}, {"opt400", 400, true, false, true}, {"opt403", 403, true, false, true}, {"opt404", 404, true, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			resp := resource.ReadResponse{State: tfsdk.State{Schema: schema.Schema{Attributes: map[string]schema.Attribute{"id": schema.StringAttribute{Computed: true}}}}}
			resp.State.Raw = tftypes.NewValue(resp.State.Schema.Type().TerraformType(ctx), nil)
			require.False(t, resp.State.SetAttribute(ctx, path.Root("id"), "existing").HasError())
			err := abstractions.NewApiError()
			err.ResponseStatusCode = tc.status
			if tc.opt {
				errors.HandleKiotaGraphErrorWithOptions(ctx, err, &resp, constants.TfOperationRead, nil, errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
			} else {
				errors.HandleKiotaGraphError(ctx, err, &resp, constants.TfOperationRead, nil)
			}
			require.Equal(t, tc.removed, resp.State.Raw.IsNull())
			require.Equal(t, tc.diagnostic, resp.Diagnostics.HasError())
		})
	}
}
