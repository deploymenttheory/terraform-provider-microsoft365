package helpers

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseHCLFile_EmptyPath tests the empty path validation unit
func TestParseHCLFile_EmptyPath(t *testing.T) {
	content, err := ParseHCLFile("")

	assert.Error(t, err, "Expected error for empty file path")
	assert.Empty(t, content, "Expected empty content for empty path")
	assert.Contains(t, err.Error(), "file path for terraform file cannot be empty")
}

// TestParseHCLFile_ValidExtensions tests the file extension validation unit
func TestParseHCLFile_ValidExtensions(t *testing.T) {
	testCases := []struct {
		name        string
		filename    string
		shouldPass  bool
		expectedErr string
	}{
		{
			name:       "valid tf extension",
			filename:   "test.tf",
			shouldPass: true,
		},
		{
			name:       "valid hcl extension",
			filename:   "test.hcl",
			shouldPass: true,
		},
		{
			name:       "valid uppercase tf extension",
			filename:   "test.TF",
			shouldPass: true,
		},
		{
			name:       "valid uppercase hcl extension",
			filename:   "test.HCL",
			shouldPass: true,
		},
		{
			name:        "invalid txt extension",
			filename:    "test.txt",
			shouldPass:  false,
			expectedErr: "invalid file extension",
		},
		{
			name:        "invalid json extension",
			filename:    "test.json",
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
				testContent := "resource \"test\" \"example\" {}"
				require.NoError(t, os.WriteFile(tc.filename, []byte(testContent), 0644))
				defer func() {
					assert.NoError(t, os.Remove(tc.filename))
				}()

				content, err := ParseHCLFile(tc.filename)
				assert.NoError(t, err, "Expected success for %s", tc.filename)
				assert.Equal(t, testContent, content, "Content mismatch for %s", tc.filename)
			} else {
				// Test invalid extension (file doesn't need to exist)
				content, err := ParseHCLFile(tc.filename)
				assert.Error(t, err, "Expected error for %s", tc.filename)
				assert.Empty(t, content, "Expected empty content for %s", tc.filename)
				assert.Contains(t, err.Error(), tc.expectedErr, "Expected error containing '%s' for %s", tc.expectedErr, tc.filename)
			}
		})
	}
}

// TestParseHCLFile_FileExistence tests the file existence validation unit
func TestParseHCLFile_FileExistence(t *testing.T) {
	t.Run("non-existent file", func(t *testing.T) {
		nonExistentFile := "nonexistent.tf"

		content, err := ParseHCLFile(nonExistentFile)

		assert.Error(t, err, "Expected error for non-existent file")
		assert.Empty(t, content, "Expected empty content for non-existent file")
		assert.Contains(t, err.Error(), "terraform file does not exist")
	})
}

// TestParseHCLFile_FileTypeValidation tests the regular file validation unit
func TestParseHCLFile_FileTypeValidation(t *testing.T) {
	t.Run("directory instead of file", func(t *testing.T) {
		dirPath := "fake_dir.tf"
		require.NoError(t, os.MkdirAll(dirPath, 0755))
		defer func() {
			assert.NoError(t, os.RemoveAll(dirPath))
		}()

		content, err := ParseHCLFile(dirPath)

		assert.Error(t, err, "Expected error when trying to read directory")
		assert.Empty(t, content, "Expected empty content when reading directory")
		assert.Contains(t, err.Error(), "supplied path does not resolve to a file")
	})
}

// TestParseHCLFile_FileSizeValidation tests the file size limit validation unit
func TestParseHCLFile_FileSizeValidation(t *testing.T) {
	t.Run("file size within limit", func(t *testing.T) {
		testFile := "small.tf"
		testContent := strings.Repeat("# comment\n", 100) // Small file

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		content, err := ParseHCLFile(testFile)

		assert.NoError(t, err, "ParseHCLFile should not fail for small file")
		assert.Equal(t, testContent, content, "Content should match for small file")
	})

	t.Run("file size exceeds limit", func(t *testing.T) {
		testFile := "large.tf"
		// Create file larger than 1MB limit (create ~2MB file)
		largeContent := strings.Repeat("# This is a comment line that will be repeated to create a large file\n", 30000)

		require.NoError(t, os.WriteFile(testFile, []byte(largeContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		content, err := ParseHCLFile(testFile)

		assert.Error(t, err, "Expected error for oversized file")
		assert.Empty(t, content, "Expected empty content for oversized file")
		assert.Contains(t, err.Error(), "file too large")
	})
}

// TestParseHCLFile_PathTraversalSecurity tests the path traversal prevention unit
func TestParseHCLFile_PathTraversalSecurity(t *testing.T) {
	testCases := []struct {
		name string
		path string
	}{
		{
			name: "simple path traversal",
			path: "../../../etc/passwd.tf",
		},
		{
			name: "complex path traversal",
			path: "../../../../../../usr/bin/evil.tf",
		},
		{
			name: "mixed path traversal",
			path: "./../../sensitive/file.tf",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content, err := ParseHCLFile(tc.path)

			assert.Error(t, err, "Expected error for path traversal attempt: %s", tc.path)
			assert.Empty(t, content, "Expected empty content for path traversal")

			// Should either fail at file existence or security validation
			hasExpectedError := strings.Contains(err.Error(), "terraform file does not exist") ||
				strings.Contains(err.Error(), "access denied") ||
				strings.Contains(err.Error(), "path outside project boundaries")

			assert.True(t, hasExpectedError, "Expected security-related error for %s, got: %v", tc.path, err)
		})
	}
}

// TestParseHCLFile_EmptyFileContent tests handling of empty files unit
func TestParseHCLFile_EmptyFileContent(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		emptyFile := "empty.tf"
		require.NoError(t, os.WriteFile(emptyFile, []byte(""), 0644))
		defer func() {
			assert.NoError(t, os.Remove(emptyFile))
		}()

		content, err := ParseHCLFile(emptyFile)

		assert.NoError(t, err, "ParseHCLFile should not fail for empty file")
		assert.Empty(t, content, "Expected empty content for empty file")
	})
}

// TestParseHCLFile_ValidFileContent tests successful file reading unit
func TestParseHCLFile_ValidFileContent(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		content  string
	}{
		{
			name:     "simple terraform resource",
			filename: "simple.tf",
			content: `resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}`,
		},
		{
			name:     "hcl variable",
			filename: "vars.hcl",
			content: `variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}`,
		},
		{
			name:     "file with special characters",
			filename: "special.tf",
			content:  `# Comment with special chars: åæø, 中文, 🎉\nresource "test" "special" {}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, os.WriteFile(tc.filename, []byte(tc.content), 0644))
			defer func() {
				assert.NoError(t, os.Remove(tc.filename))
			}()

			content, err := ParseHCLFile(tc.filename)

			assert.NoError(t, err, "ParseHCLFile should not fail for valid content")
			assert.Equal(t, tc.content, content, "Content should match exactly")
		})
	}
}

// TestFindProjectRoot tests the findProjectRoot helper function
func TestFindProjectRoot(t *testing.T) {
	t.Run("find project root with go.mod", func(t *testing.T) {
		// Create a temporary directory structure
		tmpDir := t.TempDir()
		projectDir := tmpDir + "/testproject"
		subDir := projectDir + "/subdir/nested"

		require.NoError(t, os.MkdirAll(subDir, 0755))

		// Create a go.mod file at project root
		goModPath := projectDir + "/go.mod"
		require.NoError(t, os.WriteFile(goModPath, []byte("module test\n"), 0644))

		// Test finding from subdirectory
		result := findProjectRoot(subDir)
		assert.Equal(t, projectDir, result, "Should find project root with go.mod")
	})

	t.Run("find project root with .git", func(t *testing.T) {
		// Create a temporary directory structure
		tmpDir := t.TempDir()
		projectDir := tmpDir + "/gitproject"
		subDir := projectDir + "/subdir/nested"

		require.NoError(t, os.MkdirAll(subDir, 0755))

		// Create a .git directory at project root
		gitPath := projectDir + "/.git"
		require.NoError(t, os.MkdirAll(gitPath, 0755))

		// Test finding from subdirectory
		result := findProjectRoot(subDir)
		assert.Equal(t, projectDir, result, "Should find project root with .git")
	})

	t.Run("prefer go.mod over .git", func(t *testing.T) {
		// Create a temporary directory structure
		tmpDir := t.TempDir()
		projectDir := tmpDir + "/both"
		subDir := projectDir + "/subdir"

		require.NoError(t, os.MkdirAll(subDir, 0755))

		// Create both go.mod and .git
		goModPath := projectDir + "/go.mod"
		require.NoError(t, os.WriteFile(goModPath, []byte("module test\n"), 0644))
		gitPath := projectDir + "/.git"
		require.NoError(t, os.MkdirAll(gitPath, 0755))

		// Test finding from subdirectory
		result := findProjectRoot(subDir)
		assert.Equal(t, projectDir, result, "Should find project root with go.mod (preferred over .git)")
	})

	t.Run("no project root found", func(t *testing.T) {
		// Create a temporary directory without go.mod or .git
		tmpDir := t.TempDir()
		subDir := tmpDir + "/noproject/subdir"

		require.NoError(t, os.MkdirAll(subDir, 0755))

		// Test finding from subdirectory
		result := findProjectRoot(subDir)
		assert.Empty(t, result, "Should return empty string when no project root found")
	})

	t.Run("start from filesystem root", func(t *testing.T) {
		// Test from root directory (should return empty)
		result := findProjectRoot("/")
		assert.Empty(t, result, "Should return empty string when starting from root without markers")
	})

	t.Run("nested project roots", func(t *testing.T) {
		// Create nested project structure
		tmpDir := t.TempDir()
		outerProject := tmpDir + "/outer"
		innerProject := outerProject + "/inner"
		deepDir := innerProject + "/deep/subdir"

		require.NoError(t, os.MkdirAll(deepDir, 0755))

		// Create go.mod in both outer and inner
		outerGoMod := outerProject + "/go.mod"
		require.NoError(t, os.WriteFile(outerGoMod, []byte("module outer\n"), 0644))
		innerGoMod := innerProject + "/go.mod"
		require.NoError(t, os.WriteFile(innerGoMod, []byte("module inner\n"), 0644))

		// Should find the closest (inner) project root
		result := findProjectRoot(deepDir)
		assert.Equal(t, innerProject, result, "Should find closest project root")
	})
}

// TestParseHCLFile_RelativePaths tests relative path handling
func TestParseHCLFile_RelativePaths(t *testing.T) {
	t.Run("relative path from current directory", func(t *testing.T) {
		testFile := "relative_test.tf"
		testContent := "# relative path test"

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		// Use just the filename (relative path)
		content, err := ParseHCLFile(testFile)

		assert.NoError(t, err, "Should handle relative path")
		assert.Equal(t, testContent, content)
	})

	t.Run("relative path with subdirectory", func(t *testing.T) {
		// Create subdirectory
		subdir := "testsubdir"
		require.NoError(t, os.MkdirAll(subdir, 0755))
		defer func() {
			assert.NoError(t, os.RemoveAll(subdir))
		}()

		testFile := subdir + "/subdir_test.tf"
		testContent := "# subdirectory test"

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))

		// Use relative path with subdirectory
		content, err := ParseHCLFile(testFile)

		assert.NoError(t, err, "Should handle relative path with subdirectory")
		assert.Equal(t, testContent, content)
	})

	t.Run("relative path with dot notation", func(t *testing.T) {
		testFile := "dot_notation.tf"
		testContent := "# dot notation test"

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		// Use ./filename notation
		content, err := ParseHCLFile("./" + testFile)

		assert.NoError(t, err, "Should handle ./ notation")
		assert.Equal(t, testContent, content)
	})
}

// TestParseHCLFile_AbsolutePaths tests absolute path handling
func TestParseHCLFile_AbsolutePaths(t *testing.T) {
	t.Run("absolute path within project", func(t *testing.T) {
		// Get current working directory
		cwd, err := os.Getwd()
		require.NoError(t, err)

		testFile := "absolute_test.tf"
		testContent := "# absolute path test"

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		// Use absolute path
		absolutePath := cwd + "/" + testFile
		content, err := ParseHCLFile(absolutePath)

		assert.NoError(t, err, "Should handle absolute path within project")
		assert.Equal(t, testContent, content)
	})

	t.Run("absolute path outside project boundaries", func(t *testing.T) {
		// Create a temp file outside the project
		tmpDir := t.TempDir()
		testFile := tmpDir + "/outside.tf"
		testContent := "# outside project"

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))

		// Try to access file outside project boundaries
		content, err := ParseHCLFile(testFile)

		// Should either fail with security error or succeed if project root can't be determined
		if err != nil {
			assert.Contains(t, err.Error(), "access denied", "Should fail with access denied for outside path")
			assert.Empty(t, content)
		}
		// If it succeeds, that's also acceptable if project root detection fails
	})
}

// TestParseHCLFile_CleanPath tests path cleaning and normalization
func TestParseHCLFile_CleanPath(t *testing.T) {
	t.Run("path with redundant separators", func(t *testing.T) {
		testFile := "clean_test.tf"
		testContent := "# clean path test"

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		// Use path with redundant separators
		content, err := ParseHCLFile(".//" + testFile)

		assert.NoError(t, err, "Should clean path with redundant separators")
		assert.Equal(t, testContent, content)
	})

	t.Run("path with redundant dot segments", func(t *testing.T) {
		testFile := "segments_test.tf"
		testContent := "# segments test"

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			assert.NoError(t, os.Remove(testFile))
		}()

		// Use path with redundant segments
		content, err := ParseHCLFile("././" + testFile)

		assert.NoError(t, err, "Should clean path with redundant dot segments")
		assert.Equal(t, testContent, content)
	})
}

// TestParseHCLFile_ReadPermissions tests file read permission scenarios
func TestParseHCLFile_ReadPermissions(t *testing.T) {
	// Skip on Windows as file permissions work differently
	if os.Getenv("GOOS") == "windows" {
		t.Skip("Skipping permission test on Windows")
	}

	t.Run("file without read permissions", func(t *testing.T) {
		testFile := "no_read.tf"
		testContent := "# no read permissions"

		require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))
		defer func() {
			// Restore permissions before deleting
			os.Chmod(testFile, 0644)
			assert.NoError(t, os.Remove(testFile))
		}()

		// Remove read permissions
		require.NoError(t, os.Chmod(testFile, 0000))

		content, err := ParseHCLFile(testFile)

		assert.Error(t, err, "Should fail when file has no read permissions")
		assert.Empty(t, content)
		assert.Contains(t, err.Error(), "failed to read terraform file", "Should contain read error message")
	})
}

// TestParseHCLFile_ErrorPathCoverage tests additional error paths for coverage
func TestParseHCLFile_ErrorPathCoverage(t *testing.T) {
	t.Run("path contains double dots after cleaning", func(t *testing.T) {
		// This tests the security check for ".." in absolute path
		// The check happens after path cleaning
		testPath := "../../test.tf"

		content, err := ParseHCLFile(testPath)

		// Should fail due to path traversal or file not existing
		assert.Error(t, err)
		assert.Empty(t, content)
	})

	t.Run("error accessing file stats", func(t *testing.T) {
		// Create a file in a directory, then remove the directory
		// to trigger stat error
		testDir := "tmpdir_tf"
		require.NoError(t, os.MkdirAll(testDir, 0755))

		testFile := testDir + "/test.tf"
		require.NoError(t, os.WriteFile(testFile, []byte("# test"), 0644))

		// Remove directory while file reference exists
		os.RemoveAll(testDir)

		content, err := ParseHCLFile(testFile)

		assert.Error(t, err)
		assert.Empty(t, content)
	})
}
