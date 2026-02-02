package graphBetaApplicationIdentifierUri

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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

const (
	ResourceName  = "microsoft365_graph_beta_applications_application_identifier_uri"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ApplicationIdentifierUriResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ApplicationIdentifierUriResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ApplicationIdentifierUriResource{}
)

func NewApplicationIdentifierUriResource() resource.Resource {
	return &ApplicationIdentifierUriResource{
		ReadPermissions: []string{
			"Application.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
		},
		ResourcePath: "/applications",
	}
}

type ApplicationIdentifierUriResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ApplicationIdentifierUriResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ApplicationIdentifierUriResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
// Import format: application_id/identifier_uri
func (r *ApplicationIdentifierUriResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use SplitN with limit 2 to handle identifier URIs that contain "/" (e.g., api://guid/path)
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'application_id/identifier_uri', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("identifier_uri"), parts[1])...)
}

// Schema defines the schema for the resource.
func (r *ApplicationIdentifierUriResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an identifier URI for a Microsoft Entra Application. " +
			"Identifier URIs are user-defined URIs that uniquely identify an application within its Azure AD tenant, " +
			"or within a verified custom domain if the application is multi-tenant. " +
			"This resource manages individual identifier URIs, allowing multiple URIs to be managed independently.\n\n" +
			"For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for this resource, in the format `application_id/identifier_uri`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier (UUID) for the application.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"identifier_uri": schema.StringAttribute{
				MarkdownDescription: "The user-defined URI or URI-like string that uniquely identifies an application. " +
					"Valid formats include `api://<guid>`, `api://<domain>/<path>`, `https://<domain>/<path>`, or custom URIs. " +
					"Note: Extensive validation is intentionally avoided as accepted values vary depending on tenant configuration and application type.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
