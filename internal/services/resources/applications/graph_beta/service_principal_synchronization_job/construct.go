package graphBetaServicePrincipalSynchronizationJob

import (
	"errors"
	"regexp"
	"strings"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var (
	errImportID    = errors.New("expected service_principal_id/job_id with a UUID parent and a nonempty opaque job ID without slashes")
	errJobResponse = errors.New("incomplete synchronization job response from Microsoft Graph")
)

func parseResourceID(id string) (string, string, error) {
	parent, job, ok := strings.Cut(id, "/")
	if !ok || !regexp.MustCompile(constants.GuidRegex).MatchString(parent) || strings.TrimSpace(job) == "" || strings.Contains(job, "/") {
		return "", "", errImportID
	}
	return parent, job, nil
}

func constructResource(data *ServicePrincipalSynchronizationJobResourceModel) graphmodels.SynchronizationJobable {
	body := graphmodels.NewSynchronizationJob()
	convert.FrameworkToGraphString(data.TemplateID, body.SetTemplateId)
	return body
}
