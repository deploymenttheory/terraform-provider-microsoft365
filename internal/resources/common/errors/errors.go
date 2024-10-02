package errors

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// GraphErrorInfo contains extracted information from a Graph API error
type GraphErrorInfo struct {
	StatusCode     int
	ErrorCode      string
	ErrorMessage   string
	IsODataError   bool
	AdditionalData map[string]interface{}
}

// GraphError logs detailed information about errors returned by the Graph SDK
// and returns extracted error information for further processing
func GraphError(ctx context.Context, err error) GraphErrorInfo {
	errorInfo := GraphErrorInfo{
		StatusCode:     0,
		IsODataError:   false,
		AdditionalData: make(map[string]interface{}),
	}

	if err == nil {
		return errorInfo
	}

	tflog.Error(ctx, fmt.Sprintf("Raw error: %v", err))

	var apiError abstractions.ApiErrorable
	var ok bool
	if apiError, ok = err.(abstractions.ApiErrorable); !ok {
		tflog.Error(ctx, "Error does not implement ApiErrorable interface")
		tflog.Error(ctx, fmt.Sprintf("Error type: %T", err))
		tflog.Error(ctx, fmt.Sprintf("Error message: %s", err.Error()))
		errorInfo.ErrorMessage = err.Error()
		return errorInfo
	}

	errorInfo.StatusCode = apiError.GetStatusCode()
	tflog.Info(ctx, fmt.Sprintf("HTTP Status Code: %d", errorInfo.StatusCode))

	responseHeaders := apiError.GetResponseHeaders()
	if responseHeaders != nil {
		tflog.Debug(ctx, "Response Headers:")
		for _, key := range responseHeaders.ListKeys() {
			values := responseHeaders.Get(key)
			for _, value := range values {
				tflog.Debug(ctx, fmt.Sprintf("%s: %s", key, value))
			}
		}
	}

	if odataErr, ok := err.(*odataerrors.ODataError); ok {
		errorInfo.IsODataError = true
		if mainError := odataErr.GetErrorEscaped(); mainError != nil {
			if code := mainError.GetCode(); code != nil {
				errorInfo.ErrorCode = *code
				tflog.Info(ctx, fmt.Sprintf("Error Code: %s", errorInfo.ErrorCode))
			}
			if message := mainError.GetMessage(); message != nil {
				errorInfo.ErrorMessage = *message
				tflog.Info(ctx, fmt.Sprintf("Error Message: %s", errorInfo.ErrorMessage))
			}
		}

		additionalData := odataErr.GetAdditionalData()
		if len(additionalData) > 0 {
			tflog.Debug(ctx, "Additional Data:")
			for key, value := range additionalData {
				errorInfo.AdditionalData[key] = value
				tflog.Debug(ctx, fmt.Sprintf("%s: %v", key, value))
			}
		}
	}

	return errorInfo
}
