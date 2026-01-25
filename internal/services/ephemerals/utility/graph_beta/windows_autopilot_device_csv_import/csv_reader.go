package windowsAutopilotDeviceCSVImport

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Expected CSV headers
const (
	HeaderSerialNumber   = "Device Serial Number"
	HeaderWindowsProduct = "Windows Product ID"
	HeaderHardwareHash   = "Hardware Hash"
	HeaderGroupTag       = "Group Tag"
	HeaderAssignedUser   = "Assigned User"
)

// readCSVFile reads a CSV file and returns a slice of DeviceEntry
func readCSVFile(ctx context.Context, filePath string) ([]DeviceEntry, diag.Diagnostics) {
	var diags diag.Diagnostics
	var devices []DeviceEntry

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		diags.AddError(
			"Error Opening CSV File",
			fmt.Sprintf("Could not open file at %s: %s", filePath, err),
		)
		return nil, diags
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Configure the reader to handle Windows Autopilot CSV format
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = false // Autopilot CSVs should not have quotes

	// Read the header row
	headers, err := reader.Read()
	if err != nil {
		diags.AddError(
			"Error Reading CSV Headers",
			fmt.Sprintf("Could not read headers from CSV file: %s", err),
		)
		return nil, diags
	}

	// Validate headers
	headerMap, headerDiags := validateHeaders(headers)
	diags.Append(headerDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Read each row
	lineNumber := 1 // Start at 1 for the header row
	for {
		lineNumber++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			diags.AddError(
				"Error Reading CSV Row",
				fmt.Sprintf("Error reading row %d: %s", lineNumber, err),
			)
			continue
		}

		// Process the row
		device, rowDiags := processRow(record, headerMap, lineNumber)
		diags.Append(rowDiags...)
		if !rowDiags.HasError() {
			devices = append(devices, device)
		}
	}

	// Log the number of devices read
	tflog.Info(ctx, fmt.Sprintf("Successfully read %d devices from CSV file", len(devices)))

	return devices, diags
}

// validateHeaders validates the CSV headers and returns a map of header indices
func validateHeaders(headers []string) (map[string]int, diag.Diagnostics) {
	var diags diag.Diagnostics
	headerMap := make(map[string]int)

	// Required headers
	requiredHeaders := []string{
		HeaderSerialNumber,
		HeaderHardwareHash,
	}

	// Optional headers
	optionalHeaders := []string{
		HeaderWindowsProduct,
		HeaderGroupTag,
		HeaderAssignedUser,
	}

	// Check for required headers
	for _, header := range requiredHeaders {
		found := false
		for i, h := range headers {
			if h == header {
				headerMap[header] = i
				found = true
				break
			}
		}

		if !found {
			diags.AddError(
				"Missing Required CSV Header",
				fmt.Sprintf("The CSV file is missing the required header: %s", header),
			)
		}
	}

	// Map optional headers if present
	for _, header := range optionalHeaders {
		for i, h := range headers {
			if h == header {
				headerMap[header] = i
				break
			}
		}
	}

	return headerMap, diags
}

// processRow processes a single row from the CSV file
func processRow(record []string, headerMap map[string]int, lineNumber int) (DeviceEntry, diag.Diagnostics) {
	var diags diag.Diagnostics
	var device DeviceEntry

	// Helper function to get a value from the record using the header map
	getValue := func(header string) string {
		if idx, ok := headerMap[header]; ok && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	// Get values from the record
	serialNumber := getValue(HeaderSerialNumber)
	windowsProductID := getValue(HeaderWindowsProduct)
	hardwareHash := getValue(HeaderHardwareHash)
	groupTag := getValue(HeaderGroupTag)
	assignedUser := getValue(HeaderAssignedUser)

	// Validate required fields
	if serialNumber == "" {
		diags.AddError(
			"Missing Required Value",
			fmt.Sprintf("Row %d is missing a value for %s", lineNumber, HeaderSerialNumber),
		)
	}

	if hardwareHash == "" {
		diags.AddError(
			"Missing Required Value",
			fmt.Sprintf("Row %d is missing a value for %s", lineNumber, HeaderHardwareHash),
		)
	}

	// Create the device entry
	device = DeviceEntry{
		SerialNumber:     types.StringValue(serialNumber),
		WindowsProductID: types.StringValue(windowsProductID),
		HardwareHash:     types.StringValue(hardwareHash),
		GroupTag:         types.StringValue(groupTag),
		AssignedUser:     types.StringValue(assignedUser),
	}

	return device, diags
}
