package graphBetaMobileAppDependency

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                   = &MobileAppDependencyResource{}
	_ resource.ResourceWithConfigure      = &MobileAppDependencyResource{}
	_ resource.ResourceWithImportState    = &MobileAppDependencyResource{}
	_ resource.ResourceWithIdentity       = &MobileAppDependencyResource{}
	_ resource.ResourceWithValidateConfig = &MobileAppDependencyResource{}
)

// MobileAppDependencyResource manages mobile app dependency relationships in Microsoft Graph.
type MobileAppDependencyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// NewMobileAppDependencyResource creates a new Mobile App Dependency resource.
func NewMobileAppDependencyResource() resource.Resource {
	return &MobileAppDependencyResource{
		ReadPermissions:  []string{"DeviceManagementApps.Read.All"},
		WritePermissions: []string{"DeviceManagementApps.ReadWrite.All"},
		ResourcePath:     "/deviceAppManagement/mobileApps/{mobileAppId}/relationships",
	}
}

// Metadata returns the resource type name.
func (r *MobileAppDependencyResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

// Configure sets the Microsoft Graph beta client for the resource.
func (r *MobileAppDependencyResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports a mobile app dependency using its relationship ID.
func (r *MobileAppDependencyResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the resource identity used for importing dependencies.
func (r *MobileAppDependencyResource) IdentitySchema(
	ctx context.Context,
	req resource.IdentitySchemaRequest,
	resp *resource.IdentitySchemaResponse,
) {
	resp.IdentitySchema = identityschema.Schema{Attributes: map[string]identityschema.Attribute{
		"id": identityschema.StringAttribute{RequiredForImport: true},
	}}
}

// ValidateConfig ensures that the source and target IDs identify different apps.
func (r *MobileAppDependencyResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var sourceID, targetID types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("source_id"), &sourceID)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("target_id"), &targetID)...)
	if resp.Diagnostics.HasError() || sourceID.IsNull() || sourceID.IsUnknown() ||
		targetID.IsNull() ||
		targetID.IsUnknown() {
		return
	}
	if strings.EqualFold(sourceID.ValueString(), targetID.ValueString()) {
		resp.Diagnostics.AddAttributeError(
			path.Root("target_id"),
			"Invalid app dependency",
			"An app cannot depend on itself. source_id and target_id must identify different apps.",
		)
	}
}

// Schema defines the Terraform schema for the Mobile App Dependency resource.
func (r *MobileAppDependencyResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a dependency between two Microsoft Intune Win32 apps. The source app requires the target app to be detected or automatically installed before its own installation. Each resource manages one dependency and preserves other app relationships.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The relationship identifier assigned by Microsoft Graph. Typically formatted as `<source_id>_<target_id>`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"source_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the Win32 app that requires the dependency.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid app GUID",
					),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"target_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the prerequisite Win32 app that the source app depends on.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid app GUID",
					),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"dependency_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Dependency behavior. `detect` requires the prerequisite to already be installed; `autoInstall` automatically installs it before the source app. Can be changed in place.",
				Validators: []validator.String{
					stringvalidator.OneOf("detect", "autoInstall"),
				},
			},
			"target_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The target's direction in the relationship, normally `child` for a prerequisite app.",
			},
			"source_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the source app.",
			},
			"source_display_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display version of the source app.",
			},
			"source_publisher_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publisher display name of the source app.",
			},
			"target_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the prerequisite app.",
			},
			"target_display_version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display version of the prerequisite app.",
			},
			"target_publisher": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publisher of the prerequisite app.",
			},
			"target_publisher_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The publisher display name of the prerequisite app.",
			},
			"dependent_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The number of apps that directly or indirectly depend on the parent app, as reported by Graph.",
			},
			"depends_on_app_count": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The number of apps the child app directly or indirectly depends on, as reported by Graph.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
