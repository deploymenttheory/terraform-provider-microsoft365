package graphBetaDirectoryRoleAssignment_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/directory"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

func directoryRoleAcceptancePrincipal(t *testing.T) string {
	t.Helper()
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	name := "acc-test-directory-role-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	application := graphmodels.NewApplication()
	application.SetDisplayName(&name)
	application.SetSignInAudience(new("AzureADMyOrg"))
	created, err := client.Applications().Post(ctx, application, nil)
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, created.GetId())
	require.NotNil(t, created.GetAppId())
	appID := *created.GetId()
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cleanupCancel()
		opts := directory.DefaultDeleteOptions()
		opts.ResourceID = appID
		opts.ResourceType = directory.ResourceTypeApplication
		if err := directory.ExecuteDeleteWithVerification(
			cleanupCtx,
			client,
			func(ctx context.Context) error { return client.Applications().ByApplicationId(appID).Delete(ctx, nil) },
			true,
			opts,
		); err != nil {
			t.Errorf("clean up owned role fixture application: %v", err)
		}
	})
	body := graphmodels.NewServicePrincipal()
	body.SetAppId(created.GetAppId())
	body.SetAccountEnabled(new(false))
	var principal graphmodels.ServicePrincipalable
	err = crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		var postErr error
		principal, postErr = client.ServicePrincipals().Post(ctx, body, nil)
		if postErr == nil {
			return true, nil
		}
		info := grapherrors.GraphError(ctx, postErr)
		if info.StatusCode == 404 || (info.StatusCode == 400 && (strings.Contains(strings.ToLower(info.ErrorMessage), "not found") || strings.Contains(info.ErrorMessage, "does not reference a valid application object"))) {
			return false, fmt.Errorf("waiting for fixture application replication: %w", postErr)
		}
		return false, &crud.FatalPollError{Err: postErr}
	})
	require.NoError(t, err)
	require.NotNil(t, principal)
	require.NotNil(t, principal.GetId())
	spID := *principal.GetId()
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cleanupCancel()
		opts := directory.DefaultDeleteOptions()
		opts.ResourceID = spID
		opts.ResourceType = directory.ResourceTypeServicePrincipal
		if err := directory.ExecuteDeleteWithVerification(cleanupCtx, client, func(ctx context.Context) error {
			return client.ServicePrincipals().ByServicePrincipalId(spID).Delete(ctx, nil)
		}, true, opts); err != nil {
			t.Errorf("clean up owned role fixture service principal: %v", err)
		}
	})
	t.Logf("Owned disabled directory-role fixture: application=%s service_principal=%s", appID, spID)
	return spID
}
