package customrequests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

var errEmptyJSONResponse = errors.New("graph returned an empty JSON response")

// JSONRequest uses the configured SDK adapter for endpoints without generated models.
// A nil response expects no content. Path parameters are expanded by Kiota.
func JSONRequest(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	template string,
	parameters map[string]string,
	body any,
	response any,
) error {
	params := map[string]string{"baseurl": adapter.GetBaseUrl()}
	for k, v := range parameters {
		params[k] = v
	}
	info := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(
		method,
		"{+baseurl}"+template,
		params,
	)
	info.Headers.TryAdd("Accept", "application/json")
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal Graph JSON request: %w", err)
		}
		info.Content = raw
		info.Headers.TryAdd("Content-Type", "application/json")
	}
	mappings := abstractions.ErrorMappings{
		"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
	}
	if response == nil {
		if err := adapter.SendNoContent(ctx, info, mappings); err != nil {
			return fmt.Errorf("send Graph request: %w", err)
		}
		return nil
	}
	result, err := adapter.SendPrimitive(ctx, info, "[]byte", mappings)
	if err != nil {
		return fmt.Errorf("send Graph JSON request: %w", err)
	}
	raw, ok := result.([]byte)
	if !ok || len(raw) == 0 || string(raw) == "null" {
		return errEmptyJSONResponse
	}
	if err := json.Unmarshal(raw, response); err != nil {
		return fmt.Errorf("decode Graph JSON response: %w", err)
	}
	return nil
}
