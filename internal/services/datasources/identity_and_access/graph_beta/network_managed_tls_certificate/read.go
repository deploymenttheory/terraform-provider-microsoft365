package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (d *NetworkManagedTLSCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object NetworkManagedTLSCertificateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	certificate, err := d.getManagedTLSCertificate(ctx, object.CertificateAuthorityID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	mapRemoteStateToDataSource(&object, certificate)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}
