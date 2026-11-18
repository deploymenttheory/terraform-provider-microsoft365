package graphBetaOAuth2PermissionGrant_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/directory"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

type grantFixture struct{ ClientID, ResourceID, PrincipalID string }

func cleanupGrantFixtureObject(t *testing.T, client *msgraphbetasdk.GraphServiceClient, id string, kind directory.ResourceType) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		opts := directory.DefaultDeleteOptions()
		opts.ResourceID = id
		opts.ResourceType = kind
		err := directory.ExecuteDeleteWithVerification(ctx, client, func(ctx context.Context) error {
			switch kind {
			case directory.ResourceTypeApplication:
				return client.Applications().ByApplicationId(id).Delete(ctx, nil)
			case directory.ResourceTypeServicePrincipal:
				return client.ServicePrincipals().ByServicePrincipalId(id).Delete(ctx, nil)
			default:
				return client.Users().ByUserId(id).Delete(ctx, nil)
			}
		}, true, opts)
		if err != nil {
			t.Errorf("clean up owned %s %s: %v", kind, id, err)
		}
	})
}

func grantAcceptanceFixture(t *testing.T, principal bool) *grantFixture {
	t.Helper()
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	suffix := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	fixture := &grantFixture{}
	for _, purpose := range []string{"client", "api"} {
		name := "acc-test-oauth2-grant-" + suffix + "-" + purpose
		application := graphmodels.NewApplication()
		application.SetDisplayName(&name)
		application.SetSignInAudience(new("AzureADMyOrg"))
		if purpose == "api" {
			api := graphmodels.NewApiApplication()
			scopes := []graphmodels.PermissionScopeable{}
			for _, value := range []string{"Probe.Read", "Probe.Write"} {
				scope := graphmodels.NewPermissionScope()
				scope.SetId(new(uuid.New()))
				scope.SetValue(&value)
				scope.SetTypeEscaped(new("User"))
				scope.SetIsEnabled(new(true))
				scope.SetAdminConsentDisplayName(new("Access owned acceptance probe"))
				scope.SetAdminConsentDescription(new("Access this owned acceptance API only"))
				scope.SetUserConsentDisplayName(new("Access owned acceptance probe"))
				scope.SetUserConsentDescription(new("Access this owned acceptance API only"))
				scopes = append(scopes, scope)
			}
			api.SetOauth2PermissionScopes(scopes)
			application.SetApi(api)
		}
		created, err := client.Applications().Post(ctx, application, nil)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.NotNil(t, created.GetId())
		require.NotNil(t, created.GetAppId())
		cleanupGrantFixtureObject(t, client, *created.GetId(), directory.ResourceTypeApplication)
		sp := graphmodels.NewServicePrincipal()
		sp.SetAppId(created.GetAppId())
		sp.SetAccountEnabled(new(false))
		var servicePrincipal graphmodels.ServicePrincipalable
		err = crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
			var postErr error
			servicePrincipal, postErr = client.ServicePrincipals().Post(ctx, sp, nil)
			if postErr == nil {
				return true, nil
			}
			info := grapherrors.GraphError(ctx, postErr)
			if info.StatusCode == 400 || info.StatusCode == 404 {
				return false, postErr
			}
			return false, &crud.FatalPollError{Err: postErr}
		})
		require.NoError(t, err)
		require.NotNil(t, servicePrincipal)
		require.NotNil(t, servicePrincipal.GetId())
		id := *servicePrincipal.GetId()
		cleanupGrantFixtureObject(t, client, id, directory.ResourceTypeServicePrincipal)
		if purpose == "client" {
			fixture.ClientID = id
		} else {
			fixture.ResourceID = id
			err = crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
				remote, getErr := client.ServicePrincipals().ByServicePrincipalId(id).Get(ctx, nil)
				if getErr != nil {
					return false, getErr
				}
				return remote != nil && len(remote.GetPublishedPermissionScopes()) == 2, nil
			})
			require.NoError(t, err)
		}
	}
	if principal {
		domains, err := client.Domains().Get(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, domains)
		domain := ""
		for _, candidate := range domains.GetValue() {
			if candidate.GetIsDefault() != nil && *candidate.GetIsDefault() && candidate.GetId() != nil {
				domain = *candidate.GetId()
				break
			}
		}
		require.NotEmpty(t, domain, "no default tenant domain for the owned user fixture")
		name := "acc-test-oauth2-grant-user-" + suffix
		user := graphmodels.NewUser()
		user.SetDisplayName(&name)
		user.SetMailNickname(&name)
		user.SetUserPrincipalName(new(name + "@" + domain))
		user.SetAccountEnabled(new(false))
		password := graphmodels.NewPasswordProfile()
		password.SetPassword(new("Aa1!" + rand.Text()))
		password.SetForceChangePasswordNextSignIn(new(false))
		user.SetPasswordProfile(password)
		created, err := client.Users().Post(ctx, user, nil)
		require.NoError(t, err)
		require.NotNil(t, created)
		require.NotNil(t, created.GetId())
		fixture.PrincipalID = *created.GetId()
		cleanupGrantFixtureObject(t, client, fixture.PrincipalID, directory.ResourceTypeUser)
	}
	t.Logf("Owned inactive grant fixtures: client=%s resource=%s principal=%s", fixture.ClientID, fixture.ResourceID, fixture.PrincipalID)
	return fixture
}

func (f *grantFixture) config(t *testing.T, principal bool) string {
	t.Helper()
	file := "resource_all_principals.tf"
	if principal {
		file = "resource_principal.tf"
	}
	content, err := loadAcceptanceTerraform(file)
	require.NoError(t, err)
	return fmt.Sprintf(content, f.ClientID, f.ResourceID, f.PrincipalID)
}
