package graphBetaIOSManagedDeviceAppConfigurationPolicy

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
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_ios_managed_device_app_configuration_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &IOSManagedDeviceAppConfigurationPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &IOSManagedDeviceAppConfigurationPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &IOSManagedDeviceAppConfigurationPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &IOSManagedDeviceAppConfigurationPolicyResource{}

	// Enables resource-level configuration validation
	_ resource.ResourceWithConfigValidators = &IOSManagedDeviceAppConfigurationPolicyResource{}
)

func NewIOSManagedDeviceAppConfigurationPolicyResource() resource.Resource {
	return &IOSManagedDeviceAppConfigurationPolicyResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
		},
		ResourcePath: "/deviceAppManagement/mobileAppConfigurations",
	}
}

type IOSManagedDeviceAppConfigurationPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *IOSManagedDeviceAppConfigurationPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *IOSManagedDeviceAppConfigurationPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *IOSManagedDeviceAppConfigurationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *IOSManagedDeviceAppConfigurationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages iOS managed device mobile app configurations in Microsoft Intune using the `/deviceAppManagement/mobileAppConfigurations` endpoint." +
			"Use app configuration policies in Microsoft Intune to provide custom configuration settings for an iOS/iPadOS app. These configuration settings allow an app  " +
			"to be customized based on the app suppliers direction. You must get these configuration settings (keys and values) from the supplier of the app. " +
			"To configure the app, you specify the settings as keys and values, or as XML containing the keys and values. Learn more here: https://learn.microsoft.com/en-us/intune/intune-service/apps/app-configuration-policies-use-ios",
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
			"version": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "Version of the iOS mobile app configuration.",
			},
			"encoded_setting_xml": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "MDM app configuration in Base64 encoded format. **Note:** This field is mutually exclusive with `settings` - only one can be specified.",
			},
			"settings": schema.SetNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration setting items for the iOS mobile app. **Note:** This field is mutually exclusive with `encoded_setting_xml` - only one can be specified.",
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
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
