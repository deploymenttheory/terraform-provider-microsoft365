package auditEvents

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models/managedtenants"
)

// MapRemoteStateToEphemeral maps an audit event to a model for ephemeral resource
func MapRemoteStateToEphemeral(ctx context.Context, data graphmodels.AuditEventable) AuditEventModel {
	tflog.Debug(ctx, "Starting to map remote resource state to ephemeral state", map[string]any{
		"resourceName": data.GetActivity(),
		"resourceId":   data.GetId(),
	})

	model := AuditEventModel{
		ID:                convert.GraphToFrameworkString(data.GetId()),
		Activity:          convert.GraphToFrameworkString(data.GetActivity()),
		ActivityDateTime:  convert.GraphToFrameworkTime(data.GetActivityDateTime()),
		ActivityId:        convert.GraphToFrameworkString(data.GetActivityId()),
		Category:          convert.GraphToFrameworkString(data.GetCategory()),
		HttpVerb:          convert.GraphToFrameworkString(data.GetHttpVerb()),
		InitiatedByAppId:  convert.GraphToFrameworkString(data.GetInitiatedByAppId()),
		InitiatedByUpn:    convert.GraphToFrameworkString(data.GetInitiatedByUpn()),
		InitiatedByUserId: convert.GraphToFrameworkString(data.GetInitiatedByUserId()),
		IpAddress:         convert.GraphToFrameworkString(data.GetIpAddress()),
		RequestUrl:        convert.GraphToFrameworkString(data.GetRequestUrl()),
	}

	// Handle tenant IDs - based on API response, these are single strings
	model.TenantIds = convert.GraphToFrameworkString(data.GetTenantIds())

	// Handle tenant names - based on API response, these are single strings
	model.TenantNames = convert.GraphToFrameworkString(data.GetTenantNames())

	return model
}
