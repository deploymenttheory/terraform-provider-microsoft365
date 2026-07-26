package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	commonschema "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type onPremisesIpApplicationSegmentResourceModelV0 struct {
	ID                  types.String   `tfsdk:"id"`
	ApplicationObjectID types.String   `tfsdk:"application_object_id"`
	DestinationHost     types.String   `tfsdk:"destination_host"`
	DestinationType     types.String   `tfsdk:"destination_type"`
	Ports               types.Set      `tfsdk:"ports"`
	Protocol            types.String   `tfsdk:"protocol"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

// UpgradeState returns state migrations for prior IP application segment schemas.
func (r *OnPremisesIpApplicationSegmentResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   onPremisesIpApplicationSegmentSchemaV0(ctx),
			StateUpgrader: upgradeOnPremisesIpApplicationSegmentStateV0toV1,
		},
	}
}

func upgradeOnPremisesIpApplicationSegmentStateV0toV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var priorState onPremisesIpApplicationSegmentResourceModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	protocol := types.SetNull(types.StringType)
	switch {
	case priorState.Protocol.IsUnknown():
		protocol = types.SetUnknown(types.StringType)
	case !priorState.Protocol.IsNull():
		var diags diag.Diagnostics
		protocol, diags = types.SetValueFrom(ctx, types.StringType, terraformProtocols(priorState.Protocol.ValueString()))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	upgradedState := OnPremisesIpApplicationSegmentResourceModel{
		ID:                  priorState.ID,
		ApplicationObjectID: priorState.ApplicationObjectID,
		DestinationHost:     priorState.DestinationHost,
		DestinationType:     priorState.DestinationType,
		Ports:               priorState.Ports,
		Protocol:            protocol,
		Timeouts:            priorState.Timeouts,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &upgradedState)...)
}

// onPremisesIpApplicationSegmentSchemaV0 is the immutable schema used before
// protocol changed from string to set(string).
func onPremisesIpApplicationSegmentSchemaV0(ctx context.Context) *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "Manages an IP application segment for on-premises publishing. " +
			"IP application segments define the destination hosts, ports, and protocols for applications published through Azure AD Application Proxy.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the application segment.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_object_id": schema.StringAttribute{
				MarkdownDescription: "The unique object identifier of the application.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination_host": schema.StringAttribute{
				MarkdownDescription: "Either the IP address, IP range, or FQDN of the application segment, with or without wildcards.",
				Required:            true,
			},
			"destination_type": schema.StringAttribute{
				MarkdownDescription: "The type of destination for the application segment." +
					"The supported values are: `ipAddress`, `ipRangeCidr`, and `fqdn`. " +
					"Microsoft Learn lists additional enum members for `ipApplicationSegment`, but this application-scoped Graph endpoint currently rejects `dnsSuffix` for nonweb applications and does not create a usable address range for `ipRange`.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ipAddress", "ipRangeCidr", "fqdn"),
				},
			},
			"ports": schema.SetAttribute{
				MarkdownDescription: "List of ports supported for the application segment.",
				Required:            true,
				ElementType:         types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(constants.PortRangeRegex), "Each port defined in the set must be a valid format (xxxx-xxxx) e.g 80-80, 443-443, 8080-8080, 8443-8443"),
					),
				},
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "Indicates the protocol of the network traffic acquired for the application segment." +
					"The possible values are: `tcp`, `udp`, `unknownFutureValue`.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp"),
				},
			},
			"timeouts": commonschema.ResourceTimeouts(ctx),
		},
	}
}
