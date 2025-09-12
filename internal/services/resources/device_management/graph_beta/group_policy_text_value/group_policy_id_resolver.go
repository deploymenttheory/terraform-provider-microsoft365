package graphBetaGroupPolicyTextValue

import (
	"context"
	"fmt"
	"strings"

	groupPolicyResolver "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/resolvers/graph_beta/device_management/group_policy_configurations"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// GroupPolicyIDResolver is a wrapper around the centralized resolver for text value resources
func GroupPolicyIDResolver(ctx context.Context, data *GroupPolicyTextValueResourceModel, client *msgraphbetasdk.GraphServiceClient, operation string) (error, int) {
	tflog.Debug(ctx, fmt.Sprintf("[TEXT_VALUE] GroupPolicyIDResolver: Starting %s operation", operation))

	if data.PolicyName.IsNull() || data.PolicyName.IsUnknown() ||
		data.ClassType.IsNull() || data.ClassType.IsUnknown() ||
		data.CategoryPath.IsNull() || data.CategoryPath.IsUnknown() {
		return fmt.Errorf("provide policy_name, class_type, and category_path for auto-discovery"), 0
	}

	// no presentation filter for single-value text resources
	result, err := groupPolicyResolver.GroupPolicyIDResolver(
		ctx,
		client,
		operation,
		data.PolicyName.ValueString(),
		data.ClassType.ValueString(),
		data.CategoryPath.ValueString(),
		data.GroupPolicyConfigurationID.ValueString(),
		"", // No presentation filter - accept any presentation type
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
		// For creation, store template IDs only
		data.GroupPolicyDefinitionValueID = types.StringValue(result.DefinitionTemplateID)
		data.PresentationID = types.StringValue(result.GetSinglePresentationID())
		tflog.Debug(ctx, fmt.Sprintf("[TEXT_VALUE] CREATE operation - using template IDs (definitionTemplateID: %s, presentationTemplateID: %s)",
			result.DefinitionTemplateID, result.GetSinglePresentationID()))

	case "update":
		// For update, we need both template IDs (for bindings) and instance IDs (for the update)
		resolved := result.GetSingleResolvedPresentation()
		if resolved == nil {
			return fmt.Errorf("no resolved presentation found for update operation"), 0
		}

		// Store template IDs for bindings, but we'll pass instance IDs via a different mechanism
		data.GroupPolicyDefinitionValueID = types.StringValue(result.DefinitionTemplateID) // Use template ID for binding
		data.PresentationID = types.StringValue(resolved.TemplateID)                       // Keep template ID for binding
		data.ID = types.StringValue(resolved.InstanceID)                                   // This is the presentation value instance ID

		// Store instance IDs in additional data for construct to use
		if data.AdditionalData == nil {
			data.AdditionalData = make(map[string]any)
		}
		data.AdditionalData["definitionValueInstanceID"] = result.DefinitionValueInstanceID
		data.AdditionalData["presentationValueInstanceID"] = resolved.InstanceID

		tflog.Debug(ctx, fmt.Sprintf("[TEXT_VALUE] UPDATE operation - using template IDs for bindings (definitionTemplateID: %s, presentationTemplateID: %s) and instance IDs for update (definitionValueInstanceID: %s, presentationValueInstanceID: %s)",
			result.DefinitionTemplateID, resolved.TemplateID, result.DefinitionValueInstanceID, resolved.InstanceID))

	case "read":
		// For read, use resolved instance IDs
		resolved := result.GetSingleResolvedPresentation()
		if resolved == nil {
			return fmt.Errorf("no resolved presentation found for read operation"), 0
		}

		data.GroupPolicyDefinitionValueID = types.StringValue(result.DefinitionValueInstanceID)
		data.PresentationID = types.StringValue(resolved.TemplateID) // Keep template ID for binding
		data.ID = types.StringValue(resolved.InstanceID)
		tflog.Debug(ctx, fmt.Sprintf("[TEXT_VALUE] READ operation - using instance IDs (definitionValueInstanceID: %s, presentationValueInstanceID: %s)",
			result.DefinitionValueInstanceID, resolved.InstanceID))

	default:
		return fmt.Errorf("unsupported crud operation '%s' - must be 'create', 'read', or 'update'", operation), 0
	}

	return nil, 0
}
