package xar

import (
	"io"
	"os"
)

// TempFileReader is an io.Reader with all extra io interfaces supported by a
// file on disk reader (e.g. io.ReaderAt, io.Seeker, etc.). When created with
// NewTempFileReader, it is backed by a temporary file on disk, and that file
// is deleted when Close is called.
type TempFileReader struct {
	*os.File
	keepFile bool
}

// Rewind seeks to the beginning of the file so the next read will read from
// the start of the bytes.
func (r *TempFileReader) Rewind() error {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return err
	}
	return nil
}

// Close closes the TempFileReader and deletes the underlying temp file unless
// it was instructed not to do so at creation time.
func (r *TempFileReader) Close() error {
	cerr := r.File.Close()
	var rerr error
	if !r.keepFile {
		rerr = os.Remove(r.File.Name())
	}
	if cerr != nil {
		return cerr
	}
	return rerr
}

// NewKeepFileReader creates a TempFileReader from a file path and keeps the
// file on Close, instead of deleting it.
func NewKeepFileReader(filename string) (*TempFileReader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &TempFileReader{File: f, keepFile: true}, nil
}

// NewTempFileReader creates a temp file to store the data from the provided
// reader and returns a TempFileReader that reads from that temp file, deleting
// it on close.
func NewTempFileReader(from io.Reader, tempDirFn func() string) (*TempFileReader, error) {
	if tempDirFn == nil {
		tempDirFn = os.TempDir
	}

	tempFile, err := os.CreateTemp(tempDirFn(), "terraform-provider-microsoft365-temp-file-*")
	if err != nil {
		return nil, err
	}
	tfr := &TempFileReader{File: tempFile}

	if _, err := io.Copy(tempFile, from); err != nil {
		_ = tfr.Close() // best-effort close/delete
		return nil, err
	}
	if err := tfr.Rewind(); err != nil {
		_ = tfr.Close() // best-effort close/delete
		return nil, err
	}
	return tfr, nil
}
