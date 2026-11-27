package convert

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

var errInvalidNormalizedJSON = errors.New("invalid normalized JSON")

// GraphToFrameworkNormalizedJSON preserves configured formatting for equivalent
// remote JSON while retaining actual changes for Terraform drift detection.
func GraphToFrameworkNormalizedJSON(ctx context.Context, remote json.RawMessage, prior jsontypes.Normalized) (jsontypes.Normalized, error) {
	if len(remote) == 0 || string(remote) == "null" {
		return jsontypes.NewNormalizedNull(), nil
	}
	if !json.Valid(remote) {
		return jsontypes.NewNormalizedNull(), fmt.Errorf("%w: graph returned invalid JSON", errInvalidNormalizedJSON)
	}
	var decoded any
	decoder := json.NewDecoder(bytes.NewReader(remote))
	decoder.UseNumber()
	if err := decoder.Decode(&decoded); err != nil {
		return jsontypes.NewNormalizedNull(), fmt.Errorf("decode Graph JSON: %w", err)
	}
	canonical, err := json.Marshal(decoded)
	if err != nil {
		return jsontypes.NewNormalizedNull(), fmt.Errorf("normalize Graph JSON: %w", err)
	}
	value := jsontypes.NewNormalizedValue(string(canonical))
	if !prior.IsNull() && !prior.IsUnknown() {
		equal, diagnostics := prior.StringSemanticEquals(ctx, value)
		if diagnostics.HasError() {
			return value, fmt.Errorf("%w: compare normalized JSON: %s", errInvalidNormalizedJSON, diagnostics)
		}
		if equal {
			return prior, nil
		}
	}
	return value, nil
}

// FrameworkToGraphNormalizedJSONObject omits unset values and preserves JSON value
// types in extensible Kiota additional-data objects.
func FrameworkToGraphNormalizedJSONObject(value jsontypes.Normalized, setter func(map[string]any)) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var fields map[string]any
	if d := value.Unmarshal(&fields); d.HasError() {
		return fmt.Errorf("%w: decode JSON object: %s", errInvalidNormalizedJSON, d)
	}
	if fields == nil {
		return fmt.Errorf("%w: expected a JSON object", errInvalidNormalizedJSON)
	}
	setter(fields)
	return nil
}

// GraphToJSONFields serializes a Kiota object so that additional-data values,
// including untyped dictionaries, retain their wire representation.
func GraphToJSONFields(value serialization.Parsable) (map[string]json.RawMessage, error) {
	if value == nil {
		return map[string]json.RawMessage{}, nil
	}
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", value); err != nil {
		return nil, fmt.Errorf("serialize Graph object: %w", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		return nil, fmt.Errorf("read Graph JSON: %w", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(content, &fields); err != nil {
		return nil, fmt.Errorf("decode Graph JSON fields: %w", err)
	}
	return fields, nil
}

// GraphToFrameworkNormalizedJSONObject preserves configured JSON formatting for
// equivalent Kiota objects. A missing extensible object is represented as {}.
func GraphToFrameworkNormalizedJSONObject(ctx context.Context, value serialization.Parsable, prior jsontypes.Normalized) (jsontypes.Normalized, error) {
	fields, err := GraphToJSONFields(value)
	if err != nil {
		return jsontypes.NewNormalizedNull(), err
	}
	encoded, err := json.Marshal(fields)
	if err != nil {
		return jsontypes.NewNormalizedNull(), fmt.Errorf("encode Graph JSON fields: %w", err)
	}
	return GraphToFrameworkNormalizedJSON(ctx, encoded, prior)
}
