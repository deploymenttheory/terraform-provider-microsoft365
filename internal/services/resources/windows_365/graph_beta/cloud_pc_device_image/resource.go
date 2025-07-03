package graphBetaCloudPcDeviceImage

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	ResourceName  = "graph_beta_windows_365_cloud_pc_device_image"
	CreateTimeout = 180
	UpdateTimeout = 180
	ReadTimeout   = 180
	DeleteTimeout = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &CloudPcDeviceImageResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &CloudPcDeviceImageResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &CloudPcDeviceImageResource{}

	// Enables plan modification/diff suppression
	_ resource.ResourceWithModifyPlan = &CloudPcDeviceImageResource{}
)

func NewCloudPcDeviceImageResource() resource.Resource {
	return &CloudPcDeviceImageResource{
		ReadPermissions: []string{
			"CloudPC.Read.All",
		},
		WritePermissions: []string{
			"CloudPC.ReadWrite.All",
		},
	}
}

type CloudPcDeviceImageResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ProviderTypeName string
	TypeName         string
	ReadPermissions  []string
	WritePermissions []string
}

func (r *CloudPcDeviceImageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	r.TypeName = ResourceName
	resp.TypeName = r.FullTypeName()
}

func (r *CloudPcDeviceImageResource) FullTypeName() string {
	return r.ProviderTypeName + "_" + r.TypeName
}

func (r *CloudPcDeviceImageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, constants.PROVIDER_NAME+"_"+ResourceName)
}

func (r *CloudPcDeviceImageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *CloudPcDeviceImageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Windows 365 Cloud PC Device Image using the Microsoft Graph Beta API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier (ID) of the image resource on the Cloud PC. Read-only.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"display_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The display name of this image.",
			},
			"version": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The image version. For example, 0.0.1 and 1.5.13.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 32),
				},
			},
			"source_image_resource_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the source image resource on Azure. The required ID format is: \"/subscriptions/{subscription-id}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/images/{imageName}\".",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.AzureImageResourceIDRegex), "Must be a valid Azure image resource ID: /subscriptions/{subscription-id}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/images/{imageName}"),
				},
			},
			"error_code": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The error code of the status of the image that indicates why the upload failed, if applicable. Possible values are: internalServerError, sourceImageNotFound, osVersionNotSupported, sourceImageInvalid, sourceImageNotGeneralized, unknownFutureValue, vmAlreadyAzureAdJoined, paidSourceImageNotSupport, sourceImageNotSupportCustomizeVMName, sourceImageSizeExceedsLimitation, sourceImageWithDataDiskNotSupported, sourceImageWithDiskEncryptionSetNotSupported. Read-only.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"internalServerError",
						"sourceImageNotFound",
						"osVersionNotSupported",
						"sourceImageInvalid",
						"sourceImageNotGeneralized",
						"unknownFutureValue",
						"vmAlreadyAzureAdJoined",
						"paidSourceImageNotSupport",
						"sourceImageNotSupportCustomizeVMName",
						"sourceImageSizeExceedsLimitation",
						"sourceImageWithDataDiskNotSupported",
						"sourceImageWithDiskEncryptionSetNotSupported",
					),
				},
			},
			"expiration_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The date when the image became unavailable. Read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The data and time when the image was last modified. The timestamp represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.",
			},
			"operating_system": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The operating system of the image. For example, Windows 10 Enterprise. Read-only.",
			},
			"os_build_number": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The OS build version of the image. For example, 1909. Read-only.",
			},
			"os_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The OS status of this image. Possible values are: supported, supportedWithWarning, unknown, unknownFutureValue. The default value is unknown. Read-only.",
				Validators: []validator.String{
					stringvalidator.OneOf("supported", "supportedWithWarning", "unknown", "unknownFutureValue"),
				},
			},
			"os_version_number": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The operating system version of this image. For example, 10.0.22000.296. Read-only.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the image on the Cloud PC. Possible values are: pending, ready, warning, failed, unknownFutureValue. Read-only.",
				Validators: []validator.String{
					stringvalidator.OneOf("pending", "ready", "warning", "failed", "unknownFutureValue"),
				},
			},
		},
	}
}
