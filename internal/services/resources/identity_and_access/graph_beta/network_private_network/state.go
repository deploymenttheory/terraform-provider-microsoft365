package graphBetaNetworkPrivateNetwork

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// expectedIPResolutionAttrTypes are the attribute types of a dns_resolution_identification.expected_ip_resolutions
// element.
var expectedIPResolutionAttrTypes = map[string]attr.Type{
	"type":          types.StringType,
	"value":         types.StringType,
	"begin_address": types.StringType,
	"end_address":   types.StringType,
}

func MapRemoteStateToTerraform(ctx context.Context, data *NetworkPrivateNetworkResourceModel, remoteResource *privateNetworkResponse) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.ID)
	data.Name = convert.GraphToFrameworkString(remoteResource.Name)

	appIDs := make([]attr.Value, 0, len(remoteResource.AppIDs))
	for _, appID := range remoteResource.AppIDs {
		appIDs = append(appIDs, types.StringValue(appID))
	}
	data.AppIDs = types.SetValueMust(types.StringType, appIDs)

	for _, identification := range remoteResource.NetworkIdentifications {
		if identification == nil {
			continue
		}
		if data.DNSResolutionIdentification == nil {
			data.DNSResolutionIdentification = &DNSResolutionIdentificationModel{}
		}

		dnsServers := make([]attr.Value, 0, len(identification.ServerAddresses))
		for _, server := range identification.ServerAddresses {
			if server.Value != nil {
				dnsServers = append(dnsServers, types.StringValue(*server.Value))
			}
		}
		data.DNSResolutionIdentification.DNSServers = types.SetValueMust(types.StringType, dnsServers)

		if identification.FQDNToResolve != nil {
			data.DNSResolutionIdentification.FQDNToResolve = convert.GraphToFrameworkString(identification.FQDNToResolve.Value)
		}

		expectedIPResolutions := make([]attr.Value, 0, len(identification.ExpectedIPResolutions))
		for _, resolution := range identification.ExpectedIPResolutions {
			// The @odata.type discriminates the resolution; an address is the default.
			resolutionType := expectedIPResolutionTypeIPAddress
			if resolution.ODataType != nil {
				odataType := strings.ToLower(*resolution.ODataType)
				switch {
				case strings.HasSuffix(odataType, "ipsubnet"):
					resolutionType = expectedIPResolutionTypeIPSubnet
				case strings.HasSuffix(odataType, "iprange"):
					resolutionType = expectedIPResolutionTypeIPRange
				}
			}
			attrs := map[string]attr.Value{
				"type":          types.StringValue(resolutionType),
				"value":         convert.GraphToFrameworkString(resolution.Value),
				"begin_address": convert.GraphToFrameworkString(resolution.BeginAddress),
				"end_address":   convert.GraphToFrameworkString(resolution.EndAddress),
			}
			expectedIPResolutions = append(expectedIPResolutions, types.ObjectValueMust(expectedIPResolutionAttrTypes, attrs))
		}
		data.DNSResolutionIdentification.ExpectedIPResolutions = types.SetValueMust(types.ObjectType{AttrTypes: expectedIPResolutionAttrTypes}, expectedIPResolutions)
		break
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
