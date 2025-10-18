package graphBetaGroupSettings

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_groups_group_settings"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupSettingsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupSettingsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupSettingsResource{}
)

func NewGroupSettingsResource() resource.Resource {
	return &GroupSettingsResource{
		ReadPermissions: []string{
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/groups/{group-id}/settings",
	}
}

type GroupSettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *GroupSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *GroupSettingsResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *GroupSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *GroupSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *GroupSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages group-specific directory settings for Microsoft 365 groups using the `/groups/{group-id}/settings` endpoint." +
			"This resource enables configuration of group-level settings such as guest access permissions and other group-specific policies that override tenant-wide defaults." +
			"Use this resource in conjunction with the datasource 'microsoft365_graph_beta_identity_and_access_directory_setting_templates' to get the template_id, settings and values." +
			"Use this resource in conjection with the resource 'microsoft365_graph_beta_groups_group' to get the group_id.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the group setting. Read-only.",
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier of the group for which the settings apply.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Display name of this group of settings, which comes from the associated template. Read-only.",
			},
			"template_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique identifier for the tenant-level directorySettingTemplate object that's been customized for this group-level settings object. The template named 'Group.Unified.Guest' can be used to configure group-specific settings.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
					ValidateTemplateID(ctx, r.client),
				},
			},
			"values": schema.SetNestedAttribute{
				Required:            true,
				MarkdownDescription: "Collection of name-value pairs corresponding to the name and defaultValue properties in the referenced directorySettingTemplate object.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Name of the setting from the referenced directorySettingTemplate.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.StringRegex),
									"must be a valid setting name string",
								),
							},
						},
						"value": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Value of the setting.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.StringRegex),
									"must be a valid setting value string",
								),
							},
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
