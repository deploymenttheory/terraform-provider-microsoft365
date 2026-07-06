package graphBetaCrossTenantAccessPartnerGroupSyncSettings

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
	_ resource.Resource = &CrossTenantAccessPartnerGroupSyncSettingsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CrossTenantAccessPartnerGroupSyncSettingsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CrossTenantAccessPartnerGroupSyncSettingsResource{}
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_group_sync_settings"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

func NewCrossTenantAccessPartnerGroupSyncSettingsResource() resource.Resource {
	return &CrossTenantAccessPartnerGroupSyncSettingsResource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.CrossTenantAccess",
		},
	}
}

type CrossTenantAccessPartnerGroupSyncSettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *CrossTenantAccessPartnerGroupSyncSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *CrossTenantAccessPartnerGroupSyncSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
// The import ID is the partner tenant ID
func (r *CrossTenantAccessPartnerGroupSyncSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tenant_id"), req, resp)
}

// Schema defines the resource schema.
func (r *CrossTenantAccessPartnerGroupSyncSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages cross-tenant group synchronization settings for a specific partner tenant in Microsoft Entra ID using the `/policies/crossTenantAccessPolicy/partners/{id}/identitySynchronization` endpoint. " +
			"This resource is used to configure whether groups from a partner Microsoft Entra tenant can be synchronized to your tenant.",
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
				MarkdownDescription: "Display name for the cross-tenant group synchronization policy. Use the name of the partner Microsoft Entra tenant to easily identify the policy.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_sync_inbound": schema.SingleNestedAttribute{
				MarkdownDescription: "Determines whether groups are synchronized from the partner tenant to your tenant.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"is_sync_allowed": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether groups can be synchronized from the partner tenant. " +
							"When set to `false`, any current group synchronization from the partner tenant to your tenant will stop. " +
							"This has no impact on existing groups that have already been synchronized.",
						Required: true,
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
