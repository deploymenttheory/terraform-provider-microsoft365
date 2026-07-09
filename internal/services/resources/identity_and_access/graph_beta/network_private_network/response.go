package graphBetaNetworkPrivateNetwork

import (
	"strings"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

type privateNetworkResponse struct {
	id                     *string
	name                   *string
	appIDs                 []string
	networkIdentifications []*dnsResolutionIdentificationResponse
}

func createPrivateNetworkResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &privateNetworkResponse{}, nil
}

func (r *privateNetworkResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *privateNetworkResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
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
		"appIds": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfPrimitiveValues("string")
			if err != nil {
				return err
			}
			if values == nil {
				r.appIDs = nil
				return nil
			}
			r.appIDs = make([]string, 0, len(values))
			for _, value := range values {
				switch v := value.(type) {
				case string:
					r.appIDs = append(r.appIDs, v)
				case *string:
					if v != nil {
						r.appIDs = append(r.appIDs, *v)
					}
				}
			}
			return nil
		},
		"networkIdentifications": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfObjectValues(createNetworkIdentificationResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			r.networkIdentifications = nil
			for _, value := range values {
				if typed, ok := value.(*dnsResolutionIdentificationResponse); ok {
					r.networkIdentifications = append(r.networkIdentifications, typed)
				}
			}
			return nil
		},
	}
}

type dnsResolutionIdentificationResponse struct {
	odataType             *string
	serverAddresses       []ipAddressResponse
	fqdnToResolve         *fqdnResponse
	expectedIPResolutions []ipResolutionResponse
}

func createNetworkIdentificationResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &dnsResolutionIdentificationResponse{}, nil
}

func (r *dnsResolutionIdentificationResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *dnsResolutionIdentificationResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"@odata.type": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.odataType = value
			return nil
		},
		"serverAddresses": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfObjectValues(createIPAddressResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			r.serverAddresses = nil
			for _, value := range values {
				if typed, ok := value.(*ipAddressResponse); ok {
					r.serverAddresses = append(r.serverAddresses, *typed)
				}
			}
			return nil
		},
		"fqdnToResolve": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(createFQDNResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			if typed, ok := value.(*fqdnResponse); ok {
				r.fqdnToResolve = typed
			}
			return nil
		},
		"expectedIpResolutions": func(n s.ParseNode) error {
			values, err := n.GetCollectionOfObjectValues(createIPResolutionResponseFromDiscriminatorValue)
			if err != nil {
				return err
			}
			r.expectedIPResolutions = nil
			for _, value := range values {
				if typed, ok := value.(*ipResolutionResponse); ok {
					r.expectedIPResolutions = append(r.expectedIPResolutions, *typed)
				}
			}
			return nil
		},
	}
}

type ipAddressResponse struct {
	value *string
}

func createIPAddressResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &ipAddressResponse{}, nil
}

func (r *ipAddressResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *ipAddressResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"value": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.value = value
			return nil
		},
	}
}

type fqdnResponse struct {
	value *string
}

func createFQDNResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &fqdnResponse{}, nil
}

func (r *fqdnResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *fqdnResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"value": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.value = value
			return nil
		},
	}
}

type ipResolutionResponse struct {
	odataType    *string
	value        *string
	beginAddress *string
	endAddress   *string
}

func createIPResolutionResponseFromDiscriminatorValue(parseNode s.ParseNode) (s.Parsable, error) {
	return &ipResolutionResponse{}, nil
}

func (r *ipResolutionResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *ipResolutionResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"@odata.type": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.odataType = value
			return nil
		},
		"value": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.value = value
			return nil
		},
		"beginAddress": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.beginAddress = value
			return nil
		},
		"endAddress": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return err
			}
			r.endAddress = value
			return nil
		},
	}
}

func terraformIPResolutionType(odataType *string) string {
	if odataType == nil {
		return expectedIPResolutionTypeIPAddress
	}

	value := strings.ToLower(*odataType)
	switch {
	case strings.HasSuffix(value, "ipsubnet"):
		return expectedIPResolutionTypeIPSubnet
	case strings.HasSuffix(value, "iprange"):
		return expectedIPResolutionTypeIPRange
	default:
		return expectedIPResolutionTypeIPAddress
	}
}
