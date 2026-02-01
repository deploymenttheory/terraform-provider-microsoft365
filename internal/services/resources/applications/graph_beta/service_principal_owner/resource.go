package graphBetaServicePrincipalOwner

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
	ResourceName  = "microsoft365_graph_beta_applications_service_principal_owner"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ServicePrincipalOwnerResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ServicePrincipalOwnerResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ServicePrincipalOwnerResource{}
)

func NewServicePrincipalOwnerResource() resource.Resource {
	return &ServicePrincipalOwnerResource{
		ReadPermissions: []string{
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/servicePrincipals",
	}
}

type ServicePrincipalOwnerResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ServicePrincipalOwnerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ServicePrincipalOwnerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *ServicePrincipalOwnerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *ServicePrincipalOwnerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an owner assignment for a Microsoft Entra Service Principal using the `/servicePrincipals/{id}/owners` endpoint. " +
			"Owners are users or service principals who are allowed to modify the service principal object. As a recommended best practice, " +
			"service principals should have at least two owners.\n\n" +
			"For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-owners?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this service principal owner assignment. This is a composite ID formed by combining the service principal ID and owner ID.",
			},
			"service_principal_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the service principal.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"owner_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (UUID) for the owner to be added to the service principal. This can be a user or service principal.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"owner_object_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of object being added as an owner. This determines the correct Microsoft Graph API endpoint to use. Valid values: 'User', 'ServicePrincipal'.",
				Validators: []validator.String{
					stringvalidator.OneOf("User", "ServicePrincipal"),
				},
			},
			"owner_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the owner object as returned by Microsoft Graph (e.g., 'User', 'ServicePrincipal'). Read-only.",
			},
			"owner_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the owner. Read-only.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
