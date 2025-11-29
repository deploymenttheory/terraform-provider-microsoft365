// REF: https://learn.microsoft.com/en-us/graph/api/user-get-mailboxsettings?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/user-update-mailboxsettings?view=graph-rest-beta
package graphBetaUsersUserMailboxSettings

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_users_user_mailbox_settings"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &UserMailboxSettingsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &UserMailboxSettingsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &UserMailboxSettingsResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &UserMailboxSettingsResource{}
)

func NewUserMailboxSettingsResource() resource.Resource {
	return &UserMailboxSettingsResource{
		ReadPermissions: []string{
			"MailboxSettings.Read",
			"MailboxSettings.ReadWrite",
		},
		WritePermissions: []string{
			"MailboxSettings.ReadWrite",
		},
	}
}

type UserMailboxSettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *UserMailboxSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *UserMailboxSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using the user ID
func (r *UserMailboxSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import using user_id - the ID will be generated from user_id in Read
	resource.ImportStatePassthroughID(ctx, path.Root("user_id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *UserMailboxSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 user mailbox settings using the `/users/{id}/mailboxSettings` endpoint. " +
			"This resource allows you to configure automatic replies, date/time formats, locale, time zone, working hours, " +
			"and other mailbox preferences for a user. Note: This resource manages settings that may also be modified by " +
			"users through Outlook clients. The mailbox settings always exist for a user, so 'create' and 'update' operations " +
			"both use the PATCH method, and 'delete' only removes the resource from Terraform state without affecting the " +
			"actual mailbox settings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Computed identifier for this resource (format: users/{user_id}/mailboxSettings). Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the user whose mailbox settings are being managed. This can be the user's object ID or userPrincipalName.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.Any(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID",
						),
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.UserPrincipalNameRegex),
							"must be a valid user principal name",
						),
					),
				},
			},
			"automatic_replies_setting": schema.SingleNestedAttribute{
				MarkdownDescription: "Configuration for automatic replies (also known as Out of Office or OOF) for the user's mailbox.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"status": schema.StringAttribute{
						MarkdownDescription: "The status of automatic replies. Possible values: `disabled`, `alwaysEnabled`, `scheduled`.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("disabled", "alwaysEnabled", "scheduled"),
						},
					},
					"external_audience": schema.StringAttribute{
						MarkdownDescription: "The audience that will receive external automatic reply messages. Possible values: `none`, `contactsOnly`, `all`.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "contactsOnly", "all"),
						},
					},
					"scheduled_start_date_time": schema.SingleNestedAttribute{
						MarkdownDescription: "The start date and time when automatic replies are scheduled to be sent. Required when status is `scheduled`.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"date_time": schema.StringAttribute{
								MarkdownDescription: "The date and time value in ISO 8601 format (e.g., `2016-03-19T02:00:00.0000000`). The timezone is specified separately in the `time_zone` field.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?$`),
										"must be a valid ISO 8601 datetime format (e.g., 2016-03-14T07:00:00 or 2016-03-14T07:00:00.0000000)",
									),
								},
							},
							"time_zone": schema.StringAttribute{
								MarkdownDescription: "The time zone for the date time value. Defaults to `UTC` if not specified.",
								Optional:            true,
								Computed:            true,
							},
						},
					},
					"scheduled_end_date_time": schema.SingleNestedAttribute{
						MarkdownDescription: "The end date and time when automatic replies are scheduled to stop being sent. Required when status is `scheduled`.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"date_time": schema.StringAttribute{
								MarkdownDescription: "The date and time value in ISO 8601 format (e.g., `2016-03-20T02:00:00.0000000`). The timezone is specified separately in the `time_zone` field.",
								Required:            true,
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?$`),
										"must be a valid ISO 8601 datetime format (e.g., 2016-03-14T07:00:00 or 2016-03-14T07:00:00.0000000)",
									),
								},
							},
							"time_zone": schema.StringAttribute{
								MarkdownDescription: "The time zone for the date time value. Defaults to `UTC` if not specified.",
								Optional:            true,
								Computed:            true,
							},
						},
					},
					"internal_reply_message": schema.StringAttribute{
						MarkdownDescription: "The automatic reply message to send to internal recipients. Supports HTML formatting.",
						Optional:            true,
					},
					"external_reply_message": schema.StringAttribute{
						MarkdownDescription: "The automatic reply message to send to external recipients. Supports HTML formatting.",
						Optional:            true,
					},
				},
			},
			"date_format": schema.StringAttribute{
				MarkdownDescription: "The date format for the user's mailbox. This uses [.NET standard date format patterns](https://learn.microsoft.com/en-us/dotnet/standard/base-types/standard-date-and-time-format-strings#ShortDate) that are culture-specific. Common values include: `M/d/yyyy` (US), `dd/MM/yyyy` (UK/EU), `yyyy-MM-dd` (ISO), `dd.MM.yyyy` (German). The format determines how dates are displayed in the user's mailbox.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						// US formats
						"M/d/yyyy",
						"MM/dd/yyyy",
						// European formats
						"d/M/yyyy",
						"dd/MM/yyyy",
						"dd-MM-yyyy",
						// ISO and Asian formats
						"yyyy-MM-dd",
						"yyyy/MM/dd",
						"yyyy年M月d日", // Japanese
						// German and other dot-separated formats
						"dd.MM.yyyy",
						"d.M.yyyy",
						"M.d.yyyy",
					),
				},
			},
			"delegate_meeting_message_delivery_options": schema.StringAttribute{
				MarkdownDescription: "Specifies how meeting messages and responses are delivered to delegates. Possible values: `sendToDelegateAndInformationToPrincipal`, `sendToDelegateAndPrincipal`, `sendToDelegateOnly`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"sendToDelegateAndInformationToPrincipal",
						"sendToDelegateAndPrincipal",
						"sendToDelegateOnly",
					),
				},
			},
			"language": schema.SingleNestedAttribute{
				MarkdownDescription: "The locale (language and country/region) information for the user.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"locale": schema.StringAttribute{
						MarkdownDescription: "A locale representation for the user, which includes the user's preferred language and country/region. For example, `en-us`. The language component follows 2-letter codes as defined in ISO 639-1, and the country component follows 2-letter codes as defined in ISO 3166-1 alpha-2.",
						Required:            true,
					},
					"display_name": schema.StringAttribute{
						MarkdownDescription: "The display name of the locale. Read-only.",
						Computed:            true,
					},
				},
			},
			"time_format": schema.StringAttribute{
				MarkdownDescription: "The time format for the user's mailbox. This uses [.NET standard time format patterns](https://learn.microsoft.com/en-us/dotnet/standard/base-types/standard-date-and-time-format-strings#ShortTime) that are culture-specific. Common examples include: `h:mm tt` (1:45 PM - US 12-hour), `HH:mm` (13:45 - European 24-hour). The format determines how times are displayed in the user's mailbox.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						// 12-hour formats with AM/PM (en-US style)
						"h:mm tt",
						"hh:mm tt",
						"h:mm:ss tt",
						"hh:mm:ss tt",
						// 24-hour formats (hr-HR, es-ES style)
						"H:mm",
						"HH:mm",
						"H:mm:ss",
						"HH:mm:ss",
					),
				},
			},
			"time_zone": schema.StringAttribute{
				MarkdownDescription: "The default time zone for the user's mailbox. Can be in Windows time zone format or IANA time zone format (Olson format).",
				Optional:            true,
				Computed:            true,
			},
			"working_hours": schema.SingleNestedAttribute{
				MarkdownDescription: "The working hours configured for the user's calendar.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"days_of_week": schema.SetAttribute{
						MarkdownDescription: "The days of the week on which the user works. Possible values: `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`.",
						ElementType:         types.StringType,
						Required:            true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
							),
						},
					},
					"start_time": schema.StringAttribute{
						MarkdownDescription: "The time the user starts working each day, in HH:mm:ss.fffffff format (e.g., `09:00:00.0000000`).",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^\d{2}:\d{2}:\d{2}\.\d{7}$`),
								"must be in format HH:mm:ss.fffffff (e.g., 09:00:00.0000000)",
							),
						},
					},
					"end_time": schema.StringAttribute{
						MarkdownDescription: "The time the user stops working each day, in HH:mm:ss.fffffff format (e.g., `17:00:00.0000000`).",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^\d{2}:\d{2}:\d{2}\.\d{7}$`),
								"must be in format HH:mm:ss.fffffff (e.g., 17:00:00.0000000)",
							),
						},
					},
					"time_zone": schema.SingleNestedAttribute{
						MarkdownDescription: "The time zone for the working hours. Can be a standard or custom time zone.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								MarkdownDescription: "The name of the time zone (e.g., `Pacific Standard Time`, `America/Los_Angeles`).",
								Required:            true,
							},
						},
					},
				},
			},
			"user_purpose": schema.StringAttribute{
				MarkdownDescription: "The purpose of the mailbox. Differentiates a mailbox for a single user from a shared mailbox and equipment mailbox in Exchange Online. Possible values are: user, linked, shared, room, equipment, others, unknownFutureValue. Read-only.",
				Computed:            true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
