package graphBetaIdentityAndAccessLifecycleWorkflowCustomTaskExtension

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_custom_task_extension"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

const (
	// authTypeProofOfPossession selects azureAdPopTokenAuthentication (Proof of Possession token).
	authTypeProofOfPossession = "azureAdPopTokenAuthentication"
	// authTypeAzureAd selects azureAdTokenAuthentication (normal token for resource_id).
	authTypeAzureAd = "azureAdTokenAuthentication"

	defaultMaximumRetries        = 1
	defaultTimeoutInMilliseconds = 1000
)

// isoDurationRegex accepts an ISO 8601 duration such as PT30M, PT1H or P0DT2H.
var isoDurationRegex = regexp.MustCompile(`^P(\d+D)?(T(\d+H)?(\d+M)?(\d+S)?)?$`)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &LifecycleWorkflowCustomTaskExtensionResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &LifecycleWorkflowCustomTaskExtensionResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &LifecycleWorkflowCustomTaskExtensionResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &LifecycleWorkflowCustomTaskExtensionResource{}

	// Enables cross-attribute validation of the authentication configuration
	_ resource.ResourceWithValidateConfig = &LifecycleWorkflowCustomTaskExtensionResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &LifecycleWorkflowCustomTaskExtensionResource{}
)

// NewLifecycleWorkflowCustomTaskExtensionResource returns a new lifecycle workflow custom task extension resource.
func NewLifecycleWorkflowCustomTaskExtensionResource() resource.Resource {
	return &LifecycleWorkflowCustomTaskExtensionResource{
		ReadPermissions: []string{
			"LifecycleWorkflows-CustomExt.Read.All",
		},
		WritePermissions: []string{
			"LifecycleWorkflows-CustomExt.ReadWrite.All",
		},
	}
}

type LifecycleWorkflowCustomTaskExtensionResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *LifecycleWorkflowCustomTaskExtensionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *LifecycleWorkflowCustomTaskExtensionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *LifecycleWorkflowCustomTaskExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *LifecycleWorkflowCustomTaskExtensionResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// clientConfigurationAttrTypes are the attribute types of client_configuration, used for its default.
var clientConfigurationAttrTypes = map[string]attr.Type{
	"maximum_retries":         types.Int64Type,
	"timeout_in_milliseconds": types.Int64Type,
}

// Schema defines the schema for the resource.
func (r *LifecycleWorkflowCustomTaskExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages Microsoft Entra ID Governance lifecycle workflow custom task extensions using the `/identityGovernance/lifecycleWorkflows/customTaskExtensions` endpoint. " +
			"A custom task extension lets a lifecycle workflow call an Azure Logic App. Graph does not check that the subscription, resource group or Logic App exist when the extension is created or updated.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the custom task extension.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the custom task extension.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the custom task extension. Defaults to an empty string, which is what Graph returns when no description is set.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"endpoint_configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "The Azure Logic App the custom task extension calls (`logicAppTriggerEndpointConfiguration`).",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"subscription_id": schema.StringAttribute{
						MarkdownDescription: "The ID of the Azure subscription that contains the Logic App.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID"),
						},
					},
					"resource_group_name": schema.StringAttribute{
						MarkdownDescription: "The name of the Azure resource group that contains the Logic App.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"logic_app_workflow_name": schema.StringAttribute{
						MarkdownDescription: "The name of the Logic App workflow.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"url": schema.StringAttribute{
						MarkdownDescription: "The URL of the Logic App HTTP trigger, for example `https://prod-00.westeurope.logic.azure.com:443/workflows/<id>/triggers/manual/paths/invoke?api-version=2016-10-01`. " +
							"Required when the provider authenticates as an application (client credentials): Graph rejects an app-only create without it (`Custom task extensions missing url`). " +
							"When omitted, the value Graph returns is kept in state.",
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(`^https://`), "must be an https URL"),
						},
					},
				},
			},
			"authentication_configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "How Microsoft Entra ID secures the call to the Logic App.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "The token type. `azureAdPopTokenAuthentication` sends a Proof of Possession token (the Entra admin center default); " +
							"`azureAdTokenAuthentication` sends a normal access token for `resource_id`.",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf(authTypeProofOfPossession, authTypeAzureAd),
						},
					},
					"resource_id": schema.StringAttribute{
						MarkdownDescription: "The application ID (client ID) of the Microsoft Entra application the token is issued for. Required when `type` is `azureAdTokenAuthentication`, and must not be set otherwise.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
			"client_configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "How long Microsoft Entra ID waits for the Logic App to accept the call and how often it retries. Defaults to `maximum_retries = 1` and `timeout_in_milliseconds = 1000`.",
				Optional:            true,
				Computed:            true,
				Default: objectdefault.StaticValue(types.ObjectValueMust(clientConfigurationAttrTypes, map[string]attr.Value{
					"maximum_retries":         types.Int64Value(defaultMaximumRetries),
					"timeout_in_milliseconds": types.Int64Value(defaultTimeoutInMilliseconds),
				})),
				Attributes: map[string]schema.Attribute{
					"maximum_retries": schema.Int64Attribute{
						MarkdownDescription: "The number of times Microsoft Entra ID retries a call that timed out. Lifecycle workflows accept `0` or `1`. Defaults to `1`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(defaultMaximumRetries),
						Validators: []validator.Int64{
							int64validator.Between(0, 1),
						},
					},
					"timeout_in_milliseconds": schema.Int64Attribute{
						MarkdownDescription: "How long, in milliseconds, Microsoft Entra ID waits for the Logic App to respond before it closes the connection. Must be between `200` and `2000`. Defaults to `1000`.",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(defaultTimeoutInMilliseconds),
						Validators: []validator.Int64{
							int64validator.Between(200, 2000),
						},
					},
				},
			},
			"callback_configuration": schema.SingleNestedAttribute{
				MarkdownDescription: "When set, the workflow waits for the Logic App to call back (`resume`) before it continues, and `reply_mode` is `callback`. " +
					"When omitted, the workflow does not wait and `reply_mode` is `none`.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"timeout_duration": schema.StringAttribute{
						MarkdownDescription: "How long the workflow waits for the callback, as an ISO 8601 duration. Lifecycle workflows accept five minutes to three hours, for example `PT30M` or `PT3H`. " +
							"Graph returns the duration in its shortest form, so an equivalent value such as `PT60M` is kept as written.",
						Required: true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(isoDurationRegex, "must be an ISO 8601 duration such as PT30M or PT1H"),
						},
					},
					"authorized_app_ids": schema.SetAttribute{
						MarkdownDescription: "The application IDs (client IDs) of the applications allowed to resume a task processing result for this extension, " +
							"for example the client ID of the Logic App's managed identity.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.RegexMatches(regexp.MustCompile(constants.GuidRegex), "must be a valid GUID"),
							),
						},
					},
				},
			},
			"reply_mode": schema.StringAttribute{
				MarkdownDescription: "How the custom task extension replies to the lifecycle workflows service: `callback` when `callback_configuration` is set, otherwise `none`.",
				Computed:            true,
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "When the custom task extension was created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified_date_time": schema.StringAttribute{
				MarkdownDescription: "When the custom task extension was last modified.",
				Computed:            true,
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
