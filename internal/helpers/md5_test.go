package helpers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateMd5(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("Calculate MD5 for non-empty file", func(t *testing.T) {
		content := []byte("Hello, World!")
		filePath := filepath.Join(tempDir, "test1.txt")
		err := os.WriteFile(filePath, content, 0644)
		require.NoError(t, err)

		expectedMD5 := md5.Sum(content)
		expectedMD5String := hex.EncodeToString(expectedMD5[:])

		result, err := CalculateMd5(filePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedMD5String, result)
	})

	t.Run("Calculate MD5 for empty file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "empty.txt")
		err := os.WriteFile(filePath, []byte{}, 0644)
		require.NoError(t, err)

		expectedMD5 := md5.Sum([]byte{})
		expectedMD5String := hex.EncodeToString(expectedMD5[:])

		result, err := CalculateMd5(filePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedMD5String, result)
	})

	t.Run("Calculate MD5 for non-existent file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "non_existent.txt")

		result, err := CalculateMd5(filePath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open file")
		assert.Empty(t, result)
	})

	t.Run("Calculate MD5 for large file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "large.bin")
		f, err := os.Create(filePath)
		require.NoError(t, err)
		defer f.Close()

		// Write 10MB of random data
		data := make([]byte, 1024*1024*10)
		_, err = io.ReadFull(rand.Reader, data)
		require.NoError(t, err)

		_, err = f.Write(data)
		require.NoError(t, err)

		expectedMD5 := md5.Sum(data)
		expectedMD5String := hex.EncodeToString(expectedMD5[:])

		result, err := CalculateMd5(filePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedMD5String, result)
	})
}
