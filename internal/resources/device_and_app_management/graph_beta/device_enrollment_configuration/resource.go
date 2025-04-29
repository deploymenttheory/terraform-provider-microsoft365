package graphBetaDeviceEnrollmentConfiguration

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_and_app_management_device_enrollment_configuration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceEnrollmentConfigurationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceEnrollmentConfigurationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceEnrollmentConfigurationResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceEnrollmentConfigurationResource{}
)

func NewDeviceEnrollmentConfigurationResource() resource.Resource {
	return &DeviceEnrollmentConfigurationResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceEnrollmentConfigurations",
	}
}

type DeviceEnrollmentConfigurationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *DeviceEnrollmentConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *DeviceEnrollmentConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *DeviceEnrollmentConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *DeviceEnrollmentConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Device Enrollment Configuration in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The Identifier of the entity.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of the device enrollment configuration.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The description of the device enrollment configuration.",
			},
			"priority": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Priority is used when a user exists in multiple groups that are assigned enrollment configuration. Users are subject only to the configuration with the lowest priority value.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date time that the device enrollment configuration was created.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date time that the device enrollment configuration was last modified.",
			},
			"version": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The version of the device enrollment configuration.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of Scope Tags for this entity.",
			},
			"device_enrollment_configuration_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The type of the device enrollment configuration. Possible values are: unknown, limit, platformRestrictions, windowsHelloForBusiness, defaultLimit, defaultPlatformRestrictions, defaultWindowsHelloForBusiness, defaultWindows10EnrollmentCompletionPageConfiguration, windows10EnrollmentCompletionPageConfiguration, deviceComanagementAuthorityConfiguration, singlePlatformRestriction, enrollmentNotificationsConfiguration.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"unknown",
						"limit",
						"platformRestrictions",
						"windowsHelloForBusiness",
						"defaultLimit",
						"defaultPlatformRestrictions",
						"defaultWindowsHelloForBusiness",
						"defaultWindows10EnrollmentCompletionPageConfiguration",
						"windows10EnrollmentCompletionPageConfiguration",
						"deviceComanagementAuthorityConfiguration",
						"singlePlatformRestriction",
						"enrollmentNotificationsConfiguration",
					),
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"platform_restriction": schema.SingleNestedBlock{
				MarkdownDescription: "Platform specific enrollment restrictions",
				Attributes: map[string]schema.Attribute{
					"platform_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Block the platform from enrolling.",
					},
					"personal_device_enrollment_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Block personally owned devices from enrolling.",
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum version of the platform.",
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum version of the platform.",
					},
					"blocked_manufacturers": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "List of blocked manufacturers.",
					},
				},
			},
			// Add this block to your resource.go file, in the Schema method
			// under the "Blocks" definition

			"windows_hello_for_business": schema.SingleNestedBlock{
				MarkdownDescription: "Windows Hello for Business settings for Windows 10 and later devices.",
				Attributes: map[string]schema.Attribute{
					"state": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Windows Hello for Business state. Possible values: notConfigured, enabled, disabled.",
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
					},
					"enhanced_biometrics_state": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Enhanced biometrics state. Possible values: notConfigured, enabled, disabled.",
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
					},
					"security_key_for_sign_in": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Security key for sign-in. Possible values: notConfigured, enabled, disabled.",
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
					},
					"pin_lowercase_characters_usage": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "PIN lowercase characters usage. Possible values: allowed, required, disallowed.",
						Validators: []validator.String{
							stringvalidator.OneOf("allowed", "required", "disallowed"),
						},
					},
					"pin_uppercase_characters_usage": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "PIN uppercase characters usage. Possible values: allowed, required, disallowed.",
						Validators: []validator.String{
							stringvalidator.OneOf("allowed", "required", "disallowed"),
						},
					},
					"pin_special_characters_usage": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "PIN special characters usage. Possible values: allowed, required, disallowed.",
						Validators: []validator.String{
							stringvalidator.OneOf("allowed", "required", "disallowed"),
						},
					},
					"enhanced_sign_in_security": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Setting to configure Enhanced sign-in security. Default is Not Configured.",
					},
					"pin_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum number of characters required for Windows Hello for Business PIN. Value must be between 4 and 127, inclusive.",
						Validators: []validator.Int64{
							int64validator.Between(4, 127),
						},
					},
					"pin_maximum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Maximum number of characters allowed for Windows Hello for Business PIN. Value must be between 4 and 127, inclusive.",
						Validators: []validator.Int64{
							int64validator.Between(4, 127),
						},
					},
					"pin_expiration_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the PIN expires. Value must be between 0 and 730, inclusive. If set to 0, the PIN never expires.",
						Validators: []validator.Int64{
							int64validator.Between(0, 730),
						},
					},
					"pin_previous_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous PINs that cannot be reused. Value must be between 0 and 50, inclusive.",
						Validators: []validator.Int64{
							int64validator.Between(0, 50),
						},
					},
					"remote_passport_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Controls the use of Remote Windows Hello for Business.",
					},
					"security_device_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Controls whether to require a Trusted Platform Module (TPM) for provisioning Windows Hello for Business.",
					},
					"unlock_with_biometrics_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Controls the use of biometric gestures, such as face and fingerprint, as an alternative to the Windows Hello for Business PIN.",
					},
				},
			},
			"assignment": commonschemagraphbeta.WindowsUpdateAssignments(),
		},
	}
}
