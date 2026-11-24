package sample

import "context"

// widgetNameValidator is a schema validator.
type widgetNameValidator struct{}

// Description is a framework validator method.
func (v widgetNameValidator) Description(ctx context.Context) string {
	return "widget name"
}

// MarkdownDescription is a framework validator method.
func (v widgetNameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString is a framework validator method.
func (v widgetNameValidator) ValidateString(ctx context.Context) {}
