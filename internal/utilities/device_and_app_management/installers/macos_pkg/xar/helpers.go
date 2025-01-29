package xar

// Helper function to count compressed files
func countCompressedFiles(payloads map[uint64]PayloadMetadata) int {
	count := 0
	for _, p := range payloads {
		if p.IsCompressed {
			count++
		}
	}
	return count
}
