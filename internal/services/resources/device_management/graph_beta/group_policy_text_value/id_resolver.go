package graphBetaGroupPolicyTextValue

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// groupPolicyIDResolver orchestrates the complex ID resolution process required by Microsoft Graph's
// hierarchical Group Policy architecture. This is necessary because:
//
// 1. Users author policies using human-readable names (policy_name, class_type) in HCL
// 2. Microsoft Graph requires specific GUIDs at different levels of the hierarchy:
//   - Template IDs: Reference the policy definition and presentation schemas
//   - Instance IDs: Reference actual configured values within a specific group policy configuration
//
// 3. Different CRUD operations require different ID types:
//
//   - CREATE: Uses template IDs to create new instances via updateDefinitionValues API
//
//   - READ/UPDATE: Requires instance IDs to access existing configured values
//
//     4. The Graph API hierarchy is: Configuration → DefinitionValue (instance) → PresentationValue (instance)
//     But policy templates are: Definition (template) → Presentation (template)
//
// This orchestrator eliminates the complexity of this dual-ID system from the CRUD operations,
// providing a single interface that handles the appropriate resolution strategy based on operation type.
// Without this abstraction, each CRUD operation would need to understand and implement the
// template-vs-instance ID resolution logic independently.
func groupPolicyIDResolver(ctx context.Context, data *GroupPolicyTextValueResourceModel, client *msgraphbetasdk.GraphServiceClient, operation string) error {
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] GroupPolicyIDResolver: Starting %s operation", operation))

	// Check if we have the required fields for lookup
	if data.PolicyName.IsNull() || data.PolicyName.IsUnknown() ||
		data.ClassType.IsNull() || data.ClassType.IsUnknown() ||
		data.CategoryPath.IsNull() || data.CategoryPath.IsUnknown() {
		return fmt.Errorf("provide policy_name, class_type, and category_path for auto-discovery")
	}

	// Get presentation index (default to 0)
	presentationIndex := int64(0)
	if !data.PresentationIndex.IsNull() && !data.PresentationIndex.IsUnknown() {
		presentationIndex = data.PresentationIndex.ValueInt64()
	}

	// Step 1: Always resolve policy name to template IDs
	definitionTemplateID, presentationTemplateID, err := resolvePolicyNameToTemplateIDs(
		ctx,
		data.PolicyName.ValueString(),
		data.ClassType.ValueString(),
		data.CategoryPath.ValueString(),
		presentationIndex,
		client,
	)

	if err != nil {
		return err
	}

	switch operation {
	case "create":
		// For creation, store template IDs only
		data.GroupPolicyDefinitionValueID = types.StringValue(definitionTemplateID)
		data.PresentationID = types.StringValue(presentationTemplateID)
		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] GroupPolicyIDResolver: CREATE operation - using template IDs (definitionTemplateID: %s, presentationTemplateID: %s)", definitionTemplateID, presentationTemplateID))
	case "update":
		// For update, we need both template IDs (for bindings) and instance IDs (for the update)
		definitionValueInstanceID, presentationValueInstanceID, err := resolveTemplateIDsToInstanceIDs(
			ctx,
			data.GroupPolicyConfigurationID.ValueString(),
			definitionTemplateID,
			presentationTemplateID,
			client,
		)

		if err != nil {
			return fmt.Errorf("failed to resolve template IDs to instance IDs: %w", err)
		}

		// Store template IDs for bindings, but we'll pass instance IDs via a different mechanism
		data.GroupPolicyDefinitionValueID = types.StringValue(definitionTemplateID) // Use template ID for binding
		data.PresentationID = types.StringValue(presentationTemplateID)             // Keep template ID for binding
		data.ID = types.StringValue(presentationValueInstanceID)                    // This is the presentation value instance ID

		// Store instance IDs in additional data for construct to use
		if data.AdditionalData == nil {
			data.AdditionalData = make(map[string]interface{})
		}
		data.AdditionalData["definitionValueInstanceID"] = definitionValueInstanceID
		data.AdditionalData["presentationValueInstanceID"] = presentationValueInstanceID

		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] GroupPolicyIDResolver: UPDATE operation - using template IDs for bindings (definitionTemplateID: %s, presentationTemplateID: %s) and instance IDs for update (definitionValueInstanceID: %s, presentationValueInstanceID: %s)", definitionTemplateID, presentationTemplateID, definitionValueInstanceID, presentationValueInstanceID))
	default:
		// For read, resolve template IDs to instance IDs
		definitionValueInstanceID, presentationValueInstanceID, err := resolveTemplateIDsToInstanceIDs(
			ctx,
			data.GroupPolicyConfigurationID.ValueString(),
			definitionTemplateID,
			presentationTemplateID,
			client,
		)

		if err != nil {
			return fmt.Errorf("failed to resolve template IDs to instance IDs: %w", err)
		}

		// Store instance IDs
		data.GroupPolicyDefinitionValueID = types.StringValue(definitionValueInstanceID)
		data.PresentationID = types.StringValue(presentationTemplateID) // Keep template ID for binding
		data.ID = types.StringValue(presentationValueInstanceID)
		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] GroupPolicyIDResolver: %s operation - using instance IDs (definitionValueInstanceID: %s, presentationValueInstanceID: %s)", operation, definitionValueInstanceID, presentationValueInstanceID))
	}

	return nil
}

// resolvePolicyNameToTemplateIDs bridges the gap between human-readable policy authoring and Microsoft Graph's
// GUID-based template system. This is necessary because:
//
//  1. Terraform users should author policies using intuitive names like "AD attribute containing Personal Site URL"
//     rather than memorizing opaque GUIDs like "a82c8307-85a0-499e-a9f5-ab8974337b65"
//
//  2. Microsoft Graph stores policy definitions and presentations as separate template objects, each with
//     their own GUID, requiring two separate API calls to resolve a single policy reference
//
//  3. The presentation selection logic (presentationIndex) allows users to specify which UI element to use
//     when a policy definition has multiple presentation options (text boxes, dropdowns, etc.)
//
// Without this abstraction, users would need to manually look up template GUIDs and understand the
// definition-to-presentation relationship, making the Terraform configuration brittle and user-unfriendly.
func resolvePolicyNameToTemplateIDs(
	ctx context.Context,
	policyName, classType, categoryPath string,
	presentationIndex int64,
	client *msgraphbetasdk.GraphServiceClient,
) (definitionTemplateID, presentationTemplateID string, err error) {

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Resolving policy name '%s' (classType='%s', categoryPath='%s', presentationIndex=%d) to template IDs",
		policyName, classType, categoryPath, presentationIndex))

	// Step 1: Look up the definition template ID using category path for precise matching
	definitionTemplateID, err = LookupDefinitionID(ctx, client, policyName, classType, categoryPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to find definition template: %w", err)
	}

	// Step 2: Look up the presentation template ID for single text presentations
	presentationTemplateID, err = LookupPresentationID(ctx, client, definitionTemplateID, "#microsoft.graph.groupPolicyPresentationTextBox")
	if err != nil {
		return "", "", fmt.Errorf("failed to find presentation template: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Successfully resolved policy name to template IDs - policyName='%s', categoryPath='%s', definitionTemplateID='%s', presentationTemplateID='%s'",
		policyName, categoryPath, definitionTemplateID, presentationTemplateID))

	return definitionTemplateID, presentationTemplateID, nil
}

// resolveTemplateIDsToInstanceIDs navigates from policy templates to their actual configured instances
// within a specific Group Policy Configuration. This is necessary because:
//
// 1. Template IDs represent the "schema" or "blueprint" of a policy (what CAN be configured)
// 2. Instance IDs represent actual configured values within a specific configuration (what IS configured)
// 3. Microsoft Graph's CRUD operations for reading/updating values require instance IDs, not template IDs
//
// 4. The API architecture creates a many-to-many relationship:
//   - One template can have multiple instances across different configurations
//   - One configuration can have multiple instances of the same template (edge case)
//
// 5. The search process involves:
//   - Querying all definition values in the configuration (instances)
//   - Filtering by the template definition ID to find matching instances
//   - Locating the specific presentation value within the matching definition value
//
// Without this resolution, READ/UPDATE operations would fail because they cannot directly
// access configured values using only template IDs - they need the actual instance IDs
// that were created when the policy was first configured.
func resolveTemplateIDsToInstanceIDs(ctx context.Context, configID, definitionTemplateID, presentationTemplateID string, client *msgraphbetasdk.GraphServiceClient) (definitionValueInstanceID, presentationValueInstanceID string, err error) {
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Resolving template IDs to instance IDs in config '%s' (definitionTemplateID: %s)", configID, definitionTemplateID))
	// Get all definition values for the configuration with expanded definition data
	requestConfig := &devicemanagement.GroupPolicyConfigurationsItemDefinitionValuesRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.GroupPolicyConfigurationsItemDefinitionValuesRequestBuilderGetQueryParameters{
			Expand: []string{"definition($select=id,classType,displayName,policyType,hasRelatedDefinitions,version,minUserCspVersion,minDeviceCspVersion)"},
		},
	}

	definitionValues, err := client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(configID).
		DefinitionValues().
		Get(ctx, requestConfig)

	if err != nil {
		return "", "", fmt.Errorf("failed to get definition values: %w", err)
	}

	// Check if definitionValues response is valid
	if definitionValues == nil {
		return "", "", fmt.Errorf("received nil definition values response")
	}

	values := definitionValues.GetValue()
	if values == nil {
		return "", "", fmt.Errorf("no definition values found in configuration")
	}

	// Find the definition value that matches our template definition ID
	for _, defValue := range values {
		if defValue == nil {
			continue
		}

		definition := defValue.GetDefinition()
		if definition == nil {
			continue
		}

		definitionID := definition.GetId()
		if definitionID != nil && *definitionID == definitionTemplateID {
			// Found matching definition value - now get its presentation values
			defValueID := defValue.GetId()
			if defValueID == nil {
				continue
			}

			presentationValues, err := client.
				DeviceManagement().
				GroupPolicyConfigurations().
				ByGroupPolicyConfigurationId(configID).
				DefinitionValues().
				ByGroupPolicyDefinitionValueId(*defValueID).
				PresentationValues().
				Get(ctx, nil)

			if err != nil {
				continue
			}

			if presentationValues == nil {
				continue
			}

			presValues := presentationValues.GetValue()
			if presValues == nil {
				continue
			}

			// Find the text presentation value
			for _, presValue := range presValues {
				if presValue == nil {
					continue
				}

				odataType := presValue.GetOdataType()
				if odataType != nil && *odataType == "#microsoft.graph.groupPolicyPresentationValueText" {
					presValueID := presValue.GetId()
					if presValueID != nil {
						return *defValueID, *presValueID, nil
					}
				}
			}
		}
	}

	return "", "", fmt.Errorf("no existing definition value found for definition template ID: %s", definitionTemplateID)
}
