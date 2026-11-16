package graphBetaServicePrincipalSynchronizationSchemaJson

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// synchronizationSchemaRequest is the PUT body of a synchronizationSchema. Only the writable collections are sent.
type synchronizationSchemaRequest struct {
	SynchronizationRules json.RawMessage `json:"synchronizationRules"`
	Directories          json.RawMessage `json:"directories"`
}

// constructResource builds the PUT request body from the configured JSON collections.
func constructResource(ctx context.Context, data *ServicePrincipalSynchronizationSchemaJsonResourceModel) ([]byte, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	synchronizationRules, err := constructJSONArray("synchronization_rules_json", data.SynchronizationRulesJson)
	if err != nil {
		return nil, err
	}

	directories, err := constructJSONArray("directories_json", data.DirectoriesJson)
	if err != nil {
		return nil, err
	}

	requestBody, err := json.Marshal(synchronizationSchemaRequest{SynchronizationRules: synchronizationRules, Directories: directories})
	if err != nil {
		return nil, fmt.Errorf("failed to encode synchronization schema request: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Constructed %s request body of %d bytes", ResourceName, len(requestBody)))

	return requestBody, nil
}

// constructJSONArray checks that a configured collection is a JSON array and returns it unchanged.
func constructJSONArray(attribute string, value SchemaJSONValue) (json.RawMessage, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, fmt.Errorf("%s: %w", attribute, errValueMustBeKnown)
	}

	decoded, err := decodeJSON([]byte(value.ValueString()))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", attribute, err)
	}
	if _, ok := decoded.([]any); !ok {
		return nil, fmt.Errorf("%s: %w", attribute, errNotJSONArray)
	}

	return json.RawMessage(value.ValueString()), nil
}
