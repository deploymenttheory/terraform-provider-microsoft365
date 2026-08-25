package attribute

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

// boolPresenceRequest builds a validator request against a real two attribute configuration,
// so that the dependent attribute lookup is genuinely exercised rather than skipped.
//
// A nil dependent or value means the corresponding attribute is null; an empty dependent
// string means unknown.
func boolPresenceRequest(dependent *string, dependentUnknown bool, value *bool) validator.BoolRequest {
	objectType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"intent": tftypes.String,
			"flag":   tftypes.Bool,
		},
	}

	dependentValue := tftypes.NewValue(tftypes.String, nil)
	switch {
	case dependentUnknown:
		dependentValue = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	case dependent != nil:
		dependentValue = tftypes.NewValue(tftypes.String, *dependent)
	}

	flagValue := tftypes.NewValue(tftypes.Bool, nil)
	configValue := types.BoolNull()
	if value != nil {
		flagValue = tftypes.NewValue(tftypes.Bool, *value)
		configValue = types.BoolValue(*value)
	}

	return validator.BoolRequest{
		Path:        path.Root("flag"),
		ConfigValue: configValue,
		Config: tfsdk.Config{
			Schema: fwresourceschema.Schema{
				Attributes: map[string]fwresourceschema.Attribute{
					"intent": fwresourceschema.StringAttribute{Optional: true},
					"flag":   fwresourceschema.BoolAttribute{Optional: true},
				},
			},
			Raw: tftypes.NewValue(objectType, map[string]tftypes.Value{
				"intent": dependentValue,
				"flag":   flagValue,
			}),
		},
	}
}

func stringPtr(s string) *string { return &s }
func boolPtr(b bool) *bool       { return &b }

func TestBoolSupportedOnlyWhenStringIn(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request     validator.BoolRequest
		expectError bool
	}{
		"set-with-supported-dependent-value": {
			request:     boolPresenceRequest(stringPtr("required"), false, boolPtr(true)),
			expectError: false,
		},
		"set-false-with-supported-dependent-value": {
			request:     boolPresenceRequest(stringPtr("required"), false, boolPtr(false)),
			expectError: false,
		},
		"set-true-with-unsupported-dependent-value": {
			request:     boolPresenceRequest(stringPtr("available"), false, boolPtr(true)),
			expectError: true,
		},
		// False is just as much a configured value as true: the service rejects the presence
		// of the field, not the value it carries.
		"set-false-with-unsupported-dependent-value": {
			request:     boolPresenceRequest(stringPtr("available"), false, boolPtr(false)),
			expectError: true,
		},
		"omitted-with-unsupported-dependent-value": {
			request:     boolPresenceRequest(stringPtr("available"), false, nil),
			expectError: false,
		},
		"unknown-dependent-value-defers": {
			request:     boolPresenceRequest(nil, true, boolPtr(true)),
			expectError: false,
		},
		"null-dependent-value-defers": {
			request:     boolPresenceRequest(nil, false, boolPtr(true)),
			expectError: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			response := &validator.BoolResponse{}
			BoolSupportedOnlyWhenStringIn(path.Root("intent"), []string{"required"}, "").
				ValidateBool(context.Background(), testCase.request, response)

			assert.Equal(t, testCase.expectError, response.Diagnostics.HasError())
		})
	}
}

func TestBoolUnsupportedWhenStringIn(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request     validator.BoolRequest
		expectError bool
	}{
		"set-with-rejected-dependent-value": {
			request:     boolPresenceRequest(stringPtr("uninstall"), false, boolPtr(true)),
			expectError: true,
		},
		"set-false-with-rejected-dependent-value": {
			request:     boolPresenceRequest(stringPtr("uninstall"), false, boolPtr(false)),
			expectError: true,
		},
		"set-with-accepted-dependent-value": {
			request:     boolPresenceRequest(stringPtr("available"), false, boolPtr(true)),
			expectError: false,
		},
		"omitted-with-rejected-dependent-value": {
			request:     boolPresenceRequest(stringPtr("uninstall"), false, nil),
			expectError: false,
		},
		"unknown-dependent-value-defers": {
			request:     boolPresenceRequest(nil, true, boolPtr(true)),
			expectError: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			response := &validator.BoolResponse{}
			BoolUnsupportedWhenStringIn(path.Root("intent"), []string{"uninstall"}, "").
				ValidateBool(context.Background(), testCase.request, response)

			assert.Equal(t, testCase.expectError, response.Diagnostics.HasError())
		})
	}
}

func TestBoolPresenceValidatorCustomMessage(t *testing.T) {
	t.Parallel()

	const message = "`is_removable` can only be set when `intent` is `required`."

	response := &validator.BoolResponse{}
	BoolSupportedOnlyWhenStringIn(path.Root("intent"), []string{"required"}, message).
		ValidateBool(context.Background(), boolPresenceRequest(stringPtr("available"), false, boolPtr(true)), response)

	assert.True(t, response.Diagnostics.HasError())
	assert.Equal(t, message, response.Diagnostics.Errors()[0].Detail())
}
