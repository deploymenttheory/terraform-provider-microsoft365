package graphBetaNetworkFilteringProfile

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_filtering_profile"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &NetworkFilteringProfileResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &NetworkFilteringProfileResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &NetworkFilteringProfileResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &NetworkFilteringProfileResource{}
)

// NewNetworkFilteringProfileResource creates a new instance of the NetworkFilteringProfileResource resource.
func NewNetworkFilteringProfileResource() resource.Resource {
	return &NetworkFilteringProfileResource{
		ReadPermissions: []string{
			"NetworkAccess.Read.All",
		},
		WritePermissions: []string{
			"NetworkAccess.ReadWrite.All",
		},
		ResourcePath: "/networkAccess/filteringProfiles",
	}
}

type NetworkFilteringProfileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *NetworkFilteringProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *NetworkFilteringProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *NetworkFilteringProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *NetworkFilteringProfileResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *NetworkFilteringProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra Global Secure Access security profiles, represented in Microsoft Graph beta as filtering profiles, using the `/networkAccess/filteringProfiles` endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the filtering profile.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Global Secure Access security profile.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the Global Secure Access security profile. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The priority used to order the filtering profile for processing. Microsoft Graph beta exposes this as an `Int64` on `microsoft.graph.networkaccess.filteringProfile`.",
				Required:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state of the filtering profile. Possible values are: `enabled`, `disabled`, `unknownFutureValue`. `unknownFutureValue` is accepted for Graph enum compatibility; normal configurations should use `enabled` or `disabled`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled", "unknownFutureValue"),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the filtering profile was created.",
				Computed:            true,
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the filtering profile was last modified.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The Graph version string for the filtering profile.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
