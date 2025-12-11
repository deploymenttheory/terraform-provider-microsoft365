package graphBetaAgentInstance

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the API response to the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AgentInstanceResourceModel, agentInstance graphmodels.AgentInstanceable) {
	if agentInstance == nil {
		return
	}

	tflog.Debug(ctx, "Mapping agent instance response to Terraform state")

	// Map basic properties
	data.ID = convert.GraphToFrameworkString(agentInstance.GetId())
	data.DisplayName = convert.GraphToFrameworkString(agentInstance.GetDisplayName())
	data.ManagedBy = convert.GraphToFrameworkString(agentInstance.GetManagedBy())
	data.OriginatingStore = convert.GraphToFrameworkString(agentInstance.GetOriginatingStore())
	data.CreatedBy = convert.GraphToFrameworkString(agentInstance.GetCreatedBy())
	data.SourceAgentId = convert.GraphToFrameworkString(agentInstance.GetSourceAgentId())
	data.AgentIdentityBlueprintId = convert.GraphToFrameworkString(agentInstance.GetAgentIdentityBlueprintId())
	data.AgentIdentityId = convert.GraphToFrameworkString(agentInstance.GetAgentIdentityId())
	data.AgentUserId = convert.GraphToFrameworkString(agentInstance.GetAgentUserId())
	data.Url = convert.GraphToFrameworkString(agentInstance.GetUrl())
	data.PreferredTransport = convert.GraphToFrameworkString(agentInstance.GetPreferredTransport())

	// Map timestamps
	data.CreatedDateTime = convert.GraphToFrameworkTime(agentInstance.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(agentInstance.GetLastModifiedDateTime())

	// Map owner IDs
	data.OwnerIds = convert.GraphToFrameworkStringSet(ctx, agentInstance.GetOwnerIds())

	// Map additional interfaces - only set if there are actual interfaces, otherwise keep nil
	if interfaces := agentInstance.GetAdditionalInterfaces(); len(interfaces) > 0 {
		data.AdditionalInterfaces = make([]AgentInterfaceModel, len(interfaces))
		for i, iface := range interfaces {
			data.AdditionalInterfaces[i] = AgentInterfaceModel{
				Url:       convert.GraphToFrameworkString(iface.GetUrl()),
				Transport: convert.GraphToFrameworkString(iface.GetTransport()),
			}
		}
	}

	tflog.Debug(ctx, "Finished mapping agent instance response to Terraform state")
}

// MapAgentCardManifestToTerraform maps a separately fetched agentCardManifest to the state
func MapAgentCardManifestToTerraform(ctx context.Context, data *AgentInstanceResourceModel, manifest graphmodels.AgentCardManifestable) {
	if manifest == nil {
		return
	}

	tflog.Debug(ctx, "Mapping agentCardManifest to Terraform state")
	data.AgentCardManifest = mapAgentCardManifestToTerraform(ctx, manifest)
}

// mapAgentCardManifestToTerraform maps the agent card manifest from the API response
func mapAgentCardManifestToTerraform(ctx context.Context, manifest graphmodels.AgentCardManifestable) *AgentCardManifestResourceModel {
	if manifest == nil {
		return nil
	}

	model := &AgentCardManifestResourceModel{
		ID:                                convert.GraphToFrameworkString(manifest.GetId()),
		DisplayName:                       convert.GraphToFrameworkString(manifest.GetDisplayName()),
		Description:                       convert.GraphToFrameworkString(manifest.GetDescription()),
		IconUrl:                           convert.GraphToFrameworkString(manifest.GetIconUrl()),
		OriginatingStore:                  convert.GraphToFrameworkString(manifest.GetOriginatingStore()),
		ProtocolVersion:                   convert.GraphToFrameworkString(manifest.GetProtocolVersion()),
		Version:                           convert.GraphToFrameworkString(manifest.GetVersion()),
		DocumentationUrl:                  convert.GraphToFrameworkString(manifest.GetDocumentationUrl()),
		SupportsAuthenticatedExtendedCard: convert.GraphToFrameworkBool(manifest.GetSupportsAuthenticatedExtendedCard()),
		OwnerIds:                          convert.GraphToFrameworkStringSet(ctx, manifest.GetOwnerIds()),
		DefaultInputModes:                 convert.GraphToFrameworkStringSet(ctx, manifest.GetDefaultInputModes()),
		DefaultOutputModes:                convert.GraphToFrameworkStringSet(ctx, manifest.GetDefaultOutputModes()),
	}

	// Map provider
	if provider := manifest.GetProvider(); provider != nil {
		model.Provider = &AgentCardProviderModel{
			Organization: convert.GraphToFrameworkString(provider.GetOrganization()),
			Url:          convert.GraphToFrameworkString(provider.GetUrl()),
		}
	}

	// Define the extension object type for types.List
	extensionAttrTypes := map[string]attr.Type{
		"uri":         types.StringType,
		"description": types.StringType,
		"required":    types.BoolType,
		"params":      types.MapType{ElemType: types.StringType},
	}
	extensionObjectType := types.ObjectType{AttrTypes: extensionAttrTypes}

	// Define the skill object type for types.List
	skillAttrTypes := map[string]attr.Type{
		"id":           types.StringType,
		"display_name": types.StringType,
		"description":  types.StringType,
		"tags":         types.SetType{ElemType: types.StringType},
		"examples":     types.SetType{ElemType: types.StringType},
		"input_modes":  types.SetType{ElemType: types.StringType},
		"output_modes": types.SetType{ElemType: types.StringType},
	}
	skillObjectType := types.ObjectType{AttrTypes: skillAttrTypes}

	// Map capabilities - always populate since it's required
	capsModel := &AgentCardCapabilitiesModel{
		Streaming:              convert.GraphToFrameworkBoolWithDefault(nil, false),
		PushNotifications:      convert.GraphToFrameworkBoolWithDefault(nil, false),
		StateTransitionHistory: convert.GraphToFrameworkBoolWithDefault(nil, false),
		Extensions:             types.ListNull(extensionObjectType),
	}

	if capabilities := manifest.GetCapabilities(); capabilities != nil {
		capsModel.Streaming = convert.GraphToFrameworkBoolWithDefault(capabilities.GetStreaming(), false)
		capsModel.PushNotifications = convert.GraphToFrameworkBoolWithDefault(capabilities.GetPushNotifications(), false)
		capsModel.StateTransitionHistory = convert.GraphToFrameworkBoolWithDefault(capabilities.GetStateTransitionHistory(), false)

		if extensions := capabilities.GetExtensions(); len(extensions) > 0 {
			extList := make([]attr.Value, len(extensions))
			for i, ext := range extensions {
				paramsMap := types.MapNull(types.StringType)
				if params := ext.GetParams(); params != nil {
					if additionalData := params.GetAdditionalData(); len(additionalData) > 0 {
						paramsData := make(map[string]string)
						for k, v := range additionalData {
							if strVal, ok := v.(string); ok {
								paramsData[k] = strVal
							}
						}
						if mapValue, diags := types.MapValueFrom(ctx, types.StringType, paramsData); !diags.HasError() {
							paramsMap = mapValue
						}
					}
				}

				extObj, _ := types.ObjectValue(extensionAttrTypes, map[string]attr.Value{
					"uri":         convert.GraphToFrameworkString(ext.GetUri()),
					"description": convert.GraphToFrameworkString(ext.GetDescription()),
					"required":    convert.GraphToFrameworkBool(ext.GetRequired()),
					"params":      paramsMap,
				})
				extList[i] = extObj
			}
			capsModel.Extensions, _ = types.ListValue(extensionObjectType, extList)
		} else {
			capsModel.Extensions, _ = types.ListValue(extensionObjectType, []attr.Value{})
		}
	}

	model.Capabilities = capsModel

	// Map skills
	if skills := manifest.GetSkills(); len(skills) > 0 {
		skillList := make([]attr.Value, len(skills))
		for i, skill := range skills {
			skillObj, _ := types.ObjectValue(skillAttrTypes, map[string]attr.Value{
				"id":           convert.GraphToFrameworkString(skill.GetId()),
				"display_name": convert.GraphToFrameworkString(skill.GetDisplayName()),
				"description":  convert.GraphToFrameworkString(skill.GetDescription()),
				"tags":         convert.GraphToFrameworkStringSet(ctx, skill.GetTags()),
				"examples":     convert.GraphToFrameworkStringSet(ctx, skill.GetExamples()),
				"input_modes":  convert.GraphToFrameworkStringSet(ctx, skill.GetInputModes()),
				"output_modes": convert.GraphToFrameworkStringSet(ctx, skill.GetOutputModes()),
			})
			skillList[i] = skillObj
		}
		model.Skills, _ = types.ListValue(skillObjectType, skillList)
	} else {
		model.Skills, _ = types.ListValue(skillObjectType, []attr.Value{})
	}

	return model
}
