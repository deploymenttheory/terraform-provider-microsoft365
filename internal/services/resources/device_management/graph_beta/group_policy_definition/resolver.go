package graphBetaGroupPolicyDefinition

import (
	"context"
	"fmt"

	groupPolicyResolver "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/resolvers/graph_beta/device_management/group_policy_configurations"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// resolveGroupPolicyDefinition resolves policy metadata and all presentations
func resolveGroupPolicyDefinition(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *GroupPolicyDefinitionResourceModel, operation string) error {
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Resolving group policy definition for operation: %s", operation))

	policyName := data.PolicyName.ValueString()
	classType := data.ClassType.ValueString()
	categoryPath := data.CategoryPath.ValueString()
	configID := data.GroupPolicyConfigurationID.ValueString()

	// Call the centralized resolver without a presentation filter to get ALL presentations
	result, err := groupPolicyResolver.GroupPolicyIDResolver(
		ctx,
		client,
		operation,
		policyName,
		classType,
		categoryPath,
		configID,
		"", // Empty filter = get all presentations
	)
	if err != nil {
		return fmt.Errorf("failed to resolve group policy definition: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Successfully resolved definition template ID: %s", result.DefinitionTemplateID))

	// For CREATE operation
	if operation == "create" {
		// Store the definition template ID (which is also used as the initial definition value ID)
		data.ID = data.GroupPolicyConfigurationID // Composite ID: configID
		data.AdditionalData["definitionTemplateID"] = result.DefinitionTemplateID
		tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] CREATE: definitionTemplateID = %s", result.DefinitionTemplateID))

		// Fetch ALL presentations with their types
		presentations, err := getAllPresentations(ctx, client, result.DefinitionTemplateID)
		if err != nil {
			return fmt.Errorf("failed to fetch presentations: %w", err)
		}

		data.AdditionalData["resolvedPresentations"] = presentations
		tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] CREATE: Stored %d resolved presentations", len(presentations)))
		for i, p := range presentations {
			tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] CREATE:   [%d] TemplateID=%s, Label=%s, Type=%s",
				i, p.TemplateID, p.Label, p.Type))
		}
	}

	// For READ/UPDATE operations
	if operation == "read" || operation == "update" {
		// Store instance IDs from the result
		data.AdditionalData["definitionTemplateID"] = result.DefinitionTemplateID
		data.AdditionalData["definitionValueInstanceID"] = result.DefinitionValueInstanceID

		// Set composite ID: configID/definitionValueID
		configID := data.GroupPolicyConfigurationID.ValueString()
		data.ID = types.StringValue(fmt.Sprintf("%s/%s", configID, result.DefinitionValueInstanceID))

		tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] %s: definitionTemplateID = %s, definitionValueInstanceID = %s, composite ID = %s",
			operation, result.DefinitionTemplateID, result.DefinitionValueInstanceID, data.ID.ValueString()))

		// Fetch and resolve presentation instances
		presentations, err := resolvePresentationInstances(ctx, client, configID, result.DefinitionValueInstanceID, result.DefinitionTemplateID)
		if err != nil {
			return fmt.Errorf("failed to resolve presentation instances: %w", err)
		}

		data.AdditionalData["resolvedPresentations"] = presentations
		tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] %s: Stored %d resolved presentations", operation, len(presentations)))
		for i, p := range presentations {
			tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] %s:   [%d] TemplateID=%s, InstanceID=%s, Label=%s, Type=%s",
				operation, i, p.TemplateID, p.InstanceID, p.Label, p.Type))
		}
	}

	return nil
}

// getAllPresentations fetches all presentations for a definition with their types and labels
func getAllPresentations(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, definitionTemplateID string) ([]ResolvedPresentation, error) {
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Fetching all presentations for definition: %s", definitionTemplateID))

	presentations, err := client.
		DeviceManagement().
		GroupPolicyDefinitions().
		ByGroupPolicyDefinitionId(definitionTemplateID).
		Presentations().
		Get(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to get presentations: %w", err)
	}

	if presentations == nil || presentations.GetValue() == nil || len(presentations.GetValue()) == 0 {
		tflog.Warn(ctx, fmt.Sprintf("[RESOLVER] No presentations found for definition: %s", definitionTemplateID))
		return []ResolvedPresentation{}, nil
	}

	presentationList := presentations.GetValue()
	resolvedList := make([]ResolvedPresentation, 0, len(presentationList))

	for i, pres := range presentationList {
		if pres == nil {
			tflog.Warn(ctx, fmt.Sprintf("[RESOLVER] Skipping nil presentation at index %d", i))
			continue
		}

		presID := pres.GetId()
		presLabel := pres.GetLabel()
		presType := pres.GetOdataType()

		if presID == nil {
			tflog.Warn(ctx, fmt.Sprintf("[RESOLVER] Skipping presentation at index %d with nil ID", i))
			continue
		}

		label := ""
		if presLabel != nil {
			label = *presLabel
		}

		odataType := ""
		if presType != nil {
			odataType = *presType
		}

		resolved := ResolvedPresentation{
			TemplateID: *presID,
			Label:      label,
			Type:       odataType,
			Index:      i,
		}

		resolvedList = append(resolvedList, resolved)
		tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Presentation[%d]: ID=%s, Label='%s', Type=%s", i, *presID, label, odataType))
	}

	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Successfully resolved %d presentations", len(resolvedList)))
	return resolvedList, nil
}

// resolvePresentationInstances maps presentation value instances back to their templates with types
func resolvePresentationInstances(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, configID, definitionValueInstanceID, definitionTemplateID string) ([]ResolvedPresentation, error) {
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Resolving presentation instances for definitionValue: %s", definitionValueInstanceID))

	// Get the presentation values using the collection endpoint with $expand=presentation
	presentationValuesResponse, err := client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(configID).
		DefinitionValues().
		ByGroupPolicyDefinitionValueId(definitionValueInstanceID).
		PresentationValues().
		Get(ctx, &devicemanagement.GroupPolicyConfigurationsItemDefinitionValuesItemPresentationValuesRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.GroupPolicyConfigurationsItemDefinitionValuesItemPresentationValuesRequestBuilderGetQueryParameters{
				Expand: []string{"presentation"},
			},
		})

	if err != nil {
		return nil, fmt.Errorf("failed to get presentation values: %w", err)
	}

	if presentationValuesResponse == nil {
		tflog.Warn(ctx, "[RESOLVER] Presentation values response is nil")
		return []ResolvedPresentation{}, nil
	}

	presentationValues := presentationValuesResponse.GetValue()
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] resolvePresentationInstances: Got %d presentation values from API", len(presentationValues)))
	if len(presentationValues) == 0 {
		tflog.Warn(ctx, "[RESOLVER] No presentation values found")
		return []ResolvedPresentation{}, nil
	}

	// Get all presentation templates to map instance -> template
	allPresentations, err := getAllPresentations(ctx, client, definitionTemplateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get presentation templates: %w", err)
	}

	// Build a map of template ID -> presentation info
	templateMap := make(map[string]ResolvedPresentation)
	for _, pres := range allPresentations {
		templateMap[pres.TemplateID] = pres
	}
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Built templateMap with %d entries", len(templateMap)))

	// Resolve each presentation value instance to its template
	resolvedList := make([]ResolvedPresentation, 0, len(presentationValues))

	for i, presValue := range presentationValues {
		if presValue == nil {
			tflog.Warn(ctx, fmt.Sprintf("[RESOLVER] [%d] presValue is nil", i))
			continue
		}

		instanceID := presValue.GetId()
		if instanceID == nil {
			tflog.Warn(ctx, fmt.Sprintf("[RESOLVER] [%d] instanceID is nil", i))
			continue
		}
		tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] [%d] instanceID = %s", i, *instanceID))

		// Get the presentation reference to find the template ID
		presentation := presValue.GetPresentation()
		if presentation == nil {
			tflog.Warn(ctx, fmt.Sprintf("[RESOLVER] [%d] presentation reference is nil", i))
			continue
		}

		templateID := presentation.GetId()
		if templateID == nil {
			tflog.Warn(ctx, fmt.Sprintf("[RESOLVER] [%d] templateID is nil", i))
			continue
		}
		tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] [%d] templateID = %s", i, *templateID))

		// Look up the template info
		if templateInfo, found := templateMap[*templateID]; found {
			resolved := ResolvedPresentation{
				TemplateID: *templateID,
				InstanceID: *instanceID,
				Label:      templateInfo.Label,
				Type:       templateInfo.Type,
				Index:      i,
			}
			resolvedList = append(resolvedList, resolved)
			tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] [%d] ✓ Resolved: TemplateID=%s, InstanceID=%s, Label='%s', Type=%s",
				i, *templateID, *instanceID, templateInfo.Label, templateInfo.Type))
		} else {
			tflog.Warn(ctx, fmt.Sprintf("[RESOLVER] Could not find template info for ID %s", *templateID))
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Successfully resolved %d presentation instances", len(resolvedList)))
	return resolvedList, nil
}
