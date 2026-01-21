package graphBetaGroupPolicyUploadedDefinitionFiles

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	regexvalidator "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validate/attribute"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_beta_device_management_group_policy_uploaded_definition_files"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 300 // Deletion typically takes 5 minutes to complete.
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &GroupPolicyUploadedDefinitionFileResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &GroupPolicyUploadedDefinitionFileResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &GroupPolicyUploadedDefinitionFileResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &GroupPolicyUploadedDefinitionFileResource{}
)

func NewGroupPolicyUploadedDefinitionFileResource() resource.Resource {
	return &GroupPolicyUploadedDefinitionFileResource{
		ReadPermissions: []string{
			"DeviceManagementConfiguration.Read.All",
		},
		WritePermissions: []string{
			"DeviceManagementConfiguration.ReadWrite.All",
		},
		ResourcePath: "/deviceManagement/groupPolicyUploadedDefinitionFiles",
	}
}

type GroupPolicyUploadedDefinitionFileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// Metadata returns the resource type name.
func (r *GroupPolicyUploadedDefinitionFileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *GroupPolicyUploadedDefinitionFileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState returns an error as this resource does not support import.
// ADMX/ADML file content is not returned by the API, making complete import impossible.
func (r *GroupPolicyUploadedDefinitionFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError(
		"Import Not Supported",
		"This resource does not support state import. The ADMX/ADML file content is not returned by the Microsoft Graph API, "+
			"which makes it impossible to import the complete resource state. "+
			"You must manage this resource by creating it directly in your Terraform configuration.",
	)
}

// Schema returns the schema for the resource.
func (r *GroupPolicyUploadedDefinitionFileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages group policy uploaded definition files using the `/deviceManagement/groupPolicyUploadedDefinitionFiles` endpoint. This resource is used to group policy uploaded definition files are ADMX files that define group policies that can be deployed to devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
				MarkdownDescription: "The ID of the group policy uploaded definition file.",
			},
			"display_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The display name of the group policy uploaded definition file.",
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
				Required:            true,
				MarkdownDescription: "The file name of the group policy uploaded definition file.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The content of the group policy uploaded definition file. Request is sent as raw bytes. This is a write-only field and will not be stored in state.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"default_language_code": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The default language code of the group policy uploaded definition file. Must be in the format 'xx-YY' (e.g., 'en-US', 'fr-FR').",
				Validators: []validator.String{
					regexvalidator.RegexMatches(
						regexp.MustCompile(constants.LocaleRegex),
						"must be a valid locale code in the format 'xx-YY' where xx is a 2-letter lowercase language code and YY is a 2-letter uppercase country code (e.g., 'en-US', 'fr-FR', 'de-DE')",
					),
				},
			},
			"language_codes": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "The language codes supported by the group policy uploaded definition file.",
			},
			"target_prefix": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The target prefix of the group policy uploaded definition file.",
			},
			"target_namespace": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The target namespace of the group policy uploaded definition file.",
			},
			"policy_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The policy type of the group policy uploaded definition file.",
			},
			"revision": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The revision of the group policy uploaded definition file.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the group policy uploaded definition file. Possible values are: uploadInProgress, available, uploadFailed.",
			},
			"upload_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the group policy uploaded definition file was uploaded.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date and time when the group policy uploaded definition file was last modified.",
			},
			"group_policy_uploaded_language_files": schema.SetNestedAttribute{
				Required:            true,
				MarkdownDescription: "The language file(s) associated with the group policy uploaded definition file.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"file_name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The file name of the group policy uploaded language file.",
						},
						"language_code": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The language code of the group policy uploaded language file. Must be in the format 'xx-YY' (e.g., 'en-US', 'fr-FR').",
							Validators: []validator.String{
								regexvalidator.RegexMatches(
									regexp.MustCompile(constants.LocaleRegex),
									"must be a valid locale code in the format 'xx-YY' where xx is a 2-letter lowercase language code and YY is a 2-letter uppercase country code (e.g., 'en-US', 'fr-FR', 'de-DE')",
								),
							},
						},
						"content": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The content of the group policy uploaded language file. Request is sent as raw bytes. This is a write-only field and will not be stored in state.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.RequiresReplace(),
							},
						},
					},
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
