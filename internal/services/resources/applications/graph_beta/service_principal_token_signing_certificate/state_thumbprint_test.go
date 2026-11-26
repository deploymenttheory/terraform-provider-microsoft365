package graphBetaServicePrincipalTokenSigningCertificate

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitTokenSigningCertificateThumbprintRoundTrip(t *testing.T) {
	const expected = "aabbccddeeff00112233445566778899aabbccdd"
	keyID := uuid.MustParse("11111111-1111-4111-8111-111111111111")
	for _, representation := range []string{"raw_hash", "uppercase_hex", "lowercase_hex"} {
		t.Run(representation, func(t *testing.T) {
			ctx := context.Background()
			created := ServicePrincipalTokenSigningCertificateResourceModel{ServicePrincipalID: types.StringValue("22222222-2222-4222-8222-222222222222")}
			certificate := graphmodels.NewSelfSignedCertificate()
			upper := strings.ToUpper(expected)
			certificate.SetThumbprint(&upper)
			certificate.SetKeyId(&keyID)
			MapSelfSignedCertificateToTerraform(ctx, &created, certificate)
			identifier := []byte(expected)
			if representation == "raw_hash" {
				var err error
				identifier, err = hex.DecodeString(expected)
				require.NoError(t, err)
			} else if representation == "uppercase_hex" {
				identifier = []byte(upper)
			}
			credential := graphmodels.NewKeyCredential()
			credential.SetKeyId(&keyID)
			credential.SetCustomKeyIdentifier(identifier)
			principal := graphmodels.NewServicePrincipal()
			principal.SetKeyCredentials([]graphmodels.KeyCredentialable{credential})
			imported := ServicePrincipalTokenSigningCertificateResourceModel{ServicePrincipalID: created.ServicePrincipalID, KeyId: created.KeyId, Thumbprint: types.StringNull()}
			require.True(t, MapRemoteResourceStateToTerraform(ctx, &imported, principal))
			require.Equal(t, created.Thumbprint, imported.Thumbprint, "create and import must produce identical state")
			require.Equal(t, expected, imported.Thumbprint.ValueString())
		})
	}
	t.Run("missing_thumbprint_stays_null", func(t *testing.T) {
		model := ServicePrincipalTokenSigningCertificateResourceModel{}
		MapSelfSignedCertificateToTerraform(context.Background(), &model, graphmodels.NewSelfSignedCertificate())
		require.True(t, model.Thumbprint.IsNull())
	})
}
