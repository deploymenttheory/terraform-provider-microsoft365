package graphBetaNetworkForwardingProfile

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	profileschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/identity_and_access"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_forwarding_profile"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var (
	_ resource.Resource                   = &NetworkForwardingProfileResource{}
	_ resource.ResourceWithConfigure      = &NetworkForwardingProfileResource{}
	_ resource.ResourceWithImportState    = &NetworkForwardingProfileResource{}
	_ resource.ResourceWithValidateConfig = &NetworkForwardingProfileResource{}
)

// NetworkForwardingProfileResource manages a custom traffic forwarding profile.
type NetworkForwardingProfileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func NewNetworkForwardingProfileResource() resource.Resource {
	return &NetworkForwardingProfileResource{
		ReadPermissions:  []string{"NetworkAccess.ReadWrite.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/forwardingProfiles",
	}
}

func (r *NetworkForwardingProfileResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkForwardingProfileResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

func (r *NetworkForwardingProfileResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if _, err := uuid.Parse(req.ID); err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Use the forwarding profile ID (GUID). Import default profiles with microsoft365_graph_beta_identity_and_access_network_default_forwarding_profile.",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkForwardingProfileResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config NetworkForwardingProfileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(profileschema.ValidateForwardingProfileSettings(
		ctx,
		config.TrafficForwardingType,
		config.LocalNetworkSettings,
		config.AssignmentRules,
	)...)
}

func (r *NetworkForwardingProfileResource) Schema(
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
			Required:            true,
			MarkdownDescription: "Traffic type. Values: `private`, `internet`, `m365`. Graph currently allows custom profiles for `private` only; `internet` and `m365` fail unless Microsoft enables multiple profiles for the tenant. Changing it replaces the profile.",
			Validators: []validator.String{
				stringvalidator.OneOf("private", "internet", "m365"),
			},
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
		},
		"priority": schema.Int32Attribute{
			Required:            true,
			MarkdownDescription: "Profile priority; lower numbers win when several profiles apply. Must be unique per tenant. Graph accepts 101-200 for `private` custom profiles (the default private profile uses 200).",
			Validators:          []validator.Int32{int32validator.AtLeast(1)},
		},
		"version": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Profile version.",
			PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		},
		"timeouts": commonschema.ResourceTimeouts(ctx),
	}
	maps.Copy(attributes, profileschema.ForwardingProfileSettingsAttributes())
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a custom Microsoft Entra Global Secure Access traffic forwarding profile (traffic forwarding profile v2) using Microsoft Graph beta `/networkaccess/forwardingProfiles`. " +
			"The tenant must use traffic forwarding profiles v2. Destroy deletes the profile. " +
			"Use `microsoft365_graph_beta_identity_and_access_network_default_forwarding_profile` for the built-in profiles. " +
			"Policy links and remote network associations are not managed by this resource.",
		Attributes: attributes,
	}
}
