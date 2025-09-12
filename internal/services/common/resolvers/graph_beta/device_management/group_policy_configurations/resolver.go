package groupPolicyConfigurations

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// ResolvedPresentation represents a mapping between template and instance IDs
type ResolvedPresentation struct {
	TemplateID string
	InstanceID string
	Index      int
}

// GroupPolicyIDResolver orchestrates the complex ID resolution process required by Microsoft Graph's
// hierarchical Group Policy architecture. This centralized resolver supports both single and multiple
// presentation value scenarios. This is necessary because:
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
func GroupPolicyIDResolver(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, operation string, policyName, classType, categoryPath, configID string, presentationFilter string) (*ResolverResult, error) {
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] GroupPolicyIDResolver: Starting %s operation", operation))

	if policyName == "" || classType == "" || categoryPath == "" {
		return nil, fmt.Errorf("provide policy_name, class_type, and category_path for auto-discovery")
	}

	// Step 1: Resolve policy name to definition template ID
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Resolving policy name '%s' (classType='%s', categoryPath='%s') to definition template ID",
		policyName, classType, categoryPath))

	definitionTemplateID, err := resolveGroupPolicyName(ctx, client, policyName, classType, categoryPath)
	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		return nil, fmt.Errorf("failed to find definition template: %w (status: %d)", err, errorInfo.StatusCode)
	}

	// Step 2: Resolve presentations based on filter
	var presentationTemplateIDs []string
	if presentationFilter != "" {
		// Multiple presentations with specific OData type filter
		presentationTemplateIDs, err = resolveFilteredPresentations(ctx, client, definitionTemplateID, presentationFilter)
	} else {
		// Single presentation (any OData type)
		singleID, err := resolveFirstPresentation(ctx, client, definitionTemplateID)
		if err == nil {
			presentationTemplateIDs = []string{singleID}
		}
	}

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		return nil, fmt.Errorf("failed to find presentation template(s): %w (status: %d)", err, errorInfo.StatusCode)
	}

	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Successfully resolved policy name to template IDs - policyName='%s', categoryPath='%s', definitionTemplateID='%s', presentationTemplateIDs=%v",
		policyName, categoryPath, definitionTemplateID, presentationTemplateIDs))

	result := &ResolverResult{
		DefinitionTemplateID:    definitionTemplateID,
		PresentationTemplateIDs: presentationTemplateIDs,
	}

	switch operation {
	case "create":
		// For creation, store template IDs only
		result.DefinitionValueInstanceID = definitionTemplateID // Use template for create
		tflog.Debug(ctx, "[RESOLVER] CREATE operation - using template IDs")

	case "update", "read":
		// For update/read, we need both template IDs (for bindings) and instance IDs (for the operation)
		if configID == "" {
			return nil, fmt.Errorf("configuration ID required for %s operation", operation)
		}

		definitionValueInstanceID, resolvedPresentations, err := resolveTemplateIDsToInstanceIDs(
			ctx, configID, definitionTemplateID, presentationTemplateIDs, client)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			tflog.Error(ctx, "[RESOLVER] Failed to resolve template IDs to instance IDs", map[string]any{
				"operation":        operation,
				"config_id":        configID,
				"definition_id":    definitionTemplateID,
				"presentation_ids": presentationTemplateIDs,
				"status_code":      errorInfo.StatusCode,
				"error_code":       errorInfo.ErrorCode,
				"error_message":    errorInfo.ErrorMessage,
				"request_id":       errorInfo.RequestID,
			})

			// Handle 500 errors during read scenarios - resource likely deleted
			if errorInfo.StatusCode == 500 && operation == "read" {
				tflog.Warn(ctx, "[RESOLVER] 500 error during read operation - resource appears to have been deleted from policy configuration", map[string]any{
					"operation":        operation,
					"config_id":        configID,
					"definition_id":    definitionTemplateID,
					"presentation_ids": presentationTemplateIDs,
					"status_code":      errorInfo.StatusCode,
					"error_code":       errorInfo.ErrorCode,
					"error_message":    errorInfo.ErrorMessage,
					"request_id":       errorInfo.RequestID,
					"action":           "triggering_state_removal",
				})

				// Return a specific error that indicates resource deletion for proper state handling
				return nil, fmt.Errorf("resource has been deleted from policy configuration and will be removed from state")
			}
			return nil, fmt.Errorf("failed to resolve template IDs to instance IDs: %w (status: %d)", err, errorInfo.StatusCode)
		}

		result.DefinitionValueInstanceID = definitionValueInstanceID
		result.ResolvedPresentations = resolvedPresentations
		tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] %s operation - resolved %d presentations", strings.ToUpper(operation), len(resolvedPresentations)))

	default:
		return nil, fmt.Errorf("unsupported crud operation '%s' - must be 'create', 'read', or 'update'", operation)
	}

	return result, nil
}

// ResolverResult contains all the resolved IDs needed for group policy operations
type ResolverResult struct {
	DefinitionTemplateID      string
	PresentationTemplateIDs   []string
	DefinitionValueInstanceID string
	ResolvedPresentations     []ResolvedPresentation
}

// GetSinglePresentationID returns the first presentation template ID for single-value resources
func (r *ResolverResult) GetSinglePresentationID() string {
	if len(r.PresentationTemplateIDs) > 0 {
		return r.PresentationTemplateIDs[0]
	}
	return ""
}

// GetSingleResolvedPresentation returns the first resolved presentation for single-value resources
func (r *ResolverResult) GetSingleResolvedPresentation() *ResolvedPresentation {
	if len(r.ResolvedPresentations) > 0 {
		return &r.ResolvedPresentations[0]
	}
	return nil
}

// resolveGroupPolicyName finds the definition ID based on display name, class type, and category path
// If multiple matches are found, an error is returned
func resolveGroupPolicyName(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName, classType, categoryPath string) (string, error) {
	// Trim whitespace from input parameters to handle API inconsistencies
	displayName = strings.TrimSpace(displayName)
	classType = strings.TrimSpace(classType)
	categoryPath = strings.TrimSpace(categoryPath)

	// Use flexible search with contains() and then normalize for exact matching
	flexibleFilter := fmt.Sprintf("contains(displayName, '%s') and classType eq '%s' and categoryPath eq '%s'", displayName, classType, categoryPath)
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Resolving group policy metadata using flexible OData filter: %s", flexibleFilter))

	definitions, err := client.
		DeviceManagement().
		GroupPolicyDefinitions().
		Get(ctx, &devicemanagement.GroupPolicyDefinitionsRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.GroupPolicyDefinitionsRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName", "classType", "categoryPath"},
				Filter: &[]string{flexibleFilter}[0],
			},
		})

	if err != nil {
		return "", fmt.Errorf("failed to fetch group policy definitions: %w", err)
	}

	// Use PageIterator - there are over 7000 group policy definitions
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

	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Found %d definitions with flexible search", len(allResults)))

	// Filter results for exact match with normalized whitespace
	var exactMatches []graphmodels.GroupPolicyDefinitionable
	for _, def := range allResults {
		if def != nil {
			defDisplayName := ""
			if def.GetDisplayName() != nil {
				defDisplayName = strings.TrimSpace(*def.GetDisplayName())
			}
			defClassType := ""
			if def.GetClassType() != nil {
				defClassType = def.GetClassType().String()
			}
			defCategoryPath := ""
			if def.GetCategoryPath() != nil {
				defCategoryPath = strings.TrimSpace(*def.GetCategoryPath())
			}

			if defDisplayName == displayName && defClassType == classType && defCategoryPath == categoryPath {
				exactMatches = append(exactMatches, def)
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Found %d exact matches after whitespace normalization", len(exactMatches)))

	if len(exactMatches) == 0 {
		return "", fmt.Errorf("no group policy definition found with displayName='%s', classType='%s', categoryPath='%s'", displayName, classType, categoryPath)
	}

	// If we have multiple results after complete filtering, this is an error
	if len(exactMatches) > 1 {
		var matchDetails []string
		for i, def := range exactMatches {
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
			len(exactMatches), strings.Join(matchDetails, "; "))
	}

	firstResult := exactMatches[0]
	if firstResult == nil || firstResult.GetId() == nil {
		return "", fmt.Errorf("invalid group policy definition returned from server")
	}

	definitionID := *firstResult.GetId()
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] ✅ Found single definition ID: %s", definitionID))

	return definitionID, nil
}

// resolveFirstPresentation finds the first available presentation for a policy definition.
// This simplifies the user experience by automatically selecting the first valid presentation
// without requiring users to understand presentation indices or types.
func resolveFirstPresentation(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, definitionTemplateID string) (string, error) {
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Finding first presentation for definitionID='%s'", definitionTemplateID))

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
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Found %d presentations for definition", len(presentationList)))

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

		tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Found presentation %d: type=%s, id=%s", i, *odataType, *presentationID))
		return *presentationID, nil
	}

	return "", fmt.Errorf("no valid presentations found for definition %s", definitionTemplateID)
}

// resolveFilteredPresentations finds all presentations matching a specific OData type filter
func resolveFilteredPresentations(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, definitionTemplateID, odataTypeFilter string) ([]string, error) {
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Finding presentations for definitionID='%s' with filter='%s'", definitionTemplateID, odataTypeFilter))

	// Get all presentations for this definition
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
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Found %d total presentations for definition", len(presentationList)))

	var filteredIDs []string
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

		// Filter by OData type
		if *odataType == odataTypeFilter {
			filteredIDs = append(filteredIDs, *presentationID)
			tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Added filtered presentation %d: type=%s, id=%s", i, *odataType, *presentationID))
		}
	}

	if len(filteredIDs) == 0 {
		return nil, fmt.Errorf("no presentations found with OData type %s for definition %s", odataTypeFilter, definitionTemplateID)
	}

	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Found %d filtered presentations", len(filteredIDs)))
	return filteredIDs, nil
}

// resolveTemplateIDsToInstanceIDs navigates from policy templates to their actual configured instances
// within a specific Group Policy Configuration. Supports both single and multiple presentation scenarios.
func resolveTemplateIDsToInstanceIDs(ctx context.Context, configID, definitionTemplateID string, presentationTemplateIDs []string, client *msgraphbetasdk.GraphServiceClient) (string, []ResolvedPresentation, error) {
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Resolving template IDs to instance IDs in config '%s' (definitionTemplateID: %s)", configID, definitionTemplateID))

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
		return "", nil, fmt.Errorf("failed to get definition values: %w", err)
	}

	// Check if definitionValues response is valid
	if definitionValues == nil {
		return "", nil, fmt.Errorf("received nil definition values response")
	}

	values := definitionValues.GetValue()
	if values == nil {
		return "", nil, fmt.Errorf("no definition values found in configuration")
	}

	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Found %d definition values in configuration", len(values)))

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

			tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Found matching definition value with instance ID: %s", *defValueID))

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

			tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Found %d presentation values for definition instance", len(presValues)))

			// Resolve presentation template IDs to instance IDs
			resolved := resolvePresentationInstances(ctx, presValues, presentationTemplateIDs)

			if len(resolved) > 0 {
				return *defValueID, resolved, nil
			}
		}
	}

	return "", nil, fmt.Errorf("no existing definition value found for definition template ID: %s", definitionTemplateID)
}

// resolvePresentationInstances maps presentation template IDs to their instance IDs using dual strategy
func resolvePresentationInstances(ctx context.Context, presValues []graphmodels.GroupPolicyPresentationValueable, presentationTemplateIDs []string) []ResolvedPresentation {
	// First, check if any presentation references are available
	tflog.Debug(ctx, "[RESOLVER] Checking for presentation references in API response:")
	hasAnyPresentationRefs := false
	for i, presValue := range presValues {
		if presValue != nil {
			presentation := presValue.GetPresentation()
			if presentation != nil && presentation.GetId() != nil {
				hasAnyPresentationRefs = true
				tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] - Presentation %d: ref_id='%s', instance_id='%s'", i, *presentation.GetId(), func() string {
					if presValue.GetId() != nil {
						return *presValue.GetId()
					}
					return "nil"
				}()))
			} else {
				tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] - Presentation %d: ref_id=nil, instance_id='%s'", i, func() string {
					if presValue.GetId() != nil {
						return *presValue.GetId()
					}
					return "nil"
				}()))
			}
		}
	}

	var resolved []ResolvedPresentation
	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Mapping %d presentation template IDs to instance IDs", len(presentationTemplateIDs)))

	if hasAnyPresentationRefs {
		// Strategy 1: Use presentation references (when available)
		tflog.Debug(ctx, "[RESOLVER] Using presentation reference mapping strategy")

		for i, templateID := range presentationTemplateIDs {
			tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Looking for template ID '%s' (index %d)", templateID, i))

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
					tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Comparing presentation ref ID '%s' with template ID '%s'", *presentationRefID, templateID))
					if *presentationRefID == templateID {
						presValueID := presValue.GetId()
						if presValueID != nil {
							resolved = append(resolved, ResolvedPresentation{
								TemplateID: templateID,
								InstanceID: *presValueID,
								Index:      i,
							})
							tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] ✅ Found presentation value instance ID: %s for template: %s", *presValueID, templateID))
							found = true
							break
						}
					}
				}
			}

			if !found {
				tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] ❌ No presentation value found for template ID '%s'", templateID))
			}
		}
	} else {
		// Strategy 2: Use positional mapping (fallback when presentation references are missing)
		tflog.Debug(ctx, "[RESOLVER] No presentation references found - using positional mapping strategy")

		// Map by position: presValues[0] -> templateIDs[0], presValues[1] -> templateIDs[1], etc.
		maxMappings := len(presentationTemplateIDs)
		if len(presValues) < maxMappings {
			maxMappings = len(presValues)
		}

		for i := 0; i < maxMappings; i++ {
			presValue := presValues[i]
			if presValue == nil {
				tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] ❌ Presentation value at index %d is nil", i))
				continue
			}

			presValueID := presValue.GetId()
			if presValueID == nil {
				tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] ❌ Presentation value at index %d has nil ID", i))
				continue
			}

			templateID := presentationTemplateIDs[i]
			resolved = append(resolved, ResolvedPresentation{
				TemplateID: templateID,
				InstanceID: *presValueID,
				Index:      i,
			})
			tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] ✅ Positional mapping: template[%d]='%s' -> instance='%s'", i, templateID, *presValueID))
		}

		if len(presentationTemplateIDs) > len(presValues) {
			tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] ⚠️ Warning: %d template IDs but only %d presentation values - some templates not mapped", len(presentationTemplateIDs), len(presValues)))
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("[RESOLVER] Successfully resolved %d out of %d presentation template IDs", len(resolved), len(presentationTemplateIDs)))
	return resolved
}
