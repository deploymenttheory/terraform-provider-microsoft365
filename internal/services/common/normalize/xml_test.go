package normalize

import (
	"testing"
)

func TestNormalizeXML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Remove BOM",
			input:    "\ufeff<?xml version=\"1.0\"?>",
			expected: "<?xml version=\"1.0\"?>",
		},
		{
			name:     "Unescape HTML entities",
			input:    "&lt;SiPolicy&gt;",
			expected: "<SiPolicy>",
		},
		{
			name:     "Both BOM and HTML entities",
			input:    "\ufeff&lt;SiPolicy&gt;",
			expected: "<SiPolicy>",
		},
		{
			name:     "No changes needed",
			input:    "<?xml version=\"1.0\"?><SiPolicy>",
			expected: "<?xml version=\"1.0\"?><SiPolicy>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := NormalizeXML(test.input)
			if result != test.expected {
				t.Errorf("Expected %q but got %q", test.expected, result)
			}
		})
	}
}

func TestReverseNormalizeXML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Add BOM to Windows line endings",
			input:    "<?xml version=\"1.0\"?>\r\n<SiPolicy>",
			expected: "\ufeff<?xml version=\"1.0\"?>\r\n<SiPolicy>",
		},
		{
			name:     "Don't add BOM to inline XML",
			input:    "<?xml version=\"1.0\"?><SiPolicy>",
			expected: "<?xml version=\"1.0\"?><SiPolicy>",
		},
		{
			name:     "Don't duplicate BOM",
			input:    "\ufeff<?xml version=\"1.0\"?>",
			expected: "<?xml version=\"1.0\"?>", // BOM is removed in NormalizeXML and not added back because there are no CRLF
		},
		{
			name:     "Normalize and add BOM",
			input:    "&lt;SiPolicy&gt;\r\n",
			expected: "\ufeff<SiPolicy>\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ReverseNormalizeXML(test.input)
			if result != test.expected {
				t.Errorf("Expected %q but got %q", test.expected, result)
			}
		})
	}
}

func TestLikelyFromFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Windows line endings",
			input:    "<?xml version=\"1.0\"?>\r\n<SiPolicy>",
			expected: true,
		},
		{
			name:     "Has BOM",
			input:    "\ufeff<?xml version=\"1.0\"?>",
			expected: true,
		},
		{
			name:     "Inline XML",
			input:    "<?xml version=\"1.0\"?><SiPolicy>",
			expected: false,
		},
		{
			name:     "Unix line endings",
			input:    "<?xml version=\"1.0\"?>\n<SiPolicy>",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := LikelyFromFile(test.input)
			if result != test.expected {
				t.Errorf("Expected %v but got %v", test.expected, result)
			}
		})
	}
}

func TestNormalizeAndReverseRoundTrip(t *testing.T) {
	// Test that normalizing and then reverse normalizing preserves content appropriately
	testCases := []struct {
		name          string
		input         string
		expectSame    bool
		expectedFinal string
	}{
		{
			name:       "Inline XML round trip",
			input:      "<?xml version=\"1.0\"?><SiPolicy>",
			expectSame: true,
		},
		{
			name:          "File XML with BOM",
			input:         "\ufeff<?xml version=\"1.0\"?>\r\n<SiPolicy>",
			expectSame:    false, // BOM is removed and then added back
			expectedFinal: "\ufeff<?xml version=\"1.0\"?>\r\n<SiPolicy>",
		},
		{
			name:          "File XML without BOM",
			input:         "<?xml version=\"1.0\"?>\r\n<SiPolicy>",
			expectSame:    false,
			expectedFinal: "\ufeff<?xml version=\"1.0\"?>\r\n<SiPolicy>",
		},
		{
			name:          "HTML entities",
			input:         "&lt;SiPolicy&gt;",
			expectSame:    false,
			expectedFinal: "<SiPolicy>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			normalized := NormalizeXML(tc.input)
			reversed := ReverseNormalizeXML(tc.input)

			if tc.expectSame && reversed != tc.input {
				t.Errorf("Expected round trip to preserve content exactly, but got different results\nInput: %q\nOutput: %q", tc.input, reversed)
			}

			if !tc.expectSame && tc.expectedFinal != "" && reversed != tc.expectedFinal {
				t.Errorf("Expected specific transformation\nInput: %q\nExpected: %q\nGot: %q", tc.input, tc.expectedFinal, reversed)
			}

			// Test that normalizing twice is idempotent
			doubleNormalized := NormalizeXML(normalized)
			if doubleNormalized != normalized {
				t.Errorf("Normalizing twice should be idempotent\nFirst: %q\nSecond: %q", normalized, doubleNormalized)
			}
		})
	}
}

// TestEdgeCases tests edge cases like empty strings and very large XML
func TestEdgeCases(t *testing.T) {
	// Empty string
	if NormalizeXML("") != "" {
		t.Error("NormalizeXML should return empty string for empty input")
	}
	if ReverseNormalizeXML("") != "" {
		t.Error("ReverseNormalizeXML should return empty string for empty input")
	}
	if LikelyFromFile("") {
		t.Error("Empty string should not be detected as from file")
	}

	// Very large XML (simulated)
	largeXML := "<?xml version=\"1.0\"?>"
	for i := 0; i < 1000; i++ {
		largeXML += "<element>content</element>"
	}

	normalizedLarge := NormalizeXML(largeXML)
	if normalizedLarge != largeXML {
		t.Error("Large XML should be normalized correctly")
	}

	reversedLarge := ReverseNormalizeXML(largeXML)
	if reversedLarge != largeXML {
		t.Error("Large XML should be reverse normalized correctly")
	}
}
