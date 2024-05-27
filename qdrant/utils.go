package qdrant

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	api "terraform-provider-qdrant-cloud/v1/internal/client"
)

// getClient creates a client from the provided interface
// This client already contains the Authorization needed to invoke the API
// Returns: The client to call the backend API, TF Diagnostics
func getClient(m interface{}) (*api.ClientWithResponses, diag.Diagnostics) {
	clientConfig, ok := m.(ClientConfig)
	if !ok {
		return nil, diag.FromErr(fmt.Errorf("error initializing client: provided interface cannot be casted to ClientConfig"))
	}
	optsCallback := api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("apikey %s", clientConfig.ApiKey))
		return nil
	})
	apiClient, err := api.NewClientWithResponses(clientConfig.BaseURL, optsCallback)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error initializing client: %v", err))
	}
	return apiClient, nil
}
