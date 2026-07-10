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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/networkaccess"
)

func (d *NetworkForwardingProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object NetworkForwardingProfileDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filterType := object.FilterType.ValueString()
	if filterType != "all" && (object.FilterValue.IsNull() || object.FilterValue.ValueString() == "") {
		resp.Diagnostics.AddError("Missing Required Parameter", fmt.Sprintf("filter_value must be provided when filter_type is %q", filterType))
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	expand := []string{"policies($expand=policy)"}
	items := make([]ForwardingProfileModel, 0)

	if filterType == "id" {
		profile, err := d.client.NetworkAccess().ForwardingProfiles().ByForwardingProfileId(object.FilterValue.ValueString()).Get(ctx, &networkaccess.ForwardingProfilesForwardingProfileItemRequestBuilderGetRequestConfiguration{
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
		filterValue := strings.ToLower(object.FilterValue.ValueString())
		for _, profile := range profiles.GetValue() {
			mapped := MapRemoteStateToDataSource(profile)
			switch filterType {
			case "all":
				items = append(items, mapped)
			case "name":
				if strings.Contains(strings.ToLower(mapped.Name.ValueString()), filterValue) {
					items = append(items, mapped)
				}
			case "traffic_forwarding_type":
				if strings.EqualFold(mapped.TrafficForwardingType.ValueString(), object.FilterValue.ValueString()) {
					items = append(items, mapped)
				}
			}
		}
	}

	object.Items = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished datasource read for %s, found %d items", DataSourceName, len(items)))
}
