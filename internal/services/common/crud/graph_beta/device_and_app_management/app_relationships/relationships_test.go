package apprelationships

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/microsoft/kiota-abstractions-go/authentication"
	kiotahttp "github.com/microsoft/kiota-http-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestListAndFindRelationshipOnLaterPage(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/deviceAppManagement/mobileApps/source/relationships":
			_, _ = fmt.Fprintf(
				w,
				`{"value":[{"@odata.type":"#microsoft.graph.mobileAppDependency","id":"incoming","targetId":"target","targetType":"parent","dependencyType":"detect"}],"@odata.nextLink":%q}`,
				server.URL+"/next",
			)
		case "/next":
			_, _ = fmt.Fprint(
				w,
				`{"value":[{"@odata.type":"#microsoft.graph.mobileAppDependency","id":"saved-id","targetId":"target","targetType":"child","dependencyType":"autoInstall"}]}`,
			)
		default:
			t.Errorf("unexpected request: %s", r.URL)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	t.Cleanup(server.Close)
	adapter, err := kiotahttp.NewNetHttpRequestAdapter(
		&authentication.AnonymousAuthenticationProvider{},
	)
	require.NoError(t, err)
	adapter.SetBaseUrl(server.URL)

	relationships, err := List(
		context.Background(),
		msgraphbetasdk.NewGraphServiceClient(adapter),
		"source",
	)
	require.NoError(t, err)
	require.Len(t, relationships, 2)
	found, err := Find(relationships, "", "target")
	require.NoError(t, err)
	require.Equal(t, "saved-id", *found.GetId())
	require.IsType(t, graphmodels.NewMobileAppDependency(), found)
	_, err = Find(relationships, "deleted-id", "target")
	require.ErrorIs(
		t,
		err,
		ErrNotFound,
		"a replacement relationship must not be adopted during refresh",
	)
}

func TestListDistinguishesEmptyIncompleteAndFailedResponses(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		ok     bool
	}{
		{"empty", 200, `{"value":[]}`, true},
		{"missing value", 200, `{}`, false},
		{"bad request", 400, `{"error":{"code":"BadRequest","message":"Unsupported route"}}`, false},
		{"forbidden", 403, `{"error":{"code":"Forbidden","message":"Access denied"}}`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tc.status)
					_, _ = fmt.Fprint(w, tc.body)
				}),
			)
			t.Cleanup(server.Close)
			adapter, err := kiotahttp.NewNetHttpRequestAdapter(
				&authentication.AnonymousAuthenticationProvider{},
			)
			require.NoError(t, err)
			adapter.SetBaseUrl(server.URL)
			relationships, err := List(
				context.Background(),
				msgraphbetasdk.NewGraphServiceClient(adapter),
				"source",
			)
			if tc.ok {
				require.NoError(t, err)
				require.Empty(t, relationships)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestSourceID(t *testing.T) {
	const source = "11111111-1111-1111-1111-111111111111"
	require.Equal(t, source, SourceID(source+"_22222222-2222-2222-2222-222222222222"))
	for _, id := range []string{"opaque-id", "invalid_target", source + "_invalid", source} {
		require.Empty(t, SourceID(id))
	}
}
