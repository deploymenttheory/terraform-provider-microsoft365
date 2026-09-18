package graphBetaNetworkCustomBlockPage

import (
	"context"
	"errors"
	"strings"
	"unicode/utf16"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	errInvalidState = errors.New("invalid custom block page state")
	errUnknownBody  = errors.New(
		"configuration.body must be known when applying the configuration",
	)
	errMissingState             = errors.New("custom block page response is missing its state")
	errUnsupportedConfiguration = errors.New("unsupported custom block page configuration type")
	errEmptyBody                = errors.New("markdown body cannot be empty or whitespace")
	errBodyTooLong              = errors.New(
		"markdown body length exceeds maximum allowed length of 1024 UTF-16 code units",
	)
)

// The live API counts UTF-16 code units: 512 emoji succeed, 513 fail.
type blockMessageValidator struct{}

var _ validator.String = blockMessageValidator{}

func (blockMessageValidator) Description(context.Context) string {
	return "must contain non-whitespace text and be 1–1024 UTF-16 code units long"
}

func (v blockMessageValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (blockMessageValidator) ValidateString(
	_ context.Context,
	req validator.StringRequest,
	resp *validator.StringResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if err := validateBlockMessage(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid block page body", err.Error())
	}
}

func validateBlockMessage(body string) error {
	if strings.TrimSpace(body) == "" {
		return errEmptyBody
	}
	if len(utf16.Encode([]rune(body))) > 1024 {
		return errBodyTooLong
	}
	return nil
}
