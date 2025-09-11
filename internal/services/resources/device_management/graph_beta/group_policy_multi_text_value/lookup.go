package graphBetaGroupPolicyMultiTextValue

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// LookupService handles the dynamic resolution of definition value and presentation IDs
type LookupService struct {
	client *msgraphbetasdk.GraphServiceClient
}

// NewLookupService creates a new lookup service
func NewLookupService(client *msgraphbetasdk.GraphServiceClient) *LookupService {
	return &LookupService{
		client: client,
	}
}

// ResolveDefinitionValueAndPresentation resolves both the definition value ID and presentation ID
// based on the simplified input parameters
func (s *LookupService) ResolveDefinitionValueAndPresentation(
	ctx context.Context,
	groupPolicyConfigurationID, policyName, classType string,
	presentationIndex int64,
) (definitionValueID, presentationID string, err error) {

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Resolving IDs for policy='%s', classType='%s', presentationIndex=%d",
		policyName, classType, presentationIndex))

	definitionValueID, err = s.findDefinitionValueID(ctx, groupPolicyConfigurationID, policyName, classType)
	if err != nil {
		return "", "", fmt.Errorf("failed to find definition value: %w", err)
	}

	presentationID, err = s.findPresentationID(ctx, groupPolicyConfigurationID, definitionValueID, presentationIndex)
	if err != nil {
		return "", "", fmt.Errorf("failed to find presentation: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Successfully resolved - definitionValueID='%s', presentationID='%s'",
		definitionValueID, presentationID))

	return definitionValueID, presentationID, nil
}

// findDefinitionValueID looks up the definition value ID within a specific group policy configuration
func (s *LookupService) findDefinitionValueID(ctx context.Context, configID, policyName, classType string) (string, error) {
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Finding definition value ID for policy='%s', classType='%s' in config='%s'",
		policyName, classType, configID))

	definitionValues, err := s.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(configID).
		DefinitionValues().
		Get(ctx, nil)

	if err != nil {
		return "", fmt.Errorf("failed to get definition values: %w", err)
	}

	if definitionValues == nil || definitionValues.GetValue() == nil {
		return "", fmt.Errorf("no definition values found in configuration")
	}

	normalizedClassType := strings.ToLower(classType)

	var matches []string
	for _, defValue := range definitionValues.GetValue() {
		if defValue == nil || defValue.GetDefinition() == nil {
			continue
		}

		definition := defValue.GetDefinition()
		if definition.GetDisplayName() == nil || definition.GetClassType() == nil || defValue.GetId() == nil {
			continue
		}

		defDisplayName := *definition.GetDisplayName()
		defClassType := strings.ToLower(definition.GetClassType().String())
		defValueID := *defValue.GetId()

		if defDisplayName == policyName && defClassType == normalizedClassType {
			matches = append(matches, defValueID)
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found matching definition value: ID='%s', name='%s', classType='%s'",
				defValueID, defDisplayName, defClassType))
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no definition value found for policy name '%s' with class type '%s'", policyName, classType)
	}

	if len(matches) > 1 {
		tflog.Warn(ctx, fmt.Sprintf("[LOOKUP] Multiple definition values found (%d), using first match", len(matches)))
	}

	return matches[0], nil
}

// findPresentationID finds a suitable multi-text presentation for the given definition value
func (s *LookupService) findPresentationID(ctx context.Context, configID, definitionValueID string, presentationIndex int64) (string, error) {
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Finding presentation ID for definitionValueID='%s', presentationIndex=%d",
		definitionValueID, presentationIndex))

	definitionValue, err := s.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(configID).
		DefinitionValues().
		ByGroupPolicyDefinitionValueId(definitionValueID).
		Get(ctx, nil)

	if err != nil {
		return "", fmt.Errorf("failed to get definition value: %w", err)
	}

	if definitionValue == nil || definitionValue.GetDefinition() == nil || definitionValue.GetDefinition().GetId() == nil {
		return "", fmt.Errorf("definition value or its definition is nil")
	}

	definitionID := *definitionValue.GetDefinition().GetId()

	presentations, err := s.client.
		DeviceManagement().
		GroupPolicyDefinitions().
		ByGroupPolicyDefinitionId(definitionID).
		Presentations().
		Get(ctx, nil)

	if err != nil {
		return "", fmt.Errorf("failed to get presentations: %w", err)
	}

	if presentations == nil || presentations.GetValue() == nil {
		return "", fmt.Errorf("no presentations found for definition")
	}

	// Find suitable multi-text presentations
	var suitablePresentations []string
	for _, presentation := range presentations.GetValue() {
		if presentation == nil || presentation.GetId() == nil {
			continue
		}

		// Check if this presentation supports multi-text values
		// Look for presentations that are listBox, multiTextBox, or similar multi-value types
		odataType := presentation.GetOdataType()
		if odataType != nil {
			presentationType := *odataType
			if s.isSuitableForMultiText(presentationType) {
				suitablePresentations = append(suitablePresentations, *presentation.GetId())
				tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found suitable presentation: ID='%s', type='%s'",
					*presentation.GetId(), presentationType))
			}
		}
	}

	if len(suitablePresentations) == 0 {
		return "", fmt.Errorf("no suitable multi-text presentations found for definition")
	}

	// Use the specified index or default to 0
	index := int(presentationIndex)
	if index < 0 || index >= len(suitablePresentations) {
		if presentationIndex != 0 {
			tflog.Warn(ctx, fmt.Sprintf("[LOOKUP] Presentation index %d out of range (0-%d), using index 0",
				presentationIndex, len(suitablePresentations)-1))
		}
		index = 0
	}

	return suitablePresentations[index], nil
}

// isSuitableForMultiText checks if a presentation type can handle multiple text values
func (s *LookupService) isSuitableForMultiText(odataType string) bool {
	// List of presentation types that can handle multiple text values
	multiTextTypes := []string{
		"#microsoft.graph.groupPolicyPresentationListBox",
		"#microsoft.graph.groupPolicyPresentationMultiTextBox",
		// Add other types that support multi-text as needed
	}

	for _, validType := range multiTextTypes {
		if odataType == validType {
			return true
		}
	}

	return false
}
