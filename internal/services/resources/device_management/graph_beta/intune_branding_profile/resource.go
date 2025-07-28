package graphBetaDeviceManagementIntuneBrandingProfile

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_intune_branding_profile"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &IntuneBrandingProfileResource{}
	_ resource.ResourceWithConfigure   = &IntuneBrandingProfileResource{}
	_ resource.ResourceWithImportState = &IntuneBrandingProfileResource{}
	_ resource.ResourceWithModifyPlan  = &IntuneBrandingProfileResource{}
)

func NewIntuneBrandingProfileResource() resource.Resource {
	return &IntuneBrandingProfileResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/intuneBrandingProfiles",
	}
}

type IntuneBrandingProfileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *IntuneBrandingProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

func (r *IntuneBrandingProfileResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

func (r *IntuneBrandingProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

func (r *IntuneBrandingProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *IntuneBrandingProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Intune branding profile resource in Intune.\n\n" +
			"## API Documentation\n\n" +
			"- [Graph API Endpoint](https://learn.microsoft.com/en-us/graph/api/resources/intune-wip-intunebrandingprofile?view=graph-rest-beta)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the branding profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"profile_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the branding profile.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"profile_description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the branding profile.",
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Company/organization name that is displayed to end users.",
			},
			"show_logo": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Boolean that represents whether the administrator-supplied logo images are shown or not.",
			},
			"show_display_name_next_to_logo": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Boolean that represents whether the administrator-supplied display name text is shown next to the logo image or not.",
			},
			"contact_it_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Name of the person/organization responsible for IT support.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(40),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.StringRegex),
						"must contain only letters, numbers, and spaces",
					),
				},
			},
			"contact_it_phone_number": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Phone number of the person/organization responsible for IT support.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(20),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9+\-() ]*$`),
						"must contain only numbers, spaces, and the following special characters: +, -, (, )",
					),
				},
			},
			"contact_it_email_address": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Email address of the person/organization responsible for IT support.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(40),
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.EmailRegex), "must be a valid email address"),
				},
			},
			"contact_it_notes": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Text comments regarding the person/organization responsible for IT support.",
			},
			"online_support_site_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "URL to the company/organization's IT helpdesk site.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
						"must be a valid URL",
					),
					stringvalidator.LengthAtMost(250),
				},
			},
			"online_support_site_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Display name of the company/organization's IT helpdesk site.",
			},
			"privacy_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "URL to the company/organization's privacy policy.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
						"must be a valid URL",
					),
				},
			},
			"custom_privacy_message": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Text comments regarding what the admin has access to on the device.",
			},
			"custom_can_see_privacy_message": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Text comments regarding what the admin can see on the device.",
			},
			"custom_cant_see_privacy_message": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Text comments regarding what the admin can't see on the device.",
			},
			"is_remove_device_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Boolean that represents whether the adminstrator has disabled the 'Remove Device' action on corporate owned devices.",
			},
			"is_factory_reset_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Boolean that represents whether the adminstrator has disabled the 'Factory Reset' action on corporate owned devices.",
			},
			"show_azure_ad_enterprise_apps": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Boolean that indicates if AzureAD Enterprise Apps will be shown in Company Portal.",
			},
			"show_office_web_apps": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Boolean that indicates if Office WebApps will be shown in Company Portal.",
			},
			"show_configuration_manager_apps": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Boolean that indicates if Configuration Manager Apps will be shown in Company Portal.",
			},
			"disable_device_category_selection": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Boolean that indicates if device category selection is disabled during enrollment.",
			},
			"send_device_ownership_change_push_notification": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Boolean that indicates if a push notification is sent to users when their device ownership type changes from personal to corporate.",
			},
			"enrollment_availability": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Customized device enrollment flow displayed to the end user. Possible values are: `availableWithPrompts`, `availableWithoutPrompts`, `unavailable`.",
				Validators: []validator.String{
					stringvalidator.OneOf("availableWithPrompts", "availableWithoutPrompts", "unavailable"),
				},
			},
			"disable_client_telemetry": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Applies to telemetry sent from all clients to the Intune service. When disabled, all proactive troubleshooting and issue warnings within the client are turned off, and telemetry settings appear grayed out or hidden to the device user.",
			},
			"is_default_profile": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Boolean that represents whether the profile is used as default or not.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Time when the BrandingProfile was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Time when the BrandingProfile was last modified.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Settings Catalog template profile.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"theme_color": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Primary theme color used in the Company Portal applications and web portal.",
				Attributes: map[string]schema.Attribute{
					"r": schema.Int32Attribute{
						Required:            true,
						MarkdownDescription: "Red value (0-255).",
						Validators: []validator.Int32{
							int32validator.Between(0, 255),
						},
					},
					"g": schema.Int32Attribute{
						Required:            true,
						MarkdownDescription: "Green value (0-255).",
						Validators: []validator.Int32{
							int32validator.Between(0, 255),
						},
					},
					"b": schema.Int32Attribute{
						Required:            true,
						MarkdownDescription: "Blue value (0-255).",
						Validators: []validator.Int32{
							int32validator.Between(0, 255),
						},
					},
				},
			},
			"theme_color_logo":              commonschemagraphbeta.ImageSchema("Logo image displayed in Company Portal apps which have a theme color background behind the logo."),
			"light_background_logo":         commonschemagraphbeta.ImageSchema("Logo image displayed in Company Portal apps which have a light background behind the logo."),
			"landing_page_customized_image": commonschemagraphbeta.ImageSchema("Customized image displayed in Company Portal app landing page."),
			"company_portal_blocked_actions": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Collection of blocked actions on the company portal as per platform and device ownership types.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"platform": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Device platform. Possible values are: `android`, `androidForWork`, `iOS`, `macOS`, `windowsPhone81`, `windows81AndLater`, `windows10AndLater`, `androidWorkProfile`, `unknown`, `androidAOSP`, `androidMobileApplicationManagement`, `iOSMobileApplicationManagement`, `unknownFutureValue`, `windowsMobileApplicationManagement`.",
							Validators: []validator.String{
								stringvalidator.OneOf("android", "androidForWork", "iOS", "macOS", "windowsPhone81", "windows81AndLater", "windows10AndLater", "androidWorkProfile", "unknown", "androidAOSP", "androidMobileApplicationManagement", "iOSMobileApplicationManagement", "unknownFutureValue", "windowsMobileApplicationManagement"),
							},
						},
						"owner_type": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Device ownership type. Possible values are: `unknown`, `company`, `personal`.",
							Validators: []validator.String{
								stringvalidator.OneOf("unknown", "company", "personal"),
							},
						},
						"action": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Action to block. Possible values are: `unknown`, `remove`, `reset`.",
							Validators: []validator.String{
								stringvalidator.OneOf("unknown", "remove", "reset"),
							},
						},
					},
				},
			},
			"assignments": AssignmentBlock(),
			"timeouts":    commonschema.Timeouts(ctx),
		},
	}
}

// AssignmentBlock returns the schema for the assignments block
func AssignmentBlock() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Assignments for the Windows remediation script. Each assignment specifies the target group and schedule for script execution.",
		Optional:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Type of assignment target. Must be one of: 'allDevicesAssignmentTarget', 'allLicensedUsersAssignmentTarget', 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.OneOf(
							"groupAssignmentTarget",
							"exclusionGroupAssignmentTarget",
						),
					},
				},
				"group_id": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "The Entra ID group ID to include or exclude in the assignment. Required when type is 'groupAssignmentTarget' or 'exclusionGroupAssignmentTarget'.",
					Validators: []validator.String{
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
						),
					},
				},
			},
		},
	}
}
