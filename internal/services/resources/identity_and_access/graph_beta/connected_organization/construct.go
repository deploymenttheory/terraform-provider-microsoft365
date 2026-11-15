package graphBetaIdentityAndAccessConnectedOrganization

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource converts the Terraform resource model to a Kiota SDK model.
//
// For create (isUpdate false) the body carries the identity sources, and unset optional fields are omitted.
// For update (isUpdate true) the body only carries the properties Microsoft Graph documents as updatable
// (displayName, description, state). Identity sources are omitted because changing them replaces the resource.
// An unset description is sent as an empty string on update so a removed description is cleared;
// Graph stores it as "" and the state mapping reads "" back as null.
func constructResource(ctx context.Context, data *ConnectedOrganizationResourceModel, isUpdate bool) (graphmodels.ConnectedOrganizationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, data); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := graphmodels.NewConnectedOrganization()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	} else if isUpdate {
		empty := ""
		requestBody.SetDescription(&empty)
	}

	if err := convert.FrameworkToGraphEnum(data.State, graphmodels.ParseConnectedOrganizationState, requestBody.SetState); err != nil {
		return nil, fmt.Errorf("failed to set state: %w", err)
	}

	if !isUpdate {
		identitySources, err := constructIdentitySources(data.IdentitySources)
		if err != nil {
			return nil, fmt.Errorf("failed to construct identity sources: %w", err)
		}
		if len(identitySources) > 0 {
			requestBody.SetIdentitySources(identitySources)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructIdentitySources builds the polymorphic identitySources collection.
// Attributes that do not apply to a source's type are never set, so they are omitted from the request.
func constructIdentitySources(sources []IdentitySourceModel) ([]graphmodels.IdentitySourceable, error) {
	if len(sources) == 0 {
		return nil, nil
	}

	result := make([]graphmodels.IdentitySourceable, 0, len(sources))
	for i, source := range sources {
		switch source.ODataType.ValueString() {
		case ODataTypeDomainIdentitySource:
			domain := graphmodels.NewDomainIdentitySource()
			convert.FrameworkToGraphString(source.DisplayName, domain.SetDisplayName)
			convert.FrameworkToGraphString(source.DomainName, domain.SetDomainName)
			result = append(result, domain)
		case ODataTypeAzureActiveDirectoryTenant:
			tenant := graphmodels.NewAzureActiveDirectoryTenant()
			convert.FrameworkToGraphString(source.DisplayName, tenant.SetDisplayName)
			convert.FrameworkToGraphString(source.TenantID, tenant.SetTenantId)
			result = append(result, tenant)
		case ODataTypeExternalDomainFederation:
			federation := graphmodels.NewExternalDomainFederation()
			convert.FrameworkToGraphString(source.DisplayName, federation.SetDisplayName)
			convert.FrameworkToGraphString(source.DomainName, federation.SetDomainName)
			convert.FrameworkToGraphString(source.IssuerURI, federation.SetIssuerUri)
			result = append(result, federation)
		default:
			return nil, fmt.Errorf("identity_sources[%d]: %w %q", i, errUnsupportedIdentitySourceType, source.ODataType.ValueString())
		}
	}

	return result, nil
}
