package utilityEntraIdSidConverter

import (
	"testing"
)

func TestConvertSidToObjectId(t *testing.T) {
	testCases := []struct {
		name        string
		sid         string
		expected    string
		expectError bool
	}{
		{
			name:        "Valid SID conversion - Real world example",
			sid:         "S-1-12-1-1943430372-1249052806-2496021943-3034400218",
			expected:    "73d664e4-0886-4a73-b745-c694da45ddb4",
			expectError: false,
		},
		{
			name:        "Invalid SID format - missing parts",
			sid:         "S-1-12-1-1234567890",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid SID format - wrong prefix",
			sid:         "S-1-5-21-1234567890-1234567891-1234567892-1234567893",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid SID format - non-numeric RID",
			sid:         "S-1-12-1-abc-1234567891-1234567892-1234567893",
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertSidToObjectId(tc.sid)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Expected %s, got %s", tc.expected, result)
				}
			}
		})
	}
}

func TestConvertObjectIdToSid(t *testing.T) {
	testCases := []struct {
		name        string
		objectId    string
		expected    string
		expectError bool
	}{
		{
			name:        "Valid Object ID conversion - Real world example",
			objectId:    "73d664e4-0886-4a73-b745-c694da45ddb4",
			expected:    "S-1-12-1-1943430372-1249052806-2496021943-3034400218",
			expectError: false,
		},
		{
			name:        "Valid Object ID conversion - uppercase",
			objectId:    "73D664E4-0886-4A73-B745-C694DA45DDB4",
			expected:    "S-1-12-1-1943430372-1249052806-2496021943-3034400218",
			expectError: false,
		},
		{
			name:        "Invalid Object ID format - wrong format",
			objectId:    "invalid-guid",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Invalid Object ID format - missing dashes",
			objectId:    "73d664e408864a73b745c694da45ddb4",
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertObjectIdToSid(tc.objectId)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Expected %s, got %s", tc.expected, result)
				}
			}
		})
	}
}

func TestBidirectionalConversion(t *testing.T) {
	testCases := []struct {
		name string
		sid  string
	}{
		{
			name: "Round trip conversion - Real world example",
			sid:  "S-1-12-1-1943430372-1249052806-2496021943-3034400218",
		},
		{
			name: "Round trip conversion - Random values 1",
			sid:  "S-1-12-1-1000000000-2000000000-3000000000-4000000000",
		},
		{
			name: "Round trip conversion - Random values 2",
			sid:  "S-1-12-1-123456789-987654321-111111111-222222222",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			objectId, err := convertSidToObjectId(tc.sid)
			if err != nil {
				t.Fatalf("Failed to convert SID to Object ID: %v", err)
			}

			sidResult, err := convertObjectIdToSid(objectId)
			if err != nil {
				t.Fatalf("Failed to convert Object ID back to SID: %v", err)
			}

			if sidResult != tc.sid {
				t.Errorf("Round trip conversion failed. Original: %s, Result: %s, Object ID: %s", tc.sid, sidResult, objectId)
			}
		})
	}
}
