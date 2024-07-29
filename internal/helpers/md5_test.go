package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateMd5(t *testing.T) {
	// Helper function to create a temporary file with given content
	createTempFile := func(t *testing.T, content string) (string, func()) {
		tmpfile, err := os.CreateTemp("", "testfile")
		require.NoError(t, err)

		_, err = tmpfile.Write([]byte(content))
		require.NoError(t, err)

		err = tmpfile.Close()
		require.NoError(t, err)

		return tmpfile.Name(), func() { os.Remove(tmpfile.Name()) }
	}

	t.Run("ValidFile", func(t *testing.T) {
		content := "Hello, World!"
		expectedMd5 := md5.Sum([]byte(content))
		expectedMd5String := hex.EncodeToString(expectedMd5[:])

		filePath, cleanup := createTempFile(t, content)
		defer cleanup()

		result, err := CalculateMd5(filePath)
		require.NoError(t, err)
		assert.Equal(t, expectedMd5String, result, "Expected MD5 checksum does not match")
	})

	t.Run("FileNotFound", func(t *testing.T) {
		_, err := CalculateMd5("nonexistentfile.txt")
		assert.Error(t, err, "Expected error due to non-existent file")
		assert.Contains(t, err.Error(), "failed to open file")
	})

	t.Run("ReadError", func(t *testing.T) {
		// Create a file and remove read permissions
		content := "Read error test"
		filePath, cleanup := createTempFile(t, content)
		defer cleanup()

		err := os.Chmod(filePath, 0200) // Write-only permission
		require.NoError(t, err)

		_, err = CalculateMd5(filePath)
		assert.Error(t, err, "Expected error due to read permission denied")
		assert.Contains(t, err.Error(), "failed to hash file content")
	})
}
