package graphBetaNetworkPrivateNetwork

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

const (
	privateNetworkODataType              = "#microsoft.graph.networkaccess.privateNetwork"
	dnsResolutionIdentificationODataType = "#microsoft.graph.networkaccess.dnsResolutionIdentification"
	ipAddressODataType                   = "#microsoft.graph.networkaccess.ipAddress"
	ipSubnetODataType                    = "#microsoft.graph.networkaccess.ipSubnet"
	ipRangeODataType                     = "#microsoft.graph.networkaccess.ipRange"
	fqdnODataType                        = "#microsoft.graph.networkaccess.fqdn"

	expectedIPResolutionTypeIPAddress = "ip_address"
	expectedIPResolutionTypeIPSubnet  = "ip_subnet"
	expectedIPResolutionTypeIPRange   = "ip_range"
)

// constructResource builds the portal-observed private network payload for
// /networkaccess/privateNetworks.
//
// Microsoft Graph beta currently exposes this Global Secure Access private
// network surface in the Entra portal before the generated Go SDK has a typed
// request builder for it, so this resource uses Kiota RequestInformation and
// custom Parsable request bodies while preserving the portal-observed payload
// shape.
func constructResource(ctx context.Context, data *NetworkPrivateNetworkResourceModel, includeODataType bool) (s.Parsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if data.DNSResolutionIdentification == nil {
		return nil, fmt.Errorf("dns_resolution_identification must be specified")
	}

	dnsServers, err := stringSetValues(ctx, data.DNSResolutionIdentification.DNSServers)
	if err != nil {
		return nil, fmt.Errorf("failed to read dns_resolution_identification.dns_servers: %w", err)
	}
	if len(dnsServers) == 0 {
		return nil, fmt.Errorf("at least one dns_resolution_identification.dns_servers value must be specified")
	}
	for _, server := range dnsServers {
		if _, err := netip.ParseAddr(server); err != nil {
			return nil, fmt.Errorf("dns server %q must be a valid IP address: %w", server, err)
		}
	}

	expectedIPResolutions, err := expectedIPResolutionValues(ctx, data.DNSResolutionIdentification.ExpectedIPResolutions)
	if err != nil {
		return nil, err
	}
	if len(expectedIPResolutions) == 0 {
		return nil, fmt.Errorf("at least one dns_resolution_identification.expected_ip_resolutions value must be specified")
	}

	appIDs, err := stringSetValues(ctx, data.AppIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to read app_ids: %w", err)
	}

	requestBody := &privateNetworkRequestBody{
		includeODataType: includeODataType,
		name:             data.Name.ValueStringPointer(),
		appIDs:           appIDs,
		networkIdentifications: []s.Parsable{
			&dnsResolutionIdentificationRequestBody{
				serverAddresses:       ipAddressRequestBodies(dnsServers),
				fqdnToResolve:         &fqdnRequestBody{value: data.DNSResolutionIdentification.FQDNToResolve.ValueStringPointer()},
				expectedIPResolutions: expectedIPResolutions,
			},
		},
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

func stringSetValues(ctx context.Context, value types.Set) ([]string, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	var result []string
	diags := value.ElementsAs(ctx, &result, false)
	if diags.HasError() {
		return nil, fmt.Errorf("%s", diags.Errors()[0].Detail())
	}

	return result, nil
}

func expectedIPResolutionValues(ctx context.Context, value types.Set) ([]s.Parsable, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	var terraformValues []ExpectedIPResolutionModel
	diags := value.ElementsAs(ctx, &terraformValues, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to read dns_resolution_identification.expected_ip_resolutions: %s", diags.Errors()[0].Detail())
	}

	result := make([]s.Parsable, 0, len(terraformValues))
	for _, item := range terraformValues {
		switch item.Type.ValueString() {
		case expectedIPResolutionTypeIPAddress:
			if item.Value.IsNull() || item.Value.IsUnknown() || item.Value.ValueString() == "" {
				return nil, fmt.Errorf("expected_ip_resolutions value is required when type is %q", expectedIPResolutionTypeIPAddress)
			}
			if _, err := netip.ParseAddr(item.Value.ValueString()); err != nil {
				return nil, fmt.Errorf("expected_ip_resolutions value %q must be a valid IP address: %w", item.Value.ValueString(), err)
			}
			result = append(result, &ipAddressRequestBody{value: item.Value.ValueStringPointer()})
		case expectedIPResolutionTypeIPSubnet:
			if item.Value.IsNull() || item.Value.IsUnknown() || item.Value.ValueString() == "" {
				return nil, fmt.Errorf("expected_ip_resolutions value is required when type is %q", expectedIPResolutionTypeIPSubnet)
			}
			if _, err := netip.ParsePrefix(item.Value.ValueString()); err != nil {
				return nil, fmt.Errorf("expected_ip_resolutions value %q must be a valid CIDR prefix: %w", item.Value.ValueString(), err)
			}
			result = append(result, &ipSubnetRequestBody{value: item.Value.ValueStringPointer()})
		case expectedIPResolutionTypeIPRange:
			if item.BeginAddress.IsNull() || item.BeginAddress.IsUnknown() || item.BeginAddress.ValueString() == "" {
				return nil, fmt.Errorf("expected_ip_resolutions begin_address is required when type is %q", expectedIPResolutionTypeIPRange)
			}
			if item.EndAddress.IsNull() || item.EndAddress.IsUnknown() || item.EndAddress.ValueString() == "" {
				return nil, fmt.Errorf("expected_ip_resolutions end_address is required when type is %q", expectedIPResolutionTypeIPRange)
			}
			if _, err := netip.ParseAddr(item.BeginAddress.ValueString()); err != nil {
				return nil, fmt.Errorf("expected_ip_resolutions begin_address %q must be a valid IP address: %w", item.BeginAddress.ValueString(), err)
			}
			if _, err := netip.ParseAddr(item.EndAddress.ValueString()); err != nil {
				return nil, fmt.Errorf("expected_ip_resolutions end_address %q must be a valid IP address: %w", item.EndAddress.ValueString(), err)
			}
			result = append(result, &ipRangeRequestBody{
				beginAddress: item.BeginAddress.ValueStringPointer(),
				endAddress:   item.EndAddress.ValueStringPointer(),
			})
		default:
			return nil, fmt.Errorf("expected_ip_resolutions type must be one of %q, %q, or %q", expectedIPResolutionTypeIPAddress, expectedIPResolutionTypeIPSubnet, expectedIPResolutionTypeIPRange)
		}
	}

	return result, nil
}

func ipAddressRequestBodies(values []string) []s.Parsable {
	result := make([]s.Parsable, 0, len(values))
	for _, value := range values {
		v := value
		result = append(result, &ipAddressRequestBody{value: &v})
	}
	return result
}

type privateNetworkRequestBody struct {
	includeODataType       bool
	name                   *string
	appIDs                 []string
	networkIdentifications []s.Parsable
}

func (b *privateNetworkRequestBody) Serialize(writer s.SerializationWriter) error {
	if b.includeODataType {
		odataType := privateNetworkODataType
		if err := writer.WriteStringValue("@odata.type", &odataType); err != nil {
			return err
		}
	}
	if err := writer.WriteStringValue("name", b.name); err != nil {
		return err
	}
	if err := writer.WriteCollectionOfStringValues("appIds", b.appIDs); err != nil {
		return err
	}
	if err := writer.WriteCollectionOfObjectValues("networkIdentifications", b.networkIdentifications); err != nil {
		return err
	}
	return nil
}

func (b *privateNetworkRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type dnsResolutionIdentificationRequestBody struct {
	serverAddresses       []s.Parsable
	fqdnToResolve         *fqdnRequestBody
	expectedIPResolutions []s.Parsable
}

func (b *dnsResolutionIdentificationRequestBody) Serialize(writer s.SerializationWriter) error {
	odataType := dnsResolutionIdentificationODataType
	if err := writer.WriteStringValue("@odata.type", &odataType); err != nil {
		return err
	}
	if err := writer.WriteCollectionOfObjectValues("serverAddresses", b.serverAddresses); err != nil {
		return err
	}
	if err := writer.WriteObjectValue("fqdnToResolve", b.fqdnToResolve); err != nil {
		return err
	}
	if err := writer.WriteCollectionOfObjectValues("expectedIpResolutions", b.expectedIPResolutions); err != nil {
		return err
	}
	return nil
}

func (b *dnsResolutionIdentificationRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type fqdnRequestBody struct {
	value *string
}

func (b *fqdnRequestBody) Serialize(writer s.SerializationWriter) error {
	odataType := fqdnODataType
	if err := writer.WriteStringValue("@odata.type", &odataType); err != nil {
		return err
	}
	return writer.WriteStringValue("value", b.value)
}

func (b *fqdnRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type ipAddressRequestBody struct {
	value *string
}

func (b *ipAddressRequestBody) Serialize(writer s.SerializationWriter) error {
	odataType := ipAddressODataType
	if err := writer.WriteStringValue("@odata.type", &odataType); err != nil {
		return err
	}
	return writer.WriteStringValue("value", b.value)
}

func (b *ipAddressRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type ipSubnetRequestBody struct {
	value *string
}

func (b *ipSubnetRequestBody) Serialize(writer s.SerializationWriter) error {
	odataType := ipSubnetODataType
	if err := writer.WriteStringValue("@odata.type", &odataType); err != nil {
		return err
	}
	return writer.WriteStringValue("value", b.value)
}

func (b *ipSubnetRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type ipRangeRequestBody struct {
	beginAddress *string
	endAddress   *string
}

func (b *ipRangeRequestBody) Serialize(writer s.SerializationWriter) error {
	odataType := ipRangeODataType
	if err := writer.WriteStringValue("@odata.type", &odataType); err != nil {
		return err
	}
	if err := writer.WriteStringValue("beginAddress", b.beginAddress); err != nil {
		return err
	}
	return writer.WriteStringValue("endAddress", b.endAddress)
}

func (b *ipRangeRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}
