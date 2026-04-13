package graphBetaCrossTenantAccessPartnerSettings

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	resourcevalidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/resource"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &CrossTenantAccessPartnerSettingsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CrossTenantAccessPartnerSettingsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CrossTenantAccessPartnerSettingsResource{}

	// Enables config-level validation
	_ resource.ResourceWithConfigValidators = &CrossTenantAccessPartnerSettingsResource{}
)

func NewCrossTenantAccessPartnerSettingsResource() resource.Resource {
	return &CrossTenantAccessPartnerSettingsResource{
		ReadPermissions: []string{
			"Directory.Read.All",
			"Group.ManageProtection.All",
			"Group.Read.All",
			"Policy.Read.All",
			"User.Read.All",
			"User.ReadBasic.All",
		},
		WritePermissions: []string{
			"Policy.ReadWrite.CrossTenantAccess",
		},
	}
}

type CrossTenantAccessPartnerSettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *CrossTenantAccessPartnerSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *CrossTenantAccessPartnerSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles importing the resource with an extended ID format.
//
// Supported formats:
//   - Simple:   "tenant_id" (hard_delete defaults to false)
//   - Extended: "tenant_id:hard_delete=true" or "tenant_id:hard_delete=false"
//
// Example:
//
//	terraform import microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings.example "12345678-1234-1234-1234-123456789012:hard_delete=true"
func (r *CrossTenantAccessPartnerSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	resourceID := idParts[0]
	hardDelete := false

	if len(idParts) > 1 {
		for _, part := range idParts[1:] {
			if value, found := strings.CutPrefix(part, "hard_delete="); found {
				switch strings.ToLower(value) {
				case "true":
					hardDelete = true
				case "false":
					hardDelete = false
				default:
					resp.Diagnostics.AddError(
						"Invalid Import ID",
						fmt.Sprintf("Invalid hard_delete value '%s'. Must be 'true' or 'false'.", value),
					)
					return
				}
			}
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Importing %s with ID: %s, hard_delete: %t", ResourceName, resourceID, hardDelete))

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tenant_id"), resourceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("hard_delete"), hardDelete)...)
}

// ConfigValidators returns resource-level validators applied before plan/apply.
func (r *CrossTenantAccessPartnerSettingsResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	b2bCollaborationInbound := path.Root("b2b_collaboration_inbound")
	b2bCollaborationOutbound := path.Root("b2b_collaboration_outbound")
	b2bDirectConnectOutbound := path.Root("b2b_direct_connect_outbound")
	tenantRestrictions := path.Root("tenant_restrictions")
	return []resource.ConfigValidator{
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
func (r *CrossTenantAccessPartnerSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages partner-specific cross-tenant access settings in Microsoft Entra ID using the `/policies/crossTenantAccessPolicy/partners` endpoint. This resource is used to configure B2B collaboration, B2B direct connect, inbound trust, and tenant restrictions for a specific partner organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the partner configuration. This is the same as the `tenant_id`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The tenant ID of the partner Microsoft Entra organization. This is a GUID that uniquely identifies the partner tenant.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_service_provider": schema.BoolAttribute{
				MarkdownDescription: "Identifies whether the partner-specific configuration is a cloud service provider for your organization. " +
					"**Important**: This field can only be set when using delegated (user) authentication. " +
					"When using application (client credentials) authentication, this field must be omitted entirely - " +
					"the API will reject requests with 403's that explicitly set this field to either `true` or `false`. " +
					"This is a read-only computed field when using service principal authentication.",
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_in_multi_tenant_organization": schema.BoolAttribute{
				MarkdownDescription: "Identifies whether the partner organization is part of a multi-tenant organization with the local tenant.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"hard_delete": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				MarkdownDescription: "When `true`, the partner configuration will be permanently deleted (hard delete) during destroy. " +
					"When `false` (default), the partner configuration will only be soft deleted and moved to the deleted items container where it can be restored within 30 days. " +
					"**Note**: Hard delete permanently removes the configuration and cannot be undone.",
			},
			"b2b_collaboration_inbound": schema.SingleNestedAttribute{
				MarkdownDescription: "B2B collaboration inbound access settings for the partner organization.",
				Optional:            true,
				Computed:            true,
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
				MarkdownDescription: "B2B collaboration outbound access settings for the partner organization.",
				Optional:            true,
				Computed:            true,
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
				MarkdownDescription: "B2B direct connect inbound access settings for the partner organization.",
				Optional:            true,
				Computed:            true,
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
				MarkdownDescription: "B2B direct connect outbound access settings for the partner organization.",
				Optional:            true,
				Computed:            true,
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
				MarkdownDescription: "Inbound trust settings for accepting claims from the partner organization.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"is_mfa_accepted": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether to accept MFA claims from the partner organization.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"is_compliant_device_accepted": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether to accept compliant device claims from the partner organization.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"is_hybrid_azure_ad_joined_device_accepted": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether to accept hybrid Azure AD joined device claims from the partner organization.",
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
			"automatic_user_consent_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "Automatic user consent settings for the partner organization.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"inbound_allowed": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether automatic user consent is allowed for inbound flows.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"outbound_allowed": schema.BoolAttribute{
						MarkdownDescription: "Specifies whether automatic user consent is allowed for outbound flows.",
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
			"tenant_restrictions": schema.SingleNestedAttribute{
				MarkdownDescription: "Tenant restrictions settings for the partner organization.",
				Optional:            true,
				Computed:            true,
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
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
