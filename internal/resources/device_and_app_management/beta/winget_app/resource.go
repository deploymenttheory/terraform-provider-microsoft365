package graphBetaWinGetApp

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	graphBetaMobileAppAssignment "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/mobile_app_assignment"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var _ resource.Resource = &WinGetAppResource{}
var _ resource.ResourceWithConfigure = &WinGetAppResource{}
var _ resource.ResourceWithImportState = &WinGetAppResource{}

func NewWinGetAppResource() resource.Resource {
	return &WinGetAppResource{
		ReadPermissions: []string{
			"DeviceManagementApps.Read.All",
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementApps.ReadWrite.All",
			"DeviceManagementConfiguration.ReadWrite.All",
		},
	}
}

type WinGetAppResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// GetID returns the ID of a resource from the state model.
func (s *WinGetAppResourceModel) GetID() string {
	return s.ID.ValueString()
}

// GetTypeName returns the type name of the resource from the state model.
func (r *WinGetAppResource) GetTypeName() string {
	return r.TypeName
}

// Metadata returns the resource type name.
func (r *WinGetAppResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_beta_device_and_app_management_win_get_app"
}

// Configure sets the client for the resource.
func (r *WinGetAppResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *WinGetAppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WinGetAppResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Create an instance of MobileAppAssignmentResource
	mobileAppAssignmentResource := graphBetaMobileAppAssignment.NewMobileAppAssignmentResource()

	// Get the schema for MobileAppAssignmentResource
	var mobileAppAssignmentResp resource.SchemaResponse
	mobileAppAssignmentResource.Schema(ctx, resource.SchemaRequest{}, &mobileAppAssignmentResp)

	resp.Schema = schema.Schema{
		Description: "Manages a WinGet application in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				Description: "The unique graph guid that identifies this resource." +
					"Assigned at time of resource creation. This property is read-only.",
			},
			"package_identifier": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The **unique identifier** for the WinGet/Microsoft Store app from the storefront.\n\n" +
					"For example, for the URL [https://apps.microsoft.com/detail/9nzvdkpmr9rd?hl=en-us&gl=US](https://apps.microsoft.com/detail/9nzvdkpmr9rd?hl=en-us&gl=US), " +
					"the package identifier is `9nzvdkpmr9rd`.\n\n" +
					"**Important notes:**\n" +
					"- This identifier is **required** at creation time.\n" +
					"- It **cannot be modified** for existing Terraform-deployed WinGet/Microsoft Store apps.\n\n" +
					"**Steps to find your package identifier:**\n" +
					"1. Navigate to the [Microsoft Store for Business](https://businessstore.microsoft.com/) portal.\n" +
					"2. Locate the desired WinGet app.\n" +
					"3. In the app details page, identify the package identifier in the URL or within the app information section.\n\n" +
					"**Example usage:**\n" +
					"```hcl\n" +
					"resource \"m365_graph_beta_device_and_app_management_win_get_app\" \"example\" {\n" +
					"  package_identifier = \"9nzvdkpmr9rd\"\n" +
					"\n" +
					"  install_experience {\n" +
					"    run_as_account = \"user\"\n" +
					"  }\n" +
					"}\n" +
					"```\n\n",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Z0-9]{12}$`),
						"package_identifier must be a 12-character string containing only uppercase letters and numbers.",
					),
				},
			},
			"is_featured": schema.BoolAttribute{
				Optional:    true,
				Description: "The value indicating whether the app is marked as featured by the admin.",
			},
			"privacy_information_url": schema.StringAttribute{
				Optional:    true,
				Description: "The privacy statement Url.",
			},
			"information_url": schema.StringAttribute{
				Optional:    true,
				Description: "The more information Url.",
			},
			"owner": schema.StringAttribute{
				Optional:    true,
				Description: "The owner of the app.",
			},
			"developer": schema.StringAttribute{
				Optional:    true,
				Description: "The developer of the app.",
			},
			"notes": schema.StringAttribute{
				Optional:    true,
				Description: "Notes for the app.",
			},
			"display_name": schema.StringAttribute{
				Computed: true,
				Description: "The title of the WinGet app imported from the Microsoft Store for Business." +
					"This field is automatically populated based on the package identifier.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "A detailed description of the WinGet app, automatically retrieved from the Microsoft Store for Business.",
			},
			"publisher": schema.StringAttribute{
				Computed:    true,
				Description: "The publisher of the WinGet app, automatically fetched from the Microsoft Store for Business.",
			},
			"install_experience": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"run_as_account": schema.StringAttribute{
						Required: true,
						Description: "The account type (System or User) that actions should be run as on target devices. " +
							"Required at creation time.",
						Validators: []validator.String{
							stringvalidator.OneOf("system", "user"),
						},
					},
				},
				Description: "The install experience settings associated with this application.",
			},
			"large_icon": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "The MIME type of the app's large icon, automatically determined based on the downloaded image from the Microsoft Store for Business.",
					},
					"value": schema.StringAttribute{
						Computed:    true,
						Description: "The raw byte data of the app's large icon, automatically downloaded from the Microsoft Store for Business.",
					},
				},
				Description: "The large icon for the WinGet app, automatically downloaded and set from the Microsoft Store for Business.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the app was created. This property is read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time the app was last modified. This property is read-only.",
			},
			"upload_state": schema.Int64Attribute{
				Computed:    true,
				Description: "The upload state. Possible values are: 0 - Not Ready, 1 - Ready, 2 - Processing. This property is read-only.",
			},
			"publishing_state": schema.StringAttribute{
				Computed: true,
				Description: "The publishing state for the app. The app cannot be assigned unless the app is published. " +
					"Possible values are: notPublished, processing, published.",
			},
			"is_assigned": schema.BoolAttribute{
				Computed:    true,
				Description: "The value indicating whether the app is assigned to at least one group. This property is read-only.",
			},
			"role_scope_tag_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of scope tag ids for this mobile app.",
			},
			"dependent_app_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total number of dependencies the child app has. This property is read-only.",
			},
			"superseding_app_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total number of apps this app directly or indirectly supersedes. This property is read-only.",
			},
			"superseded_app_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total number of apps this app is directly or indirectly superseded by. This property is read-only.",
			},
			"manifest_hash": schema.StringAttribute{
				Computed:    true,
				Description: "Hash of package metadata properties used to validate that the application matches the metadata in the source repository.",
			},
			"assignments": schema.ListNestedAttribute{
				Optional:    true,
				Description: "The list of group assignments for this mobile app.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: mobileAppAssignmentResp.Schema.Attributes,
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
