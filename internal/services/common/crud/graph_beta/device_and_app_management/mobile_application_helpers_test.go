package helpers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTempFileInfo(t *testing.T) {
	t.Run("TempFileInfo struct initialization", func(t *testing.T) {
		fileInfo := TempFileInfo{
			FilePath:      "/tmp/test.msi",
			ShouldCleanup: true,
		}

		assert.Equal(t, "/tmp/test.msi", fileInfo.FilePath)
		assert.True(t, fileInfo.ShouldCleanup)
	})
}

func TestSetInstallerSourcePath(t *testing.T) {
	ctx := context.Background()

	t.Run("Null metadata object", func(t *testing.T) {
		nullObj := types.ObjectNull(map[string]attr.Type{
			"installer_file_path_source": types.StringType,
			"installer_url_source":       types.StringType,
		})

		path, fileInfo, err := SetInstallerSourcePath(ctx, nullObj)

		assert.NoError(t, err)
		assert.Empty(t, path)
		assert.Equal(t, TempFileInfo{}, fileInfo)
	})

	t.Run("Unknown installer sources during plan", func(t *testing.T) {
		// Create metadata with unknown values
		metadataValue, diags := types.ObjectValue(
			map[string]attr.Type{
				"installer_file_path_source": types.StringType,
				"installer_url_source":       types.StringType,
			},
			map[string]attr.Value{
				"installer_file_path_source": types.StringUnknown(),
				"installer_url_source":       types.StringUnknown(),
			},
		)
		require.False(t, diags.HasError())

		path, fileInfo, err := SetInstallerSourcePath(ctx, metadataValue)

		assert.NoError(t, err)
		assert.Empty(t, path)
		assert.Equal(t, TempFileInfo{}, fileInfo)
	})

	t.Run("Local file path provided", func(t *testing.T) {
		// Create a temporary test file
		tempFile, err := os.CreateTemp("", "test-installer-*.msi")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		metadataValue, diags := types.ObjectValue(
			map[string]attr.Type{
				"installer_file_path_source": types.StringType,
				"installer_url_source":       types.StringType,
			},
			map[string]attr.Value{
				"installer_file_path_source": types.StringValue(tempFile.Name()),
				"installer_url_source":       types.StringNull(),
			},
		)
		require.False(t, diags.HasError())

		path, fileInfo, err := SetInstallerSourcePath(ctx, metadataValue)

		assert.NoError(t, err)
		assert.Equal(t, tempFile.Name(), path)
		assert.Equal(t, tempFile.Name(), fileInfo.FilePath)
		assert.False(t, fileInfo.ShouldCleanup)
	})

	t.Run("Empty file path falls back to URL", func(t *testing.T) {
		// Create a test HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("fake installer content"))
		}))
		defer server.Close()

		metadataValue, diags := types.ObjectValue(
			map[string]attr.Type{
				"installer_file_path_source": types.StringType,
				"installer_url_source":       types.StringType,
			},
			map[string]attr.Value{
				"installer_file_path_source": types.StringValue(""),
				"installer_url_source":       types.StringValue(server.URL + "/test-installer.msi"),
			},
		)
		require.False(t, diags.HasError())

		path, fileInfo, err := SetInstallerSourcePath(ctx, metadataValue)

		assert.NoError(t, err)
		assert.NotEmpty(t, path)
		assert.NotEmpty(t, fileInfo.FilePath)
		assert.True(t, fileInfo.ShouldCleanup)

		// Verify the downloaded file exists
		_, err = os.Stat(path)
		assert.NoError(t, err)

		// Clean up the downloaded file
		os.Remove(path)
	})

	t.Run("URL download with null file path", func(t *testing.T) {
		// Create a test HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("fake installer content"))
		}))
		defer server.Close()

		metadataValue, diags := types.ObjectValue(
			map[string]attr.Type{
				"installer_file_path_source": types.StringType,
				"installer_url_source":       types.StringType,
			},
			map[string]attr.Value{
				"installer_file_path_source": types.StringNull(),
				"installer_url_source":       types.StringValue(server.URL + "/test-installer.msi"),
			},
		)
		require.False(t, diags.HasError())

		path, fileInfo, err := SetInstallerSourcePath(ctx, metadataValue)

		assert.NoError(t, err)
		assert.NotEmpty(t, path)
		assert.NotEmpty(t, fileInfo.FilePath)
		assert.True(t, fileInfo.ShouldCleanup)

		// Verify the downloaded file exists
		_, err = os.Stat(path)
		assert.NoError(t, err)

		// Clean up the downloaded file
		os.Remove(path)
	})

	t.Run("URL download fails", func(t *testing.T) {
		metadataValue, diags := types.ObjectValue(
			map[string]attr.Type{
				"installer_file_path_source": types.StringType,
				"installer_url_source":       types.StringType,
			},
			map[string]attr.Value{
				"installer_file_path_source": types.StringNull(),
				"installer_url_source":       types.StringValue("http://invalid-url-that-does-not-exist.com/installer.msi"),
			},
		)
		require.False(t, diags.HasError())

		path, fileInfo, err := SetInstallerSourcePath(ctx, metadataValue)

		assert.Error(t, err)
		assert.Empty(t, path)
		assert.Equal(t, TempFileInfo{}, fileInfo)
		assert.Contains(t, err.Error(), "failed to download installer file")
	})

	t.Run("No installer sources provided", func(t *testing.T) {
		metadataValue, diags := types.ObjectValue(
			map[string]attr.Type{
				"installer_file_path_source": types.StringType,
				"installer_url_source":       types.StringType,
			},
			map[string]attr.Value{
				"installer_file_path_source": types.StringNull(),
				"installer_url_source":       types.StringNull(),
			},
		)
		require.False(t, diags.HasError())

		path, fileInfo, err := SetInstallerSourcePath(ctx, metadataValue)

		assert.Error(t, err)
		assert.Empty(t, path)
		assert.Equal(t, TempFileInfo{}, fileInfo)
		assert.Contains(t, err.Error(), "installer file not provided")
	})

	t.Run("Empty installer sources", func(t *testing.T) {
		metadataValue, diags := types.ObjectValue(
			map[string]attr.Type{
				"installer_file_path_source": types.StringType,
				"installer_url_source":       types.StringType,
			},
			map[string]attr.Value{
				"installer_file_path_source": types.StringValue(""),
				"installer_url_source":       types.StringValue(""),
			},
		)
		require.False(t, diags.HasError())

		path, fileInfo, err := SetInstallerSourcePath(ctx, metadataValue)

		assert.Error(t, err)
		assert.Empty(t, path)
		assert.Equal(t, TempFileInfo{}, fileInfo)
		assert.Contains(t, err.Error(), "installer file not provided")
	})

	t.Run("File path takes precedence over URL", func(t *testing.T) {
		// Create a temporary test file
		tempFile, err := os.CreateTemp("", "test-installer-*.msi")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		// Create a test HTTP server (should not be called)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("Server should not be called when file path is provided")
		}))
		defer server.Close()

		metadataValue, diags := types.ObjectValue(
			map[string]attr.Type{
				"installer_file_path_source": types.StringType,
				"installer_url_source":       types.StringType,
			},
			map[string]attr.Value{
				"installer_file_path_source": types.StringValue(tempFile.Name()),
				"installer_url_source":       types.StringValue(server.URL + "/should-not-be-downloaded.msi"),
			},
		)
		require.False(t, diags.HasError())

		path, fileInfo, err := SetInstallerSourcePath(ctx, metadataValue)

		assert.NoError(t, err)
		assert.Equal(t, tempFile.Name(), path)
		assert.Equal(t, tempFile.Name(), fileInfo.FilePath)
		assert.False(t, fileInfo.ShouldCleanup)
	})
}

func TestCleanupTempFile(t *testing.T) {
	ctx := context.Background()

	t.Run("Clean up file that should be cleaned up", func(t *testing.T) {
		// Create a temporary file
		tempFile, err := os.CreateTemp("", "test-cleanup-*.tmp")
		require.NoError(t, err)
		tempFile.Close()

		// Verify file exists
		_, err = os.Stat(tempFile.Name())
		require.NoError(t, err)

		fileInfo := TempFileInfo{
			FilePath:      tempFile.Name(),
			ShouldCleanup: true,
		}

		CleanupTempFile(ctx, fileInfo)

		// Verify file is removed
		_, err = os.Stat(tempFile.Name())
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("Do not clean up file that should not be cleaned up", func(t *testing.T) {
		// Create a temporary file
		tempFile, err := os.CreateTemp("", "test-no-cleanup-*.tmp")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		// Verify file exists
		_, err = os.Stat(tempFile.Name())
		require.NoError(t, err)

		fileInfo := TempFileInfo{
			FilePath:      tempFile.Name(),
			ShouldCleanup: false,
		}

		CleanupTempFile(ctx, fileInfo)

		// Verify file still exists
		_, err = os.Stat(tempFile.Name())
		assert.NoError(t, err)
	})

	t.Run("Handle empty file path", func(t *testing.T) {
		fileInfo := TempFileInfo{
			FilePath:      "",
			ShouldCleanup: true,
		}

		// Should not panic
		CleanupTempFile(ctx, fileInfo)
	})

	t.Run("Handle non-existent file", func(t *testing.T) {
		fileInfo := TempFileInfo{
			FilePath:      "/path/to/non/existent/file.tmp",
			ShouldCleanup: true,
		}

		// Should not panic and should handle error gracefully
		CleanupTempFile(ctx, fileInfo)
	})

	t.Run("Handle file with insufficient permissions", func(t *testing.T) {
		// Create a temporary file in a directory
		tempDir, err := os.MkdirTemp("", "test-cleanup-dir-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		tempFile := filepath.Join(tempDir, "test-file.tmp")
		err = os.WriteFile(tempFile, []byte("test"), 0644)
		require.NoError(t, err)

		// Remove write permissions from directory (on Unix systems)
		if strings.Contains(strings.ToLower(os.Getenv("GOOS")), "windows") {
			t.Skip("Skipping permission test on Windows")
		}

		err = os.Chmod(tempDir, 0444) // Read-only
		require.NoError(t, err)
		defer os.Chmod(tempDir, 0755) // Restore permissions for cleanup

		fileInfo := TempFileInfo{
			FilePath:      tempFile,
			ShouldCleanup: true,
		}

		// Should handle permission error gracefully (logs warning but doesn't panic)
		CleanupTempFile(ctx, fileInfo)

		// File should still exist due to permission error
		_, err = os.Stat(tempFile)
		if !os.IsNotExist(err) {
			// File exists, which is expected due to permission error
			assert.NoError(t, err)
		}
	})
}

// Helper function tests
func TestSetInstallerSourcePathIntegration(t *testing.T) {
	ctx := context.Background()

	t.Run("Full workflow with local file", func(t *testing.T) {
		// Create a temporary test file
		tempFile, err := os.CreateTemp("", "integration-test-*.msi")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		content := "fake MSI content for testing"
		_, err = tempFile.WriteString(content)
		require.NoError(t, err)
		tempFile.Close()

		// Create metadata object
		metadataValue, diags := types.ObjectValue(
			map[string]attr.Type{
				"installer_file_path_source": types.StringType,
				"installer_url_source":       types.StringType,
			},
			map[string]attr.Value{
				"installer_file_path_source": types.StringValue(tempFile.Name()),
				"installer_url_source":       types.StringNull(),
			},
		)
		require.False(t, diags.HasError())

		// Test the function
		path, fileInfo, err := SetInstallerSourcePath(ctx, metadataValue)
		require.NoError(t, err)

		// Verify results
		assert.Equal(t, tempFile.Name(), path)
		assert.Equal(t, tempFile.Name(), fileInfo.FilePath)
		assert.False(t, fileInfo.ShouldCleanup)

		// Verify file still exists and has correct content
		readContent, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, content, string(readContent))

		// Cleanup should be no-op for local files
		CleanupTempFile(ctx, fileInfo)

		// File should still exist
		_, err = os.Stat(tempFile.Name())
		assert.NoError(t, err)
	})

	t.Run("Full workflow with URL download and cleanup", func(t *testing.T) {
		content := "fake installer content from URL"

		// Create a test HTTP server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", "attachment; filename=test-installer.msi")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(content))
		}))
		defer server.Close()

		// Create metadata object
		metadataValue, diags := types.ObjectValue(
			map[string]attr.Type{
				"installer_file_path_source": types.StringType,
				"installer_url_source":       types.StringType,
			},
			map[string]attr.Value{
				"installer_file_path_source": types.StringNull(),
				"installer_url_source":       types.StringValue(server.URL + "/test-installer.msi"),
			},
		)
		require.False(t, diags.HasError())

		// Test the function
		path, fileInfo, err := SetInstallerSourcePath(ctx, metadataValue)
		require.NoError(t, err)

		// Verify results
		assert.NotEmpty(t, path)
		assert.Equal(t, path, fileInfo.FilePath)
		assert.True(t, fileInfo.ShouldCleanup)

		// Verify downloaded file exists and has correct content
		readContent, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, content, string(readContent))

		// Test cleanup
		CleanupTempFile(ctx, fileInfo)

		// File should no longer exist
		_, err = os.Stat(path)
		assert.True(t, os.IsNotExist(err))
	})
}

// Benchmark tests
func BenchmarkSetInstallerSourcePathLocal(b *testing.B) {
	ctx := context.Background()

	// Create a temporary test file
	tempFile, err := os.CreateTemp("", "benchmark-test-*.msi")
	require.NoError(b, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	metadataValue, diags := types.ObjectValue(
		map[string]attr.Type{
			"installer_file_path_source": types.StringType,
			"installer_url_source":       types.StringType,
		},
		map[string]attr.Value{
			"installer_file_path_source": types.StringValue(tempFile.Name()),
			"installer_url_source":       types.StringNull(),
		},
	)
	require.False(b, diags.HasError())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := SetInstallerSourcePath(ctx, metadataValue)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCleanupTempFile(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Create a temporary file for each iteration
		tempFile, err := os.CreateTemp("", "benchmark-cleanup-*.tmp")
		if err != nil {
			b.Fatal(err)
		}
		tempFile.Close()

		fileInfo := TempFileInfo{
			FilePath:      tempFile.Name(),
			ShouldCleanup: true,
		}
		b.StartTimer()

		CleanupTempFile(ctx, fileInfo)
	}
}
