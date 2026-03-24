package graphBetaChangeNotificationsSubscription

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(
	ctx context.Context,
	data *SubscriptionResourceModel,
	remote graphmodels.Subscriptionable,
) {
	if remote == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remote.GetId())
	data.ChangeType = convert.GraphToFrameworkString(remote.GetChangeType())
	data.NotificationURL = convert.GraphToFrameworkString(remote.GetNotificationUrl())
	data.Resource = convert.GraphToFrameworkString(remote.GetResource())
	data.ExpirationDateTime = convert.GraphToFrameworkTime(remote.GetExpirationDateTime())
	data.ClientState = convert.GraphToFrameworkString(remote.GetClientState())
	data.LifecycleNotificationURL = convert.GraphToFrameworkString(
		remote.GetLifecycleNotificationUrl(),
	)
	data.LatestSupportedTLSVersion = convert.GraphToFrameworkString(
		remote.GetLatestSupportedTlsVersion(),
	)
	data.NotificationURLAppID = convert.GraphToFrameworkString(remote.GetNotificationUrlAppId())
	data.NotificationQueryOptions = convert.GraphToFrameworkString(
		remote.GetNotificationQueryOptions(),
	)
	data.NotificationContentType = convert.GraphToFrameworkString(
		remote.GetNotificationContentType(),
	)
	data.IncludeResourceData = convert.GraphToFrameworkBoolWithDefault(
		remote.GetIncludeResourceData(),
		false,
	)
	data.EncryptionCertificate = convert.GraphToFrameworkString(remote.GetEncryptionCertificate())
	data.EncryptionCertificateID = convert.GraphToFrameworkString(
		remote.GetEncryptionCertificateId(),
	)
	data.ApplicationID = convert.GraphToFrameworkString(remote.GetApplicationId())
	data.CreatorID = convert.GraphToFrameworkString(remote.GetCreatorId())
}
