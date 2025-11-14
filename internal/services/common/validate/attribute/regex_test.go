package attribute

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexValidator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		validator   validator.String
		value       types.String
		expectError bool
	}{
		{
			name: "valid uuid",
			validator: RegexMatches(
				regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`),
				"must be a valid UUID",
			),
			value:       types.StringValue("550e8400-e29b-41d4-a716-446655440000"),
			expectError: false,
		},
		{
			name: "invalid uuid",
			validator: RegexMatches(
				regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`),
				"must be a valid UUID",
			),
			value:       types.StringValue("invalid-uuid"),
			expectError: true,
		},
		{
			name: "null value",
			validator: RegexMatches(
				regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`),
				"must be a valid UUID",
			),
			value:       types.StringNull(),
			expectError: false,
		},
		{
			name: "unknown value",
			validator: RegexMatches(
				regexp.MustCompile(`^[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}$`),
				"must be a valid UUID",
			),
			value:       types.StringUnknown(),
			expectError: false,
		},
		{
			name: "custom error message",
			validator: RegexMatches(
				regexp.MustCompile(`^test\d+$`),
				"value must start with 'test' followed by numbers",
			),
			value:       types.StringValue("invalid"),
			expectError: true,
		},
		{
			name: "empty string",
			validator: RegexMatches(
				regexp.MustCompile(`^test\d+$`),
				"value must start with 'test' followed by numbers",
			),
			value:       types.StringValue(""),
			expectError: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				ConfigValue: test.value,
			}
			response := validator.StringResponse{}

			test.validator.ValidateString(context.Background(), request, &response)

			if test.expectError {
				assert.NotEmpty(t, response.Diagnostics)
			} else {
				assert.Empty(t, response.Diagnostics)
			}
		})
	}
}

func TestRegexValidatorDescription(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		regexp           *regexp.Regexp
		message          string
		wantDesc         string
		wantMarkdownDesc string
	}{
		{
			name:             "with custom message",
			regexp:           regexp.MustCompile(`^test\d+$`),
			message:          "custom validation message",
			wantDesc:         "custom validation message",
			wantMarkdownDesc: "custom validation message",
		},
		{
			name:             "without custom message",
			regexp:           regexp.MustCompile(`^test\d+$`),
			message:          "",
			wantDesc:         "value must match regular expression '^test\\d+$'",
			wantMarkdownDesc: "value must match regular expression '^test\\d+$'",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			validator := RegexMatches(tc.regexp, tc.message)

			gotDesc := validator.Description(context.Background())
			gotMarkdownDesc := validator.MarkdownDescription(context.Background())

			assert.Equal(t, tc.wantDesc, gotDesc)
			assert.Equal(t, tc.wantMarkdownDesc, gotMarkdownDesc)
		})
	}
}

func TestRegexMatches(t *testing.T) {
	t.Parallel()

	testRegex := regexp.MustCompile(`^test\d+$`)
	testMessage := "test message"

	validator := RegexMatches(testRegex, testMessage)
	require.NotNil(t, validator)

	regexValidator, ok := validator.(regexValidator)
	require.True(t, ok)

	assert.Equal(t, testRegex, regexValidator.regexp)
	assert.Equal(t, testMessage, regexValidator.message)
}
