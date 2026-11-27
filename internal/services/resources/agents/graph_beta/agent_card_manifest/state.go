package graphBetaAgentCardManifest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *AgentCardManifestResourceModel, remote graphmodels.AgentCardManifestable) error {
	if remote == nil || remote.GetId() == nil {
		return fmt.Errorf("%w: Graph returned no manifest identity", errInvalidManifest)
	}
	payload, err := convert.GraphToJSONFields(remote)
	if err != nil {
		return fmt.Errorf("serialize manifest response: %w", err)
	}

	data.ID = convert.GraphToFrameworkString(remote.GetId())
	if data.AgentCardManifestID.IsNull() || data.AgentCardManifestID.IsUnknown() {
		data.AgentCardManifestID = data.ID
	}
	data.DisplayName = convert.GraphToFrameworkString(remote.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.IconURL = convert.GraphToFrameworkString(remote.GetIconUrl())
	data.ProtocolVersion = convert.GraphToFrameworkString(remote.GetProtocolVersion())
	data.Version = convert.GraphToFrameworkString(remote.GetVersion())
	data.DocumentationURL = convert.GraphToFrameworkString(remote.GetDocumentationUrl())
	data.ManagedBy = convert.GraphToFrameworkString(remote.GetManagedBy())
	data.OriginatingStore = convert.GraphToFrameworkString(remote.GetOriginatingStore())
	data.OwnerIDs = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, append([]string{}, remote.GetOwnerIds()...))
	data.DefaultInputModes = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, append([]string{}, remote.GetDefaultInputModes()...))
	data.DefaultOutputModes = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, append([]string{}, remote.GetDefaultOutputModes()...))
	data.SupportsAuthenticatedExtendedCard = convert.GraphToFrameworkBool(remote.GetSupportsAuthenticatedExtendedCard())
	data.CreatedBy = convert.GraphToFrameworkString(remote.GetCreatedBy())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remote.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remote.GetLastModifiedDateTime())
	data.Provider = types.ObjectNull(providerAttrTypes)
	if provider := remote.GetProvider(); provider != nil {
		data.Provider = types.ObjectValueMust(
			providerAttrTypes,
			map[string]attr.Value{"organization": convert.GraphToFrameworkString(provider.GetOrganization()), "url": convert.GraphToFrameworkString(provider.GetUrl())},
		)
	}
	var priorCaps CapabilitiesModel
	priorExtensions := []ExtensionModel{}
	if !data.Capabilities.IsNull() && !data.Capabilities.IsUnknown() {
		if d := data.Capabilities.As(ctx, &priorCaps, basetypes.ObjectAsOptions{}); d.HasError() {
			return fmt.Errorf("%w: decode previous capabilities: %s", errInvalidManifest, d)
		}
		if !priorCaps.Extensions.IsNull() && !priorCaps.Extensions.IsUnknown() {
			if d := priorCaps.Extensions.ElementsAs(ctx, &priorExtensions, false); d.HasError() {
				return fmt.Errorf("%w: decode previous extensions: %s", errInvalidManifest, d)
			}
		}
	}
	priorParams := make(map[string]jsontypes.Normalized, len(priorExtensions))
	for _, extension := range priorExtensions {
		priorParams[extension.URI.ValueString()] = extension.ParamsJSON
	}
	data.Capabilities = types.ObjectNull(capabilitiesAttrTypes)
	if caps := remote.GetCapabilities(); caps != nil {

		entries := make([]ExtensionModel, 0, len(caps.GetExtensions()))
		for _, extension := range caps.GetExtensions() {
			if extension == nil {
				return fmt.Errorf("%w: Graph returned a null capability extension", errInvalidManifest)
			}
			prior := priorParams[convert.GraphToFrameworkString(extension.GetUri()).ValueString()]
			params, err := convert.GraphToFrameworkNormalizedJSONObject(ctx, extension.GetParams(), prior)
			if err != nil {
				return fmt.Errorf("map extension params: %w", err)
			}
			entries = append(
				entries,
				ExtensionModel{
					URI:         convert.GraphToFrameworkString(extension.GetUri()),
					Description: convert.GraphToFrameworkString(extension.GetDescription()),
					Required:    convert.GraphToFrameworkBool(extension.GetRequired()),
					ParamsJSON:  params,
				},
			)
		}
		extensions, diagnostics := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: extensionAttrTypes}, entries)
		if diagnostics.HasError() {
			return fmt.Errorf("%w: map extensions: %s", errInvalidManifest, diagnostics)
		}
		model := CapabilitiesModel{
			Streaming:              convert.GraphToFrameworkBool(caps.GetStreaming()),
			PushNotifications:      convert.GraphToFrameworkBool(caps.GetPushNotifications()),
			StateTransitionHistory: convert.GraphToFrameworkBool(caps.GetStateTransitionHistory()),
			Extensions:             extensions,
		}
		value, diagnostics := types.ObjectValueFrom(ctx, capabilitiesAttrTypes, model)
		if diagnostics.HasError() {
			return fmt.Errorf("%w: map capabilities: %s", errInvalidManifest, diagnostics)
		}
		data.Capabilities = value
	}
	entries := make([]SkillModel, 0, len(remote.GetSkills()))
	for _, skill := range remote.GetSkills() {
		if skill == nil {
			return fmt.Errorf("%w: Graph returned a null skill", errInvalidManifest)
		}
		entries = append(entries, SkillModel{
			ID: convert.GraphToFrameworkString(
				skill.GetId(),
			),
			DisplayName: convert.GraphToFrameworkString(skill.GetDisplayName()),
			Description: convert.GraphToFrameworkString(skill.GetDescription()),
			Tags:        convert.GraphToFrameworkStringSetPreserveEmpty(ctx, append([]string{}, skill.GetTags()...)),
			Examples:    convert.GraphToFrameworkStringSetPreserveEmpty(ctx, append([]string{}, skill.GetExamples()...)),
			InputModes:  convert.GraphToFrameworkStringSetPreserveEmpty(ctx, append([]string{}, skill.GetInputModes()...)),
			OutputModes: convert.GraphToFrameworkStringSetPreserveEmpty(ctx, append([]string{}, skill.GetOutputModes()...)),
		})
	}
	skills, diagnostics := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: skillAttrTypes}, entries)
	if diagnostics.HasError() {
		return fmt.Errorf("%w: map skills: %s", errInvalidManifest, diagnostics)
	}
	data.Skills = skills
	schemes := payload["securitySchemes"]
	if len(schemes) == 0 || string(schemes) == "null" {
		schemes = json.RawMessage(`{}`)
	}
	data.SecuritySchemesJSON, err = convert.GraphToFrameworkNormalizedJSON(ctx, schemes, data.SecuritySchemesJSON)
	if err != nil {
		return fmt.Errorf("map security schemes: %w", err)
	}
	security := payload["security"]
	if len(security) == 0 || string(security) == "null" {
		security = json.RawMessage(`[]`)
	}
	data.SecurityJSON, err = convert.GraphToFrameworkNormalizedJSON(ctx, security, data.SecurityJSON)
	if err != nil {
		return fmt.Errorf("map security requirements: %w", err)
	}
	return nil
}
