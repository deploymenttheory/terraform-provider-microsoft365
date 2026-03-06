package graphBetaWindowsAutopatchDeployment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
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
	ResourceName  = "microsoft365_graph_beta_device_management_windows_autopatch_deployment"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &WindowsAutopatchDeploymentResource{}
	_ resource.ResourceWithConfigure   = &WindowsAutopatchDeploymentResource{}
	_ resource.ResourceWithImportState = &WindowsAutopatchDeploymentResource{}
	_ resource.ResourceWithIdentity    = &WindowsAutopatchDeploymentResource{}
)

func NewWindowsAutopatchDeploymentResource() resource.Resource {
	return &WindowsAutopatchDeploymentResource{
		ReadPermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		WritePermissions: []string{
			"WindowsUpdates.ReadWrite.All",
		},
		ResourcePath: "/admin/windows/updates/deployments",
	}
}

type WindowsAutopatchDeploymentResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *WindowsAutopatchDeploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *WindowsAutopatchDeploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *WindowsAutopatchDeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *WindowsAutopatchDeploymentResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *WindowsAutopatchDeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Windows Update deployments for deploying content to a set of devices. " +
			"Deployments define which update content (feature or quality updates) should be deployed, to which audience, " +
			"and with what settings (schedule, monitoring, etc.). " +
			"See the [Microsoft Graph API documentation](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-deployment) for more information.",
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
						Required:            true,
						MarkdownDescription: "The ID of the catalog entry to deploy. This should reference a feature or quality update from the Windows Update catalog.",
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
				Optional:            true,
				MarkdownDescription: "Settings specified on the deployment governing how to deploy content.",
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
							"monitoring_rules": schema.ListNestedAttribute{
								Optional:            true,
								MarkdownDescription: "Rules for monitoring the deployment and taking action based on signals.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"signal": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "The signal to monitor. Valid values are: `rollback`.",
											Validators: []validator.String{
												stringvalidator.OneOf("rollback"),
											},
										},
										"threshold": schema.Int32Attribute{
											Required:            true,
											MarkdownDescription: "The threshold value that triggers the action.",
											Validators: []validator.Int32{
												int32validator.AtLeast(1),
											},
										},
										"action": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "The action to take when the threshold is met. Valid values are: `pauseDeployment`, `alertError`.",
											Validators: []validator.String{
												stringvalidator.OneOf("pauseDeployment", "alertError"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"state": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Execution status of the deployment.",
				Attributes: map[string]schema.Attribute{
					"requested_value": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The requested state value. Valid values are: `none`, `paused`, `archived`. Use `paused` to pause the deployment.",
						Validators: []validator.String{
							stringvalidator.OneOf("none", "paused", "archived"),
						},
					},
					"effective_value": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The effective state value of the deployment. Possible values: `scheduled`, `offering`, `paused`, `faulted`, `archived` (read-only).",
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
