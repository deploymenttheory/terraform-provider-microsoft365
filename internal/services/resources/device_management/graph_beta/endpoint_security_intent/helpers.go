package graphBetaEndpointSecurityIntent

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
)

const (
	settingTypeBoolean    = "boolean"
	settingTypeInteger    = "integer"
	settingTypeString     = "string"
	settingTypeComplex    = "complex"
	settingTypeCollection = "collection"

	noFilterId = "00000000-0000-0000-0000-000000000000"
)

var guidPattern = regexp.MustCompile(constants.GuidRegex)

var (
	errInvalidSettings    = errors.New("invalid settings")
	errInvalidAssignments = errors.New("invalid assignments")
	errInvalidImportId    = errors.New("invalid import ID")
)

// settingAttrTypes are the attribute types of one element of the settings set.
var settingAttrTypes = map[string]attr.Type{
	"definition_id": types.StringType,
	"setting_type":  types.StringType,
	"value_json":    types.StringType,
}

// assignmentAttrTypes are the attribute types of one element of the assignments set.
var assignmentAttrTypes = map[string]attr.Type{
	"type":        types.StringType,
	"group_id":    types.StringType,
	"filter_id":   types.StringType,
	"filter_type": types.StringType,
}

// newSettingInstance returns the Graph setting instance for a setting type.
// Graph's updateSettings rejects the abstract deviceManagementSettingInstance type with a 500, so the concrete type is always sent.
func newSettingInstance(settingType string) (graphmodels.DeviceManagementSettingInstanceable, error) {
	switch settingType {
	case settingTypeBoolean:
		return graphmodels.NewDeviceManagementBooleanSettingInstance(), nil
	case settingTypeInteger:
		return graphmodels.NewDeviceManagementIntegerSettingInstance(), nil
	case settingTypeString:
		return graphmodels.NewDeviceManagementStringSettingInstance(), nil
	case settingTypeComplex:
		return graphmodels.NewDeviceManagementComplexSettingInstance(), nil
	case settingTypeCollection:
		return graphmodels.NewDeviceManagementCollectionSettingInstance(), nil
	default:
		return nil, fmt.Errorf("%w: unsupported setting type %q", errInvalidSettings, settingType)
	}
}

// settingTypeOf returns the setting type of a Graph setting instance, or an empty string for types the resource cannot write.
func settingTypeOf(instance graphmodels.DeviceManagementSettingInstanceable) string {
	switch instance.(type) {
	case *graphmodels.DeviceManagementBooleanSettingInstance:
		return settingTypeBoolean
	case *graphmodels.DeviceManagementIntegerSettingInstance:
		return settingTypeInteger
	case *graphmodels.DeviceManagementStringSettingInstance:
		return settingTypeString
	case *graphmodels.DeviceManagementComplexSettingInstance:
		return settingTypeComplex
	case *graphmodels.DeviceManagementCollectionSettingInstance:
		return settingTypeCollection
	default:
		return ""
	}
}

// jsonEqual reports whether two JSON documents hold the same value, ignoring formatting and key order.
func jsonEqual(a, b string) bool {
	if a == b {
		return true
	}

	var left, right any
	if err := json.Unmarshal([]byte(a), &left); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &right); err != nil {
		return false
	}

	return reflect.DeepEqual(left, right)
}

// parseImportId splits an import ID of the form `<intent_id>` or `<intent_id>/<definition_id>,<definition_id>`.
func parseImportId(importId string) (string, []string, error) {
	intentId, definitionList, hasSettings := strings.Cut(importId, "/")
	if !guidPattern.MatchString(intentId) {
		return "", nil, fmt.Errorf("%w: expected '<intent_id>' or '<intent_id>/<definition_id>,<definition_id>' with a GUID intent_id, got %q", errInvalidImportId, importId)
	}

	if !hasSettings {
		return intentId, nil, nil
	}

	definitionIds := make([]string, 0)
	seen := make(map[string]bool)
	for _, definitionId := range strings.Split(definitionList, ",") {
		definitionId = strings.TrimSpace(definitionId)
		if definitionId == "" || seen[definitionId] {
			continue
		}
		seen[definitionId] = true
		definitionIds = append(definitionIds, definitionId)
	}

	if len(definitionIds) == 0 {
		return "", nil, fmt.Errorf("%w: no setting definition IDs found after '/' in %q", errInvalidImportId, importId)
	}

	return intentId, definitionIds, nil
}

// importedSettings builds placeholder settings for an import. Read replaces the type and value of each one with what Graph returns.
func importedSettings(definitionIds []string) types.Set {
	elements := make([]attr.Value, 0, len(definitionIds))
	for _, definitionId := range definitionIds {
		elements = append(elements, types.ObjectValueMust(settingAttrTypes, map[string]attr.Value{
			"definition_id": types.StringValue(definitionId),
			"setting_type":  types.StringValue(""),
			"value_json":    types.StringValue(""),
		}))
	}

	return types.SetValueMust(types.ObjectType{AttrTypes: settingAttrTypes}, elements)
}
