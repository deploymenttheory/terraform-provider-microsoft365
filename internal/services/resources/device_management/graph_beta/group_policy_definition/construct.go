package graphBetaGroupPolicyDefinition

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
)

// constructResource constructs the updateDefinitionValues request for group policy definition operations
func constructResource(
	ctx context.Context,
	data *GroupPolicyDefinitionResourceModel,
	operation string,
) (*devicemanagement.GroupPolicyConfigurationsItemUpdateDefinitionValuesPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing updateDefinitionValues request for %s", ResourceName))

	requestBody := devicemanagement.NewGroupPolicyConfigurationsItemUpdateDefinitionValuesPostRequestBody()

	definitionValue := models.NewGroupPolicyDefinitionValue()

	convert.FrameworkToGraphBool(data.Enabled, definitionValue.SetEnabled)

	definitionTemplateID := data.AdditionalData["definitionTemplateID"].(string)
	definitionBindURL := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')", definitionTemplateID)

	definitionValue.SetAdditionalData(map[string]any{
		"definition@odata.bind": definitionBindURL,
	})

	var presentationValues []models.GroupPolicyPresentationValueable

	// For delete operations, we don't need to construct presentation values
	if operation != constants.TfOperationDelete {
		// Get resolved presentations from AdditionalData
		resolvedPresentations, ok := data.AdditionalData["resolvedPresentations"].([]ResolvedPresentation)
		if !ok {
			return nil, fmt.Errorf("missing resolved presentations in AdditionalData")
		}

		// Build a map of label -> resolved presentation for quick lookup
		labelToPresentation := make(map[string]ResolvedPresentation)
		for _, resolved := range resolvedPresentations {
			if resolved.Label != "" {
				labelToPresentation[resolved.Label] = resolved
			}
		}

		// Get user-provided values
		var userValues []PresentationValue
		if !data.Values.IsNull() && !data.Values.IsUnknown() {
			data.Values.ElementsAs(ctx, &userValues, false)
		}

		for _, userValue := range userValues {
			label := userValue.Label.ValueString()
			value := userValue.Value.ValueString()

			// Find the matching presentation
			resolved, found := labelToPresentation[label]
			if !found {
				return nil, fmt.Errorf("no presentation found with label '%s'", label)
			}

			// Convert the value based on presentation type
			presValue, err := ConvertValueForType(ctx, label, value, resolved.Type)
			if err != nil {
				return nil, fmt.Errorf("failed to convert value for label '%s': %w", label, err)
			}

			// Set the presentation binding
			presentationBindURL := fmt.Sprintf("https://graph.microsoft.com/beta/deviceManagement/groupPolicyDefinitions('%s')/presentations('%s')",
				definitionTemplateID, resolved.TemplateID)

			presValue.SetAdditionalData(map[string]any{
				"presentation@odata.bind": presentationBindURL,
			})

			// For update operations, set the instance ID
			if operation == constants.TfOperationUpdate && resolved.InstanceID != "" {
				presValue.SetId(&resolved.InstanceID)
			}

			presentationValues = append(presentationValues, presValue)

			tflog.Debug(ctx, fmt.Sprintf("[CONSTRUCT] Added presentation value: label='%s', type='%s', value='%s'", label, resolved.Type, value))
		}
	}

	definitionValue.SetPresentationValues(presentationValues)

	switch operation {
	case constants.TfOperationCreate:
		addedValues := []models.GroupPolicyDefinitionValueable{definitionValue}
		requestBody.SetAdded(addedValues)
		requestBody.SetUpdated([]models.GroupPolicyDefinitionValueable{})
		requestBody.SetDeletedIds([]string{})

	case constants.TfOperationUpdate:
		// For update, use the instance ID
		definitionValueInstanceID, ok := data.AdditionalData["definitionValueInstanceID"].(string)
		if !ok {
			return nil, fmt.Errorf("missing definitionValueInstanceID in AdditionalData for update operation")
		}

		definitionValue.SetId(&definitionValueInstanceID)

		updatedValues := []models.GroupPolicyDefinitionValueable{definitionValue}
		requestBody.SetAdded([]models.GroupPolicyDefinitionValueable{})
		requestBody.SetUpdated(updatedValues)
		requestBody.SetDeletedIds([]string{})

	case constants.TfOperationDelete:
		definitionValueInstanceID := data.AdditionalData["definitionValueInstanceID"].(string)
		requestBody.SetAdded([]models.GroupPolicyDefinitionValueable{})
		requestBody.SetUpdated([]models.GroupPolicyDefinitionValueable{})
		requestBody.SetDeletedIds([]string{definitionValueInstanceID})
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final updateDefinitionValues JSON for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing updateDefinitionValues request for %s", ResourceName))
	return requestBody, nil
}

// ConvertValueForType converts a string value to the appropriate type for the Graph API
func ConvertValueForType(ctx context.Context, label, value, odataType string) (models.GroupPolicyPresentationValueable, error) {
	// Map presentation types to presentation VALUE types
	var valueOdataType string

	switch odataType {
	case "#microsoft.graph.groupPolicyPresentationCheckBox":
		valueOdataType = "#microsoft.graph.groupPolicyPresentationValueBoolean"
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean value for checkbox '%s': %w", label, err)
		}
		presValue := models.NewGroupPolicyPresentationValueBoolean()
		presValue.SetOdataType(&valueOdataType)
		presValue.SetValue(&boolVal)
		return presValue, nil

	case "#microsoft.graph.groupPolicyPresentationListBox":
		pairs, err := parseListBoxValue(value)
		if err != nil {
			return nil, fmt.Errorf("invalid ListBox value for '%s': %w", label, err)
		}
		presValue := models.NewGroupPolicyPresentationValueList()
		valueOdataType = "#microsoft.graph.groupPolicyPresentationValueList"
		presValue.SetOdataType(&valueOdataType)
		values := make([]models.KeyValuePairable, 0, len(pairs))
		for _, pair := range pairs {
			entry := models.NewKeyValuePair()
			entry.SetName(pair.Name)
			entry.SetValue(pair.Value)
			values = append(values, entry)
		}
		presValue.SetValues(values)
		return presValue, nil
	case "#microsoft.graph.groupPolicyPresentationTextBox":
		valueOdataType = "#microsoft.graph.groupPolicyPresentationValueText"
		presValue := models.NewGroupPolicyPresentationValueText()
		presValue.SetOdataType(&valueOdataType)
		presValue.SetValue(&value)
		return presValue, nil

	case "#microsoft.graph.groupPolicyPresentationMultiTextBox":
		valueOdataType = "#microsoft.graph.groupPolicyPresentationValueMultiText"
		presValue := models.NewGroupPolicyPresentationValueMultiText()
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
		presValue := models.NewGroupPolicyPresentationValueDecimal()
		presValue.SetOdataType(&valueOdataType)
		presValue.SetValue(&intVal)
		return presValue, nil

	case "#microsoft.graph.groupPolicyPresentationDropdownList",
		"#microsoft.graph.groupPolicyPresentationComboBox":
		// Both dropdown and combobox use text value type
		valueOdataType = "#microsoft.graph.groupPolicyPresentationValueText"
		presValue := models.NewGroupPolicyPresentationValueText()
		presValue.SetOdataType(&valueOdataType)
		presValue.SetValue(&value)
		return presValue, nil

	default:
		return nil, fmt.Errorf("%w: '%s' for label '%s'", sentinels.ErrUnsupportedPresentationType, odataType, label)
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

var errInvalidListBox = errors.New("expected a JSON array of name/value string pairs")

// listBoxPair preserves the order and text of Graph ListBox name/value entries.
type listBoxPair struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

func parseListBoxValue(value string) ([]listBoxPair, error) {
	var pairs []listBoxPair
	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&pairs); err != nil {
		return nil, fmt.Errorf("expected a JSON array of name/value string pairs: %w", err)
	}
	if !json.Valid([]byte(value)) || pairs == nil {
		return nil, errInvalidListBox
	}
	for i, pair := range pairs {
		if pair.Name == nil || pair.Value == nil {
			return nil, fmt.Errorf("%w: entry %d requires name and value strings", errInvalidListBox, i)
		}
	}
	return pairs, nil
}
