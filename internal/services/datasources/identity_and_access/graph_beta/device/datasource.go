package graphBetaDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	DataSourceName = "microsoft365_graph_beta_identity_and_access_device"
	ReadTimeout    = 180
)

var (
	_ datasource.DataSource              = &DeviceDataSource{}
	_ datasource.DataSourceWithConfigure = &DeviceDataSource{}
)

func NewDeviceDataSource() datasource.DataSource {
	return &DeviceDataSource{
		ReadPermissions: []string{
			"Directory.Read.All",
		},
	}
}

type DeviceDataSource struct {
	client *msgraphbetasdk.GraphServiceClient

	ReadPermissions []string
}

func (d *DeviceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = DataSourceName
}

func (d *DeviceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = client.SetGraphBetaClientForDataSource(ctx, req, resp, DataSourceName)
}

func (d *DeviceDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves Microsoft Entra Devices using the `/devices` endpoint. " +
			"Supports flexible lookup by object ID, display name, device ID, or custom OData queries. " +
			"Can also retrieve device memberships, registered owners, and registered users.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for the data source. This is a placeholder attribute required by Terraform.",
			},
			"object_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique object identifier of the device in Microsoft Entra ID. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("display_name"),
						path.MatchRoot("device_id"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("device_id"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The display name of the device. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"device_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The unique device identifier set by Azure Device Registration Service. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("list_all"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"list_all": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Retrieve all devices in the tenant. Conflicts with specific lookup attributes.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("device_id"),
						path.MatchRoot("odata_query"),
					),
				},
			},
			"list_member_of": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "When true and combined with object_id, retrieves the groups and administrative units that the device is a direct member of. Requires object_id to be specified.",
			},
			"list_registered_owners": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "When true and combined with object_id, retrieves the registered owners of the device. Requires object_id to be specified.",
			},
			"list_registered_users": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "When true and combined with object_id, retrieves the registered users of the device. Requires object_id to be specified.",
			},
			"odata_query": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Custom OData filter expression for advanced queries (e.g., `operatingSystem eq 'Windows' and accountEnabled eq true`). Conflicts with specific lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("object_id"),
						path.MatchRoot("display_name"),
						path.MatchRoot("device_id"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of devices matching the query criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the device object.",
						},
						"account_enabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "true if the account is enabled; otherwise, false.",
						},
						"alternative_security_ids": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Alternative security identifiers for the device.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.Int64Attribute{
										Computed:            true,
										MarkdownDescription: "The type of the alternative security identifier.",
									},
									"identity_provider": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The identity provider for the alternative security identifier.",
									},
									"key": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The key value of the alternative security identifier.",
									},
								},
							},
						},
						"approximate_last_sign_in_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The timestamp of the last sign-in activity.",
						},
						"compliance_expiration_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The timestamp when the device is no longer deemed compliant.",
						},
						"device_category": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "User-defined property set by Intune to automatically add devices to groups.",
						},
						"device_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique identifier set by Azure Device Registration Service.",
						},
						"device_metadata": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Metadata for the device.",
						},
						"device_ownership": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Ownership of the device (unknown, company, personal).",
						},
						"device_version": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Version of the device.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name for the device.",
						},
						"domain_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The on-premises domain name of the device.",
						},
						"enrollment_profile_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Enrollment profile applied to the device.",
						},
						"enrollment_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Enrollment type of the device.",
						},
						"extension_attributes": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Extension attributes 1-15 for the device.",
							Attributes: map[string]schema.Attribute{
								"extension_attribute1":  schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 1."},
								"extension_attribute2":  schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 2."},
								"extension_attribute3":  schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 3."},
								"extension_attribute4":  schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 4."},
								"extension_attribute5":  schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 5."},
								"extension_attribute6":  schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 6."},
								"extension_attribute7":  schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 7."},
								"extension_attribute8":  schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 8."},
								"extension_attribute9":  schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 9."},
								"extension_attribute10": schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 10."},
								"extension_attribute11": schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 11."},
								"extension_attribute12": schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 12."},
								"extension_attribute13": schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 13."},
								"extension_attribute14": schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 14."},
								"extension_attribute15": schema.StringAttribute{Computed: true, MarkdownDescription: "Extension attribute 15."},
							},
						},
						"is_compliant": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "true if the device complies with Mobile Device Management (MDM) policies.",
						},
						"is_managed": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "true if the device is managed by a Mobile Device Management (MDM) app.",
						},
						"is_management_restricted": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates whether the device is a member of a restricted management administrative unit.",
						},
						"is_rooted": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "true if the device is rooted or jail-broken.",
						},
						"management_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The management channel of the device.",
						},
						"manufacturer": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Manufacturer of the device.",
						},
						"mdm_app_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Application identifier used to register device into MDM.",
						},
						"model": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Model of the device.",
						},
						"on_premises_last_sync_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The last time the object was synced with the on-premises directory.",
						},
						"on_premises_security_identifier": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The on-premises security identifier (SID) for the user who was synchronized from on-premises.",
						},
						"on_premises_sync_enabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "true if this object is synced from an on-premises directory.",
						},
						"operating_system": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of operating system on the device.",
						},
						"operating_system_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The version of the operating system on the device.",
						},
						"physical_ids": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "Physical identifiers for the device.",
						},
						"profile_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The profile type of the device (RegisteredDevice, SecureVM, Printer, Shared, IoT).",
						},
						"registration_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Date and time when the device was registered.",
						},
						"system_labels": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of labels applied to the device by the system.",
						},
						"trust_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Type of trust for the joined device (Workplace, AzureAd, ServerAd).",
						},
					},
				},
			},
			"member_of": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Groups and administrative units that the device is a direct member of. Only populated when list_member_of is true.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the directory object.",
						},
						"odata_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The OData type of the directory object (e.g., #microsoft.graph.group).",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the directory object.",
						},
					},
				},
			},
			"registered_owners": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The registered owners of the device. Only populated when list_registered_owners is true.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the directory object.",
						},
						"odata_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The OData type of the directory object (e.g., #microsoft.graph.user).",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the directory object.",
						},
					},
				},
			},
			"registered_users": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The registered users of the device. Only populated when list_registered_users is true.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the directory object.",
						},
						"odata_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The OData type of the directory object (e.g., #microsoft.graph.user).",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the directory object.",
						},
					},
				},
			},
			"timeouts": commonschema.DatasourceTimeouts(ctx),
		},
	}
}
