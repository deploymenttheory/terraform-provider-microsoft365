package graphCloudPcDeviceImage

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

const (
	ResourceName  = "microsoft365_graph_device_and_app_management_cloud_pc_device_image"
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
	client           *msgraphsdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
}

// Metadata returns the resource type name.
func (r *CloudPcDeviceImageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *CloudPcDeviceImageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = client.SetGraphStableClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *CloudPcDeviceImageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Schema returns the schema for the resource.
func (r *CloudPcDeviceImageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier (ID) of the image resource on the Cloud PC. Read-only.",
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"display_name": schema.StringAttribute{
				Required: true,
				Description: "The display name of the associated device image. " +
					"The device image display name and the version are used to " +
					"uniquely identify the Cloud PC device image.",
			},
			"error_code": schema.StringAttribute{
				Computed: true,
				Description: "The error code of the status of the image that indicates why " +
					"the upload failed, if applicable. Possible values are: " +
					"internalServerError, sourceImageNotFound, osVersionNotSupported, " +
					"sourceImageInvalid, sourceImageNotGeneralized, unknownFutureValue, " +
					"vmAlreadyAzureAdJoined, paidSourceImageNotSupport, " +
					"sourceImageNotSupportCustomizeVMName, sourceImageSizeExceedsLimitation. " +
					"Note that you must use the Prefer: include-unknown-enum-members " +
					"request header to get the following values from this evolvable enum: " +
					"vmAlreadyAzureAdJoined, paidSourceImageNotSupport, " +
					"sourceImageNotSupportCustomizeVMName, sourceImageSizeExceedsLimitation. Read-only.",
			},
			"expiration_date": schema.StringAttribute{
				Computed:    true,
				Description: "The date when the image became unavailable. Read-only.",
			},
			"last_modified_date_time": schema.StringAttribute{
				Computed: true,
				Description: "The date and time when the image was last modified. " +
					"The timestamp represents date and time information using ISO 8601 format " +
					"and is always in UTC. For example, midnight UTC on Jan 1, 2014 " +
					"is 2014-01-01T00:00:00Z. Read-only.",
			},
			"operating_system": schema.StringAttribute{
				Computed:    true,
				Description: "The operating system (OS) of the image. For example, Windows 10 Enterprise. Read-only.",
			},
			"os_build_number": schema.StringAttribute{
				Computed:    true,
				Description: "The OS build version of the image. For example, 1909. Read-only.",
			},
			"os_status": schema.StringAttribute{
				Computed: true,
				Description: "The OS status of this image. Possible values are: supported, supportedWithWarning, " +
					"unknown, unknownFutureValue. The default value is unknown. Read-only.",
				Validators: []validator.String{
					stringvalidator.OneOf("supported", "supportedWithWarning", "unknown", "unknownFutureValue"),
				},
			},
			"source_image_resource_id": schema.StringAttribute{
				Required: true,
				Description: "The unique identifier (ID) of the source image resource on Azure. " +
					"The required ID format is: /subscriptions/{subscription-id}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/images/{imageName}.",
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "The status of the image on the Cloud PC. Possible values are: pending, ready, failed, " +
					"unknownFutureValue. Read-only.",
				Validators: []validator.String{
					stringvalidator.OneOf("pending", "ready", "failed", "unknownFutureValue"),
				},
			},
			"version": schema.StringAttribute{
				Required:    true,
				Description: "The image version. For example, 0.0.1 and 1.5.13.",
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
