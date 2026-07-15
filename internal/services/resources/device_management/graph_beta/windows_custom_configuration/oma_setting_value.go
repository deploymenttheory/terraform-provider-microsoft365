package graphBetaWindowsCustomConfiguration

import (
	"fmt"
	"strconv"
	"time"
)

// parsedOmaSettingValue holds the typed representation of an OMA setting value together with
// the canonical string form that mapOmaSetting produces when reading the value back from the
// Graph API. Plan-time validation requires configured values to already be in canonical form,
// otherwise the post-apply Read would rewrite the value and Terraform would fail with
// "Provider produced inconsistent result after apply".
type parsedOmaSettingValue struct {
	canonical  string
	intValue   int32
	boolValue  bool
	floatValue float32
	timeValue  time.Time
}

// parseOmaSettingValue parses and canonicalizes an OMA setting value according to its OData type.
// It is the single source of truth for value parsing, shared by ModifyPlan (validation) and
// constructOmaSetting (request construction).
func parseOmaSettingValue(odataType, value string) (parsedOmaSettingValue, error) {
	switch odataType {
	case "#microsoft.graph.omaSettingString", "#microsoft.graph.omaSettingBase64", "#microsoft.graph.omaSettingStringXml":
		return parsedOmaSettingValue{canonical: value}, nil

	case "#microsoft.graph.omaSettingInteger":
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return parsedOmaSettingValue{}, fmt.Errorf("value %q is not a valid integer", value)
		}
		return parsedOmaSettingValue{
			canonical: strconv.FormatInt(intValue, 10),
			intValue:  int32(intValue),
		}, nil

	case "#microsoft.graph.omaSettingBoolean":
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return parsedOmaSettingValue{}, fmt.Errorf("value %q is not a valid boolean", value)
		}
		return parsedOmaSettingValue{
			canonical: strconv.FormatBool(boolValue),
			boolValue: boolValue,
		}, nil

	case "#microsoft.graph.omaSettingDateTime":
		timeValue, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return parsedOmaSettingValue{}, fmt.Errorf("value %q is not a valid RFC3339 timestamp (e.g. 2024-01-01T00:00:00Z)", value)
		}
		return parsedOmaSettingValue{
			canonical: timeValue.UTC().Format(time.RFC3339),
			timeValue: timeValue,
		}, nil

	case "#microsoft.graph.omaSettingFloatingPoint":
		floatValue, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return parsedOmaSettingValue{}, fmt.Errorf("value %q is not a valid floating point number", value)
		}
		return parsedOmaSettingValue{
			canonical:  strconv.FormatFloat(floatValue, 'f', -1, 32),
			floatValue: float32(floatValue),
		}, nil

	default:
		return parsedOmaSettingValue{}, fmt.Errorf("unsupported oma setting odata type: %s", odataType)
	}
}
