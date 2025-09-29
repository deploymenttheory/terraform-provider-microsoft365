package helpers

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// DecodeBase64ToString decodes a base64-encoded string and returns a Terraform Framework string.
// If decoding fails, it logs a warning and returns the original string as fallback.
func DecodeBase64ToString(ctx context.Context, encoded string) types.String {
	decodedContent, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		tflog.Warn(ctx, "Failed to decode base64 content", map[string]any{
			"error": err.Error(),
		})
		return types.StringValue(encoded)
	}
	return types.StringValue(string(decodedContent))
}

// ByteStringToBase64 converts a byte slice to a base64-encoded string
func ByteStringToBase64(data []byte) string {
	if data == nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(data)
}

// StringToBase64 encodes the input string to base64.
func StringToBase64(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("input string is empty")
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	return encoded, nil
}
