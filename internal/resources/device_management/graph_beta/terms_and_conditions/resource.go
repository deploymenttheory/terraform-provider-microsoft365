package graphBetaTermsAndConditions

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_terms_and_conditions"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &TermsAndConditionsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &TermsAndConditionsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &TermsAndConditionsResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &TermsAndConditionsResource{}
)

func NewTermsAndConditionsResource() resource.Resource {
	return &TermsAndConditionsResource{
		ReadPermissions: []string{
			"DeviceManagementServiceConfig.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementServiceConfig.ReadWrite.All",
		},
		ResourcePath: "deviceManagement/termsAndConditions",
	}
}

type TermsAndConditionsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *TermsAndConditionsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + ResourceName
}

// Configure sets the client for the resource.
func (r *TermsAndConditionsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.SetGraphBetaClientForResource(ctx, req, resp, r.TypeName)
}

// ImportState imports the resource state.
func (r *TermsAndConditionsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *TermsAndConditionsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Terms and Conditions in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the terms and conditions policy.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Administrator-supplied name for the terms and conditions policy.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "Administrator-supplied description of the terms and conditions policy. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"title": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Administrator-supplied title of the terms and conditions. This is shown to the user on prompts to accept the terms and conditions policy.",
			},
			"body_text": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Administrator-supplied body text of the terms and conditions, typically the terms themselves. This is shown to the user on prompts to accept the terms and conditions policy.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(60000),
				},
			},
			"acceptance_statement": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Administrator-supplied explanation of the terms and conditions, typically describing what it means to accept the terms and conditions set out in the policy. This is shown to the user on prompts to accept the terms and conditions policy.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
			},
			"version": schema.Int32Attribute{
				Optional:            true,
				Computed:            true,
				Default:             int32default.StaticInt32(1),
				MarkdownDescription: "Integer indicating the current version of the terms. Incremented when an administrator makes a change to the terms and wishes to require users to re-accept the modified terms and conditions policy.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this Entity instance.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "DateTime the object was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_date_time": schema.StringAttribute{
				MarkdownDescription: "DateTime the object was last modified.",
				Computed:            true,
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
