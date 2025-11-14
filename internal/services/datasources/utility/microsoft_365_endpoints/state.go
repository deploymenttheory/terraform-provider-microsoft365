package utilityMicrosoft365Endpoints

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// mapResponseToState converts API response data into Terraform state
func mapResponseToState(ctx context.Context, data *Microsoft365EndpointsDataSourceModel, endpoints []Microsoft365EndpointResponse) error {

	endpointObjects := make([]attr.Value, 0, len(endpoints))

	for _, endpoint := range endpoints {
		urls := make([]attr.Value, 0, len(endpoint.URLs))
		for _, url := range endpoint.URLs {
			urls = append(urls, types.StringValue(url))
		}
		urlsList, diag := types.ListValue(types.StringType, urls)
		if diag.HasError() {
			return fmt.Errorf("failed to create URLs list: %v", diag.Errors())
		}

		if len(endpoint.URLs) == 0 {
			urlsList = types.ListNull(types.StringType)
		}

		ips := make([]attr.Value, 0, len(endpoint.IPs))
		for _, ip := range endpoint.IPs {
			ips = append(ips, types.StringValue(ip))
		}
		ipsList, diag := types.ListValue(types.StringType, ips)
		if diag.HasError() {
			return fmt.Errorf("failed to create IPs list: %v", diag.Errors())
		}

		if len(endpoint.IPs) == 0 {
			ipsList = types.ListNull(types.StringType)
		}

		endpointObj, diag := types.ObjectValue(
			map[string]attr.Type{
				"id":                        types.Int64Type,
				"service_area":              types.StringType,
				"service_area_display_name": types.StringType,
				"urls":                      types.ListType{ElemType: types.StringType},
				"ips":                       types.ListType{ElemType: types.StringType},
				"tcp_ports":                 types.StringType,
				"udp_ports":                 types.StringType,
				"express_route":             types.BoolType,
				"category":                  types.StringType,
				"required":                  types.BoolType,
				"notes":                     types.StringType,
			},
			map[string]attr.Value{
				"id":                        types.Int64Value(endpoint.ID),
				"service_area":              types.StringValue(endpoint.ServiceArea),
				"service_area_display_name": types.StringValue(endpoint.ServiceAreaDisplayName),
				"urls":                      urlsList,
				"ips":                       ipsList,
				"tcp_ports":                 stringOrNull(endpoint.TCPPorts),
				"udp_ports":                 stringOrNull(endpoint.UDPPorts),
				"express_route":             types.BoolValue(endpoint.ExpressRoute),
				"category":                  types.StringValue(endpoint.Category),
				"required":                  types.BoolValue(endpoint.Required),
				"notes":                     stringOrNull(endpoint.Notes),
			},
		)
		if diag.HasError() {
			return fmt.Errorf("failed to create endpoint object: %v", diag.Errors())
		}

		endpointObjects = append(endpointObjects, endpointObj)
	}

	endpointsList, diag := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":                        types.Int64Type,
				"service_area":              types.StringType,
				"service_area_display_name": types.StringType,
				"urls":                      types.ListType{ElemType: types.StringType},
				"ips":                       types.ListType{ElemType: types.StringType},
				"tcp_ports":                 types.StringType,
				"udp_ports":                 types.StringType,
				"express_route":             types.BoolType,
				"category":                  types.StringType,
				"required":                  types.BoolType,
				"notes":                     types.StringType,
			},
		},
		endpointObjects,
	)
	if diag.HasError() {
		return fmt.Errorf("failed to create endpoints list: %v", diag.Errors())
	}

	data.Endpoints = endpointsList

	tflog.Trace(ctx, fmt.Sprintf("Populated state with %d endpoints", len(endpointObjects)))

	return nil
}

// stringOrNull returns a StringValue if the string is non-empty, otherwise returns StringNull
func stringOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
