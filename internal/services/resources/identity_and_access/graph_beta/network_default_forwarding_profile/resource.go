package graphBetaNetworkDefaultForwardingProfile

import (
	"context"
	"maps"
	"slices"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	profileschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema/graph_beta/identity_and_access"
)

const (
	ResourceName  = "microsoft365_graph_beta_identity_and_access_network_default_forwarding_profile"
	CreateTimeout = 180
	ReadTimeout   = 180
	UpdateTimeout = 180
	DeleteTimeout = 180
)

var trafficTypes = []string{"private", "internet", "m365"}

var (
	_ resource.Resource                   = &NetworkDefaultForwardingProfileResource{}
	_ resource.ResourceWithConfigure      = &NetworkDefaultForwardingProfileResource{}
	_ resource.ResourceWithImportState    = &NetworkDefaultForwardingProfileResource{}
	_ resource.ResourceWithValidateConfig = &NetworkDefaultForwardingProfileResource{}
)

// NetworkDefaultForwardingProfileResource manages a built-in forwarding profile.
type NetworkDefaultForwardingProfileResource struct {
	client           *msgraphbetasdk.GraphServiceClient
	ReadPermissions  []string
	WritePermissions []string
	ResourcePath     string
}

func NewNetworkDefaultForwardingProfileResource() resource.Resource {
	return &NetworkDefaultForwardingProfileResource{
		ReadPermissions:  []string{"NetworkAccess.ReadWrite.All"},
		WritePermissions: []string{"NetworkAccess.ReadWrite.All"},
		ResourcePath:     "/networkaccess/forwardingProfiles",
	}
}

func (r *NetworkDefaultForwardingProfileResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = ResourceName
}

func (r *NetworkDefaultForwardingProfileResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = client.SetGraphBetaClientForResource(ctx, req, resp, ResourceName)
}

// ImportState accepts a traffic type or a profile ID. Read resolves the rest.
func (r *NetworkDefaultForwardingProfileResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	if slices.Contains(trafficTypes, req.ID) {
		resp.Diagnostics.Append(
			resp.State.SetAttribute(ctx, path.Root("traffic_forwarding_type"), req.ID)...)
		return
	}
	if _, err := uuid.Parse(req.ID); err != nil {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Use a traffic type (private, internet or m365) or the default profile ID (GUID).",
		)
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *NetworkDefaultForwardingProfileResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config NetworkDefaultForwardingProfileResourceModel
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

func (r *NetworkDefaultForwardingProfileResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	computedString := func(description string) schema.StringAttribute {
		return schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: description,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
		}
	}
	attributes := map[string]schema.Attribute{
		"id": computedString("Default forwarding profile ID. It differs per tenant."),
		"traffic_forwarding_type": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Traffic type of the built-in profile to manage. Values: `private`, `internet`, `m365`. Changing it moves management to another default profile.",
			Validators:          []validator.String{stringvalidator.OneOf(trafficTypes...)},
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
		},
		"name": computedString(
			"Profile name. Graph does not allow changing it on default profiles.",
		),
		"description": computedString(
			"Profile description. Graph does not allow changing it on default profiles.",
		),
		"priority": schema.Int32Attribute{
			Computed:            true,
			MarkdownDescription: "Profile priority. Graph does not allow changing it on default profiles.",
			PlanModifiers:       []planmodifier.Int32{int32planmodifier.UseStateForUnknown()},
		},
		"version":  computedString("Profile version."),
		"timeouts": commonschema.ResourceTimeouts(ctx),
	}
	maps.Copy(attributes, profileschema.ForwardingProfileSettingsAttributes())
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a built-in Microsoft Entra Global Secure Access traffic forwarding profile (traffic forwarding profile v2) using Microsoft Graph beta `/networkaccess/forwardingProfiles`. " +
			"The profile is selected by `traffic_forwarding_type`. " +
			"Initial apply updates the existing default profile; destroy removes Terraform state only and leaves the remote profile unchanged. " +
			"Manage one instance per traffic type per tenant. Use `microsoft365_graph_beta_identity_and_access_network_forwarding_profile` for custom profiles. " +
			"Policy links and remote network associations are not managed by this resource.",
		Attributes: attributes,
	}
}
