package api

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_qdrant_cloud_programmatic_access_ClustersAPIService(t *testing.T) {

	apiKey := "YOUR_API_KEY_HERE"
	accountId := "YOUR_ACCOUNT_ID_HERE"

	apiClient, err := NewClientWithResponses("https://cloud.qdrant.io/public/v1", WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("apikey %s", apiKey))
		return nil
	}))
	require.NoError(t, err)

	t.Run("Test ClustersAPI ListClusters", func(t *testing.T) {

		t.Skip("skip test") // remove to run test

		accountGuid, err := uuid.Parse(accountId)
		require.NoError(t, err)

		ctx := context.Background()
		resp, err := apiClient.ListClustersWithResponse(ctx, accountGuid, &ListClustersParams{})

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.JSON200)
		for _, c := range *resp.JSON200 {
			require.Equal(t, "TestClusterTF", c.Name)
		}

	})
}
