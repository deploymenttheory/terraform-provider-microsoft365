package crud

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
)

// newImportStateResponse returns a response whose state has a single "id" string attribute
// set to null, which is how the framework hands state to ImportState.
func newImportStateResponse(t *testing.T) *resource.ImportStateResponse {
	t.Helper()

	s := schema.Schema{Attributes: map[string]schema.Attribute{"id": schema.StringAttribute{Computed: true}}}
	objectType := s.Type().TerraformType(context.Background())
	raw := tftypes.NewValue(objectType, map[string]tftypes.Value{"id": tftypes.NewValue(tftypes.String, nil)})

	return &resource.ImportStateResponse{State: tfsdk.State{Schema: s, Raw: raw}}
}

func importedID(t *testing.T, resp *resource.ImportStateResponse) types.String {
	t.Helper()

	var id types.String
	diags := resp.State.GetAttribute(context.Background(), path.Root("id"), &id)
	require.False(t, diags.HasError(), "reading id from state: %v", diags)

	return id
}

func TestImportStatePassthroughValidatedID(t *testing.T) {
	prefixedGUID := ImportIDFormat{Pattern: regexp.MustCompile(constants.PrefixedGuidRegex), Description: "a prefixed GUID such as A_00000000-0000-0000-0000-000000000000"}

	tests := []struct {
		name   string
		id     string
		format ImportIDFormat
		valid  bool
	}{
		{name: "valid GUID", id: "d557c813-b8e5-4efc-b00e-9c0bd5fd10df", format: ImportIDFormatGUID, valid: true},
		{name: "valid upper-case GUID", id: "D557C813-B8E5-4EFC-B00E-9C0BD5FD10DF", format: ImportIDFormatGUID, valid: true},
		{name: "valid prefixed GUID", id: "A_d557c813-b8e5-4efc-b00e-9c0bd5fd10df", format: prefixedGUID, valid: true},
		{name: "prefixed GUID against GUID format", id: "A_d557c813-b8e5-4efc-b00e-9c0bd5fd10df", format: ImportIDFormatGUID, valid: false},
		{name: "empty", id: "", format: ImportIDFormatGUID, valid: false},
		{name: "query string", id: "d557c813-b8e5-4efc-b00e-9c0bd5fd10df?$select=id", format: ImportIDFormatGUID, valid: false},
		{name: "ampersand", id: "d557c813-b8e5-4efc-b00e-9c0bd5fd10df&x=1", format: ImportIDFormatGUID, valid: false},
		{name: "path traversal", id: "../d557c813-b8e5-4efc-b00e-9c0bd5fd10df", format: ImportIDFormatGUID, valid: false},
		{name: "slash", id: "d557c813-b8e5-4efc-b00e-9c0bd5fd10df/settings", format: ImportIDFormatGUID, valid: false},
		{name: "single quote", id: "d557c813-b8e5-4efc-b00e-9c0bd5fd10df')", format: ImportIDFormatGUID, valid: false},
		{name: "curly brace", id: "{d557c813-b8e5-4efc-b00e-9c0bd5fd10df}", format: ImportIDFormatGUID, valid: false},
		{name: "percent", id: "d557c813-b8e5-4efc-b00e-9c0bd5fd10d%66", format: ImportIDFormatGUID, valid: false},
		{name: "trailing newline", id: "d557c813-b8e5-4efc-b00e-9c0bd5fd10df\n", format: ImportIDFormatGUID, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := newImportStateResponse(t)
			req := resource.ImportStateRequest{ID: tt.id}

			ImportStatePassthroughValidatedID(context.Background(), req, resp, tt.format)

			if tt.valid {
				require.False(t, resp.Diagnostics.HasError(), "unexpected diagnostics: %v", resp.Diagnostics)
				assert.Equal(t, types.StringValue(tt.id), importedID(t, resp))

				return
			}

			require.True(t, resp.Diagnostics.HasError())
			require.Len(t, resp.Diagnostics.Errors(), 1)
			assert.Equal(t, "Invalid Import ID", resp.Diagnostics.Errors()[0].Summary())
			assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.format.Description)
			assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "import ID")
			assert.True(t, importedID(t, resp).IsNull(), "state must not be set for an invalid ID")
		})
	}
}

// TestImportStatePassthroughValidatedID_MatchesPassthrough checks that a valid ID leaves the
// response exactly as resource.ImportStatePassthroughID would.
func TestImportStatePassthroughValidatedID_MatchesPassthrough(t *testing.T) {
	id := "d557c813-b8e5-4efc-b00e-9c0bd5fd10df"

	want := newImportStateResponse(t)
	resource.ImportStatePassthroughID(context.Background(), path.Root("id"), resource.ImportStateRequest{ID: id}, want)

	got := newImportStateResponse(t)
	ImportStatePassthroughValidatedID(context.Background(), resource.ImportStateRequest{ID: id}, got, ImportIDFormatGUID)

	assert.Equal(t, want.Diagnostics, got.Diagnostics)
	assert.True(t, want.State.Raw.Equal(got.State.Raw))
}

func TestValidateImportIDPart(t *testing.T) {
	t.Run("valid part", func(t *testing.T) {
		var diags diag.Diagnostics

		ok := ValidateImportIDPart(&diags, "agent_identity_blueprint_id", "11111111-1111-1111-1111-111111111111", ImportIDFormatGUID)

		assert.True(t, ok)
		assert.False(t, diags.HasError())
	})

	t.Run("invalid part names the part and format", func(t *testing.T) {
		var diags diag.Diagnostics

		ok := ValidateImportIDPart(&diags, "agent_identity_blueprint_id", "abc?$select=x", ImportIDFormatGUID)

		assert.False(t, ok)
		require.Len(t, diags.Errors(), 1)
		assert.Equal(t, "Invalid Import ID", diags.Errors()[0].Summary())
		assert.Equal(t, `The agent_identity_blueprint_id "abc?$select=x" is not valid. It must be a GUID such as 00000000-0000-0000-0000-000000000000.`, diags.Errors()[0].Detail())
	})

	t.Run("enrollment notification configuration format", func(t *testing.T) {
		format := ImportIDFormatEnrollmentNotificationConfiguration
		var diags diag.Diagnostics

		assert.True(t, ValidateImportIDPart(&diags, "import ID", "d557c813-b8e5-4efc-b00e-9c0bd5fd10df_EnrollmentNotificationsConfiguration", format))
		assert.False(t, ValidateImportIDPart(&diags, "import ID", "d557c813-b8e5-4efc-b00e-9c0bd5fd10df", format))
		assert.False(t, ValidateImportIDPart(&diags, "import ID", "d557c813-b8e5-4efc-b00e-9c0bd5fd10df_EnrollmentNotificationsConfiguration/x", format))
	})
}
