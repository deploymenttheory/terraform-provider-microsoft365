package graphBetaIdentityAndAccessConnectedOrganization

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote connected organization from Kiota SDK to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data ConnectedOrganizationResourceModel, org graphmodels.ConnectedOrganizationable) ConnectedOrganizationResourceModel {
	if org == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName": org.GetDisplayName(),
		"resourceId":   org.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(org.GetId())
	data.DisplayName = convert.GraphToFrameworkString(org.GetDisplayName())
	data.State = convert.GraphToFrameworkEnum(org.GetState())
	data.CreatedBy = convert.GraphToFrameworkString(org.GetCreatedBy())
	data.CreatedDateTime = convert.GraphToFrameworkTime(org.GetCreatedDateTime())
	data.ModifiedBy = convert.GraphToFrameworkString(org.GetModifiedBy())
	data.ModifiedDateTime = convert.GraphToFrameworkTime(org.GetModifiedDateTime())

	// Graph keeps a cleared description as "" (see constructResource); both "" and null mean no description.
	if description := org.GetDescription(); description != nil && *description != "" {
		data.Description = types.StringValue(*description)
	} else {
		data.Description = types.StringNull()
	}

	data.IdentitySources = mapIdentitySources(ctx, org.GetIdentitySources())

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"resourceId": data.ID.ValueString(),
	})

	return data
}

// mapIdentitySources maps the polymorphic identitySources collection to Terraform state.
// Attributes that do not apply to a source's type are null.
func mapIdentitySources(ctx context.Context, sources []graphmodels.IdentitySourceable) []IdentitySourceModel {
	if len(sources) == 0 {
		return nil
	}

	result := make([]IdentitySourceModel, 0, len(sources))
	for _, source := range sources {
		if source == nil {
			continue
		}

		model := IdentitySourceModel{
			ODataType:   convert.GraphToFrameworkString(source.GetOdataType()),
			DisplayName: types.StringNull(),
			DomainName:  types.StringNull(),
			TenantID:    types.StringNull(),
			IssuerURI:   types.StringNull(),
		}

		switch typed := source.(type) {
		case *graphmodels.DomainIdentitySource:
			model.ODataType = types.StringValue(ODataTypeDomainIdentitySource)
			model.DisplayName = convert.GraphToFrameworkString(typed.GetDisplayName())
			model.DomainName = convert.GraphToFrameworkString(typed.GetDomainName())
		case *graphmodels.AzureActiveDirectoryTenant:
			model.ODataType = types.StringValue(ODataTypeAzureActiveDirectoryTenant)
			model.DisplayName = convert.GraphToFrameworkString(typed.GetDisplayName())
			model.TenantID = convert.GraphToFrameworkString(typed.GetTenantId())
		case *graphmodels.ExternalDomainFederation:
			model.ODataType = types.StringValue(ODataTypeExternalDomainFederation)
			model.DisplayName = convert.GraphToFrameworkString(typed.GetDisplayName())
			model.DomainName = convert.GraphToFrameworkString(typed.GetDomainName())
			model.IssuerURI = convert.GraphToFrameworkString(typed.GetIssuerUri())
		default:
			tflog.Warn(ctx, "Unsupported identity source type returned by Microsoft Graph", map[string]any{
				"odataType": model.ODataType.ValueString(),
			})
		}

		result = append(result, model)
	}

	return result
}
