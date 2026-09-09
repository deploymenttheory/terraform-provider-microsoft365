package graphBetaNetworkPromptPolicyRule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// promptPolicyRuleConsistencyPredicate verifies that a successful write is visible on readback.
// The portal-backed endpoint can accept a PATCH while a stale replica still returns old values.
func promptPolicyRuleConsistencyPredicate(
	expected *NetworkPromptPolicyRuleResourceModel,
) func(context.Context, tfsdk.State) bool {
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual NetworkPromptPolicyRuleResourceModel
		if diags := state.Get(ctx, &actual); diags.HasError() {
			return false
		}
		if actual.ID.IsNull() || actual.ID.IsUnknown() || actual.ID.ValueString() == "" {
			return false
		}
		return (expected.Name.IsUnknown() || expected.Name.Equal(actual.Name)) &&
			(expected.Description.IsUnknown() || expected.Description.Equal(actual.Description)) &&
			(expected.Action.IsUnknown() || expected.Action.Equal(actual.Action)) &&
			(expected.Priority.IsUnknown() || expected.Priority.Equal(actual.Priority)) &&
			(expected.Enabled.IsUnknown() || expected.Enabled.Equal(actual.Enabled)) &&
			(expected.PromptLogging.IsUnknown() || expected.PromptLogging.Equal(actual.PromptLogging)) &&
			(expected.ScanResult.IsUnknown() || expected.ScanResult.Equal(actual.ScanResult)) &&
			conversationSchemesConsistent(ctx, expected.ConversationSchemes, actual.ConversationSchemes)
	}
}

func conversationSchemesConsistent(ctx context.Context, expected, actual types.List) bool {
	if expected.IsUnknown() {
		return true
	}
	if expected.IsNull() || actual.IsNull() || actual.IsUnknown() {
		return expected.Equal(actual)
	}

	var expectedSchemes, actualSchemes []ConversationSchemeModel
	if expectedDiags := expected.ElementsAs(ctx, &expectedSchemes, false); expectedDiags.HasError() {
		return false
	}
	if actualDiags := actual.ElementsAs(ctx, &actualSchemes, false); actualDiags.HasError() {
		return false
	}
	if len(expectedSchemes) != len(actualSchemes) {
		return false
	}

	for i := range expectedSchemes {
		expectedScheme := expectedSchemes[i]
		actualScheme := actualSchemes[i]
		if !stringValueConsistent(expectedScheme.Type, actualScheme.Type) {
			return false
		}
		if expectedScheme.Type.IsUnknown() {
			continue
		}

		switch expectedScheme.Type.ValueString() {
		case schemeTypeCustom:
			if !stringValueConsistent(expectedScheme.URL, actualScheme.URL) ||
				!jsonPathConsistent(expectedScheme.JSONPath, actualScheme.JSONPath) {
				return false
			}
		case schemeTypePredefined:
			if !stringValueConsistent(expectedScheme.SchemeName, actualScheme.SchemeName) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func stringValueConsistent(expected, actual types.String) bool {
	return expected.IsUnknown() || expected.Equal(actual)
}

func jsonPathConsistent(expected, actual types.String) bool {
	if stringValueConsistent(expected, actual) {
		return true
	}
	return expected.IsNull() && !actual.IsNull() && !actual.IsUnknown() && actual.ValueString() == ""
}
