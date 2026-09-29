package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

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
	ResourceName  = "microsoft365_graph_beta_identity_and_access_b2b_management_policy_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &B2bManagementPolicyAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &B2bManagementPolicyAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &B2bManagementPolicyAssignmentResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &B2bManagementPolicyAssignmentResource{}
)

func NewB2bManagementPolicyAssignmentResource() resource.Resource {
	return &B2bManagementPolicyAssignmentResource{
		ReadPermissions: []string{
			"Policy.Read.B2BManagementPolicy",
			"Application.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.B2BManagementPolicy",
			"Application.ReadWrite.OwnedBy",
		},
	}
}

type B2bManagementPolicyAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

func (r *B2bManagementPolicyAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *B2bManagementPolicyAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using a composite ID: {b2b_management_policy_id}/{directory_object_id}
func (r *B2bManagementPolicyAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Import ID must be in format: b2b_management_policy_id/directory_object_id. Got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("b2b_management_policy_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("directory_object_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), compositeID(parts[0], parts[1]))...)
}

func (r *B2bManagementPolicyAssignmentResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *B2bManagementPolicyAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Applies a Microsoft Entra B2B management policy to an application or service principal. " +
			"Import ID format: `b2b_management_policy_id/directory_object_id`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Composite ID `b2b_management_policy_id/directory_object_id`; Microsoft Graph has no assignment ID.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"b2b_management_policy_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the B2B management policy to assign.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "b2b_management_policy_id must be a valid UUID"),
				},
			},
			"directory_object_id": schema.StringAttribute{
				MarkdownDescription: "The object ID of the application or service principal to apply the policy to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "directory_object_id must be a valid UUID"),
				},
			},
			"directory_object_type": schema.StringAttribute{
				MarkdownDescription: "The detected type of `directory_object_id`: `application` or `servicePrincipal`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
