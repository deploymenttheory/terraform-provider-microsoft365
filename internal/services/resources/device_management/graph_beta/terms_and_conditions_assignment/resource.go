package graphBetaTermsAndConditionsAssignment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_terms_and_conditions_assignment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &TermsAndConditionsAssignmentResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &TermsAndConditionsAssignmentResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &TermsAndConditionsAssignmentResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &TermsAndConditionsAssignmentResource{}
)

func NewTermsAndConditionsAssignmentResource() resource.Resource {
	return &TermsAndConditionsAssignmentResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
		ResourcePath: "deviceManagement/termsAndConditions",
	}
}

type TermsAndConditionsAssignmentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *TermsAndConditionsAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *TermsAndConditionsAssignmentResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *TermsAndConditionsAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *TermsAndConditionsAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *TermsAndConditionsAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages terms and conditions policy assignments in Microsoft Intune using the `/deviceManagement/termsAndConditions/{termsAndConditionsId}/assignments` endpoint. This resource controls which users or groups must accept specific terms and conditions policies before device enrollment or resource access, enabling targeted compliance enforcement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the terms and conditions assignment.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"terms_and_conditions_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the terms and conditions to attach the assignment to.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"target": schema.SingleNestedAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"target_type": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "The target group type for the assignment. Possible values are:\n\n" +
							"- **allLicensedUsers**: Target all users with Intune licenses\n" +
							"- **groupAssignment**: Target a specific Entra ID group\n" +
							"- **configurationManagerCollection**: Target System Center Configuration Manager collection",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"allLicensedUsers",
								"groupAssignment",
								"configurationManagerCollection",
							),
						},
					},
					"group_id": schema.StringAttribute{
						MarkdownDescription: "The Entra ID group ID for the assignment target. Required when target_type is 'groupAssignment' or 'exclusionGroupAssignment'.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(constants.GuidRegex),
								"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
							),
						},
					},
					"collection_id": schema.StringAttribute{
						MarkdownDescription: "The SCCM group collection ID for the assignment target. Default collections start with 'SMS', while custom collections start with your site code (e.g., 'MEM').",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(
								regexp.MustCompile(`^[A-Za-z]{2,8}[0-9A-Za-z]{8}$`),
								"Must be a valid SCCM collection ID format. Default collections start with 'SMS' followed by an alphanumeric ID. Custom collections start with your site code (e.g., 'MEM') followed by an alphanumeric ID.",
							),
						},
					},
				},
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
