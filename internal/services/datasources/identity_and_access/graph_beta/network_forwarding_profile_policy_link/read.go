package graphBetaNetworkForwardingProfilePolicyLink

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	profileds "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/identity_and_access/graph_beta/network_forwarding_profile"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/networkaccess"
)

type lookupMethod int

const (
	lookupByForwardingProfileID lookupMethod = iota
	lookupByForwardingProfileName
	lookupByTrafficForwardingType
	lookupUnset
)

func (d *NetworkForwardingProfilePolicyLinkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object NetworkForwardingProfilePolicyLinkDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	method := determineLookupMethod(object)
	if method == lookupUnset {
		resp.Diagnostics.AddError(
			"Missing Forwarding Profile Query Criteria",
			"One of forwarding_profile_id, forwarding_profile_name, or traffic_forwarding_type must be specified.",
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	expand := []string{"policies($expand=policy)"}
	profiles := make([]profileds.ForwardingProfileModel, 0)

	if method == lookupByForwardingProfileID {
		profile, err := d.client.NetworkAccess().ForwardingProfiles().ByForwardingProfileId(object.ForwardingProfileID.ValueString()).Get(ctx, &networkaccess.ForwardingProfilesForwardingProfileItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &networkaccess.ForwardingProfilesForwardingProfileItemRequestBuilderGetQueryParameters{Expand: expand},
		})
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}
		profiles = append(profiles, profileds.MapRemoteStateToDataSource(profile))
	} else {
		response, err := d.client.NetworkAccess().ForwardingProfiles().Get(ctx, &networkaccess.ForwardingProfilesRequestBuilderGetRequestConfiguration{
			QueryParameters: &networkaccess.ForwardingProfilesRequestBuilderGetQueryParameters{Expand: expand},
		})
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}
		for _, profile := range response.GetValue() {
			mapped := profileds.MapRemoteStateToDataSource(profile)
			if matchesForwardingProfile(mapped, object, method) {
				profiles = append(profiles, mapped)
			}
		}
	}

	matches := findPolicyLinkMatches(profiles, object.PolicyName.ValueString())
	if len(matches) == 0 {
		resp.Diagnostics.AddError(
			"Forwarding Profile Policy Link Not Found",
			fmt.Sprintf("No forwarding profile policy link matched policy_name %q for the selected forwarding profile criteria.", object.PolicyName.ValueString()),
		)
		return
	}
	if len(matches) > 1 {
		resp.Diagnostics.AddError(
			"Multiple Forwarding Profile Policy Links Found",
			fmt.Sprintf("Found %d forwarding profile policy links matching policy_name %q. Use a more specific forwarding profile selector.", len(matches), object.PolicyName.ValueString()),
		)
		return
	}

	timeouts := object.Timeouts
	object = matches[0].toModel()
	object.Timeouts = timeouts
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished datasource read for %s, policy link id %s", DataSourceName, object.PolicyLinkID.ValueString()))
}

func determineLookupMethod(object NetworkForwardingProfilePolicyLinkDataSourceModel) lookupMethod {
	if hasStringValue(object.ForwardingProfileID) {
		return lookupByForwardingProfileID
	}
	if hasStringValue(object.ForwardingProfileName) {
		return lookupByForwardingProfileName
	}
	if hasStringValue(object.TrafficForwardingType) {
		return lookupByTrafficForwardingType
	}
	return lookupUnset
}

func hasStringValue(value interface {
	IsNull() bool
	IsUnknown() bool
	ValueString() string
}) bool {
	return !value.IsNull() && !value.IsUnknown() && strings.TrimSpace(value.ValueString()) != ""
}

func matchesForwardingProfile(profile profileds.ForwardingProfileModel, object NetworkForwardingProfilePolicyLinkDataSourceModel, method lookupMethod) bool {
	switch method {
	case lookupByForwardingProfileName:
		return strings.EqualFold(profile.Name.ValueString(), object.ForwardingProfileName.ValueString())
	case lookupByTrafficForwardingType:
		return strings.EqualFold(profile.TrafficForwardingType.ValueString(), object.TrafficForwardingType.ValueString())
	default:
		return true
	}
}

type policyLinkMatch struct {
	profile profileds.ForwardingProfileModel
	link    profileds.ForwardingProfilePolicyLinkModel
}

func findPolicyLinkMatches(profiles []profileds.ForwardingProfileModel, policyName string) []policyLinkMatch {
	matches := make([]policyLinkMatch, 0)
	for _, profile := range profiles {
		for _, link := range profile.Policies {
			if strings.EqualFold(link.PolicyName.ValueString(), policyName) {
				matches = append(matches, policyLinkMatch{profile: profile, link: link})
			}
		}
	}
	return matches
}
