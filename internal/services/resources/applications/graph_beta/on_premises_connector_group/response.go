package graphBetaApplicationsOnPremisesConnectorGroup

import (
	stderrors "errors"
	"fmt"

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

// decodeCreateConnectorGroupResponse validates and decodes the Microsoft Graph response.
func decodeCreateConnectorGroupResponse(result any, err error) (*connectorGroupResponse, error) {
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errCreateResponseNil
	}
	connectorGroup, ok := result.(*connectorGroupResponse)
	if !ok {
		return nil, fmt.Errorf("create connector group returned %T, %w", result, errUnexpectedResponseType)
	}
	return connectorGroup, nil
}

// decodeGetConnectorGroupResponse validates and decodes the Microsoft Graph response.
func decodeGetConnectorGroupResponse(result any, err error) (*connectorGroupResponse, error) {
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil //nolint:nilnil // Empty responses are handled by the lifecycle method.
	}
	connectorGroup, ok := result.(*connectorGroupResponse)
	if !ok {
		return nil, fmt.Errorf("get connector group returned %T, %w", result, errUnexpectedResponseType)
	}
	return connectorGroup, nil
}

// decodeUpdateConnectorGroupResponse validates and decodes the Microsoft Graph response.
func decodeUpdateConnectorGroupResponse(result any, err error) (*connectorGroupResponse, error) {
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil //nolint:nilnil // Empty responses are handled by the lifecycle method.
	}
	connectorGroup, ok := result.(*connectorGroupResponse)
	if !ok {
		return nil, fmt.Errorf("update connector group returned %T, %w", result, errUnexpectedResponseType)
	}
	return connectorGroup, nil
}

var (
	errCreateResponseNil      = stderrors.New("create connector group returned nil response")
	errUnexpectedResponseType = stderrors.New("expected connectorGroupResponse")
)
