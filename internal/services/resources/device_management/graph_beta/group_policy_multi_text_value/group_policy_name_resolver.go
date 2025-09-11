package graphBetaGroupPolicyMultiTextValue

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// GroupPolicyDefinition represents a simplified group policy definition for lookups
type GroupPolicyDefinition struct {
	ID          string
	DisplayName string
	ClassType   string
}

// groupPolicyNameResolver finds the definition ID based on display name, class type, and optional category path
func groupPolicyNameResolver(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName, classType, categoryPath string) (string, error) {
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Looking up definition ID for displayName='%s', classType='%s', categoryPath='%s'", displayName, classType, categoryPath))

	// Parse class type using SDK enum helper
	classTypeEnum, err := graphmodels.ParseGroupPolicyDefinitionClassType(strings.ToLower(classType))
	if err != nil || classTypeEnum == nil {
		tflog.Error(ctx, fmt.Sprintf("[LOOKUP] Failed to parse class type '%s': %v", classType, err))
		return "", fmt.Errorf("invalid class_type '%s', must be 'user' or 'machine'", classType)
	}

	// Convert back to string for filter
	normalizedClassType := (*classTypeEnum.(*graphmodels.GroupPolicyDefinitionClassType)).String()
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Parsed classType '%s' -> normalized '%s'", classType, normalizedClassType))

	// Get all group policy definitions with a filter for the display name
	// Note: We'll need to fetch all definitions and filter client-side since Graph API
	// filtering on displayName might not work reliably
	filterQuery := fmt.Sprintf("displayName eq '%s' and classType eq '%s'", displayName, normalizedClassType)
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Attempting API filter query: %s", filterQuery))

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
		tflog.Error(ctx, fmt.Sprintf("[LOOKUP] Failed to fetch group policy definitions with filter: %s", err.Error()))
		return "", fmt.Errorf("failed to fetch group policy definitions: %w", err)
	}

	if definitions == nil || definitions.GetValue() == nil || len(definitions.GetValue()) == 0 {
		tflog.Debug(ctx, "[LOOKUP] No definitions found with exact filter (count=0), trying client-side filtering")

		// If no exact match found, try fetching all and filtering client-side
		allDefinitions, err := getAllDefinitions(ctx, client)
		if err != nil {
			return "", err
		}

		return findMatchingDefinition(ctx, allDefinitions, displayName, normalizedClassType, categoryPath)
	}

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found %d definitions with filter, checking for exact matches", len(definitions.GetValue())))

	// Check for matches
	var exactMatches []GroupPolicyDefinition
	var nameAndTypeMatches []GroupPolicyDefinition

	for i, def := range definitions.GetValue() {
		if def != nil &&
			def.GetDisplayName() != nil &&
			def.GetClassType() != nil &&
			def.GetId() != nil {

			defDisplayName := *def.GetDisplayName()
			defClassType := def.GetClassType().String()
			defID := *def.GetId()
			defCategoryPath := ""
			if def.GetCategoryPath() != nil {
				defCategoryPath = *def.GetCategoryPath()
			}

			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Checking definition %d: ID=%s, DisplayName='%s', ClassType='%s', CategoryPath='%s'", i, defID, defDisplayName, defClassType, defCategoryPath))

			displayNameMatch := strings.EqualFold(defDisplayName, displayName)
			classTypeMatch := defClassType == normalizedClassType
			categoryPathMatch := categoryPath == "" || strings.EqualFold(defCategoryPath, categoryPath)

			// Collect all matches by displayName and classType (regardless of categoryPath)
			if displayNameMatch && classTypeMatch {
				nameAndTypeMatches = append(nameAndTypeMatches, GroupPolicyDefinition{
					ID:          defID,
					DisplayName: defDisplayName,
					ClassType:   defClassType,
				})
			}

			// Only exact matches (including categoryPath if provided)
			if displayNameMatch && classTypeMatch && categoryPathMatch {
				tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ✅ EXACT MATCH FOUND: ID=%s", defID))
				exactMatches = append(exactMatches, GroupPolicyDefinition{
					ID:          defID,
					DisplayName: defDisplayName,
					ClassType:   defClassType,
				})
			} else {
				tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ❌ No exact match: DisplayName=%v, ClassType=%v, CategoryPath=%v",
					displayNameMatch, classTypeMatch, categoryPathMatch))
			}
		} else {
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ⚠️  Skipping definition %d: missing required fields", i))
		}
	}

	// Handle results - check for ambiguous case FIRST
	if len(nameAndTypeMatches) > 1 && categoryPath == "" {
		// Multiple matches found but no category_path provided - this is the error case
		tflog.Warn(ctx, fmt.Sprintf("[LOOKUP] Found %d matches for displayName='%s', classType='%s' but no category_path provided", len(nameAndTypeMatches), displayName, classType))

		// Collect category paths from all matches
		var categoryPaths []string
		for _, match := range nameAndTypeMatches {
			// We need to get the category path for each match - let's get it from the original definitions
			for _, def := range definitions.GetValue() {
				if def != nil && def.GetId() != nil && *def.GetId() == match.ID {
					if def.GetCategoryPath() != nil {
						categoryPaths = append(categoryPaths, *def.GetCategoryPath())
					} else {
						categoryPaths = append(categoryPaths, "(no category)")
					}
					break
				}
			}
		}

		return "", fmt.Errorf("multiple group policy matches found for displayName='%s' and classType='%s', please add the category_path and include one of the following: %s",
			displayName, classType, strings.Join(categoryPaths, ", "))
	}

	// Now handle successful matches
	if len(exactMatches) == 1 {
		// Single exact match found
		return exactMatches[0].ID, nil
	} else if len(exactMatches) > 1 {
		// Multiple exact matches (shouldn't happen, but handle gracefully)
		tflog.Warn(ctx, fmt.Sprintf("[LOOKUP] Multiple exact matches found: %d", len(exactMatches)))
		return exactMatches[0].ID, nil
	} else if len(nameAndTypeMatches) == 1 {
		// Single name+type match (no categoryPath needed)
		return nameAndTypeMatches[0].ID, nil
	}

	return "", fmt.Errorf("no group policy definition found with displayName='%s' and classType='%s'", displayName, classType)
}

// getAllDefinitions fetches all group policy definitions (with pagination handling)
func getAllDefinitions(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) ([]graphmodels.GroupPolicyDefinitionable, error) {
	tflog.Debug(ctx, "[LOOKUP] Fetching all group policy definitions for client-side filtering")

	var allDefinitions []graphmodels.GroupPolicyDefinitionable

	// Initial request
	tflog.Debug(ctx, "[LOOKUP] Making API call to fetch all definitions with select=id,displayName,classType,categoryPath and top=999")
	definitions, err := client.
		DeviceManagement().
		GroupPolicyDefinitions().
		Get(ctx, &devicemanagement.GroupPolicyDefinitionsRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.GroupPolicyDefinitionsRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName", "classType", "categoryPath"},
				Top:    &[]int32{999}[0], // Get as many as possible in first request
			},
		})

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("[LOOKUP] Failed to fetch all group policy definitions: %s", err.Error()))
		return nil, fmt.Errorf("failed to fetch group policy definitions: %w", err)
	}

	if definitions != nil && definitions.GetValue() != nil {
		allDefinitions = append(allDefinitions, definitions.GetValue()...)
		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Successfully fetched %d group policy definitions", len(allDefinitions)))

		// Log first few definitions for debugging
		for i, def := range allDefinitions {
			if i >= 5 { // Only log first 5 to avoid spam
				tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ... and %d more definitions", len(allDefinitions)-5))
				break
			}
			if def != nil && def.GetDisplayName() != nil && def.GetClassType() != nil && def.GetId() != nil {
				categoryPath := ""
				if def.GetCategoryPath() != nil {
					categoryPath = *def.GetCategoryPath()
				}
				tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Definition %d: ID=%s, DisplayName='%s', ClassType='%s', CategoryPath='%s'",
					i, *def.GetId(), *def.GetDisplayName(), def.GetClassType().String(), categoryPath))
			}
		}
	} else {
		tflog.Warn(ctx, "[LOOKUP] API returned empty definitions list")
	}

	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Total fetched: %d group policy definitions", len(allDefinitions)))
	return allDefinitions, nil
}

// findMatchingDefinition searches for a matching definition in the provided list
func findMatchingDefinition(ctx context.Context, definitions []graphmodels.GroupPolicyDefinitionable, displayName, classType, categoryPath string) (string, error) {
	tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Client-side searching %d definitions for displayName='%s', classType='%s', categoryPath='%s'", len(definitions), displayName, classType, categoryPath))

	var exactMatches []GroupPolicyDefinition
	var nameAndTypeMatches []GroupPolicyDefinition

	for i, def := range definitions {
		if def == nil || def.GetDisplayName() == nil || def.GetClassType() == nil || def.GetId() == nil {
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Skipping definition %d: missing required fields", i))
			continue
		}

		defDisplayName := *def.GetDisplayName()
		defClassTypeString := def.GetClassType().String()
		defID := *def.GetId()
		defCategoryPath := ""
		if def.GetCategoryPath() != nil {
			defCategoryPath = *def.GetCategoryPath()
		}

		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Client-side checking definition %d: ID=%s, DisplayName='%s', ClassType='%s', CategoryPath='%s'", i, defID, defDisplayName, defClassTypeString, defCategoryPath))

		displayNameMatch := strings.EqualFold(defDisplayName, displayName)
		classTypeMatch := defClassTypeString == classType
		categoryPathMatch := categoryPath == "" || strings.EqualFold(defCategoryPath, categoryPath)

		// Collect all matches by displayName and classType (regardless of categoryPath)
		if displayNameMatch && classTypeMatch {
			nameAndTypeMatches = append(nameAndTypeMatches, GroupPolicyDefinition{
				ID:          defID,
				DisplayName: defDisplayName,
				ClassType:   defClassTypeString,
			})
		}

		// Exact match (including category path if provided)
		if displayNameMatch && classTypeMatch && categoryPathMatch {
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ✅ CLIENT-SIDE EXACT MATCH FOUND: ID=%s", defID))
			exactMatches = append(exactMatches, GroupPolicyDefinition{
				ID:          defID,
				DisplayName: defDisplayName,
				ClassType:   defClassTypeString,
			})
		} else {
			tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] ❌ Not exact match: DisplayName=%v, ClassType=%v, CategoryPath=%v",
				displayNameMatch, classTypeMatch, categoryPathMatch))
		}
	}

	// Handle results - check for ambiguous case FIRST
	if len(nameAndTypeMatches) > 1 && categoryPath == "" {
		// Multiple matches found but no category_path provided - this is the error case
		tflog.Warn(ctx, fmt.Sprintf("[LOOKUP] Found %d matches for displayName='%s', classType='%s' but no category_path provided", len(nameAndTypeMatches), displayName, classType))

		// Collect category paths from all matches
		var categoryPaths []string
		for _, match := range nameAndTypeMatches {
			// We need to get the category path for each match - let's get it from the original definitions
			for _, def := range definitions {
				if def != nil && def.GetId() != nil && *def.GetId() == match.ID {
					if def.GetCategoryPath() != nil {
						categoryPaths = append(categoryPaths, *def.GetCategoryPath())
					} else {
						categoryPaths = append(categoryPaths, "(no category)")
					}
					break
				}
			}
		}

		return "", fmt.Errorf("multiple group policy matches found for displayName='%s' and classType='%s', please add the category_path and include one of the following: %s",
			displayName, classType, strings.Join(categoryPaths, ", "))
	}

	// Now handle successful matches
	if len(exactMatches) == 1 {
		// Single exact match found
		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Returning single exact match: %s", exactMatches[0].ID))
		return exactMatches[0].ID, nil
	} else if len(exactMatches) > 1 {
		// Multiple exact matches (shouldn't happen, but handle gracefully)
		tflog.Warn(ctx, fmt.Sprintf("[LOOKUP] Multiple exact matches found: %d", len(exactMatches)))
		return exactMatches[0].ID, nil
	} else if len(nameAndTypeMatches) == 1 {
		// Single match found
		tflog.Debug(ctx, fmt.Sprintf("[LOOKUP] Found single match: ID=%s, DisplayName='%s'", nameAndTypeMatches[0].ID, nameAndTypeMatches[0].DisplayName))
		return nameAndTypeMatches[0].ID, nil
	}

	// No matches found
	return "", fmt.Errorf("no group policy definition found with displayName='%s' and classType='%s'", displayName, classType)
}
