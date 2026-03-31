package graphBetaCrossTenantAccessPartnerUserSyncSettings

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &CrossTenantAccessPartnerUserSyncSettingsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CrossTenantAccessPartnerUserSyncSettingsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CrossTenantAccessPartnerUserSyncSettingsResource{}
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_user_sync_settings"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

func NewCrossTenantAccessPartnerUserSyncSettingsResource() resource.Resource {
	return &CrossTenantAccessPartnerUserSyncSettingsResource{
		ReadPermissions: []string{
			"CrossTenantInformation.ReadBasic.All",
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.CrossTenantAccess",
		},
	}
}

type CrossTenantAccessPartnerUserSyncSettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *CrossTenantAccessPartnerUserSyncSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *CrossTenantAccessPartnerUserSyncSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
// The import ID is the partner tenant ID
func (r *CrossTenantAccessPartnerUserSyncSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tenant_id"), req, resp)
}

// Schema defines the resource schema.
func (r *CrossTenantAccessPartnerUserSyncSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages cross-tenant user synchronization settings for a specific partner tenant in Microsoft Entra ID using the `/policies/crossTenantAccessPolicy/partners/{id}/identitySynchronization` endpoint. " +
			"This resource is used to configure whether users from a partner Microsoft Entra tenant can be synchronized to your tenant.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the identity synchronization policy. This is the same as the `tenant_id`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The tenant ID of the partner Microsoft Entra organization. This is a GUID that uniquely identifies the partner tenant.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name for the cross-tenant user synchronization policy. Use the name of the partner Microsoft Entra tenant to easily identify the policy.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_sync_inbound": schema.SingleNestedAttribute{
				MarkdownDescription: "Determines whether users are synchronized from the partner tenant to your tenant.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"is_sync_allowed": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether users can be synchronized from the partner tenant. " +
							"When set to `false`, any current user synchronization from the partner tenant to your tenant will stop. " +
							"This has no impact on existing users who have already been synchronized.",
						Required: true,
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
