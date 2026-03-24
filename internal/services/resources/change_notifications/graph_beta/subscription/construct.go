package graphBetaChangeNotificationsSubscription

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func constructResource(ctx context.Context, data *SubscriptionResourceModel) (graphmodels.Subscriptionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	body := graphmodels.NewSubscription()

	convert.FrameworkToGraphString(data.ChangeType, body.SetChangeType)
	convert.FrameworkToGraphString(data.NotificationURL, body.SetNotificationUrl)
	convert.FrameworkToGraphString(data.Resource, body.SetResource)
	if err := convert.FrameworkToGraphTime(data.ExpirationDateTime, body.SetExpirationDateTime); err != nil {
		return nil, fmt.Errorf("expiration_date_time: %w", err)
	}
	convert.FrameworkToGraphString(data.ClientState, body.SetClientState)
	convert.FrameworkToGraphString(data.LifecycleNotificationURL, body.SetLifecycleNotificationUrl)
	convert.FrameworkToGraphString(data.LatestSupportedTLSVersion, body.SetLatestSupportedTlsVersion)
	convert.FrameworkToGraphString(data.NotificationURLAppID, body.SetNotificationUrlAppId)
	convert.FrameworkToGraphString(data.NotificationQueryOptions, body.SetNotificationQueryOptions)
	convert.FrameworkToGraphBool(data.IncludeResourceData, body.SetIncludeResourceData)
	convert.FrameworkToGraphString(data.EncryptionCertificate, body.SetEncryptionCertificate)
	convert.FrameworkToGraphString(data.EncryptionCertificateID, body.SetEncryptionCertificateId)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), body); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return body, nil
}

// constructPatch builds the body for PATCH /subscriptions/{id}. Only expirationDateTime and notificationUrl are supported.
func constructPatch(
	_ context.Context,
	data *SubscriptionResourceModel,
) (graphmodels.Subscriptionable, error) {
	body := graphmodels.NewSubscription()
	if err := convert.FrameworkToGraphTime(
		data.ExpirationDateTime, body.SetExpirationDateTime,
	); err != nil {
		return nil, fmt.Errorf("expiration_date_time: %w", err)
	}
	convert.FrameworkToGraphString(data.NotificationURL, body.SetNotificationUrl)
	return body, nil
}
