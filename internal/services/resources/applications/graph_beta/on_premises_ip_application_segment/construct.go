package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"fmt"

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

	if !data.Protocol.IsNull() && !data.Protocol.IsUnknown() {
		requestBody.protocol = data.Protocol.ValueString()
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
	// Microsoft Learn and the Graph metadata model this value as "ipAddress",
	// but the beta endpoint currently rejects that value and accepts "ip".
	if destinationType == "ipAddress" {
		return "ip"
	}

	return destinationType
}

type ipApplicationSegmentRequestBody struct {
	destinationHost string
	destinationType string
	// The beta endpoint expects the legacy scalar port field to be present even
	// when the modern ports collection carries the actual range values.
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
