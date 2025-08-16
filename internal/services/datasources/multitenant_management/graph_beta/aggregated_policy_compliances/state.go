package graphBetaAggregatedPolicyCompliances

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models/managedtenants"
)

// MapRemoteStateToDataSource maps an aggregated policy compliance to a model
func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.AggregatedPolicyComplianceable) AggregatedPolicyComplianceModel {
	tflog.Debug(ctx, "Starting to map remote resource state to Terraform state", map[string]interface{}{
		"resourceName": data.GetCompliancePolicyName(),
		"resourceId":   data.GetId(),
	})

	model := AggregatedPolicyComplianceModel{
		ID:                              convert.GraphToFrameworkString(data.GetId()),
		CompliancePolicyId:              convert.GraphToFrameworkString(data.GetCompliancePolicyId()),
		CompliancePolicyName:            convert.GraphToFrameworkString(data.GetCompliancePolicyName()),
		CompliancePolicyPlatform:        convert.GraphToFrameworkString(data.GetCompliancePolicyPlatform()),
		CompliancePolicyType:            convert.GraphToFrameworkString(data.GetCompliancePolicyType()),
		LastRefreshedDateTime:           convert.GraphToFrameworkTime(data.GetLastRefreshedDateTime()),
		NumberOfCompliantDevices:        convert.GraphToFrameworkInt64(data.GetNumberOfCompliantDevices()),
		NumberOfErrorDevices:            convert.GraphToFrameworkInt64(data.GetNumberOfErrorDevices()),
		NumberOfNonCompliantDevices:     convert.GraphToFrameworkInt64(data.GetNumberOfNonCompliantDevices()),
		PolicyModifiedDateTime:          convert.GraphToFrameworkTime(data.GetPolicyModifiedDateTime()),
		TenantDisplayName:               convert.GraphToFrameworkString(data.GetTenantDisplayName()),
		TenantId:                        convert.GraphToFrameworkString(data.GetTenantId()),
	}

	return model
}