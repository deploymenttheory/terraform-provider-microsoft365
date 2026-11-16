package graphBetaIdentityAndAccessMobilityManagementPolicy

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_mobility_management_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// policyTypeDevice selects the /policies/mobileDeviceManagementPolicies collection (MDM user scope).
	policyTypeDevice = "device"

	// policyTypeApp selects the /policies/mobileAppManagementPolicies collection (MAM user scope).
	policyTypeApp = "app"

	// appliesToNone, appliesToAll and appliesToSelected are the supported policyScope values.
	appliesToNone     = "none"
	appliesToAll      = "all"
	appliesToSelected = "selected"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &MobilityManagementPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &MobilityManagementPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &MobilityManagementPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &MobilityManagementPolicyResource{}

	// Enables common validations
	_ resource.ResourceWithConfigValidators = &MobilityManagementPolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &MobilityManagementPolicyResource{}
)

// https://learn.microsoft.com/en-us/graph/api/mobiledevicemanagementpolicies-get?view=graph-rest-beta#permissions
// https://learn.microsoft.com/en-us/graph/api/mobiledevicemanagementpolicies-update?view=graph-rest-beta#permissions
// https://learn.microsoft.com/en-us/graph/api/mobiledevicemanagementpolicies-post-includedgroups?view=graph-rest-beta#permissions
// Microsoft Graph documents no application permission for these endpoints; only delegated (work or school account) permissions are supported.
func NewMobilityManagementPolicyResource() resource.Resource {
	return &MobilityManagementPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.MobilityManagement",
		},
	}
}

type MobilityManagementPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *MobilityManagementPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *MobilityManagementPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using the composite ID {policy_type}/{policy_id}.
func (r *MobilityManagementPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	policyType, policyID, err := parseResourceID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_type"), policyType)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_id"), policyID)...)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *MobilityManagementPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *MobilityManagementPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Microsoft Entra mobility management policy (the MDM or MAM automatic enrollment user scope) using the Microsoft Graph beta " +
			"`/policies/mobileDeviceManagementPolicies` and `/policies/mobileAppManagementPolicies` endpoints. A mobility management policy is created by the service " +
			"when a mobility management application such as Microsoft Intune is added to the tenant, so this resource adopts an existing policy and never creates or deletes one.\n\n" +
			"`create` adopts the policy and applies the configured settings; `update` applies changes. Settings that are not configured keep their current tenant value and are never written. " +
			"`applies_to` is written in its own request because Microsoft Graph rejects it together with other properties. `selected` is applied by adding the `included_groups` references, " +
			"and only the group references that differ are added or removed. On `destroy`, Terraform only removes the resource from state and leaves the tenant policy as it is.\n\n" +
			"**Permissions:** Microsoft Graph does not support application permissions for these endpoints. The provider must authenticate with a delegated (work or school account) " +
			"identity that holds `Policy.ReadWrite.MobilityManagement` and a role such as Cloud Device Administrator.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The resource identifier in the format `{policy_type}/{policy_id}`, for example `device/0000000a-0000-0000-c000-000000000000`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"policy_type": schema.StringAttribute{
				MarkdownDescription: "The kind of mobility management policy. `device` manages the MDM user scope (`/policies/mobileDeviceManagementPolicies`) and `app` manages the MAM user scope (`/policies/mobileAppManagementPolicies`). Changing this forces a new resource.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(policyTypeDevice, policyTypeApp),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the mobility management policy, which is the application ID of the mobility management application. For Microsoft Intune this is `0000000a-0000-0000-c000-000000000000`. Changing this forces a new resource.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the mobility management application. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the mobility management application. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_valid": schema.BoolAttribute{
				MarkdownDescription: "Whether the policy is valid. Invalid policies cannot be updated. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"applies_to": schema.StringAttribute{
				MarkdownDescription: "The user scope of the policy. Possible values are: `none`, `all`, `selected`. " +
					"`selected` requires at least one entry in `included_groups`. `all` and `none` remove every included group. " +
					"**Warning:** `all` makes every user in the tenant in scope for automatic enrollment. When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(appliesToNone, appliesToAll, appliesToSelected),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"included_groups": schema.SetAttribute{
				MarkdownDescription: "The IDs of the Microsoft Entra groups in scope of the policy. Only valid when `applies_to` is `selected`. " +
					"Group references are added and removed individually, so only the groups that differ are changed. When not configured, the current tenant groups are kept.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
						),
					),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"compliance_url": schema.StringAttribute{
				MarkdownDescription: "The compliance URL of the mobility management application. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"discovery_url": schema.StringAttribute{
				MarkdownDescription: "The discovery URL of the mobility management application. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"terms_of_use_url": schema.StringAttribute{
				MarkdownDescription: "The terms of use URL of the mobility management application. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_mdm_enrollment_during_registration_disabled": schema.BoolAttribute{
				MarkdownDescription: "Whether users in an automatic enrollment configuration on Microsoft Entra registered devices are no longer prompted to enroll their device in MDM " +
					"during the Entra account registration flow. Only valid when `policy_type` is `device`. When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
