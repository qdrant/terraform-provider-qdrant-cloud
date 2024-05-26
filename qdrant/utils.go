package qdrant

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/oapi-codegen/runtime/types"
	"log"
	"net/http"
	"terraform-provider-qdrant-cloud/v1/internal/client"
)

func ExecuteClientAction(ctx context.Context, actionName string, apiKey string, parameters map[string]interface{}) (interface{}, error) {
	// Initialize the client
	apiClient, err := api.NewClientWithResponses("https://cloud.qdrant.io/public/v1", api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("apikey %s", apiKey))
		return nil
	}))
	if err != nil {
		log.Printf("Error initializing client: %v", err)
		return nil, err
	}

	var (
		accountId  types.UUID
		clusterId  types.UUID
		apiKeyId   string
		body       interface{}
		reqEditors []api.RequestEditorFn
		response   interface{}
	)

	if v, ok := parameters["accountId"].(types.UUID); ok {
		accountId = v
	}

	if v, ok := parameters["clusterId"].(types.UUID); ok {
		clusterId = v
	}

	if v, ok := parameters["apiKeyId"].(string); ok {
		apiKeyId = v
	}

	if v, ok := parameters["body"]; ok {
		body = v
	}

	if v, ok := parameters["reqEditors"].([]api.RequestEditorFn); ok {
		reqEditors = v
	}

	switch actionName {
	case "ListApiKeys":
		response, err = apiClient.ListApiKeysWithResponse(ctx, accountId, reqEditors...)
	case "CreateApiKey":
		response, err = apiClient.CreateApiKeyWithResponse(ctx, accountId, body.(api.CreateApiKeyJSONRequestBody), reqEditors...)
	case "DeleteApiKey":
		response, err = apiClient.DeleteApiKeyWithResponse(ctx, accountId, apiKeyId, reqEditors...)
	case "ListClusters":
		params := parameters["params"].(*api.ListClustersParams)
		response, err = apiClient.ListClustersWithResponse(ctx, accountId, params, reqEditors...)
	case "CreateCluster":
		response, err = apiClient.CreateClusterWithResponse(ctx, accountId, body.(api.CreateClusterJSONRequestBody), reqEditors...)
	case "DeleteCluster":
		params := parameters["params"].(*api.DeleteClusterParams)
		response, err = apiClient.DeleteClusterWithResponse(ctx, accountId, clusterId, params, reqEditors...)
	case "GetCluster":
		response, err = apiClient.GetClusterWithResponse(ctx, accountId, clusterId, reqEditors...)
	case "UpdateCluster":
		response, err = apiClient.UpdateClusterWithResponse(ctx, accountId, clusterId, body.(api.UpdateClusterJSONRequestBody), reqEditors...)
	case "GetPackages":
		params := parameters["params"].(*api.GetPackagesParams)
		response, err = apiClient.GetPackagesWithResponse(ctx, params, reqEditors...)
	default:
		err = fmt.Errorf("invalid action name: %s", actionName)
	}

	if err != nil {
		log.Printf("Error executing action %s: %v", actionName, err)
		return nil, err
	}

	return response, nil
}

func GetClient(m interface{}) (*api.ClientWithResponses, error, diag.Diagnostics, bool) {
	apiKey := m.(ClientConfig).ApiKey

	opts := api.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("apikey %s", apiKey))
		return nil
	})

	apiClient, err := api.NewClientWithResponses(m.(ClientConfig).BaseURL, opts)
	if err != nil {
		d := diag.FromErr(fmt.Errorf("error initializing client: %v", err))
		if d.HasError() {
			return nil, nil, d, true
		}
	}
	return apiClient, err, nil, false
}
