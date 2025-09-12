package graphBetaGroupPolicyBooleanValue

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
func GroupPolicyIDResolver(ctx context.Context, data *GroupPolicyBooleanValueResourceModel, client *msgraphbetasdk.GraphServiceClient, operation string) (error, int) {
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

	// Step 2: Resolve ALL presentations for this definition (extends the single presentation logic)
	// This keeps the same API call structure but collects all presentations instead of just the first
	presentationTemplateIDs, err := resolveAllPresentations(ctx, client, definitionTemplateID)
	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		return fmt.Errorf("failed to find presentation templates: %w", err), errorInfo.StatusCode
	}

	if len(presentationTemplateIDs) == 0 {
		return fmt.Errorf("no presentations found for definition %s", definitionTemplateID), 0
	}

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Successfully resolved policy name to template IDs - policyName='%s', categoryPath='%s', definitionTemplateID='%s', foundPresentations=%d",
		data.PolicyName.ValueString(), data.CategoryPath.ValueString(), definitionTemplateID, len(presentationTemplateIDs)))

	switch operation {
	case "create":
		// For creation, store template IDs and resolved presentations (extends single presentation logic)
		data.GroupPolicyDefinitionValueID = types.StringValue(definitionTemplateID)
		data.ID = types.StringValue(definitionTemplateID) // Use definition ID as the main ID

		// Store resolved presentations for construct to use
		if data.AdditionalData == nil {
			data.AdditionalData = make(map[string]any)
		}

		var resolvedPresentations []ResolvedPresentation
		for i, templateID := range presentationTemplateIDs {
			resolvedPresentations = append(resolvedPresentations, ResolvedPresentation{
				TemplateID: templateID,
				Index:      i,
			})
		}
		data.AdditionalData["resolvedPresentations"] = resolvedPresentations

		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] GroupPolicyIDResolver: CREATE operation - using template IDs (definitionTemplateID: %s, presentations: %d)", definitionTemplateID, len(presentationTemplateIDs)))

	case "update":
		// For update, extend the single resolution logic to handle multiple presentations
		// This uses the same API call pattern but resolves multiple template-to-instance mappings
		definitionValueInstanceID, resolvedPresentations, err := resolveMultipleTemplateIDsToInstanceIDs(
			ctx,
			data.GroupPolicyConfigurationID.ValueString(),
			definitionTemplateID,
			presentationTemplateIDs,
			client,
		)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			return fmt.Errorf("failed to resolve template IDs to instance IDs: %w", err), errorInfo.StatusCode
		}

		// Store template IDs for bindings and instance IDs for updates
		data.GroupPolicyDefinitionValueID = types.StringValue(definitionTemplateID) // Use template ID for binding
		data.ID = types.StringValue(definitionValueInstanceID)                      // This is the definition value instance ID

		// Store instance IDs in additional data for construct to use
		if data.AdditionalData == nil {
			data.AdditionalData = make(map[string]any)
		}
		data.AdditionalData["definitionValueInstanceID"] = definitionValueInstanceID
		data.AdditionalData["resolvedPresentations"] = resolvedPresentations

		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] GroupPolicyIDResolver: UPDATE operation - using template IDs for bindings and instance IDs for update (definitionValueInstanceID: %s, presentations: %d)", definitionValueInstanceID, len(resolvedPresentations)))

	case "read":
		// For read, extend the single resolution logic to handle multiple presentations
		// This uses the same API call pattern but resolves multiple template-to-instance mappings
		definitionValueInstanceID, resolvedPresentations, err := resolveMultipleTemplateIDsToInstanceIDs(
			ctx,
			data.GroupPolicyConfigurationID.ValueString(),
			definitionTemplateID,
			presentationTemplateIDs,
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

				// Build concise error message based on what API actually returns
				errorMsg := "resource no longer exists in policy configuration (HTTP 500)"
				if errorInfo.ErrorMessage != "" {
					// Use the API's error message as it's usually descriptive
					errorMsg += fmt.Sprintf(": %s", errorInfo.ErrorMessage)
				}

				return fmt.Errorf("%s: %w", errorMsg, err), 500
			}
			return fmt.Errorf("failed to resolve template IDs to instance IDs: %w", err), errorInfo.StatusCode
		}

		// Store instance IDs
		data.GroupPolicyDefinitionValueID = types.StringValue(definitionValueInstanceID)
		data.ID = types.StringValue(definitionValueInstanceID)

		// Store resolved presentations for state mapping
		if data.AdditionalData == nil {
			data.AdditionalData = make(map[string]any)
		}
		data.AdditionalData["resolvedPresentations"] = resolvedPresentations

		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] GroupPolicyIDResolver: READ operation - using instance IDs (definitionValueInstanceID: %s, presentations: %d)", definitionValueInstanceID, len(resolvedPresentations)))

	default:
		return fmt.Errorf("unsupported crud operation '%s' - must be 'create', 'read', or 'update'", operation), 0
	}

	return nil, 0
}

// groupPolicyNameResolver finds the definition ID based on display name, class type, and category path
// If multiple matches are found, an error is returned
func groupPolicyNameResolver(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName, classType, categoryPath string) (string, error) {

	// Trim whitespace from input parameters to handle API inconsistencies
	displayName = strings.TrimSpace(displayName)
	classType = strings.TrimSpace(classType)
	categoryPath = strings.TrimSpace(categoryPath)

	tflog.Debug(ctx, "[LOOKUP] Starting policy name resolution with exact criteria:")
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] - displayName: '%s' (length: %d)", displayName, len(displayName)))
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] - classType: '%s'", classType))
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] - categoryPath: '%s' (length: %d)", categoryPath, len(categoryPath)))

	filterQuery := fmt.Sprintf("displayName eq '%s' and classType eq '%s' and categoryPath eq '%s'", displayName, classType, categoryPath)
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] OData filter query: %s", filterQuery))

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

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found %d definitions with OData filter", len(allResults)))

	if len(allResults) == 0 {
		tflog.Debug(ctx, "[LOOKUP] No exact match found - trying flexible matching for whitespace differences...")

		// Try a broader search to find policies that might match with whitespace differences
		broadFilterQuery := fmt.Sprintf("contains(displayName, '%s') and classType eq '%s' and categoryPath eq '%s'", displayName, classType, categoryPath)
		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Flexible search filter: %s", broadFilterQuery))

		broadDefinitions, err := client.
			DeviceManagement().
			GroupPolicyDefinitions().
			Get(ctx, &devicemanagement.GroupPolicyDefinitionsRequestBuilderGetRequestConfiguration{
				QueryParameters: &devicemanagement.GroupPolicyDefinitionsRequestBuilderGetQueryParameters{
					Select: []string{"id", "displayName", "classType", "categoryPath"},
					Filter: &[]string{broadFilterQuery}[0],
					Top:    &[]int32{20}[0], // Limit to 20 results
				},
			})

		if err == nil && broadDefinitions != nil {
			broadPageIterator, err := graphcore.NewPageIterator[graphmodels.GroupPolicyDefinitionable](
				broadDefinitions,
				client.GetAdapter(),
				graphmodels.CreateGroupPolicyDefinitionCollectionResponseFromDiscriminatorValue,
			)

			if err == nil {
				var broadResults []graphmodels.GroupPolicyDefinitionable
				broadPageIterator.Iterate(ctx, func(item graphmodels.GroupPolicyDefinitionable) bool {
					if item != nil {
						broadResults = append(broadResults, item)
					}
					return true
				})

				tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found %d policies with flexible search:", len(broadResults)))

				// Check if any of the broad results match exactly when whitespace is normalized
				var exactMatches []graphmodels.GroupPolicyDefinitionable
				for i, def := range broadResults {
					if def != nil {
						defDisplayName := ""
						if def.GetDisplayName() != nil {
							defDisplayName = strings.TrimSpace(*def.GetDisplayName())
						}
						defClassType := ""
						if def.GetClassType() != nil {
							defClassType = strings.TrimSpace(def.GetClassType().String())
						}
						defCategoryPath := ""
						if def.GetCategoryPath() != nil {
							defCategoryPath = strings.TrimSpace(*def.GetCategoryPath())
						}

						tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] %d. DisplayName='%s', ClassType='%s', CategoryPath='%s'", i+1, defDisplayName, defClassType, defCategoryPath))

						// Check for exact match after trimming whitespace
						if defDisplayName == displayName && defClassType == classType && defCategoryPath == categoryPath {
							tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ✅ Found exact match after whitespace normalization: %s", *def.GetId()))
							exactMatches = append(exactMatches, def)
						}
					}
				}

				// If we found exact matches after whitespace normalization, use them
				if len(exactMatches) == 1 {
					tflog.Debug(ctx, "[LOOKUP] Using whitespace-normalized exact match")
					allResults = exactMatches
				} else if len(exactMatches) > 1 {
					tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found %d exact matches after whitespace normalization - this is ambiguous", len(exactMatches)))
					allResults = exactMatches
				}
			}
		}

		// If still no matches found after flexible search, return error
		if len(allResults) == 0 {
			return "", fmt.Errorf("no group policy definition found with displayName='%s', classType='%s', categoryPath='%s'", displayName, classType, categoryPath)
		}
	}

	// Log all results found for debugging
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Exact match results found: %d", len(allResults)))
	for i, def := range allResults {
		if def != nil {
			defID := ""
			if def.GetId() != nil {
				defID = *def.GetId()
			}
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
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Result %d: ID=%s, DisplayName='%s', ClassType='%s', CategoryPath='%s'", i+1, defID, defDisplayName, defClassType, defCategoryPath))
		}
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
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ✅ Found single definition ID: %s", definitionID))

	return definitionID, nil
}

// resolveAllPresentations extends resolveFirstPresentation to find ALL presentations for a policy definition.
// This keeps the same API call structure but collects all presentations instead of just returning the first.
// This supports policies with multiple presentations (like boolean values with multiple checkboxes).
func resolveAllPresentations(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	definitionTemplateID string,
) (presentationTemplateIDs []string, err error) {

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Finding all presentations for definitionID='%s'", definitionTemplateID))

	// Get all presentations for this definition (SAME API CALL as resolveFirstPresentation)
	presentations, err := client.
		DeviceManagement().
		GroupPolicyDefinitions().
		ByGroupPolicyDefinitionId(definitionTemplateID).
		Presentations().
		Get(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to get presentations for definition %s: %w", definitionTemplateID, err)
	}

	if presentations == nil || presentations.GetValue() == nil || len(presentations.GetValue()) == 0 {
		return nil, fmt.Errorf("no presentations found for definition %s", definitionTemplateID)
	}

	presentationList := presentations.GetValue()
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found %d presentations for definition", len(presentationList)))

	var allPresentationIDs []string

	// Collect only CHECKBOX presentation IDs for boolean values (filter out text and other types)
	for i, presentation := range presentationList {
		if presentation == nil {
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Skipping nil presentation at index %d", i))
			continue
		}

		presentationID := presentation.GetId()
		if presentationID == nil || *presentationID == "" {
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Skipping presentation %d with missing ID", i))
			continue
		}

		odataType := presentation.GetOdataType()
		odataTypeStr := ""
		if odataType != nil {
			odataTypeStr = *odataType
		}

		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found presentation %d: ID=%s, type=%s", i, *presentationID, odataTypeStr))

		// Only include checkbox presentations for boolean values - filter out text and other types
		if odataType != nil && *odataType == "#microsoft.graph.groupPolicyPresentationCheckBox" {
			allPresentationIDs = append(allPresentationIDs, *presentationID)
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ✅ Added CHECKBOX presentation template ID: %s", *presentationID))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ⏭️  Skipping non-checkbox presentation: %s (type: %s)", *presentationID, odataTypeStr))
		}
	}

	if len(allPresentationIDs) == 0 {
		return nil, fmt.Errorf("no valid presentations found for definition %s", definitionTemplateID)
	}

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Successfully resolved %d presentation template IDs", len(allPresentationIDs)))
	return allPresentationIDs, nil
}

// resolveMultipleTemplateIDsToInstanceIDs extends resolveTemplateIDsToInstanceIDs to handle multiple presentations.
// This keeps the same API call structure but maps multiple template IDs to their instance IDs.
// It uses the exact same API calls as the single resolution but handles multiple presentation mappings.
func resolveMultipleTemplateIDsToInstanceIDs(ctx context.Context, configID, definitionTemplateID string, presentationTemplateIDs []string, client *msgraphbetasdk.GraphServiceClient) (definitionValueInstanceID string, resolvedPresentations []ResolvedPresentation, err error) {
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Resolving multiple template IDs to instance IDs in config '%s' (definitionTemplateID: %s, presentations: %d)", configID, definitionTemplateID, len(presentationTemplateIDs)))

	// Get all definition values for the configuration with expanded definition data (SAME API CALL as single resolution)
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
		return "", nil, fmt.Errorf("failed to get definition values: %w", err)
	}

	// Check if definitionValues response is valid (SAME LOGIC as single resolution)
	if definitionValues == nil {
		return "", nil, fmt.Errorf("received nil definition values response")
	}

	values := definitionValues.GetValue()
	if values == nil {
		return "", nil, fmt.Errorf("no definition values found in configuration")
	}

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found %d definition values in configuration '%s'", len(values), configID))

	// Find the definition value that matches our template definition ID (SAME LOGIC as single resolution)
	for i, defValue := range values {
		if defValue == nil {
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Skipping nil definition value at index %d", i))
			continue
		}

		definition := defValue.GetDefinition()
		if definition == nil {
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Skipping definition value %d with nil definition", i))
			continue
		}

		definitionID := definition.GetId()
		defValueInstanceID := ""
		if defValue.GetId() != nil {
			defValueInstanceID = *defValue.GetId()
		}

		if definitionID != nil {
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Checking definition value %d: definitionTemplateID='%s', instanceID='%s'", i, *definitionID, defValueInstanceID))
		}

		if definitionID != nil && *definitionID == definitionTemplateID {
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ✅ Found matching definition value: instanceID='%s'", defValueInstanceID))
			// Found matching definition value - now get its presentation values (SAME API CALL as single resolution)
			defValueID := defValue.GetId()
			if defValueID == nil {
				continue
			}

			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Getting presentation values for definition value instance ID: %s", *defValueID))

			presentationValues, err := client.
				DeviceManagement().
				GroupPolicyConfigurations().
				ByGroupPolicyConfigurationId(configID).
				DefinitionValues().
				ByGroupPolicyDefinitionValueId(*defValueID).
				PresentationValues().
				Get(ctx, nil)

			if err != nil {
				tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Failed to get presentation values: %v", err))
				continue
			}

			if presentationValues == nil {
				tflog.Debug(ctx, "[LOOKUP] Received nil presentation values response")
				continue
			}

			presValues := presentationValues.GetValue()
			if presValues == nil {
				tflog.Debug(ctx, "[LOOKUP] No presentation values found in response")
				continue
			}

			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found %d presentation values for definition instance", len(presValues)))

			// First, log all presentation reference IDs found to understand what we're working with
			tflog.Debug(ctx, "[LOOKUP] All presentation reference IDs found:")
			hasAnyPresentationRefs := false
			for i, presValue := range presValues {
				if presValue != nil {
					presentation := presValue.GetPresentation()
					if presentation != nil && presentation.GetId() != nil {
						hasAnyPresentationRefs = true
						tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] - Presentation %d: ref_id='%s', instance_id='%s'", i, *presentation.GetId(), func() string {
							if presValue.GetId() != nil {
								return *presValue.GetId()
							}
							return "nil"
						}()))
					} else {
						tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] - Presentation %d: ref_id=nil, instance_id='%s'", i, func() string {
							if presValue.GetId() != nil {
								return *presValue.GetId()
							}
							return "nil"
						}()))
					}
				}
			}

			// Map multiple presentation template IDs to their instance IDs
			var resolved []ResolvedPresentation
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Mapping %d presentation template IDs to instance IDs", len(presentationTemplateIDs)))

			if hasAnyPresentationRefs {
				// Strategy 1: Use presentation references (when available)
				tflog.Debug(ctx, "[LOOKUP] Using presentation reference mapping strategy")

				for i, templateID := range presentationTemplateIDs {
					tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Looking for template ID '%s' (index %d)", templateID, i))

					// Find the presentation value that corresponds to this template
					found := false
					for _, presValue := range presValues {
						if presValue == nil {
							continue
						}

						// Get the presentation reference from the presentation value
						presentation := presValue.GetPresentation()
						if presentation == nil {
							continue
						}

						presentationRefID := presentation.GetId()
						if presentationRefID != nil {
							tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Comparing presentation ref ID '%s' with template ID '%s'", *presentationRefID, templateID))
							if *presentationRefID == templateID {
								presValueID := presValue.GetId()
								if presValueID != nil {
									resolved = append(resolved, ResolvedPresentation{
										TemplateID: templateID,
										InstanceID: *presValueID,
										Index:      i,
									})
									tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ✅ Found presentation value instance ID: %s for template: %s", *presValueID, templateID))
									found = true
									break
								}
							}
						}
					}

					if !found {
						tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ❌ No presentation value found for template ID '%s'", templateID))
					}
				}
			} else {
				// Strategy 2: Use positional mapping (fallback when presentation references are missing)
				tflog.Debug(ctx, "[LOOKUP] No presentation references found - using positional mapping strategy")

				// Map by position: presValues[0] -> templateIDs[0], presValues[1] -> templateIDs[1], etc.
				maxMappings := len(presentationTemplateIDs)
				if len(presValues) < maxMappings {
					maxMappings = len(presValues)
				}

				for i := 0; i < maxMappings; i++ {
					presValue := presValues[i]
					if presValue == nil {
						tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ❌ Presentation value at index %d is nil", i))
						continue
					}

					presValueID := presValue.GetId()
					if presValueID == nil {
						tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ❌ Presentation value at index %d has nil ID", i))
						continue
					}

					templateID := presentationTemplateIDs[i]
					resolved = append(resolved, ResolvedPresentation{
						TemplateID: templateID,
						InstanceID: *presValueID,
						Index:      i,
					})
					tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ✅ Positional mapping: template[%d]='%s' -> instance='%s'", i, templateID, *presValueID))
				}

				if len(presentationTemplateIDs) > len(presValues) {
					tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ⚠️ Warning: %d template IDs but only %d presentation values - some templates not mapped", len(presentationTemplateIDs), len(presValues)))
				}
			}

			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Successfully resolved %d out of %d presentation template IDs", len(resolved), len(presentationTemplateIDs)))

			if len(resolved) > 0 {
				return *defValueID, resolved, nil
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ❌ No matching definition value found for template ID '%s' in configuration '%s'", definitionTemplateID, configID))
	return "", nil, fmt.Errorf("no existing definition value found for definition template ID: %s", definitionTemplateID)
}
