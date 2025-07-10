package utilityWindowsMSIAppMetadata

// This data source is based on Fleet's implementation of the MSI database parser
// https://github.com/fleetdm/fleet/blob/main/pkg/file/msi.go
// Credit to Fleet for the original implementation

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sassoftware/relic/v8/lib/comdoc"
	"golang.org/x/text/encoding/charmap"
)

// ExtractMSIMetadata extracts metadata from an MSI file (simplified approach based on Fleet's implementation)
func ExtractMSIMetadata(reader io.ReaderAt, fileSize int64) (*MetadataDataSourceModel, error) {

	metadata := &MetadataDataSourceModel{}
	if err := calculateFileChecksums(reader, fileSize, metadata); err != nil {
		return nil, fmt.Errorf("calculating checksums: %w", err)
	}

	doc, err := comdoc.ReadFile(reader)
	if err != nil {
		return nil, fmt.Errorf("reading MSI file: %w", err)
	}
	defer doc.Close()

	entries, err := doc.ListDir(nil)
	if err != nil {
		return nil, fmt.Errorf("listing MSI directory: %w", err)
	}

	// Extract the required tables
	tables := map[string]io.Reader{
		"Table._StringData": nil,
		"Table._StringPool": nil,
		"Table._Columns":    nil,
		"Table.Property":    nil,
		"Table.File":        nil,
		"Table.Feature":     nil,
	}

	for _, entry := range entries {
		if entry.Type != comdoc.DirStream {
			continue
		}

		name := decodeMSIName(entry.Name())
		if _, ok := tables[name]; ok {
			reader, err := doc.ReadStream(entry)
			if err != nil {
				return nil, fmt.Errorf("opening stream %s: %w", name, err)
			}
			tables[name] = reader
		}
	}

	requiredTables := []string{"Table._StringData", "Table._StringPool", "Table._Columns", "Table.Property"}
	for _, tableName := range requiredTables {
		if tables[tableName] == nil {
			return nil, fmt.Errorf("required table %s not found in MSI", tableName)
		}
	}

	allStrings, err := decodeStrings(tables["Table._StringData"], tables["Table._StringPool"])
	if err != nil {
		return nil, fmt.Errorf("decoding strings: %w", err)
	}

	propertyTable, err := parsePropertyTableStructure(tables["Table._Columns"], allStrings)
	if err != nil {
		return nil, fmt.Errorf("parsing property table structure: %w", err)
	}

	properties, err := extractProperties(tables["Table.Property"], propertyTable, allStrings)
	if err != nil {
		return nil, fmt.Errorf("extracting properties: %w", err)
	}

	if err := buildMetadataFromProperties(metadata, properties); err != nil {
		return nil, fmt.Errorf("building metadata: %w", err)
	}

	// Extract files if File table exists
	if tables["Table.File"] != nil {
		files, err := extractFiles(tables["Table.File"], tables["Table._Columns"], allStrings)
		if err == nil { // Don't fail if file extraction fails
			metadata.Files = createStringList(files)
		}
	}

	// Extract features if Feature table exists
	if tables["Table.Feature"] != nil {
		features, err := extractFeatures(tables["Table.Feature"], tables["Table._Columns"], allStrings)
		if err == nil { // Don't fail if feature extraction fails
			metadata.RequiredFeatures = createStringList(features)
		}
	}

	// Set null values for optional fields if they weren't populated
	setDefaultValues(metadata)

	return metadata, nil
}

// calculateFileChecksums calculates SHA256 and MD5 checksums of the file
func calculateFileChecksums(reader io.ReaderAt, fileSize int64, metadata *MetadataDataSourceModel) error {
	sha256Hash := sha256.New()
	md5Hash := md5.New()

	const chunkSize = 64 * 1024 // 64KB chunks
	buffer := make([]byte, chunkSize)
	var offset int64

	for offset < fileSize {
		readSize := chunkSize
		if offset+int64(chunkSize) > fileSize {
			readSize = int(fileSize - offset)
		}

		n, err := reader.ReadAt(buffer[:readSize], offset)
		if err != nil && err != io.EOF {
			return fmt.Errorf("reading file at offset %d: %w", offset, err)
		}

		if n == 0 {
			break
		}

		sha256Hash.Write(buffer[:n])
		md5Hash.Write(buffer[:n])
		offset += int64(n)

		if n < readSize {
			break
		}
	}

	metadata.SHA256Checksum = types.StringValue(fmt.Sprintf("%x", sha256Hash.Sum(nil)))
	metadata.MD5Checksum = types.StringValue(fmt.Sprintf("%x", md5Hash.Sum(nil)))
	metadata.SizeMB = types.Float64Value(float64(fileSize) / (1024 * 1024))

	return nil
}

// buildMetadataFromProperties builds the metadata structure from MSI properties
func buildMetadataFromProperties(metadata *MetadataDataSourceModel, properties map[string]string) error {
	// Decode product name (may be in Windows-1252 encoding)
	productName := properties["ProductName"]
	if decoded, err := charmap.Windows1252.NewDecoder().String(productName); err == nil {
		productName = decoded
	}

	// Core properties
	metadata.ProductCode = types.StringValue(strings.TrimSpace(properties["ProductCode"]))
	metadata.ProductVersion = types.StringValue(strings.TrimSpace(properties["ProductVersion"]))
	metadata.ProductName = types.StringValue(strings.TrimSpace(productName))
	metadata.Publisher = types.StringValue(strings.TrimSpace(properties["Manufacturer"]))

	// Additional properties
	metadata.UpgradeCode = types.StringValue(strings.TrimSpace(properties["UpgradeCode"]))
	metadata.Language = types.StringValue(properties["ProductLanguage"])
	metadata.PackageType = types.StringValue("Application") // Default for MSI
	metadata.InstallLocation = types.StringValue(properties["TARGETDIR"])
	metadata.MinOSVersion = types.StringValue(properties["MinVersion"])

	// Generate install/uninstall commands
	productCode := properties["ProductCode"]
	if productCode != "" {
		metadata.InstallCommand = types.StringValue(fmt.Sprintf("msiexec /i \"%s\" /quiet", productCode))
		metadata.UninstallCommand = types.StringValue(fmt.Sprintf("msiexec /x \"%s\" /quiet", productCode))
	}

	// Determine architecture from Template property
	template := properties["Template"]
	arch := "Unknown"
	if strings.Contains(template, "x64") || strings.Contains(template, "Intel64") {
		arch = "x64"
	} else if strings.Contains(template, "Intel") {
		arch = "x86"
	} else if strings.Contains(template, "Arm64") {
		arch = "ARM64"
	}
	metadata.Architecture = types.StringValue(arch)

	propMap := make(map[string]attr.Value)
	for key, value := range properties {
		propMap[key] = types.StringValue(value)
	}
	metadata.Properties = types.MapValueMust(types.StringType, propMap)

	return nil
}

// extractFiles extracts file information from the File table
func extractFiles(fileReader io.Reader, columnsReader io.Reader, allStrings []string) ([]string, error) {

	fileStructure, err := parseTableStructure(columnsReader, allStrings, "File")
	if err != nil {
		return nil, err
	}

	if fileStructure == nil {
		return []string{}, nil
	}

	// Find FileName column
	fileNameCol := -1
	for i, col := range fileStructure.Columns {
		if col.Name == "FileName" {
			fileNameCol = i
			break
		}
	}

	if fileNameCol == -1 {
		return []string{}, nil
	}

	data, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}

	// Simple parsing - this assumes 2-byte columns (simplified)
	rowSize := len(fileStructure.Columns) * 2
	if rowSize == 0 {
		return []string{}, nil
	}

	rowCount := len(data) / rowSize
	var files []string

	reader := bytes.NewReader(data)
	for row := 0; row < rowCount; row++ {
		// Skip to the filename column
		reader.Seek(int64(row*rowSize+fileNameCol*2), io.SeekStart)

		var stringID uint16
		if err := binary.Read(reader, binary.LittleEndian, &stringID); err != nil {
			continue
		}

		if stringID > 0 && int(stringID) <= len(allStrings) {
			fileName := allStrings[stringID-1]
			if fileName != "" {
				files = append(files, fileName)
			}
		}
	}

	return files, nil
}

// extractFeatures extracts feature information from the Feature table
func extractFeatures(featureReader io.Reader, columnsReader io.Reader, allStrings []string) ([]string, error) {
	featureStructure, err := parseTableStructure(columnsReader, allStrings, "Feature")
	if err != nil {
		return nil, err
	}

	if featureStructure == nil {
		return []string{}, nil
	}

	// Find Feature column
	featureCol := -1
	for i, col := range featureStructure.Columns {
		if col.Name == "Feature" {
			featureCol = i
			break
		}
	}

	if featureCol == -1 {
		return []string{}, nil
	}

	data, err := io.ReadAll(featureReader)
	if err != nil {
		return nil, err
	}

	// Simple parsing - assumes 2-byte columns
	rowSize := len(featureStructure.Columns) * 2
	if rowSize == 0 {
		return []string{}, nil
	}

	rowCount := len(data) / rowSize
	var features []string

	reader := bytes.NewReader(data)
	for row := 0; row < rowCount; row++ {
		reader.Seek(int64(row*rowSize+featureCol*2), io.SeekStart)

		var stringID uint16
		if err := binary.Read(reader, binary.LittleEndian, &stringID); err != nil {
			continue
		}

		if stringID > 0 && int(stringID) <= len(allStrings) {
			featureName := allStrings[stringID-1]
			if featureName != "" {
				features = append(features, featureName)
			}
		}
	}

	return features, nil
}

type msiTable struct {
	Name    string
	Columns []msiColumn
}

type msiColumn struct {
	Number     int
	Name       string
	Attributes uint16
}

func (c msiColumn) Type() msiType {
	if c.Attributes&0x0F00 < 0x800 {
		return msiType(c.Attributes & 0xFFF)
	}
	return msiType(c.Attributes & 0xF00)
}

type msiType uint16

const (
	msiLong            msiType = 0x104
	msiShort           msiType = 0x502
	msiBinary          msiType = 0x900
	msiString          msiType = 0xD00
	msiStringLocalized msiType = 0xF00
	msiUnknown         msiType = 0
)

// extractProperties extracts properties from the Property table (from Fleet)
func extractProperties(propReader io.Reader, table *msiTable, strings []string) (map[string]string, error) {
	if len(table.Columns) != 2 || table.Columns[0].Type() != msiString || table.Columns[1].Type() != msiStringLocalized {
		return nil, errors.New("unexpected Property table structure")
	}

	const propTableRowSize = 4 // 2 uint16s

	b, err := io.ReadAll(propReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read property table: %w", err)
	}

	rowCount := len(b) / propTableRowSize
	propReader = bytes.NewReader(b)

	cols := [][]uint16{
		make([]uint16, 0, rowCount),
		make([]uint16, 0, rowCount),
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < rowCount; j++ {
			var v uint16
			err := binary.Read(propReader, binary.LittleEndian, &v)
			if err != nil {
				return nil, fmt.Errorf("failed to read column %d: %w", i, err)
			}
			cols[i] = append(cols[i], v)
		}
	}

	kv := make(map[string]string, rowCount)
	for i := 0; i < rowCount; i++ {
		if cols[0][i] > 0 && int(cols[0][i]) <= len(strings) &&
			cols[1][i] > 0 && int(cols[1][i]) <= len(strings) {
			kv[strings[cols[0][i]-1]] = strings[cols[1][i]-1]
		}
	}

	return kv, nil
}

// parsePropertyTableStructure parses the Property table structure (from Fleet)
func parsePropertyTableStructure(colReader io.Reader, strings []string) (*msiTable, error) {
	return parseTableStructure(colReader, strings, "Property")
}

// parseTableStructure parses any table structure from the _Columns table
func parseTableStructure(colReader io.Reader, strings []string, targetTable string) (*msiTable, error) {
	const colTableRowSize = 8 // 4 uint16s

	b, err := io.ReadAll(colReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read columns table: %w", err)
	}

	rowCount := len(b) / colTableRowSize
	colReader = bytes.NewReader(b)

	cols := [][]uint16{
		make([]uint16, 0, rowCount),
		make([]uint16, 0, rowCount),
		make([]uint16, 0, rowCount),
		make([]uint16, 0, rowCount),
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < rowCount; j++ {
			var v uint16
			err := binary.Read(colReader, binary.LittleEndian, &v)
			if err != nil {
				return nil, fmt.Errorf("failed to read column %d: %w", i, err)
			}
			cols[i] = append(cols[i], v)
		}
	}

	var tbl msiTable
	for i := 0; i < rowCount; i++ {
		tblID, colNum, colNameID, colAttr := cols[0][i], cols[1][i], cols[2][i], cols[3][i]

		if tblID > 0 && int(tblID) <= len(strings) && colNameID > 0 && int(colNameID) <= len(strings) {
			tableName := strings[tblID-1]
			if tableName == targetTable {
				tbl.Name = tableName
				tbl.Columns = append(tbl.Columns, msiColumn{
					Number:     int(colNum),
					Name:       strings[colNameID-1],
					Attributes: colAttr,
				})
			}
		}
	}

	if tbl.Name == "" {
		return nil, nil // Table not found, but don't error
	}

	return &tbl, nil
}

// decodeStrings decodes strings from MSI string pool (from Fleet)
func decodeStrings(dataReader, poolReader io.Reader) ([]string, error) {
	type header struct {
		Codepage uint16
		Unknown  uint16
	}

	var poolHeader header
	err := binary.Read(poolReader, binary.LittleEndian, &poolHeader)
	if err != nil {
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
		return nil, fmt.Errorf("failed to read pool header: %w", err)
	}

	type entry struct {
		Size     uint16
		RefCount uint16
	}

	var stringEntry entry
	var stringTable []string
	var buf bytes.Buffer

	for {
		err := binary.Read(poolReader, binary.LittleEndian, &stringEntry)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to read pool entry: %w", err)
		}

		stringEntrySize := uint32(stringEntry.Size)

		// Handle large strings
		if stringEntry.Size == 0 && stringEntry.RefCount != 0 {
			err := binary.Read(poolReader, binary.LittleEndian, &stringEntrySize)
			if err != nil {
				return nil, fmt.Errorf("failed to read size of large string: %w", err)
			}
		}

		buf.Reset()
		buf.Grow(int(stringEntrySize))
		_, err = io.CopyN(&buf, dataReader, int64(stringEntrySize))
		if err != nil {
			return nil, fmt.Errorf("failed to read string data: %w", err)
		}
		stringTable = append(stringTable, buf.String())
	}

	return stringTable, nil
}

// decodeMSIName decodes MSI names (from Fleet)
func decodeMSIName(msiName string) string {
	out := ""
	for _, x := range msiName {
		switch {
		case x >= 0x3800 && x < 0x4800:
			x -= 0x3800
			out += string(decodeMSIRune(x&0x3f)) + string(decodeMSIRune(x>>6))
		case x >= 0x4800 && x < 0x4840:
			x -= 0x4800
			out += string(decodeMSIRune(x))
		case x == 0x4840:
			out += "Table."
		default:
			out += string(x)
		}
	}
	return out
}

// decodeMSIRune decodes a single MSI rune (from Fleet)
func decodeMSIRune(x rune) rune {
	if x < 10 {
		return x + '0'
	} else if x < 10+26 {
		return x - 10 + 'A'
	} else if x < 10+26+26 {
		return x - 10 - 26 + 'a'
	} else if x == 10+26+26 {
		return '.'
	}
	return '_'
}

// Helper functions for building Terraform types

// createStringList creates a Terraform list from a string slice
func createStringList(items []string) types.List {
	if len(items) == 0 {
		return types.ListNull(types.StringType)
	}

	elements := make([]attr.Value, len(items))
	for i, item := range items {
		elements[i] = types.StringValue(item)
	}

	return types.ListValueMust(types.StringType, elements)
}

// setDefaultValues sets null values for optional fields that weren't populated
func setDefaultValues(metadata *MetadataDataSourceModel) {
	if metadata.UpgradeCode.IsNull() || metadata.UpgradeCode.ValueString() == "" {
		metadata.UpgradeCode = types.StringNull()
	}
	if metadata.Language.IsNull() || metadata.Language.ValueString() == "" {
		metadata.Language = types.StringNull()
	}
	if metadata.InstallLocation.IsNull() || metadata.InstallLocation.ValueString() == "" {
		metadata.InstallLocation = types.StringNull()
	}
	if metadata.MinOSVersion.IsNull() || metadata.MinOSVersion.ValueString() == "" {
		metadata.MinOSVersion = types.StringNull()
	}
	if metadata.InstallCommand.IsNull() || metadata.InstallCommand.ValueString() == "" {
		metadata.InstallCommand = types.StringNull()
	}
	if metadata.UninstallCommand.IsNull() || metadata.UninstallCommand.ValueString() == "" {
		metadata.UninstallCommand = types.StringNull()
	}

	// Set transform paths to empty list for now
	metadata.TransformPaths = types.ListNull(types.StringType)

	// Set defaults for lists if they're null
	if metadata.Files.IsNull() {
		metadata.Files = types.ListNull(types.StringType)
	}
	if metadata.RequiredFeatures.IsNull() {
		metadata.RequiredFeatures = types.ListNull(types.StringType)
	}
}
