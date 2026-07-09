package graphBetaNetworkPrivateNetwork

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MapRemoteStateToTerraform(ctx context.Context, data *NetworkPrivateNetworkResourceModel, remoteResource *privateNetworkResponse) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.id)
	data.Name = convert.GraphToFrameworkString(remoteResource.name)
	data.AppIDs = stringSliceToSet(remoteResource.appIDs)

	for _, identification := range remoteResource.networkIdentifications {
		if identification == nil {
			continue
		}
		if data.DNSResolutionIdentification == nil {
			data.DNSResolutionIdentification = &DNSResolutionIdentificationModel{}
		}
		data.DNSResolutionIdentification.DNSServers = dnsServerSet(identification.serverAddresses)
		if identification.fqdnToResolve != nil {
			data.DNSResolutionIdentification.FQDNToResolve = convert.GraphToFrameworkString(identification.fqdnToResolve.value)
		}
		data.DNSResolutionIdentification.ExpectedIPResolutions = expectedIPResolutionSet(identification.expectedIPResolutions)
		break
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

func stringSliceToSet(values []string) types.Set {
	if values == nil {
		return types.SetValueMust(types.StringType, []attr.Value{})
	}

	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.StringValue(value))
	}

	return types.SetValueMust(types.StringType, elements)
}

func dnsServerSet(values []ipAddressResponse) types.Set {
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		if value.value != nil {
			elements = append(elements, types.StringValue(*value.value))
		}
	}

	return types.SetValueMust(types.StringType, elements)
}

func expectedIPResolutionSet(values []ipResolutionResponse) types.Set {
	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		attrs := map[string]attr.Value{
			"type":          types.StringValue(terraformIPResolutionType(value.odataType)),
			"value":         types.StringNull(),
			"begin_address": types.StringNull(),
			"end_address":   types.StringNull(),
		}
		if value.value != nil {
			attrs["value"] = types.StringValue(*value.value)
		}
		if value.beginAddress != nil {
			attrs["begin_address"] = types.StringValue(*value.beginAddress)
		}
		if value.endAddress != nil {
			attrs["end_address"] = types.StringValue(*value.endAddress)
		}

		elements = append(elements, types.ObjectValueMust(expectedIPResolutionAttrTypes(), attrs))
	}

	return types.SetValueMust(expectedIPResolutionObjectType(), elements)
}

func expectedIPResolutionObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: expectedIPResolutionAttrTypes()}
}

func expectedIPResolutionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":          types.StringType,
		"value":         types.StringType,
		"begin_address": types.StringType,
		"end_address":   types.StringType,
	}
}
