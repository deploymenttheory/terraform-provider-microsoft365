package graphBetaNetworkTLSInspectionPolicy

import (
	"fmt"

	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

type tlsInspectionPolicyResponse struct {
	id                   *string
	name                 *string
	description          *string
	defaultAction        *string
	version              *string
	lastModifiedDateTime *string
}

func createTLSInspectionPolicyResponseFromDiscriminatorValue(
	parseNode s.ParseNode,
) (s.Parsable, error) {
	return &tlsInspectionPolicyResponse{}, nil
}

func (r *tlsInspectionPolicyResponse) Serialize(writer s.SerializationWriter) error { return nil }

func (r *tlsInspectionPolicyResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"id": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return wrapResponseError(err)
			}
			r.id = value
			return nil
		},
		"name": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return wrapResponseError(err)
			}
			r.name = value
			return nil
		},
		"description": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return wrapResponseError(err)
			}
			r.description = value
			return nil
		},
		"version": func(n s.ParseNode) error { v, e := n.GetStringValue(); r.version = v; return wrapResponseError(e) },
		"lastModifiedDateTime": func(n s.ParseNode) error {
			v, e := n.GetStringValue()
			r.lastModifiedDateTime = v
			return wrapResponseError(e)
		},
		"settings": func(n s.ParseNode) error {
			value, err := n.GetObjectValue(
				createTLSInspectionPolicySettingsResponseFromDiscriminatorValue,
			)
			if err != nil || value == nil {
				return wrapResponseError(err)
			}
			if settings, ok := value.(*tlsInspectionPolicySettingsResponse); ok {
				r.defaultAction = settings.defaultAction
			}
			return nil
		},
	}
}

type tlsInspectionPolicySettingsResponse struct{ defaultAction *string }

func createTLSInspectionPolicySettingsResponseFromDiscriminatorValue(
	parseNode s.ParseNode,
) (s.Parsable, error) {
	return &tlsInspectionPolicySettingsResponse{}, nil
}

func (r *tlsInspectionPolicySettingsResponse) Serialize(writer s.SerializationWriter) error {
	return nil
}

func (r *tlsInspectionPolicySettingsResponse) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{
		"defaultAction": func(n s.ParseNode) error {
			value, err := n.GetStringValue()
			if err != nil {
				return wrapResponseError(err)
			}
			r.defaultAction = value
			return nil
		},
	}
}

func wrapResponseError(err error) error {
	if err != nil {
		return fmt.Errorf("parse TLS inspection policy response: %w", err)
	}
	return nil
}
