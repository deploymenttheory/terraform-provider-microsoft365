package helpers

// StringPtrToString converts a string pointer to a string.
func StringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
