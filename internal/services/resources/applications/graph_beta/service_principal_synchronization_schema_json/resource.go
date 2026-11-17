package graphBetaServicePrincipalSynchronizationSchemaJson

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_applications_service_principal_synchronization_schema_json"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ServicePrincipalSynchronizationSchemaJsonResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ServicePrincipalSynchronizationSchemaJsonResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ServicePrincipalSynchronizationSchemaJsonResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &ServicePrincipalSynchronizationSchemaJsonResource{}
)

// https://learn.microsoft.com/en-us/graph/api/synchronization-synchronizationschema-update?view=graph-rest-beta
func NewServicePrincipalSynchronizationSchemaJsonResource() resource.Resource {
	return &ServicePrincipalSynchronizationSchemaJsonResource{
		ReadPermissions: []string{
			"Synchronization.Read.All",
		},
		WritePermissions: []string{
			"Application.ReadWrite.OwnedBy",
			"Synchronization.ReadWrite.All",
		},
	}
}

type ServicePrincipalSynchronizationSchemaJsonResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *ServicePrincipalSynchronizationSchemaJsonResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ServicePrincipalSynchronizationSchemaJsonResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state from an ID in the format {service_principal_id}/{job_id}.
func (r *ServicePrincipalSynchronizationSchemaJsonResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	servicePrincipalId, jobId, ok := strings.Cut(req.ID, "/")
	if !ok || servicePrincipalId == "" || jobId == "" || strings.Contains(jobId, "/") {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			fmt.Sprintf("Import ID must be in the format service_principal_id/job_id. Got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_principal_id"), servicePrincipalId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("job_id"), jobId)...)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *ServicePrincipalSynchronizationSchemaJsonResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *ServicePrincipalSynchronizationSchemaJsonResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the synchronization (provisioning) schema of a service principal's synchronization job using the " +
			"`/servicePrincipals/{id}/synchronization/jobs/{jobId}/schema` endpoint. The schema defines which objects are provisioned, " +
			"the attribute mappings and scoping filters, and the directory definitions of the source and target systems. " +
			"The schema is sent in full with `PUT` on every create and update, so the configuration must contain the complete " +
			"`synchronizationRules` and `directories` collections. Destroying the resource calls `DELETE` on the schema, which resets it " +
			"to the default schema of the job's synchronization template.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the resource in the format `{service_principal_id}/{job_id}`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_principal_id": schema.StringAttribute{
				MarkdownDescription: "The object ID of the service principal that owns the synchronization job. Changing it replaces the resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'"),
				},
			},
			"job_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the synchronization job whose schema is managed, for example `scim.{tenantId}.{appId}`. Changing it replaces the resource.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[^/]+$`), "must not contain '/'"),
				},
			},
			"synchronization_rules_json": schema.StringAttribute{
				MarkdownDescription: "The complete `synchronizationRules` collection of the schema as a JSON array, for example from `jsonencode()` or `file()`. " +
					"Each rule defines the source and target directories, object mappings, attribute mappings and scoping filters. " +
					"The value is compared semantically: whitespace and property order are ignored, and so are the `Unsynchronized` object mapping " +
					"metadata entries that Microsoft Graph derives from `enabled`.",
				Required:   true,
				CustomType: SchemaJSONType{},
				Validators: []validator.String{
					jsonArrayValidator{},
				},
			},
			"directories_json": schema.StringAttribute{
				MarkdownDescription: "The complete `directories` collection of the schema as a JSON array, for example from `jsonencode()` or `file()`. " +
					"Each directory definition lists the objects and attributes that can take part in the synchronization rules. " +
					"The value is compared semantically: whitespace and property order are ignored, and so are the `version` properties that " +
					"Microsoft Graph updates on every schema change.",
				Required:   true,
				CustomType: SchemaJSONType{},
				Validators: []validator.String{
					jsonArrayValidator{},
				},
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The version of the schema, updated by Microsoft Graph with every schema change.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
