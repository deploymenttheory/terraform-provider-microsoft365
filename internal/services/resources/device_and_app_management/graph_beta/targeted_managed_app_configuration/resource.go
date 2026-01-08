package graphBetaTargetedManagedAppConfigurations

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	sharedschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_targeted_managed_app_configuration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &TargetedManagedAppConfigurationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &TargetedManagedAppConfigurationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &TargetedManagedAppConfigurationResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &TargetedManagedAppConfigurationResource{}
)

func NewTargetedManagedAppConfigurationResource() resource.Resource {
	return &TargetedManagedAppConfigurationResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/targetedManagedAppConfigurations",
	}
}

type TargetedManagedAppConfigurationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *TargetedManagedAppConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *TargetedManagedAppConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *TargetedManagedAppConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *TargetedManagedAppConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages targeted managed app configurations in Microsoft Intune using the `/deviceAppManagement/targetedManagedAppConfigurations` endpoint. Configuration used to deliver a set of custom settings as-is to all users in the targeted security group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Key of the entity",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "Policy display name",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Required. The description of the resource. Maximum length is 10000 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(10000),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "Date and time the policy was created",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last time the policy was modified",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "Version of the entity",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of Scope Tags for this entity",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"custom_settings": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "A set of string key and string value pairs to be sent to apps",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The configuration setting name",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "The configuration setting value",
						},
					},
				},
			},
			"app_group_type": schema.StringAttribute{
				Required: true,
				Description: "Public Apps selection scope. Indicates a collection of apps to target which can be one of several " +
					"pre-defined lists of apps or a manually selected list of apps. Valid values are: `selectedPublicApps`, " +
					"`allCoreMicrosoftApps`, `allMicrosoftApps`, `allApps`",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"selectedPublicApps",
						"allCoreMicrosoftApps",
						"allMicrosoftApps",
						"allApps",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"apps": schema.SetNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of apps to which the policy is deployed",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"mobile_app_identifier": schema.SingleNestedAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The mobile app identifier",
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Required:    true,
									Description: "The type of mobile app identifier. Valid values are: `android_mobile_app`, `ios_mobile_app`, `windows_app`",
									Validators: []validator.String{
										stringvalidator.OneOf(
											"android_mobile_app",
											"ios_mobile_app",
											"windows_app",
										),
									},
								},
								"bundle_id": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "The bundle identifier for iOS apps",
								},
								"package_id": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "The package identifier for Android apps",
								},
								"windows_app_id": schema.StringAttribute{
									Optional:    true,
									Computed:    true,
									Description: "The app identifier for Windows apps",
								},
							},
						},
						"version": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The version of the mobile app",
						},
					},
				},
			},
			"assignments": sharedschema.InclusionGroupAndExclusionGroupAssignmentsSchema(),
			"settings_catalog": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings Catalog (configuration policy) mobile app configuration settings",
				Attributes:          DeviceConfigV2Attributes(),
			},
			"deployed_app_count": schema.Int32Attribute{
				Computed:    true,
				Description: "Count of apps to which the current policy is deployed",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates if the policy is deployed to any inclusion groups or not",
			},
			"targeted_app_management_levels": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Intended app management levels for this configuration. Valid values are: `unspecified`, `unmanaged`, `mdm`, `androidEnterprise`, `androidEnterpriseDedicatedDevicesWithAzureAdSharedMode`, `androidOpenSourceProjectUserAssociated`, `androidOpenSourceProjectUserless`",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"unspecified",
						"unmanaged",
						"mdm",
						"androidEnterprise",
						"androidEnterpriseDedicatedDevicesWithAzureAdSharedMode",
						"androidOpenSourceProjectUserAssociated",
						"androidOpenSourceProjectUserless",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
		},
	}

	resp.Schema.Attributes["timeouts"] = commonschema.ResourceTimeouts(ctx)
}
