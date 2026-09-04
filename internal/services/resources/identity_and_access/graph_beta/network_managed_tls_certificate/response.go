package graphBetaNetworkManagedTLSCertificate

import s "github.com/microsoft/kiota-abstractions-go/serialization"

type managedTLSCertificateResponse struct {
	id               *string
	name             *string
	commonName       *string
	organizationName *string
	validityMonths   *int32
	status           *string
	createdDateTime  *string
	certificate      *string
	validity         *managedTLSCertificateValidityResponse
}

func createManagedTLSCertificateResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &managedTLSCertificateResponse{}, nil
}

func (r *managedTLSCertificateResponse) Serialize(writer s.SerializationWriter) error { return nil }

func (r *managedTLSCertificateResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"id": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			r.id = value
			return err
		},
		"name": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			r.name = value
			return err
		},
		"commonName": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			r.commonName = value
			return err
		},
		"organizationName": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			r.organizationName = value
			return err
		},
		"validityMonths": func(n s.ParseNode) error {
			value, err := n.GetInt32Value()
			r.validityMonths = value
			return err
		},
		"status": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			r.status = value
			return err
		},
		"createdDateTime": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			r.createdDateTime = value
			return err
		},
		"certificate": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			r.certificate = value
			return err
		},
		"validity": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createManagedTLSCertificateValidityResponseFromDiscriminatorValue)
			if err != nil || value == nil {
				return err
			}
			if validity, ok := value.(*managedTLSCertificateValidityResponse); ok {
				r.validity = validity
			}
			return nil
		},
	}
}

type managedTLSCertificateValidityResponse struct {
	startDateTime *string
	endDateTime   *string
}

func createManagedTLSCertificateValidityResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &managedTLSCertificateValidityResponse{}, nil
}

func (r *managedTLSCertificateValidityResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *managedTLSCertificateValidityResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"startDateTime": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			r.startDateTime = value
			return err
		},
		"endDateTime": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			r.endDateTime = value
			return err
		},
	}
}
