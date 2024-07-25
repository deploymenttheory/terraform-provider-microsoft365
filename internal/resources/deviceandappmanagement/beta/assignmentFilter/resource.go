// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-policyset-deviceandappmanagementassignmentfilter?view=graph-rest-beta
package graphBetaAssignmentFilter

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

var _ resource.Resource = &AssignmentFilterResource{}
var _ resource.ResourceWithConfigure = &AssignmentFilterResource{}
var _ resource.ResourceWithImportState = &AssignmentFilterResource{}

func NewAssignmentFilterResource() resource.Resource {
	return &AssignmentFilterResource{}
}

type AssignmentFilterResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
}

type AssignmentFilterResourceModel struct {
	ID                             types.String   `tfsdk:"id"`
	DisplayName                    types.String   `tfsdk:"display_name"`
	Description                    types.String   `tfsdk:"description"`
	Platform                       types.String   `tfsdk:"platform"`
	Rule                           types.String   `tfsdk:"rule"`
	AssignmentFilterManagementType types.String   `tfsdk:"assignment_filter_management_type"`
	CreatedDateTime                types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime           types.String   `tfsdk:"last_modified_date_time"`
	RoleScopeTags                  types.List     `tfsdk:"role_scope_tags"`
	Timeouts                       timeouts.Value `tfsdk:"timeouts"`
}

// GetID returns the ID of a resource from the state model.
func (s *AssignmentFilterResourceModel) GetID() string {
	return s.ID.ValueString()
}

// GetTypeName returns the type name of the resource from the state model.
func (r *AssignmentFilterResource) GetTypeName() string {
	return r.TypeName
}

// Metadata returns the resource type name.
func (r *AssignmentFilterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph_beta_device_and_app_management_assignment_filter"
}

// Configure sets the client for the resource.
func (r *AssignmentFilterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring AssignmentFilterResource")

	if req.ProviderData == nil {
		tflog.Warn(ctx, "Provider data is nil, skipping resource configuration")
		return
	}

	clients, ok := req.ProviderData.(*client.GraphClients)
	if !ok {
		tflog.Error(ctx, "Unexpected Provider Data Type", map[string]interface{}{
			"expected": "*client.GraphClients",
			"actual":   fmt.Sprintf("%T", req.ProviderData),
		})
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *client.GraphClients, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	if clients.BetaClient == nil {
		tflog.Warn(ctx, "BetaClient is nil, resource may not be fully configured")
		return
	}

	r.client = clients.BetaClient
	tflog.Debug(ctx, "Initialized graphBetaAssignmentFilter resource with Graph Beta Client")
}

// ImportState imports the resource state.
func (r *AssignmentFilterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *AssignmentFilterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the assignment filter.",
			},
			"display_name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the assignment filter.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The optional description of the assignment filter.",
			},
			"platform": schema.StringAttribute{
				Required: true,
				Description: fmt.Sprintf(
					"The Intune device management type (platform) for the assignment filter. "+
						"Must be one of the following values: %s. "+
						"This specifies the OS platform type for which the assignment filter will be applied.",
					strings.Join(validPlatformTypes, ", ")),
				Validators: []validator.String{
					stringvalidator.OneOf(validPlatformTypes...),
				},
			},
			"rule": schema.StringAttribute{
				Required:    true,
				Description: "Rule definition of the assignment filter.",
			},
			"assignment_filter_management_type": schema.StringAttribute{
				Optional:    true,
				Description: fmt.Sprintf("Indicates filter is applied to either 'devices' or 'apps' management type. Possible values are: %s. Default filter will be applied to 'devices'.", strings.Join(validAssignmentFilterManagementTypes, ", ")),
				Validators: []validator.String{
					stringvalidator.OneOf(validAssignmentFilterManagementTypes...),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the assignment filter.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:    true,
				Description: "Last modified time of the assignment filter.",
			},
			"role_scope_tags": schema.ListAttribute{
				Optional:    true,
				Description: "Indicates role scope tags assigned for the assignment filter.",
				ElementType: types.StringType,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
	}
}
