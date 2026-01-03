package graphBetaTenantInformation

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.TenantInformationable) TenantInformationDataSourceModel {
	model := TenantInformationDataSourceModel{
		TenantID:            convert.GraphToFrameworkString(data.GetTenantId()),
		DisplayName:         convert.GraphToFrameworkString(data.GetDisplayName()),
		DefaultDomainName:   convert.GraphToFrameworkString(data.GetDefaultDomainName()),
		FederationBrandName: convert.GraphToFrameworkString(data.GetFederationBrandName()),
	}

	return model
}
