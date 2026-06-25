package graphBetaOperationApprovalPolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_operation_approval_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &OperationApprovalPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &OperationApprovalPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &OperationApprovalPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &OperationApprovalPolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &OperationApprovalPolicyResource{}
)

func NewOperationApprovalPolicyResource() resource.Resource {
	return &OperationApprovalPolicyResource{
		ReadPermissions: []string{
			"DeviceManagementRBAC.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementRBAC.ReadWrite.All",
		},
		ResourcePath: "deviceManagement/operationApprovalPolicies",
	}
}

type OperationApprovalPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *OperationApprovalPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *OperationApprovalPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *OperationApprovalPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *OperationApprovalPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema returns the schema for the resource.
func (r *OperationApprovalPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Operation Approval Policies in Microsoft Intune.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the policy. This ID is assigned at when the policy is created. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Indicates the display name of the policy. Maximum length of the display name is 128 characters. This property is required when the policy is created, and is defined by the IT Admins to identify the policy.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "Indicates the last DateTime that the policy was modified. The value cannot be modified and is automatically populated whenever values in the request are updated. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_type": schema.StringAttribute{
				MarkdownDescription: "The policy type for the OperationApprovalPolicy. Possible values are:\n\n" +
					"- **unknown**: Default. Unknown policy type.\n" +
					"- **deviceWipe**: Device wipe policy type.\n" +
					"- **deviceRetire**: Device retire policy type.\n" +
					"- **deviceDelete**: Device delete policy type.\n" +
					"- **app**: Application policy type.\n" +
					"- **script**: Script policy type.\n" +
					"- **role**: Role policy type.\n" +
					"- **unknownFutureValue**: Evolvable enumeration sentinel value. Do not use.\n" +
					"- **tenantConfiguration**: Tenant configuration policy type.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"unknown",
						"deviceWipe",
						"deviceRetire",
						"deviceDelete",
						"app",
						"script",
						"role",
						"unknownFutureValue",
						"tenantConfiguration",
					),
				},
			},
			"policy_platform": schema.StringAttribute{
				MarkdownDescription: "Indicates the applicable platform for the policy. Possible values are:\n\n" +
					"- **notApplicable**: Not applicable to any platform\n" +
					"- **androidDeviceAdministrator**: Android device administrator platform\n" +
					"- **androidEnterprise**: Android enterprise platform\n" +
					"- **iOSiPadOS**: iOS/iPadOS platform\n" +
					"- **macOS**: macOS platform\n" +
					"- **windows10AndLater**: Windows 10 and later platform\n" +
					"- **windows81AndLater**: Windows 8.1 and later platform\n" +
					"- **windows10X**: Windows 10X platform\n" +
					"- **unknownFutureValue**: Evolvable enumeration sentinel value. Do not use.",
				Optional: true,
				Computed: true,
				//Default:  stringdefault.StaticString("notApplicable"),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"notApplicable",
						"androidDeviceAdministrator",
						"androidEnterprise",
						"iOSiPadOS",
						"macOS",
						"windows10AndLater",
						"windows81AndLater",
						"windows10X",
						"unknownFutureValue",
					),
				},
			},
			"policy_set": schema.SingleNestedAttribute{
				MarkdownDescription: "Indicates areas of the Intune UX that could support MAA UX for the current logged in IT Admin. This property is required, and is defined by the IT Admins in order to correctly show the expected experience.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"policy_type": schema.StringAttribute{
						MarkdownDescription: "The policy type for the policy set. Possible values are:\n\n" +
							"- **unknown**: Default. Unknown policy type.\n" +
							"- **deviceWipe**: Device wipe policy type.\n" +
							"- **deviceRetire**: Device retire policy type.\n" +
							"- **deviceDelete**: Device delete policy type.\n" +
							"- **app**: Application policy type.\n" +
							"- **script**: Script policy type.\n" +
							"- **role**: Role policy type.\n" +
							"- **unknownFutureValue**: Evolvable enumeration sentinel value. Do not use.\n" +
							"- **tenantConfiguration**: Tenant configuration policy type.",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"unknown",
								"deviceWipe",
								"deviceRetire",
								"deviceDelete",
								"app",
								"script",
								"role",
								"unknownFutureValue",
								"tenantConfiguration",
							),
						},
					},
					"policy_platform": schema.StringAttribute{
						MarkdownDescription: "The policy platform for the policy set.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("notApplicable"),
						Validators: []validator.String{
							stringvalidator.OneOf(
								"notApplicable",
								"androidDeviceAdministrator",
								"androidEnterprise",
								"iOSiPadOS",
								"macOS",
								"windows10AndLater",
								"windows81AndLater",
								"windows10X",
								"unknownFutureValue",
							),
						},
					},
				},
			},
			"approver_group_ids": schema.SetAttribute{
				MarkdownDescription: "The Microsoft Entra ID (Azure AD) security group IDs for the approvers for the policy. This property is required when the policy is created, and is defined by the IT Admins to define the possible approvers for the policy.",
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
							"Must be a valid GUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
						),
					),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
