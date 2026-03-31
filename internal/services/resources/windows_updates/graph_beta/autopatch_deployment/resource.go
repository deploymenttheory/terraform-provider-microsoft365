package graphBetaWindowsUpdatesAutopatchDeployment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	resourcevalidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/resource"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_windows_updates_autopatch_deployment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsUpdatesAutopatchDeploymentResource{}
	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsUpdatesAutopatchDeploymentResource{}
	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsUpdatesAutopatchDeploymentResource{}
	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &WindowsUpdatesAutopatchDeploymentResource{}
	// Enables resource-level config validation
	_ resource.ResourceWithConfigValidators = &WindowsUpdatesAutopatchDeploymentResource{}
)

func NewWindowsUpdatesAutopatchDeploymentResource() resource.Resource {
	return &WindowsUpdatesAutopatchDeploymentResource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/deployments",
	}
}

type WindowsUpdatesAutopatchDeploymentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsUpdatesAutopatchDeploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsUpdatesAutopatchDeploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsUpdatesAutopatchDeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsUpdatesAutopatchDeploymentResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *WindowsUpdatesAutopatchDeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Update deployments for deploying content to a set of devices using the `/admin/windows/updates/deployments` endpoint. " +
			"Deployments define which update content (feature or quality updates) should be deployed, to which audience, " +
			"and with what settings (schedule, monitoring, etc.). ",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Unique identifier for the deployment.",
			},
			"content": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Specifies what content to deploy. Cannot be changed after creation.",
				Attributes: map[string]schema.Attribute{
					"catalog_entry_id": schema.StringAttribute{
						Required: true,
						MarkdownDescription: "The ID of the catalog entry to deploy. This should reference a feature or " +
							"quality update from the Windows Update catalog. Cannot be changed after creation and will require a replacement to update.",
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							planmodifiers.RequiresReplaceString(),
						},
					},
					"catalog_entry_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The type of catalog entry being deployed. Valid values are: `featureUpdate`, `qualityUpdate`.",
						Validators: []validator.String{
							stringvalidator.OneOf("featureUpdate", "qualityUpdate"),
						},
						PlanModifiers: []planmodifier.String{
							planmodifiers.RequiresReplaceString(),
						},
					},
				},
			},
			"settings": schema.SingleNestedAttribute{
				Optional: true,
				MarkdownDescription: "Settings specified on the deployment governing how to deploy content. " +
					"Settings may be added to a deployment that was created without them, but once configured " +
					"they cannot be modified in place — changes will require the resource to be replaced.",
				PlanModifiers: []planmodifier.Object{
					planmodifiers.RequiresReplaceIfStateNonNullObject(),
				},
				Attributes: map[string]schema.Attribute{
					"schedule": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Schedule settings for the deployment.",
						Attributes: map[string]schema.Attribute{
							"start_date_time": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The date and time when the deployment should start. Must be in ISO 8601 format.",
								Validators: []validator.String{
									stringvalidator.RegexMatches(
										regexp.MustCompile(constants.ISO8601DateTimeRegex),
										"must be a valid ISO 8601 datetime",
									),
								},
							},
							"gradual_rollout": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "Settings for gradual rollout of the deployment.",
								Attributes: map[string]schema.Attribute{
									"duration_between_offers": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "The duration between each offer in ISO 8601 format (e.g., `P7D` for 7 days).",
									},
									"devices_per_offer": schema.Int32Attribute{
										Optional:            true,
										MarkdownDescription: "The number of devices to offer the update to in each rollout wave.",
										Validators: []validator.Int32{
											int32validator.AtLeast(1),
										},
									},
									"end_date_time": schema.StringAttribute{
										Optional:            true,
										MarkdownDescription: "The date and time when the gradual rollout should complete. Must be in ISO 8601 format.",
										Validators: []validator.String{
											stringvalidator.RegexMatches(
												regexp.MustCompile(constants.ISO8601DateTimeRegex),
												"must be a valid ISO 8601 datetime",
											),
										},
									},
								},
							},
						},
					},
					"monitoring": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Monitoring settings for the deployment.",
						Attributes: map[string]schema.Attribute{
							"monitoring_rules": schema.SetNestedAttribute{
								Optional:            true,
								MarkdownDescription: "Rules for monitoring the deployment and the action that should be taken when the signal threshold is met.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"signal": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "The signal to monitor. Valid values are: `rollback` or `ineligible`.",
											Validators: []validator.String{
												stringvalidator.OneOf("rollback", "ineligible"),
											},
										},
										"threshold": schema.Int32Attribute{
											Optional: true,
											MarkdownDescription: "The percentage of devices that trigger the action. " +
												"Required for `rollback` signal. Must not be set when `action` is `offerFallback` " +
												"(the fallback is offered to all ineligible devices unconditionally).",
											Validators: []validator.Int32{
												int32validator.AtLeast(1),
											},
										},
										"action": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "The action to take when the monitoring threshold is met. Valid values are: `alertError`, `offerFallback`, `pauseDeployment`.",
											Validators: []validator.String{
												stringvalidator.OneOf("alertError", "offerFallback", "pauseDeployment"),
											},
										},
									},
								},
							},
						},
					},
					"user_experience": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "User experience settings for the deployment. These settings control how the update is presented to end users.",
						Attributes: map[string]schema.Attribute{
							"days_until_forced_reboot": schema.Int32Attribute{
								Optional:            true,
								MarkdownDescription: "Number of days after installation before the device is forced to reboot. Valid for expedited quality updates. If not specified, uses device policy defaults.",
								Validators: []validator.Int32{
									int32validator.AtLeast(0),
								},
							},
							"offer_as_optional": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether to offer the update as optional to end users. When `true`, users can choose when to install. When `false` (default), the update is offered as recommended.",
							},
						},
					},
					"expedite": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Expedite settings for quality updates. Used to bypass Windows Update for Business deferral policies and deploy updates as quickly as possible.",
						Attributes: map[string]schema.Attribute{
							"is_expedited": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether this deployment should be expedited. When `true`, the update overrides device policies and installs as quickly as possible.",
							},
							"is_readiness_test": schema.BoolAttribute{
								Optional:            true,
								MarkdownDescription: "Whether this is a readiness test for expedited updates. Used to verify devices meet prerequisites before actual expedited deployment.",
							},
						},
					},
					"content_applicability": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Content applicability settings for the deployment, including safeguard configurations.",
						Attributes: map[string]schema.Attribute{
							"safeguard": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "Safeguard settings to control which safeguard holds are applied to the deployment.",
								Attributes: map[string]schema.Attribute{
									"disabled_safeguard_profiles": schema.SetNestedAttribute{
										Optional:            true,
										MarkdownDescription: "List of safeguard profiles to disable for this deployment. By default, all safeguards are applied.",
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"category": schema.StringAttribute{
													Required:            true,
													MarkdownDescription: "The category of safeguard to disable. Valid values are: `likelyIssues`.",
													Validators: []validator.String{
														stringvalidator.OneOf("likelyIssues"),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the deployment was created.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the deployment was last modified.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}

// ConfigValidators returns resource-level validators applied before plan/apply.
func (r *WindowsUpdatesAutopatchDeploymentResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	monitoringRules := path.Root("settings").AtName("monitoring").AtName("monitoring_rules")
	return []resource.ConfigValidator{
		// The "ineligible" signal must be paired with the "offerFallback" action (and vice-versa).
		// All other signal/action combinations are unconstrained.
		resourcevalidator.SetNestedRequiredPairs(
			path.Root("settings"),
			monitoringRules,
			"signal",
			"action",
			map[string]string{
				"ineligible": "offerFallback",
			},
		),
		// "offerFallback" applies to all ineligible devices unconditionally — threshold is not accepted.
		resourcevalidator.SetNestedFieldNullWhen(
			path.Root("settings"),
			monitoringRules,
			"action",
			"offerFallback",
			"threshold",
		),
	}
}
