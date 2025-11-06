package graphBetaAzureNetworkConnection

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_windows_365_azure_network_connection"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &CloudPcOnPremisesConnectionResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CloudPcOnPremisesConnectionResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CloudPcOnPremisesConnectionResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &CloudPcOnPremisesConnectionResource{}
)

func NewCloudPcOnPremisesConnectionResource() resource.Resource {
	return &CloudPcOnPremisesConnectionResource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
		},
		WritePermissions: []string{
			"CloudPC.ReadWrite.All",
		},
	}
}

type CloudPcOnPremisesConnectionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *CloudPcOnPremisesConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

// FullTypeName returns the full resource type name in the format "providername_resourcename".
func (r *CloudPcOnPremisesConnectionResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

// Configure sets the client for the resource.
func (r *CloudPcOnPremisesConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

// ImportState imports the resource state.
func (r *CloudPcOnPremisesConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *CloudPcOnPremisesConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the Azure network connection. Read-only.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name for the Azure network connection.",
			},
			"connection_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Specifies the method by which a provisioned Cloud PC is joined to Microsoft Entra. The azureADJoin option indicates the absence of an on-premises Active Directory (AD) in the current tenant that results in the Cloud PC device only joining to Microsoft Entra. The hybridAzureADJoin option indicates the presence of an on-premises AD in the current tenant and that the Cloud PC joins both the on-premises AD and Microsoft Entra. The selected option also determines the types of users who can be assigned and can sign into a Cloud PC. The azureADJoin option allows both cloud-only and hybrid users to be assigned and sign in, whereas hybridAzureADJoin is restricted to hybrid users only. Default: hybridAzureADJoin. Possible values: hybridAzureADJoin, azureADJoin, unknownFutureValue.",
				Validators: []validator.String{
					stringvalidator.OneOf("hybridAzureADJoin", "azureADJoin", "unknownFutureValue"),
				},
			},
			"ad_domain_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The fully qualified domain name (FQDN) of the Active Directory domain you want to join. Optional.",
			},
			"ad_domain_username": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The username of an Active Directory account (user or service account) that has permissions to create computer objects in Active Directory. Required format: admin@contoso.com. Optional.",
			},
			"ad_domain_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "The password associated with adDomainUsername.",
			},
			"organizational_unit": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The organizational unit (OU) in which the computer account is created. If left null, the OU configured as the default (a well-known computer object container) in your Active Directory domain (OU) is used. Optional.",
			},
			"resource_group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the target resource group. Required format: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^/subscriptions/[a-f0-9\-]+/resourceGroups/[a-z0-9\-_.()]+$`),
						"Must be a valid Azure resource group ID in lowercase (e.g., /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName})",
					),
				},
			},
			"subnet_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the target subnet. Required format: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{virtualNetworkId}/subnets/{subnetName}.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^/subscriptions/[a-f0-9\-]+/resourceGroups/[a-z0-9\-_.()]+/providers/Microsoft.Network/virtualNetworks/[a-z0-9\-_.()]+/subnets/[a-z0-9\-_.()]+$`),
						"Must be a valid Azure subnet ID in lowercase (e.g., /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{vnet}/subnets/{subnet})",
					),
				},
			},
			"subscription_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the target Azure subscription associated with your tenant.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format 00000000-0000-0000-0000-000000000000",
					),
				},
			},
			"virtual_network_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the target virtual network. Required format: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{virtualNetworkName}.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^/subscriptions/[a-f0-9\-]+/resourceGroups/[a-z0-9\-_.()]+/providers/Microsoft.Network/virtualNetworks/[a-z0-9\-_.()]+$`),
						"Must be a valid Azure virtual network ID in lowercase (e.g., /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualNetworks/{vnet})",
					),
				},
			},
			"health_check_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the most recent health check done on the Azure network connection. For example, if the status is passed, the Azure network connection passed all checks run by the service. Possible values: pending, running, passed, failed, warning, informational, unknownFutureValue. Read-only.",
				Validators: []validator.String{
					stringvalidator.OneOf("pending", "running", "passed", "failed", "warning", "informational", "unknownFutureValue"),
				},
			},
			"managed_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies which services manage the Azure network connection. Possible values: windows365, devBox, unknownFutureValue, rpaBox. Read-only.",
				Validators: []validator.String{
					stringvalidator.OneOf("windows365", "devBox", "unknownFutureValue", "rpaBox"),
				},
			},
			"in_use": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Indicates whether a Cloud PC is using this on-premises network connection. true if at least one Cloud PC is using it. Otherwise, false. Read-only. Default is false.",
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
