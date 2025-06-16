package errors

// Enhanced comprehensive OData error codes
var comprehensiveODataErrorCodes = map[string]string{
	// Query operation errors
	"ExpandNotSupported":  "The $expand operation is not supported for the specified property or relationship.",
	"FilterNotSupported":  "The $filter operation is not supported for the specified property.",
	"SelectNotSupported":  "The $select operation is not supported for the specified property.",
	"OrderByNotSupported": "The $orderby operation is not supported for the specified property.",
	"SearchNotSupported":  "The $search operation is not supported for this endpoint.",
	"InvalidFilter":       "The specified filter syntax or expression is invalid.",
	"InvalidQuery":        "The OData query contains invalid syntax or unsupported operations.",
	"QueryTooComplex":     "The OData query is too complex for the service to process.",
	"NotImplemented":      "The requested OData feature is not implemented for this endpoint.",
	"SyntaxError":         "There is a syntax error in the OData query string.",

	// Authentication and authorization errors
	"InvalidAuthenticationToken": "The access token is invalid or expired.",
	"Forbidden":                  "The caller does not have permission to perform the operation.",
	"InvalidRequest":             "The request is invalid.",
	"ResourceNotFound":           "The specified resource does not exist.",
	"Unauthorized":               "Authentication is required to access this resource.",
	"CompactTokenParsingFailure": "Failed to parse the compact token.",
	"TokenExpired":               "The access token has expired.",
	"InvalidToken":               "The provided token is invalid.",

	// Request processing errors
	"BadRequest":            "The request could not be understood by the server due to malformed syntax.",
	"Conflict":              "The request could not be completed due to a conflict with the current state of the resource.",
	"PreconditionFailed":    "The precondition given in one or more of the request header fields evaluated to false.",
	"RequestEntityTooLarge": "The request entity is larger than the server is willing or able to process.",
	"InvalidRange":          "The specified range is invalid.",
	"UnsupportedMediaType":  "The media type is not supported.",
	"MethodNotAllowed":      "The method specified in the request is not allowed for the resource.",

	// Rate limiting and throttling
	"TooManyRequests":    "The user has sent too many requests in a given amount of time.",
	"ServiceUnavailable": "The server is currently unable to handle the request due to temporary overloading or maintenance.",
	"RequestThrottled":   "The request has been throttled.",
	"QuotaLimitExceeded": "The quota limit has been exceeded.",

	// Specific Microsoft Graph errors
	"ItemNotFound":                    "The requested item was not found.",
	"NameAlreadyExists":               "The specified name already exists.",
	"LockNotFoundOrAlreadyExpired":    "The lock was not found or has already expired.",
	"UnknownError":                    "An unknown error has occurred.",
	"ActivityLimitReached":            "The activity limit has been reached.",
	"GeneralException":                "A general exception occurred.",
	"NotAllowed":                      "The operation is not allowed.",
	"ResourceModified":                "The resource has been modified.",
	"ResyncRequired":                  "Resynchronization is required.",
	"ErrorInvalidIdMalformed":         "The provided ID is malformed.",
	"ErrorInvalidUser":                "The specified user is invalid.",
	"ErrorItemNotFound":               "The specified item was not found.",
	"ErrorInvalidOperation":           "The operation is invalid.",
	"ErrorAccessDenied":               "Access is denied.",
	"ErrorUnsupportedOperation":       "The operation is not supported.",
	"ErrorInternalServerUnknownError": "An unknown internal server error occurred.",

	// Tenant and licensing errors
	"TenantNotFound":       "The specified tenant was not found.",
	"FeatureNotEnabled":    "The required feature is not enabled for this tenant.",
	"LicenseRequired":      "A license is required to access this feature.",
	"SubscriptionRequired": "A subscription is required to access this resource.",

	// Validation errors
	"InvalidValue":            "The provided value is invalid.",
	"MissingRequiredProperty": "A required property is missing.",
	"PropertyNotUpdatable":    "The specified property cannot be updated.",
	"InvalidPropertyValue":    "The property value is invalid.",
	"DuplicateValue":          "The provided value already exists.",
	"ValueOutOfRange":         "The provided value is out of the acceptable range.",
}
