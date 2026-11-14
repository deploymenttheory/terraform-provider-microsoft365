package graphBetaIdentityAndAccessConnectedOrganization

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_connected_organization"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

// OData types of the identity sources a connected organization can be created with.
const (
	ODataTypeDomainIdentitySource       = "#microsoft.graph.domainIdentitySource"
	ODataTypeAzureActiveDirectoryTenant = "#microsoft.graph.azureActiveDirectoryTenant"
	ODataTypeExternalDomainFederation   = "#microsoft.graph.externalDomainFederation"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &ConnectedOrganizationResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &ConnectedOrganizationResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &ConnectedOrganizationResource{}

	// Enables plan-time validation of the polymorphic identity sources
	_ resource.ResourceWithValidateConfig = &ConnectedOrganizationResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &ConnectedOrganizationResource{}
)

// https://learn.microsoft.com/en-us/graph/api/entitlementmanagement-post-connectedorganizations?view=graph-rest-beta#permissions
func NewConnectedOrganizationResource() resource.Resource {
	return &ConnectedOrganizationResource{
		ReadPermissions: []string{
			"EntitlementManagement.Read.All",
		},
		WritePermissions: []string{
			"EntitlementManagement.ReadWrite.All",
		},
	}
}

type ConnectedOrganizationResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *ConnectedOrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *ConnectedOrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *ConnectedOrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *ConnectedOrganizationResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *ConnectedOrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages entitlement management connected organizations in Microsoft Entra ID Governance using the `/identityGovernance/entitlementManagement/connectedOrganizations` endpoint. " +
			"A connected organization is an external Microsoft Entra tenant or domain whose users can request access packages.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the connected organization.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.GuidRegex),
						"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
					),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the connected organization.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the connected organization. Removing it from configuration clears the description in Microsoft Graph.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state of the connected organization. It defines whether assignment policies with the requestor scope type `AllConfiguredConnectedOrganizationSubjects` apply to it. " +
					"Possible values are: `configured`, `proposed`. Defaults to `configured`.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("configured"),
				Validators: []validator.String{
					stringvalidator.OneOf("configured", "proposed"),
				},
			},
			"identity_sources": schema.SetNestedAttribute{
				MarkdownDescription: "The identity sources of the connected organization. Microsoft Graph only accepts identity sources when the connected organization is created, " +
					"so changing this attribute destroys and recreates the connected organization.",
				Required: true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"odata_type": schema.StringAttribute{
							MarkdownDescription: "The type of identity source. Possible values are: " +
								"`#microsoft.graph.domainIdentitySource` (an external domain; requires `domain_name`), " +
								"`#microsoft.graph.azureActiveDirectoryTenant` (a Microsoft Entra tenant; requires `tenant_id`), " +
								"`#microsoft.graph.externalDomainFederation` (a domain federated with a SAML or WS-Fed identity provider configured in this tenant; requires `domain_name` and `issuer_uri`).",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									ODataTypeDomainIdentitySource,
									ODataTypeAzureActiveDirectoryTenant,
									ODataTypeExternalDomainFederation,
								),
							},
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the identity source. For a domain identity source this is typically the domain name.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"domain_name": schema.StringAttribute{
							MarkdownDescription: "The domain name. Required for `#microsoft.graph.domainIdentitySource` and `#microsoft.graph.externalDomainFederation`.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"tenant_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the Microsoft Entra tenant. Required for `#microsoft.graph.azureActiveDirectoryTenant`.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(constants.GuidRegex),
									"must be a valid GUID in the format '00000000-0000-0000-0000-000000000000'",
								),
							},
						},
						"issuer_uri": schema.StringAttribute{
							MarkdownDescription: "The issuer URI of the federation server. Required for `#microsoft.graph.externalDomainFederation`.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
					},
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "The user principal name or application name of the identity that created the connected organization.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the connected organization was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_by": schema.StringAttribute{
				MarkdownDescription: "The user principal name or application name of the identity that last modified the connected organization.",
				Computed:            true,
			},
			"modified_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the connected organization was last modified.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
