package graphBetaSettingsCatalogConfigurationPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ListResourceName = "microsoft365_graph_beta_device_management_settings_catalog_configuration_policy"
)

var (
	_ list.ListResource              = &SettingsCatalogListResource{}
	_ list.ListResourceWithConfigure = &SettingsCatalogListResource{}
)

func NewSettingsCatalogListResource() list.ListResource {
	return &SettingsCatalogListResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		ResourcePath: "/deviceManagement/configurationPolicies",
	}
}

type SettingsCatalogListResource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
	ResourcePath    string
}

// Metadata returns the list resource type name.
func (r *SettingsCatalogListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ListResourceName
}

// Configure sets the client for the list resource.
func (r *SettingsCatalogListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ListResourceName)
	if r.client != nil {
		tflog.Debug(ctx, "Successfully configured list resource client", map[string]any{
			"list_resource": ListResourceName,
		})
	} else {
		tflog.Error(ctx, "Failed to configure list resource client - client is nil", map[string]any{
			"list_resource": ListResourceName,
		})
	}
}

// ListResourceConfigSchema defines the schema for the list resource configuration.
func (r *SettingsCatalogListResource) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		MarkdownDescription: "Configuration for listing Settings Catalog configuration policies in Microsoft Intune. " +
			"The list resource automatically retrieves ALL policies (across all pages) and returns the complete dataset. " +
			"Filter parameters narrow down which policies are returned from Microsoft Graph, optimizing the discovery process. " +
			"For full resource details including assignments and settings, use Terraform's import functionality with `terraform plan -generate-config-out`.",
		Attributes: map[string]listschema.Attribute{
			"name_filter": listschema.StringAttribute{
				MarkdownDescription: "Filter policies by name using partial matching. Supports the OData `contains` operator. " +
					"Example: `name_filter = \"Kerberos\"` will match \"[Base] Prod | Windows - Settings Catalog | Kerberos ver1.0\".",
				Optional: true,
			},
			"platform_filter": listschema.ListAttribute{
				ElementType: types.StringType,
				MarkdownDescription: "Filter policies by platform(s). Valid values: " +
					"`none`, `android`, `iOS`, `macOS`, `windows10X`, `windows10`, `linux`, `unknownFutureValue`, `androidEnterprise`, `aosp`. " +
					"Multiple platforms use OR logic. Example: `platform_filter = [\"windows10\", \"macOS\"]`.",
				Optional: true,
			},
			"template_family_filter": listschema.StringAttribute{
				MarkdownDescription: "Filter policies by template family. Valid values: " +
					"`none`, `endpointSecurityAntivirus`, `endpointSecurityDiskEncryption`, `endpointSecurityFirewall`, " +
					"`endpointSecurityEndpointDetectionAndResponse`, `endpointSecurityAttackSurfaceReduction`, `endpointSecurityAccountProtection`, " +
					"`endpointSecurityApplicationControl`, `endpointSecurityEndpointPrivilegeManagement`, `enrollmentConfiguration`, " +
					"`appQuietTime`, `baseline`, `unknownFutureValue`, `deviceConfigurationScripts`, `deviceConfigurationPolicies`, " +
					"`windowsOsRecoveryPolicies`, `companyPortal`. Example: `template_family_filter = \"baseline\"`.",
				Optional: true,
			},
			"is_assigned_filter": listschema.BoolAttribute{
				MarkdownDescription: "Filter policies by assignment status. Set to `true` to return only policies with " +
					"assignments, `false` for policies without assignments. This filter queries the assignments endpoint " +
					"for each policy (the API's `isAssigned` field is unreliable) and may take 20-30 seconds for large tenants.",
				Optional: true,
			},
			"odata_filter": listschema.StringAttribute{
				MarkdownDescription: "Advanced: Custom OData $filter query for complex filtering scenarios. " +
					"Allows direct control over the API filter expression. " +
					"Example: `odata_filter = \"platforms eq 'windows10' and isAssigned eq true\"`. " +
					"When specified, this overrides individual filter parameters. " +
					"See Microsoft Graph API documentation for supported operators and syntax.",
				Optional: true,
			},
		},
	}
}
