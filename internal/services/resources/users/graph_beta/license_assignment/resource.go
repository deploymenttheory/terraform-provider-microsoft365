package graphBetaUserLicenseAssignment

import (
	"context"
	"fmt"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_users_user_license_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &UserLicenseAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &UserLicenseAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &UserLicenseAssignmentResource{}
)

func NewUserLicenseAssignmentResource() resource.Resource {
	return &UserLicenseAssignmentResource{
		ReadPermissions: []string{
			"User.Read.All",
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"User.ReadWrite.All",
			"Directory.ReadWrite.All",
		},
		ResourcePath: "/users",
	}
}

type UserLicenseAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *UserLicenseAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *UserLicenseAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
// Expected format: {user_id}_{sku_id}
// Example: 00000000-0000-0000-0000-000000000001_11111111-1111-1111-1111-111111111111
func (r *UserLicenseAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Split the ID into user_id and sku_id by finding the last underscore that separates two UUIDs
	// UUIDs are 36 characters long (including hyphens)
	// So we need to find the underscore at position len - 37
	id := req.ID
	if len(id) < 73 { // 36 (uuid) + 1 (_) + 36 (uuid) = 73
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: {user_id}_{sku_id}, got: %s", req.ID),
		)
		return
	}

	// Find the underscore that separates the two UUIDs
	separatorIndex := len(id) - 37 // 36 chars for second UUID + 1 for underscore
	if separatorIndex <= 0 || id[separatorIndex] != '_' {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: {user_id}_{sku_id}, got: %s", req.ID),
		)
		return
	}

	userId := id[0:separatorIndex]
	skuId := id[separatorIndex+1:]

	// Validate UUIDs using the constant from the regex package
	guidRegex := regexp.MustCompile(constants.GuidRegex)
	if !guidRegex.MatchString(userId) || !guidRegex.MatchString(skuId) {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID format: {user_id}_{sku_id}, both must be valid UUIDs. Got: %s", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), userId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("sku_id"), skuId)...)
}

// Schema returns the schema for the resource.
func (r *UserLicenseAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a single Microsoft 365 license assignment for an individual user using the `/users/{userId}/assignLicense` endpoint. " +
			"This resource allows management to Add or remove licenses for the user to enable or disable their use of Microsoft cloud offerings that the company " +
			"has licenses to. For example, an organization can have a Microsoft 365 Enterprise E3 subscription with 100 licenses, and this request assigns one of " +
			"those licenses to a specific user. You can also enable and disable specific plans associated with a subscription. ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The unique identifier for this license assignment resource. Format: `{user_id}_{sku_id}`.",
			},
			"user_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier for the user. Can be either the object ID (UUID) or user principal name (UPN).",
				Validators: []validator.String{
					stringvalidator.Any(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"Must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
						),
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.UserPrincipalNameRegex),
							"Must be a valid User Principal Name format (user@domain.com)",
						),
					),
				},
			},
			"user_principal_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The user principal name (UPN) of the user. This is computed and read-only.",
			},
			"sku_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique identifier (GUID) for the license SKU to assign to the user.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.RequiresReplaceString(),
				},
			},
		"disabled_plans": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "A collection of the unique identifiers for service plans to disable for this license.",
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"Each disabled plan must be a valid UUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
					),
				),
			},
		},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
