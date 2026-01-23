package attribute

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.List = int64ListSumEqualsValidator{}

// int64ListSumEqualsValidator validates that the sum of int64 values in a list equals a specific value.
// It can optionally emit a warning instead of an error when the sum doesn't match.
type int64ListSumEqualsValidator struct {
	expectedSum int64
	warningOnly bool
}

// Description provides a human-readable description of the validator's purpose.
func (v int64ListSumEqualsValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

// MarkdownDescription provides a Markdown description of the validator's purpose.
func (v int64ListSumEqualsValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("sum of values must equal %d", v.expectedSum)
}

// ValidateList performs the validation logic on the list attribute.
func (v int64ListSumEqualsValidator) ValidateList(ctx context.Context, request validator.ListRequest, response *validator.ListResponse) {
	// Skip validation if list is null or unknown
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	// Calculate the sum
	var sum int64
	for _, elem := range request.ConfigValue.Elements() {
		int64Elem, ok := elem.(types.Int64)
		if !ok {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid List Element Type",
				"Expected all elements to be int64 values.",
			)
			return
		}

		// Skip unknown or null values
		if int64Elem.IsNull() || int64Elem.IsUnknown() {
			return
		}

		sum += int64Elem.ValueInt64()
	}

	// Check if sum matches expected value
	if sum != v.expectedSum {
		message := fmt.Sprintf("Sum of values is %d but expected %d. Distribution may not be exact.", sum, v.expectedSum)
		
		if v.warningOnly {
			response.Diagnostics.Append(diag.NewAttributeWarningDiagnostic(
				request.Path,
				"List Sum Warning",
				message,
			))
		} else {
			response.Diagnostics.AddAttributeError(
				request.Path,
				"Invalid List Sum",
				message,
			)
		}
	}
}

// Int64ListSumEquals returns a validator that ensures the sum of int64 values in a list equals the expected value.
// If the sum doesn't match, it returns an error.
func Int64ListSumEquals(expectedSum int64) validator.List {
	return int64ListSumEqualsValidator{
		expectedSum: expectedSum,
		warningOnly: false,
	}
}

// Int64ListSumEqualsWarning returns a validator that ensures the sum of int64 values in a list equals the expected value.
// If the sum doesn't match, it returns a warning instead of an error.
func Int64ListSumEqualsWarning(expectedSum int64) validator.List {
	return int64ListSumEqualsValidator{
		expectedSum: expectedSum,
		warningOnly: true,
	}
}
