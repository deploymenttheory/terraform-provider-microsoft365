package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	khttp "github.com/microsoft/kiota-http-go"
)

// Keep the SDK's compression and 415 fallback, but make the fallback body
// seekable for the downstream SDK retry handler. Kiota 1.5.6 otherwise replaces
// it with io.NopCloser(bytes.Buffer), which cannot be rewound after a 429.
type replayableCompressionHandler struct {
	*khttp.CompressionHandler
}

//nolint:wrapcheck // Preserve SDK errors for the shared Graph error handler.
func (h *replayableCompressionHandler) Intercept(pipeline khttp.Pipeline, index int, req *http.Request) (*http.Response, error) {
	return h.CompressionHandler.Intercept(&compressionFallbackPipeline{Pipeline: pipeline}, index, req)
}

type compressionFallbackPipeline struct {
	khttp.Pipeline
	forwarded bool
}

//nolint:wrapcheck // Preserve errors returned by the SDK pipeline.
func (p *compressionFallbackPipeline) Next(req *http.Request, index int) (*http.Response, error) {
	// Compression calls Next a second time only for its uncompressed 415 fallback.
	// The SDK has already buffered this payload; do not buffer ordinary streams.
	if p.forwarded && req.Body != nil {
		if _, seekable := req.Body.(io.Seeker); !seekable {
			body, err := io.ReadAll(req.Body)
			if err != nil {
				return nil, fmt.Errorf("read uncompressed Graph request for retry: %w", err)
			}
			if err := req.Body.Close(); err != nil {
				return nil, fmt.Errorf("close uncompressed Graph request body: %w", err)
			}
			req.Body = khttp.NopCloser(bytes.NewReader(body))
		}
	}
	p.forwarded = true
	return p.Pipeline.Next(req, index)
}
