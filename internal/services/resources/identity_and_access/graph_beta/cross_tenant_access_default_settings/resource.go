package graphBetaCrossTenantAccessDefaultSettings

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	resourcevalidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/resource"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// singletonID is the static identifier used for this singleton resource in Terraform state.
	// The crossTenantAccessPolicyConfigurationDefault API has no server-assigned ID; it is addressed by a fixed path.
	singletonID = "crossTenantAccessDefaultSettings"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &CrossTenantAccessDefaultSettingsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CrossTenantAccessDefaultSettingsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CrossTenantAccessDefaultSettingsResource{}

	// Enables config-level validation
	_ resource.ResourceWithConfigValidators = &CrossTenantAccessDefaultSettingsResource{}
)

func NewCrossTenantAccessDefaultSettingsResource() resource.Resource {
	return &CrossTenantAccessDefaultSettingsResource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.CrossTenantAccess",
		},
	}
}

type CrossTenantAccessDefaultSettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *CrossTenantAccessDefaultSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *CrossTenantAccessDefaultSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state. Because this is a singleton resource the ID portion of
// the import address is always normalised to the static identifier "crossTenantAccessDefaultSettings".
//
// restore_defaults_on_destroy is a Terraform-only flag that is never returned by the API. It must
// therefore be supplied explicitly at import time via the extended ID format:
//
//	terraform import <address> crossTenantAccessDefaultSettings:restore_defaults_on_destroy=true
//	terraform import <address> crossTenantAccessDefaultSettings:restore_defaults_on_destroy=false
//
// If the flag is omitted the value defaults to false (state-only removal on destroy).
func (r *CrossTenantAccessDefaultSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	restoreDefaultsOnDestroy := false // safe default

	if len(idParts) > 1 {
		for _, part := range idParts[1:] {
			if strings.HasPrefix(part, "restore_defaults_on_destroy=") {
				value := strings.TrimPrefix(part, "restore_defaults_on_destroy=")
				switch strings.ToLower(value) {
				case "true":
					restoreDefaultsOnDestroy = true
				case "false":
					restoreDefaultsOnDestroy = false
				default:
					resp.Diagnostics.AddError(
						"Invalid Import ID",
						fmt.Sprintf("Invalid restore_defaults_on_destroy value %q. Must be 'true' or 'false'.", value),
					)
					return
				}
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Importing %s with restore_defaults_on_destroy: %t", ResourceName, restoreDefaultsOnDestroy))

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), singletonID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("restore_defaults_on_destroy"), restoreDefaultsOnDestroy)...)
}

// ConfigValidators returns resource-level validators applied before plan/apply.
func (r *CrossTenantAccessDefaultSettingsResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	b2bCollaborationInbound := path.Root("b2b_collaboration_inbound")
	b2bCollaborationOutbound := path.Root("b2b_collaboration_outbound")
	b2bDirectConnectOutbound := path.Root("b2b_direct_connect_outbound")
	tenantRestrictions := path.Root("tenant_restrictions")
	return []resource.ConfigValidator{
		// access_type must be consistent between users_and_groups and applications within each block.
		resourcevalidator.ConsistentNestedStringAttributes(
			b2bCollaborationInbound,
			b2bCollaborationInbound.AtName("users_and_groups").AtName("access_type"),
			b2bCollaborationInbound.AtName("applications").AtName("access_type"),
		),
		resourcevalidator.ConsistentNestedStringAttributes(
			b2bCollaborationOutbound,
			b2bCollaborationOutbound.AtName("users_and_groups").AtName("access_type"),
			b2bCollaborationOutbound.AtName("applications").AtName("access_type"),
		),
		resourcevalidator.ConsistentNestedStringAttributes(
			tenantRestrictions,
			tenantRestrictions.AtName("users_and_groups").AtName("access_type"),
			tenantRestrictions.AtName("applications").AtName("access_type"),
		),
		// "AllUsers" cannot be mixed with specific user or group GUIDs in the same targets set.
		// b2b_direct_connect_inbound is excluded: it only supports "AllUsers" for users_and_groups.
		resourcevalidator.NoAllUsersMix(
			b2bCollaborationInbound,
			b2bCollaborationInbound.AtName("users_and_groups").AtName("targets"),
		),
		resourcevalidator.NoAllUsersMix(
			b2bCollaborationOutbound,
			b2bCollaborationOutbound.AtName("users_and_groups").AtName("targets"),
		),
		resourcevalidator.NoAllUsersMix(
			b2bDirectConnectOutbound,
			b2bDirectConnectOutbound.AtName("users_and_groups").AtName("targets"),
		),
		resourcevalidator.NoAllUsersMix(
			tenantRestrictions,
			tenantRestrictions.AtName("users_and_groups").AtName("targets"),
		),
	}
}

// Schema defines the schema for the resource.
func (r *CrossTenantAccessDefaultSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the default configuration for cross-tenant access policy in Microsoft Entra ID using the `/policies/crossTenantAccessPolicy/default` endpoint.\n\n" +
			"This is a **singleton resource** — one default configuration exists per tenant and cannot be created or deleted via the Microsoft Graph API. " +
			"The `create` operation uses an UPDATE (PATCH) request to configure the default settings. " +
			"On `destroy`, the resource can optionally restore the default configuration to system defaults " +
			"by setting `restore_defaults_on_destroy = true` (using the resetToSystemDefault API), or simply remove it from Terraform state while leaving the configuration in place (the default behaviour).\n\n" +
			"See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyconfigurationdefault?view=graph-rest-beta) for details.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the cross-tenant access default settings. This is a singleton resource; the value is always `crossTenantAccessDefaultSettings`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_service_default": schema.BoolAttribute{
				MarkdownDescription: "If `true`, the default configuration is set to the system default configuration. If `false`, the default settings are customized. This is a read-only computed value.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"b2b_collaboration_inbound": schema.SingleNestedAttribute{
				MarkdownDescription: "B2B collaboration inbound access settings lets you collaborate with people outside of your organization by allowing them to sign in using their own identities. " +
					"These users become guests in your Microsoft Entra tenant. You can invite external users directly or you can set up self-service sign-up so they can request access to your resources.",
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"users_and_groups": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether to allow or block access for users and groups.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"access_type": schema.StringAttribute{
								MarkdownDescription: "The access type. Possible values: `allowed`, `blocked`.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("allowed", "blocked"),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"targets": schema.SetNestedAttribute{
								MarkdownDescription: "The set of user and group targets to allow or block.",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"target": schema.StringAttribute{
											MarkdownDescription: "The unique identifier of the user or group. Can be a user/group GUID, or special value `AllUsers`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.Any(
													stringvalidator.OneOf("AllUsers"),
													stringvalidator.RegexMatches(
														regexp.MustCompile(constants.GuidRegex),
														"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
													),
												),
											},
										},
										"target_type": schema.StringAttribute{
											MarkdownDescription: "The type of target. Possible values: `user`, `group`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("user", "group"),
											},
										},
									},
								},
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
					"applications": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether to allow or block access for applications.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"access_type": schema.StringAttribute{
								MarkdownDescription: "The access type. Possible values: `allowed`, `blocked`.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("allowed", "blocked"),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"targets": schema.SetNestedAttribute{
								MarkdownDescription: "The set of application targets to allow or block.",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"target": schema.StringAttribute{
											MarkdownDescription: "The unique identifier of the application. Can be an application GUID, or special values: `AllApplications`, `Office365`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.Any(
													stringvalidator.OneOf("AllApplications", "Office365"),
													stringvalidator.RegexMatches(
														regexp.MustCompile(constants.GuidRegex),
														"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
													),
												),
											},
										},
										"target_type": schema.StringAttribute{
											MarkdownDescription: "The type of target. Must be `application`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("application"),
											},
										},
									},
								},
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"b2b_collaboration_outbound": schema.SingleNestedAttribute{
				MarkdownDescription: "B2B collaboration outbound access settings determine whether your users can be invited to external Microsoft Entra tenants for B2B collaboration and added to " +
					"their directories as guests. These default settings apply to all external Microsoft Entra tenants except those with organization-specific settings. Below, specify whether your " +
					"users and groups can be invited for B2B collaboration and the external applications they can access.",
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"users_and_groups": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether to allow or block access for users and groups.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"access_type": schema.StringAttribute{
								MarkdownDescription: "The access type. Possible values: `allowed`, `blocked`.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("allowed", "blocked"),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"targets": schema.SetNestedAttribute{
								MarkdownDescription: "The set of user and group targets to allow or block.",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"target": schema.StringAttribute{
											MarkdownDescription: "The unique identifier of the user or group. Can be a user/group GUID, or special value `AllUsers`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.Any(
													stringvalidator.OneOf("AllUsers"),
													stringvalidator.RegexMatches(
														regexp.MustCompile(constants.GuidRegex),
														"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
													),
												),
											},
										},
										"target_type": schema.StringAttribute{
											MarkdownDescription: "The type of target. Possible values: `user`, `group`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("user", "group"),
											},
										},
									},
								},
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
					"applications": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether to allow or block access for applications.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"access_type": schema.StringAttribute{
								MarkdownDescription: "The access type. Possible values: `allowed`, `blocked`.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("allowed", "blocked"),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"targets": schema.SetNestedAttribute{
								MarkdownDescription: "The set of application targets to allow or block.",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"target": schema.StringAttribute{
											MarkdownDescription: "The unique identifier of the application. Can be an application GUID, or special values: `AllApplications`, `Office365`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.Any(
													stringvalidator.OneOf("AllApplications", "Office365"),
													stringvalidator.RegexMatches(
														regexp.MustCompile(constants.GuidRegex),
														"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
													),
												),
											},
										},
										"target_type": schema.StringAttribute{
											MarkdownDescription: "The type of target. Must be `application`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("application"),
											},
										},
									},
								},
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"b2b_direct_connect_inbound": schema.SingleNestedAttribute{
				MarkdownDescription: "B2B direct connect inbound access settings determine whether users from external Microsoft Entra tenants can access your " +
					"resources without being added to your tenant as guests. By selecting 'Allow access' below, you're permitting users and groups from other " +
					"organizations to connect with you. To establish a connection, an admin from the other organization must also enable B2B direct connect.",
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"users_and_groups": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether to allow or block access for users and groups.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"access_type": schema.StringAttribute{
								MarkdownDescription: "The access type. Possible values: `allowed`, `blocked`.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("allowed", "blocked"),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"targets": schema.SetNestedAttribute{
								MarkdownDescription: "The set of user and group targets to allow or block.",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"target": schema.StringAttribute{
											MarkdownDescription: "The unique identifier of the user or group. Can be a user/group GUID, or special value `AllUsers`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.Any(
													stringvalidator.OneOf("AllUsers"),
													stringvalidator.RegexMatches(
														regexp.MustCompile(constants.GuidRegex),
														"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
													),
												),
											},
										},
										"target_type": schema.StringAttribute{
											MarkdownDescription: "The type of target. Possible values: `user`, `group`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("user", "group"),
											},
										},
									},
								},
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
					"applications": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether to allow or block access for applications.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"access_type": schema.StringAttribute{
								MarkdownDescription: "The access type. Possible values: `allowed`, `blocked`.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("allowed", "blocked"),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"targets": schema.SetNestedAttribute{
								MarkdownDescription: "The set of application targets to allow or block.",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"target": schema.StringAttribute{
											MarkdownDescription: "The unique identifier of the application. Can be an application GUID, or special values: `AllApplications`, `Office365`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.Any(
													stringvalidator.OneOf("AllApplications", "Office365"),
													stringvalidator.RegexMatches(
														regexp.MustCompile(constants.GuidRegex),
														"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
													),
												),
											},
										},
										"target_type": schema.StringAttribute{
											MarkdownDescription: "The type of target. Must be `application`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("application"),
											},
										},
									},
								},
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"b2b_direct_connect_outbound": schema.SingleNestedAttribute{
				MarkdownDescription: "Outbound access settings determine how your users and groups can interact with apps and resources in external organizations. " +
					"The default settings apply to all your cross-tenant scenarios unless you configure organizational settings to override them for a specific " +
					"organization. Default settings can be modified but not deleted.",
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"users_and_groups": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether to allow or block access for users and groups.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"access_type": schema.StringAttribute{
								MarkdownDescription: "The access type. Possible values: `allowed`, `blocked`.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("allowed", "blocked"),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"targets": schema.SetNestedAttribute{
								MarkdownDescription: "The set of user and group targets to allow or block.",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"target": schema.StringAttribute{
											MarkdownDescription: "The unique identifier of the user or group. Can be a user/group GUID, or special value `AllUsers`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.Any(
													stringvalidator.OneOf("AllUsers"),
													stringvalidator.RegexMatches(
														regexp.MustCompile(constants.GuidRegex),
														"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
													),
												),
											},
										},
										"target_type": schema.StringAttribute{
											MarkdownDescription: "The type of target. Possible values: `user`, `group`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("user", "group"),
											},
										},
									},
								},
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
					"applications": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether to allow or block access for applications.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"access_type": schema.StringAttribute{
								MarkdownDescription: "The access type. Possible values: `allowed`, `blocked`.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("allowed", "blocked"),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"targets": schema.SetNestedAttribute{
								MarkdownDescription: "The set of application targets to allow or block.",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"target": schema.StringAttribute{
											MarkdownDescription: "The unique identifier of the application. Can be an application GUID, or special values: `AllApplications`, `Office365`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.Any(
													stringvalidator.OneOf("AllApplications", "Office365"),
													stringvalidator.RegexMatches(
														regexp.MustCompile(constants.GuidRegex),
														"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
													),
												),
											},
										},
										"target_type": schema.StringAttribute{
											MarkdownDescription: "The type of target. Must be `application`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("application"),
											},
										},
									},
								},
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"inbound_trust": schema.SingleNestedAttribute{
				MarkdownDescription: "Configure whether your Conditional Access policies will accept claims from other Microsoft Entra tenants when external users access your resources. " +
					"The default settings apply to all external Microsoft Entra tenants except those with organization-specific settings. You'll first need to configure Conditional Access for " +
					"guest users on all cloud apps if you want to require multifactor authentication or require a device to be compliant or Microsoft Entra hybrid joined.",
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"is_mfa_accepted": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether to trust MFA claims from external Microsoft Entra organizations.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"is_compliant_device_accepted": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether to trust compliant device claims from external Microsoft Entra organizations.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"is_hybrid_azure_ad_joined_device_accepted": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether to trust hybrid Azure AD joined device claims from external Microsoft Entra organizations.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"invitation_redemption_identity_provider_configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "Defines the priority order based on which an identity provider is selected during invitation redemption for a guest user.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"primary_identity_provider_precedence_order": schema.ListAttribute{
						MarkdownDescription: "Users will redeem their invitations using the default order set by Microsoft. You can enable and specify the order of " +
							"identity providers that your guest users can sign in with when they redeem their invitation. Possible values are: `externalFederation`, " +
							"`azureActiveDirectory`, `socialIdentityProviders`. By not specifying an invitation redemption identity provider type it will set set to disabled.",
						ElementType: types.StringType,
						Required:    true,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(
								stringvalidator.OneOf(
									"externalFederation",
									"azureActiveDirectory",
									"socialIdentityProviders",
								),
							),
						},
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"fallback_identity_provider": schema.StringAttribute{
						MarkdownDescription: "Fallback identity providers are used when none of the primary identity providers are applicable. " +
							"You must always have at least one fallback provider set to prevent users from being blocked while redeeming an invitation. Possible values: `defaultConfiguredIdp`, `emailOneTimePasscode`.",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"defaultConfiguredIdp",
								"emailOneTimePasscode",
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_restrictions": schema.SingleNestedAttribute{
				MarkdownDescription: "Tenant restrictions lets you control whether your users can access external applications from your network or " +
					"devices using external accounts, including accounts issued to them by external organizations and accounts they've created in unknown " +
					"tenants. Below, select which external applications to allow or block. These default settings apply to all external Microsoft Entra " +
					"tenants except those with organization-specific settings.",
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"users_and_groups": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether to allow or block access for users and groups.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"access_type": schema.StringAttribute{
								MarkdownDescription: "The access type. Possible values: `allowed`, `blocked`.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("allowed", "blocked"),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"targets": schema.SetNestedAttribute{
								MarkdownDescription: "The set of user and group targets to allow or block.",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"target": schema.StringAttribute{
											MarkdownDescription: "The unique identifier of the user or group. Can be a user/group GUID, or special value `AllUsers`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.Any(
													stringvalidator.OneOf("AllUsers"),
													stringvalidator.RegexMatches(
														regexp.MustCompile(constants.GuidRegex),
														"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
													),
												),
											},
										},
										"target_type": schema.StringAttribute{
											MarkdownDescription: "The type of target. Possible values: `user`, `group`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("user", "group"),
											},
										},
									},
								},
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
					"applications": schema.SingleNestedAttribute{
						MarkdownDescription: "Specifies whether to allow or block access for applications.",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"access_type": schema.StringAttribute{
								MarkdownDescription: "The access type. Possible values: `allowed`, `blocked`.",
								Optional:            true,
								Computed:            true,
								Validators: []validator.String{
									stringvalidator.OneOf("allowed", "blocked"),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"targets": schema.SetNestedAttribute{
								MarkdownDescription: "The set of application targets to allow or block.",
								Required:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"target": schema.StringAttribute{
											MarkdownDescription: "The unique identifier of the application. Can be an application GUID, or special values: `AllApplications`, `Office365`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.Any(
													stringvalidator.OneOf("AllApplications", "Office365"),
													stringvalidator.RegexMatches(
														regexp.MustCompile(constants.GuidRegex),
														"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
													),
												),
											},
										},
										"target_type": schema.StringAttribute{
											MarkdownDescription: "The type of target. Must be `application`.",
											Required:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("application"),
											},
										},
									},
								},
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.UseStateForUnknown(),
								},
							},
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"automatic_user_consent_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "Determines the default configuration for automatic user consent settings. The `inbound_allowed` and `outbound_allowed` properties are always `false` and can't be updated in the default configuration. Read-only.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"inbound_allowed": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether inbound automatic user consent is allowed. This is always `false` in the default configuration.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"outbound_allowed": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether outbound automatic user consent is allowed. This is always `false` in the default configuration.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"restore_defaults_on_destroy": schema.BoolAttribute{
				MarkdownDescription: "Controls behaviour when this resource is destroyed. " +
					"When `true`, Terraform will issue a POST request to `/policies/crossTenantAccessPolicy/default/resetToSystemDefault` to reset the default configuration to system defaults, " +
					"then verify that `is_service_default` is `true` before removing from state. " +
					"When `false` (the default), Terraform removes the resource from state only — the existing default configuration is left unchanged in Microsoft Entra ID.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
