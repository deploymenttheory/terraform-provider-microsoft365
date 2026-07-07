package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ipApplicationSegmentResponse intentionally parses only the fields this
// Terraform resource manages. The generated SDK model is not used here because
// its destinationType enum is based on Learn / metadata values, while the real
// beta endpoint returns "ip" for a single IP address segment.
//
// Learn says destinationType supports "ipAddress":
// https://learn.microsoft.com/en-us/graph/api/resources/ipapplicationsegment?view=graph-rest-beta
//
// Observed response after creating a segment with Terraform:
//
//	{
//	  "@odata.context": ".../applications('{applicationObjectId}')/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments/$entity",
//	  "action": "tunnel",
//	  "destinationHost": "10.10.10.10",
//	  "destinationType": "ip",
//	  "exclusions": null,
//	  "id": "{ipApplicationSegmentId}",
//	  "inclusions": null,
//	  "port": 0,
//	  "ports": ["443-443"],
//	  "protocol": "tcp"
//	}
//
// Parsing destinationType as a raw string preserves that observed "ip" literal
// so state mapping can translate it back to Terraform's documented "ipAddress"
// value instead of losing it through the generated enum parser.
type ipApplicationSegmentResponse struct {
	id              *string
	destinationHost *string
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
