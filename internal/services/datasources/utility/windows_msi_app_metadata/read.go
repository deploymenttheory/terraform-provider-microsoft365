package utilityWindowsMSIAppMetadata

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read fetches and extracts metadata from an MSI file
func (d *WindowsMSIAppMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config WindowsMSIAppMetadataDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	filePathProvided := !config.InstallerFilePathSource.IsNull() && config.InstallerFilePathSource.ValueString() != ""
	urlProvided := !config.InstallerURLSource.IsNull() && config.InstallerURLSource.ValueString() != ""

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with file path provided: %t, URL provided: %t",
		DataSourceName, filePathProvided, urlProvided))

	// Validate inputs - must have either a file path or URL, but not both
	if !filePathProvided && !urlProvided {
		resp.Diagnostics.AddError(
			"Missing Input Parameter",
			"Either installer_file_path_source or installer_url_source must be provided",
		)
		return
	}

	if filePathProvided && urlProvided {
		resp.Diagnostics.AddError(
			"Multiple Input Parameters",
			"Only one of installer_file_path_source or installer_url_source can be provided",
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, config.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Create MSI reader
	var msiReader io.ReaderAt
	var fileSize int64
	var err error

	if filePathProvided {
		msiReader, fileSize, err = d.createFileReader(config.InstallerFilePathSource.ValueString())
	} else {
		msiReader, fileSize, err = d.createURLReader(ctx, config.InstallerURLSource.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading MSI File",
			fmt.Sprintf("Unable to read MSI file: %s", err),
		)
		return
	}

	// Extract metadata from MSI (includes checksum calculation)
	metadata, err := ExtractMSIMetadata(msiReader, fileSize)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Extracting MSI Metadata",
			fmt.Sprintf("Unable to extract metadata from MSI file: %s", err),
		)
		return
	}

	// Create state model
	var state WindowsMSIAppMetadataDataSourceModel
	state.InstallerFilePathSource = config.InstallerFilePathSource
	state.InstallerURLSource = config.InstallerURLSource
	state.Timeouts = config.Timeouts
	state.Metadata = metadata

	// Generate unique ID based on source
	var idSource string
	if filePathProvided {
		idSource = fmt.Sprintf("file:%s", config.InstallerFilePathSource.ValueString())
	} else {
		idSource = fmt.Sprintf("url:%s", config.InstallerURLSource.ValueString())
	}
	hasher := sha256.New()
	hasher.Write([]byte(idSource))
	id := fmt.Sprintf("%x", hasher.Sum(nil))
	state.ID = types.StringValue(id)

	tflog.Debug(ctx, "Successfully extracted MSI metadata")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", DataSourceName))
}

// createFileReader creates a reader for a local file
func (d *WindowsMSIAppMetadataDataSource) createFileReader(filePath string) (io.ReaderAt, int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, fmt.Errorf("opening file %s: %w", filePath, err)
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, 0, fmt.Errorf("getting file info for %s: %w", filePath, err)
	}

	return &fileReaderAt{file: file}, stat.Size(), nil
}

// createURLReader creates a reader for a URL
func (d *WindowsMSIAppMetadataDataSource) createURLReader(ctx context.Context, url string) (io.ReaderAt, int64, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request for %s: %w", url, err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("downloading file from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("HTTP %d when downloading from %s", resp.StatusCode, url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("reading response body from %s: %w", url, err)
	}

	return &bytesReaderAt{data: data}, int64(len(data)), nil
}

// fileReaderAt wraps an os.File to implement io.ReaderAt
type fileReaderAt struct {
	file *os.File
}

func (f *fileReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	return f.file.ReadAt(p, off)
}

func (f *fileReaderAt) Close() error {
	return f.file.Close()
}

// bytesReaderAt wraps a byte slice to implement io.ReaderAt
type bytesReaderAt struct {
	data []byte
}

func (b *bytesReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if off >= int64(len(b.data)) {
		return 0, io.EOF
	}

	n = copy(p, b.data[off:])
	if n < len(p) {
		err = io.EOF
	}

	return n, err
}
