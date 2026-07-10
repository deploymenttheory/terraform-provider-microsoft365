package graphBetaNetworkForwardingProfilePolicyLink

import s "github.com/microsoft/kiota-abstractions-go/serialization"

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
