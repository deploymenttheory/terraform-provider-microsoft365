package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

func constructResource(ctx context.Context, data *OnPremisesIpApplicationSegmentResourceModel) (s.Parsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := &ipApplicationSegmentRequestBody{
		port: 0,
	}

	if !data.DestinationHost.IsNull() && !data.DestinationHost.IsUnknown() {
		requestBody.destinationHost = data.DestinationHost.ValueString()
	}

	if !data.DestinationType.IsNull() && !data.DestinationType.IsUnknown() {
		requestBody.destinationType = graphDestinationType(data.DestinationType.ValueString())
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.Ports, func(ports []string) {
		requestBody.ports = ports
	}); err != nil {
		return nil, fmt.Errorf("failed to set ports: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.Protocol, func(protocols []string) {
		requestBody.protocol = graphProtocols(protocols)
	}); err != nil {
		return nil, fmt.Errorf("failed to set protocol: %w", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

func graphDestinationType(destinationType string) string {
	// Microsoft Learn documents ipApplicationSegment.destinationType values as
	// "ipAddress", "ipRange", "ipRangeCidr", "fqdn", and "dnsSuffix":
	// https://learn.microsoft.com/en-us/graph/api/resources/ipapplicationsegment?view=graph-rest-beta
	// The create API page also lists "ipAddress" for POST:
	// https://learn.microsoft.com/en-us/graph/api/onpremisespublishingprofile-post-applicationsegments?view=graph-rest-beta
	//
	// The beta endpoint currently behaves differently for a single IP segment:
	// POSTing destinationType "ipAddress" returns 400 InvalidJson_BadRequest
	// ("Valid JSON content expected."), while destinationType "ip" succeeds.
	// Direct API checks also showed dnsSuffix returning 400
	// Invalid_AppSegments_NonwebApp_Duplicate ("DNS suffix can only be added to
	// Quick Access configuration") and ipRange with a host range returning
	// DestinationHost_InvalidIP. The schema only allows values observed to create
	// and read back through this application-scoped endpoint. Wildcard hosts such
	// as "*.internal.example.com" are supported when sent as destinationType fqdn;
	// they should not be modeled as dnsSuffix for this endpoint.
	// Keep Terraform's public schema aligned with Learn, but send the literal
	// accepted by Graph.
	if destinationType == "ipAddress" {
		return "ip"
	}

	return destinationType
}

func graphProtocols(protocols []string) string {
	sort.Strings(protocols)
	return strings.Join(protocols, ",")
}

type ipApplicationSegmentRequestBody struct {
	destinationHost string
	destinationType string
	// Learn marks "port" as deprecated / DO NOT USE, but its create example
	// still includes "port": 0. The real beta endpoint also returns "port": 0
	// together with "ports": ["443-443"] for segments created by Terraform.
	// Include the scalar field to match the endpoint's accepted wire shape while
	// using "ports" for the actual range values in Terraform state.
	port     int32
	ports    []string
	protocol string
}

func (b *ipApplicationSegmentRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("destinationHost", &b.destinationHost); err != nil {
		return err
	}
	if err := writer.WriteStringValue("destinationType", &b.destinationType); err != nil {
		return err
	}
	if err := writer.WriteInt32Value("port", &b.port); err != nil {
		return err
	}
	if err := writer.WriteCollectionOfStringValues("ports", b.ports); err != nil {
		return err
	}
	if err := writer.WriteStringValue("protocol", &b.protocol); err != nil {
		return err
	}

	return nil
}

func (b *ipApplicationSegmentRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}
