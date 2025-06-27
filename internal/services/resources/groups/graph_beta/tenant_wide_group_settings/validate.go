package graphBetaTenantWideGroupSettings

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// TemplateIDValidator validates that the template_id is a valid directory setting template ID
type TemplateIDValidator struct {
	client *msgraphbetasdk.GraphServiceClient
}

// NewTemplateIDValidator creates a new TemplateIDValidator
func NewTemplateIDValidator(client *msgraphbetasdk.GraphServiceClient) *TemplateIDValidator {
	return &TemplateIDValidator{
		client: client,
	}
}

// Description returns the validation description
func (v *TemplateIDValidator) Description(ctx context.Context) string {
	return "Template ID must be a valid directory setting template ID from Microsoft Graph"
}

// MarkdownDescription returns the validation description in markdown format
func (v *TemplateIDValidator) MarkdownDescription(ctx context.Context) string {
	return "Template ID must be a valid directory setting template ID from Microsoft Graph"
}

// ValidateString performs the validation
func (v *TemplateIDValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	templateID := req.ConfigValue.ValueString()
	if templateID == "" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Template ID",
			"Template ID cannot be empty",
		)
		return
	}

	// Fetch available templates from Microsoft Graph API
	availableTemplates, err := v.fetchDirectorySettingTemplates(ctx)
	if err != nil {
		tflog.Warn(ctx, "Failed to fetch directory setting templates for validation, skipping validation", map[string]interface{}{
			"error": err.Error(),
		})
		// Don't fail validation if we can't fetch templates - this allows offline usage
		return
	}

	// Check if the provided template ID exists in the available templates
	var validTemplateIDs []string
	var templateNames []string
	found := false

	for _, template := range availableTemplates {
		validTemplateIDs = append(validTemplateIDs, *template.GetId())
		templateNames = append(templateNames, *template.GetDisplayName())

		if *template.GetId() == templateID {
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Template ID",
			fmt.Sprintf("Template ID '%s' is not valid. Available templates: %s",
				templateID,
				strings.Join(templateNames, ", ")),
		)
		return
	}

	tflog.Debug(ctx, "Template ID validation passed", map[string]interface{}{
		"template_id":         templateID,
		"available_templates": templateNames,
	})
}

// fetchDirectorySettingTemplates fetches the list of available directory setting templates
func (v *TemplateIDValidator) fetchDirectorySettingTemplates(ctx context.Context) ([]graphmodels.DirectorySettingTemplateable, error) {
	if v.client == nil {
		return nil, fmt.Errorf("Graph client is not initialized")
	}

	templates, err := v.client.DirectorySettingTemplates().Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch directory setting templates: %w", err)
	}

	if templates == nil || templates.GetValue() == nil {
		return nil, fmt.Errorf("no directory setting templates returned from API")
	}

	return templates.GetValue(), nil
}

// ValidateTemplateID is a helper function to validate template ID with a client
func ValidateTemplateID(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) validator.String {
	return NewTemplateIDValidator(client)
}
