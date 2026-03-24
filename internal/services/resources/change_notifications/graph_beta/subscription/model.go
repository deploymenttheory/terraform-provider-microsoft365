package graphBetaChangeNotificationsSubscription

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SubscriptionResourceModel struct {
	ID types.String `tfsdk:"id"`

	ChangeType                types.String `tfsdk:"change_type"`
	NotificationURL           types.String `tfsdk:"notification_url"`
	Resource                  types.String `tfsdk:"resource"`
	ExpirationDateTime        types.String `tfsdk:"expiration_date_time"`
	ClientState               types.String `tfsdk:"client_state"`
	LifecycleNotificationURL  types.String `tfsdk:"lifecycle_notification_url"`
	LatestSupportedTLSVersion types.String `tfsdk:"latest_supported_tls_version"`
	NotificationURLAppID      types.String `tfsdk:"notification_url_app_id"`
	NotificationQueryOptions  types.String `tfsdk:"notification_query_options"`
	NotificationContentType   types.String `tfsdk:"notification_content_type"`
	IncludeResourceData       types.Bool   `tfsdk:"include_resource_data"`
	EncryptionCertificate     types.String `tfsdk:"encryption_certificate"`
	EncryptionCertificateID   types.String `tfsdk:"encryption_certificate_id"`

	ApplicationID types.String `tfsdk:"application_id"`
	CreatorID     types.String `tfsdk:"creator_id"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
