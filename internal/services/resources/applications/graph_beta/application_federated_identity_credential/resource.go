package graphBetaApplicationFederatedIdentityCredential

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_application_federated_identity_credential"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ApplicationFederatedIdentityCredentialResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ApplicationFederatedIdentityCredentialResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ApplicationFederatedIdentityCredentialResource{}
)

func NewApplicationFederatedIdentityCredentialResource() resource.Resource {
	return &ApplicationFederatedIdentityCredentialResource{
		ReadPermissions: []string{
			"Application.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/applications",
	}
}

type ApplicationFederatedIdentityCredentialResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ApplicationFederatedIdentityCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ApplicationFederatedIdentityCredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles the import functionality.
func (r *ApplicationFederatedIdentityCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID format: application_id/credential_id
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in format 'application_id/credential_id', got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}

// Schema returns the schema for the resource.
func (r *ApplicationFederatedIdentityCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Federated Identity Credential for an Application in Microsoft Entra ID. This resource is used to configure a trust relationship between your application and the identity provider for your compute platform, enabling token-based authentication with Microsoft identity platform. By configuring this trust relationship, you can use tokens issued by that platform to authenticate and call APIs in the Microsoft ecosystem. Maximum of 20 federated identity credentials can be added to an application.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the federated identity credential. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_id": schema.StringAttribute{
				MarkdownDescription: "The Object ID (id) of the Application to which this federated identity credential belongs. " +
					"This is required and cannot be changed after creation.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the federated identity credential, which has a limit of 120 characters and must be URL friendly. " +
					"It is immutable once created.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 120),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"issuer": schema.StringAttribute{
				MarkdownDescription: "The URL of the external identity provider and must match the issuer claim of the external token being exchanged. " +
					"The combination of the values of issuer and subject must be unique on the app. It has a limit of 600 characters.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(600),
				},
			},
			"subject": schema.StringAttribute{
				MarkdownDescription: "Nullable. Defaults to null if not set. The identifier of the external software workload within the external identity provider. " +
					"Like the audience value, it has no fixed format, as each identity provider uses their own - sometimes a GUID, sometimes a colon delimited identifier, " +
					"sometimes arbitrary strings. The value here must match the sub claim within the token presented to Microsoft Entra ID. " +
					"It has a limit of 600 characters. The combination of issuer and subject must be unique on the app. " +
					"If subject is defined, claims_matching_expression must be null.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(600),
				},
			},
			"audiences": schema.SetAttribute{
				MarkdownDescription: "The audience that can appear in the external token. This field is mandatory and should be set to `api://AzureADTokenExchange` " +
					"for Microsoft Entra ID. It says what Microsoft identity platform should accept in the aud claim in the incoming token. " +
					"This value represents Microsoft Entra ID in your external identity provider and has no fixed value across identity providers - " +
					"you may need to create a new application registration in your identity provider to serve as the audience of this token. " +
					"This field can only accept a single value and has a limit of 600 characters.",
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.LengthAtMost(600)),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "An optional description of the federated identity credential.",
				Optional:            true,
				Computed:            true,
			},
			"claims_matching_expression": schema.StringAttribute{
				MarkdownDescription: "Nullable. Defaults to null if not set. Enables the use of claims matching expressions against specified claims. " +
					"If claims_matching_expression is defined, subject must be null. For the list of supported expression syntax and claims, " +
					"visit the [Flexible FIC reference](https://aka.ms/flexiblefic).",
				Optional: true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
