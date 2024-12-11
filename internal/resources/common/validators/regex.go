package validators

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// RegexValidator provides string validation against a regular expression pattern,
// skipping validation for null values. It implements the validator.String interface
// for use with Terraform Framework attributes.
type regexValidator struct {
	regexp  *regexp.Regexp
	message string
}

func (v regexValidator) Description(ctx context.Context) string {
	if v.message != "" {
		return v.message
	}
	return fmt.Sprintf("value must match regular expression '%s'", v.regexp)
}

func (v regexValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v regexValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() {
		return
	}

	if request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if !v.regexp.MatchString(value) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			value,
		))
	}
}

// RegexMatches returns a validator.String that ensures any configured value
// matches the provided regular expression pattern. Null and unknown values are skipped.
// An optional message can be provided to customize the error output.
func RegexMatches(regexp *regexp.Regexp, message string) validator.String {
	return regexValidator{
		regexp:  regexp,
		message: message,
	}
}
