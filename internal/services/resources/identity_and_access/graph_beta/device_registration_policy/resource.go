package graphBetaIdentityAndAccessDeviceRegistrationPolicy

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_device_registration_policy"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// singletonID is the fixed identifier for the tenant-wide device registration policy.
	singletonID = "deviceRegistrationPolicy"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &DeviceRegistrationPolicyResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &DeviceRegistrationPolicyResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &DeviceRegistrationPolicyResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &DeviceRegistrationPolicyResource{}

	// Enables configuration validation
	_ resource.ResourceWithValidateConfig = &DeviceRegistrationPolicyResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &DeviceRegistrationPolicyResource{}
)

// https://learn.microsoft.com/en-us/graph/api/deviceregistrationpolicy-get?view=graph-rest-beta#permissions
// https://learn.microsoft.com/en-us/graph/api/deviceregistrationpolicy-update?view=graph-rest-beta#permissions
func NewDeviceRegistrationPolicyResource() resource.Resource {
	return &DeviceRegistrationPolicyResource{
		ReadPermissions: []string{
			"Policy.Read.DeviceConfiguration",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.DeviceConfiguration",
		},
	}
}

type DeviceRegistrationPolicyResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *DeviceRegistrationPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *DeviceRegistrationPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using the fixed singleton identifier.
func (r *DeviceRegistrationPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import this singleton using the literal ID deviceRegistrationPolicy.",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *DeviceRegistrationPolicyResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *DeviceRegistrationPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the tenant-wide Microsoft Entra device registration policy using the Microsoft Graph beta `/policies/deviceRegistrationPolicy` endpoint. " +
			"The policy controls the device quota per user, multifactor authentication for device registration, Windows LAPS, and who can register (Microsoft Entra registered) " +
			"or join (Microsoft Entra joined) devices.\n\n" +
			"This is a **singleton resource**: one policy exists per tenant. `create` adopts the existing policy and `create`/`update` apply the configuration with a full-replace PUT. " +
			"Settings that are not configured keep their current tenant values, because the provider reads the policy before every PUT and sends unconfigured values back unchanged. " +
			"On `destroy`, Terraform only removes the resource from state and leaves the tenant policy as it is. Manage only one instance per tenant across all Terraform states.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The fixed singleton identifier `deviceRegistrationPolicy`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the device registration policy. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the device registration policy. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_device_quota": schema.Int32Attribute{
				MarkdownDescription: "The maximum number of devices that a user can have in the organization before new device registrations are blocked. " +
					"When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			"multi_factor_auth_configuration": schema.StringAttribute{
				MarkdownDescription: "Whether multifactor authentication is required to register or join devices to Microsoft Entra ID. " +
					"Possible values are: `notRequired`, `required`. When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("notRequired", "required"),
				},
			},
			"local_admin_password": schema.SingleNestedAttribute{
				MarkdownDescription: "Windows Local Administrator Password Solution (LAPS) settings. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"is_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether Windows LAPS is enabled for Microsoft Entra joined devices.",
						Required:            true,
					},
				},
			},
			"azure_ad_registration": schema.SingleNestedAttribute{
				MarkdownDescription: "Settings for Microsoft Entra registered devices. When not configured, the current tenant value is kept. " +
					"When Microsoft Intune is enabled in the tenant, `is_admin_configurable` is `false` and Graph does not accept changes to these settings; " +
					"configuring a different `allowed_to_register` value then fails with a diagnostic.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"is_admin_configurable": schema.BoolAttribute{
						MarkdownDescription: "Whether an administrator can change these settings. `false` when Microsoft Intune manages device registration. Read-only.",
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"allowed_to_register": membershipAttribute("Who can register devices as Microsoft Entra registered.", true),
				},
			},
			"azure_ad_join": schema.SingleNestedAttribute{
				MarkdownDescription: "Settings for Microsoft Entra joined devices. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"is_admin_configurable": schema.BoolAttribute{
						MarkdownDescription: "Whether an administrator can change these settings. Read-only.",
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"allowed_to_join": membershipAttribute("Who can join devices to Microsoft Entra ID. When not configured, the current tenant value is kept.", false),
					"local_admins": schema.SingleNestedAttribute{
						MarkdownDescription: "Local administrator settings for Microsoft Entra joined devices. When not configured, the current tenant value is kept.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"enable_global_admins": schema.BoolAttribute{
								MarkdownDescription: "Whether global administrators are local administrators on all Microsoft Entra joined devices. When not configured, the current tenant value is kept.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"registering_users": membershipAttribute(registeringUsersDescription, false),
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}

// registeringUsersDescription describes azure_ad_join.local_admins.registering_users.
const registeringUsersDescription = "Users who are local administrators on the Microsoft Entra joined devices they register. When not configured, the current tenant value is kept."

// membershipAttribute returns the schema for a polymorphic deviceRegistrationMembership.
func membershipAttribute(description string, required bool) schema.SingleNestedAttribute {
	attribute := schema.SingleNestedAttribute{
		MarkdownDescription: description + " `type` selects the membership kind: `all` (everyone), `none` (nobody) or `enumerated` (the listed `users` and `groups`).",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				MarkdownDescription: "The membership kind. Possible values are: `all`, `none`, `enumerated`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(membershipTypeAll, membershipTypeNone, membershipTypeEnumerated),
				},
			},
			"users": schema.SetAttribute{
				MarkdownDescription: "User object IDs. Only valid when `type` is `enumerated`.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
						),
					),
				},
			},
			"groups": schema.SetAttribute{
				MarkdownDescription: "Group object IDs. Only valid when `type` is `enumerated`.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
						),
					),
				},
			},
		},
	}
	if required {
		attribute.Required = true
		return attribute
	}
	attribute.Optional = true
	attribute.Computed = true
	attribute.PlanModifiers = []planmodifier.Object{
		objectplanmodifier.UseStateForUnknown(),
	}
	return attribute
}
