package graphBetaGroupPolicyDefinition

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
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
		return fmt.Errorf("validation failed: no presentations found for policy '%s'", data.PolicyName.ValueString())
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
			return fmt.Errorf("validation failed for value[%d]: label '%s' not found in policy. Available labels: %s",
				i, label, strings.Join(availableLabels, ", "))
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
			return fmt.Errorf("checkbox '%s' requires a boolean value ('true' or 'false'), got '%s'", label, value)
		}

	case "#microsoft.graph.groupPolicyPresentationTextBox":
		// Any string is valid
		if value == "" {
			tflog.Warn(ctx, fmt.Sprintf("[VALIDATE] Textbox '%s' has empty value", label))
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
			return fmt.Errorf("dropdown/combobox '%s' requires a non-empty value", label)
		}

	case "#microsoft.graph.groupPolicyPresentationText",
		"#microsoft.graph.groupPolicyPresentationLabel":
		// Read-only presentation types - users shouldn't provide values for these
		return fmt.Errorf("presentation '%s' is read-only (type: %s) and cannot have a value set", label, odataType)

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

// ConvertValueForType converts a string value to the appropriate type for the Graph API
func ConvertValueForType(ctx context.Context, label, value, odataType string) (graphmodels.GroupPolicyPresentationValueable, error) {
	// Map presentation types to presentation VALUE types
	var valueOdataType string

	switch odataType {
	case "#microsoft.graph.groupPolicyPresentationCheckBox":
		valueOdataType = "#microsoft.graph.groupPolicyPresentationValueBoolean"
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean value for checkbox '%s': %w", label, err)
		}
		presValue := graphmodels.NewGroupPolicyPresentationValueBoolean()
		presValue.SetOdataType(&valueOdataType)
		presValue.SetValue(&boolVal)
		return presValue, nil

	case "#microsoft.graph.groupPolicyPresentationTextBox":
		valueOdataType = "#microsoft.graph.groupPolicyPresentationValueText"
		presValue := graphmodels.NewGroupPolicyPresentationValueText()
		presValue.SetOdataType(&valueOdataType)
		presValue.SetValue(&value)
		return presValue, nil

	case "#microsoft.graph.groupPolicyPresentationMultiTextBox":
		valueOdataType = "#microsoft.graph.groupPolicyPresentationValueMultiText"
		presValue := graphmodels.NewGroupPolicyPresentationValueMultiText()
		presValue.SetOdataType(&valueOdataType)
		// Split by newline or comma for multi-text
		values := splitMultiTextValue(value)
		presValue.SetValues(values)
		return presValue, nil

	case "#microsoft.graph.groupPolicyPresentationDecimalTextBox":
		valueOdataType = "#microsoft.graph.groupPolicyPresentationValueDecimal"
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid numeric value for decimal textbox '%s': %w", label, err)
		}
		presValue := graphmodels.NewGroupPolicyPresentationValueDecimal()
		presValue.SetOdataType(&valueOdataType)
		presValue.SetValue(&intVal)
		return presValue, nil

	case "#microsoft.graph.groupPolicyPresentationDropdownList",
		"#microsoft.graph.groupPolicyPresentationComboBox":
		// Both dropdown and combobox use text value type
		valueOdataType = "#microsoft.graph.groupPolicyPresentationValueText"
		presValue := graphmodels.NewGroupPolicyPresentationValueText()
		presValue.SetOdataType(&valueOdataType)
		presValue.SetValue(&value)
		return presValue, nil

	default:
		return nil, fmt.Errorf("unsupported presentation type '%s' for label '%s'", odataType, label)
	}
}

// splitMultiTextValue splits a multi-text value by newlines or commas
func splitMultiTextValue(value string) []string {
	// First try splitting by newline
	if strings.Contains(value, "\n") {
		lines := strings.Split(value, "\n")
		result := make([]string, 0, len(lines))
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}

	// Fall back to comma-separated
	if strings.Contains(value, ",") {
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}

	// Single value
	return []string{value}
}

// ExtractValueFromPresentation extracts the string value from a presentation value for state mapping
func ExtractValueFromPresentation(presValue graphmodels.GroupPolicyPresentationValueable) string {
	if presValue == nil {
		return ""
	}

	odataType := presValue.GetOdataType()
	if odataType == nil {
		return ""
	}

	switch *odataType {
	case "#microsoft.graph.groupPolicyPresentationValueBoolean":
		if boolVal, ok := presValue.(graphmodels.GroupPolicyPresentationValueBooleanable); ok {
			if val := boolVal.GetValue(); val != nil {
				return strconv.FormatBool(*val)
			}
		}

	case "#microsoft.graph.groupPolicyPresentationValueText":
		if textVal, ok := presValue.(graphmodels.GroupPolicyPresentationValueTextable); ok {
			if val := textVal.GetValue(); val != nil {
				return *val
			}
		}

	case "#microsoft.graph.groupPolicyPresentationValueMultiText":
		if multiTextVal, ok := presValue.(graphmodels.GroupPolicyPresentationValueMultiTextable); ok {
			if vals := multiTextVal.GetValues(); vals != nil {
				return strings.Join(vals, "\n")
			}
		}

	case "#microsoft.graph.groupPolicyPresentationValueDecimal":
		if decimalVal, ok := presValue.(graphmodels.GroupPolicyPresentationValueDecimalable); ok {
			if val := decimalVal.GetValue(); val != nil {
				return strconv.FormatInt(int64(*val), 10)
			}
		}
	}

	return ""
}
