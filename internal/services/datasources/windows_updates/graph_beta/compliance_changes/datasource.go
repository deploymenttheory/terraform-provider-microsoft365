package graphBetaWindowsUpdatesComplianceChanges

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_windows_updates_compliance_changes"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &ComplianceChangesDataSource{}
	_ datasource.DataSourceWithConfigure = &ComplianceChangesDataSource{}
)

func NewComplianceChangesDataSource() datasource.DataSource {
	return &ComplianceChangesDataSource{
		ReadPermissions: []string{
			"WindowsUpdates.Read.All",
		},
	}
}

type ComplianceChangesDataSource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
}

func (d *ComplianceChangesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *ComplianceChangesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *ComplianceChangesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves compliance changes (content approvals) for a Windows Update policy using the `/admin/windows/updates/updatePolicies/{updatePolicyId}/complianceChanges` endpoint. " +
			"This data source lists all content approvals that have been created for a policy, including their revocation status and deployment settings.",
		Attributes: map[string]schema.Attribute{
			"update_policy_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the update policy to query for compliance changes.",
			},
			"compliance_changes": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of compliance changes (content approvals) for the policy.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the compliance change.",
						},
						"created_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the compliance change was created.",
						},
						"is_revoked": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the compliance change has been revoked.",
						},
						"revoked_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The date and time when the compliance change was revoked.",
						},
						"content": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The content that was approved for deployment.",
							Attributes: map[string]schema.Attribute{
								"catalog_entry_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The ID of the catalog entry that was approved.",
								},
								"catalog_entry_type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The type of catalog entry (featureUpdate, qualityUpdate, driverUpdate).",
								},
							},
						},
						"deployment_settings": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Settings for how the approved content should be deployed.",
							Attributes: map[string]schema.Attribute{
								"schedule": schema.SingleNestedAttribute{
									Computed:            true,
									MarkdownDescription: "Schedule settings for the deployment.",
									Attributes: map[string]schema.Attribute{
										"start_date_time": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "The date and time when deployment should start.",
										},
										"gradual_rollout": schema.SingleNestedAttribute{
											Computed:            true,
											MarkdownDescription: "Gradual rollout settings.",
											Attributes: map[string]schema.Attribute{
												"duration_between_offers": schema.StringAttribute{
													Computed:            true,
													MarkdownDescription: "Duration between offers in ISO 8601 format.",
												},
												"devices_per_offer": schema.Int32Attribute{
													Computed:            true,
													MarkdownDescription: "Number of devices per offer.",
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
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
