package qdrant

import (
	"context"
	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	api "terraform-provider-qdrant-cloud/v1/internal/client"
	"testing"
)

// MockClient is a mock implementation of the ClientWithResponsesInterface
type MockClient struct {
	mock.Mock
}

func (m *MockClient) ListApiKeysWithResponse(ctx context.Context, accountId types.UUID, reqEditors ...api.RequestEditorFn) (*api.ListApiKeysResponse, error) {
	args := m.Called(ctx, accountId, reqEditors)
	return args.Get(0).(*api.ListApiKeysResponse), args.Error(1)
}

func (m *MockClient) CreateApiKeyWithResponse(ctx context.Context, accountId types.UUID, body api.CreateApiKeyJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.CreateApiKeyResponse, error) {
	args := m.Called(ctx, accountId, body, reqEditors)
	return args.Get(0).(*api.CreateApiKeyResponse), args.Error(1)
}

func (m *MockClient) DeleteApiKeyWithResponse(ctx context.Context, accountId types.UUID, apiKeyId string, reqEditors ...api.RequestEditorFn) (*api.DeleteApiKeyResponse, error) {
	args := m.Called(ctx, accountId, apiKeyId, reqEditors)
	return args.Get(0).(*api.DeleteApiKeyResponse), args.Error(1)
}

func (m *MockClient) ListClustersWithResponse(ctx context.Context, accountId types.UUID, params *api.ListClustersParams, reqEditors ...api.RequestEditorFn) (*api.ListClustersResponse, error) {
	args := m.Called(ctx, accountId, params, reqEditors)
	return args.Get(0).(*api.ListClustersResponse), args.Error(1)
}

func (m *MockClient) CreateClusterWithResponse(ctx context.Context, accountId types.UUID, body api.CreateClusterJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.CreateClusterResponse, error) {
	args := m.Called(ctx, accountId, body, reqEditors)
	return args.Get(0).(*api.CreateClusterResponse), args.Error(1)
}

func (m *MockClient) DeleteClusterWithResponse(ctx context.Context, accountId types.UUID, clusterId types.UUID, params *api.DeleteClusterParams, reqEditors ...api.RequestEditorFn) (*api.DeleteClusterResponse, error) {
	args := m.Called(ctx, accountId, clusterId, params, reqEditors)
	return args.Get(0).(*api.DeleteClusterResponse), args.Error(1)
}

func (m *MockClient) GetClusterWithResponse(ctx context.Context, accountId types.UUID, clusterId types.UUID, reqEditors ...api.RequestEditorFn) (*api.GetClusterResponse, error) {
	args := m.Called(ctx, accountId, clusterId, reqEditors)
	return args.Get(0).(*api.GetClusterResponse), args.Error(1)
}

func (m *MockClient) UpdateClusterWithResponse(ctx context.Context, accountId types.UUID, clusterId types.UUID, body api.UpdateClusterJSONRequestBody, reqEditors ...api.RequestEditorFn) (*api.UpdateClusterResponse, error) {
	args := m.Called(ctx, accountId, clusterId, body, reqEditors)
	return args.Get(0).(*api.UpdateClusterResponse), args.Error(1)
}

func (m *MockClient) GetPackagesWithResponse(ctx context.Context, params *api.GetPackagesParams, reqEditors ...api.RequestEditorFn) (*api.GetPackagesResponse, error) {
	args := m.Called(ctx, params, reqEditors)
	return args.Get(0).(*api.GetPackagesResponse), args.Error(1)
}

func TestExecuteClientAction(t *testing.T) {
	// Create a mock client
	mockClient := new(MockClient)

	// Set up test cases
	testCases := []struct {
		name       string
		actionName string
		parameters map[string]interface{}
		setupMocks func()
		expected   interface{}
	}{
		{
			name:       "ListApiKeys",
			actionName: "ListApiKeys",
			parameters: map[string]interface{}{
				"accountId": "550e8400-e29b-41d4-a716-446655440000",
			},
			setupMocks: func() {
				mockResponse := &api.ListApiKeysResponse{
					Body:         []byte(`[]`),
					HTTPResponse: &http.Response{StatusCode: 200},
					JSON200:      &[]api.GetApiKeyOut{},
				}
				mockClient.On("ListApiKeysWithResponse", mock.Anything, "550e8400-e29b-41d4-a716-446655440000", mock.Anything).Return(mockResponse, nil)
			},
			expected: &api.ListApiKeysResponse{
				Body:         []byte(`[]`),
				HTTPResponse: &http.Response{StatusCode: 200},
				JSON200:      &[]api.GetApiKeyOut{},
			},
		},
		// Add more test cases for each action...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mocks
			tc.setupMocks()

			// Execute the function
			response, err := ExecuteClientAction(context.TODO(), tc.actionName, "test-api-key", tc.parameters)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, response)

			// Verify the mock client was called as expected
			mockClient.AssertExpectations(t)
		})
	}
}
