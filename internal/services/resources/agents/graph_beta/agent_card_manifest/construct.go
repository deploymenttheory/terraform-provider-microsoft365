package graphBetaAgentCardManifest

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource sends configured writable properties only, never values
// populated by an earlier Read for unconfigured Optional+Computed attributes.
func constructResource(ctx context.Context, data *AgentCardManifestResourceModel) (graphmodels.AgentCardManifestable, error) {
	body := graphmodels.NewAgentCardManifest()
	convert.FrameworkToGraphString(data.DisplayName, body.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, body.SetDescription)
	convert.FrameworkToGraphString(data.IconURL, body.SetIconUrl)
	convert.FrameworkToGraphString(data.ProtocolVersion, body.SetProtocolVersion)
	convert.FrameworkToGraphString(data.Version, body.SetVersion)
	convert.FrameworkToGraphString(data.DocumentationURL, body.SetDocumentationUrl)
	convert.FrameworkToGraphString(data.ManagedBy, body.SetManagedBy)
	convert.FrameworkToGraphString(data.OriginatingStore, body.SetOriginatingStore)
	if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, data.OwnerIDs, body.SetOwnerIds); err != nil {
		return nil, fmt.Errorf("construct owner_ids: %w", err)
	}
	if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, data.DefaultInputModes, body.SetDefaultInputModes); err != nil {
		return nil, fmt.Errorf("construct default_input_modes: %w", err)
	}
	if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, data.DefaultOutputModes, body.SetDefaultOutputModes); err != nil {
		return nil, fmt.Errorf("construct default_output_modes: %w", err)
	}
	convert.FrameworkToGraphBool(data.SupportsAuthenticatedExtendedCard, body.SetSupportsAuthenticatedExtendedCard)
	if !data.Provider.IsNull() && !data.Provider.IsUnknown() {
		var model ProviderModel
		if d := data.Provider.As(ctx, &model, basetypes.ObjectAsOptions{}); d.HasError() {
			return nil, fmt.Errorf("%w: decode provider: %s", errInvalidManifest, d)
		}
		provider := graphmodels.NewAgentProvider()
		convert.FrameworkToGraphString(model.Organization, provider.SetOrganization)
		convert.FrameworkToGraphString(model.URL, provider.SetUrl)
		body.SetProvider(provider)
	}
	if !data.Capabilities.IsNull() && !data.Capabilities.IsUnknown() {
		var model CapabilitiesModel
		if d := data.Capabilities.As(ctx, &model, basetypes.ObjectAsOptions{}); d.HasError() {
			return nil, fmt.Errorf("%w: decode capabilities: %s", errInvalidManifest, d)
		}
		caps := graphmodels.NewAgentCapabilities()
		convert.FrameworkToGraphBool(model.Streaming, caps.SetStreaming)
		convert.FrameworkToGraphBool(model.PushNotifications, caps.SetPushNotifications)
		convert.FrameworkToGraphBool(model.StateTransitionHistory, caps.SetStateTransitionHistory)
		if !model.Extensions.IsNull() && !model.Extensions.IsUnknown() {
			var entries []ExtensionModel
			if d := model.Extensions.ElementsAs(ctx, &entries, false); d.HasError() {
				return nil, fmt.Errorf("%w: decode extensions: %s", errInvalidManifest, d)
			}
			extensions := make([]graphmodels.AgentExtensionable, 0, len(entries))
			for _, entry := range entries {
				extension := graphmodels.NewAgentExtension()
				convert.FrameworkToGraphString(entry.URI, extension.SetUri)
				convert.FrameworkToGraphString(entry.Description, extension.SetDescription)
				convert.FrameworkToGraphBool(entry.Required, extension.SetRequired)
				if err := convert.FrameworkToGraphNormalizedJSONObject(entry.ParamsJSON, func(value map[string]any) {
					params := graphmodels.NewAgentExtensionParams()
					params.SetAdditionalData(value)
					extension.SetParams(params)
				}); err != nil {
					return nil, fmt.Errorf("construct extension params_json: %w", err)
				}

				extensions = append(extensions, extension)
			}
			caps.SetExtensions(extensions)
		}
		body.SetCapabilities(caps)
	}
	if !data.Skills.IsNull() && !data.Skills.IsUnknown() {
		var entries []SkillModel
		if d := data.Skills.ElementsAs(ctx, &entries, false); d.HasError() {
			return nil, fmt.Errorf("%w: decode skills: %s", errInvalidManifest, d)
		}
		skills := make([]graphmodels.AgentSkillable, 0, len(entries))
		for _, entry := range entries {
			skill := graphmodels.NewAgentSkill()
			convert.FrameworkToGraphString(entry.ID, skill.SetId)
			convert.FrameworkToGraphString(entry.DisplayName, skill.SetDisplayName)
			convert.FrameworkToGraphString(entry.Description, skill.SetDescription)
			if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, entry.Tags, skill.SetTags); err != nil {
				return nil, fmt.Errorf("construct skill Tags: %w", err)
			}
			if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, entry.Examples, skill.SetExamples); err != nil {
				return nil, fmt.Errorf("construct skill Examples: %w", err)
			}
			if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, entry.InputModes, skill.SetInputModes); err != nil {
				return nil, fmt.Errorf("construct skill InputModes: %w", err)
			}
			if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, entry.OutputModes, skill.SetOutputModes); err != nil {
				return nil, fmt.Errorf("construct skill OutputModes: %w", err)
			}
			skills = append(skills, skill)
		}
		body.SetSkills(skills)
	}
	if err := convert.FrameworkToGraphNormalizedJSONObject(data.SecuritySchemesJSON, func(value map[string]any) {
		schemes := graphmodels.NewSecuritySchemes()
		schemes.SetAdditionalData(value)
		body.SetSecuritySchemes(schemes)
	}); err != nil {
		return nil, fmt.Errorf("construct security_schemes_json: %w", err)
	}

	if !data.SecurityJSON.IsNull() && !data.SecurityJSON.IsUnknown() {
		var values []map[string]any
		if d := data.SecurityJSON.Unmarshal(&values); d.HasError() {
			return nil, fmt.Errorf("%w: decode security_json: %s", errInvalidManifest, d)
		}
		security := make([]graphmodels.SecurityRequirementable, 0, len(values))
		for _, value := range values {
			if value == nil {
				return nil, fmt.Errorf("%w: security_json entries must be objects, got null", errInvalidManifest)
			}
			requirement := graphmodels.NewSecurityRequirement()
			requirement.SetAdditionalData(value)
			security = append(security, requirement)
		}
		body.SetSecurity(security)
	}
	return body, nil
}
