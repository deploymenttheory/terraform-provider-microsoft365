package graphBetaApplicationsOnPremisesConnectorGroup

import (
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// connectorGroupResponse intentionally parses connectorGroupType and region as
// raw strings instead of using the generated msgraph-beta-sdk-go ConnectorGroup
// model. The generated request builders exist for:
//
//	/onPremisesPublishingProfiles/{id}/connectorGroups
//
// but the generated ConnectorGroupRegion enum comes from the Microsoft Graph
// beta OData CSDL metadata endpoint and only contains nam, eur, aus, asia, ind,
// and unknownFutureValue. Direct API verification on 2026-07-05 using:
//
//	GET /beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups
//
// returned the default connector group with "region": "japan". Microsoft Learn
// also documents connectorGroupType values that are not present in the current
// beta metadata/SDK enum:
// https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta
//
// Using a dedicated response parser preserves actual API values that the SDK
// enum parser would otherwise drop. The same live verification showed newly
// created connector groups return isDefault=false while the imported tenant
// default group returns isDefault=true. Requests still use Kiota
// RequestInformation and the provider's Graph adapter, so authentication,
// middleware, retry, and OData error handling stay on the normal SDK pipeline.
type connectorGroupResponse struct {
	id                 *string
	name               *string
	connectorGroupType *string
	isDefault          *bool
	region             *string
}

func newConnectorGroupResponse() *connectorGroupResponse {
	return &connectorGroupResponse{}
}

func createConnectorGroupResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return newConnectorGroupResponse(), nil
}

func (r *connectorGroupResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *connectorGroupResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"id": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.id = value
			return nil
		},
		"name": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.name = value
			return nil
		},
		"connectorGroupType": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.connectorGroupType = value
			return nil
		},
		"isDefault": func(n s.ParseNode) error {
			value, err := n.GetBoolValue()
			if err != nil {
				return err
			}
			r.isDefault = value
			return nil
		},
		"region": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.region = value
			return nil
		},
	}
}
