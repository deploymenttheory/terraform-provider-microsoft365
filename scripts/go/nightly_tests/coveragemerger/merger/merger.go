package merger

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	// CoverageFileExt is the file extension for coverage files.
	CoverageFileExt = ".txt"

	// CoverageMode is the coverage mode used in Go coverage files.
	CoverageMode = "mode: atomic"

	// DefaultFilePermission is the default file permission for created files.
	DefaultFilePermission = 0644
)

// Merger merges multiple Go coverage files into a single output file.
type Merger struct {
	inputDir   string
	outputPath string
	stdout     io.Writer
}

// New creates a new coverage merger.
func New(inputDir, outputPath string) *Merger {
	return &Merger{
		inputDir:   inputDir,
		outputPath: outputPath,
		stdout:     os.Stdout,
	}
}

// WithOutput sets a custom stdout writer.
func (m *Merger) WithOutput(w io.Writer) *Merger {
	m.stdout = w
	return m
}

// Merge performs the coverage file merge operation.
func (m *Merger) Merge() error {
	coverageFiles, err := m.findCoverageFiles()
	if err != nil {
		return fmt.Errorf("failed to find coverage files: %w", err)
	}

	if len(coverageFiles) == 0 {
		fmt.Fprintln(m.stdout, "⚠️  No coverage files found")
		return m.createEmptyOutput()
	}

	fmt.Fprintf(m.stdout, "Found %d coverage files to merge\n", len(coverageFiles))

	outputFile, err := m.createOutputFile()
	if err != nil {
		return err
	}
	defer outputFile.Close()

	totalLines, err := m.mergeCoverageData(outputFile, coverageFiles)
	if err != nil {
		return err
	}

	fmt.Fprintf(m.stdout, "Total coverage lines: %d\n", totalLines)
	fmt.Fprintf(m.stdout, "✅ Successfully merged coverage files to %s\n", m.outputPath)

	return nil
}

// findCoverageFiles recursively finds all coverage files in the input directory.
func (m *Merger) findCoverageFiles() ([]string, error) {
	var files []string

	err := filepath.WalkDir(m.inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, CoverageFileExt) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// createOutputFile creates and initializes the output coverage file.
func (m *Merger) createOutputFile() (*os.File, error) {
	f, err := os.Create(m.outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}

	if _, err := f.WriteString(CoverageMode + "\n"); err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to write coverage mode: %w", err)
	}

	return f, nil
}

// mergeCoverageData merges coverage data from all input files into the output file.
func (m *Merger) mergeCoverageData(output *os.File, files []string) (int, error) {
	totalLines := 0

	for _, filePath := range files {
		lines, err := m.appendCoverageFile(output, filePath)
		if err != nil {
			return totalLines, fmt.Errorf("failed to merge %s: %w", filePath, err)
		}

		totalLines += lines
		fmt.Fprintf(m.stdout, "  Merged %s (%d lines)\n", filepath.Base(filePath), lines)
	}

	return totalLines, nil
}

// appendCoverageFile appends coverage data from a single file, skipping the mode line.
func (m *Merger) appendCoverageFile(output *os.File, path string) (int, error) {
	input, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	lineCount := 0
	firstLine := true

	for scanner.Scan() {
		line := scanner.Text()

		// Skip the mode line (first line)
		if firstLine {
			firstLine = false
			if strings.HasPrefix(line, "mode:") {
				continue
			}
		}

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		if _, err := output.WriteString(line + "\n"); err != nil {
			return lineCount, fmt.Errorf("failed to write line: %w", err)
		}

		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return lineCount, fmt.Errorf("error reading file: %w", err)
	}

	return lineCount, nil
}

// createEmptyOutput creates an empty coverage file with just the mode line.
func (m *Merger) createEmptyOutput() error {
	return os.WriteFile(m.outputPath, []byte(CoverageMode+"\n"), DefaultFilePermission)
}
