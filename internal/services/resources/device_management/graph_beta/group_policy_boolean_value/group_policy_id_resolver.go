package graphBetaGroupPolicyBooleanValue

import (
	"context"
	"fmt"
	"strings"

	groupPolicyResolver "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/resolvers/graph_beta/device_management/group_policy_configurations"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// GroupPolicyIDResolver is a wrapper around the centralized resolver for boolean value resources
func GroupPolicyIDResolver(ctx context.Context, data *GroupPolicyBooleanValueResourceModel, client *msgraphbetasdk.GraphServiceClient, operation string) (error, int) {
	tflog.Debug(ctx, fmt.Sprintf("[BOOLEAN_VALUE] GroupPolicyIDResolver: Starting %s operation", operation))

	if data.PolicyName.IsNull() || data.PolicyName.IsUnknown() ||
		data.ClassType.IsNull() || data.ClassType.IsUnknown() ||
		data.CategoryPath.IsNull() || data.CategoryPath.IsUnknown() {
		return fmt.Errorf("provide policy_name, class_type, and category_path for auto-discovery"), 0
	}

	// Use centralized resolver with checkbox presentation filter for boolean resources
	result, err := groupPolicyResolver.GroupPolicyIDResolver(
		ctx,
		client,
		operation,
		data.PolicyName.ValueString(),
		data.ClassType.ValueString(),
		data.CategoryPath.ValueString(),
		data.GroupPolicyConfigurationID.ValueString(),
		"#microsoft.graph.groupPolicyPresentationCheckBox", // Filter for checkbox presentations only
	)

	if err != nil {
		statusCode := 0
		if strings.Contains(err.Error(), "status: ") {
			fmt.Sscanf(err.Error(), "%*s status: %d", &statusCode)
		}
		return err, statusCode
	}

	switch operation {
	case "create":
		// For creation, store template IDs and resolved presentations
		data.GroupPolicyDefinitionValueID = types.StringValue(result.DefinitionTemplateID)
		data.ID = types.StringValue(result.DefinitionTemplateID) // Use definition ID as the main ID

		// Store resolved presentations for construct to use
		if data.AdditionalData == nil {
			data.AdditionalData = make(map[string]any)
		}

		// Convert centralized ResolvedPresentation to local ResolvedPresentation
		var resolvedPresentations []ResolvedPresentation
		for i, templateID := range result.PresentationTemplateIDs {
			resolvedPresentations = append(resolvedPresentations, ResolvedPresentation{
				TemplateID: templateID,
				Index:      i,
			})
		}
		data.AdditionalData["resolvedPresentations"] = resolvedPresentations

		tflog.Debug(ctx, fmt.Sprintf("[BOOLEAN_VALUE] CREATE operation - using template IDs (definitionTemplateID: %s, presentations: %d)",
			result.DefinitionTemplateID, len(result.PresentationTemplateIDs)))

	case "update":
		// For update, we need both template IDs (for bindings) and instance IDs (for the update)
		data.GroupPolicyDefinitionValueID = types.StringValue(result.DefinitionTemplateID) // Use template ID for binding
		data.ID = types.StringValue(result.DefinitionValueInstanceID)                      // This is the definition value instance ID

		// Store instance IDs in additional data for construct to use
		if data.AdditionalData == nil {
			data.AdditionalData = make(map[string]any)
		}
		data.AdditionalData["definitionValueInstanceID"] = result.DefinitionValueInstanceID

		// Convert centralized ResolvedPresentation to local ResolvedPresentation
		var resolvedPresentations []ResolvedPresentation
		for _, resolved := range result.ResolvedPresentations {
			resolvedPresentations = append(resolvedPresentations, ResolvedPresentation{
				TemplateID: resolved.TemplateID,
				InstanceID: resolved.InstanceID,
				Index:      resolved.Index,
			})
		}
		data.AdditionalData["resolvedPresentations"] = resolvedPresentations

		tflog.Debug(ctx, fmt.Sprintf("[BOOLEAN_VALUE] UPDATE operation - using template IDs for bindings (definitionTemplateID: %s) and instance IDs for update (definitionValueInstanceID: %s, presentations: %d)",
			result.DefinitionTemplateID, result.DefinitionValueInstanceID, len(resolvedPresentations)))

	case "read":
		// For read, use resolved instance IDs
		data.GroupPolicyDefinitionValueID = types.StringValue(result.DefinitionValueInstanceID)
		data.ID = types.StringValue(result.DefinitionValueInstanceID) // Use definition value instance ID

		// Store resolved presentations in additional data for state mapping
		if data.AdditionalData == nil {
			data.AdditionalData = make(map[string]any)
		}

		// Convert centralized ResolvedPresentation to local ResolvedPresentation
		var resolvedPresentations []ResolvedPresentation
		for _, resolved := range result.ResolvedPresentations {
			resolvedPresentations = append(resolvedPresentations, ResolvedPresentation{
				TemplateID: resolved.TemplateID,
				InstanceID: resolved.InstanceID,
				Index:      resolved.Index,
			})
		}
		data.AdditionalData["resolvedPresentations"] = resolvedPresentations

		tflog.Debug(ctx, fmt.Sprintf("[BOOLEAN_VALUE] READ operation - using instance IDs (definitionValueInstanceID: %s, presentations: %d)",
			result.DefinitionValueInstanceID, len(resolvedPresentations)))

	default:
		return fmt.Errorf("unsupported crud operation '%s' - must be 'create', 'read', or 'update'", operation), 0
	}

	return nil, 0
}
