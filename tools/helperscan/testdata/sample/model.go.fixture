package sample

import "example.com/kiota/serialization"

// widgetResponse is a hand-written Kiota model.
type widgetResponse struct {
	additionalData map[string]any
}

// createWidgetResponseFromDiscriminatorValue is a Kiota factory.
func createWidgetResponseFromDiscriminatorValue(parseNode serialization.ParseNode) (serialization.Parsable, error) {
	return &widgetResponse{}, nil
}

// Serialize is a Kiota interface method.
func (m *widgetResponse) Serialize(writer serialization.SerializationWriter) error {
	return nil
}

// GetFieldDeserializers is a Kiota interface method.
func (m *widgetResponse) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return nil
}

// GetAdditionalData is a Kiota interface method.
func (m *widgetResponse) GetAdditionalData() map[string]any {
	return m.additionalData
}
