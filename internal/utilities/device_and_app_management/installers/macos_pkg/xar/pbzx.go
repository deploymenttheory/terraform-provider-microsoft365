package xar

import (
	"bytes"
	"container/heap"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/xi2/xz"
)

const (
	XBSZ = 4 * 1024    // Block size for reading (4KB)
	ZBSZ = 1024 * XBSZ // Buffer size for decompressed data (4MB)
)

// Chunk represents a single PBZX chunk with metadata
type Chunk struct {
	index    int    // Chunk sequence number
	inflated int    // Expected size after inflation
	data     []byte // Raw chunk data
	result   []byte // Inflated result (if processed)
}

// ChunkHeap implements heap.Interface for ordered chunk writing
type ChunkHeap []*Chunk

func (h ChunkHeap) Len() int            { return len(h) }
func (h ChunkHeap) Less(i, j int) bool  { return h[i].index < h[j].index }
func (h ChunkHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *ChunkHeap) Push(x interface{}) { *h = append(*h, x.(*Chunk)) }
func (h *ChunkHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// ExtractOptions configures the PBZX extraction process
type ExtractOptions struct {
	NumWorkers int  // Number of parallel decompression workers
	CPIOCheck  bool // Whether to verify CPIO header
}

func extractPBZXFiles(ctx context.Context, reader io.Reader) (io.Reader, error) {
	tflog.Debug(ctx, "Starting PBZX extraction")

	// Create pipe to convert our writer-based process into a reader
	pipeReader, pipeWriter := io.Pipe()

	// Start extraction in a goroutine
	go func() {
		// Always close the writer when we're done
		defer pipeWriter.Close()

		opts := ExtractOptions{
			NumWorkers: runtime.NumCPU(),
			CPIOCheck:  false, // Caller will do CPIO check
		}

		// If extraction fails, close the pipe with an error
		if err := extractPBZXFilesWithworkers(ctx, reader, pipeWriter, opts); err != nil {
			pipeWriter.CloseWithError(fmt.Errorf("PBZX extraction failed: %w", err))
			return
		}
	}()

	return pipeReader, nil
}

// extractPBZXFiles processes a PBZX formatted payload and streams the decompressed content
func extractPBZXFilesWithworkers(ctx context.Context, reader io.Reader, writer io.Writer, opts ExtractOptions) error {
	if opts.NumWorkers == 0 {
		opts.NumWorkers = runtime.NumCPU()
	}

	tflog.Debug(ctx, "Starting PBZX extraction", map[string]interface{}{
		"workers": opts.NumWorkers,
	})

	// Validate PBZX header
	if err := validateHeader(ctx, reader); err != nil {
		return err
	}

	// Setup channels for processing pipeline
	errCh := make(chan error, 1)
	inflateCh := make(chan *Chunk, opts.NumWorkers*2)
	writeCh := make(chan *Chunk, opts.NumWorkers*2)

	// Create cancellable context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Error handler function
	handleError := func(err error) {
		select {
		case errCh <- err:
		default:
		}
		cancel()
	}

	var wg sync.WaitGroup

	// Start chunk reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(inflateCh)
		if err := readChunks(ctx, reader, inflateCh, writeCh); err != nil {
			handleError(err)
		}
	}()

	// Start decompression workers
	wg.Add(opts.NumWorkers)
	for i := 0; i < opts.NumWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			if err := inflateChunks(ctx, inflateCh, writeCh, workerID); err != nil {
				handleError(err)
			}
		}(i)
	}

	// Start ordered writer goroutine
	var writeWg sync.WaitGroup
	writeWg.Add(1)
	go func() {
		defer writeWg.Done()
		if err := writeOrderedChunks(ctx, writeCh, writer, opts.CPIOCheck); err != nil {
			handleError(err)
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	close(writeCh)
	writeWg.Wait()

	// Check for errors
	select {
	case err := <-errCh:
		return fmt.Errorf("PBZX extraction failed: %w", err)
	default:
		tflog.Info(ctx, "PBZX extraction completed successfully")
		return nil
	}
}

func validateHeader(ctx context.Context, reader io.Reader) error {
	magic := make([]byte, 4)
	if _, err := io.ReadFull(reader, magic); err != nil {
		return fmt.Errorf("failed to read PBZX magic number: %w", err)
	}

	if !bytes.Equal(magic, []byte("pbzx")) {
		return fmt.Errorf("invalid PBZX magic number: expected 'pbzx', got %x", magic)
	}

	// Read initial flags
	flags := make([]byte, 8)
	if _, err := io.ReadFull(reader, flags); err != nil {
		return fmt.Errorf("failed to read PBZX flags: %w", err)
	}

	return nil
}

func readChunks(ctx context.Context, reader io.Reader, inflateCh, writeCh chan<- *Chunk) error {
	var chunkIndex int

	for {
		// Read chunk metadata
		var inflateSize, deflateSize uint64
		if err := binary.Read(reader, binary.BigEndian, &inflateSize); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to read chunk size: %w", err)
		}

		if err := binary.Read(reader, binary.BigEndian, &deflateSize); err != nil {
			return fmt.Errorf("failed to read chunk deflate size: %w", err)
		}

		// Sanity checks
		if uint64(int(inflateSize)) != inflateSize {
			return fmt.Errorf("chunk size too large: %d", inflateSize)
		}

		// Read chunk data
		data := make([]byte, deflateSize)
		if _, err := io.ReadFull(reader, data); err != nil {
			return fmt.Errorf("failed to read chunk data: %w", err)
		}

		chunk := &Chunk{
			index:    chunkIndex,
			inflated: int(inflateSize),
			data:     data,
		}

		// Route chunk based on compression
		switch {
		case deflateSize < inflateSize:
			// Compressed chunk
			select {
			case <-ctx.Done():
				return ctx.Err()
			case inflateCh <- chunk:
			}
		case deflateSize == inflateSize:
			// Uncompressed chunk
			chunk.result = data
			select {
			case <-ctx.Done():
				return ctx.Err()
			case writeCh <- chunk:
			}
		default:
			return fmt.Errorf("invalid chunk sizes: deflate %d > inflate %d", deflateSize, inflateSize)
		}

		chunkIndex++
	}
}

func inflateChunks(ctx context.Context, inflateCh <-chan *Chunk, writeCh chan<- *Chunk, workerID int) error {
	tflog.Debug(ctx, "Starting decompression worker", map[string]interface{}{
		"worker_id": workerID,
	})

	for chunk := range inflateCh {
		// Verify XZ header
		if !bytes.HasPrefix(chunk.data, []byte{0xfd, '7', 'z', 'X', 'Z', 0x00}) {
			return fmt.Errorf("invalid XZ header in chunk %d", chunk.index)
		}

		// Create XZ reader
		xzr, err := xz.NewReader(bytes.NewReader(chunk.data), 0)
		if err != nil {
			return fmt.Errorf("failed to create XZ reader for chunk %d: %w", chunk.index, err)
		}
		xzr.Multistream(false)

		// Decompress chunk
		result := make([]byte, chunk.inflated)
		n, err := io.ReadFull(xzr, result)
		if err != nil && err != io.ErrUnexpectedEOF {
			return fmt.Errorf("decompression failed for chunk %d: %w", chunk.index, err)
		}

		chunk.result = result[:n]
		chunk.data = nil // Free the compressed data

		select {
		case <-ctx.Done():
			return ctx.Err()
		case writeCh <- chunk:
		}
	}

	return nil
}

func writeOrderedChunks(ctx context.Context, writeCh <-chan *Chunk, writer io.Writer, cpioCheck bool) error {
	var nextIndex int
	h := &ChunkHeap{}
	heap.Init(h)

	var extractedFiles []string // List to track extracted file names

	for chunk := range writeCh {
		heap.Push(h, chunk)

		// Write chunks in order
		for h.Len() > 0 && (*h)[0].index == nextIndex {
			chunk := heap.Pop(h).(*Chunk)

			// Check CPIO header if this is the first chunk
			if nextIndex == 0 && cpioCheck {
				if !bytes.HasPrefix(chunk.result, []byte("070701")) {
					return fmt.Errorf("invalid CPIO header")
				}
			}

			// Extract filenames from CPIO archive (if applicable)
			if cpioCheck {
				files, err := extractCPIOFilenames(chunk.result)
				if err == nil {
					extractedFiles = append(extractedFiles, files...)
				}
			}

			if _, err := writer.Write(chunk.result); err != nil {
				return fmt.Errorf("failed to write chunk %d: %w", chunk.index, err)
			}

			nextIndex++
		}
	}

	// Write any remaining chunks
	for h.Len() > 0 {
		chunk := heap.Pop(h).(*Chunk)
		if _, err := writer.Write(chunk.result); err != nil {
			return fmt.Errorf("failed to write remaining chunk %d: %w", chunk.index, err)
		}
	}

	// Log extracted files
	if len(extractedFiles) > 0 {
		tflog.Info(ctx, "Extracted files", map[string]interface{}{
			"files": extractedFiles,
		})
	}

	return nil
}

func extractCPIOFilenames(data []byte) ([]string, error) {
	var filenames []string
	reader := bytes.NewReader(data)

	for {
		header := make([]byte, 110) // CPIO new ASCII header is 110 bytes
		_, err := io.ReadFull(reader, header)
		if err != nil {
			break
		}

		// Read filename length
		nameLenHex := string(header[94:102]) // Offset 94-101 (8 chars for filename length)
		var nameLen int
		_, err = fmt.Sscanf(nameLenHex, "%x", &nameLen)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CPIO filename length: %w", err)
		}

		// Read filename
		name := make([]byte, nameLen)
		_, err = io.ReadFull(reader, name)
		if err != nil {
			break
		}
		filenames = append(filenames, string(bytes.Trim(name, "\x00")))

		// Skip file data (size from header)
		fileSizeHex := string(header[54:62]) // Offset 54-61 (8 chars for file size)
		var fileSize int
		_, err = fmt.Sscanf(fileSizeHex, "%x", &fileSize)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CPIO file size: %w", err)
		}

		// Align to next header (CPIO aligns to 4-byte boundaries)
		offset := fileSize
		if offset%4 != 0 {
			offset += 4 - (offset % 4)
		}
		_, _ = reader.Seek(int64(offset), io.SeekCurrent)
	}

	return filenames, nil
}
