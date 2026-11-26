package graphBetaGroupPolicyDefinition

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
)

// ValidateValues validates user-provided values against the Graph API catalog
// to ensure labels exist and values are appropriate for their presentation types
func ValidateValues(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *GroupPolicyDefinitionResourceModel) error {
	if data.Values.IsNull() || data.Values.IsUnknown() {
		return nil
	}

	// Get resolved presentations from AdditionalData (set by resolver)
	resolvedPresentations, ok := data.AdditionalData["resolvedPresentations"].([]ResolvedPresentation)
	if !ok || len(resolvedPresentations) == 0 {
		return fmt.Errorf("%w for policy '%s'", sentinels.ErrNoPresentationsFound, data.PolicyName.ValueString())
	}

	tflog.Debug(ctx, fmt.Sprintf("[VALIDATE] Validating %d resolved presentations", len(resolvedPresentations)))

	// Build a map of label -> presentation for validation
	labelToPres := make(map[string]ResolvedPresentation)
	for _, pres := range resolvedPresentations {
		labelToPres[pres.Label] = pres
	}

	// Get user-provided values
	var userValues []PresentationValue
	data.Values.ElementsAs(ctx, &userValues, false)

	tflog.Debug(ctx, fmt.Sprintf("[VALIDATE] Validating %d user-provided values", len(userValues)))

	// Validate each user value
	for i, userValue := range userValues {
		label := userValue.Label.ValueString()
		value := userValue.Value.ValueString()

		// Check if label exists
		pres, found := labelToPres[label]
		if !found {
			availableLabels := make([]string, 0, len(labelToPres))
			for l := range labelToPres {
				availableLabels = append(availableLabels, l)
			}
			return fmt.Errorf("validation failed for value[%d]: %w: '%s'. Available labels: %s",
				i, sentinels.ErrLabelNotFound, label, strings.Join(availableLabels, ", "))
		}

		// Validate value format based on presentation type
		if err := validateValueForType(ctx, label, value, pres.Type); err != nil {
			return fmt.Errorf("validation failed for value[%d] (label='%s'): %w", i, label, err)
		}

		tflog.Debug(ctx, fmt.Sprintf("[VALIDATE] ✓ Value[%d]: label='%s', type='%s', value='%s'", i, label, pres.Type, value))
	}

	tflog.Debug(ctx, "[VALIDATE] All values validated successfully")
	return nil
}

// validateValueForType validates that a value string is appropriate for the given presentation type
func validateValueForType(ctx context.Context, label, value, odataType string) error {
	switch odataType {
	case "#microsoft.graph.groupPolicyPresentationCheckBox":
		// Must be a boolean
		if _, err := strconv.ParseBool(value); err != nil {
			return fmt.Errorf("%w for checkbox '%s', got '%s'", sentinels.ErrCheckboxRequiresBooleanValue, label, value)
		}

	case "#microsoft.graph.groupPolicyPresentationTextBox":
		// Any string is valid
		if value == "" {
			tflog.Warn(ctx, fmt.Sprintf("[VALIDATE] Textbox '%s' has empty value", label))
		}

	case "#microsoft.graph.groupPolicyPresentationListBox":
		if _, err := parseListBoxValue(value); err != nil {
			return fmt.Errorf("invalid ListBox value for '%s': %w", label, err)
		}
	case "#microsoft.graph.groupPolicyPresentationMultiTextBox":
		// For multi-text, we expect comma-separated or newline-separated values
		// Any string is technically valid
		if value == "" {
			tflog.Warn(ctx, fmt.Sprintf("[VALIDATE] Multi-textbox '%s' has empty value", label))
		}

	case "#microsoft.graph.groupPolicyPresentationDecimalTextBox":
		// Must be numeric
		if _, err := strconv.ParseInt(value, 10, 64); err != nil {
			return fmt.Errorf("decimal textbox '%s' requires a numeric value, got '%s': %w", label, value, err)
		}

	case "#microsoft.graph.groupPolicyPresentationDropdownList",
		"#microsoft.graph.groupPolicyPresentationComboBox":
		// Would need to validate against allowed options, but that requires fetching dropdown items
		// For now, just ensure non-empty
		if value == "" {
			return fmt.Errorf("%w: '%s'", sentinels.ErrDropdownRequiresNonEmptyValue, label)
		}

	case "#microsoft.graph.groupPolicyPresentationText",
		"#microsoft.graph.groupPolicyPresentationLabel":
		// Read-only presentation types - users shouldn't provide values for these
		return fmt.Errorf("%w: '%s' (type: %s)", sentinels.ErrPresentationReadOnly, label, odataType)

	default:
		// Unknown type - log warning but don't fail
		tflog.Warn(ctx, fmt.Sprintf("[VALIDATE] Unknown presentation type '%s' for label '%s', skipping type validation", odataType, label))
	}

	return nil
}

// GetPresentationTypeName returns a human-readable name for a presentation type
func GetPresentationTypeName(odataType string) string {
	switch odataType {
	case "#microsoft.graph.groupPolicyPresentationCheckBox":
		return "CheckBox (boolean)"
	case "#microsoft.graph.groupPolicyPresentationTextBox":
		return "TextBox (string)"
	case "#microsoft.graph.groupPolicyPresentationMultiTextBox":
		return "MultiTextBox (string array)"
	case "#microsoft.graph.groupPolicyPresentationDecimalTextBox":
		return "DecimalTextBox (integer)"
	case "#microsoft.graph.groupPolicyPresentationDropdownList":
		return "DropdownList (select option)"
	case "#microsoft.graph.groupPolicyPresentationComboBox":
		return "ComboBox (select or type)"
	case "#microsoft.graph.groupPolicyPresentationText":
		return "Text (display only)"
	case "#microsoft.graph.groupPolicyPresentationLabel":
		return "Label (display only)"
	case "#microsoft.graph.groupPolicyPresentationListBox":
		return "ListBox (key-value pairs)"
	default:
		return "Unknown"
	}
}
