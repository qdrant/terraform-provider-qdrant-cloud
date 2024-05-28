package qdrant

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

// getClient creates a client from the provided interface
// This client already contains the Authorization needed to invoke the API
// Returns: The client to call the backend API, TF Diagnostics
func getClient(m interface{}) (*qc.ClientWithResponses, diag.Diagnostics) {
	clientConfig, ok := m.(*ClientConfig)
	if !ok {
		return nil, diag.FromErr(fmt.Errorf("error initializing client: provided interface cannot be casted to ClientConfig"))
	}
	optsCallback := qc.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("apikey %s", clientConfig.ApiKey))
		return nil
	})
	apiClient, err := qc.NewClientWithResponses(clientConfig.BaseURL, optsCallback)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error initializing client: %v", err))
	}
	return apiClient, nil
}

// getError returns a human readable error composed from the given HTTP validation error
func getError(error *qc.HTTPValidationError) string {
	if error == nil {
		return "no error"
	}
	details := error.Detail
	if details == nil {
		return "no details"
	}
	var result []string
	for _, ve := range *details {
		result = append(result, fmt.Sprintf("%s:%s", ve.Type, ve.Msg))
	}
	return strings.Join(result, ",")
}
