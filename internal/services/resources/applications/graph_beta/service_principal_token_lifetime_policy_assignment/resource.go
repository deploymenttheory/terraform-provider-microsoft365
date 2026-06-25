package graphBetaApplicationsServicePrincipalTokenLifetimePolicyAssignment

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
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_service_principal_token_lifetime_policy_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ServicePrincipalTokenLifetimePolicyAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ServicePrincipalTokenLifetimePolicyAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ServicePrincipalTokenLifetimePolicyAssignmentResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &ServicePrincipalTokenLifetimePolicyAssignmentResource{}
)

func NewServicePrincipalTokenLifetimePolicyAssignmentResource() resource.Resource {
	return &ServicePrincipalTokenLifetimePolicyAssignmentResource{
		ReadPermissions: []string{
			"Policy.Read.All",
			"Application.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.ApplicationConfiguration",
			"Application.ReadWrite.All",
		},
		ResourcePath: "/servicePrincipals/{id}/tokenLifetimePolicies/$ref",
	}
}

type ServicePrincipalTokenLifetimePolicyAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *ServicePrincipalTokenLifetimePolicyAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *ServicePrincipalTokenLifetimePolicyAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using a composite ID: {service_principal_id}/{token_lifetime_policy_id}
func (r *ServicePrincipalTokenLifetimePolicyAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Import ID must be in format: service_principal_id/token_lifetime_policy_id. Got: %s", req.ID),
		)
		return
	}

	servicePrincipalID := parts[0]
	tokenLifetimePolicyID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_principal_id"), servicePrincipalID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("token_lifetime_policy_id"), tokenLifetimePolicyID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), servicePrincipalID+"/"+tokenLifetimePolicyID)...)
}

func (r *ServicePrincipalTokenLifetimePolicyAssignmentResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *ServicePrincipalTokenLifetimePolicyAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the assignment of a token lifetime policy to a service principal using the " +
			"`/servicePrincipals/{id}/tokenLifetimePolicies/$ref` endpoint. " +
			"Only one token lifetime policy can be assigned to a service principal at a time. " +
			"To import this resource, use the format: `service_principal_id/token_lifetime_policy_id`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "A locally generated composite identifier for this assignment in the format `service_principal_id/token_lifetime_policy_id`. " +
					"The Microsoft Graph API does not return an assignment-specific ID for this resource; this value is constructed by the provider " +
					"to uniquely identify the assignment within Terraform state. Use this format when importing the resource.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_principal_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier (Object ID) of the service principal to assign the token lifetime policy to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "service_principal_id must be a valid UUID"),
				},
			},
			"token_lifetime_policy_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the token lifetime policy to assign to the service principal.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "token_lifetime_policy_id must be a valid UUID"),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
