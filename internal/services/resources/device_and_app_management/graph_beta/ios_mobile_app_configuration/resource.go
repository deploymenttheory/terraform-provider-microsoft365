package graphBetaIOSMobileAppConfiguration

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	ResourceName  = "graph_beta_device_and_app_management_ios_mobile_app_configuration"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &IOSMobileAppConfigurationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &IOSMobileAppConfigurationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &IOSMobileAppConfigurationResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &IOSMobileAppConfigurationResource{}
)

func NewIOSMobileAppConfigurationResource() resource.Resource {
	return &IOSMobileAppConfigurationResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/mobileAppConfigurations",
	}
}

type IOSMobileAppConfigurationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *IOSMobileAppConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *IOSMobileAppConfigurationResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *IOSMobileAppConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *IOSMobileAppConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *IOSMobileAppConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages iOS mobile app configurations in Microsoft Intune using the `/deviceAppManagement/mobileAppConfigurations` endpoint. iOS mobile app configurations enable administrators to configure managed iOS apps with specific settings and values.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this iOS mobile app configuration",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the iOS mobile app configuration",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The optional description of the iOS mobile app configuration",
			},
			"targeted_mobile_apps": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Set of mobile app IDs that this configuration targets.",
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"Must be a valid GUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
						),
					),
				},
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this iOS mobile app configuration.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"version": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Version of the iOS mobile app configuration.",
			},
			"encoded_setting_xml": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "MDM app configuration in Base64 encoded format.",
			},
			"settings": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration setting items for the iOS mobile app.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"app_config_key": schema.StringAttribute{
							Required:    true,
							Description: "The configuration key name",
						},
						"app_config_key_type": schema.StringAttribute{
							Required:    true,
							Description: "The configuration key type (e.g., \"stringType\", \"integerType\", \"realType\", \"booleanType\", \"tokenType\")",
							Validators: []validator.String{
								stringvalidator.OneOf(
									"stringType",
									"integerType",
									"realType",
									"booleanType",
									"tokenType",
								),
							},
						},
						"app_config_key_value": schema.StringAttribute{
							Required:    true,
							Description: "The configuration key value",
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
