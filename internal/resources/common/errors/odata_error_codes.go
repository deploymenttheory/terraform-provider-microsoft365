package errors

var commonODataErrorCodes = map[string]string{
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
}
