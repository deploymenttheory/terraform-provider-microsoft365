package graphBetaNetworkDeviceForwardingProfile

import (
	"context"
	"maps"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	profileschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/identity_and_access"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

const (
	ResourceName = "microsoft365_graph_beta_identity_and_access_network_device_forwarding_profile"

	trafficForwardingTypePrivate = "private"
	CreateTimeout                = 180
	ReadTimeout                  = 180
	UpdateTimeout                = 180
	DeleteTimeout                = 180
)

var (
	// Basic resource interface (CRUD operations)
	_ resource.Resource = &NetworkDeviceForwardingProfileResource{}

	// Allows the resource to be configured with the provider client
	_ resource.ResourceWithConfigure = &NetworkDeviceForwardingProfileResource{}

	// Enables import functionality
	_ resource.ResourceWithImportState = &NetworkDeviceForwardingProfileResource{}

	// Enables resource-specific configuration validation
	_ resource.ResourceWithValidateConfig = &NetworkDeviceForwardingProfileResource{}
)

// NetworkDeviceForwardingProfileResource manages a pre-logon (device) forwarding profile.
type NetworkDeviceForwardingProfileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

// NewNetworkDeviceForwardingProfileResource creates the network device forwarding profile resource.
func NewNetworkDeviceForwardingProfileResource() resource.Resource {
	return &NetworkDeviceForwardingProfileResource{
		ReadPermissions:  []string{"NetworkAccess.ReadWrite.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/deviceForwardingProfiles",
	}
}

// Metadata returns the resource type name.
func (r *NetworkDeviceForwardingProfileResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

// Configure sets the client for the resource.
func (r *NetworkDeviceForwardingProfileResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState imports the resource state.
func (r *NetworkDeviceForwardingProfileResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if _, err := uuid.Parse(req.ID); err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Use the device forwarding profile ID (GUID).",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// ValidateConfig validates the resource configuration.
func (r *NetworkDeviceForwardingProfileResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config NetworkDeviceForwardingProfileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(profileschema.ValidateForwardingProfileSettings(
		ctx,
		types.StringValue(trafficForwardingTypePrivate),
		config.LocalNetworkSettings,
		config.AssignmentRules,
	)...)
	// Graph refuses user assignments on device profiles ("Only device and group assignments are allowed").
	if config.AssignmentRules.IsNull() || config.AssignmentRules.IsUnknown() {
		return
	}
	var rules sharedmodels.AssignmentRulesModel
	if d := config.AssignmentRules.As(ctx, &rules, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true}); d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}
	if !rules.UserIDs.IsNull() && !rules.UserIDs.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("assignment_rules").AtName("user_ids"),
			"User assignments are not supported",
			"Pre-logon (device) forwarding profiles can be assigned to devices and groups only.",
		)
	}
}

// Schema defines the schema for the resource.
func (r *NetworkDeviceForwardingProfileResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	attributes := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Forwarding profile ID.",
			PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Profile name.",
			Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Profile description.",
		},
		"traffic_forwarding_type": schema.StringAttribute{
			Computed:            true,
			Default:             stringdefault.StaticString(trafficForwardingTypePrivate),
			MarkdownDescription: "Traffic type. Always `private`; Graph supports pre-logon profiles for Private Access only.",
		},
		"priority": schema.Int32Attribute{
			Required:            true,
			MarkdownDescription: "Profile priority, 101-200; lower wins. Graph accepts other values on update, but they would fail on re-create.",
			Validators:          []validator.Int32{int32validator.Between(101, 200)},
		},
		"version": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Profile version.",
			PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"timeouts": commonschema.ResourceTimeouts(ctx),
	}
	maps.Copy(attributes, profileschema.ForwardingProfileSettingsAttributes())
	// Device profiles reject source matching, so there is no attribute.
	delete(attributes, "forwarding_source_matching")
	// The shared assignment descriptions mention users, so rewrite them for device profiles.
	if rules, ok := attributes["assignment_rules"].(schema.SingleNestedAttribute); ok {
		rules.MarkdownDescription = "Group, device and platform assignments. Graph replaces assignments as a whole. " +
			"Omitting it leaves the current assignments unchanged."
		nested := maps.Clone(rules.Attributes)
		if m, ok := nested["membership"].(schema.StringAttribute); ok {
			m.MarkdownDescription = "`all` assigns all devices, `enumerated` assigns the listed IDs, `none` assigns nothing."
			nested["membership"] = m
		}
		if u, ok := nested["user_ids"].(schema.SetAttribute); ok {
			u.MarkdownDescription = "Not supported on pre-logon profiles; Graph rejects user assignments."
			nested["user_ids"] = u
		}
		rules.Attributes = nested
		attributes["assignment_rules"] = rules
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Microsoft Entra Global Secure Access pre-logon (device) forwarding profile using Microsoft Graph beta `/networkaccess/deviceForwardingProfiles`. " +
			"Pre-logon profiles connect a device to Private Access before a user signs in to Windows. They support Private Access only and are assigned to devices or groups, not users. " +
			"Destroy deletes the profile. Policy links are not managed by this resource.",
		Attributes: attributes,
	}
}
