package graphBetaMacosDeviceComplianceScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_macos_device_compliance_script"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                = &DeviceComplianceScriptResource{}
	_ resource.ResourceWithConfigure   = &DeviceComplianceScriptResource{}
	_ resource.ResourceWithImportState = &DeviceComplianceScriptResource{}
	_ resource.ResourceWithIdentity    = &DeviceComplianceScriptResource{}
)

func NewDeviceComplianceScriptResource() resource.Resource {
	return &DeviceComplianceScriptResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/deviceComplianceScripts",
	}
}

type DeviceComplianceScriptResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func (r *DeviceComplianceScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

func (r *DeviceComplianceScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *DeviceComplianceScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *DeviceComplianceScriptResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (r *DeviceComplianceScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages macOS device compliance scripts in Microsoft Intune using the `/deviceManagement/deviceComplianceScripts` endpoint " +
			"(the same endpoint the Windows device compliance script resource uses -- the object is distinguished internally by its `platform` property, " +
			"which this resource always sets to `macOS`). This resource uploads a POSIX-compliant bash discovery script that Intune runs on enrolled " +
			"macOS devices; pair it with a `microsoft365_graph_beta_device_management_macos_device_compliance_policy` resource's " +
			"`device_compliance_policy_script` block to build a custom compliance policy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "Unique identifier for the device compliance script.",
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the device compliance script.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional description of the resource. Maximum length is 1500 characters.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1500),
				},
			},
			"publisher": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the device compliance script publisher.",
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Version of the device compliance script.",
			},
			"run_as_account": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Indicates the type of execution context. Possible values are: system, user.",
				Validators: []validator.String{
					stringvalidator.OneOf("system", "user"),
				},
			},
			"enforce_signature_check": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				MarkdownDescription: "Indicate whether the script signature needs to be checked. Default is false.",
			},
			"detection_script_content": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The entire content of the detection script. Must be a POSIX-compliant bash script " +
					"(`#!/bin/bash` shebang, UTF-8 encoded with no BOM, targeting bash 3.2 -- that's what ships on macOS, " +
					"so no associative arrays or `${var,,}`), must emit a single-line JSON object on stdout keyed on your " +
					"setting names, and must exit 0 on success / non-zero on failure. Max 1 MB script content, max 1 MB output, " +
					"must finish in under 10 minutes.",
			},
			"role_scope_tag_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Set of scope tag IDs for this device compliance script.",
				PlanModifiers: []planmodifier.Set{
					planmodifiers.DefaultSetValue(
						[]attr.Value{types.StringValue("0")},
					),
				},
			},
			"created_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp of when the device compliance script was created. This property is read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp of when the device compliance script was modified. This property is read-only.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}