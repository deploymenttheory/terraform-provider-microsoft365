package constants

import "time"

const (
	// GuidRegex matches a standard GUID/UUID.
	// Example: "123e4567-e89b-12d3-a456-426614174000"
	GuidRegex = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"

	// GuidOrEmptyValueRegex matches a standard GUID/UUID or an empty string.
	// Example: "123e4567-e89b-12d3-a456-426614174000" or ""
	GuidOrEmptyValueRegex = "^(?:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})?$"

	// UrlValidStringRegex matches a valid URL string (letters, numbers, and URL-safe characters).
	// Example: "https://example.com/path?query=1"
	UrlValidStringRegex = "(?i)^[A-Za-z0-9-._~%/:/?=]+$"

	// HttpOrHttpsUrlRegex matches a URL that starts with either http:// or https://
	// Example: "https://example.com" or "http://example.org"
	HttpOrHttpsUrlRegex = "^https?://.*$"

	// ApiIdRegex matches API IDs consisting of alphanumeric characters, slashes, dots, or underscores.
	// Example: "api/v1/resource_1"
	ApiIdRegex = "^[0-9a-zA-Z/._]*$"

	// StringRegex matches any string (including empty).
	// Example: "any string here"
	StringRegex = "^.*$"

	// VersionRegex matches a version string in the format "X.Y.Z.W".
	// Example: "1.0.0.0"
	VersionRegex = "^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$"

	// TimeFormatUTCTimeStampRegex matches a UTC timestamp in the format "YYYY-MM-DDTHH:MM:SSZ".
	// Example: "2023-05-01T13:45:30Z"
	TimeFormatUTCTimeStampRegex = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z)$"

	// BooleanRegex matches the string "true" or "false".
	// Example: "true"
	BooleanRegex = "^(true|false)$"

	// TimeFormatRFC3339Regex is the time format for RFC3339.
	// Example: "2023-05-01T13:45:30Z"
	TimeFormatRFC3339Regex = time.RFC3339

	// ISO8601DurationRegex matches an ISO 8601 duration format.
	// Examples: "P1D" (1 day), "PT1H" (1 hour), "P1W" (1 week), "P1Y2M3DT4H5M6S" (1 year, 2 months, 3 days, 4 hours, 5 minutes, 6 seconds)
	ISO8601DurationRegex = `^P(?:\d+Y)?(?:\d+M)?(?:\d+W)?(?:\d+D)?(?:T(?:\d+H)?(?:\d+M)?(?:\d+S)?)?$`

	// AzureImageResourceIDRegex matches a valid Azure image resource ID for a custom image.
	// Example: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.Compute/images/myimage"
	AzureImageResourceIDRegex = `^/subscriptions/[^/]+/resourceGroups/[^/]+/providers/Microsoft\.Compute/images/[^/]+$`

	// DateFormatYYYYMMDDRegex matches a date string in the format "YYYY-MM-DD".
	// Example: "2023-12-31"
	DateFormatYYYYMMDDRegex = "^\\d{4}-\\d{2}-\\d{2}$"

	// TimeFormatHHMMSSRegex matches a time string in the format "HH:MM:SS" (24-hour clock).
	// Example: "23:59:59"
	TimeFormatHHMMSSRegex = "^([01]\\d|2[0-3]):([0-5]\\d):([0-5]\\d)$"
)
