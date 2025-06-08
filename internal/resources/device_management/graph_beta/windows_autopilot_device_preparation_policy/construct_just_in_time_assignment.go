package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructJustInTimeAssignmentBody constructs the request body for assigning a just-in-time configuration
func constructJustInTimeAssignmentBody(ctx context.Context, deviceSecurityGroupID string) (serialization.Parsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing just-in-time assignment with security group: %s", deviceSecurityGroupID))

	// Create a custom entity to represent the request body
	entity := models.NewEntity()

	// Define the JSON structure manually
	requestBodyMap := map[string]interface{}{
		"justInTimeAssignments": map[string]interface{}{
			"targetType": "entraSecurityGroup",
			"target": []string{
				deviceSecurityGroupID,
			},
		},
	}

	// Convert the map to JSON and add to the entity's additionalData
	jsonData, err := json.Marshal(requestBodyMap)
	if err != nil {
		return nil, fmt.Errorf("error marshaling just-in-time assignment body: %v", err)
	}

	// Parse the JSON into a map
	var rawData map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawData); err != nil {
		return nil, fmt.Errorf("error unmarshaling just-in-time assignment body: %v", err)
	}

	// Add the raw data to the entity's additionalData
	additionalData := entity.GetAdditionalData()
	for k, v := range rawData {
		additionalData[k] = v
	}
	entity.SetAdditionalData(additionalData)

	// Log the final JSON payload
	if err := constructors.DebugLogGraphObject(ctx, "Final JSON to be sent to Graph API for just-in-time assignment", entity); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, "Finished constructing just-in-time assignment body")

	return entity, nil
}
