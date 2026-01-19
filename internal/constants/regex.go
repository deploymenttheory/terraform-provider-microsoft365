package constants

import (
	"time"
)

const (
	// GuidRegex matches a standard GUID/UUID.
	// Example: "123e4567-e89b-12d3-a456-426614174000"
	GuidRegex = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"

	// EntraIdSidRegex matches a Microsoft Entra ID (Azure AD) Security Identifier (SID).
	// Example: "S-1-12-1-1943430372-1249052806-2496021943-3034400218"
	EntraIdSidRegex = `^S-1-12-1-\d+-\d+-\d+-\d+$`

	// GuidOrEmptyValueRegex matches a standard GUID/UUID or an empty string.
	// Example: "123e4567-e89b-12d3-a456-426614174000" or ""
	GuidOrEmptyValueRegex = "^(?:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})?$"

	// PrefixedGuidRegex matches a GUID with a single character prefix followed by underscore.
	// Example: "A_123e4567-e89b-12d3-a456-426614174000" or "T_00000000-0000-0000-0000-000000000000"
	PrefixedGuidRegex = "^[0-9a-zA-Z]_[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"

	// GuidOrPrefixedGuidRegex matches either a standard GUID or a GUID with a single character prefix.
	// Example: "123e4567-e89b-12d3-a456-426614174000" or "A_123e4567-e89b-12d3-a456-426614174000"
	GuidOrPrefixedGuidRegex = "^(?:[0-9a-zA-Z]_)?[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"

	// UrlValidStringRegex matches a valid URL string (letters, numbers, and URL-safe characters).
	// Example: "https://example.com/path?query=1"
	UrlValidStringRegex = "(?i)^[A-Za-z0-9-._~%/:/?=]+$"

	// HttpOrHttpsUrlRegex matches a URL that starts with either http:// or https://
	// Example: "https://example.com" or "http://example.org"
	HttpOrHttpsUrlRegex = "^https?://.*$"

	// EmailRegex matches a valid email address format
	// Example: "user@example.com"
	EmailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// UserPrincipalNameRegex matches a valid Microsoft 365 User Principal Name (UPN).
	// The UPN format follows RFC 822 and only allows specific characters in the alias portion.
	// Allowed characters: A-Z, a-z, 0-9, ' . - _ ! # ^ ~
	// Example: "user.name@contoso.com" or "first-last@contoso.onmicrosoft.com"
	UserPrincipalNameRegex = `^[A-Za-z0-9'.\-_!#^~]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`

	// ApiIdRegex matches API IDs consisting of alphanumeric characters, slashes, dots, or underscores.
	// Example: "api/v1/resource_1"
	ApiIdRegex = "^[0-9a-zA-Z/._]*$"

	// StringRegex matches any string (including empty).
	// Example: "any string here"
	StringRegex = "^.*$"

	// VersionRegex matches a version string in the format "X.Y.Z.W".
	// Example: "1.0.0.0"
	VersionRegex = "^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$"

	// OSVersionRegex matches an operating system version string in the format "X.Y.Z.W" with any number of digits.
	// Example: "10.0.22631.9999" or "1.1.1.1"
	OSVersionRegex = `^\d+\.\d+\.\d+\.\d+$`

	// SemVerRegex matches a Semantic Versioning string in the format "X.Y.Z" (Major.Minor.Patch).
	// Examples: "1.0.0", "2.1.3", "10.20.30"
	SemVerRegex = `^[0-9]+\.[0-9]+\.[0-9]+$`

	// MajorMinorVersionRegex matches a version string in the format "X.Y" (Major.Minor).
	// Examples: "1.0", "2.1", "10.20"
	MajorMinorVersionRegex = `^[0-9]+\.[0-9]+$`

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

	// IPv4CIDRRegex matches a valid IPv4 CIDR range.
	// Example: "192.168.1.0/24"
	IPv4CIDRRegex = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)/(3[0-2]|[12]?[0-9])$`

	// IPv6CIDRRegex matches a valid IPv6 CIDR range.
	// Example: "2001:db8::/32"
	IPv6CIDRRegex = `^([0-9a-fA-F]{0,4}:){1,7}[0-9a-fA-F]{0,4}(\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?$`

	// PortRangeRegex matches a valid port range 0-65535, hyphen separated
	// Example: "80-80", "443-443", "8080-8080", "8443-8443"
	PortRangeRegex = `^(?:0|[1-9]\d{0,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])-(?:0|[1-9]\d{0,3}|[1-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])$`

	// SubjectKeyIdentifierRegex matches a Subject Key Identifier (SKI) in hexadecimal format (40 hex characters).
	// Example: "1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0B"
	SubjectKeyIdentifierRegex = `^[0-9A-Fa-f]{40}$`

	// X509CertificateIssuerRegex matches a custom certificate issuer identifier format (CUSTOMIDENTIFIER: followed by 40-character hex SKI).
	// Example: "CUSTOMIDENTIFIER:1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A0B"
	X509CertificateIssuerRegex = `^CUSTOMIDENTIFIER:[0-9A-Fa-f]{40}$`

	// OIDRegex matches a valid Object Identifier in dotted decimal notation.
	// Example: "1.3.6.1.4.1.311.21.8.1.1"
	OIDRegex = `^[0-9]+(\.[0-9]+)+$`

	// ActiveDirectoryDNRegex matches a valid Active Directory distinguished name (DN).
	// Example: "OU=Computers,DC=contoso,DC=com" or "CN=Server1,OU=Servers,DC=contoso,DC=com"
	ActiveDirectoryDNRegex = `^(OU=|CN=)[^,]+(,(OU=|CN=|DC=)[^,]+)*$`

	// ISO8601DateTimeRegex matches an ISO 8601 datetime format with optional milliseconds and timezone.
	// Examples: "2023-05-01T13:45:30Z", "2023-05-01T13:45:30.123Z", "2023-05-01T13:45:30+00:00", "2023-05-01T13:45:30.123456+05:30"
	ISO8601DateTimeRegex = `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?(Z|[+-]\d{2}:\d{2})$`

	// LocaleRegex matches a locale string in the format "xx-YY" where xx is a 2-letter ISO 639-1 language code (lowercase)
	// and YY is a 2-letter ISO 3166-1 alpha-2 country code (uppercase).
	// Examples: "en-US", "fr-FR", "de-DE", "ja-JP", "es-ES"
	LocaleRegex = `^[a-z]{2}-[A-Z]{2}$`

	// IdentifierUriRegex matches a valid Microsoft Entra ID application identifier URI.
	// Common formats: "api://<guid>", "api://<domain>/<path>", "https://<domain>/<path>", "urn:<namespace>:<identifier>"
	// Examples: "api://123e4567-e89b-12d3-a456-426614174000", "api://contoso.com/myapp", "https://contoso.com/api"
	IdentifierUriRegex = `^(api://|https://|urn:)[^\s]+$`

	// DayOfWeekRegex matches lowercase day of week names.
	// Valid values: "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"
	DayOfWeekRegex = `^(monday|tuesday|wednesday|thursday|friday|saturday|sunday)$`
)
