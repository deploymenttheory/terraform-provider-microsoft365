package graphBetaMacOSPKGApp

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

// The tiny XAR fixture makes unit tests independent of downloads. Live acceptance
// and curl probes use the checksum-verified official Nudge application package.
func writeUnitPKG(t *testing.T) string {
	t.Helper()
	packageInfo := []byte(`<pkg-info identifier="com.example.unit" version="1.2.3" install-location="/Applications"/>`)
	toc := []byte(fmt.Sprintf(`<xar><toc><file><name>PackageInfo</name><data><offset>0</offset><length>%d</length><size>%d</size><encoding style="application/octet-stream"/></data></file></toc></xar>`, len(packageInfo), len(packageInfo)))
	var compressed bytes.Buffer
	writer := zlib.NewWriter(&compressed)
	_, err := writer.Write(toc)
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	var archive bytes.Buffer
	header := struct {
		Magic                            uint32
		HeaderSize, Version              uint16
		CompressedSize, UncompressedSize int64
		HashType                         uint32
	}{0x78617221, 28, 1, int64(compressed.Len()), int64(len(toc)), 0}
	require.NoError(t, binary.Write(&archive, binary.BigEndian, header))
	archive.Write(compressed.Bytes())
	archive.Write(packageInfo)
	path := filepath.Join(t.TempDir(), "Unit.pkg")
	require.NoError(t, os.WriteFile(path, archive.Bytes(), 0600))
	return path
}

func constructorModel(t *testing.T) *MacOSPKGAppResourceModel {
	t.Helper()
	r := NewMacOSPKGAppResource().(*MacOSPKGAppResource)
	state := metadataUpdateState(t, r, "", "", true)
	var model MacOSPKGAppResourceModel
	require.False(t, state.Get(context.Background(), &model).HasError())
	model.MacOSPkgApp.PreInstallScript = &MacOSAppScriptResourceModel{ScriptContent: types.StringValue("ZXhpdCAwCg==")}
	model.MacOSPkgApp.PostInstallScript = &MacOSAppScriptResourceModel{ScriptContent: types.StringValue("ZXhpdCAwCg==")}
	return &model
}

func TestUnitResourceMacOSPkgApp_03_CreateInstallerMetadata(t *testing.T) {
	path := writeUnitPKG(t)
	body, err := constructResource(context.Background(), constructorModel(t), path, false)
	require.NoError(t, err)
	pkg := body.(graphmodels.MacOSPkgAppable)
	require.Equal(t, "Unit.pkg", *pkg.GetFileName())
	require.Equal(t, "com.example.unit", *pkg.GetPrimaryBundleId())
	require.Equal(t, "1.2.3", *pkg.GetPrimaryBundleVersion())
	require.NotEmpty(t, pkg.GetIncludedApps())
	require.Equal(t, "com.example.unit", *pkg.GetIncludedApps()[0].GetBundleId())
	require.Equal(t, "1.2.3", *pkg.GetIncludedApps()[0].GetBundleVersion())
}

func TestUnitResourceMacOSPkgApp_04_CreateStillRequiresInstaller(t *testing.T) {
	for _, path := range []string{"", filepath.Join(t.TempDir(), "missing.pkg")} {
		_, err := constructResource(context.Background(), constructorModel(t), path, false)
		require.Error(t, err)
	}
}

func TestUnitResourceMacOSPkgApp_05_UpdateDetectionMetadata(t *testing.T) {
	model := constructorModel(t)
	appType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"bundle_id": types.StringType, "bundle_version": types.StringType,
	}}
	model.MacOSPkgApp.IncludedApps = types.SetValueMust(appType, []attr.Value{
		types.ObjectValueMust(appType.AttrTypes, map[string]attr.Value{
			"bundle_id":      types.StringValue("com.example.detection"),
			"bundle_version": types.StringValue("2.0"),
		}),
	})
	body, err := constructResource(context.Background(), model, filepath.Join(t.TempDir(), "missing.pkg"), true)
	require.NoError(t, err)
	pkg := body.(graphmodels.MacOSPkgAppable)
	require.Nil(t, pkg.GetFileName())
	require.Nil(t, pkg.GetPrimaryBundleId())
	require.Nil(t, pkg.GetPrimaryBundleVersion())
	require.Nil(t, pkg.GetCommittedContentVersion())
	require.Len(t, pkg.GetIncludedApps(), 1)
	require.Equal(t, "com.example.detection", *pkg.GetIncludedApps()[0].GetBundleId())
	require.Equal(t, "2.0", *pkg.GetIncludedApps()[0].GetBundleVersion())
	require.True(t, *pkg.GetIgnoreVersionDetection())
	require.True(t, *pkg.GetMinimumSupportedOperatingSystem().GetV120())
	require.Equal(t, "ZXhpdCAwCg==", *pkg.GetPreInstallScript().GetScriptContent())
	require.Equal(t, "ZXhpdCAwCg==", *pkg.GetPostInstallScript().GetScriptContent())
}
