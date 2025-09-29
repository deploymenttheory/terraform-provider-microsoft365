package constructors

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

// DebugLogGraphObject is a helper function to serialize and debug log Microsoft Graph objects that
// implement serialization.Parsable.
// This function takes a Microsoft Graph object (like a policy, device configuration, etc.) and converts it
// into a human-readable JSON format before it gets sent to Microsoft's API. It's like taking a snapshot of
// what we're about to send.
//
// Parameters:
//   - ctx: The context for logging
//   - message: A descriptive message that will prefix the JSON in the logs
//   - object: Any Microsoft Graph object that implements serialization.Parsable
//
// Returns:
//   - error: Any error encountered during serialization or logging
//
// Usage example:
//
//	if err := debugLogGraphObject(ctx, "Final JSON to be sent to Graph API", profile); err != nil {
//	    tflog.Error(ctx, "Failed to debug log object", map[string]any{
//	        "error": err.Error(),
//	    })
//	}
func DebugLogGraphObject(ctx context.Context, message string, object serialization.Parsable) error {
	factory := jsonserialization.NewJsonSerializationWriterFactory()
	writer, err := factory.GetSerializationWriter("application/json")
	if err != nil {
		return err
	}

	err = writer.WriteObjectValue("", object)
	if err != nil {
		return err
	}

	jsonBytes, err := writer.GetSerializedContent()
	if err != nil {
		return err
	}

	var rawJSON interface{}
	if err := json.Unmarshal(jsonBytes, &rawJSON); err != nil {
		return err
	}

	debugJSON, err := json.MarshalIndent(rawJSON, "", "    ")
	if err != nil {
		return err
	}

	tflog.Debug(ctx, message, map[string]any{
		"json": "\n" + string(debugJSON),
	})

	return nil
}
