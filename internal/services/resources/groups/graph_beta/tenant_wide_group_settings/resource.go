package graphBetaTenantWideGroupSettings

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
	ResourceName  = "microsoft365_graph_beta_groups_tenant_wide_group_settings"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &TenantWideGroupSettingsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &TenantWideGroupSettingsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &TenantWideGroupSettingsResource{}
)

func NewTenantWideGroupSettingsResource() resource.Resource {
	return &TenantWideGroupSettingsResource{
		ReadPermissions: []string{
			"Group.Read.All",
		},
		WritePermissions: []string{
			"Directory.ReadWrite.All",
			"Group.ReadWrite.All",
		},
		ResourcePath: "/settings",
	}
}

type TenantWideGroupSettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *TenantWideGroupSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *TenantWideGroupSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *TenantWideGroupSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *TenantWideGroupSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages tenant-wide directory settings for Microsoft 365 groups using the `/settings` endpoint. This resource enables a collection of configurations that allow admins to manage behaviors for specific Microsoft Entra objects like Microsoft 365 groups." +
			"This resource applies settings tenant-wide, enabling admins to control various aspects of group functionality." +
			"Use this resource in conjunction with the datasource 'microsoft365_graph_beta_identity_and_access_directory_setting_templates' to get the template_id, settings and values.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the tenant-wide setting. Read-only.",
			},
			"display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Display name of this group of settings, which comes from the associated template. Read-only.",
			},
			"template_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique identifier for the tenant-level directorySettingTemplate object that's been customized for this tenant-level settings object. The template options can be found at 'https://learn.microsoft.com/en-us/graph/group-directory-settings?tabs=http'.",
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
