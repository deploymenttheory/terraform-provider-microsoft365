package graphBetaWindowsTrustedRootCertificate

import (
	"context"
	"fmt"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func constructResource(
	ctx context.Context,
	data *WindowsTrustedRootCertificateResourceModel,
) (graphmodels.Windows81TrustedRootCertificateable, error) {
	der, err := decodeCertificate(data.TrustedRootCertificate.ValueString())
	if err != nil {
		return nil, err
	}
	body := graphmodels.NewWindows81TrustedRootCertificate()
	body.SetTrustedRootCertificate(der)
	convert.FrameworkToGraphString(data.DisplayName, body.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, body.SetDescription)
	convert.FrameworkToGraphString(data.CertFileName, body.SetCertFileName)
	if err := convert.FrameworkToGraphEnum(
		data.DestinationStore,
		graphmodels.ParseCertificateDestinationStore,
		body.SetDestinationStore,
	); err != nil {
		return nil, fmt.Errorf("set destination store: %w", err)
	}
	// Leave omitted scope tags unset so Intune supplies the default tag.
	// The SDK backing store serializes an explicit nil setter as JSON null.
	if !data.RoleScopeTagIds.IsNull() && !data.RoleScopeTagIds.IsUnknown() {
		if err := convert.FrameworkToGraphStringSet(
			ctx,
			data.RoleScopeTagIds,
			body.SetRoleScopeTagIds,
		); err != nil {
			return nil, fmt.Errorf("set scope tags: %w", err)
		}
	}
	return body, nil
}
