package utilityGroupPolicyValueReference

import (
	"context"
	"fmt"
	"strings"
	"time"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// Read implements the datasource.DataSource interface
func (d *groupPolicyValueReferenceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config groupPolicyValueReferenceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := config.Timeouts.Read(ctx, ReadTimeout*time.Second)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	policyName := config.PolicyName.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Searching for group policy definitions with name: %s", policyName))

	// Normalize the policy name for fuzzy matching
	normalizedPolicyName := normalizeString(policyName)

	// Query Microsoft Graph API for group policy definitions
	// We use contains() for a broader search since we'll do fuzzy matching in code
	filter := fmt.Sprintf("contains(displayName, '%s')", escapeSingleQuotes(policyName))

	definitions, err := d.client.
		DeviceManagement().
		GroupPolicyDefinitions().
		Get(ctx, &devicemanagement.GroupPolicyDefinitionsRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.GroupPolicyDefinitionsRequestBuilderGetQueryParameters{
				Select: []string{"id", "displayName", "classType", "categoryPath", "explainText", "supportedOn", "policyType"},
				Filter: &filter,
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return
	}

	// Use PageIterator to collect exact and fuzzy matching definitions
	var exactMatches []graphmodels.GroupPolicyDefinitionable
	var fuzzyMatches []fuzzyMatchResult

	pageIterator, err := graphcore.NewPageIterator[graphmodels.GroupPolicyDefinitionable](
		definitions,
		d.client.GetAdapter(),
		graphmodels.CreateGroupPolicyDefinitionCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Page Iterator",
			fmt.Sprintf("Could not create page iterator for group policy definitions: %s", err.Error()),
		)
		return
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.GroupPolicyDefinitionable) bool {
		if item.GetDisplayName() == nil {
			return true // continue iteration
		}

		itemDisplayName := strings.TrimSpace(*item.GetDisplayName())
		normalizedItemName := normalizeString(itemDisplayName)

		// Check for exact match first
		if normalizedPolicyName == normalizedItemName {
			exactMatches = append(exactMatches, item)
			tflog.Debug(ctx, fmt.Sprintf("Exact match found: '%s'", itemDisplayName))
		} else if fuzzy.Match(normalizedPolicyName, normalizedItemName) {
			// Calculate similarity score for fuzzy matches
			score := calculateSimilarityScore(normalizedPolicyName, normalizedItemName)
			fuzzyMatches = append(fuzzyMatches, fuzzyMatchResult{
				definition: item,
				score:      score,
				name:       itemDisplayName,
			})
			tflog.Debug(ctx, fmt.Sprintf("Fuzzy match found: '%s' (score: %d)", itemDisplayName, score))
		}

		return true // continue iteration
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Iterate Group Policy Definitions",
			fmt.Sprintf("Error while iterating through group policy definitions: %s", err.Error()),
		)
		return
	}

	// Handle exact matches
	var allDefinitions []graphmodels.GroupPolicyDefinitionable
	if len(exactMatches) > 0 {
		allDefinitions = exactMatches
		tflog.Debug(ctx, fmt.Sprintf("Found %d exact matches for '%s'", len(exactMatches), policyName))

		// If multiple exact matches found, add informational warning
		if len(exactMatches) > 1 {
			matchDetails := formatMultipleMatchDetails(exactMatches)
			resp.Diagnostics.AddWarning(
				"Multiple Definitions Found",
				fmt.Sprintf("Found %d group policy definitions with the exact name '%s'.\n\n"+
					"This is common when policies exist for both User and Machine configurations.\n\n"+
					"Matched definitions:\n%s\n"+
					"All matching definitions are included in the results. "+
					"Use the class_type, category_path, and other attributes to distinguish between them in your configuration.",
					len(exactMatches), policyName, matchDetails),
			)
		}
	} else if len(fuzzyMatches) > 0 {
		// No exact match, but fuzzy matches found - return error with suggestions
		sortedSuggestions := sortAndFormatSuggestions(fuzzyMatches)
		resp.Diagnostics.AddError(
			"No Exact Match Found",
			fmt.Sprintf("No exact match found for policy name '%s'.\n\n"+
				"Did you mean one of these? (ranked by similarity):\n%s\n\n"+
				"Please use the exact policy name from the list above.",
				policyName, sortedSuggestions),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Using %d definitions for policy '%s'", len(allDefinitions), policyName))

	if len(allDefinitions) == 0 {
		resp.Diagnostics.AddWarning(
			"No Matching Definitions Found",
			fmt.Sprintf("No group policy definitions found with display name '%s'. Verify the policy name is correct and matches exactly.", policyName),
		)

		// Set empty results
		config.Id = types.StringValue(fmt.Sprintf("policy_name:%s", policyName))
		config.Definitions = types.ListNull(types.ObjectType{
			AttrTypes: definitionAttrTypes(),
		})

		resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
		return
	}

	// For each definition, fetch its presentations
	presentationsMap := make(map[string][]PresentationModel)
	for _, def := range allDefinitions {
		definitionID := ""
		if def.GetId() != nil {
			definitionID = *def.GetId()
		}

		// Fetch presentations for this definition
		presentations, err := d.client.
			DeviceManagement().
			GroupPolicyDefinitions().
			ByGroupPolicyDefinitionId(definitionID).
			Presentations().
			Get(ctx, nil)

		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to fetch presentations for definition %s: %s", definitionID, err.Error()))
			// Continue without presentations rather than failing the entire operation
			presentationsMap[definitionID] = []PresentationModel{}
			continue
		}

		if presentations != nil && presentations.GetValue() != nil {
			presentationsMap[definitionID] = mapPresentationsToState(presentations.GetValue())
		} else {
			presentationsMap[definitionID] = []PresentationModel{}
		}
	}

	// Map all definitions to state
	definitionsList := mapDefinitionsToState(ctx, allDefinitions, presentationsMap, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	config.Id = types.StringValue(fmt.Sprintf("policy_name:%s", policyName))
	config.Definitions = definitionsList

	tflog.Debug(ctx, fmt.Sprintf("Successfully retrieved %d definitions for policy '%s'", len(allDefinitions), policyName))

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// normalizeString normalizes a string for fuzzy matching by:
// - Converting to lowercase
// - Removing extra whitespace
// - Collapsing multiple spaces to single space
func normalizeString(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Trim leading/trailing whitespace
	s = strings.TrimSpace(s)

	// Replace multiple spaces with single space
	s = strings.Join(strings.Fields(s), " ")

	return s
}

// escapeSingleQuotes escapes single quotes in strings for OData filter queries
func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// fuzzyMatchResult holds a fuzzy match result with its similarity score
type fuzzyMatchResult struct {
	definition graphmodels.GroupPolicyDefinitionable
	score      int
	name       string
}

// calculateSimilarityScore calculates a similarity score between two strings
// Lower score means more similar (Levenshtein distance)
func calculateSimilarityScore(query, target string) int {
	return fuzzy.LevenshteinDistance(query, target)
}

// sortAndFormatSuggestions sorts fuzzy matches by score and formats them for display
func sortAndFormatSuggestions(matches []fuzzyMatchResult) string {
	// Deduplicate matches by name (keep best score for each unique name)
	uniqueMatches := make(map[string]fuzzyMatchResult)
	for _, match := range matches {
		existing, exists := uniqueMatches[match.name]
		if !exists || match.score < existing.score {
			uniqueMatches[match.name] = match
		}
	}

	// Convert map back to slice
	deduplicatedMatches := make([]fuzzyMatchResult, 0, len(uniqueMatches))
	for _, match := range uniqueMatches {
		deduplicatedMatches = append(deduplicatedMatches, match)
	}

	// Sort by score (lower is better)
	// Use a simple bubble sort since we typically have few matches
	for i := 0; i < len(deduplicatedMatches)-1; i++ {
		for j := 0; j < len(deduplicatedMatches)-i-1; j++ {
			if deduplicatedMatches[j].score > deduplicatedMatches[j+1].score {
				deduplicatedMatches[j], deduplicatedMatches[j+1] = deduplicatedMatches[j+1], deduplicatedMatches[j]
			}
		}
	}

	// Format suggestions (limit to top 10)
	var suggestions strings.Builder
	limit := len(deduplicatedMatches)
	if limit > 10 {
		limit = 10
	}

	for i := 0; i < limit; i++ {
		suggestions.WriteString(fmt.Sprintf("  %d. \"%s\"\n", i+1, deduplicatedMatches[i].name))
	}

	if len(deduplicatedMatches) > 10 {
		suggestions.WriteString(fmt.Sprintf("  ... and %d more\n", len(deduplicatedMatches)-10))
	}

	return suggestions.String()
}

// formatMultipleMatchDetails formats multiple exact matches for display
func formatMultipleMatchDetails(definitions []graphmodels.GroupPolicyDefinitionable) string {
	var details strings.Builder

	for i, def := range definitions {
		displayName := "Unknown"
		if def.GetDisplayName() != nil {
			displayName = *def.GetDisplayName()
		}

		classType := "unknown"
		if def.GetClassType() != nil {
			classType = def.GetClassType().String()
		}

		categoryPath := "N/A"
		if def.GetCategoryPath() != nil {
			categoryPath = *def.GetCategoryPath()
		}

		details.WriteString(fmt.Sprintf("  %d. \"%s\"\n", i+1, displayName))
		details.WriteString(fmt.Sprintf("     - Class Type: %s\n", classType))
		details.WriteString(fmt.Sprintf("     - Category: %s\n", categoryPath))

		if i < len(definitions)-1 {
			details.WriteString("\n")
		}
	}

	return details.String()
}
