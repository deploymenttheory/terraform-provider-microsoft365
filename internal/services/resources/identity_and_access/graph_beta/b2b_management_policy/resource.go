package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	attributevalidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_b2b_management_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &B2bManagementPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &B2bManagementPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &B2bManagementPolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &B2bManagementPolicyResource{}
)

func NewB2bManagementPolicyResource() resource.Resource {
	return &B2bManagementPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.B2BManagementPolicy",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.B2BManagementPolicy",
		},
	}
}

type B2bManagementPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (r *B2bManagementPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *B2bManagementPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *B2bManagementPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *B2bManagementPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *B2bManagementPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra B2B management policies using the `/policies/b2bManagementPolicies` endpoint. " +
			"A B2B management policy controls B2B collaboration settings such as the domains that can be invited.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the B2B management policy.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the B2B management policy.",
				Required:            true,
			},
			"definition": schema.ListAttribute{
				MarkdownDescription: "A collection containing the policy definition as a JSON string, for example " +
					"`{\"B2BManagementPolicy\":{\"InvitationsAllowedAndBlockedDomainsPolicy\":{\"BlockedDomains\":[\"example.com\"]}}}`.",
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(attributevalidator.JSONSchemaValidator()),
				},
			},
			"is_organization_default": schema.BoolAttribute{
				MarkdownDescription: "If `true`, the policy applies tenant-wide as the organization default. Only one policy can be the " +
					"organization default. Microsoft Graph cannot change this on an existing policy, so changing it replaces the policy. " +
					"Defaults to `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
