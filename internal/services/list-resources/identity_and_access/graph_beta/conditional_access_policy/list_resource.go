package graphBetaConditionalAccessPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ListResourceName = "microsoft365_graph_beta_identity_and_access_conditional_access_policy"
)

var (
	_ list.ListResource              = &ConditionalAccessPolicyListResource{}
	_ list.ListResourceWithConfigure = &ConditionalAccessPolicyListResource{}
)

func NewConditionalAccessPolicyListResource() list.ListResource {
	return &ConditionalAccessPolicyListResource{
		ReadPermissions: []string{
			"Policy.Read.All",
		},
		ResourcePath: "/identity/conditionalAccess/policies",
	}
}

type ConditionalAccessPolicyListResource struct {
	client          *msgraphbetasdk.GraphServiceClient
	ReadPermissions []string
	ResourcePath    string
}

// Metadata returns the list resource type name.
func (r *ConditionalAccessPolicyListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ListResourceName
}

// Configure sets the client for the list resource.
func (r *ConditionalAccessPolicyListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ConditionalAccessPolicyListResource) ListResourceConfigSchema(ctx context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		MarkdownDescription: "Lists Conditional Access policies from Microsoft Entra ID using the `/identity/conditionalAccess/policies` endpoint. " +
			"This list resource is used to automatically retrieve all policies across multiple pages with advanced filtering capabilities for policy discovery and import. " +
			"For full resource details, use Terraform's import functionality with `terraform plan -generate-config-out`.",
		Attributes: map[string]listschema.Attribute{
			"display_name_filter": listschema.StringAttribute{
				MarkdownDescription: "Filter policies by display name using partial matching. Supports the OData `contains` operator. " +
					"Example: `display_name_filter = \"MFA\"` will match \"Require MFA for Admins\".",
				Optional: true,
			},
			"state_filter": listschema.StringAttribute{
				MarkdownDescription: "Filter policies by state. Valid values: `enabled`, `disabled`, `enabledForReportingButNotEnforced`. " +
					"Example: `state_filter = \"enabled\"` returns only enabled policies.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled", "enabledForReportingButNotEnforced"),
				},
			},
			"odata_filter": listschema.StringAttribute{
				MarkdownDescription: "Advanced: Custom OData $filter query for complex filtering scenarios. " +
					"Allows direct control over the API filter expression. " +
					"Example: `odata_filter = \"state eq 'enabled' and displayName eq 'MFA Policy'\"`. " +
					"When specified, this overrides individual filter parameters. " +
					"See Microsoft Graph API documentation for supported operators and syntax.",
				Optional: true,
			},
		},
	}
}
