package graphBetaAgentUser

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_agents_agent_user"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &AgentUserResource{}
	_ resource.ResourceWithConfigure   = &AgentUserResource{}
	_ resource.ResourceWithImportState = &AgentUserResource{}
)

func NewAgentUserResource() resource.Resource {
	return &AgentUserResource{
		ReadPermissions: []string{
			"User.Read.All",
			"AgentIdUser.ReadWrite.All", // There's no such thing as AgentIdUser.Read.All
			"Directory.Read.All",
		},
		WritePermissions: []string{
			"AgentIdUser.ReadWrite.IdentityParentedBy",
			"User.ReadWrite.All",
			"AgentIdUser.ReadWrite.All",
			"User.DeleteRestore.All",
			"CustomSecAttributeAssignment.Read.All ",
		},
		ResourcePath: "/users",
	}
}

type AgentUserResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *AgentUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *AgentUserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState handles importing the resource with an extended ID format.
//
// Supported formats:
//   - Simple:   "resource_id" (hard_delete defaults to false)
//   - Extended: "resource_id:hard_delete=true" or "resource_id:hard_delete=false"
//
// Example:
//
//	terraform import microsoft365_graph_beta_agents_agent_user.example "12345678-1234-1234-1234-123456789012:hard_delete=true"
func (r *AgentUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")
	resourceID := idParts[0]
	hardDelete := false // Default to soft delete for safety

	if len(idParts) > 1 {
		for _, part := range idParts[1:] {
			if strings.HasPrefix(part, "hard_delete=") {
				value := strings.TrimPrefix(part, "hard_delete=")
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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), resourceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("hard_delete"), hardDelete)...)
}

func (r *AgentUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft 365 users using the `/users/microsoft.graph.agentUser` endpoint. " +
			"Represents a specialized subtype of user identity in Microsoft Entra ID designed for AI-powered applications (agents) " +
			"that need to function as digital workers. Agent users enable agents to access APIs and services that specifically require user identities, " +
			"receiving tokens with `idtyp=user` claims. Agent users are distinct from human users and they only interlinked to users through relationships " +
			"such as owner, sponsor, and manager.\n\n" +
			"Each agent user maintains a one-to-one relationship with a parent agent identity and is authenticated through that parent's credentials. " +
			"Agent users have user-like capabilities such as being added to groups, assigned licenses, and accessing collaborative features like mailboxes and chat, " +
			"while operating under security constraints including no password authentication, no privileged admin role assignments, and permissions similar to guest users.\n\n" +
			"For more information, see the [agentUser resource type](https://learn.microsoft.com/en-us/graph/api/resources/agentuser?view=graph-rest-beta).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the agent user. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The name displayed in the address book for the agent user. This value is usually the combination of the user's first name, middle initial, and last name. " +
					"This property is required when an agent user is created and it cannot be cleared during updates. Maximum length is 256 characters.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"agent_identity_id": schema.StringAttribute{
				MarkdownDescription: "The object ID of the agent identity that this agent user is associated with. " +
					"This creates a one-to-one relationship with the agent identity. The agent user authenticates through the " +
					"parent agent identity's credentials. Required.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_enabled": schema.BoolAttribute{
				MarkdownDescription: "Set to `true` if the agent user account is enabled; otherwise, `false`. This property is required when a user is created.",
				Required:            true,
			},
			"user_principal_name": schema.StringAttribute{
				MarkdownDescription: "The user principal name (UPN) of the agent user. The UPN is an Internet-style sign-in name for the user based on the Internet standard RFC 822. " +
					"By convention, this should map to the agent user's email name. The general format is alias@domain, where the domain must be present in the tenant's verified domain collection. " +
					"This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of organization. " +
					"NOTE: This property can't contain accent characters. Only the following characters are allowed A - Z, a - z, 0 - 9, ' . - _ ! # ^ ~.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.UserPrincipalNameRegex),
						"must be a valid user principal name in the format alias@domain. Only the following characters are allowed in the alias: A-Z, a-z, 0-9, ' . - _ ! # ^ ~",
					),
				},
			},
			"mail_nickname": schema.StringAttribute{
				MarkdownDescription: "The mail alias for the agent user. This property must be specified when a user is created. Maximum length is 64 characters.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 64),
				},
			},
			"mail": schema.StringAttribute{
				MarkdownDescription: "The SMTP address for the agent user, for example, jeff@contoso.com. Read-only.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"given_name": schema.StringAttribute{
				MarkdownDescription: "The given name (first name) of the agent user. Maximum length is 64 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
				},
			},
			"surname": schema.StringAttribute{
				MarkdownDescription: "The user's surname (family name or last name). Maximum length is 64 characters.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
				},
			},
			"job_title": schema.StringAttribute{
				MarkdownDescription: "The agent user's job title. Maximum length is 128 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"department": schema.StringAttribute{
				MarkdownDescription: "The name of the department in which the agent user works. Maximum length is 64 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
				},
			},
			"company_name": schema.StringAttribute{
				MarkdownDescription: "The company name which the agent user is associated. This property can be useful for describing the company that an external user comes from. Maximum length is 64 characters.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
				},
			},
			"office_location": schema.StringAttribute{
				MarkdownDescription: "The office location in the agent user's place of business. Maximum length is 128 characters.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"city": schema.StringAttribute{
				MarkdownDescription: "The city in which the agent user is located. Maximum length is 128 characters.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state or province in the agent user's address. Maximum length is 128 characters.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"country": schema.StringAttribute{
				MarkdownDescription: "The country/region in which the agent user is located; for example, US or UK.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(128),
				},
			},
			"postal_code": schema.StringAttribute{
				MarkdownDescription: "The postal code for the agent user's postal address. The postal code is specific to the user's country/region. Maximum length is 40 characters.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(40),
				},
			},
			"street_address": schema.StringAttribute{
				MarkdownDescription: "The street address of the agent user's place of business. Maximum length is 1024 characters.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1024),
				},
			},
			"usage_location": schema.StringAttribute{
				MarkdownDescription: "A two-letter country code (ISO standard 3166). Required for users that are assigned " +
					"licenses due to legal requirements to check for availability of services in countries. " +
					"Examples include: `US`, `JP`, and `GB`. Not nullable.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"AD", "AE", "AF", "AG", "AI", "AL", "AM", "AO", "AQ", "AR", "AS", "AT", "AU", "AW", "AX", "AZ",
						"BA", "BB", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BL", "BM", "BN", "BO", "BQ", "BR", "BS", "BT", "BV", "BW", "BY", "BZ",
						"CA", "CC", "CD", "CF", "CG", "CH", "CI", "CK", "CL", "CM", "CN", "CO", "CR", "CU", "CV", "CW", "CX", "CY", "CZ",
						"DE", "DJ", "DK", "DM", "DO", "DZ",
						"EC", "EE", "EG", "EH", "ER", "ES", "ET",
						"FI", "FJ", "FK", "FM", "FO", "FR",
						"GA", "GB", "GD", "GE", "GF", "GG", "GH", "GI", "GL", "GM", "GN", "GP", "GQ", "GR", "GS", "GT", "GU", "GW", "GY",
						"HK", "HM", "HN", "HR", "HT", "HU",
						"ID", "IE", "IL", "IM", "IN", "IO", "IQ", "IR", "IS", "IT",
						"JE", "JM", "JO", "JP",
						"KE", "KG", "KH", "KI", "KM", "KN", "KP", "KR", "KW", "KY", "KZ",
						"LA", "LB", "LC", "LI", "LK", "LR", "LS", "LT", "LU", "LV", "LY",
						"MA", "MC", "MD", "ME", "MF", "MG", "MH", "MK", "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ",
						"NA", "NC", "NE", "NF", "NG", "NI", "NL", "NO", "NP", "NR", "NU", "NZ",
						"OM",
						"PA", "PE", "PF", "PG", "PH", "PK", "PL", "PM", "PN", "PR", "PS", "PT", "PW", "PY",
						"QA",
						"RE", "RO", "RS", "RU", "RW",
						"SA", "SB", "SC", "SD", "SE", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SR", "SS", "ST", "SV", "SX", "SY", "SZ",
						"TC", "TD", "TF", "TG", "TH", "TJ", "TK", "TL", "TM", "TN", "TO", "TR", "TT", "TV", "TW", "TZ",
						"UA", "UG", "UM", "US", "UY", "UZ",
						"VA", "VC", "VE", "VG", "VI", "VN", "VU",
						"WF", "WS",
						"YE", "YT",
						"ZA", "ZM", "ZW",
					),
				},
			},
			"preferred_language": schema.StringAttribute{
				MarkdownDescription: "The preferred language for the agent user. The preferred language format is based on ISO 639-1 Code; for example en-US.",
				Optional:            true,
				Computed:            true,
			},
			"user_type": schema.StringAttribute{
				MarkdownDescription: "A string value that can be used to classify user types in your directory, such as `Member` and `Guest`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Member", "Guest"),
				},
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the agent user was created. The value cannot be modified and is automatically populated when the entity is created. Read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_type": schema.StringAttribute{
				MarkdownDescription: "Indicates whether the agent user account was created through one of the following methods: As a regular school or work account (null), As an external account (Invitation), " +
					"As a local account for an Azure Active Directory B2C tenant (LocalAccount), Through self-service sign-up by an internal user using email verification (EmailVerified), " +
					"Through self-service sign-up by an external user signing up through a link that is part of a user flow (SelfServiceSignUp). Read-only.",
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sponsor_ids": schema.SetAttribute{
				MarkdownDescription: "The users and groups responsible for this agent user's privileges in the tenant and keep the agent user's information and access updated. Required.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"hard_delete": schema.BoolAttribute{
				MarkdownDescription: "When set to `true`, the resource will be permanently deleted from the Entra ID (hard delete) " +
					"rather than being moved to deleted items (soft delete). This prevents the resource from being restored " +
					"and immediately frees up the resource name for reuse. When `false` (default), the resource is soft deleted and can be restored within 30 days. " +
					"Note: This field defaults to `false` on import since the API does not return this value.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"timeouts": commonschema.Timeouts(ctx),
		},
	}
}
