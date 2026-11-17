package graphBetaServicePrincipalSynchronizationSchemaJson

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/serviceprincipals"
)

// unsynchronizedMetadataKey is the object mapping metadata entry Graph sets from the mapping's enabled flag.
const unsynchronizedMetadataKey = "Unsynchronized"

var (
	errValueMustBeKnown        = errors.New("value must be known")
	errNotJSONArray            = errors.New("value must be a JSON array")
	errUnexpectedValueType     = errors.New("unexpected value type")
	errTrailingJSON            = errors.New("unexpected content after the JSON value")
	errEmptySchemaResponse     = errors.New("graph returned an empty synchronization schema response")
	errNoSynchronizationSchema = errors.New("graph returned no synchronization schema")
)

var graphErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

var (
	_ basetypes.StringTypable                    = SchemaJSONType{}
	_ basetypes.StringValuableWithSemanticEquals = SchemaJSONValue{}
)

// SchemaJSONType is a string type holding a synchronization schema collection as JSON.
type SchemaJSONType struct {
	basetypes.StringType
}

// Equal returns true if the given type is equivalent.
func (t SchemaJSONType) Equal(o attr.Type) bool {
	other, ok := o.(SchemaJSONType)
	if !ok {
		return false
	}
	return t.StringType.Equal(other.StringType)
}

// String returns a human readable string of the type name.
func (t SchemaJSONType) String() string {
	return "SchemaJSONType"
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t SchemaJSONType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return SchemaJSONValue{StringValue: in}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.
func (t SchemaJSONType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("convert synchronization schema JSON value: %w", err)
	}
	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("%w %T", errUnexpectedValueType, attrValue)
	}
	return SchemaJSONValue{StringValue: stringValue}, nil
}

// ValueType returns the Value type.
func (t SchemaJSONType) ValueType(ctx context.Context) attr.Value {
	return SchemaJSONValue{}
}

// SchemaJSONValue is a synchronization schema collection as JSON. Values are semantically equal when they hold the
// same JSON after the properties Graph manages itself are removed (see normaliseSchemaJSON).
type SchemaJSONValue struct {
	basetypes.StringValue
}

// NewSchemaJSONValue returns a known SchemaJSONValue.
func NewSchemaJSONValue(value string) SchemaJSONValue {
	return SchemaJSONValue{StringValue: basetypes.NewStringValue(value)}
}

// NewSchemaJSONNull returns a null SchemaJSONValue.
func NewSchemaJSONNull() SchemaJSONValue {
	return SchemaJSONValue{StringValue: basetypes.NewStringNull()}
}

// Equal returns true if the given value is equivalent.
func (v SchemaJSONValue) Equal(o attr.Value) bool {
	other, ok := o.(SchemaJSONValue)
	if !ok {
		return false
	}
	return v.StringValue.Equal(other.StringValue)
}

// Type returns the value's type.
func (v SchemaJSONValue) Type(ctx context.Context) attr.Type {
	return SchemaJSONType{}
}

// StringSemanticEquals keeps the prior value when the new value only differs in formatting or Graph managed properties.
func (v SchemaJSONValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(SchemaJSONValue)
	if !ok {
		diags.AddError("Semantic Equality Check Error", fmt.Sprintf("Expected value type %T but got value type %T.", v, newValuable))
		return false, diags
	}

	if v.IsNull() || v.IsUnknown() || newValue.IsNull() || newValue.IsUnknown() {
		return false, diags
	}

	return schemaJSONSemanticallyEqual(v.ValueString(), newValue.ValueString()), diags
}

// schemaJSONSemanticallyEqual compares two schema collections after normalisation. Invalid JSON is never equal.
func schemaJSONSemanticallyEqual(a, b string) bool {
	left, err := normaliseSchemaJSON(a)
	if err != nil {
		return false
	}
	right, err := normaliseSchemaJSON(b)
	if err != nil {
		return false
	}
	return reflect.DeepEqual(left, right)
}

// normaliseSchemaJSON decodes a schema collection and removes what Graph manages itself:
//   - `version` properties, which Graph sets to a new value on every schema change;
//   - OData annotations such as `@odata.type`;
//   - `Unsynchronized` metadata entries, which Graph derives from an object mapping's `enabled` flag.
//
// Array order is kept, because the order of expression parameters is significant.
func normaliseSchemaJSON(value string) (any, error) {
	decoded, err := decodeJSON([]byte(value))
	if err != nil {
		return nil, err
	}
	return normaliseValue(decoded), nil
}

func normaliseValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		normalised := make(map[string]any, len(typed))
		for key, nested := range typed {
			if key == "version" || strings.HasPrefix(key, "@odata.") {
				continue
			}
			normalised[key] = normaliseValue(nested)
		}
		return normalised
	case []any:
		normalised := make([]any, 0, len(typed))
		for _, item := range typed {
			if entry, ok := item.(map[string]any); ok && entry["key"] == unsynchronizedMetadataKey && len(entry) == 2 {
				if _, hasValue := entry["value"]; hasValue {
					continue
				}
			}
			normalised = append(normalised, normaliseValue(item))
		}
		return normalised
	default:
		return value
	}
}

// decodeJSON decodes JSON keeping numbers as json.Number so large or decimal values round-trip unchanged.
func decodeJSON(content []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.UseNumber()

	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return nil, fmt.Errorf("decode JSON: %w", err)
	}
	if decoder.More() {
		return nil, fmt.Errorf("decode JSON: %w", errTrailingJSON)
	}
	return decoded, nil
}

// canonicalJSONArray returns a JSON array in the form Terraform's jsonencode() produces: compact, with object keys
// sorted and <, > and & escaped. A null or missing collection is returned as an empty array.
func canonicalJSONArray(content json.RawMessage) (string, error) {
	if len(content) == 0 || string(content) == "null" {
		return "[]", nil
	}
	decoded, err := decodeJSON(content)
	if err != nil {
		return "", err
	}
	if _, ok := decoded.([]any); !ok {
		return "", fmt.Errorf("%w, got %T", errNotJSONArray, decoded)
	}
	encoded, err := json.Marshal(decoded)
	if err != nil {
		return "", fmt.Errorf("encode JSON: %w", err)
	}
	return string(encoded), nil
}

// schemaRequestBuilder returns the SDK request builder for the synchronization job schema.
func (r *ServicePrincipalSynchronizationSchemaJsonResource) schemaRequestBuilder(
	servicePrincipalId string,
	jobId string,
) *serviceprincipals.ItemSynchronizationJobsItemSchemaRequestBuilder {
	return r.client.ServicePrincipals().ByServicePrincipalId(servicePrincipalId).Synchronization().Jobs().BySynchronizationJobId(jobId).Schema()
}

// getSchema reads the synchronization job schema as raw JSON.
//
// The SDK model would drop properties it does not know, so the response is read as bytes through the SDK adapter.
// Authentication, middleware and OData error mapping still go through the adapter.
func (r *ServicePrincipalSynchronizationSchemaJsonResource) getSchema(ctx context.Context, servicePrincipalId, jobId string) (*synchronizationSchemaResponse, error) {
	requestInfo, err := r.schemaRequestBuilder(servicePrincipalId, jobId).ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("build synchronization schema GET request: %w", err)
	}

	result, err := r.client.GetAdapter().SendPrimitive(ctx, requestInfo, "[]byte", graphErrorMapping)
	if err != nil {
		return nil, err //nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
	}

	content, ok := result.([]byte)
	if !ok || len(content) == 0 {
		return nil, errEmptySchemaResponse
	}

	var response synchronizationSchemaResponse
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, fmt.Errorf("decode synchronization schema response: %w", err)
	}
	return &response, nil
}

// putSchema replaces the synchronization job schema.
//
// The SDK request builder exposes PATCH, but the documented update method is PUT, so the PUT reuses the builder's
// request information with the method changed and the JSON body sent unchanged.
func (r *ServicePrincipalSynchronizationSchemaJsonResource) putSchema(ctx context.Context, servicePrincipalId, jobId string, requestBody []byte) error {
	requestInfo, err := r.schemaRequestBuilder(servicePrincipalId, jobId).ToGetRequestInformation(ctx, nil)
	if err != nil {
		return fmt.Errorf("build synchronization schema PUT request: %w", err)
	}

	requestInfo.Method = abstractions.PUT
	requestInfo.Content = requestBody
	requestInfo.Headers.TryAdd("Content-Type", "application/json")

	return r.client.GetAdapter().SendNoContent(ctx, requestInfo, graphErrorMapping) //nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
}

// resourceId returns the Terraform ID for a service principal and synchronization job.
func resourceId(servicePrincipalId, jobId string) string {
	return servicePrincipalId + "/" + jobId
}
