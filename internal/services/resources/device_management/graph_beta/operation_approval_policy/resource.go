package graphBetaOperationApprovalPolicy

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_device_management_operation_approval_policy"
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
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *OperationApprovalPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full type name of the resource for logging purposes.
func (r *OperationApprovalPolicyResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *OperationApprovalPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *OperationApprovalPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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
				MarkdownDescription: "Indicates the description of the policy. Maximum length of the description is 1024 characters. This property is not required, but can be used by the IT Admin to describe the policy.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
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
					"- **unknown**: Unknown policy type\n" +
					"- **app**: Application policy\n" +
					"- **script**: Script policy\n" +
					"- **role**: Role policy",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"unknown",
						"app",
						"script",
						"role",
					),
				},
			},
			// SDK supports only the above fields. MSFT docs suggest the below. weird
			// "policy_type": schema.StringAttribute{
			// 	MarkdownDescription: "The policy type for the OperationApprovalPolicy. Possible values are:\n\n" +
			// 		"- **unknown**: Unknown policy type\n" +
			// 		"- **deviceAction**: Device action policy\n" +
			// 		"- **deviceWipe**: Device wipe policy\n" +
			// 		"- **deviceRetire**: Device retire policy\n" +
			// 		"- **deviceRetireNonCompliant**: Device retire non-compliant policy\n" +
			// 		"- **deviceDelete**: Device delete policy\n" +
			// 		"- **deviceLock**: Device lock policy\n" +
			// 		"- **deviceErase**: Device erase policy\n" +
			// 		"- **deviceDisableActivationLock**: Device disable activation lock policy\n" +
			// 		"- **windowsEnrollment**: Windows enrollment policy\n" +
			// 		"- **compliancePolicy**: Compliance policy\n" +
			// 		"- **configurationPolicy**: Configuration policy\n" +
			// 		"- **appProtectionPolicy**: App protection policy\n" +
			// 		"- **policySet**: Policy set\n" +
			// 		"- **filter**: Filter policy\n" +
			// 		"- **endpointSecurityPolicy**: Endpoint security policy\n" +
			// 		"- **app**: App policy\n" +
			// 		"- **script**: Script policy\n" +
			// 		"- **role**: Role policy\n" +
			// 		"- **deviceResetPasscode**: Device reset passcode policy\n" +
			// 		"- **operationApprovalPolicy**: Operation approval policy",
			// 	Required: true,
			// 	Validators: []validator.String{
			// 		stringvalidator.OneOf(
			// 			"unknown",
			// 			"deviceAction",
			// 			"deviceWipe",
			// 			"deviceRetire",
			// 			"deviceRetireNonCompliant",
			// 			"deviceDelete",
			// 			"deviceLock",
			// 			"deviceErase",
			// 			"deviceDisableActivationLock",
			// 			"windowsEnrollment",
			// 			"compliancePolicy",
			// 			"configurationPolicy",
			// 			"appProtectionPolicy",
			// 			"policySet",
			// 			"filter",
			// 			"endpointSecurityPolicy",
			// 			"app",
			// 			"script",
			// 			"role",
			// 			"deviceResetPasscode",
			// 			"operationApprovalPolicy",
			// 		),
			// 	},
			// },
			"policy_platform": schema.StringAttribute{
				MarkdownDescription: "Indicates the applicable platform for the policy. Possible values are:\n\n" +
					"- **notApplicable**: Not applicable to any platform\n" +
					"- **androidDeviceAdministrator**: Android device administrator platform\n" +
					"- **androidEnterprise**: Android enterprise platform\n" +
					"- **iOSiPadOS**: iOS/iPadOS platform\n" +
					"- **macOS**: macOS platform\n" +
					"- **windows10AndLater**: Windows 10 and later platform\n" +
					"- **windows81AndLater**: Windows 8.1 and later platform\n" +
					"- **windows10X**: Windows 10X platform",
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
					),
				},
			},
			"policy_set": schema.SingleNestedAttribute{
				MarkdownDescription: "Indicates areas of the Intune UX that could support MAA UX for the current logged in IT Admin. This property is required, and is defined by the IT Admins in order to correctly show the expected experience.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"policy_type": schema.StringAttribute{
						MarkdownDescription: "The policy type for the OperationApprovalPolicy. Possible values are:\n\n" +
							"- **unknown**: Unknown policy type\n" +
							"- **app**: Application policy\n" +
							"- **script**: Script policy\n" +
							"- **role**: Role policy",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"unknown",
								"app",
								"script",
								"role",
							),
						},
					},
					// SDK supports only the above fields. MSFT docs suggest the below. weird
					// "policy_type": schema.StringAttribute{
					// 	MarkdownDescription: "The policy type for the policy set.",
					// 	Required:            true,
					// 	Validators: []validator.String{
					// 		stringvalidator.OneOf(
					// 			"unknown",
					// 			"deviceAction",
					// 			"deviceWipe",
					// 			"deviceRetire",
					// 			"deviceRetireNonCompliant",
					// 			"deviceDelete",
					// 			"deviceLock",
					// 			"deviceErase",
					// 			"deviceDisableActivationLock",
					// 			"windowsEnrollment",
					// 			"compliancePolicy",
					// 			"configurationPolicy",
					// 			"appProtectionPolicy",
					// 			"policySet",
					// 			"filter",
					// 			"endpointSecurityPolicy",
					// 			"app",
					// 			"script",
					// 			"role",
					// 			"deviceResetPasscode",
					// 			"operationApprovalPolicy",
					// 		),
					// 	},
					// },
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
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
