package errors

// standardErrorDescriptions provides consistent error messaging across the provider
var standardErrorDescriptions = map[int]ErrorDescription{
	400: {
		Summary: "Bad Request - 400",
		Detail:  "The request was invalid or malformed. Please check the request parameters and try again.",
	},
	401: {
		Summary: "Unauthorized - 401",
		Detail:  "Authentication failed. Please check your Entra ID credentials and permissions.",
	},
	403: {
		Summary: "Forbidden - 403",
		Detail:  "Your credentials lack sufficient authorisation to perform this operation. Grant the required Microsoft Graph permissions to your Entra ID authentication method.",
	},
	404: {
		Summary: "Not Found - 404",
		Detail:  "The requested resource was not found.",
	},
	405: {
		Summary: "Method Not Allowed - 405",
		Detail:  "The HTTP method in the request isn't allowed on the resource.",
	},
	406: {
		Summary: "Not Acceptable - 406",
		Detail:  "The Microsoft Graph API doesn't support the format requested in the Accept header.",
	},
	409: {
		Summary: "Conflict - 409",
		Detail:  "The operation failed due to a conflict with the current state of the target resource. This might be due to multiple clients modifying the same resource simultaneously, the requested resource may not be in the state that was expected, or the request itself may create a conflict if it is completed.",
	},
	410: {
		Summary: "Gone - 410",
		Detail:  "The requested resource is no longer available at the Microsoft Graph API server.",
	},
	411: {
		Summary: "Length Required - 411",
		Detail:  "A Content-Length header is required on the request to the Microsoft Graph API.",
	},
	412: {
		Summary: "Precondition Failed - 412",
		Detail:  "A precondition provided in the request (such as an if-match header) doesn't match the resource's current state.",
	},
	413: {
		Summary: "Request Entity Too Large - 413",
		Detail:  "The request size exceeds the maximum limit allowed by the Microsoft Graph API.",
	},
	415: {
		Summary: "Unsupported Media Type - 415",
		Detail:  "The content type of the request is a format that isn't supported by the Microsoft Graph API.",
	},
	416: {
		Summary: "Requested Range Not Satisfiable - 416",
		Detail:  "The specified byte range in the request is invalid or unavailable.",
	},
	422: {
		Summary: "Unprocessable Entity - 422",
		Detail:  "The Microsoft Graph API can't process the request because it is semantically incorrect.",
	},
	423: {
		Summary: "Locked - 423",
		Detail:  "The resource that is being accessed is locked.",
	},
	425: {
		Summary: "Too Early - 425",
		Detail:  "The server is unwilling to risk processing a request that might be replayed.",
	},
	428: {
		Summary: "Precondition Required - 428",
		Detail:  "The origin server requires the request to be conditional.",
	},
	429: {
		Summary: "Too Many Requests - 429",
		Detail:  "Request throttled by Microsoft Graph API rate limits. Please try again later.",
	},
	431: {
		Summary: "Request Header Fields Too Large - 431",
		Detail:  "The server is unwilling to process the request because its header fields are too large.",
	},
	500: {
		Summary: "Internal Server Error - 500",
		Detail:  "Microsoft Graph API encountered an internal error. Please try again later.",
	},
	501: {
		Summary: "Not Implemented - 501",
		Detail:  "The requested feature isn't implemented in the Microsoft Graph API.",
	},
	502: {
		Summary: "Bad Gateway - 502",
		Detail:  "The server received an invalid response from an upstream server while trying to fulfill the request.",
	},
	503: {
		Summary: "Service Unavailable - 503",
		Detail:  "The Microsoft Graph API service is temporarily unavailable or overloaded. This is typically a transient condition that will be automatically resolved after a short time.",
	},
	504: {
		Summary: "Gateway Timeout - 504",
		Detail:  "The server, while acting as a proxy, didn't receive a timely response from the upstream server it needed to access in attempting to complete the request. This is typically due to the complexity of the query or high server load. Try simplifying your query or retry later.",
	},
	507: {
		Summary: "Insufficient Storage - 507",
		Detail:  "The maximum storage quota for the Microsoft Graph API has been reached.",
	},
	509: {
		Summary: "Bandwidth Limit Exceeded - 509",
		Detail:  "Your application has been throttled for exceeding the maximum bandwidth cap on the Microsoft Graph API. Your application can retry the request after more time has elapsed.",
	},
}
