package graphBetaNetworkForwardingProfilePolicyLink

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

func constructResource(data *NetworkForwardingProfilePolicyLinkResourceModel) s.Parsable {
	return &forwardingPolicyLinkStateRequestBody{
		ODataType: forwardingPolicyLinkODataType,
		State:     data.State.ValueString(),
	}
}

type forwardingPolicyLinkStateRequestBody struct {
	ODataType string
	State     string
}

func (b *forwardingPolicyLinkStateRequestBody) Serialize(writer s.SerializationWriter) error {
	if b.ODataType != "" {
		if err := writer.WriteStringValue("@odata.type", &b.ODataType); err != nil {
			return err
		}
	}
	if b.State != "" {
		if err := writer.WriteStringValue("state", &b.State); err != nil {
			return err
		}
	}
	return nil
}

func (b *forwardingPolicyLinkStateRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

const (
	forwardingProfilePolicyLinkURLTemplate = "{+baseurl}/networkAccess/forwardingProfiles/{forwardingProfileId}/policies/{policyLinkId}"
)

func newForwardingPolicyLinkRequestInformation(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	forwardingProfileID, policyLinkID string,
	requestBody s.Parsable,
) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{
		"baseurl":             adapter.GetBaseUrl(),
		"forwardingProfileId": forwardingProfileID,
		"policyLinkId":        policyLinkID,
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, forwardingProfilePolicyLinkURLTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set forwarding policy link request content: %w", err)
		}
	}
	return requestInfo, nil
}
