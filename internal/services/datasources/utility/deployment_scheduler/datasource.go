package utilityDeploymentScheduler

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	DataSourceName = "microsoft365_utility_deployment_scheduler"
	ReadTimeout    = 30
)

var (
	// Basic datasource interface (Read operations)
	_ datasource.DataSource = &deploymentSchedulerDataSource{}

	// Allows the datasource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &deploymentSchedulerDataSource{}
)

func NewDeploymentSchedulerDataSource() datasource.DataSource {
	return &deploymentSchedulerDataSource{}
}

type deploymentSchedulerDataSource struct{}

func (d *deploymentSchedulerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

// Configure implements the DataSourceWithConfigure interface.
// For utility datasources that perform local computations (like deployment scheduling),
// this method doesn't need to extract Microsoft Graph clients from ProviderData. However, it's still
// required for interface compliance and maintains consistency across all datasources in the provider.
// This pattern allows for future flexibility if the datasource later needs access to provider configuration.
func (d *deploymentSchedulerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
}

func (d *deploymentSchedulerDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A conditional gate valve for phased deployments. Returns scope ID(s) when specified time-based conditions are met, " +
			"enabling controlled rollout of policies and updates. All times are in UTC (RFC3339 format). " +
			"When conditions are not met, returns null, preventing deployment to target groups.\n\n" +
			"This datasource is evaluated on every Terraform plan/apply, allowing gates to automatically open when conditions are satisfied.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of this deployment scheduler instance.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "A descriptive name for this deployment phase (e.g., 'Phase 2 - Production Rollout'). Used in status messages.",
			},
			"deployment_start_time": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The deployment campaign start date and time in UTC (RFC3339 format, e.g., '2024-01-15T00:00:00Z'). " +
					"All time-based conditions are calculated relative to this timestamp, similar to how Unix epoch time works as a reference point. " +
					"If not provided, uses the current time on each evaluation (not recommended for time-based conditions). " +
					"Explicitly setting this value allows coordinating multiple deployment phases to a single campaign start time.",
				Validators: []validator.String{
					attribute.RegexMatches(
						regexp.MustCompile(constants.ISO8601DateTimeRegex),
						"Must be a valid RFC3339 date/time in UTC (e.g., '2024-01-15T00:00:00Z')",
					),
				},
			},
			"scope_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A single scope ID (typically a group ID) to release when conditions are met. Use this when deploying to one group. Use either `scope_id` or `scope_ids`, not both.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("scope_ids")),
					stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "Must be a valid GUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)"),
				},
			},
			"scope_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "List of multiple scope IDs (user GUIDs, device GUIDs, etc.) to release when conditions are met. Use this when deploying to multiple individual entities. Use either `scope_id` or `scope_ids`, not both.",
				Validators: []validator.List{
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(
							regexp.MustCompile(constants.GuidRegex),
							"Each scope ID must be a valid GUID format (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)",
						),
					),
				},
			},
			"time_condition": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Time-based condition that must be satisfied before releasing scope ID(s).",
				Attributes: map[string]schema.Attribute{
					"delay_start_time_by": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Number of hours to delay after `deployment_start_time` before allowing the gate to open. Must be >= 0. Set to 0 for immediate release at deployment start time.",
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
					"absolute_earliest": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Absolute earliest time (UTC RFC3339 format) when gate can open, regardless of `delay_start_time_by`. Use this to prevent deployment before a specific date/time (e.g., wait for Patch Tuesday). If specified, gate cannot open before this time even if delay_start_time_by has elapsed.",
						Validators: []validator.String{
							attribute.RegexMatches(
								regexp.MustCompile(constants.ISO8601DateTimeRegex),
								"Must be a valid RFC3339 date/time in UTC (e.g., '2024-01-15T00:00:00Z')",
							),
						},
					},
					"absolute_latest": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Absolute deadline (UTC RFC3339 format) when gate must close and will never open again. Use this for time-limited deployment campaigns or change freeze deadlines. If current time exceeds this, gate permanently closes.",
						Validators: []validator.String{
							attribute.RegexMatches(
								regexp.MustCompile(constants.ISO8601DateTimeRegex),
								"Must be a valid RFC3339 date/time in UTC (e.g., '2024-01-15T00:00:00Z')",
							),
						},
					},
					"max_open_duration_hours": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Maximum number of hours the gate can remain open after it first opens. Must be >= 0. Use this for pilot/temporary deployments that should automatically expire (e.g., 2-week pilot = 336 hours). When duration expires, gate auto-closes and scope IDs are retracted. Set to 0 for unlimited duration (default behavior).",
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
				},
			},
			"inclusion_time_windows": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Defines one or more time windows when the gate is allowed to open. " +
					"The current time must fall within at least one of the defined windows for the gate to open. " +
					"Use this for office hours restrictions, maintenance windows, etc. " +
					"Multiple windows are evaluated with OR logic (any window matches = condition passes).",
				Attributes: map[string]schema.Attribute{
					"window": schema.ListNestedAttribute{
						Required:            true,
						MarkdownDescription: "List of time windows when deployment is allowed.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"days_of_week": schema.ListAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Days of the week when this window is active. Valid values: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`. If not specified, all days are included.",
									Validators: []validator.List{
										listvalidator.ValueStringsAre(
											attribute.RegexMatches(
												regexp.MustCompile(constants.DayOfWeekRegex),
												"Must be a lowercase day of week: monday, tuesday, wednesday, thursday, friday, saturday, or sunday",
											),
										),
									},
								},
								"time_of_day_start": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Start time in UTC (HH:MM:SS format, e.g., '09:00:00'). If not specified, starts at 00:00:00.",
									Validators: []validator.String{
										attribute.RegexMatches(
											regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
											"Must be a valid time in HH:MM:SS format (24-hour clock)",
										),
									},
								},
								"time_of_day_end": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "End time in UTC (HH:MM:SS format, e.g., '17:00:00'). If not specified, ends at 23:59:59.",
									Validators: []validator.String{
										attribute.RegexMatches(
											regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
											"Must be a valid time in HH:MM:SS format (24-hour clock)",
										),
									},
								},
								"date_start": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Absolute start date/time in UTC (RFC3339 format, e.g., '2024-01-15T00:00:00Z'). Use for specific date ranges. If not specified, no start date limit.",
									Validators: []validator.String{
										attribute.RegexMatches(
											regexp.MustCompile(constants.ISO8601DateTimeRegex),
											"Must be a valid RFC3339 date/time in UTC (e.g., '2024-01-15T00:00:00Z')",
										),
									},
								},
								"date_end": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Absolute end date/time in UTC (RFC3339 format, e.g., '2024-01-31T23:59:59Z'). Use for specific date ranges. If not specified, no end date limit.",
									Validators: []validator.String{
										attribute.RegexMatches(
											regexp.MustCompile(constants.ISO8601DateTimeRegex),
											"Must be a valid RFC3339 date/time in UTC (e.g., '2024-01-31T23:59:59Z')",
										),
									},
								},
							},
						},
					},
				},
			},
			"exclusion_time_windows": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Defines one or more time windows when the gate must remain closed, even if other conditions are met. " +
					"Use this for holiday freezes, blackout periods, etc. " +
					"Exclusions take precedence over inclusions. " +
					"Multiple windows are evaluated with OR logic (any window matches = deployment blocked).",
				Attributes: map[string]schema.Attribute{
					"window": schema.ListNestedAttribute{
						Required:            true,
						MarkdownDescription: "List of time windows when deployment is blocked.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"days_of_week": schema.ListAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Days of the week when this window blocks deployment. Valid values: `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, `saturday`, `sunday`. If not specified, all days are included.",
									Validators: []validator.List{
										listvalidator.ValueStringsAre(
											attribute.RegexMatches(
												regexp.MustCompile(constants.DayOfWeekRegex),
												"Must be a lowercase day of week: monday, tuesday, wednesday, thursday, friday, saturday, or sunday",
											),
										),
									},
								},
								"time_of_day_start": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Start time in UTC (HH:MM:SS format, e.g., '00:00:00'). If not specified, starts at 00:00:00.",
									Validators: []validator.String{
										attribute.RegexMatches(
											regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
											"Must be a valid time in HH:MM:SS format (24-hour clock)",
										),
									},
								},
								"time_of_day_end": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "End time in UTC (HH:MM:SS format, e.g., '23:59:59'). If not specified, ends at 23:59:59.",
									Validators: []validator.String{
										attribute.RegexMatches(
											regexp.MustCompile(constants.TimeFormatHHMMSSRegex),
											"Must be a valid time in HH:MM:SS format (24-hour clock)",
										),
									},
								},
								"date_start": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Absolute start date/time in UTC (RFC3339 format, e.g., '2024-12-20T00:00:00Z'). Use for specific date ranges like holiday freezes. If not specified, no start date limit.",
									Validators: []validator.String{
										attribute.RegexMatches(
											regexp.MustCompile(constants.ISO8601DateTimeRegex),
											"Must be a valid RFC3339 date/time in UTC (e.g., '2024-12-20T00:00:00Z')",
										),
									},
								},
								"date_end": schema.StringAttribute{
									Optional:            true,
									Computed:            true,
									MarkdownDescription: "Absolute end date/time in UTC (RFC3339 format, e.g., '2025-01-05T23:59:59Z'). Use for specific date ranges. If not specified, no end date limit.",
									Validators: []validator.String{
										attribute.RegexMatches(
											regexp.MustCompile(constants.ISO8601DateTimeRegex),
											"Must be a valid RFC3339 date/time in UTC (e.g., '2025-01-05T23:59:59Z')",
										),
									},
								},
							},
						},
					},
				},
			},
			"manual_override": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Emergency override to immediately release scope ID(s), bypassing all time conditions and windows. " +
					"Set to `true` to force-release the gate. Useful for emergency deployments. " +
					"When enabled, all other conditions are ignored and scope ID(s) are immediately released. " +
					"Default: false.",
			},
			"depends_on_scheduler": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "Dependency gate that requires another scheduler to have been open for a minimum duration before this gate can open. " +
					"Useful for sequential phased rollouts where Phase 3 shouldn't start until Phase 2 has been running successfully for a period. " +
					"This gate calculates when the prerequisite would have opened and ensures sufficient time has passed.",
				Attributes: map[string]schema.Attribute{
					"prerequisite_delay_start_time_by": schema.Int64Attribute{
						Required: true,
						MarkdownDescription: "The delay_start_time_by value of the prerequisite scheduler. " +
							"This is when the prerequisite gate would have opened (hours after deployment_start_time). " +
							"For example, if Phase 2 has `delay_start_time_by = 168`, set this to 168.",
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
					"minimum_open_hours": schema.Int64Attribute{
						Required: true,
						MarkdownDescription: "Minimum number of hours the prerequisite gate must have been open before this gate can open. " +
							"For example, if set to 48, this gate won't open until 48 hours after the prerequisite gate opened. " +
							"Must be >= 0.",
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
					},
				},
			},
			"require_all_conditions": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, all specified conditions must be met (AND logic). When false, any condition passing will release scope ID(s) (OR logic). Defaults to true. Reserved for future use when multiple condition types are supported.",
			},
			"released_scope_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The single scope ID released by this gate when conditions are met. Returns the value from `scope_id` when the gate opens, or null when conditions are not met. Use this in policy assignments when you provided `scope_id`.",
			},
			"released_scope_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of scope IDs released by this gate when conditions are met. Returns the full `scope_ids` list when the gate opens, or null when conditions are not met. Use this in policy assignments when you provided `scope_ids`.",
			},
			"condition_met": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Boolean indicating whether all required conditions are satisfied. True means the gate is open and scope IDs are released.",
			},
			"status_message": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: "Human-readable status message describing the current state of all conditions. " +
					"Visible in Terraform plan output. Examples:\n" +
					"- `Conditions met: Time condition met (50h/48h required)`\n" +
					"- `Waiting: Time condition not met (22h/48h required)`",
			},
			"conditions_detail": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Detailed breakdown of each condition's evaluation for debugging and monitoring.",
				Attributes: map[string]schema.Attribute{
					"time_condition_detail": schema.SingleNestedAttribute{
						Computed:            true,
						MarkdownDescription: "Detailed time condition evaluation.",
						Attributes: map[string]schema.Attribute{
							"required": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Whether a time condition was specified.",
							},
							"delay_start_time_by": schema.Int64Attribute{
								Computed:            true,
								MarkdownDescription: "Required delay in hours from deployment start time.",
							},
							"deployment_start_time": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The deployment start time used for calculations (UTC RFC3339).",
							},
							"current_time": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The current evaluation time (UTC RFC3339).",
							},
							"hours_elapsed": schema.Float64Attribute{
								Computed:            true,
								MarkdownDescription: "Hours elapsed since deployment start time.",
							},
							"condition_met": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Whether the time condition is satisfied.",
							},
						},
					},
				},
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
