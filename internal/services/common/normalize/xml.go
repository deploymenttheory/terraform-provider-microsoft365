package normalize

import (
	"html"
	"strings"
)

// NormalizeXML normalizes XML content by removing BOM, normalizing whitespace, and unescaping HTML entities
// for create and update constructor operations
func NormalizeXML(xmlContent string) string {
	// Remove BOM if present
	xmlContent = strings.TrimPrefix(xmlContent, "\ufeff")

	// Unescape HTML entities like &lt; &gt;
	xmlContent = html.UnescapeString(xmlContent)

	return xmlContent
}

// reverseNormalizeXMLContent reverses normalization to match original Terraform format
// for read stating operations
func ReverseNormalizeXML(xmlContent string) string {
	// First normalize the content to ensure consistent format
	xmlContent = NormalizeXML(xmlContent)

	// Only add BOM if the xml content is likely came to have come from a file
	// This prevents adding BOM to .tf files with inline XML which causes Terraform diff issues
	if LikelyFromFile(xmlContent) {
		if !strings.HasPrefix(xmlContent, "\ufeff") {
			xmlContent = "\ufeff" + xmlContent
		}
	}

	// Note: We don't re-escape HTML entities as they should stay as proper XML
	return xmlContent
}

// likelyFromFile determines if the XML content likely came from a file
// This is a heuristic approach that checks for common file as source indicators
func LikelyFromFile(xmlContent string) bool {
	// Check for Windows line endings (CRLF) which are common in files
	if strings.Contains(xmlContent, "\r\n") {
		return true
	}

	// Check if content already has a BOM, which typically happens with files
	if strings.HasPrefix(xmlContent, "\ufeff") {
		return true
	}

	// Additional heuristics could be added here
	// For example, checking for consistent indentation patterns common in files

	return false
}
