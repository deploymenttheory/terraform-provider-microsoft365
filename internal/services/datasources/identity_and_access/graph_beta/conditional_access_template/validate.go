package graphBetaConditionalAccessTemplate

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/lithammer/fuzzysearch/fuzzy"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

// validateRequest validates that the provided template_id or name exists in the available templates.
// If not found, it returns an error with suggestions for similar valid values using fuzzy matching.
func (d *ConditionalAccessTemplateDataSource) validateRequest(ctx context.Context, templateID, name string, resp *datasource.ReadResponse) (graphmodels.ConditionalAccessTemplateable, bool) {
	var templates struct {
		Value []*templateResponse `json:"value"`
	}
	err := customrequests.JSONRequest(ctx, d.client.GetAdapter(), abstractions.GET,
		"/identity/conditionalAccess/templates", nil, nil, &templates)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Retrieve Templates",
			fmt.Sprintf("Unable to retrieve conditional access templates from Microsoft Graph API: %s", err.Error()),
		)
		return nil, false
	}

	if len(templates.Value) == 0 {
		resp.Diagnostics.AddError(
			"No Templates Found",
			"The API returned no conditional access templates",
		)
		return nil, false
	}

	templateList := templates.Value

	// Build lists of all valid IDs and names for error messaging
	var validIDs []string
	var validNames []string
	idToTemplate := make(map[string]graphmodels.ConditionalAccessTemplateable)
	nameToTemplate := make(map[string]graphmodels.ConditionalAccessTemplateable)

	for _, template := range templateList {
		if template == nil {
			continue
		}
		if id := template.GetId(); id != nil && *id != "" {
			validIDs = append(validIDs, *id)
			idToTemplate[*id] = template
		}
		if templateName := template.GetName(); templateName != nil && *templateName != "" {
			validNames = append(validNames, *templateName)
			nameToTemplate[*templateName] = template
		}
	}

	// Try to find exact match
	if templateID != "" {
		if template, exists := idToTemplate[templateID]; exists {
			return template, true
		}

		// Not found - provide helpful error with fuzzy suggestions
		suggestions := fuzzy.RankFindFold(templateID, validIDs)
		sort.Sort(suggestions)

		var suggestionText string
		if len(suggestions) > 0 && len(suggestions) <= 5 {
			suggestionText = fmt.Sprintf("\n\nDid you mean one of these?\n  - %s", strings.Join(extractTargets(suggestions, 5), "\n  - "))
		}

		resp.Diagnostics.AddError(
			"Invalid Template ID",
			fmt.Sprintf("No conditional access template found with ID: %s%s\n\nValid template IDs:\n  - %s",
				templateID,
				suggestionText,
				strings.Join(validIDs, "\n  - ")),
		)
		return nil, false
	}

	// Search by name
	if name != "" {
		if template, exists := nameToTemplate[name]; exists {
			return template, true
		}

		// Not found - provide helpful error with fuzzy suggestions
		suggestions := fuzzy.RankFindFold(name, validNames)
		sort.Sort(suggestions)

		var suggestionText string
		if len(suggestions) > 0 && len(suggestions) <= 5 {
			suggestionText = fmt.Sprintf("\n\nDid you mean one of these?\n  - %s", strings.Join(extractTargets(suggestions, 5), "\n  - "))
		}

		resp.Diagnostics.AddError(
			"Invalid Template Name",
			fmt.Sprintf("No conditional access template found with name: %s%s\n\nValid template names:\n  - %s",
				name,
				suggestionText,
				strings.Join(validNames, "\n  - ")),
		)
		return nil, false
	}

	// Should never reach here due to schema validation
	resp.Diagnostics.AddError(
		"Invalid Request",
		"Either template_id or name must be provided",
	)
	return nil, false
}

// extractTargets extracts the target strings from fuzzy.Rank results, limited to maxCount
func extractTargets(ranks fuzzy.Ranks, maxCount int) []string {
	var targets []string
	count := 0
	for _, rank := range ranks {
		if count >= maxCount {
			break
		}
		targets = append(targets, rank.Target)
		count++
	}
	return targets
}
