package graphBetaAgentInstance

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource builds the request body for creating or updating an agent instance.
// When isCreate is true, the agentCardManifest is included in the request body.
// When isCreate is false (update), the agentCardManifest is excluded as it must be updated separately.
func constructResource(ctx context.Context, data *AgentInstanceResourceModel, isCreate bool) (graphmodels.AgentInstanceable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model (isCreate: %t)", ResourceName, isCreate))

	requestBody := graphmodels.NewAgentInstance()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	if err := convert.FrameworkToGraphStringSet(ctx, data.OwnerIds, requestBody.SetOwnerIds); err != nil {
		return nil, fmt.Errorf("failed to set owner_ids: %w", err)
	}

	convert.FrameworkToGraphString(data.ManagedBy, requestBody.SetManagedBy)
	convert.FrameworkToGraphString(data.OriginatingStore, requestBody.SetOriginatingStore)
	convert.FrameworkToGraphString(data.SourceAgentId, requestBody.SetSourceAgentId)
	convert.FrameworkToGraphString(data.AgentIdentityBlueprintId, requestBody.SetAgentIdentityBlueprintId)
	convert.FrameworkToGraphString(data.AgentIdentityId, requestBody.SetAgentIdentityId)
	convert.FrameworkToGraphString(data.Url, requestBody.SetUrl)
	convert.FrameworkToGraphString(data.PreferredTransport, requestBody.SetPreferredTransport)

	if len(data.AdditionalInterfaces) > 0 {
		interfaces := make([]graphmodels.AgentInterfaceable, len(data.AdditionalInterfaces))
		for i, iface := range data.AdditionalInterfaces {
			agentInterface := graphmodels.NewAgentInterface()
			convert.FrameworkToGraphString(iface.Url, agentInterface.SetUrl)
			convert.FrameworkToGraphString(iface.Transport, agentInterface.SetTransport)
			interfaces[i] = agentInterface
		}
		requestBody.SetAdditionalInterfaces(interfaces)
	}

	// Only include agentCardManifest during creation - updates are handled separately
	if isCreate && data.AgentCardManifest != nil {
		manifest, err := constructAgentCardManifest(ctx, data.AgentCardManifest)
		if err != nil {
			return nil, fmt.Errorf("failed to construct agent_card_manifest: %w", err)
		}
		requestBody.SetAgentCardManifest(manifest)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructAgentCardManifest builds the agent card manifest portion of the request body
func constructAgentCardManifest(ctx context.Context, data *AgentCardManifestResourceModel) (graphmodels.AgentCardManifestable, error) {
	manifest := graphmodels.NewAgentCardManifest()

	convert.FrameworkToGraphString(data.DisplayName, manifest.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, manifest.SetDescription)
	convert.FrameworkToGraphString(data.IconUrl, manifest.SetIconUrl)
	convert.FrameworkToGraphString(data.OriginatingStore, manifest.SetOriginatingStore)
	convert.FrameworkToGraphString(data.ProtocolVersion, manifest.SetProtocolVersion)
	convert.FrameworkToGraphString(data.Version, manifest.SetVersion)
	convert.FrameworkToGraphString(data.DocumentationUrl, manifest.SetDocumentationUrl)
	convert.FrameworkToGraphBool(data.SupportsAuthenticatedExtendedCard, manifest.SetSupportsAuthenticatedExtendedCard)

	if err := convert.FrameworkToGraphStringSet(ctx, data.OwnerIds, manifest.SetOwnerIds); err != nil {
		return nil, fmt.Errorf("failed to set manifest owner_ids: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.DefaultInputModes, manifest.SetDefaultInputModes); err != nil {
		return nil, fmt.Errorf("failed to set default_input_modes: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.DefaultOutputModes, manifest.SetDefaultOutputModes); err != nil {
		return nil, fmt.Errorf("failed to set default_output_modes: %w", err)
	}

	if data.Provider != nil {
		provider := graphmodels.NewAgentProvider()
		convert.FrameworkToGraphString(data.Provider.Organization, provider.SetOrganization)
		convert.FrameworkToGraphString(data.Provider.Url, provider.SetUrl)
		manifest.SetProvider(provider)
	}

	if data.Capabilities != nil {
		capabilities := graphmodels.NewAgentCapabilities()
		convert.FrameworkToGraphBool(data.Capabilities.Streaming, capabilities.SetStreaming)
		convert.FrameworkToGraphBool(data.Capabilities.PushNotifications, capabilities.SetPushNotifications)
		convert.FrameworkToGraphBool(data.Capabilities.StateTransitionHistory, capabilities.SetStateTransitionHistory)

		// Handle extensions as types.List
		if !data.Capabilities.Extensions.IsNull() && !data.Capabilities.Extensions.IsUnknown() {
			var extModels []AgentCardCapabilityExtensionModel
			diags := data.Capabilities.Extensions.ElementsAs(ctx, &extModels, false)
			if !diags.HasError() && len(extModels) > 0 {
				extensions := make([]graphmodels.AgentExtensionable, len(extModels))
				for i, ext := range extModels {
					extension := graphmodels.NewAgentExtension()
					convert.FrameworkToGraphString(ext.Uri, extension.SetUri)
					convert.FrameworkToGraphString(ext.Description, extension.SetDescription)
					convert.FrameworkToGraphBool(ext.Required, extension.SetRequired)

					if !ext.Params.IsNull() && !ext.Params.IsUnknown() {
						var params map[string]string
						ext.Params.ElementsAs(ctx, &params, false)
						if len(params) > 0 {
							paramsObj := graphmodels.NewAgentExtensionParams()
							paramsData := make(map[string]any)
							for k, v := range params {
								paramsData[k] = v
							}
							paramsObj.SetAdditionalData(paramsData)
							extension.SetParams(paramsObj)
						}
					}
					extensions[i] = extension
				}
				capabilities.SetExtensions(extensions)
			}
		}
		manifest.SetCapabilities(capabilities)
	}

	// Handle skills as types.List
	if !data.Skills.IsNull() && !data.Skills.IsUnknown() {
		var skillModels []AgentCardSkillModel
		diags := data.Skills.ElementsAs(ctx, &skillModels, false)
		if !diags.HasError() && len(skillModels) > 0 {
			skills := make([]graphmodels.AgentSkillable, len(skillModels))
			for i, skillData := range skillModels {
				skill := graphmodels.NewAgentSkill()
				convert.FrameworkToGraphString(skillData.ID, skill.SetId)
				convert.FrameworkToGraphString(skillData.DisplayName, skill.SetDisplayName)
				convert.FrameworkToGraphString(skillData.Description, skill.SetDescription)

				if err := convert.FrameworkToGraphStringSet(ctx, skillData.Tags, skill.SetTags); err != nil {
					return nil, fmt.Errorf("failed to set skill tags: %w", err)
				}
				if err := convert.FrameworkToGraphStringSet(ctx, skillData.Examples, skill.SetExamples); err != nil {
					return nil, fmt.Errorf("failed to set skill examples: %w", err)
				}
				if err := convert.FrameworkToGraphStringSet(ctx, skillData.InputModes, skill.SetInputModes); err != nil {
					return nil, fmt.Errorf("failed to set skill input_modes: %w", err)
				}
				if err := convert.FrameworkToGraphStringSet(ctx, skillData.OutputModes, skill.SetOutputModes); err != nil {
					return nil, fmt.Errorf("failed to set skill output_modes: %w", err)
				}

				skills[i] = skill
			}
			manifest.SetSkills(skills)
		}
	}

	return manifest, nil
}
