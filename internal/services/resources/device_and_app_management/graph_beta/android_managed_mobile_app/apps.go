package graphBetaDeviceAndAppManagementAndroidManagedMobileApp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// The androidManagedAppProtections/{id}/apps collection supports only a GET of the
// whole list; per-app GET, POST and PATCH return "No OData route exists". Writes go
// through managedAppPolicies/{id}/targetApps, which REPLACES the policy's entire
// app list. Every write therefore reads the list, changes one entry and posts the
// full result while holding a per-policy lock, so resources that share a policy
// cannot overwrite each other's apps.

var (
	errAppAlreadyTargeted = errors.New("the app is already targeted by this policy")
	errAppWaitTimeout     = errors.New(
		"timed out waiting for the policy app list to reflect the change",
	)
	policyLocks sync.Map

	// appPollInterval is how often the app list is re-read after a write.
	appPollInterval = 2 * time.Second
)

// lockPolicy serializes app-list writes for one policy until the returned func is called.
func lockPolicy(ctx context.Context, policyID string) (func(), error) {
	value, _ := policyLocks.LoadOrStore(strings.ToLower(policyID), make(chan struct{}, 1))
	lock := value.(chan struct{})
	select {
	case lock <- struct{}{}:
		return func() { <-lock }, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("wait for policy app lock: %w", ctx.Err())
	}
}

// listPolicyApps returns every app currently targeted by the policy.
func (r *AndroidManagedMobileAppResource) listPolicyApps(
	ctx context.Context,
	policyID string,
) ([]graphmodels.ManagedMobileAppable, error) {
	response, err := r.client.
		DeviceAppManagement().
		AndroidManagedAppProtections().
		ByAndroidManagedAppProtectionId(policyID).
		Apps().
		Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("list policy apps: %w", err)
	}
	return response.GetValue(), nil
}

// setPolicyApps replaces the policy's targeted apps with the given list.
func (r *AndroidManagedMobileAppResource) setPolicyApps(
	ctx context.Context,
	policyID string,
	apps []graphmodels.ManagedMobileAppable,
) error {
	body := deviceappmanagement.NewManagedAppPoliciesItemTargetAppsPostRequestBody()
	body.SetApps(apps)
	err := r.client.
		DeviceAppManagement().
		ManagedAppPolicies().
		ByManagedAppPolicyId(policyID).
		TargetApps().
		Post(ctx, body, nil)
	if err != nil {
		return fmt.Errorf("set policy target apps: %w", err)
	}
	return nil
}

// packageID returns the Android package of an app entry, or "" for other identifiers.
func packageID(app graphmodels.ManagedMobileAppable) string {
	if app == nil {
		return ""
	}
	identifier, ok := app.GetMobileAppIdentifier().(graphmodels.AndroidMobileAppIdentifierable)
	if !ok || identifier.GetPackageId() == nil {
		return ""
	}
	return *identifier.GetPackageId()
}

// writableApp rebuilds an entry without the server-assigned id and version, which
// targetApps derives itself.
func writableApp(app graphmodels.ManagedMobileAppable) graphmodels.ManagedMobileAppable {
	pkg := packageID(app)
	if pkg == "" {
		return app
	}
	identifier := graphmodels.NewAndroidMobileAppIdentifier()
	identifier.SetPackageId(&pkg)
	entry := graphmodels.NewManagedMobileApp()
	entry.SetMobileAppIdentifier(identifier)
	return entry
}

// findByPackage returns the entry targeting the package, matched case-insensitively.
func findByPackage(
	apps []graphmodels.ManagedMobileAppable,
	pkg string,
) graphmodels.ManagedMobileAppable {
	for _, app := range apps {
		if strings.EqualFold(packageID(app), pkg) {
			return app
		}
	}
	return nil
}

// findByID returns the entry with the given id.
func findByID(apps []graphmodels.ManagedMobileAppable, id string) graphmodels.ManagedMobileAppable {
	for _, app := range apps {
		if app != nil && app.GetId() != nil && strings.EqualFold(*app.GetId(), id) {
			return app
		}
	}
	return nil
}

// waitForApps polls the policy's app list until done reports true or ctx expires.
func (r *AndroidManagedMobileAppResource) waitForApps(
	ctx context.Context,
	policyID string,
	done func([]graphmodels.ManagedMobileAppable) bool,
) ([]graphmodels.ManagedMobileAppable, error) {
	for {
		apps, err := r.listPolicyApps(ctx, policyID)
		if err != nil {
			return nil, err
		}
		if done(apps) {
			return apps, nil
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%w: %w", errAppWaitTimeout, ctx.Err())
		case <-time.After(appPollInterval):
		}
	}
}
