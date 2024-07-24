package graphBetaAssignmentFilter

// import (
// 	"fmt"

// 	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
// )

// // StringToDevicePlatformType converts a string to DevicePlatformType.
// func StringToDevicePlatformType(platformStr string, supportedPlatformTypes map[string]models.DevicePlatformType) (*models.DevicePlatformType, error) {
// 	platform, exists := supportedPlatformTypes[platformStr]
// 	if !exists {
// 		return nil, fmt.Errorf("unsupported platform type: %s", platformStr)
// 	}
// 	return &platform, nil
// }

// // DevicePlatformTypeToString converts a DevicePlatformType to its string representation.
// func DevicePlatformTypeToString(platform *models.DevicePlatformType) (string, error) {
// 	if platform == nil {
// 		return "", fmt.Errorf("platform is nil")
// 	}
// 	return platform.String(), nil
// }

// // supportedPlatformTypes is a map of string representations to their corresponding platform types.
// var supportedPlatformTypes = map[string]models.DevicePlatformType{
// 	"android":                            models.ANDROID_DEVICEPLATFORMTYPE,
// 	"androidForWork":                     models.ANDROIDFORWORK_DEVICEPLATFORMTYPE,
// 	"iOS":                                models.IOS_DEVICEPLATFORMTYPE,
// 	"macOS":                              models.MACOS_DEVICEPLATFORMTYPE,
// 	"windowsPhone81":                     models.WINDOWSPHONE81_DEVICEPLATFORMTYPE,
// 	"windows81AndLater":                  models.WINDOWS81ANDLATER_DEVICEPLATFORMTYPE,
// 	"windows10AndLater":                  models.WINDOWS10ANDLATER_DEVICEPLATFORMTYPE,
// 	"androidWorkProfile":                 models.ANDROIDWORKPROFILE_DEVICEPLATFORMTYPE,
// 	"unknown":                            models.UNKNOWN_DEVICEPLATFORMTYPE,
// 	"androidAOSP":                        models.ANDROIDAOSP_DEVICEPLATFORMTYPE,
// 	"androidMobileApplicationManagement": models.ANDROIDMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE,
// 	"iOSMobileApplicationManagement":     models.IOSMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE,
// 	"unknownFutureValue":                 models.UNKNOWNFUTUREVALUE_DEVICEPLATFORMTYPE,
// 	"windowsMobileApplicationManagement": models.WINDOWSMOBILEAPPLICATIONMANAGEMENT_DEVICEPLATFORMTYPE,
// }
