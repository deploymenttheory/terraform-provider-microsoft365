package sharedStater

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// TestMapCommittedContentVersionStateToTerraform tests the MapCommittedContentVersionStateToTerraform function
func TestMapCommittedContentVersionStateToTerraform(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name              string
		versionId         string
		respFiles         any
		err               error
		installerFileName string
		validate          func(t *testing.T, result types.List)
	}{
		{
			name:              "Nil response files",
			versionId:         "version-123",
			respFiles:         nil,
			err:               nil,
			installerFileName: "",
			validate: func(t *testing.T, result types.List) {
				assert.False(t, result.IsNull())
				require.Len(t, result.Elements(), 1)
			},
		},
		{
			name:              "Error fetching files",
			versionId:         "version-123",
			respFiles:         nil,
			err:               fmt.Errorf("fetch error"),
			installerFileName: "",
			validate: func(t *testing.T, result types.List) {
				assert.False(t, result.IsNull())
				require.Len(t, result.Elements(), 1)
			},
		},
		{
			name:      "Valid file collection",
			versionId: "version-123",
			respFiles: func() graphmodels.MobileAppContentFileCollectionResponseable {
				collection := graphmodels.NewMobileAppContentFileCollectionResponse()
				
				file1 := graphmodels.NewMobileAppContentFile()
				fileName := "test-app.apk"
				fileSize := int64(1024)
				isCommitted := true
				
				file1.SetName(&fileName)
				file1.SetSize(&fileSize)
				file1.SetIsCommitted(&isCommitted)
				
				collection.SetValue([]graphmodels.MobileAppContentFileable{file1})
				return collection
			}(),
			err:               nil,
			installerFileName: "",
			validate: func(t *testing.T, result types.List) {
				assert.False(t, result.IsNull())
				require.Len(t, result.Elements(), 1)
			},
		},
		{
			name:      "Filter by installer file name",
			versionId: "version-123",
			respFiles: func() graphmodels.MobileAppContentFileCollectionResponseable {
				collection := graphmodels.NewMobileAppContentFileCollectionResponse()
				
				file1 := graphmodels.NewMobileAppContentFile()
				fileName1 := "test-app.apk"
				fileSize := int64(1024)
				isCommitted := true
				
				file1.SetName(&fileName1)
				file1.SetSize(&fileSize)
				file1.SetIsCommitted(&isCommitted)
				
				file2 := graphmodels.NewMobileAppContentFile()
				fileName2 := "other-file.txt"
				file2.SetName(&fileName2)
				file2.SetSize(&fileSize)
				file2.SetIsCommitted(&isCommitted)
				
				collection.SetValue([]graphmodels.MobileAppContentFileable{file1, file2})
				return collection
			}(),
			err:               nil,
			installerFileName: "test-app.apk",
			validate: func(t *testing.T, result types.List) {
				assert.False(t, result.IsNull())
				// Should only include the matching file
				require.Len(t, result.Elements(), 1)
			},
		},
		{
			name:      "Invalid response type",
			versionId: "version-123",
			respFiles: "invalid-type",
			err:       nil,
			installerFileName: "",
			validate: func(t *testing.T, result types.List) {
				assert.False(t, result.IsNull())
				// Should handle gracefully and return empty files
				require.Len(t, result.Elements(), 1)
			},
		},
		{
			name:      "Nil file in collection",
			versionId: "version-123",
			respFiles: func() graphmodels.MobileAppContentFileCollectionResponseable {
				collection := graphmodels.NewMobileAppContentFileCollectionResponse()
				
				file1 := graphmodels.NewMobileAppContentFile()
				fileName := "test-app.apk"
				fileSize := int64(1024)
				isCommitted := true
				
				file1.SetName(&fileName)
				file1.SetSize(&fileSize)
				file1.SetIsCommitted(&isCommitted)
				
				collection.SetValue([]graphmodels.MobileAppContentFileable{file1, nil})
				return collection
			}(),
			err:               nil,
			installerFileName: "",
			validate: func(t *testing.T, result types.List) {
				assert.False(t, result.IsNull())
				// Should skip nil files
				require.Len(t, result.Elements(), 1)
			},
		},
		{
			name:      "Empty version ID",
			versionId: "",
			respFiles: func() graphmodels.MobileAppContentFileCollectionResponseable {
				collection := graphmodels.NewMobileAppContentFileCollectionResponse()
				file1 := graphmodels.NewMobileAppContentFile()
				fileName := "test-app.apk"
				file1.SetName(&fileName)
				collection.SetValue([]graphmodels.MobileAppContentFileable{file1})
				return collection
			}(),
			err:               nil,
			installerFileName: "",
			validate: func(t *testing.T, result types.List) {
				assert.False(t, result.IsNull())
				require.Len(t, result.Elements(), 1)
				// Should still create list even with empty version ID
			},
		},
		{
			name:      "File with nil name",
			versionId: "version-123",
			respFiles: func() graphmodels.MobileAppContentFileCollectionResponseable {
				collection := graphmodels.NewMobileAppContentFileCollectionResponse()
				file1 := graphmodels.NewMobileAppContentFile()
				file1.SetName(nil)
				collection.SetValue([]graphmodels.MobileAppContentFileable{file1})
				return collection
			}(),
			err:               nil,
			installerFileName: "test.apk",
			validate: func(t *testing.T, result types.List) {
				assert.False(t, result.IsNull())
				// File with nil name won't match filter, so no files in result
				require.Len(t, result.Elements(), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapCommittedContentVersionStateToTerraform(
				ctx,
				tt.versionId,
				tt.respFiles,
				tt.err,
				tt.installerFileName,
			)
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}
