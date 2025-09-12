package graphBetaGroupPolicyTextValue

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
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
func GroupPolicyIDResolver(ctx context.Context, data *GroupPolicyTextValueResourceModel, client *msgraphbetasdk.GraphServiceClient, operation string) (error, int) {
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] GroupPolicyIDResolver: Starting %s operation", operation))

	// Check if we have the required fields for lookup
	if data.PolicyName.IsNull() || data.PolicyName.IsUnknown() ||
		data.ClassType.IsNull() || data.ClassType.IsUnknown() ||
		data.CategoryPath.IsNull() || data.CategoryPath.IsUnknown() {
		return fmt.Errorf("provide policy_name, class_type, and category_path for auto-discovery"), 0
	}

	// Step 1: Resolve policy name to definition template ID
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Resolving policy name '%s' (classType='%s', categoryPath='%s') to definition template ID",
		data.PolicyName.ValueString(), data.ClassType.ValueString(), data.CategoryPath.ValueString()))

	definitionTemplateID, err := groupPolicyNameResolver(ctx, client, data.PolicyName.ValueString(), data.ClassType.ValueString(), data.CategoryPath.ValueString())
	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		return fmt.Errorf("failed to find definition template: %w", err), errorInfo.StatusCode
	}

	// Step 2: Use the first presentation available (any OData type)
	// The system will automatically handle all presentations during CRUD operations
	presentationTemplateID, err := resolveFirstPresentation(ctx, client, definitionTemplateID)
	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		return fmt.Errorf("failed to find presentation template: %w", err), errorInfo.StatusCode
	}

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Successfully resolved policy name to template IDs - policyName='%s', categoryPath='%s', definitionTemplateID='%s', presentationTemplateID='%s'",
		data.PolicyName.ValueString(), data.CategoryPath.ValueString(), definitionTemplateID, presentationTemplateID))

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
			errorInfo := errors.GraphError(ctx, err)
			return fmt.Errorf("failed to resolve template IDs to instance IDs: %w", err), errorInfo.StatusCode
		}

		// Store template IDs for bindings, but we'll pass instance IDs via a different mechanism
		data.GroupPolicyDefinitionValueID = types.StringValue(definitionTemplateID) // Use template ID for binding
		data.PresentationID = types.StringValue(presentationTemplateID)             // Keep template ID for binding
		data.ID = types.StringValue(presentationValueInstanceID)                    // This is the presentation value instance ID

		// Store instance IDs in additional data for construct to use
		if data.AdditionalData == nil {
			data.AdditionalData = make(map[string]any)
		}
		data.AdditionalData["definitionValueInstanceID"] = definitionValueInstanceID
		data.AdditionalData["presentationValueInstanceID"] = presentationValueInstanceID

		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] GroupPolicyIDResolver: UPDATE operation - using template IDs for bindings (definitionTemplateID: %s, presentationTemplateID: %s) and instance IDs for update (definitionValueInstanceID: %s, presentationValueInstanceID: %s)", definitionTemplateID, presentationTemplateID, definitionValueInstanceID, presentationValueInstanceID))
	case "read":
		// For read, resolve template IDs to instance IDs
		definitionValueInstanceID, presentationValueInstanceID, err := resolveTemplateIDsToInstanceIDs(
			ctx,
			data.GroupPolicyConfigurationID.ValueString(),
			definitionTemplateID,
			presentationTemplateID,
			client,
		)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			// Handle 500 errors during read scenarios - resource likely deleted
			if errorInfo.StatusCode == 500 && operation == "read" {
				tflog.Warn(ctx, "500 error during read operation indicates resource has been deleted from policy configuration", map[string]any{
					"status_code":   errorInfo.StatusCode,
					"error_code":    errorInfo.ErrorCode,
					"error_message": errorInfo.ErrorMessage,
					"request_id":    errorInfo.RequestID,
				})

				errorMsg := "resource no longer exists in policy configuration (HTTP 500)"
				if errorInfo.ErrorMessage != "" {
					errorMsg += fmt.Sprintf(": %s", errorInfo.ErrorMessage)
				}

				return fmt.Errorf("%s: %w", errorMsg, err), 500
			}
			return fmt.Errorf("failed to resolve template IDs to instance IDs: %w", err), errorInfo.StatusCode
		}

		data.GroupPolicyDefinitionValueID = types.StringValue(definitionValueInstanceID)
		data.PresentationID = types.StringValue(presentationTemplateID) // Keep template ID for binding
		data.ID = types.StringValue(presentationValueInstanceID)
		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] GroupPolicyIDResolver: READ operation - using instance IDs (definitionValueInstanceID: %s, presentationValueInstanceID: %s)", definitionValueInstanceID, presentationValueInstanceID))
	default:
		return fmt.Errorf("unsupported crud operation '%s' - must be 'create', 'read', or 'update'", operation), 0
	}

	return nil, 0
}

// groupPolicyNameResolver finds the definition ID based on display name, class type, and category path
// If multiple matches are found, an error is returned
func groupPolicyNameResolver(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName, classType, categoryPath string) (string, error) {

	filterQuery := fmt.Sprintf("displayName eq '%s' and classType eq '%s' and categoryPath eq '%s'", displayName, classType, categoryPath)
	tflog.Debug(ctx, fmt.Sprintf(" Resolving supplied group policy metadata to resolve group policy definition id using odata filter: %s", filterQuery))

	definitions, err := client.
		DeviceManagement().
		GroupPolicyDefinitions().
		Get(ctx, &devicemanagement.GroupPolicyDefinitionsRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.GroupPolicyDefinitionsRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName", "classType", "categoryPath"},
				Filter: &[]string{filterQuery}[0],
			},
		})

	if err != nil {
		return "", fmt.Errorf("failed to fetch group policy definitions: %w", err)
	}

	// Use PageIterator there are over 7000 group policy definitions
	var allResults []graphmodels.GroupPolicyDefinitionable

	pageIterator, err := graphcore.NewPageIterator[graphmodels.GroupPolicyDefinitionable](
		definitions,
		client.GetAdapter(),
		graphmodels.CreateGroupPolicyDefinitionCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return "", fmt.Errorf("failed to create page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.GroupPolicyDefinitionable) bool {
		if item != nil {
			allResults = append(allResults, item)
		}
		return true
	})

	if err != nil {
		return "", fmt.Errorf("failed to iterate pages: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf(" Found %d definitions with odata filter with paginated GET", len(allResults)))

	if len(allResults) == 0 {
		return "", fmt.Errorf("no group policy definition found with displayName='%s', classType='%s', categoryPath='%s'", displayName, classType, categoryPath)
	}

	// If we have multiple results after complete filtering, this is an error
	if len(allResults) > 1 {

		var matchDetails []string
		for i, def := range allResults {
			if def != nil && def.GetId() != nil {
				defID := *def.GetId()
				defDisplayName := ""
				if def.GetDisplayName() != nil {
					defDisplayName = *def.GetDisplayName()
				}
				defClassType := ""
				if def.GetClassType() != nil {
					defClassType = def.GetClassType().String()
				}
				defCategoryPath := ""
				if def.GetCategoryPath() != nil {
					defCategoryPath = *def.GetCategoryPath()
				}
				matchDetails = append(matchDetails, fmt.Sprintf("Match %d: ID=%s, DisplayName='%s', ClassType='%s', CategoryPath='%s'", i+1, defID, defDisplayName, defClassType, defCategoryPath))
			}
		}

		return "", fmt.Errorf("group policy name resolution failed to resolve to a singular definition, got %d matches: %s",
			len(allResults), strings.Join(matchDetails, "; "))
	}

	firstResult := allResults[0]
	if firstResult == nil || firstResult.GetId() == nil {
		return "", fmt.Errorf("invalid group policy definition returned from server")
	}

	definitionID := *firstResult.GetId()
	tflog.Debug(ctx, fmt.Sprintf(" ✅ Found single definition ID: %s", definitionID))

	return definitionID, nil
}

// resolveFirstPresentation finds the first available presentation for a policy definition.
// This simplifies the user experience by automatically selecting the first valid presentation
// without requiring users to understand presentation indices or types.
//
// The function searches through all presentations for a definition and returns the first one that
// has a valid OData type, supporting all presentation types (TextBox, MultiTextBox, Text, CheckBox, etc.).
func resolveFirstPresentation(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	definitionTemplateID string,
) (presentationTemplateID string, err error) {

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Finding first presentation for definitionID='%s'", definitionTemplateID))

	// Get all presentations for this definition
	presentations, err := client.
		DeviceManagement().
		GroupPolicyDefinitions().
		ByGroupPolicyDefinitionId(definitionTemplateID).
		Presentations().
		Get(ctx, nil)

	if err != nil {
		return "", fmt.Errorf("failed to get presentations for definition %s: %w", definitionTemplateID, err)
	}

	if presentations == nil || presentations.GetValue() == nil || len(presentations.GetValue()) == 0 {
		return "", fmt.Errorf("no presentations found for definition %s", definitionTemplateID)
	}

	presentationList := presentations.GetValue()
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found %d presentations for definition", len(presentationList)))

	// Return the first presentation with a valid OData type and ID
	for i, presentation := range presentationList {
		if presentation == nil {
			continue
		}

		odataType := presentation.GetOdataType()
		if odataType == nil {
			continue
		}

		presentationID := presentation.GetId()
		if presentationID == nil {
			continue
		}

		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found presentation %d: type=%s, id=%s", i, *odataType, *presentationID))
		return *presentationID, nil
	}

	return "", fmt.Errorf("no valid presentations found for definition %s", definitionTemplateID)
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

			// Find the first presentation value (any type)
			for _, presValue := range presValues {
				if presValue == nil {
					continue
				}

				odataType := presValue.GetOdataType()
				if odataType != nil {
					presValueID := presValue.GetId()
					if presValueID != nil {
						tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found presentation value instance ID: %s (type: %s)", *presValueID, *odataType))
						return *defValueID, *presValueID, nil
					}
				}
			}
		}
	}

	return "", "", fmt.Errorf("no existing definition value found for definition template ID: %s", definitionTemplateID)
}
