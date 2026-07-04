package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ipApplicationSegmentResponse struct {
	id              *string
	destinationHost *string
	// Parse destinationType as a raw string because the generated enum parser
	// drops the endpoint's observed "ip" value while docs describe "ipAddress".
	destinationType *string
	ports           []string
	protocol        *string
}

func newIpApplicationSegmentResponse() *ipApplicationSegmentResponse {
	return &ipApplicationSegmentResponse{}
}

func createIpApplicationSegmentResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return newIpApplicationSegmentResponse(), nil
}

func (r *ipApplicationSegmentResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *ipApplicationSegmentResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"id": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.id = value
			return nil
		},
		"destinationHost": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.destinationHost = value
			return nil
		},
		"destinationType": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.destinationType = value
			return nil
		},
		"ports": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfPrimitiveValues("string")
			if err != nil {
				return err
			}
			if values == nil {
				return nil
			}

			r.ports = make([]string, 0, len(values))
			for _, value := range values {
				if value == nil {
					continue
				}
				r.ports = append(r.ports, *(value.(*string)))
			}
			return nil
		},
		"protocol": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.protocol = value
			return nil
		},
	}
}
