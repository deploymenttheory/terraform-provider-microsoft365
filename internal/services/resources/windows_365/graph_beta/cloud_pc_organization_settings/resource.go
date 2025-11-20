package graphBetaCloudPcOrganizationSettings

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_windows_365_cloud_pc_organization_settings"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
	SingletonID   = "singleton"
)

var (
	_ resource.Resource                = &CloudPcOrganizationSettingsResource{}
	_ resource.ResourceWithConfigure   = &CloudPcOrganizationSettingsResource{}
	_ resource.ResourceWithImportState = &CloudPcOrganizationSettingsResource{}
)

func NewCloudPcOrganizationSettingsResource() resource.Resource {
	return &CloudPcOrganizationSettingsResource{
		ReadPermissions:  []string{"CloudPC.Read.All"},
		WritePermissions: []string{"CloudPC.ReadWrite.All"},
	}
}

type CloudPcOrganizationSettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (r *CloudPcOrganizationSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *CloudPcOrganizationSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *CloudPcOrganizationSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CloudPcOrganizationSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the singleton Cloud PC organization settings for a tenant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The singleton ID for the Cloud PC organization settings.",
			},
			"enable_mem_auto_enroll": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Specifies whether new Cloud PCs will be automatically enrolled in Microsoft Endpoint Manager (MEM). The default value is false.",
			},
			"enable_single_sign_on": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "True if the provisioned Cloud PC can be accessed by single sign-on. False indicates that the provisioned Cloud PC doesn't support this feature. " +
					"Default value is false. Windows 365 users can use single sign-on to authenticate to Microsoft Entra ID with passwordless options (for example, FIDO keys) to access their Cloud PC. Optional.",
			},
			"os_version": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The version of the operating system (OS) to provision on Cloud PCs. Possible values: windows10, windows11, unknownFutureValue.",
				Validators: []validator.String{
					stringvalidator.OneOf("windows10", "windows11", "unknownFutureValue"),
				},
			},
			"user_account_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The account type of the user on provisioned Cloud PCs. Possible values: standardUser, administrator, unknownFutureValue.",
				Validators: []validator.String{
					stringvalidator.OneOf("standardUser", "administrator", "unknownFutureValue"),
				},
			},
			"windows_settings": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Represents the Cloud PC organization settings for a tenant. A tenant has only one cloudPcOrganizationSettings object. The default language value en-US.",
				Attributes: map[string]schema.Attribute{
					"language": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("en-US"),
						MarkdownDescription: "The Windows language/region tag to use for language pack configuration and localization of the Cloud PC. The default value is en-US, which corresponds to English (United States).",
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
