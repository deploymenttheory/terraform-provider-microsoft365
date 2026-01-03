package graphBetaServicePrincipalAppRoleAssignedTo

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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
	ResourceName  = "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ServicePrincipalAppRoleAssignedToResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ServicePrincipalAppRoleAssignedToResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ServicePrincipalAppRoleAssignedToResource{}
)

func NewServicePrincipalAppRoleAssignedToResource() resource.Resource {
	return &ServicePrincipalAppRoleAssignedToResource{
		ReadPermissions: []string{
			"Application.ReadWrite.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.All",
			"AppRoleAssignment.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/servicePrincipals/{servicePrincipal-id}/appRoleAssignedTo",
	}
}

type ServicePrincipalAppRoleAssignedToResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *ServicePrincipalAppRoleAssignedToResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ServicePrincipalAppRoleAssignedToResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
// Import ID format: {resource_object_id}/{app_role_assignment_id}
func (r *ServicePrincipalAppRoleAssignedToResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Import ID must be in format: resource_object_id/app_role_assignment_id. Got: %s", req.ID),
		)
		return
	}

	resourceObjectID := parts[0]
	assignmentID := parts[1]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_object_id"), resourceObjectID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), assignmentID)...)
}

// Schema defines the schema for the resource.
func (r *ServicePrincipalAppRoleAssignedToResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages app role assignments granted for a service principal using the `/servicePrincipals/{id}/appRoleAssignedTo` endpoint. " +
			"This resource enables assigning app roles defined by a resource service principal to users, groups, or client service principals.\n\n" +
			"App roles assigned to service principals are also known as **application permissions**. These can be granted directly with app role assignments " +
			"or through a consent experience.\n\n" +
			"To grant an app role assignment, you need three identifiers:\n" +
			"- `target_service_principal_object_id`: The Object ID of the user, group, or client service principal to which you are assigning the app role\n" +
			"- `resource_object_id`: The Object ID of the resource service principal which has defined the app role\n" +
			"- `app_role_id`: The ID of the appRole (defined on the resource service principal) to assign\n\n" +
			"For more information, see the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-approleassignedto?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for the app role assignment.",
			},
			"resource_object_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The Object ID of the service principal that exposes the app roles (permissions). " +
					"This is the API whose permissions you are granting. For Microsoft 365 permissions, this is typically " +
					"the Microsoft Graph service principal (appId: 00000003-0000-0000-c000-000000000000). Other examples " +
					"include SharePoint Online, Exchange Online, or your own custom APIs.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
				},
			},
			"app_role_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The identifier (ID) for the app role which is assigned to the principal. This app " +
					"role must be exposed in the `appRoles` property on the resource application's service principal " +
					"(`resource_object_id`). If the resource application has not declared any app roles, a default app role ID " +
					"of `00000000-0000-0000-0000-000000000000` can be specified to signal that the principal is assigned to the resource app without any specific app roles.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
				},
			},
			"target_service_principal_object_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The Object ID of the service principal being granted the app role. This is the " +
					"enterprise app (service principal) that will receive the permission.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
				},
			},
			"principal_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The type of the assigned principal. This can be either `User`, `Group`, or `ServicePrincipal`. Read-only.",
			},
			"principal_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the user, group, or service principal that was granted the app role assignment. Read-only.",
			},
			"resource_display_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The display name of the resource app's service principal to which the assignment is made. Read-only.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The time when the app role assignment was created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. Read-only.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
