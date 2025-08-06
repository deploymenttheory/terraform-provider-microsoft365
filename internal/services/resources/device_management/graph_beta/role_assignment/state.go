package graphBetaRoleDefinitionAssignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote object state to the Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data *RoleAssignmentResourceModel, remoteResource graphmodels.DeviceAndAppManagementRoleAssignmentable) {
	if remoteResource == nil {
		return
	}

	if remoteResource.GetId() != nil {
		data.ID = types.StringValue(*remoteResource.GetId())
	}

	if remoteResource.GetDisplayName() != nil {
		data.DisplayName = types.StringValue(*remoteResource.GetDisplayName())
	} else {
		data.DisplayName = types.StringNull()
	}

	if remoteResource.GetDescription() != nil {
		data.Description = types.StringValue(*remoteResource.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Map members
	if remoteResource.GetMembers() != nil {
		members := remoteResource.GetMembers()
		memberElements := make([]attr.Value, len(members))
		for i, member := range members {
			memberElements[i] = types.StringValue(member)
		}
		membersSet, _ := types.SetValue(types.StringType, memberElements)
		data.Members = membersSet
	} else {
		data.Members = types.SetNull(types.StringType)
	}

	// Map role definition ID from additional data
	if additionalData := remoteResource.GetAdditionalData(); additionalData != nil {
		if roleDefBind, ok := additionalData["roleDefinition@odata.bind"].(string); ok {
			// Extract role definition ID from the odata.bind URL
			// Format: "https://graph.microsoft.com/beta/deviceManagement/roleDefinitions('id')"
			if len(roleDefBind) > 0 {
				// Extract ID from the URL - look for the last occurrence of single quotes
				start := -1
				end := -1
				for i := len(roleDefBind) - 1; i >= 0; i-- {
					if roleDefBind[i] == '\'' {
						if end == -1 {
							end = i
						} else {
							start = i + 1
							break
						}
					}
				}
				if start != -1 && end != -1 && start < end {
					data.RoleDefinitionId = types.StringValue(roleDefBind[start:end])
				}
			}
		}
	}

	// Map scope configuration
	scopeConfig := ScopeConfigurationResourceModel{}

	// Check scope type
	if remoteResource.GetScopeType() != nil {
		switch *remoteResource.GetScopeType() {
		case graphmodels.ALLLICENSEDUSERS_ROLEASSIGNMENTSCOPETYPE:
			scopeConfig.Type = types.StringValue("AllLicensedUsers")
			scopeConfig.ResourceScopes = types.SetNull(types.StringType)
		case graphmodels.ALLDEVICES_ROLEASSIGNMENTSCOPETYPE:
			scopeConfig.Type = types.StringValue("AllDevices")
			scopeConfig.ResourceScopes = types.SetNull(types.StringType)
		default:
			// Default to ResourceScopes if unknown type
			scopeConfig.Type = types.StringValue("ResourceScopes")
			mapResourceScopes(remoteResource, &scopeConfig)
		}
	} else {
		// No scope type means it's resource scopes
		scopeConfig.Type = types.StringValue("ResourceScopes")
		mapResourceScopes(remoteResource, &scopeConfig)
	}

	data.ScopeConfig = []ScopeConfigurationResourceModel{scopeConfig}
}

// mapResourceScopes maps the resource scopes from the remote resource
func mapResourceScopes(remoteResource graphmodels.DeviceAndAppManagementRoleAssignmentable, scopeConfig *ScopeConfigurationResourceModel) {
	if remoteResource.GetResourceScopes() != nil {
		resourceScopes := remoteResource.GetResourceScopes()
		scopeElements := make([]attr.Value, len(resourceScopes))
		for i, scope := range resourceScopes {
			scopeElements[i] = types.StringValue(scope)
		}
		resourceScopesSet, _ := types.SetValue(types.StringType, scopeElements)
		scopeConfig.ResourceScopes = resourceScopesSet
	} else {
		scopeConfig.ResourceScopes = types.SetNull(types.StringType)
	}
}
