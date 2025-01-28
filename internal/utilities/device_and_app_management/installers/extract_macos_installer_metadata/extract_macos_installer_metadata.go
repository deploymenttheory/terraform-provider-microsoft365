package extract

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

const (
	xarHeaderSize  = 28
	xarHeaderMagic = 0x78617221 // 'xar!'
)

// XAR Header structure
type xarHeader struct {
	Magic         uint32
	Size          uint16
	Version       uint16
	TocLengthZlib uint64
	TocLengthRaw  uint64
	ChecksumType  uint32
}

// XML structures for TOC
type xmlToc struct {
	XMLName xml.Name  `xml:"toc"`
	Files   []xmlFile `xml:"file"`
}

type xmlFile struct {
	Name     string    `xml:"name"`
	Type     string    `xml:"type"`
	Data     *xmlData  `xml:"data"`
	Children []xmlFile `xml:"file"`
}

type xmlData struct {
	Length   int64       `xml:"length"`
	Offset   int64       `xml:"offset"`
	Size     int64       `xml:"size"`
	Encoding xmlEncoding `xml:"encoding"`
}

type xmlEncoding struct {
	Style string `xml:"style,attr"`
}

// PkgReader handles pkg file reading
type PkgReader struct {
	file       *os.File
	fileSize   int64
	header     xarHeader
	heapOffset int64
}

func OpenPkg(filePath string) (*PkgReader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	reader := &PkgReader{
		file:     file,
		fileSize: fileInfo.Size(),
	}

	if err := reader.readHeader(); err != nil {
		file.Close()
		return nil, err
	}

	return reader, nil
}

func (r *PkgReader) Close() error {
	return r.file.Close()
}

func (r *PkgReader) readHeader() error {
	headerData := make([]byte, xarHeaderSize)
	_, err := r.file.ReadAt(headerData, 0)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	r.header.Magic = binary.BigEndian.Uint32(headerData[0:4])
	if r.header.Magic != xarHeaderMagic {
		return fmt.Errorf("invalid XAR magic number")
	}

	r.header.Size = binary.BigEndian.Uint16(headerData[4:6])
	r.header.Version = binary.BigEndian.Uint16(headerData[6:8])
	r.header.TocLengthZlib = binary.BigEndian.Uint64(headerData[8:16])
	r.header.TocLengthRaw = binary.BigEndian.Uint64(headerData[16:24])
	r.header.ChecksumType = binary.BigEndian.Uint32(headerData[24:28])

	r.heapOffset = xarHeaderSize + int64(r.header.TocLengthZlib)
	return nil
}

func (r *PkgReader) readTOC() (*xmlToc, error) {
	// Read compressed TOC
	compressedTOC := make([]byte, r.header.TocLengthZlib)
	_, err := r.file.ReadAt(compressedTOC, xarHeaderSize)
	if err != nil {
		return nil, fmt.Errorf("failed to read TOC: %w", err)
	}

	// Decompress TOC
	zr, err := zlib.NewReader(bytes.NewReader(compressedTOC))
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer zr.Close()

	// Decode XML
	var toc xmlToc
	decoder := xml.NewDecoder(zr)
	decoder.Strict = false
	if err := decoder.Decode(&toc); err != nil {
		return nil, fmt.Errorf("failed to decode TOC XML: %w", err)
	}

	return &toc, nil
}

func (r *PkgReader) findPackageInfo(files []xmlFile) (*xmlFile, error) {
	for _, file := range files {
		if file.Name == "PackageInfo" {
			return &file, nil
		}
		if len(file.Children) > 0 {
			if found, err := r.findPackageInfo(file.Children); err == nil {
				return found, nil
			}
		}
	}
	return nil, fmt.Errorf("PackageInfo not found")
}

func (r *PkgReader) ExtractPackageInfo() ([]byte, error) {
	// Read TOC
	toc, err := r.readTOC()
	if err != nil {
		return nil, err
	}

	// Find PackageInfo file
	pkgInfo, err := r.findPackageInfo(toc.Files)
	if err != nil {
		return nil, err
	}

	// Calculate absolute offset
	offset := r.heapOffset + pkgInfo.Data.Offset

	// Read the raw data
	rawData := make([]byte, pkgInfo.Data.Size)
	_, err = r.file.ReadAt(rawData, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to read PackageInfo data: %w", err)
	}

	// Handle compression if needed
	switch pkgInfo.Data.Encoding.Style {
	case "application/x-gzip":
		zr, err := zlib.NewReader(bytes.NewReader(rawData))
		if err != nil {
			return nil, fmt.Errorf("failed to create zlib reader: %w", err)
		}
		defer zr.Close()
		return io.ReadAll(zr)
	case "application/octet-stream":
		return rawData, nil
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", pkgInfo.Data.Encoding.Style)
	}
}
