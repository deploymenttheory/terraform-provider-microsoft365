package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	commonschemagraphbeta "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/device_management"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_windows_bios_configurations_and_other_settings_template"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &WindowsBiosConfigurationsAndOtherSettingsTemplateResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &WindowsBiosConfigurationsAndOtherSettingsTemplateResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &WindowsBiosConfigurationsAndOtherSettingsTemplateResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &WindowsBiosConfigurationsAndOtherSettingsTemplateResource{}

	// Enables identity schema for list resource support
	_ resource.ResourceWithIdentity = &WindowsBiosConfigurationsAndOtherSettingsTemplateResource{}
)

func NewWindowsBiosConfigurationsAndOtherSettingsTemplateResource() resource.Resource {
	return &WindowsBiosConfigurationsAndOtherSettingsTemplateResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/hardwareConfigurations",
	}
}

type WindowsBiosConfigurationsAndOtherSettingsTemplateResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *WindowsBiosConfigurationsAndOtherSettingsTemplateResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *WindowsBiosConfigurationsAndOtherSettingsTemplateResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *WindowsBiosConfigurationsAndOtherSettingsTemplateResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// IdentitySchema defines the identity schema for this resource, used by list operations to uniquely identify instances
func (r *WindowsBiosConfigurationsAndOtherSettingsTemplateResource) IdentitySchema(
	ctx context.Context,
	req resource.IdentitySchemaRequest,
	resp *resource.IdentitySchemaResponse,
) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *WindowsBiosConfigurationsAndOtherSettingsTemplateResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the Intune 'BIOS configuration and other settings' template using the `/deviceManagement/hardwareConfigurations` endpoint. " +
			"This template applies an OEM generated hardware/BIOS configuration file (for example a Dell Command | Configure `.cctk` package) to enrolled Windows 10/11 " +
			"Microsoft Entra joined devices, allowing administrators to remotely control device hardware properties such as Secure Boot or TPM state. Dell is the only " +
			"OEM supported by the service at this time.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for this Intune BIOS configuration and other settings template.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the BIOS configuration and other settings template.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"file_name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Name of the OEM configuration file being uploaded, including its extension (for example `bios-config.cctk`). " +
					"The service infers `hardware_configuration_format` from this file name when the format is not supplied explicitly.",
			},
			"configuration_file_content": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
				MarkdownDescription: "The base64 encoded content of the OEM generated configuration file. Use Terraform's `filebase64()` function to read the file " +
					"from disk, for example `filebase64(\"${path.module}/bios-config.cctk\")`. Note that a Dell Command | Configure package may contain BIOS setup " +
					"passwords, so this attribute is treated as sensitive.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.Base64Regex),
						"must be a base64 encoded string, for example the output of filebase64()",
					),
				},
			},
			"hardware_configuration_format": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The OEM the configuration file targets. Possible values are: `dell`, `surface`, `surfaceDock`. When omitted, the service infers the format from `file_name`.",
				Validators: []validator.String{
					stringvalidator.OneOf("dell", "surface", "surfaceDock"),
				},
			},
			"per_device_password_disabled": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When `true`, Intune does not generate and manage a unique BIOS password for each targeted device. Leave this `false` to let Intune manage per device BIOS passwords.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this BIOS configuration and other settings template.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"version": schema.Int32Attribute{
				Computed:            true,
				MarkdownDescription: "The version of the BIOS configuration and other settings template. Read-only.",
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the BIOS configuration and other settings template was created. Read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time the BIOS configuration and other settings template was last modified. Read-only.",
			},
			"assignments": commonschemagraphbeta.HardwareConfigurationAssignmentsSchema(),
			"timeouts":    commonschema.ResourceTimeouts(ctx),
		},
	}
}
