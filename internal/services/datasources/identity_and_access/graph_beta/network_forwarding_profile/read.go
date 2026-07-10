package graphBetaNetworkForwardingProfile

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/networkaccess"
)

type lookupMethod int

const (
	lookupByForwardingProfileID lookupMethod = iota
	lookupByName
	lookupByTrafficForwardingType
	lookupListAll
	lookupUnset
)

func (d *NetworkForwardingProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object NetworkForwardingProfileDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	method := determineLookupMethod(object)
	if method == lookupUnset {
		resp.Diagnostics.AddError(
			"Missing Query Criteria",
			"One of forwarding_profile_id, name, traffic_forwarding_type, or list_all must be specified.",
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	expand := []string{"policies($expand=policy)"}
	items := make([]ForwardingProfileModel, 0)

	if method == lookupByForwardingProfileID {
		profile, err := d.client.NetworkAccess().ForwardingProfiles().ByForwardingProfileId(object.ForwardingProfileID.ValueString()).Get(ctx, &networkaccess.ForwardingProfilesForwardingProfileItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &networkaccess.ForwardingProfilesForwardingProfileItemRequestBuilderGetQueryParameters{Expand: expand},
		})
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}
		items = append(items, MapRemoteStateToDataSource(profile))
	} else {
		profiles, err := d.client.NetworkAccess().ForwardingProfiles().Get(ctx, &networkaccess.ForwardingProfilesRequestBuilderGetRequestConfiguration{
			QueryParameters: &networkaccess.ForwardingProfilesRequestBuilderGetQueryParameters{Expand: expand},
		})
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
			return
		}
		for _, profile := range profiles.GetValue() {
			mapped := MapRemoteStateToDataSource(profile)
			switch method {
			case lookupListAll:
				items = append(items, mapped)
			case lookupByName:
				if strings.EqualFold(mapped.Name.ValueString(), object.Name.ValueString()) {
					items = append(items, mapped)
				}
			case lookupByTrafficForwardingType:
				if strings.EqualFold(mapped.TrafficForwardingType.ValueString(), object.TrafficForwardingType.ValueString()) {
					items = append(items, mapped)
				}
			}
		}
	}

	object.ID = types.StringValue(dataSourceID(object, method))
	object.Items = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished datasource read for %s, found %d items", DataSourceName, len(items)))
}

func determineLookupMethod(object NetworkForwardingProfileDataSourceModel) lookupMethod {
	if hasStringValue(object.ForwardingProfileID) {
		return lookupByForwardingProfileID
	}
	if hasStringValue(object.Name) {
		return lookupByName
	}
	if hasStringValue(object.TrafficForwardingType) {
		return lookupByTrafficForwardingType
	}
	if !object.ListAll.IsNull() && !object.ListAll.IsUnknown() && object.ListAll.ValueBool() {
		return lookupListAll
	}
	return lookupUnset
}

func hasStringValue(value types.String) bool {
	return !value.IsNull() && !value.IsUnknown() && strings.TrimSpace(value.ValueString()) != ""
}

func dataSourceID(object NetworkForwardingProfileDataSourceModel, method lookupMethod) string {
	switch method {
	case lookupByForwardingProfileID:
		return "forwarding_profile_id/" + object.ForwardingProfileID.ValueString()
	case lookupByName:
		return "name/" + object.Name.ValueString()
	case lookupByTrafficForwardingType:
		return "traffic_forwarding_type/" + object.TrafficForwardingType.ValueString()
	case lookupListAll:
		return "list_all"
	default:
		return "unknown"
	}
}
