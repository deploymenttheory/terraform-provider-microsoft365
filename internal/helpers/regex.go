package helpers

import "time"

const (
	GuidRegex             = "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
	GuidOrEmptyValueRegex = "^(?:[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})?$"
	UrlValidStringRegex   = "(?i)^[A-Za-z0-9-._~%/:/?=]+$"
	ApiIdRegex            = "^[0-9a-zA-Z/._]*$"
	StringRegex           = "^.*$"
	VersionRegex          = "^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$"
	TimeRegex             = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z)$"
	BooleanRegex          = "^(true|false)$"
	TimeFormatRFC3339     = time.RFC3339
)
