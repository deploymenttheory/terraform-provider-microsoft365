package graphBetaIdentityAndAccessLifecycleWorkflowSettings

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_settings"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180

	// singletonID is the fixed identifier for the tenant-wide lifecycle workflow settings.
	// Graph does not return an id for this singleton.
	singletonID = "lifecycleWorkflowSettings"
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &LifecycleWorkflowSettingsResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &LifecycleWorkflowSettingsResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &LifecycleWorkflowSettingsResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &LifecycleWorkflowSettingsResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &LifecycleWorkflowSettingsResource{}
)

// https://learn.microsoft.com/en-us/graph/api/identitygovernance-lifecyclemanagementsettings-get?view=graph-rest-beta#permissions
// https://learn.microsoft.com/en-us/graph/api/identitygovernance-lifecyclemanagementsettings-update?view=graph-rest-beta#permissions
func NewLifecycleWorkflowSettingsResource() resource.Resource {
	return &LifecycleWorkflowSettingsResource{
		ReadPermissions: []string{
			"LifecycleWorkflows.Read.All",
		},
		WritePermissions: []string{
			"LifecycleWorkflows.ReadWrite.All",
		},
	}
}

type LifecycleWorkflowSettingsResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *LifecycleWorkflowSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *LifecycleWorkflowSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state using the fixed singleton identifier.
func (r *LifecycleWorkflowSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != singletonID {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Import this singleton using the literal ID "+singletonID+".",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *LifecycleWorkflowSettingsResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

// Schema defines the schema for the resource.
func (r *LifecycleWorkflowSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the tenant-wide Microsoft Entra ID Governance lifecycle workflow settings using the Microsoft Graph beta " +
			"`/identityGovernance/lifecycleWorkflows/settings` endpoint. The settings control how often scheduled lifecycle workflows run " +
			"and how lifecycle workflow email notifications are sent.\n\n" +
			"This is a **singleton resource**: one settings object exists per tenant. `create` adopts the existing settings and `create`/`update` " +
			"apply the configuration with PATCH. Settings that are not configured keep their current tenant values. On `destroy`, Terraform only " +
			"removes the resource from state and leaves the tenant settings as they are. Manage only one instance per tenant across all Terraform states.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The fixed singleton identifier `" + singletonID + "`.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workflow_schedule_interval_in_hours": schema.Int32Attribute{
				MarkdownDescription: "The interval in hours at which all scheduled workflows run. Must be between `1` and `24`. " +
					"When not configured, the current tenant value is kept.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int32{
					int32validator.Between(1, 24),
				},
			},
			"email_settings": schema.SingleNestedAttribute{
				MarkdownDescription: "The settings for email notifications sent by lifecycle workflow tasks. When not configured, the current tenant value is kept.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"sender_domain": schema.StringAttribute{
						MarkdownDescription: "The domain used to send email notifications, for example `microsoft.com` (the default Microsoft sender). " +
							"Any other value must be a verified domain of the tenant, written exactly as it appears in the tenant domain list; " +
							"Graph compares it case-sensitively and rejects other values. When not configured, the current tenant value is kept.",
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"use_company_branding": schema.BoolAttribute{
						MarkdownDescription: "Whether the tenant's company branding is used in email notifications. When not configured, the current tenant value is kept.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
