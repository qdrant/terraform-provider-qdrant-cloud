package qdrant

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
)

type mockClusterServiceClient struct {
	qcCluster.ClusterServiceClient
	responses []mockGetClusterResponse
	callCount int
}

type mockGetClusterResponse struct {
	cluster *qcCluster.Cluster
	err     error
}

func (m *mockClusterServiceClient) getCluster(_ context.Context) (*qcCluster.GetClusterResponse, error) {
	idx := m.callCount
	m.callCount++
	if idx >= len(m.responses) {
		idx = len(m.responses) - 1
	}
	r := m.responses[idx]
	if r.err != nil {
		return nil, r.err
	}
	return &qcCluster.GetClusterResponse{Cluster: r.cluster}, nil
}

func buildRefreshFunc(mock *mockClusterServiceClient) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := mock.getCluster(context.Background())
		if err != nil {
			return nil, "", err
		}

		cluster := resp.GetCluster()
		phase := cluster.GetState().GetPhase()

		if phase == qcCluster.ClusterPhase_CLUSTER_PHASE_FAILED_TO_CREATE {
			return nil, "", errors.New("cluster creation failed (phase=\"CLUSTER_PHASE_FAILED_TO_CREATE\" reason=\"" + cluster.GetState().GetReason() + "\")")
		}

		url := strings.TrimSpace(cluster.GetState().GetEndpoint().GetUrl())
		if url != "" {
			return cluster, clusterWaitReady, nil
		}

		return cluster, clusterWaitPending, nil
	}
}

func TestClusterEndpointRefresh_ReturnsReadyWhenURLPresent(t *testing.T) {
	mock := &mockClusterServiceClient{
		responses: []mockGetClusterResponse{
			{cluster: &qcCluster.Cluster{
				State: &qcCluster.ClusterState{
					Endpoint: &qcCluster.ClusterEndpoint{Url: "https://cluster.example.com"},
				},
			}},
		},
	}

	result, state, err := buildRefreshFunc(mock)()

	require.NoError(t, err)
	assert.Equal(t, clusterWaitReady, state)
	cluster := result.(*qcCluster.Cluster)
	assert.Equal(t, "https://cluster.example.com", cluster.GetState().GetEndpoint().GetUrl())
}

func TestClusterEndpointRefresh_ReturnsPendingWhenURLEmpty(t *testing.T) {
	mock := &mockClusterServiceClient{
		responses: []mockGetClusterResponse{
			{cluster: &qcCluster.Cluster{
				State: &qcCluster.ClusterState{
					Phase:    qcCluster.ClusterPhase_CLUSTER_PHASE_CREATING,
					Endpoint: &qcCluster.ClusterEndpoint{Url: ""},
				},
			}},
		},
	}

	result, state, err := buildRefreshFunc(mock)()

	require.NoError(t, err)
	assert.Equal(t, clusterWaitPending, state)
	assert.NotNil(t, result)
}

func TestClusterEndpointRefresh_FailsFastOnFailedToCreate(t *testing.T) {
	mock := &mockClusterServiceClient{
		responses: []mockGetClusterResponse{
			{cluster: &qcCluster.Cluster{
				State: &qcCluster.ClusterState{
					Phase:  qcCluster.ClusterPhase_CLUSTER_PHASE_FAILED_TO_CREATE,
					Reason: "insufficient resources",
				},
			}},
		},
	}

	result, _, err := buildRefreshFunc(mock)()

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cluster creation failed")
	assert.Contains(t, err.Error(), "CLUSTER_PHASE_FAILED_TO_CREATE")
	assert.Contains(t, err.Error(), "insufficient resources")
}

func TestClusterEndpointRefresh_PropagatesAPIErrors(t *testing.T) {
	mock := &mockClusterServiceClient{
		responses: []mockGetClusterResponse{
			{err: errors.New("backend unavailable")},
		},
	}

	result, _, err := buildRefreshFunc(mock)()

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "backend unavailable")
}

func TestClusterEndpointRefresh_EventualSuccess(t *testing.T) {
	mock := &mockClusterServiceClient{
		responses: []mockGetClusterResponse{
			{cluster: &qcCluster.Cluster{
				State: &qcCluster.ClusterState{
					Phase:    qcCluster.ClusterPhase_CLUSTER_PHASE_CREATING,
					Endpoint: &qcCluster.ClusterEndpoint{Url: ""},
				},
			}},
			{cluster: &qcCluster.Cluster{
				State: &qcCluster.ClusterState{
					Phase:    qcCluster.ClusterPhase_CLUSTER_PHASE_CREATING,
					Endpoint: &qcCluster.ClusterEndpoint{Url: ""},
				},
			}},
			{cluster: &qcCluster.Cluster{
				State: &qcCluster.ClusterState{
					Phase:    qcCluster.ClusterPhase_CLUSTER_PHASE_HEALTHY,
					Endpoint: &qcCluster.ClusterEndpoint{Url: "https://ready.example.com"},
				},
			}},
		},
	}

	refreshFunc := buildRefreshFunc(mock)

	// First two calls return pending
	_, state, err := refreshFunc()
	require.NoError(t, err)
	assert.Equal(t, clusterWaitPending, state)

	_, state, err = refreshFunc()
	require.NoError(t, err)
	assert.Equal(t, clusterWaitPending, state)

	// Third call returns ready
	result, state, err := refreshFunc()
	require.NoError(t, err)
	assert.Equal(t, clusterWaitReady, state)
	cluster := result.(*qcCluster.Cluster)
	assert.Equal(t, "https://ready.example.com", cluster.GetState().GetEndpoint().GetUrl())
	assert.Equal(t, 3, mock.callCount)
}
