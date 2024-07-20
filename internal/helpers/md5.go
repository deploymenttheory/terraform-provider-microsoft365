package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// CalculateMd5 calculates the MD5 checksum of a file.
func CalculateMd5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file content: %w", err)
	}

	checksum := hash.Sum(nil)
	md5String := hex.EncodeToString(checksum)

	return md5String, nil
}
