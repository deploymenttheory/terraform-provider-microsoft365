package crud

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
)

// ImportIDFormat is the format an import ID, or one part of a composite import ID, must match.
//
// Resources whose ID is written into a custom_requests URL validate it on import, because
// custom_requests embeds the ID in the URL without escaping it. A malformed ID such as
// `abc?$select=x` or `../foo` would otherwise change the request path or query.
type ImportIDFormat struct {
	// Pattern is the anchored regular expression the ID must match.
	Pattern *regexp.Regexp
	// Description names the expected format in the error diagnostic, for example
	// "a GUID such as 00000000-0000-0000-0000-000000000000".
	Description string
}

// ImportIDFormatGUID is the format for resources whose Graph ID is a GUID.
var ImportIDFormatGUID = ImportIDFormat{
	Pattern:     regexp.MustCompile(constants.GuidRegex),
	Description: "a GUID such as 00000000-0000-0000-0000-000000000000",
}

// ImportIDFormatEnrollmentNotificationConfiguration is the format Intune uses for device enrollment
// notification configuration IDs: a GUID followed by "_EnrollmentNotificationsConfiguration".
var ImportIDFormatEnrollmentNotificationConfiguration = ImportIDFormat{
	Pattern:     regexp.MustCompile(constants.EnrollmentNotificationConfigurationIdRegex),
	Description: "a GUID followed by _EnrollmentNotificationsConfiguration, such as 00000000-0000-0000-0000-000000000000_EnrollmentNotificationsConfiguration",
}

// ValidateImportIDPart checks that value matches format. name identifies the value in the
// diagnostic, for example "import ID" or "agent_identity_blueprint_id". On a mismatch it adds
// an error diagnostic naming the expected format and returns false.
func ValidateImportIDPart(diags *diag.Diagnostics, name, value string, format ImportIDFormat) bool {
	if format.Pattern.MatchString(value) {
		return true
	}

	diags.AddError(
		"Invalid Import ID",
		fmt.Sprintf("The %s %q is not valid. It must be %s.", name, value, format.Description),
	)

	return false
}

// ImportStatePassthroughValidatedID validates req.ID against format and, if it matches,
// imports it exactly as resource.ImportStatePassthroughID does for the "id" attribute.
// If it doesn't match, it adds an error diagnostic naming the expected format and leaves
// the state unset.
func ImportStatePassthroughValidatedID(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse, format ImportIDFormat) {
	if !ValidateImportIDPart(&resp.Diagnostics, "import ID", req.ID, format) {
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
