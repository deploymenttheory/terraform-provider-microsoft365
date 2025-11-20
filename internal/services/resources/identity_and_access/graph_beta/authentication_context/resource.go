package graphBetaAuthenticationContext

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_authentication_context"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &AuthenticationContextResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &AuthenticationContextResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &AuthenticationContextResource{}
)

func NewAuthenticationContextResource() resource.Resource {
	return &AuthenticationContextResource{
		ReadPermissions: []string{
			"AuthenticationContext.ReadWrite.All",
			"Policy.ReadWrite.ConditionalAccess",
		},
		WritePermissions: []string{
			"AuthenticationContext.ReadWrite.All",
			"Policy.ReadWrite.ConditionalAccess",
		},
		ResourcePath: "/identity/conditionalAccess/authenticationContextClassReferences",
	}
}

type AuthenticationContextResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *AuthenticationContextResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *AuthenticationContextResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *AuthenticationContextResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema defines the schema for the resource.
func (r *AuthenticationContextResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 Authentication Context Class References using the `/identity/conditionalAccess/authenticationContextClassReferences` endpoint. " +
			"Authentication context is used to trigger step-up authentication in scenarios and applications. " +
			"It allows you to require additional verification when users access sensitive resources or perform sensitive actions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier used to reference authentication context class. Supported values are `c8` through `c99`. " +
					"This value is used to trigger step-up authentication and will be issued in the `acrs` claim of the access token.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^c([8-9]|[1-9][0-9])$`),
						"must be in the format 'c' followed by a number from 8 to 99 (e.g., c8, c10, c99)",
					),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "A friendly name that identifies the authentication context class reference. This value should be used to identify the authentication context class reference when building user-facing admin experiences. For example, a selection UX.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A short explanation of the policies that are enforced by authenticationContextClassReference. This value should be used to provide secondary text to describe the authentication context class reference when building user-facing admin experiences. For example, a selection UX.",
				Optional:            true,
			},
			"is_available": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the authentication context class reference is available for use by apps. The default value is `false`. When `isAvailable` is set to `false`, the authentication context class reference is not shown in the authentication context UX elements and may not be used by applications.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
