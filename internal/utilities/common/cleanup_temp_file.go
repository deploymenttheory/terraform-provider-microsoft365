package common

// cleanupTempFile removes a temporary file if it exists and should be cleaned up
// func cleanupTempFile(ctx context.Context, fileInfo TempFileInfo) {
// 	if !fileInfo.ShouldCleanup || fileInfo.FilePath == "" {
// 		return
// 	}

// 	if err := os.Remove(fileInfo.FilePath); err != nil {
// 		tflog.Warn(ctx, fmt.Sprintf("Failed to remove temporary file %s: %v", fileInfo.FilePath, err))
// 	} else {
// 		tflog.Debug(ctx, fmt.Sprintf("Successfully removed temporary file: %s", fileInfo.FilePath))
// 	}
// }
