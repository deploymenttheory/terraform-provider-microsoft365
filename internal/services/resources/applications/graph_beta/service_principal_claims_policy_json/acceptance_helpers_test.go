package graphBetaServicePrincipalClaimsPolicyJSON_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/jarcoal/httpmock"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/applicationtemplates"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/directory"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

var (
	errFixtureOwnership = errors.New("fixture ownership prefix mismatch")
	errFixtureResponse  = errors.New("instantiate response omitted fixture object IDs")
)

func instantiateClaimsFixture(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, name string) (string, string, error) {
	body := applicationtemplates.NewItemInstantiatePostRequestBody()
	body.SetDisplayName(&name)
	request, err := client.ApplicationTemplates().ByApplicationTemplateId("8adf8e6e-67b2-4cf2-a259-e3dc5476c621").Instantiate().ToPostRequestInformation(ctx, body, nil)
	if err != nil {
		return "", "", fmt.Errorf("build fixture instantiate request: %w", err)
	}
	response, err := client.GetAdapter().SendPrimitive(ctx, request, "[]byte", abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue})
	if err != nil {
		return "", "", fmt.Errorf("instantiate fixture: %w", err)
	}
	content, ok := response.([]byte)
	if !ok {
		return "", "", errFixtureResponse
	}
	var pair struct {
		Application struct {
			ID       string `json:"id"`
			ObjectID string `json:"objectId"`
		} `json:"application"`
		ServicePrincipal struct {
			ID       string `json:"id"`
			ObjectID string `json:"objectId"`
		} `json:"servicePrincipal"`
	}
	if err := json.Unmarshal(content, &pair); err != nil {
		return "", "", fmt.Errorf("decode fixture identities: %w", err)
	}
	appID, spID := pair.Application.ID, pair.ServicePrincipal.ID
	if appID == "" {
		appID = pair.Application.ObjectID
	}
	if spID == "" {
		spID = pair.ServicePrincipal.ObjectID
	}
	if appID == "" || spID == "" {
		return appID, spID, errFixtureResponse
	}
	return appID, spID, nil
}

func cleanupClaimsFixture(t *testing.T, client *msgraphbetasdk.GraphServiceClient, id, name string, kind directory.ResourceType) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		softDelete := func(ctx context.Context) error {
			if kind == directory.ResourceTypeApplication {
				return client.Applications().ByApplicationId(id).Delete(ctx, nil)
			}
			return client.ServicePrincipals().ByServicePrincipalId(id).Delete(ctx, nil)
		}
		opts := directory.DefaultDeleteOptions()
		opts.ResourceID = id
		opts.ResourceName = name
		opts.ResourceType = kind
		if err := directory.ExecuteDeleteWithVerification(ctx, client, softDelete, true, opts); err != nil {
			t.Errorf("cleanup owned %s %s: %v", kind, id, err)
			return
		}
		// Deleted items can appear after an immediate 404. Recheck both collections
		// across multiple observations and purge only the exact created fixture ID.
		attempts, absent := 0, 0
		err := crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
			attempts++
			active, activeErr := client.DirectoryObjects().ByDirectoryObjectId(id).Get(ctx, nil)
			if activeErr != nil && grapherrors.GraphError(ctx, activeErr).StatusCode != 404 {
				return false, &crud.FatalPollError{Err: activeErr}
			}
			deleted, deletedErr := client.Directory().DeletedItems().ByDirectoryObjectId(id).Get(ctx, nil)
			if deletedErr != nil && grapherrors.GraphError(ctx, deletedErr).StatusCode != 404 {
				return false, &crud.FatalPollError{Err: deletedErr}
			}
			if activeErr == nil || deletedErr == nil {
				absent = 0
				if !strings.HasPrefix(name, "acc-test-claims-policy-") {
					return false, &crud.FatalPollError{Err: errFixtureOwnership}
				}
				if active != nil && active.GetId() != nil && *active.GetId() == id {
					if err := softDelete(ctx); err != nil && grapherrors.GraphError(ctx, err).StatusCode != 404 {
						return false, &crud.FatalPollError{Err: err}
					}
				}
				if deleted != nil && deleted.GetId() != nil && *deleted.GetId() == id {
					if err := client.Directory().DeletedItems().ByDirectoryObjectId(id).Delete(ctx, nil); err != nil && grapherrors.GraphError(ctx, err).StatusCode != 404 {
						return false, &crud.FatalPollError{Err: err}
					}
				}
				return false, nil
			}
			absent++
			return attempts >= 10 && absent >= 3, nil
		})
		if err != nil {
			t.Errorf("confirm complete owned fixture cleanup %s: %v", id, err)
		}
	})
}

func TestUnitClaimsPolicyFixtureLegacyObjectIDs(t *testing.T) {
	transport := httpmock.NewMockTransport()
	transport.RegisterResponder("POST", "https://graph.microsoft.com/beta/applicationTemplates/8adf8e6e-67b2-4cf2-a259-e3dc5476c621/instantiate", httpmock.NewStringResponder(201,
		`{"application":{"objectId":"11111111-1111-4111-8111-111111111111"},"servicePrincipal":{"objectId":"22222222-2222-4222-8222-222222222222"}}`))
	client := providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient
	appID, spID, err := instantiateClaimsFixture(context.Background(), client, "acc-test-claims-policy-unit")
	require.NoError(t, err)
	require.Equal(t, "11111111-1111-4111-8111-111111111111", appID)
	require.Equal(t, "22222222-2222-4222-8222-222222222222", spID)
}

func claimsPolicyAcceptancePrincipal(t *testing.T) string {
	t.Helper()
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	name := "acc-test-claims-policy-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	t.Logf("Creating owned fixture %s", name)
	appID, spID, err := instantiateClaimsFixture(ctx, client, name)
	if appID != "" {
		cleanupClaimsFixture(t, client, appID, name, directory.ResourceTypeApplication)
	}
	if spID != "" {
		cleanupClaimsFixture(t, client, spID, name, directory.ResourceTypeServicePrincipal)
	}
	require.NoError(t, err)
	require.NotEmpty(t, appID)
	require.NotEmpty(t, spID)
	t.Logf("Owned fixture application=%s service_principal=%s", appID, spID)
	disabled := graphmodels.NewServicePrincipal()
	disabled.SetAccountEnabled(new(false))
	disabled.SetPreferredSingleSignOnMode(new("saml"))
	disabled.SetAppRoleAssignmentRequired(new(true))
	err = crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		_, patchErr := client.ServicePrincipals().ByServicePrincipalId(spID).Patch(ctx, disabled, nil)
		if patchErr != nil {
			if grapherrors.GraphError(ctx, patchErr).StatusCode == 404 {
				return false, fmt.Errorf("waiting for fixture replication: %w", patchErr)
			}
			return false, &crud.FatalPollError{Err: patchErr}
		}
		remote, getErr := client.ServicePrincipals().ByServicePrincipalId(spID).Get(ctx, nil)
		if getErr != nil {
			if grapherrors.GraphError(ctx, getErr).StatusCode == 404 {
				return false, fmt.Errorf("waiting for disabled fixture: %w", getErr)
			}
			return false, &crud.FatalPollError{Err: getErr}
		}
		return remote != nil && remote.GetAccountEnabled() != nil && !*remote.GetAccountEnabled() && remote.GetPreferredSingleSignOnMode() != nil &&
			*remote.GetPreferredSingleSignOnMode() == "saml" &&
			remote.GetAppRoleAssignmentRequired() != nil &&
			*remote.GetAppRoleAssignmentRequired(), nil
	})
	require.NoError(t, err)
	return spID
}
