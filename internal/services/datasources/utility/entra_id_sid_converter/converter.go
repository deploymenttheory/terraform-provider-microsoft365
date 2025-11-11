package entra_id_sid_converter

import (
	"encoding/binary"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
)

var (
	sidPattern  = regexp.MustCompile(constants.EntraIdSidRegex)
	guidPattern = regexp.MustCompile(constants.GuidRegex)
)

func convertSidToObjectId(sid string) (string, error) {
	if !sidPattern.MatchString(sid) {
		return "", fmt.Errorf("invalid SID format: %s", sid)
	}

	parts := strings.Split(sid, "-")
	if len(parts) != 8 {
		return "", fmt.Errorf("invalid SID format: expected 8 parts, got %d", len(parts))
	}

	rids := make([]uint32, 4)
	for i := range 4 {
		rid, err := strconv.ParseUint(parts[4+i], 10, 32)
		if err != nil {
			return "", fmt.Errorf("invalid RID component at position %d: %s", i, parts[4+i])
		}
		rids[i] = uint32(rid)
	}

	var bytes [16]byte
	for i, rid := range rids {
		binary.LittleEndian.PutUint32(bytes[i*4:], rid)
	}

	objectId := fmt.Sprintf(
		"%08x-%04x-%04x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		binary.LittleEndian.Uint32(bytes[0:4]),
		binary.LittleEndian.Uint16(bytes[4:6]),
		binary.LittleEndian.Uint16(bytes[6:8]),
		bytes[8], bytes[9],
		bytes[10], bytes[11], bytes[12], bytes[13], bytes[14], bytes[15],
	)

	return objectId, nil
}

func convertObjectIdToSid(objectId string) (string, error) {
	if !guidPattern.MatchString(objectId) {
		return "", fmt.Errorf("invalid Object ID format: %s", objectId)
	}

	parts := strings.Split(objectId, "-")
	if len(parts) != 5 {
		return "", fmt.Errorf("invalid GUID format: expected 5 parts separated by dashes")
	}

	var bytes [16]byte

	// Parse Data1 (first 8 hex chars) as little-endian uint32
	var data1 uint32
	_, err := fmt.Sscanf(parts[0], "%08x", &data1)
	if err != nil {
		return "", fmt.Errorf("invalid Data1 component: %s", parts[0])
	}
	binary.LittleEndian.PutUint32(bytes[0:4], data1)

	// Parse Data2 (next 4 hex chars) as little-endian uint16
	var data2 uint16
	_, err = fmt.Sscanf(parts[1], "%04x", &data2)
	if err != nil {
		return "", fmt.Errorf("invalid Data2 component: %s", parts[1])
	}
	binary.LittleEndian.PutUint16(bytes[4:6], data2)

	// Parse Data3 (next 4 hex chars) as little-endian uint16
	var data3 uint16
	_, err = fmt.Sscanf(parts[2], "%04x", &data3)
	if err != nil {
		return "", fmt.Errorf("invalid Data3 component: %s", parts[2])
	}
	binary.LittleEndian.PutUint16(bytes[6:8], data3)

	// Parse Data4[0-1] (next 4 hex chars) as big-endian
	_, err = fmt.Sscanf(parts[3], "%02x%02x", &bytes[8], &bytes[9])
	if err != nil {
		return "", fmt.Errorf("invalid Data4 component (part 1): %s", parts[3])
	}

	// Parse Data4[2-7] (last 12 hex chars) as big-endian
	_, err = fmt.Sscanf(parts[4], "%02x%02x%02x%02x%02x%02x",
		&bytes[10], &bytes[11], &bytes[12], &bytes[13], &bytes[14], &bytes[15])
	if err != nil {
		return "", fmt.Errorf("invalid Data4 component (part 2): %s", parts[4])
	}

	rids := make([]uint32, 4)
	for i := range 4 {
		rids[i] = binary.LittleEndian.Uint32(bytes[i*4 : (i+1)*4])
	}

	sid := fmt.Sprintf("S-1-12-1-%d-%d-%d-%d", rids[0], rids[1], rids[2], rids[3])

	return sid, nil
}
