package graphBetaChangeNotificationsSubscription

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_change_notifications_subscription"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &SubscriptionResource{}
	_ resource.ResourceWithConfigure   = &SubscriptionResource{}
	_ resource.ResourceWithImportState = &SubscriptionResource{}
	_ resource.ResourceWithModifyPlan  = &SubscriptionResource{}
	_ resource.ResourceWithIdentity    = &SubscriptionResource{}
)

func NewSubscriptionResource() resource.Resource {
	return &SubscriptionResource{
		ReadPermissions:  subscriptionPermissionHints,
		WritePermissions: subscriptionPermissionHints,
		ResourcePath:     "/subscriptions",
	}
}

type SubscriptionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *SubscriptionResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *SubscriptionResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *SubscriptionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *SubscriptionResource) IdentitySchema(
	ctx context.Context,
	req resource.IdentitySchemaRequest,
	resp *resource.IdentitySchemaResponse,
) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *SubscriptionResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Microsoft Graph [subscription](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta) resource (change notifications / webhooks) on the beta endpoint. Attribute semantics follow that resource type unless noted. " +
			"Creating a subscription requires an HTTPS `notification_url` that completes [notification endpoint validation](https://learn.microsoft.com/en-us/graph/webhooks#notification-endpoint-validation). " +
			"**Permissions** are not generic: they depend on the subscription **resource** path (`resource`), on whether you use delegated or application access, and sometimes on the exact scenario (for example, Copilot [aiInteraction](https://learn.microsoft.com/en-us/microsoft-365/copilot/extensibility/api/ai-services/interaction-export/resources/aiinteraction) paths such as `copilot/users/{userId}/interactionHistory/getAllEnterpriseInteractions` vs `copilot/interactionHistory/getAllEnterpriseInteractions`). " +
			"Use the least-privileged permissions in [Create subscription](https://learn.microsoft.com/en-us/graph/api/subscription-post-subscriptions?view=graph-rest-beta) for your path. Microsoft Graph does not allow granting write-only permissions when read permissions are sufficient. " +
			"Only `expiration_date_time` and `notification_url` can be updated via the Graph API ([Update subscription](https://learn.microsoft.com/en-us/graph/api/subscription-update?view=graph-rest-beta)); changing other attributes forces a new resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Optional. Unique identifier for the subscription. Read-only. See [subscription resource type](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta).",
			},
			"change_type": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				MarkdownDescription: "Required. Indicates the type of change in the subscribed resource that raises a change notification. Supported values are `created`, `updated`, and `deleted`. Multiple values can be combined using a comma-separated list (optional spaces after commas are allowed). " +
					"**Note:** Some resources restrict which change types are valid (for example, drive root item and list change notifications support only `updated`; user and group notifications support `updated` and `deleted` with different semantics). See [changeType](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta) on the subscription resource.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.SubscriptionChangeTypeRegex), "must be a comma-separated list of created, updated, and/or deleted (case-insensitive)"),
				},
			},
			"notification_url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Required. The URL of the endpoint that receives the change notifications. This URL must use the HTTPS protocol. Any query string parameter included in **notificationUrl** is included in the HTTP POST request when Microsoft Graph sends the change notifications. See [notificationUrl](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta).",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.HttpsUrlRegex), "must be a non-empty HTTPS URL without whitespace"),
				},
			},
			"resource": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				MarkdownDescription: "Required. Specifies the resource that is monitored for changes (Graph **resource**). Do not include the base URL (`https://graph.microsoft.com/beta/` or `https://graph.microsoft.com/v1.0/`). See [resource](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta) and [supported resources](https://learn.microsoft.com/en-us/graph/change-notifications-overview#supported-resources). " +
					"This path determines which Microsoft Graph permissions your app needs; see [Create subscription](https://learn.microsoft.com/en-us/graph/api/subscription-post-subscriptions?view=graph-rest-beta). A leading `/` is allowed when the documented path uses it (for example `/communications/presences/{id}`).",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 8192),
					subscriptionResourcePath(),
				},
			},
			"expiration_date_time": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Required. Specifies the date and time when the webhook subscription expires. The time is in UTC. Any value under 45 minutes after the time of the request is automatically extended to 45 minutes after the request time. For maximum supported subscription length per resource, see [Subscription lifetime](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta#subscription-lifetime). " +
					"Use an ISO 8601 / RFC 3339-style timestamp with timezone, for example `2030-01-01T12:00:00Z`.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.ISO8601DateTimeRegex), "must be an ISO 8601 datetime with timezone (e.g. 2030-01-01T12:00:00Z)"),
				},
			},
			"client_state": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				MarkdownDescription: "Optional. Specifies the value of the **clientState** property sent by the service in each change notification. The maximum length is 255 characters. The client can verify that the change notification came from the service by comparing this value with **clientState** on each [changeNotification](https://learn.microsoft.com/en-us/graph/api/resources/changenotification?view=graph-rest-beta). See [clientState](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta) on the subscription resource.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"lifecycle_notification_url": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				MarkdownDescription: "The URL of the endpoint that receives lifecycle notifications, including `subscriptionRemoved`, `reauthorizationRequired`, and `missed` notifications. This URL must use the HTTPS protocol. Required for Teams resources if **expirationDateTime** is more than one hour from now; optional otherwise. See [lifecycleNotificationUrl](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta) and [lifecycle events](https://learn.microsoft.com/en-us/graph/change-notifications-lifecycle-events).",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.HttpsUrlOrEmptyRegex), "must be empty or a non-empty HTTPS URL without whitespace"),
				},
			},
			"latest_supported_tls_version": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				MarkdownDescription: "Optional. Specifies the latest version of Transport Layer Security (TLS) that the notification endpoint, specified by **notificationUrl**, supports. The possible values are: `v1_0`, `v1_1`, `v1_2`, `v1_3`. For endpoints that already support TLS 1.2, setting this property is optional; Microsoft Graph may default it to `v1_2`. See [latestSupportedTlsVersion](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta).",
				Validators: []validator.String{
					stringvalidator.OneOf("v1_0", "v1_1", "v1_2", "v1_3"),
				},
			},
			"notification_url_app_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				MarkdownDescription: "Optional. The app ID that the subscription service can use to generate the validation token. The value allows the client to validate the authenticity of the notification received. See [notificationUrlAppId](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta).",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidOrEmptyValueRegex), "must be empty or a GUID"),
				},
			},
			"notification_query_options": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				MarkdownDescription: "Optional. OData query options for specifying the value for the targeting resource. Supported only for specific workloads such as Universal Print. See [notificationQueryOptions](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta).",
			},
			"include_resource_data": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.NewRequiresReplaceIfChangedBool(),
				},
				MarkdownDescription: "Optional. When set to `true`, change notifications can [include resource data](https://learn.microsoft.com/en-us/graph/change-notifications-with-resource-data). **encryptionCertificate** and **encryptionCertificateId** are required when **includeResourceData** is `true`. See [includeResourceData](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta).",
			},
			"encryption_certificate": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				MarkdownDescription: "Optional. A base64-encoded representation of a certificate with a public key used to encrypt resource data in change notifications. Required when **includeResourceData** is `true`. See [encryptionCertificate](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta).",
			},
			"encryption_certificate_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				MarkdownDescription: "Optional. A custom app-provided identifier to help identify the certificate needed to decrypt resource data. See [encryptionCertificateId](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta).",
			},
			"application_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Optional. Identifier of the application used to create the subscription. Read-only. See [applicationId](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta).",
			},
			"creator_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Optional. Identifier of the user or service principal that created the subscription. If the app used delegated permissions to create the subscription, this field contains the ID of the signed-in user the app called on behalf of. If the app used application permissions, this field contains the ID of the service principal corresponding to the app. Read-only. See [creatorId](https://learn.microsoft.com/en-us/graph/api/resources/subscription?view=graph-rest-beta).",
			},
			"notification_content_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Content type for notification payloads as returned by Microsoft Graph (read-only).",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
