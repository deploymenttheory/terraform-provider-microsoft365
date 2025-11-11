package entra_id_sid_converter

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type sidRidRangeValidator struct{}

func (v sidRidRangeValidator) Description(ctx context.Context) string {
	return "Validates that all RID components in an Entra ID SID are within the valid uint32 range (0 to 4,294,967,295)"
}

func (v sidRidRangeValidator) MarkdownDescription(ctx context.Context) string {
	return "Validates that all RID components in an Entra ID SID are within the valid uint32 range (0 to 4,294,967,295)"
}

func (v sidRidRangeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	sidValue := req.ConfigValue.ValueString()

	if !strings.HasPrefix(sidValue, "S-1-12-1-") {
		return
	}

	parts := strings.Split(sidValue, "-")
	if len(parts) != 8 {
		return
	}

	const maxUint32 = uint64(4294967295)

	for i := 4; i < 8; i++ {
		ridValue, err := strconv.ParseUint(parts[i], 10, 64)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid RID Component",
				fmt.Sprintf("RID component at position %d ('%s') is not a valid number: %s", i-3, parts[i], err),
			)
			return
		}

		if ridValue > maxUint32 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"RID Component Out of Range",
				fmt.Sprintf("RID component at position %d has value %d, which exceeds the maximum uint32 value of 4,294,967,295. Each RID component must fit within a 32-bit unsigned integer.", i-3, ridValue),
			)
			return
		}
	}
}

func ValidateSidRidRange() validator.String {
	return sidRidRangeValidator{}
}
