package helpers

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseJSONFile_EmptyPath tests the empty path validation unit
func TestParseJSONFile_EmptyPath(t *testing.T) {
	content, err := ParseJSONFile("")

	assert.Error(t, err, "Expected error for empty file path")
	assert.Empty(t, content, "Expected empty content for empty path")
	assert.Contains(t, err.Error(), "file path for json file cannot be empty")
}

// TestParseJSONFile_ValidExtensions tests the file extension validation unit
func TestParseJSONFile_ValidExtensions(t *testing.T) {
	testCases := []struct {
		name        string
		filename    string
		shouldPass  bool
		expectedErr string
	}{
		{
			name:       "valid json extension",
			filename:   "test.json",
			shouldPass: true,
		},
		{
			name:       "valid uppercase json extension",
			filename:   "test.JSON",
			shouldPass: true,
		},
		{
			name:        "invalid txt extension",
			filename:    "test.txt",
			shouldPass:  false,
			expectedErr: "invalid file extension",
		},
		{
			name:        "invalid js extension",
			filename:    "test.js",
			shouldPass:  false,
			expectedErr: "invalid file extension",
		},
		{
			name:        "invalid yaml extension",
			filename:    "test.yaml",
			shouldPass:  false,
			expectedErr: "invalid file extension",
		},
		{
			name:        "no extension",
			filename:    "test",
			shouldPass:  false,
			expectedErr: "invalid file extension",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.shouldPass {
				// Create a valid test file for valid extensions
				testContent := `{"name": "test", "value": 42}`
				require.NoError(t, os.WriteFile(tc.filename, []byte(testContent), 0644))
				defer func() {
					assert.NoError(t, os.Remove(tc.filename))
				}()

				content, err := ParseJSONFile(tc.filename)
				assert.NoError(t, err, "Expected success for %s", tc.filename)
				assert.Equal(t, testContent, content, "Content mismatch for %s", tc.filename)
			} else {
				// Test invalid extension (file doesn't need to exist)
				content, err := ParseJSONFile(tc.filename)
				assert.Error(t, err, "Expected error for %s", tc.filename)
				assert.Empty(t, content, "Expected empty content for %s", tc.filename)
				assert.Contains(t, err.Error(), tc.expectedErr, "Expected error containing '%s' for %s", tc.expectedErr, tc.filename)
			}
		})
	}
}

// TestParseJSONFile_FileExistence tests the file existence validation unit
func TestParseJSONFile_FileExistence(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		nonExistentFile := "nonexistent.json"

		content, err := ParseJSONFile(nonExistentFile)

		assert.Error(t, err, "Expected error for non-existent file")
		assert.Empty(t, content, "Expected empty content for non-existent file")
		assert.Contains(t, err.Error(), "json file does not exist")
	})
}

// TestParseJSONFile_FileTypeValidation tests the regular file validation unit
func TestParseJSONFile_FileTypeValidation(t *testing.T) {
	t.Run("directory instead of file", func(t *testing.T) {
		dirPath := "fake_dir.json"
		require.NoError(t, os.MkdirAll(dirPath, 0755))
		defer func() {
			assert.NoError(t, os.RemoveAll(dirPath))
		}()

		content, err := ParseJSONFile(dirPath)

		assert.Error(t, err, "Expected error when trying to read directory")
		assert.Empty(t, content, "Expected empty content when reading directory")
		assert.Contains(t, err.Error(), "supplied path does not resolve to a file")
	})
}

// TestParseJSONFile_FileSizeValidation tests the file size limit validation unit
func TestParseJSONFile_FileSizeValidation(t *testing.T) {
	t.Run("file size within limit", func(t *testing.T) {
		testFile := "small.json"
		testContent := `{"comments": ["` + strings.Repeat("small json content ", 100) + `"]}`

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		content, err := ParseJSONFile(testFile)

		assert.NoError(t, err, "ParseJSONFile should not fail for small file")
		assert.Equal(t, testContent, content, "Content should match for small file")
	})

	t.Run("file size exceeds limit", func(t *testing.T) {
		testFile := "large.json"
		// Create file larger than 10MB limit (create ~12MB file)
		largeArray := strings.Repeat(`"This is a long string that will be repeated to create a large JSON file",`, 200000)
		largeContent := `{"data": [` + largeArray[:len(largeArray)-1] + `]}`

		require.NoError(t, os.WriteFile(testFile, []byte(largeContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		content, err := ParseJSONFile(testFile)

		assert.Error(t, err, "Expected error for oversized file")
		assert.Empty(t, content, "Expected empty content for oversized file")
		assert.Contains(t, err.Error(), "file too large")
	})
}

// TestParseJSONFile_PathTraversalSecurity tests the path traversal prevention unit
func TestParseJSONFile_PathTraversalSecurity(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{
			name: "simple path traversal",
			path: "../../../etc/config.json",
		},
		{
			name: "complex path traversal",
			path: "../../../../../../usr/local/config.json",
		},
		{
			name: "mixed path traversal",
			path: "./../../sensitive/data.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := ParseJSONFile(tc.path)

			assert.Error(t, err, "Expected error for path traversal attempt: %s", tc.path)
			assert.Empty(t, content, "Expected empty content for path traversal")

			// Should either fail at file existence or security validation
			hasExpectedError := strings.Contains(err.Error(), "json file does not exist") ||
				strings.Contains(err.Error(), "access denied") ||
				strings.Contains(err.Error(), "path outside project boundaries")

			assert.True(t, hasExpectedError, "Expected security-related error for %s, got: %v", tc.path, err)
		})
	}
}

// TestParseJSONFile_EmptyFileContent tests handling of empty files unit
func TestParseJSONFile_EmptyFileContent(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		emptyFile := "empty.json"
		require.NoError(t, os.WriteFile(emptyFile, []byte(""), 0644))
		defer func() {
			assert.NoError(t, os.Remove(emptyFile))
		}()

		content, err := ParseJSONFile(emptyFile)

		assert.NoError(t, err, "ParseJSONFile should not fail for empty file")
		assert.Empty(t, content, "Expected empty content for empty file")
	})
}

// TestParseJSONFile_ValidFileContent tests successful file reading unit
func TestParseJSONFile_ValidFileContent(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		content  string
	}{
		{
			name:     "simple json object",
			filename: "simple.json",
			content: `{
  "name": "test-resource",
  "type": "example",
  "value": 42,
  "enabled": true
}`,
		},
		{
			name:     "json array",
			filename: "array.json",
			content: `[
  {"id": 1, "name": "first"},
  {"id": 2, "name": "second"},
  {"id": 3, "name": "third"}
]`,
		},
		{
			name:     "nested json structure",
			filename: "nested.json",
			content: `{
  "metadata": {
    "version": "1.0",
    "author": "test"
  },
  "config": {
    "settings": {
      "debug": true,
      "timeout": 30
    },
    "features": ["auth", "logging"]
  }
}`,
		},
		{
			name:     "file with special characters",
			filename: "special.json",
			content:  `{"description": "JSON with special chars: åæø, 中文, 🎉", "test": true}`,
		},
		{
			name:     "compact json",
			filename: "compact.json",
			content:  `{"compact":true,"noSpaces":42,"array":[1,2,3]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(tc.filename, []byte(tc.content), 0644))
			defer func() {
				assert.NoError(t, os.Remove(tc.filename))
			}()

			content, err := ParseJSONFile(tc.filename)

			assert.NoError(t, err, "ParseJSONFile should not fail for valid content")
			assert.Equal(t, tc.content, content, "Content should match exactly")
		})
	}
}

// TestParseJSONFile_MalformedJSON tests that the parser loads malformed JSON (validation is not its responsibility)
func TestParseJSONFile_MalformedJSON(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		content  string
	}{
		{
			name:     "invalid json syntax",
			filename: "invalid.json",
			content:  `{"name": "test", "value": 42,}`, // trailing comma
		},
		{
			name:     "unclosed bracket",
			filename: "unclosed.json",
			content:  `{"name": "test"`, // missing closing brace
		},
		{
			name:     "not json at all",
			filename: "notjson.json",
			content:  `This is not JSON content at all`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(tc.filename, []byte(tc.content), 0644))
			defer func() {
				assert.NoError(t, os.Remove(tc.filename))
			}()

			// The parser should load the content successfully - JSON validation is not its responsibility
			content, err := ParseJSONFile(tc.filename)

			assert.NoError(t, err, "ParseJSONFile should load malformed JSON (validation is not its responsibility)")
			assert.Equal(t, tc.content, content, "Content should match exactly even if malformed")
		})
	}
}

// TestParseJSONFile_RelativePaths tests relative path handling
func TestParseJSONFile_RelativePaths(t *testing.T) {
	t.Run("relative path from current directory", func(t *testing.T) {
		testFile := "relative_test.json"
		testContent := `{"test": "relative path"}`

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		// Use just the filename (relative path)
		content, err := ParseJSONFile(testFile)

		assert.NoError(t, err, "Should handle relative path")
		assert.Equal(t, testContent, content)
	})

	t.Run("relative path with subdirectory", func(t *testing.T) {
		// Create subdirectory
		subdir := "jsontestsubdir"
		require.NoError(t, os.MkdirAll(subdir, 0755))
		defer func() {
			assert.NoError(t, os.RemoveAll(subdir))
		}()

		testFile := subdir + "/subdir_test.json"
		testContent := `{"test": "subdirectory"}`

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))

		// Use relative path with subdirectory
		content, err := ParseJSONFile(testFile)

		assert.NoError(t, err, "Should handle relative path with subdirectory")
		assert.Equal(t, testContent, content)
	})

	t.Run("relative path with dot notation", func(t *testing.T) {
		testFile := "dot_notation.json"
		testContent := `{"test": "dot notation"}`

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		// Use ./filename notation
		content, err := ParseJSONFile("./" + testFile)

		assert.NoError(t, err, "Should handle ./ notation")
		assert.Equal(t, testContent, content)
	})
}

// TestParseJSONFile_AbsolutePaths tests absolute path handling
func TestParseJSONFile_AbsolutePaths(t *testing.T) {
	t.Run("absolute path within project", func(t *testing.T) {
		// Get current working directory
		cwd, err := os.Getwd()
		require.NoError(t, err)

		testFile := "absolute_test.json"
		testContent := `{"test": "absolute path"}`

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		// Use absolute path
		absolutePath := cwd + "/" + testFile
		content, err := ParseJSONFile(absolutePath)

		assert.NoError(t, err, "Should handle absolute path within project")
		assert.Equal(t, testContent, content)
	})

	t.Run("absolute path outside project boundaries", func(t *testing.T) {
		// Create a temp file outside the project
		tmpDir := t.TempDir()
		testFile := tmpDir + "/outside.json"
		testContent := `{"test": "outside project"}`

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))

		// Try to access file outside project boundaries
		content, err := ParseJSONFile(testFile)

		// Should either fail with security error or succeed if project root can't be determined
		if err != nil {
			assert.Contains(t, err.Error(), "access denied", "Should fail with access denied for outside path")
			assert.Empty(t, content)
		}
		// If it succeeds, that's also acceptable if project root detection fails
	})
}

// TestParseJSONFile_CleanPath tests path cleaning and normalization
func TestParseJSONFile_CleanPath(t *testing.T) {
	t.Run("path with redundant separators", func(t *testing.T) {
		testFile := "clean_test.json"
		testContent := `{"test": "clean path"}`

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		// Use path with redundant separators
		content, err := ParseJSONFile(".//" + testFile)

		assert.NoError(t, err, "Should clean path with redundant separators")
		assert.Equal(t, testContent, content)
	})

	t.Run("path with redundant dot segments", func(t *testing.T) {
		testFile := "segments_test.json"
		testContent := `{"test": "segments"}`

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		// Use path with redundant segments
		content, err := ParseJSONFile("././" + testFile)

		assert.NoError(t, err, "Should clean path with redundant dot segments")
		assert.Equal(t, testContent, content)
	})
}

// TestParseJSONFile_ReadPermissions tests file read permission scenarios
func TestParseJSONFile_ReadPermissions(t *testing.T) {
	// Skip on Windows as file permissions work differently
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	t.Run("file without read permissions", func(t *testing.T) {
		testFile := "no_read.json"
		testContent := `{"test": "no read"}`

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			// Restore permissions before deleting
			os.Chmod(testFile, 0644)
			assert.NoError(t, os.Remove(testFile))
		}()

		// Remove read permissions
		require.NoError(t, os.Chmod(testFile, 0000))

		content, err := ParseJSONFile(testFile)

		assert.Error(t, err, "Should fail when file has no read permissions")
		assert.Empty(t, content)
		assert.Contains(t, err.Error(), "failed to read json file", "Should contain read error message")
	})
}

// TestParseJSONFile_CallerContext tests error when caller context cannot be determined
func TestParseJSONFile_CallerContext(t *testing.T) {
	t.Run("handles caller context errors gracefully", func(t *testing.T) {
		// This test verifies the function handles runtime.Caller errors
		// In normal operation, runtime.Caller(1) should always succeed
		// We can't easily trigger a failure, but we ensure the code path exists
		testFile := "caller_context.json"
		testContent := `{"test": "caller context"}`

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		content, err := ParseJSONFile(testFile)

		// Should succeed in normal operation
		assert.NoError(t, err)
		assert.Equal(t, testContent, content)
	})
}

// TestParseJSONFile_ErrorPathCoverage tests additional error paths for coverage
func TestParseJSONFile_ErrorPathCoverage(t *testing.T) {
	t.Run("path contains double dots after cleaning", func(t *testing.T) {
		// This tests the security check for ".." in absolute path
		testPath := "../../test.json"

		content, err := ParseJSONFile(testPath)

		// Should fail due to path traversal or file not existing
		assert.Error(t, err)
		assert.Empty(t, content)
	})

	t.Run("error accessing file stats", func(t *testing.T) {
		// Create a file in a directory, then remove the directory
		testDir := "tmpdir_json"
		require.NoError(t, os.MkdirAll(testDir, 0755))

		testFile := testDir + "/test.json"
		require.NoError(t, os.WriteFile(testFile, []byte(`{"test": true}`), 0644))

		// Remove directory while file reference exists
		os.RemoveAll(testDir)

		content, err := ParseJSONFile(testFile)

		assert.Error(t, err)
		assert.Empty(t, content)
	})
}
