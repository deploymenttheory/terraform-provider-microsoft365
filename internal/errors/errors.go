package errors

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

// HandleGraphError handles errors returned by the Graph SDK
func HandleGraphError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	var apiError abstractions.ApiErrorable
	var ok bool

	tflog.Error(ctx, fmt.Sprintf("Raw error: %v", err))

	if apiError, ok = err.(abstractions.ApiErrorable); !ok {
		tflog.Error(ctx, "Error does not implement ApiErrorable interface")
		tflog.Error(ctx, fmt.Sprintf("Error type: %T", err))
		tflog.Error(ctx, fmt.Sprintf("Error message: %s", err.Error()))
		return
	}

	// Return HTTP status code and resp headers
	statusCode := apiError.GetStatusCode()
	tflog.Info(ctx, fmt.Sprintf("HTTP Status Code: %d", statusCode))

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

	// If it's an ODataError, we can get more detailed information
	if odataErr, ok := err.(*odataerrors.ODataError); ok {
		if mainError := odataErr.GetErrorEscaped(); mainError != nil {
			if code := mainError.GetCode(); code != nil {
				tflog.Info(ctx, fmt.Sprintf("Error Code: %s", *code))
			}
			if message := mainError.GetMessage(); message != nil {
				tflog.Info(ctx, fmt.Sprintf("Error Message: %s", *message))
			}
		}

		additionalData := odataErr.GetAdditionalData()
		if len(additionalData) > 0 {
			tflog.Debug(ctx, "Additional Data:")
			for key, value := range additionalData {
				tflog.Debug(ctx, fmt.Sprintf("%s: %v", key, value))
			}
		}
	}
}

// HandleGraphErrorStructured handles errors returned by the Graph SDK and logs them using structured logging
func HandleGraphErrorStructured(ctx context.Context, err error) {
	if err == nil {
		return
	}
	var apiError abstractions.ApiErrorable
	var ok bool

	tflog.Error(ctx, "Raw error", map[string]interface{}{"error": err})

	if apiError, ok = err.(abstractions.ApiErrorable); !ok {
		tflog.Error(ctx, "Error does not implement ApiErrorable interface")
		tflog.Error(ctx, "Error details", map[string]interface{}{
			"type":    fmt.Sprintf("%T", err),
			"message": err.Error(),
		})
		return
	}

	statusCode := apiError.GetStatusCode()
	tflog.Info(ctx, "HTTP Status Code", map[string]interface{}{"statusCode": statusCode})

	responseHeaders := apiError.GetResponseHeaders()
	if responseHeaders != nil {
		headerMap := make(map[string]interface{})
		for _, key := range responseHeaders.ListKeys() {
			values := responseHeaders.Get(key)
			if len(values) == 1 {
				headerMap[key] = values[0]
			} else if len(values) > 1 {
				headerMap[key] = values
			}
		}
		tflog.Debug(ctx, "Response Headers", map[string]interface{}{"headers": headerMap})
	}

	// If it's an ODataError, we can get more detailed information
	if odataErr, ok := err.(*odataerrors.ODataError); ok {
		if mainError := odataErr.GetErrorEscaped(); mainError != nil {
			if code := mainError.GetCode(); code != nil {
				tflog.Info(ctx, "Error Code", map[string]interface{}{"code": *code})
			}
			if message := mainError.GetMessage(); message != nil {
				tflog.Info(ctx, "Error Message", map[string]interface{}{"message": *message})
			}
		}
	}
}

// HandleGraphErrorUnstructured handles errors returned by the Graph SDK and logs them using unstructured logging
func HandleGraphErrorUnstructured(ctx context.Context, err error) {
	if err == nil {
		return
	}
	var apiError abstractions.ApiErrorable
	var ok bool

	tflog.Error(ctx, fmt.Sprintf("Raw error: %v", err))

	if apiError, ok = err.(abstractions.ApiErrorable); !ok {
		tflog.Error(ctx, "Error does not implement ApiErrorable interface")
		tflog.Error(ctx, fmt.Sprintf("Error type: %T", err))
		tflog.Error(ctx, fmt.Sprintf("Error message: %s", err.Error()))
		return
	}

	statusCode := apiError.GetStatusCode()
	tflog.Info(ctx, fmt.Sprintf("HTTP Status Code: %d", statusCode))

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
		if mainError := odataErr.GetErrorEscaped(); mainError != nil {
			if code := mainError.GetCode(); code != nil {
				tflog.Info(ctx, fmt.Sprintf("Error Code: %s", *code))
			}
			if message := mainError.GetMessage(); message != nil {
				tflog.Info(ctx, fmt.Sprintf("Error Message: %s", *message))
			}
		}
	}
}
