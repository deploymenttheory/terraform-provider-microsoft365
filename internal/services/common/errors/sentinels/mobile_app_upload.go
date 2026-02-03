package sentinels

import "errors"

// Mobile app upload-related sentinel errors (special value used to signify that no
// further processing is possible) for consistent error handling across all mobile
// app resource types (Win32, macOS LOB, macOS DMG, macOS PKG, etc.)
//
// These errors are used during the content file upload workflow which includes:
// - Creating content versions
// - Uploading files to Azure Storage
// - Committing files with encryption metadata
// - Waiting for upload state changes
var (
	// ErrFileStatusFailed indicates that retrieving the file upload status from the API failed
	ErrFileStatusFailed = errors.New("failed to get file status")

	// ErrUploadStateNil indicates that the upload state is unexpectedly nil
	ErrUploadStateNil = errors.New("upload state is nil")

	// ErrAzureStorageURIRequestFailed indicates that the Azure Storage URI request failed
	ErrAzureStorageURIRequestFailed = errors.New("azure storage URI request failed")

	// ErrWaitingForAzureStorageURI indicates that we're waiting for Azure Storage URI to be ready
	ErrWaitingForAzureStorageURI = errors.New("waiting for Azure Storage URI")

	// ErrAzureStorageURINil indicates that the Azure Storage URI is unexpectedly nil
	ErrAzureStorageURINil = errors.New("azure Storage URI is nil")

	// ErrCommitRequestConstruction indicates failure to construct the commit request
	ErrCommitRequestConstruction = errors.New("failed to construct commit request")

	// ErrFileCommitFailed indicates that the file commit operation failed
	ErrFileCommitFailed = errors.New("failed to commit file")

	// ErrReadAfterCreate indicates failure to read resource state after creation
	ErrReadAfterCreate = errors.New("error reading resource state after Create Method")
)
